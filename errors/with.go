package errors

import (
	"github.com/hinoguma/go-fault"
	"time"
)

func With(err error) *WithWrapper {
	if err == nil {
		return &WithWrapper{err: nil}
	}
	return &WithWrapper{err: ToFault(err)}
}

type WithWrapper struct {
	err fault.Fault
}

func (w *WithWrapper) Err() error {
	return w.err
}

// set stack trace starting from caller of StackTrace method
func (w *WithWrapper) StackTrace() *WithWrapper {
	return w.StackTraceWithSkipDepth(2, fault.MaxStackTraceDepth)
}

// if skip is negative, it will be treated as 0
// so starting from the caller of StackTraceWithSkipDepth method
func (w *WithWrapper) StackTraceWithSkipDepth(skip, depth int) *WithWrapper {
	if w.err == nil {
		return w
	}
	if skip < 0 {
		skip = 0
	}
	_ = w.err.SetStackTraceWithSkipMaxDepth(skip+1, depth)
	return w
}

func (w *WithWrapper) Type(t fault.ErrorType) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.SetType(t)
	return w
}

func (w *WithWrapper) RequestID(id string) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.SetRequestID(id)
	return w
}

func (w *WithWrapper) When(t time.Time) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.SetWhen(t)
	return w
}

func (w *WithWrapper) AddTagSafe(key string, value fault.TagValue) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, value)
	return w
}

func (w *WithWrapper) AddTagString(key string, value string) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, fault.StringTagValue(value))
	return w
}

func (w *WithWrapper) AddTagInt(key string, value int) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, fault.IntTagValue(value))
	return w
}

func (w *WithWrapper) AddTagFloat(key string, value float64) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, fault.FloatTagValue(value))
	return w
}

func (w *WithWrapper) AddTagBool(key string, value bool) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, fault.BoolTagValue(value))
	return w
}

func (w *WithWrapper) DeleteTag(key string) *WithWrapper {
	if w.err == nil {
		return w
	}
	_ = w.err.DeleteTag(key)
	return w
}
