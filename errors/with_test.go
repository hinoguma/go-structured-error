package errors

import (
	"errors"
	"github.com/hinoguma/go-fault"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	LineTraceLevel1             = 25
	LineTraceLevel2             = 33
	LineTraceLevel3             = 39
	LineTraceLevel4             = 46
	LineTraceLevel5             = 50
	LineStackTraceWithSkipDepth = 46
)

func traceLevel1(w *WithWrapper, skip int, depth int, level int) {
	if level <= 1 {
		w.StackTraceWithSkipDepth(skip, depth)
	}
	traceLevel2(w, skip, depth, level)
}

func traceLevel2(w *WithWrapper, skip int, depth int, level int) {
	if level <= 2 {
		w.StackTraceWithSkipDepth(skip, depth)
	}
	traceLevel3(w, skip, depth, level)
}

func traceLevel3(w *WithWrapper, skip int, depth int, level int) {
	if level <= 3 {
		w.StackTraceWithSkipDepth(skip, depth)
	}
	traceLevel4(w, skip, depth, level)
}

func traceLevel4(w *WithWrapper, skip int, depth int, level int) {
	if level <= 4 {
		w.StackTraceWithSkipDepth(skip, depth)
	}
	traceLevel5(w, skip, depth)
}

func traceLevel5(w *WithWrapper, skip int, depth int) {
	w.StackTraceWithSkipDepth(skip, depth)
}

func TestWithWrapper_StackTraceWithSkipDepth(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		level    int
		skip     int
		depth    int
		expected fault.StackTrace
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			level:    5,
			skip:     2,
			depth:    5,
			expected: nil,
		},
		{
			label:   "go standard error skip 0",
			wrapper: With(stdErr),
			level:   5,
			skip:    0,
			depth:   1,
			expected: fault.StackTrace{
				{
					File:     "ignored",
					Line:     LineStackTraceWithSkipDepth,
					Function: "github.com/hinoguma/go-fault/errors.(*WithWrapper).StackTraceWithSkipDepth",
				},
			},
		},
		{
			label:   "go standard error skip 1",
			wrapper: With(stdErr),
			level:   5,
			skip:    1,
			depth:   1,
			expected: fault.StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel5,
					Function: "github.com/hinoguma/go-fault/errors.traceLevel5",
				},
			},
		},
		{
			label:   "go standard error skip 2",
			wrapper: With(stdErr),
			level:   5,
			skip:    2,
			depth:   1,
			expected: fault.StackTrace{
				{
					File:     "ignored",
					Line:     LineTraceLevel4,
					Function: "github.com/hinoguma/go-fault/errors.traceLevel4",
				},
			},
		},
		{
			label:   "go standard error skip -1",
			wrapper: With(stdErr),
			level:   5,
			skip:    -1,
			depth:   1,
			expected: fault.StackTrace{
				{
					File:     "ignored",
					Line:     LineStackTraceWithSkipDepth,
					Function: "github.com/hinoguma/go-fault/errors.(*WithWrapper).StackTraceWithSkipDepth",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			traceLevel1(tc.wrapper, tc.skip, tc.depth, tc.level)
			got := tc.wrapper.Err()
			if got == nil {
				if tc.expected != nil {
					t.Errorf("expected %v, got nil", tc.expected)
				}
				return
			}
			gotf, ok := got.(fault.Fault)
			if !ok {
				t.Errorf("expected fault.Fault, got %T", got)
				return
			}
			assertEqualsStackTrace(t, gotf.StackTrace(), tc.expected, "github.com/hinoguma/go-fault")
		})
	}
}

func TestWithWrapper_StackTrace(t *testing.T) {
	with := With(stdErr)
	with.StackTrace()
	got := with.Err()
	if got == nil {
		t.Errorf("expected error, got nil")
		return
	}
	gotf, ok := got.(fault.Fault)
	if !ok {
		t.Errorf("expected fault.Fault, got %T", got)
		return
	}
	expected := fault.StackTrace{
		{
			File:     "ignored",
			Line:     150,
			Function: "github.com/hinoguma/go-fault/errors.TestWithWrapper_StackTrace",
		},
		{
			File:     "ignored",
			Line:     0,
			Function: "testing.tRunner",
		},
		{
			File:     "ignored",
			Line:     0,
			Function: "runtime.goexit",
		},
	}
	assertEqualsStackTrace(t, gotf.StackTrace(), expected, "github.com/hinoguma/go-fault")

	// nil error
	withNil := With(nil)
	withNil.StackTrace()
	gotNil := withNil.Err()
	if gotNil != nil {
		t.Errorf("expected nil, got %v", gotNil)
	}
}

var stdErr = errors.New("standard error")

type testCustomFaultError struct {
	fault.FaultError
}

func assertEqualsStackTraceItem(t *testing.T, got, expected fault.StackTraceItem, filterPrefix string) {
	// only check traces from this package
	// runtime and file system depends on environment
	if !strings.HasPrefix(got.Function, filterPrefix) {
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

func assertEqualsStackTrace(t *testing.T, got, expected fault.StackTrace, filterPrefix string) {
	if len(got) != len(expected) {
		t.Errorf("expected stack trace length %v, got %v expected:%v got :%v", len(expected), len(got), expected, got)
		return
	}
	for i := range got {
		assertEqualsStackTraceItem(t, got[i], expected[i], filterPrefix)
	}
}

func TestWithWrapper_convertToFault(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		expected fault.Fault
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			expected: nil,
		},
		{
			label:    "go standard error",
			wrapper:  With(stdErr),
			expected: fault.NewRawFaultError(stdErr),
		},
		{
			label:    "fault error",
			wrapper:  With(fault.NewRawFaultError(stdErr)),
			expected: fault.NewRawFaultError(stdErr),
		},
		{
			label:    "custom fault error",
			wrapper:  With(&testCustomFaultError{}),
			expected: &testCustomFaultError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.convertToFault()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestWithWrapper_Err(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			expected: nil,
		},
		{
			label:    "go standard error",
			wrapper:  With(stdErr),
			expected: stdErr,
		},
		{
			label:    "fault error",
			wrapper:  With(fault.NewRawFaultError(stdErr)),
			expected: fault.NewRawFaultError(stdErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.Err()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestWithWrapper_RequestID(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		value    string
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			value:    "12345",
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: With(stdErr),
			value:   "12345",
			expected: fault.NewRawFaultError(stdErr).
				SetRequestID("12345"),
		},
		{
			label:   "fault error",
			wrapper: With(fault.NewRawFaultError(stdErr)),
			value:   "12345",
			expected: fault.NewRawFaultError(stdErr).
				SetRequestID("12345"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.RequestID(tc.value).Err()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestWithWrapper_When(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		value    time.Time
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			value:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: With(stdErr),
			value:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: fault.NewRawFaultError(stdErr).
				SetWhen(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			label:   "fault error",
			wrapper: With(fault.NewRawFaultError(stdErr)),
			value:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: fault.NewRawFaultError(stdErr).
				SetWhen(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.When(tc.value).Err()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestWithWrapper_AddTagSafe(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		key      string
		value    fault.TagValue
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			key:      "key1",
			value:    fault.StringTagValue("value1"),
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: With(stdErr),
			key:     "key1",
			value:   fault.StringTagValue("value1"),
			expected: fault.NewRawFaultError(stdErr).
				AddTagSafe("key1", fault.StringTagValue("value1")),
		},
		{
			label:   "fault error",
			wrapper: With(fault.NewRawFaultError(stdErr)),
			key:     "key1",
			value:   fault.StringTagValue("value1"),
			expected: fault.NewRawFaultError(stdErr).
				AddTagSafe("key1", fault.StringTagValue("value1")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.AddTagSafe(tc.key, tc.value).Err()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestWithWrapper_DeleteTag(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *WithWrapper
		key      string
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  With(nil),
			key:      "key1",
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: With(stdErr),
			key:     "key1",
			expected: fault.NewRawFaultError(stdErr).
				DeleteTag("key1"),
		},
		{
			label:   "fault error",
			wrapper: With(fault.NewRawFaultError(stdErr)),
			key:     "key1",
			expected: fault.NewRawFaultError(stdErr).
				DeleteTag("key1"),
		},
		{
			label:    "fault error with existing tag",
			wrapper:  With(fault.NewRawFaultError(stdErr).AddTagSafe("key1", fault.StringTagValue("value1"))),
			key:      "key1",
			expected: fault.NewRawFaultError(stdErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.DeleteTag(tc.key).Err()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}
