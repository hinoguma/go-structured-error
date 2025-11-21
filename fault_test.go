package fault

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func assertFaultError(t *testing.T, got, expected *FaultError) {
	if got.err != expected.err {
		t.Errorf("expected err %v, got %v", expected.err, got.err)
	}
	if got.errorType != expected.errorType {
		t.Errorf("expected errorType %v, got %v", expected.errorType, got.errorType)
	}
	if got.when == nil || expected.when == nil {
		if got.when != expected.when {
			t.Errorf("expected when %v, got %v", expected.when, got.when)
		}
	} else {
		if !got.when.Equal(*expected.when) {
			t.Errorf("expected when %v, got %v", *expected.when, *got.when)
		}
	}
	if got.requestId != expected.requestId {
		t.Errorf("expected requestId %v, got %v", expected.requestId, got.requestId)
	}
	assertEqualsStackTrace(t, got.stacktrace, expected.stacktrace, "github.com/hinoguma/go-fault")
	assertEqualsTags(t, got.tags, expected.tags)
	if len(got.subErrors) != len(expected.subErrors) {
		t.Errorf("expected subErrors length %v, got %v", len(expected.subErrors), len(got.subErrors))
	} else {
		for i := range got.subErrors {
			if got.subErrors[i] != expected.subErrors[i] {
				t.Errorf("expected subError %v, got %v", expected.subErrors[i], got.subErrors[i])
			}
		}
	}
}

func assertFaultErrorWithErrorValue(t *testing.T, got, expected *FaultError) {
	if got.err == nil || expected.err == nil {
		if got.err != expected.err {
			t.Errorf("expected err %v, got %v", expected.err, got.err)
		}
	} else {
		if got.err.Error() != expected.err.Error() {
			t.Errorf("expected err %v, got %v", expected.err, got.err)
		}
	}
	if got.errorType != expected.errorType {
		t.Errorf("expected errorType %v, got %v", expected.errorType, got.errorType)
	}
	if got.when == nil || expected.when == nil {
		if got.when != expected.when {
			t.Errorf("expected when %v, got %v", expected.when, got.when)
		}
	} else {
		if !got.when.Equal(*expected.when) {
			t.Errorf("expected when %v, got %v", *expected.when, *got.when)
		}
	}
	if got.requestId != expected.requestId {
		t.Errorf("expected requestId %v, got %v", expected.requestId, got.requestId)
	}
	assertEqualsStackTrace(t, got.stacktrace, expected.stacktrace, "github.com/hinoguma/go-fault")
	assertEqualsTags(t, got.tags, expected.tags)
	if len(got.subErrors) != len(expected.subErrors) {
		t.Errorf("expected subErrors length %v, got %v", len(expected.subErrors), len(got.subErrors))
	} else {
		for i := range got.subErrors {
			gotSubErr := got.subErrors[i]
			expectedSubErr := expected.subErrors[i]
			if gotSubErr == nil || expectedSubErr == nil {
				if gotSubErr != expectedSubErr {
					t.Errorf("expected subError %v, got %v", expectedSubErr, gotSubErr)
				}
			} else {
				if gotSubErr.Error() != expectedSubErr.Error() {
					t.Errorf("expected subError %v, got %v", expectedSubErr, gotSubErr)
				}
			}
		}
	}
}

func TestFaultError_Type(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		expected ErrorType
	}{
		{
			label:    "initial type",
			err:      &FaultError{},
			expected: ErrorTypeNone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			if tc.err.Type() != tc.expected {
				t.Errorf("expected type %v, got %v", tc.expected, tc.err.Type())
			}
		})
	}
}

func TestFaultError_When(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		expected *time.Time
	}{
		{
			label:    "when is nil",
			err:      &FaultError{},
			expected: nil,
		},
		{
			label: "when is set",
			err: &FaultError{
				when: func() *time.Time {
					t := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
					return &t
				}(),
			},
			expected: func() *time.Time {
				t := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
				return &t
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.When()
			if got == nil || tc.expected == nil {
				if got != tc.expected {
					t.Errorf("expected when %v, got %v", tc.expected, got)
				}
				return
			}
			if !got.Equal(*tc.expected) {
				t.Errorf("expected when %v, got %v", *tc.expected, *got)
			}
		})
	}
}

func TestFaultError_RequestID(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		expected string
	}{
		{
			label:    "empty request ID",
			err:      &FaultError{},
			expected: "",
		},
		{
			label: "set request ID",
			err: &FaultError{
				requestId: "12345",
			},
			expected: "12345",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.RequestID()
			if got != tc.expected {
				t.Errorf("expected request ID %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestFaultError_StackTrace(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		expected StackTrace
	}{
		{
			label:    "empty stack trace",
			err:      &FaultError{},
			expected: StackTrace{},
		},
		{
			label: "stack trace with items",
			err: &FaultError{
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-fault.TestFaultError_StackTrace",
					},
				},
			},
			expected: StackTrace{
				{
					File:     "fault_test.go",
					Line:     75,
					Function: "github.com/hinoguma/go-fault.TestFaultError_StackTrace",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.StackTrace()
			assertEqualsStackTrace(t, got, tc.expected, "github.com/hinoguma/go-fault")
		})
	}
}

func TestFaultError_Setters(t *testing.T) {
	stdErr := errors.New("go standard error")
	testCases := []struct {
		label    string
		err      *FaultError
		setFunc  func(err *FaultError)
		expected *FaultError
	}{
		{
			label: "set when, requestId and error",
			err:   &FaultError{},
			setFunc: func(err *FaultError) {
				t := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
				_ = err.SetWhen(t).
					SetType("testType").
					SetRequestID("12345").
					SetErr(stdErr)
			},
			expected: &FaultError{
				errorType: "testType",
				when: func() *time.Time {
					t := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
					return &t
				}(),
				requestId: "12345",
				err:       stdErr,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			tc.setFunc(tc.err)
			assertFaultError(t, tc.err, tc.expected)
		})
	}
}

func TestFaultError_AddTag(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		addFunc  func(err *FaultError)
		expected *FaultError
	}{
		{
			label: "add tags",
			err:   &FaultError{},
			addFunc: func(err *FaultError) {
				_ = err.AddTagString("tag1", "value1")
				_ = err.AddTagInt("tag2", 42)
				_ = err.AddTagBool("tag3", true)
				_ = err.AddTagFloat("tag4", 3.14)
				_ = err.AddTagSafe("tag5", StringTagValue("safeValue"))
			},
			expected: &FaultError{
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
						{Key: "tag2", Value: IntTagValue(42)},
						{Key: "tag3", Value: BoolTagValue(true)},
						{Key: "tag4", Value: FloatTagValue(3.14)},
						{Key: "tag5", Value: StringTagValue("safeValue")},
					},
					keyMap: map[string]int{
						"tag1": 0,
						"tag2": 1,
						"tag3": 2,
						"tag4": 3,
						"tag5": 4,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			tc.addFunc(tc.err)
			assertFaultError(t, tc.err, tc.expected)
		})
	}
}

func TestFaultError_DeleteTag(t *testing.T) {
	testCases := []struct {
		label     string
		err       *FaultError
		deleteKey string
		expected  *FaultError
	}{
		{
			label: "delete existing tag",
			err: &FaultError{
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
						{Key: "tag2", Value: IntTagValue(42)},
					},
					keyMap: map[string]int{
						"tag1": 0,
						"tag2": 1,
					},
				},
			},
			deleteKey: "tag1",
			expected: &FaultError{
				tags: Tags{
					tags: []Tag{
						{Key: "tag2", Value: IntTagValue(42)},
					},
					keyMap: map[string]int{
						"tag2": 0,
					},
				},
			},
		},
		{
			label: "delete non-existing tag",
			err: &FaultError{
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
					},
					keyMap: map[string]int{
						"tag1": 0,
					},
				},
			},
			deleteKey: "tag2",
			expected: &FaultError{
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
					},
					keyMap: map[string]int{
						"tag1": 0,
					},
				},
			},
		},
		{
			label:     "delete tag from empty tags",
			err:       &FaultError{},
			deleteKey: "tag1",
			expected:  &FaultError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			_ = tc.err.DeleteTag(tc.deleteKey)
			assertFaultError(t, tc.err, tc.expected)
		})
	}
}

func TestFaultError_AddSubError(t *testing.T) {
	sub1 := errors.New("sub error 1")
	sub2 := errors.New("sub error 2")
	testCases := []struct {
		label    string
		err      *FaultError
		subErr   error
		expected *FaultError
	}{
		{
			label:  "add sub-error",
			err:    &FaultError{},
			subErr: sub1,
			expected: &FaultError{
				subErrors: []error{sub1},
			},
		},
		{
			label: "add another sub-error",
			err: &FaultError{
				subErrors: []error{sub1},
			},
			subErr: sub2,
			expected: &FaultError{
				subErrors: []error{sub1, sub2},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			_ = tc.err.AddSubError(tc.subErr)
			assertFaultError(t, tc.err, tc.expected)
		})
	}
}

func TestFaultError_Error(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		expected string
	}{
		{
			label: "basic error message",
			err: &FaultError{
				err: errors.New("basic error"),
			},
			expected: "[Type: none] basic error",
		},
		{
			label:    "has no underlying error",
			err:      &FaultError{},
			expected: "[Type: none] <no error>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.Error()
			if got != tc.expected {
				t.Errorf("expected error string %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestNewRawFaultError(t *testing.T) {
	stdErr := errors.New("standard error")
	expected := &FaultError{
		errorType:  ErrorTypeNone,
		err:        stdErr,
		stacktrace: make(StackTrace, 0),
		when:       nil,
		requestId:  "",
		tags:       NewTags(),
		subErrors:  make([]error, 0),
	}

	got := NewRawFaultError(stdErr)
	assertFaultError(t, got, expected)
}

func TestNew(t *testing.T) {
	testCases := []struct {
		label    string
		message  string
		expected *FaultError
	}{
		{
			label:   "basic fault error",
			message: "fault occurred",
			expected: &FaultError{
				errorType: ErrorTypeNone,
				err:       errors.New("fault occurred"),
				stacktrace: []StackTraceItem{
					{
						File:     "ignored",
						Line:     500,
						Function: "github.com/hinoguma/go-fault.TestNew.func1",
					},
					{
						File:     "ignored",
						Line:     -1,
						Function: "testing.tRunner",
					},
					{
						File:     "ignored",
						Line:     -1,
						Function: "runtime.goexit",
					},
				},
				when:      nil,
				requestId: "",
				tags:      NewTags(),
				subErrors: make([]error, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := New(tc.message) // 500
			assertFaultErrorWithErrorValue(t, got, tc.expected)
		})
	}
}

func TestFaultError_JsonFormatter(t *testing.T) {
	stdErr := errors.New("go standard error")
	when := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	faultErr := &FaultError{
		errorType: ErrorTypeNone,
		err:       stdErr,
	}

	testCases := []struct {
		label    string
		err      *FaultError
		expected ErrorFormatter
	}{
		{
			label: "all props",
			err: &FaultError{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-fault.TestFaultError_JsonFormat",
					},
				},
				when:      &when,
				requestId: "12345",
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
					},
					keyMap: map[string]int{
						"tag1": 0,
					},
				},
				subErrors: []error{stdErr, faultErr},
			},
			expected: JsonFormatter{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-fault.TestFaultError_JsonFormat",
					},
				},
				when:      &when,
				requestId: "12345",
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
					},
					keyMap: map[string]int{
						"tag1": 0,
					},
				},
				subErrors: []error{stdErr, faultErr},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.JsonFormatter()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected JsonFormat %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestFaultError_JsonString(t *testing.T) {
	stdErr := errors.New("go standard error")
	when := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	faultErr := &FaultError{
		errorType: ErrorTypeNone,
		err:       errors.New("go standard error2"),
	}

	testCases := []struct {
		label    string
		err      *FaultError
		expected string
	}{
		{
			label: "all props",
			err: &FaultError{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-fault.TestFaultError_JsonString",
					},
				},
				when:      &when,
				requestId: "12345",
				tags: Tags{
					tags: []Tag{
						{Key: "tag1", Value: StringTagValue("value1")},
					},
					keyMap: map[string]int{
						"tag1": 0,
					},
				},
				subErrors: []error{stdErr, faultErr},
			},
			expected: `{"type":"none","message":"go standard error","when":"2024-06-01T12:00:00Z","request_id":"12345","tags":{"tag1":"value1"},"stacktrace":[{"file":"fault_test.go","line":75,"function":"github.com/hinoguma/go-fault.TestFaultError_JsonString"}],"sub_errors":[{"type":"none","message":"go standard error","stacktrace":[]},{"type":"none","message":"go standard error2","stacktrace":[]}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.JsonString()
			if got != tc.expected {
				t.Errorf("expected JsonString %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestFaultError_Format(t *testing.T) {

	testCases := []struct {
		label         string
		err           *FaultError
		expectedS     string
		expectedV     string
		expectedVPlus string
		expectedQ     string
	}{
		{
			label: "basic format",
			err: &FaultError{
				errorType: ErrorTypeNone,
				err:       errors.New("basic error"),
			},
			expectedS: "[Type: none] basic error",
			expectedV: "[Type: none] basic error",
			expectedVPlus: `main_error:
    message: basic error
    type: none`,
			expectedQ: `"[Type: none] basic error"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			gotS := fmt.Sprintf("%s", tc.err)
			if gotS != tc.expectedS {
				t.Errorf("expectedS format %v, got %v", tc.expectedS, gotS)
			}
			gotV := fmt.Sprintf("%v", tc.err)
			if gotV != tc.expectedV {
				t.Errorf("expectedV format %v, got %v", tc.expectedV, gotV)
			}
			gotVPlus := fmt.Sprintf("%+v", tc.err)
			if gotVPlus != tc.expectedVPlus {
				t.Errorf("expectedVPlus format %v, got %v", tc.expectedVPlus, gotVPlus)
			}
			gotT := fmt.Sprintf("%q", tc.err)
			if gotT != tc.expectedQ {
				t.Errorf("expectedQ format %v, got %v", tc.expectedQ, gotT)
			}
		})
	}
}
