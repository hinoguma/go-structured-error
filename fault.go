package fault

import (
	"errors"
	"fmt"
	"time"
)

func New(message string) *FaultError {
	err := NewRawFaultError(errors.New(message))
	// set stack trace starting from caller of NewFaultError
	_ = err.SetStackTraceWithSkipMaxDepth(2, MaxStackTraceDepth)
	return err
}

func NewRawFaultError(err error) *FaultError {
	return &FaultError{
		errorType:  ErrorTypeNone,
		err:        err,
		stacktrace: make(StackTrace, 0),

		when:      nil,
		requestId: "",
		tags:      NewTags(),
		subErrors: make([]error, 0),
	}
}

const (
	ErrorTypeNone ErrorType = ""
)

type ErrorType string

func (value ErrorType) String() string {
	return string(value)
}

func (value ErrorType) StringWithDefaultNone() string {
	if value == "" {
		return "none"
	}
	return string(value)
}

type FaultError struct {
	// required
	errorType  ErrorType
	err        error
	stacktrace StackTrace

	// optional
	when      *time.Time
	requestId string
	tags      Tags
	subErrors []error
}

func (e *FaultError) Error() string {
	m := NoErrStr
	if e.err != nil {
		m = e.err.Error()
	}
	return fmt.Sprintf("[Type: %s] %s", e.errorType.StringWithDefaultNone(), m)
}

func (e FaultError) Unwrap() error {
	return e.err
}

func (e *FaultError) Is(target error) bool {
	if target == nil {
		return false
	}
	targetFe, ok := target.(Fault)
	if !ok {
		return false
	}
	return e.Type() == targetFe.Type() && errors.Is(e.Unwrap(), targetFe.Unwrap())
}

func (e FaultError) Type() ErrorType {
	return e.errorType
}

func (e FaultError) StackTrace() StackTrace {
	if e.stacktrace == nil {
		return make([]StackTraceItem, 0)
	}
	return e.stacktrace
}

func (e FaultError) When() *time.Time {
	return e.when
}

func (e FaultError) RequestID() string {
	return e.requestId
}

func (e *FaultError) SetErr(err error) Fault {
	e.err = err
	return e
}

func (e *FaultError) SetType(errorType ErrorType) Fault {
	e.errorType = errorType
	return e
}

func (e *FaultError) SetWhen(t time.Time) Fault {
	e.when = &t
	return e
}

func (e *FaultError) SetRequestID(requestID string) Fault {
	e.requestId = requestID
	return e
}

// WithStackTrace sets stack trace starting from caller of WithStackTrace
func (e *FaultError) WithStackTrace() Fault {
	return e.SetStackTraceWithSkipMaxDepth(2, MaxStackTraceDepth) // skip 2 to start at caller of WithStackTrace
}

func (e *FaultError) SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) Fault {
	e.stacktrace = NewStackTrace(skip, maxDepth)
	return e
}

func (e *FaultError) AddTagString(key string, value string) Fault {
	return e.AddTagSafe(key, StringTagValue(value))
}

func (e *FaultError) AddTagInt(key string, value int) Fault {
	return e.AddTagSafe(key, IntTagValue(value))
}

func (e *FaultError) AddTagBool(key string, value bool) Fault {
	return e.AddTagSafe(key, BoolTagValue(value))
}

func (e *FaultError) AddTagFloat(key string, value float64) Fault {
	return e.AddTagSafe(key, FloatTagValue(value))
}

func (e *FaultError) AddTagSafe(key string, value TagValue) Fault {
	e.tags.SetValueSafe(key, value)
	return e
}

// It`s planned to be implemented later
//func (e *FaultError) AddTag(key string, value any) bool {
//
//	e.tags.SetValueSafe(key, InterfaceTagValue(value))
//	return true
//}

func (e *FaultError) DeleteTag(key string) Fault {
	e.tags.Delete(key)
	return e
}

func (e *FaultError) AddSubError(errs ...error) Fault {
	if len(errs) == 0 {
		return e
	}
	filtered := make([]error, 0)
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	if len(filtered) == 0 {
		return e
	}
	if e.subErrors == nil {
		e.subErrors = make([]error, 0)
	}
	e.subErrors = append(e.subErrors, filtered...)
	return e
}

func (e *FaultError) JsonString() string {
	return e.JsonFormatter().Format()
}

func (e *FaultError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s", e.VerboseFormatter().Format())
			return
		}
		_, _ = fmt.Fprintf(f, "%s", e.Error())
	case 's':
		_, _ = fmt.Fprintf(f, "%s", e.Error())
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", e.Error())
	}
}

func (e *FaultError) JsonFormatter() ErrorFormatter {
	return JsonFormatter{
		errorType:  e.errorType,
		err:        e.err,
		stacktrace: e.stacktrace,
		when:       e.when,
		requestId:  e.requestId,
		tags:       e.tags,
		subErrors:  e.subErrors,
	}
}

func (e *FaultError) VerboseFormatter() ErrorFormatter {
	return VerboseFormatter{
		title:      "main_error",
		errorType:  e.errorType,
		err:        e.err,
		stacktrace: e.stacktrace,
		when:       e.when,
		requestId:  e.requestId,
		tags:       e.tags,
		subErrors:  e.subErrors,
	}
}

// IsTYpe() checks whether the given error or any of its wrapped errors is of the specified ErrorType.
// errors.Is() checks for error equality, but this function checks for error type.
func IsType(err error, t ErrorType) bool {
	if err == nil {
		return false
	}
	fe, ok := err.(Typer)
	if ok && fe.Type() == t {
		return true
	}

	switch x := err.(type) {
	case interface{ Unwrap() error }:
		return IsType(x.Unwrap(), t)
	case interface{ Unwrap() []error }:
		for _, subErr := range x.Unwrap() {
			if IsType(subErr, t) {
				return true
			}
		}
	}
	return false
}

/*********************
	Interfaces
 *********************/

type Fault interface {
	error
	Unwrap() error
	Type() ErrorType
	When() *time.Time
	RequestID() string
	StackTrace() StackTrace

	SetErr(err error) Fault
	SetType(errorType ErrorType) Fault
	SetWhen(t time.Time) Fault
	SetRequestID(requestID string) Fault
	WithStackTrace() Fault // auto set stack trace
	SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) Fault
	AddTagSafe(key string, value TagValue) Fault
	DeleteTag(key string) Fault
	AddSubError(errs ...error) Fault

	JsonString() string
}

type Typer interface {
	Type() ErrorType
}

type StackTracer interface {
	StackTrace() StackTrace
}

type JsonStringer interface {
	JsonString() string
}
