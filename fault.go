package fault

import (
	"errors"
	"time"
)

type Fault interface {
	error
	Unwrap() error
	Type() FaultType
	When() *time.Time
	RequestID() string

	SetErr(err error) Fault
	SetType(faultType FaultType) Fault
	SetWhen(t time.Time) Fault
	SetRequestID(requestID string) Fault
	WithStackTrace() Fault // auto set stack trace
	SetStackTraceWithSkipMaxDepth(skip int, maxDepth int)
	AddTagSafe(key string, value TagValue) Fault
	DeleteTag(key string) Fault
}

type FaultType string

const (
	FaultTypeUtil FaultType = ""
)

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
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e FaultError) Unwrap() error {
	return e.err
}

func (e FaultError) Is(target error) bool {
	if target == nil {
		return false
	}
	switch x := target.(type) {
	case Fault:
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

func (e *FaultError) SetErr(err error) *FaultError {
	e.err = err
	return e
}

func (e *FaultError) SetWhen(t time.Time) *FaultError {
	e.when = &t
	return e
}

func (e *FaultError) SetRequestID(requestID string) *FaultError {
	e.requestId = requestID
	return e
}

func (e *FaultError) WithStackTrace() *FaultError {
	return e.SetStackTraceWithSkipMaxDepth(4, GetMaxDepthStackTrace()) // skip 4 to start at caller of WithStackTrace
}

func (e *FaultError) SetStackTraceWithSkipMaxDepth(skip int, maxDepth int) *FaultError {
	e.stacktrace = NewStackTrace(skip, maxDepth)
	return e
}

func (e *FaultError) AddTagString(key string, value string) *FaultError {
	e.AddTagSafe(key, StringTagValue(value))
	return e
}

func (e *FaultError) AddTagInt(key string, value int) *FaultError {
	e.AddTagSafe(key, IntTagValue(value))
	return e
}

func (e *FaultError) AddTagBool(key string, value bool) *FaultError {
	e.AddTagSafe(key, BoolTagValue(value))
	return e
}

func (e *FaultError) AddTagFloat(key string, value float64) *FaultError {
	e.AddTagSafe(key, FloatTagValue(value))
	return e
}

func (e *FaultError) AddTagSafe(key string, value TagValue) {
	e.tags.SetValueSafe(key, value)
}

// It`s planned to be implemented later
//func (e *FaultError) AddTag(key string, value any) bool {
//
//	e.tags.SetValueSafe(key, InterfaceTagValue(value))
//	return true
//}

func (e *FaultError) DeleteTag(key string) *FaultError {
	e.tags.Delete(key)
	return e
}

func (e *FaultError) AddSubError(errs ...error) *FaultError {
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
