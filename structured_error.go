package serrors

import (
	"errors"
	"fmt"
	"time"
)

func NewRawStructuredError(err error) *StructuredError {
	return &StructuredError{
		errorType:  ErrorTypeNone,
		err:        err,
		stacktrace: make(StackTrace, 0),

		when:      nil,
		requestId: "",
		tags:      NewTags(),
		subErrors: make([]error, 0),
	}
}

func NewWithSkipAndDepth(err error, skip int, maxDepth int) *StructuredError {
	if skip < 0 {
		skip = 0
	}
	fe := NewRawStructuredError(err)
	_ = fe.SetStackTraceWithSkipMaxDepth(skip+1, maxDepth) // skip +1 to start at caller of NewWithSkipAndDepth
	return fe
}

const NoErrStr string = "<no error>"

const ErrorTypeNone ErrorType = ""

// ErrorType represents the type/category of the error
// You can define your own error types as needed
// IsType() method uses this type for comparison
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

// StructuredError has structured information about an error
// It can be JsonString for logging
type StructuredError struct {
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

func (e *StructuredError) Error() string {
	m := NoErrStr
	if e.err != nil {
		m = e.err.Error()
	}
	return fmt.Sprintf("[Type: %s] %s", e.errorType.StringWithDefaultNone(), m)
}

func (e StructuredError) Unwrap() error {
	return e.err
}

func (e *StructuredError) Is(target error) bool {
	if target == nil {
		return false
	}
	targetFe, ok := target.(SError)
	if !ok {
		return false
	}
	return e.Type() == targetFe.Type() && errors.Is(e.Unwrap(), targetFe.Unwrap())
}

func (e StructuredError) Type() ErrorType {
	return e.errorType
}

func (e StructuredError) StackTrace() StackTrace {
	if e.stacktrace == nil {
		return make([]StackTraceItem, 0)
	}
	return e.stacktrace
}

func (e StructuredError) When() *time.Time {
	return e.when
}

func (e StructuredError) RequestID() string {
	return e.requestId
}

func (e *StructuredError) SetErr(err error) SError {
	e.err = err
	return e
}

func (e *StructuredError) SetType(errorType ErrorType) SError {
	e.errorType = errorType
	return e
}

func (e *StructuredError) SetWhen(t time.Time) SError {
	e.when = &t
	return e
}

func (e *StructuredError) SetRequestID(requestID string) SError {
	e.requestId = requestID
	return e
}

// WithStackTrace sets stack trace starting from caller of WithStackTrace
func (e *StructuredError) WithStackTrace() SError {
	return e.SetStackTraceWithSkipMaxDepth(2, MaxStackTraceDepth) // skip 2 to start at caller of WithStackTrace
}

func (e *StructuredError) SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) SError {
	e.stacktrace = NewStackTrace(skip, maxDepth)
	return e
}

func (e *StructuredError) AddTagString(key string, value string) SError {
	return e.AddTagSafe(key, StringTagValue(value))
}

func (e *StructuredError) AddTagInt(key string, value int) SError {
	return e.AddTagSafe(key, IntTagValue(value))
}

func (e *StructuredError) AddTagBool(key string, value bool) SError {
	return e.AddTagSafe(key, BoolTagValue(value))
}

func (e *StructuredError) AddTagFloat(key string, value float64) SError {
	return e.AddTagSafe(key, FloatTagValue(value))
}

func (e *StructuredError) AddTagSafe(key string, value TagValue) SError {
	e.tags.SetValueSafe(key, value)
	return e
}

// It`s planned to be implemented later
//func (e *StructuredError) AddTag(key string, value any) bool {
//
//	e.tags.SetValueSafe(key, InterfaceTagValue(value))
//	return true
//}

func (e *StructuredError) DeleteTag(key string) SError {
	e.tags.Delete(key)
	return e
}

func (e *StructuredError) AddSubError(errs ...error) SError {
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

func (e *StructuredError) JsonString() string {
	return e.JsonPrinter().Print()
}

func (e *StructuredError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s", e.VerbosePrinter().Print())
			return
		}
		_, _ = fmt.Fprintf(f, "%s", e.Error())
	case 's':
		_, _ = fmt.Fprintf(f, "%s", e.Error())
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", e.Error())
	}
}

func (e *StructuredError) JsonPrinter() JsonPrinter {
	return ErrorJsonPrinter{
		errorType:  e.errorType,
		err:        e.err,
		stacktrace: e.stacktrace,
		when:       e.when,
		requestId:  e.requestId,
		tags:       e.tags,
		subErrors:  e.subErrors,
	}
}

func (e *StructuredError) VerbosePrinter() VerbosePrinter {
	return ErrorVerbosePrinter{
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

/*********************
	Interfaces
 *********************/

type SError interface {
	error
	Unwrap() error
	Type() ErrorType
	When() *time.Time
	RequestID() string
	StackTrace() StackTrace

	SetErr(err error) SError
	SetType(errorType ErrorType) SError
	SetWhen(t time.Time) SError
	SetRequestID(requestID string) SError
	WithStackTrace() SError // auto set stack trace
	SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) SError
	AddTagSafe(key string, value TagValue) SError
	DeleteTag(key string) SError
	AddSubError(errs ...error) SError
}

// IsType() use this interface
// if your custom error implements it, It`s comparable in IsType()
type HasType interface {
	Type() ErrorType
}

type JsonStringer interface {
	JsonString() string
}

type HasJsonPrinter interface {
	JsonPrinter() JsonPrinter
}
