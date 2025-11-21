package fault

import (
	"fmt"
	"runtime"
	"strconv"
)

type StackTrace []StackTraceItem

func (st StackTrace) JsonValueString() string {
	jv := "["
	for i, item := range st {
		if i > 0 {
			jv += ","
		}
		jv += "{"
		jv += "\"file\":\"" + item.File + "\","
		jv += "\"line\":" + strconv.Itoa(item.Line) + ","
		jv += "\"function\":\"" + item.Function + "\""
		jv += "}"
	}
	jv += "]"
	return jv
}

// NewStackTrace captures the current stack trace
// NewStackTrace() starts capturing from the caller of NewStackTrace()
// skip: number of stack frames to skip from the top
// maxDepth: maximum number of stack frames to capture
// if maxDepth <= 0, an empty StackTrace is returned
func NewStackTrace(skip int, maxDepth int) StackTrace {
	if skip < 0 {
		skip = 0
	}
	skip += 2 // skip Callers and NewStackTrace
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

func (item StackTraceItem) String() string {
	return fmt.Sprintf("%s() %s:%d", item.Function, item.File, item.Line)
}

func NewStackTraceItem(f runtime.Frame) StackTraceItem {
	return StackTraceItem{
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	}
}
