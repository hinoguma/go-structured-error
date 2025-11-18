package fault

import (
	"errors"
	"fmt"
	"time"
)

type Fault interface {
	error
	Unwrap() error
	Type() FaultType
	When() *time.Time
	RequestID() string

	SetErr(err error) Fault
	SetWhen(t time.Time) Fault
	SetRequestID(requestID string) Fault
	WithStackTrace() Fault // auto set stack trace
	SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) Fault
	AddTagSafe(key string, value TagValue) Fault
	DeleteTag(key string) Fault
}

type FaultType string

func (value FaultType) String() string {
	return string(value)
}

const (
	FaultTypeNone FaultType = ""
)

func NewRawFaultError(err error) *FaultError {
	return &FaultError{
		faultType:  FaultTypeNone,
		err:        err,
		stacktrace: make(StackTrace, 0),

		when:      nil,
		requestId: "",
		tags:      NewTags(),
		subErrors: make([]error, 0),
	}
}

func New(message string) *FaultError {
	err := NewRawFaultError(errors.New(message))
	// set stack trace starting from caller of NewFaultError
	err.SetStackTraceWithSkipMaxDepth(2, GetMaxDepthStackTrace())
	return err
}

type FaultError struct {
	// required
	faultType  FaultType
	err        error
	stacktrace StackTrace

	// optional
	when      *time.Time
	requestId string
	tags      Tags
	subErrors []error
}

func (e *FaultError) Error() string {
	t := e.faultType.String()
	if t == "" {
		t = "none"
	}
	m := "<no error>"
	if e.err != nil {
		m = e.err.Error()
	}
	return fmt.Sprintf("[Type: %s] %s", t, m)
}

func (e FaultError) Unwrap() error {
	return e.err
}

func (e FaultError) Is(target error) bool {
	if target == nil {
		return false
	}
	switch x := target.(type) {
	case interface {
		Type() FaultType
		Unwrap() error
	}:
		if x.Type() == e.Type() {
			return errors.Is(e.Unwrap(), x.Unwrap())
		}
	}
	return false
}

func (e FaultError) Type() FaultType {
	return e.faultType
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
	return e.SetStackTraceWithSkipMaxDepth(2, GetMaxDepthStackTrace()) // skip 4 to start at caller of WithStackTrace
}

func (e *FaultError) SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) Fault {
	e.stacktrace = NewStackTrace(skip, maxDepth)
	return e
}

func (e *FaultError) AddTagString(key string, value string) Fault {
	e.AddTagSafe(key, StringTagValue(value))
	return e
}

func (e *FaultError) AddTagInt(key string, value int) Fault {
	e.AddTagSafe(key, IntTagValue(value))
	return e
}

func (e *FaultError) AddTagBool(key string, value bool) Fault {
	e.AddTagSafe(key, BoolTagValue(value))
	return e
}

func (e *FaultError) AddTagFloat(key string, value float64) Fault {
	e.AddTagSafe(key, FloatTagValue(value))
	return e
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
