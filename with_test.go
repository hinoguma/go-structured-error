package go_fault

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestWith(t *testing.T) {
	tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		label   string
		err     error
		options []WithFunc
		expect  error
	}{
		{
			label: "nil error",
			err:   nil,
			options: []WithFunc{
				WithRequestID("request-123"),
				WithWhen(tm),
			},
			expect: &StructuredError{
				errorType:  "",
				err:        nil,
				stacktrace: make(StackTrace, 0),
				when:       &tm,
				requestId:  "request-123",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			options: []WithFunc{
				WithType(ErrorTypeNone),
				WithRequestID("request-456"),
				WithTagSafe("key1", StringTagValue("value1")),
			},
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("some error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "request-456",
				tags: Tags{
					tags: []Tag{
						{
							Key:   "key1",
							Value: StringTagValue("value1"),
						},
					},
					keyMap: map[string]int{
						"key1": 0,
					},
				},
				subErrors: make([]error, 0),
			},
		},
		{
			label: "existing fault error",
			err: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			options: []WithFunc{
				WithType(ErrorTypeNone),
				WithRequestID("new-request-id"),
				WithWhen(tm),
			},
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       &tm,
				requestId:  "new-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			wrapped := With(tc.err, tc.options...)
			if !reflect.DeepEqual(wrapped, tc.expect) {
				t.Errorf("With() = %+v, want %+v", wrapped, tc.expect)
			}
		})
	}

}

func TestWithType(t *testing.T) {
	testCases := []struct {
		label  string
		err    error
		value  ErrorType
		expect error
	}{
		{
			label: "nil error",
			err:   nil,
			value: ErrorTypeNone,
			expect: &StructuredError{
				errorType:  "",
				err:        nil,
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			value: ErrorTypeNone,
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("some error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "existing fault error",
			err: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			value: ErrorTypeNone,
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			wrapped := WithType(tc.value)(tc.err)
			if !reflect.DeepEqual(wrapped, tc.expect) {
				t.Errorf("WithType() = %+v, want %+v", wrapped, tc.expect)
			}
		})
	}
}

func TestWithRequestID(t *testing.T) {
	testCases := []struct {
		label  string
		err    error
		id     string
		expect error
	}{
		{
			label: "nil error",
			err:   nil,
			id:    "request-123",
			expect: &StructuredError{
				errorType:  "",
				err:        nil,
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "request-123",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			id:    "request-456",
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("some error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "request-456",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "existing fault error",
			err: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			id: "new-request-id",
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "new-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			wrapped := WithRequestID(tc.id)(tc.err)
			if !reflect.DeepEqual(wrapped, tc.expect) {
				t.Errorf("WithRequestID() = %+v, want %+v", wrapped, tc.expect)
			}
		})
	}
}

func TestWithWhen(t *testing.T) {
	tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		label  string
		err    error
		value  time.Time
		expect error
	}{
		{
			label: "nil error",
			err:   nil,
			value: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: &StructuredError{
				errorType:  "",
				err:        nil,
				stacktrace: make(StackTrace, 0),
				when:       &tm,
				requestId:  "",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			value: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("some error"),
				stacktrace: make(StackTrace, 0),
				when:       &tm,
				requestId:  "",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "existing fault error",
			err: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			value: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       &tm,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			wrapped := WithWhen(tc.value)(tc.err)
			if !reflect.DeepEqual(wrapped, tc.expect) {
				t.Errorf("WithWhen() = %+v, want %+v", wrapped, tc.expect)
			}
		})
	}
}

func TestWithTagSafe(t *testing.T) {
	testCases := []struct {
		label  string
		err    error
		key    string
		value  TagValue
		expect error
	}{
		{
			label: "nil error",
			err:   nil,
			key:   "key1",
			value: StringTagValue("value1"),
			expect: &StructuredError{
				err:        nil,
				stacktrace: make(StackTrace, 0),
				subErrors:  make([]error, 0),
				tags: Tags{
					tags: []Tag{
						{
							Key:   "key1",
							Value: StringTagValue("value1"),
						},
					},
					keyMap: map[string]int{
						"key1": 0,
					},
				},
			},
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			key:   "key2",
			value: IntTagValue(42),
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("some error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags: Tags{
					tags: []Tag{
						{
							Key:   "key2",
							Value: IntTagValue(42),
						},
					},
					keyMap: map[string]int{
						"key2": 0,
					},
				},
				subErrors: make([]error, 0),
			},
		},
		{
			label: "existing fault error",
			err: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			key:   "key3",
			value: BoolTagValue(true),
			expect: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags: Tags{
					tags: []Tag{
						{
							Key:   "key3",
							Value: BoolTagValue(true),
						},
					},
					keyMap: map[string]int{
						"key3": 0,
					},
				},
				subErrors: make([]error, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			wrapped := WithTagSafe(tc.key, tc.value)(tc.err)
			if wrapped == nil || tc.expect == nil {
				if wrapped != tc.expect {
					t.Errorf("WithTagSafe() = %+v, want %+v", wrapped, tc.expect)
				}
				return
			}
			if !reflect.DeepEqual(wrapped, tc.expect) {
				t.Errorf("WithTagSafe() = %+v, want %+v", wrapped, tc.expect)
			}
		})
	}
}
