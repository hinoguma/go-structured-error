package errors

import (
	"github.com/hinoguma/go-fault"
	"time"
)

func With(err error) *WithWrapper {
	return &WithWrapper{err: err}
}

type WithWrapper struct {
	err error
}

func (w *WithWrapper) Err() error {
	return w.err
}

func (w *WithWrapper) convertToFault() fault.Fault {
	if w.err == nil {
		return nil
	}
	err, ok := w.err.(fault.Fault)
	if !ok {
		err = fault.NewRawFaultError(w.err)
	}
	return err
}

// set stack trace starting from caller of StackTrace method
func (w *WithWrapper) StackTrace() *WithWrapper {
	return w.StackTraceWithSkipDepth(2, fault.GetMaxDepthStackTrace())
}

// if skip is negative, it will be treated as 0
// so starting from the caller of StackTraceWithSkipDepth method
func (w *WithWrapper) StackTraceWithSkipDepth(skip, depth int) *WithWrapper {
	if w.err == nil {
		return w
	}
	err := w.convertToFault()
	if skip < 0 {
		skip = 0
	}
	w.err = err.SetStackTraceWithSkipMaxDepth(skip+1, depth)
	return w
}

func (w *WithWrapper) Type(t fault.FaultType) *WithWrapper {
	if w.err == nil {
		return w
	}
	err := w.convertToFault()
	w.err = err.SetType(t)
	return w
}

func (w *WithWrapper) RequestID(id string) *WithWrapper {
	if w.err == nil {
		return w
	}
	err := w.convertToFault()
	w.err = err.SetRequestID(id)
	return w
}

func (w *WithWrapper) When(t time.Time) *WithWrapper {
	if w.err == nil {
		return w
	}
	err := w.convertToFault()
	w.err = err.SetWhen(t)
	return w
}

func (w *WithWrapper) AddTagSafe(key string, value fault.TagValue) *WithWrapper {
	if w.err == nil {
		return w
	}
	err := w.convertToFault()
	w.err = err.AddTagSafe(key, value)
	return w
}

func (w *WithWrapper) DeleteTag(key string) *WithWrapper {
	if w.err == nil {
		return w
	}
	err := w.convertToFault()
	w.err = err.DeleteTag(key)
	return w
}
