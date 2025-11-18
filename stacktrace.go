package fault

import "runtime"

type StackTrace []StackTraceItem

func NewStackTrace(skip int, maxDepth int) StackTrace {
	if skip+2 < 2 {
		skip = 2 // skip for NewStackTrace and runtime.Callers
	}
	if maxDepth <= 0 {
		return make(StackTrace, 0)
	}
	var trace StackTrace
	pc := make([]uintptr, maxDepth)
	cnt := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:cnt])
	for {
		frame, more := frames.Next()
		item := NewStackTraceItem(frame)
		trace = append(trace, item)
		if !more {
			break
		}
	}
	return trace
}

type StackTraceItem struct {
	File     string
	Line     int
	Function string
}

func NewStackTraceItem(f runtime.Frame) StackTraceItem {
	return StackTraceItem{
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	}
}
