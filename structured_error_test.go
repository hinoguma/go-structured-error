package go_fault

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestStructuredError_Type(t *testing.T) {
	testCases := []struct {
		label    string
		err      *StructuredError
		expected ErrorType
	}{
		{
			label:    "initial type",
			err:      &StructuredError{},
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

func TestStructuredError_When(t *testing.T) {
	testCases := []struct {
		label    string
		err      *StructuredError
		expected *time.Time
	}{
		{
			label:    "when is nil",
			err:      &StructuredError{},
			expected: nil,
		},
		{
			label: "when is set",
			err: &StructuredError{
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

func TestStructuredError_RequestID(t *testing.T) {
	testCases := []struct {
		label    string
		err      *StructuredError
		expected string
	}{
		{
			label:    "empty request ID",
			err:      &StructuredError{},
			expected: "",
		},
		{
			label: "set request ID",
			err: &StructuredError{
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

func TestStructuredError_StackTrace(t *testing.T) {
	testCases := []struct {
		label    string
		err      *StructuredError
		expected StackTrace
	}{
		{
			label:    "empty stack trace",
			err:      &StructuredError{},
			expected: StackTrace{},
		},
		{
			label: "stack trace with items",
			err: &StructuredError{
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-structured-error.TestStructuredError_StackTrace",
					},
				},
			},
			expected: StackTrace{
				{
					File:     "fault_test.go",
					Line:     75,
					Function: "github.com/hinoguma/go-structured-error.TestStructuredError_StackTrace",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.StackTrace()
			assertEqualsStackTrace(t, got, tc.expected, "github.com/hinoguma/go-structured-error.")
		})
	}
}

func TestStructuredError_Setters(t *testing.T) {
	stdErr := errors.New("go standard error")
	testCases := []struct {
		label    string
		err      *StructuredError
		setFunc  func(err *StructuredError)
		expected *StructuredError
	}{
		{
			label: "set when, requestId and error",
			err:   &StructuredError{},
			setFunc: func(err *StructuredError) {
				t := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
				_ = err.SetWhen(t).
					SetType("testType").
					SetRequestID("12345").
					SetErr(stdErr)
			},
			expected: &StructuredError{
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
			assertStructuredError(t, tc.err, tc.expected)
		})
	}
}

func TestStructuredError_AddTag(t *testing.T) {
	testCases := []struct {
		label    string
		err      *StructuredError
		addFunc  func(err *StructuredError)
		expected *StructuredError
	}{
		{
			label: "add tags",
			err:   &StructuredError{},
			addFunc: func(err *StructuredError) {
				_ = err.AddTagString("tag1", "value1")
				_ = err.AddTagInt("tag2", 42)
				_ = err.AddTagBool("tag3", true)
				_ = err.AddTagFloat("tag4", 3.14)
				_ = err.AddTagSafe("tag5", StringTagValue("safeValue"))
			},
			expected: &StructuredError{
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
			assertStructuredError(t, tc.err, tc.expected)
		})
	}
}

func TestStructuredError_DeleteTag(t *testing.T) {
	testCases := []struct {
		label     string
		err       *StructuredError
		deleteKey string
		expected  *StructuredError
	}{
		{
			label: "delete existing tag",
			err: &StructuredError{
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
			expected: &StructuredError{
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
			err: &StructuredError{
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
			expected: &StructuredError{
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
			err:       &StructuredError{},
			deleteKey: "tag1",
			expected:  &StructuredError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			_ = tc.err.DeleteTag(tc.deleteKey)
			assertStructuredError(t, tc.err, tc.expected)
		})
	}
}

func TestStructuredError_AddSubError(t *testing.T) {
	sub1 := errors.New("sub error 1")
	sub2 := errors.New("sub error 2")
	testCases := []struct {
		label    string
		err      *StructuredError
		subErr   error
		expected *StructuredError
	}{
		{
			label:  "add sub-error",
			err:    &StructuredError{},
			subErr: sub1,
			expected: &StructuredError{
				subErrors: []error{sub1},
			},
		},
		{
			label: "add another sub-error",
			err: &StructuredError{
				subErrors: []error{sub1},
			},
			subErr: sub2,
			expected: &StructuredError{
				subErrors: []error{sub1, sub2},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			_ = tc.err.AddSubError(tc.subErr)
			assertStructuredError(t, tc.err, tc.expected)
		})
	}
}

func TestStructuredError_Error(t *testing.T) {
	testCases := []struct {
		label    string
		err      *StructuredError
		expected string
	}{
		{
			label: "basic error message",
			err: &StructuredError{
				err: errors.New("basic error"),
			},
			expected: "[Type: none] basic error",
		},
		{
			label:    "has no underlying error",
			err:      &StructuredError{},
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

func TestNewRawStructuredError(t *testing.T) {
	stdErr := errors.New("standard error")
	expected := &StructuredError{
		errorType:  ErrorTypeNone,
		err:        stdErr,
		stacktrace: make(StackTrace, 0),
		when:       nil,
		requestId:  "",
		tags:       NewTags(),
		subErrors:  make([]error, 0),
	}

	got := NewRawStructuredError(stdErr)
	assertStructuredError(t, got, expected)
}

func TestNewWithSkipAndDepth(t *testing.T) {
	stdErr := errors.New("go standard error")
	testCases := []struct {
		label    string
		skip     int
		depth    int
		expected *StructuredError
	}{
		{
			label: "skip 0 starts capturing from NewWithSkipAndDepth",
			skip:  0,
			depth: 1,
			expected: &StructuredError{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						Function: "github.com/hinoguma/go-structured-error.NewWithSkipAndDepth",
					},
				},
				tags: NewTags(),
			},
		},
		{
			label: "skip -1 treated as skip 0",
			skip:  -1,
			depth: 1,
			expected: &StructuredError{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						Function: "github.com/hinoguma/go-structured-error.NewWithSkipAndDepth",
					},
				},
				tags: NewTags(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := NewWithSkipAndDepth(stdErr, tc.skip, tc.depth)
			assertStructuredErrorWithErrorValue(t, got, tc.expected)
		})
	}
}

func TestStructuredError_JsonPrinter(t *testing.T) {
	stdErr := errors.New("go standard error")
	when := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	faultErr := &StructuredError{
		errorType: ErrorTypeNone,
		err:       stdErr,
	}

	testCases := []struct {
		label    string
		err      *StructuredError
		expected JsonPrinter
	}{
		{
			label: "all props",
			err: &StructuredError{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-structured-error.TestStructuredError_JsonFormat",
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
			expected: ErrorJsonPrinter{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-structured-error.TestStructuredError_JsonFormat",
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
			got := tc.err.JsonPrinter()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected JsonFormat %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestStructuredError_JsonString(t *testing.T) {
	stdErr := errors.New("go standard error")
	when := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	faultErr := &StructuredError{
		errorType: ErrorTypeNone,
		err:       errors.New("go standard error2"),
	}

	testCases := []struct {
		label    string
		err      *StructuredError
		expected string
	}{
		{
			label: "all props",
			err: &StructuredError{
				errorType: ErrorTypeNone,
				err:       stdErr,
				stacktrace: StackTrace{
					{
						File:     "fault_test.go",
						Line:     75,
						Function: "github.com/hinoguma/go-structured-error.TestStructuredError_JsonString",
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
			expected: `{"type":"none","message":"go standard error","when":"2024-06-01T12:00:00Z","request_id":"12345","tags":{"tag1":"value1"},"stacktrace":[{"file":"fault_test.go","line":75,"function":"github.com/hinoguma/go-structured-error.TestStructuredError_JsonString"}],"sub_errors":[{"type":"none","message":"go standard error","stacktrace":[]},{"type":"none","message":"go standard error2","stacktrace":[]}]}`,
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

func TestStructuredError_Format(t *testing.T) {

	testCases := []struct {
		label         string
		err           *StructuredError
		expectedS     string
		expectedV     string
		expectedVPlus string
		expectedQ     string
	}{
		{
			label: "basic format",
			err: &StructuredError{
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
