package serrors

import "testing"

func TestStructuredError_WithStackTrace(t *testing.T) {
	// WithStackTrace should start capturing from the caller of WithStackTrace
	err := &StructuredError{}
	_ = err.WithStackTrace() // 8
	if len(err.stacktrace) == 0 {
		t.Errorf("expected stacktrace to be set, but it was empty")
	}
	firstFrame := err.stacktrace[0]
	if firstFrame.Function != "github.com/hinoguma/go-structured-error.TestStructuredError_WithStackTrace" {
		t.Errorf("expected top stack frame to be TestStructuredError_WithStackTrace, but got %s", err.stacktrace[0].Function)
	}
	if firstFrame.Line != 8 {
		t.Errorf("expected top stack frame line to be 9, but got %d", err.stacktrace[0].Line)
	}
}

func setStackTraceWithSkipMaxDepth1(err *StructuredError, skip, depth int) {
	setStackTraceWithSkipMaxDepth2(err, skip, depth)
}
func setStackTraceWithSkipMaxDepth2(err *StructuredError, skip, depth int) {
	setStackTraceWithSkipMaxDepth3(err, skip, depth)
}
func setStackTraceWithSkipMaxDepth3(err *StructuredError, skip, depth int) {
	_ = err.SetStackTraceWithSkipMaxDepth(skip, depth)
}

func TestStructuredError_SetStackTraceWithSkipMaxDepth(t *testing.T) {
	testCases := []struct {
		label    string
		skip     int
		depth    int
		expected StackTrace
	}{
		{
			label: "skip 0 starts capturing from SetStackTraceWithSkipMaxDepth",
			skip:  0,
			depth: 2,
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     150,
					Function: "github.com/hinoguma/go-structured-error.(*StructuredError).SetStackTraceWithSkipMaxDepth",
				},
				{
					File:     "ignored",
					Line:     28,
					Function: "github.com/hinoguma/go-structured-error.setStackTraceWithSkipMaxDepth3",
				},
			},
		},
		{
			label: "skip 1 starts capturing from setStackTraceWithSkipMaxDepth3, the caller of SetStackTraceWithSkipMaxDepth",
			skip:  1,
			depth: 2,
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     28,
					Function: "github.com/hinoguma/go-structured-error.setStackTraceWithSkipMaxDepth3",
				},
				{
					File:     "ignored",
					Line:     25,
					Function: "github.com/hinoguma/go-structured-error.setStackTraceWithSkipMaxDepth2",
				},
			},
		},
		{
			label: "skip -1 treated as skip 0",
			skip:  -1,
			depth: 1,
			expected: StackTrace{
				{
					File:     "ignored",
					Line:     150,
					Function: "github.com/hinoguma/go-structured-error.(*StructuredError).SetStackTraceWithSkipMaxDepth",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := &StructuredError{}
			setStackTraceWithSkipMaxDepth1(err, tc.skip, tc.depth)
			assertEqualsStackTrace(t, err.stacktrace, tc.expected, "github.com/hinoguma/go-structured-error.")
		})
	}

}
