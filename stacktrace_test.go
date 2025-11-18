package fault

import (
	"runtime"
	"strings"
	"testing"
)

func assertEqualsStackTraceItem(t *testing.T, got, expected StackTraceItem) {
	// only check traces from this package
	// runtime and file system depends on environment
	if !strings.HasPrefix(got.Function, "app/exception.traceLevel") {
		return
	}
	//if got.File != expected.File {
	//	t.Errorf("expected file %v, got %v", expected.File, got.File)
	//}
	if got.Line != expected.Line {
		t.Errorf("expected line %v, got %v expected:%v got:%v", expected.Line, got.Line, expected, got)
	}
	if got.Function != expected.Function {
		t.Errorf("expected function %v, got %v expected:%v got:%v", expected.Function, got.Function, expected, got)
	}
}

func assertEqualsStackTrace(t *testing.T, got, expected StackTrace) {
	if len(got) != len(expected) {
		t.Errorf("expected stack trace length %v, got %v expected:%v got :%v", len(expected), len(got), expected, got)
		return
	}
	for i := range got {
		assertEqualsStackTraceItem(t, got[i], expected[i])
	}
}

const (
	LineTraceLevel1 = 48
	LineTraceLevel2 = 55
	LineTraceLevel3 = 62
	LineTraceLevel4 = 69
	LineTraceLevel5 = 73
)

func traceLevel1(skip int, maxDepth int, level int) StackTrace {
	if level <= 1 {
		return NewStackTrace(skip, maxDepth)
	}
	return traceLevel2(skip, maxDepth, level)
}

func traceLevel2(skip int, maxDepth int, level int) StackTrace {
	if level <= 2 {
		return NewStackTrace(skip, maxDepth)
	}
	return traceLevel3(skip, maxDepth, level)
}

func traceLevel3(skip int, maxDepth int, level int) StackTrace {
	if level <= 3 {
		return NewStackTrace(skip, maxDepth)
	}
	return traceLevel4(skip, maxDepth, level)
}

func traceLevel4(skip int, maxDepth int, level int) StackTrace {
	if level <= 4 {
		return NewStackTrace(skip, maxDepth)
	}
	return traceLevel5(skip, maxDepth)
}

func traceLevel5(skip int, maxDepth int) StackTrace {
	return NewStackTrace(skip, maxDepth)
}

func TestNewStackTrace(t *testing.T) {

	testCases := []struct {
		label    string
		caller   func() StackTrace
		expected StackTrace
	}{
		{
			label: "max depth is not beyond stack",
			caller: func() StackTrace {
				return traceLevel1(0, 3, 5)
			},
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel5,
					Function: "app/exception.traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "app/exception.traceLevel4",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "app/exception.traceLevel3",
				},
			},
		},
		{
			label: "max depth is beyond stack",
			caller: func() StackTrace {
				return traceLevel1(0, 100, 5)
			},
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel5,
					Function: "app/exception.traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "app/exception.traceLevel4",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "app/exception.traceLevel3",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel2,
					Function: "app/exception.traceLevel2",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel1,
					Function: "app/exception.traceLevel1",
				},

				// The exact file and line number will vary, so we use placeholders
				{
					File:     "ignored",
					Line:     0,
					Function: "",
				},
				{
					File:     "ignored",
					Line:     0,
					Function: "",
				},
				{
					File:     "ignored",
					Line:     0,
					Function: "",
				},
				{
					File:     "ignored",
					Line:     0,
					Function: "",
				},
			},
		},
		{
			label: "skip -1",
			caller: func() StackTrace {
				return traceLevel1(-1, 3, 5)
			},
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel5,
					Function: "app/exception.traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "app/exception.traceLevel4",
				},

				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "app/exception.traceLevel3",
				},
			},
		},
		{
			label: "skip 2",
			caller: func() StackTrace {
				return traceLevel1(2, 3, 5)
			},
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "app/exception.traceLevel3",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel2,
					Function: "app/exception.traceLevel2",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel1,
					Function: "app/exception.traceLevel1",
				},
			},
		},
		{
			label: "skip -10",
			caller: func() StackTrace {
				return traceLevel1(-10, 3, 5)
			},
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel5,
					Function: "app/exception.traceLevel5",
				},
			},
		},
		{
			label: "max depth is 0",
			caller: func() StackTrace {
				return traceLevel1(0, 0, 5)
			},
			expected: make(StackTrace, 0),
		},
		{
			label: "max depth is -1",
			caller: func() StackTrace {
				return traceLevel1(0, -1, 5)
			},
			expected: make(StackTrace, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			st := tc.caller()
			assertEqualsStackTrace(t, st, tc.expected)
		})
	}

}

func TestNewStackTraceItem(t *testing.T) {
	testCases := []struct {
		label    string
		fr       runtime.Frame
		expected StackTraceItem
	}{
		{
			label: "basic stack trace item",
			fr: runtime.Frame{
				File:     "exception/stacktrace_test.go",
				Line:     75,
				Function: "exception.TestNewStackTraceItem",
			},
			expected: StackTraceItem{
				File:     "exception/stacktrace_test.go",
				Line:     75,
				Function: "exception.TestNewStackTraceItem",
			},
		},
		{
			label: "empty frame",
			fr:    runtime.Frame{},
			expected: StackTraceItem{
				File:     "",
				Line:     0,
				Function: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			item := NewStackTraceItem(tc.fr)
			assertEqualsStackTraceItem(t, item, tc.expected)
		})
	}
}
