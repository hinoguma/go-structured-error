package go_fault

import (
	"runtime"
	"testing"
)

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
					Function: "github.com/hinoguma/go-traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "github.com/hinoguma/go-traceLevel4",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "github.com/hinoguma/go-traceLevel3",
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
					Function: "github.com/hinoguma/go-traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "github.com/hinoguma/go-traceLevel4",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "github.com/hinoguma/go-traceLevel3",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel2,
					Function: "github.com/hinoguma/go-traceLevel2",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel1,
					Function: "github.com/hinoguma/go-traceLevel1",
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
					Function: "github.com/hinoguma/go-traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "github.com/hinoguma/go-traceLevel4",
				},

				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "github.com/hinoguma/go-traceLevel3",
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
					Function: "github.com/hinoguma/go-traceLevel3",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel2,
					Function: "github.com/hinoguma/go-traceLevel2",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel1,
					Function: "github.com/hinoguma/go-traceLevel1",
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
					Function: "github.com/hinoguma/go-traceLevel5",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "github.com/hinoguma/go-traceLevel4",
				},
				{
					File:     "ignored",
					Line:     LineTraceLevel3,
					Function: "github.com/hinoguma/go-traceLevel3",
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
			assertEqualsStackTrace(t, st, tc.expected, "github.com/hinoguma/go-traceLevel")
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
				Function: "github.com/hinoguma/go-TestNewStackTraceItem",
			},
			expected: StackTraceItem{
				File:     "exception/stacktrace_test.go",
				Line:     75,
				Function: "github.com/hinoguma/go-TestNewStackTraceItem",
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
			assertEqualsStackTraceItem(t, item, tc.expected, "github.com/hinoguma/go-fault/structurederror.")
		})
	}
}

func TestStackTrace_JsonValueString(t *testing.T) {

	testCases := []struct {
		label    string
		trace    StackTrace
		expected string
	}{
		{
			label: "basic stack trace",
			trace: StackTrace{
				{
					File:     "file1.go",
					Line:     10,
					Function: "function1",
				},
				{
					File:     "file2.go",
					Line:     20,
					Function: "function2",
				},
			},
			expected: `[{"file":"file1.go","line":10,"function":"function1"},{"file":"file2.go","line":20,"function":"function2"}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.trace.JsonValueString()
			if got != tc.expected {
				t.Errorf("expected json value string %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestStackTraceItem_String(t *testing.T) {
	testCases := []struct {
		label    string
		item     StackTraceItem
		expected string
	}{
		{
			label: "basic stack trace item",
			item: StackTraceItem{
				File:     "file.go",
				Line:     42,
				Function: "myFunction",
			},
			expected: "myFunction() file.go:42",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.item.String()
			if got != tc.expected {
				t.Errorf("expected string %v, got %v", tc.expected, got)
			}
		})
	}
}
