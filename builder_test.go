package go_fault

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func builderTraceLevel1(w *StructuredErrorBuilder, skip int, depth int, level int) {
	if level <= 1 {
		w.StackTraceBuilderSkipDepth(skip, depth)
	}
	builderTraceLevel2(w, skip, depth, level)
}

func builderTraceLevel2(w *StructuredErrorBuilder, skip int, depth int, level int) {
	if level <= 2 {
		w.StackTraceBuilderSkipDepth(skip, depth)
	}
	builderTraceLevel3(w, skip, depth, level)
}

func builderTraceLevel3(w *StructuredErrorBuilder, skip int, depth int, level int) {
	if level <= 3 {
		w.StackTraceBuilderSkipDepth(skip, depth)
	}
	builderTraceLevel4(w, skip, depth, level)
}

func builderTraceLevel4(w *StructuredErrorBuilder, skip int, depth int, level int) {
	if level <= 4 {
		w.StackTraceBuilderSkipDepth(skip, depth)
	}
	builderTraceLevel5(w, skip, depth)
}

func builderTraceLevel5(w *StructuredErrorBuilder, skip int, depth int) {
	w.StackTraceBuilderSkipDepth(skip, depth)
}

func TestStructuredErrorBuilder_StackTraceBuilderSkipDepth(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		level    int
		skip     int
		depth    int
		expected StackTrace
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			level:    5,
			skip:     2,
			depth:    5,
			expected: nil,
		},
		{
			label:   "go standard error skip 0",
			wrapper: Builder(errStd),
			level:   5,
			skip:    0,
			depth:   1,
			expected: StackTrace{
				{
					Function: "github.com/hinoguma/go-fault.(*StructuredErrorBuilder).StackTraceBuilderSkipDepth",
				},
			},
		},
		{
			label:   "go standard error skip 1",
			wrapper: Builder(errStd),
			level:   5,
			skip:    1,
			depth:   1,
			expected: StackTrace{
				{
					Function: "github.com/hinoguma/go-fault.builderTraceLevel5",
				},
			},
		},
		{
			label:   "go standard error skip 2",
			wrapper: Builder(errStd),
			level:   5,
			skip:    2,
			depth:   1,
			expected: StackTrace{
				{
					Function: "github.com/hinoguma/go-fault.builderTraceLevel4",
				},
			},
		},
		{
			label:   "go standard error skip -1",
			wrapper: Builder(errStd),
			level:   5,
			skip:    -1,
			depth:   1,
			expected: StackTrace{
				{
					Function: "github.com/hinoguma/go-fault.(*StructuredErrorBuilder).StackTraceBuilderSkipDepth",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			builderTraceLevel1(tc.wrapper, tc.skip, tc.depth, tc.level)
			got := tc.wrapper.Build()
			if got == nil {
				if tc.expected != nil {
					t.Errorf("expected %v, got nil", tc.expected)
				}
				return
			}
			gotf, ok := got.(Structured)
			if !ok {
				t.Errorf("expected Structured, got %T", got)
				return
			}
			assertEqualsStackTrace(t, gotf.StackTrace(), tc.expected, "github.com/hinoguma/go-fault.")
		})
	}
}

func TestStructuredErrorBuilder_StackTrace(t *testing.T) {
	with := Builder(errStd)
	with.StackTrace()
	got := with.Build()
	if got == nil {
		t.Errorf("expected error, got nil")
		return
	}
	gotf, ok := got.(Structured)
	if !ok {
		t.Errorf("expected Structured, got %T", got)
		return
	}
	expected := StackTrace{
		{
			Function: "github.com/hinoguma/go-fault.TestStructuredErrorBuilder_StackTrace",
		},
		{
			Function: "testing.tRunner",
		},
		{
			Function: "runtime.goexit",
		},
	}
	assertEqualsStackTrace(t, gotf.StackTrace(), expected, "github.com/hinoguma/go-fault.")

	// nil error
	withNil := Builder(nil)
	withNil.StackTrace()
	gotNil := withNil.Build()
	if gotNil != nil {
		t.Errorf("expected nil, got %v", gotNil)
	}
}

var errStd = errors.New("standard error")
var errC3 = newTestCustomError3()

func TestStructuredErrorBuilder_Err(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			expected: nil,
		},
		{
			label:    "go standard error",
			wrapper:  Builder(errStd),
			expected: NewRawStructuredError(errStd),
		},
		{
			label:    "fault error",
			wrapper:  Builder(NewRawStructuredError(errStd)),
			expected: NewRawStructuredError(errStd),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_RequestID(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		value    string
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			value:    "12345",
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: Builder(errStd),
			value:   "12345",
			expected: NewRawStructuredError(errStd).
				SetRequestID("12345"),
		},
		{
			label:   "fault error",
			wrapper: Builder(NewRawStructuredError(errStd)),
			value:   "12345",
			expected: NewRawStructuredError(errStd).
				SetRequestID("12345"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.RequestID(tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_Type(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		value    ErrorType
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			value:    "12345",
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: Builder(errStd),
			value:   "12345",
			expected: NewRawStructuredError(errStd).
				SetType("12345"),
		},
		{
			label:   "fault error",
			wrapper: Builder(NewRawStructuredError(errStd)),
			value:   "12345",
			expected: NewRawStructuredError(errStd).
				SetType("12345"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.Type(tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_When(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		value    time.Time
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			value:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: Builder(errStd),
			value:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: NewRawStructuredError(errStd).
				SetWhen(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			label:   "fault error",
			wrapper: Builder(NewRawStructuredError(errStd)),
			value:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: NewRawStructuredError(errStd).
				SetWhen(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.When(tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_AddTagSafe(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		key      string
		value    TagValue
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			key:      "key1",
			value:    StringTagValue("value1"),
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: Builder(errStd),
			key:     "key1",
			value:   StringTagValue("value1"),
			expected: NewRawStructuredError(errStd).
				AddTagSafe("key1", StringTagValue("value1")),
		},
		{
			label:   "fault error",
			wrapper: Builder(NewRawStructuredError(errStd)),
			key:     "key1",
			value:   StringTagValue("value1"),
			expected: NewRawStructuredError(errStd).
				AddTagSafe("key1", StringTagValue("value1")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.AddTagSafe(tc.key, tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_AddTagString(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		key      string
		value    string
		expected error
	}{
		{
			label:   "string",
			wrapper: Builder(errStd),
			key:     "key1",
			value:   "value1",
			expected: NewRawStructuredError(errStd).
				AddTagSafe("key1", StringTagValue("value1")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.AddTagString(tc.key, tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_AddTagInt(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		key      string
		value    int
		expected error
	}{
		{
			label:   "int",
			wrapper: Builder(errStd),
			key:     "key1",
			value:   42,
			expected: NewRawStructuredError(errStd).
				AddTagSafe("key1", IntTagValue(42)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.AddTagInt(tc.key, tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_AddTagFloat(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		key      string
		value    float64
		expected error
	}{
		{
			label:   "int",
			wrapper: Builder(errStd),
			key:     "key1",
			value:   42,
			expected: NewRawStructuredError(errStd).
				AddTagSafe("key1", FloatTagValue(42)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.AddTagFloat(tc.key, tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_AddTagBool(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		key      string
		value    bool
		expected error
	}{
		{
			label:   "bool",
			wrapper: Builder(errStd),
			key:     "key1",
			value:   true,
			expected: NewRawStructuredError(errStd).
				AddTagSafe("key1", BoolTagValue(true)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.AddTagBool(tc.key, tc.value).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestStructuredErrorBuilder_DeleteTag(t *testing.T) {
	testCases := []struct {
		label    string
		wrapper  *StructuredErrorBuilder
		key      string
		expected error
	}{
		{
			label:    "nil error",
			wrapper:  Builder(nil),
			key:      "key1",
			expected: nil,
		},
		{
			label:   "go standard error",
			wrapper: Builder(errStd),
			key:     "key1",
			expected: NewRawStructuredError(errStd).
				DeleteTag("key1"),
		},
		{
			label:   "fault error",
			wrapper: Builder(NewRawStructuredError(errStd)),
			key:     "key1",
			expected: NewRawStructuredError(errStd).
				DeleteTag("key1"),
		},
		{
			label:    "fault error with existing tag",
			wrapper:  Builder(NewRawStructuredError(errStd).AddTagSafe("key1", StringTagValue("value1"))),
			key:      "key1",
			expected: NewRawStructuredError(errStd),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.wrapper.DeleteTag(tc.key).Build()
			if !reflect.DeepEqual(tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}
