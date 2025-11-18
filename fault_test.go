package fault

import (
	"errors"
	"testing"
	"time"
)

func assertFaultError(t *testing.T, got, expected *FaultError) {
	if got.err != expected.err {
		t.Errorf("expected err %v, got %v", expected.err, got.err)
	}
	if got.faultType != expected.faultType {
		t.Errorf("expected faultType %v, got %v", expected.faultType, got.faultType)
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
	assertEqualsStackTrace(t, got.stacktrace, expected.stacktrace)
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

func TestFaultError_Type(t *testing.T) {
	testCases := []struct {
		label    string
		err      *FaultError
		expected FaultType
	}{
		{
			label:    "initial type is util",
			err:      &FaultError{},
			expected: FaultTypeUtil,
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
						Function: "fault.TestFaultError_StackTrace",
					},
				},
			},
			expected: StackTrace{
				{
					File:     "fault_test.go",
					Line:     75,
					Function: "fault.TestFaultError_StackTrace",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.err.StackTrace()
			assertEqualsStackTrace(t, got, tc.expected)
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
				err.SetWhen(t).
					SetRequestID("12345").
					SetErr(stdErr)
			},
			expected: &FaultError{
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
				err.AddTagString("tag1", "value1").
					AddTagInt("tag2", 42).
					AddTagBool("tag3", true).
					AddTagFloat("tag4", 3.14).
					AddTagSafe("tag5", StringTagValue("safeValue"))
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
			tc.err.DeleteTag(tc.deleteKey)
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
			tc.err.AddSubError(tc.subErr)
			assertFaultError(t, tc.err, tc.expected)
		})
	}
}

func TestFaultError_Error(t *testing.T) {

}

// todo stacktrace, Is, As test
