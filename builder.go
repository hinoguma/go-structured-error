package serrors

import "time"

func Builder(err error) *StructuredErrorBuilder {
	if err == nil {
		return &StructuredErrorBuilder{err: nil}
	}
	return &StructuredErrorBuilder{err: ToStructured(err)}
}

type StructuredErrorBuilder struct {
	err SError
}

func (w *StructuredErrorBuilder) Build() error {
	return w.err
}

// set stack trace starting from caller of StackTrace method
func (w *StructuredErrorBuilder) StackTrace() *StructuredErrorBuilder {
	return w.StackTraceBuilderSkipDepth(2, MaxStackTraceDepth)
}

// if skip is negative, it will be treated as 0
// so starting from the caller of StackTraceBuilderSkipDepth method
func (w *StructuredErrorBuilder) StackTraceBuilderSkipDepth(skip, depth int) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	if skip < 0 {
		skip = 0
	}
	_ = w.err.SetStackTraceWithSkipMaxDepth(skip+1, depth)
	return w
}

func (w *StructuredErrorBuilder) Type(t ErrorType) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.SetType(t)
	return w
}

func (w *StructuredErrorBuilder) RequestID(id string) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.SetRequestID(id)
	return w
}

func (w *StructuredErrorBuilder) When(t time.Time) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.SetWhen(t)
	return w
}

func (w *StructuredErrorBuilder) AddTagSafe(key string, value TagValue) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, value)
	return w
}

func (w *StructuredErrorBuilder) AddTagString(key string, value string) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, StringTagValue(value))
	return w
}

func (w *StructuredErrorBuilder) AddTagInt(key string, value int) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, IntTagValue(value))
	return w
}

func (w *StructuredErrorBuilder) AddTagFloat(key string, value float64) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, FloatTagValue(value))
	return w
}

func (w *StructuredErrorBuilder) AddTagBool(key string, value bool) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.AddTagSafe(key, BoolTagValue(value))
	return w
}

func (w *StructuredErrorBuilder) DeleteTag(key string) *StructuredErrorBuilder {
	if w.err == nil {
		return w
	}
	_ = w.err.DeleteTag(key)
	return w
}
