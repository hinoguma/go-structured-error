package fault

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestWith(t *testing.T) {
	testCases := []struct {
		label   string
		err     error
		options []WithFunc
		expect  Fault
	}{}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			wrapped := With(tc.err, tc.options...)
			if tc.err == nil || tc.expect == nil {
				if wrapped != tc.expect {
					t.Errorf("With() = %+v, want %+v", wrapped, tc.expect)
				}
				return
			}
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
			label:  "nil error",
			err:    nil,
			value:  ErrorTypeNone,
			expect: nil,
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			value: ErrorTypeNone,
			expect: &FaultError{
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
			err: &FaultError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			value: ErrorTypeNone,
			expect: &FaultError{
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
			if tc.err == nil || tc.expect == nil {
				if wrapped != tc.expect {
					t.Errorf("WithType() = %+v, want %+v", wrapped, tc.expect)
				}
				return
			}
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
			label:  "nil error",
			err:    nil,
			id:     "request-123",
			expect: nil,
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			id:    "request-456",
			expect: &FaultError{
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
			err: &FaultError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			id: "new-request-id",
			expect: &FaultError{
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
			if tc.err == nil || tc.expect == nil {
				if wrapped != tc.expect {
					t.Errorf("WithRequestID() = %+v, want %+v", wrapped, tc.expect)
				}
				return
			}
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
			label:  "nil error",
			err:    nil,
			value:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: nil,
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			value: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: &FaultError{
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
			err: &FaultError{
				errorType:  ErrorTypeNone,
				err:        errors.New("existing fault error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "old-request-id",
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			value: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: &FaultError{
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
			if tc.err == nil || tc.expect == nil {
				if wrapped != tc.expect {
					t.Errorf("WithWhen() = %+v, want %+v", wrapped, tc.expect)
				}
				return
			}
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
			label:  "nil error",
			err:    nil,
			key:    "key1",
			value:  StringTagValue("value1"),
			expect: nil,
		},
		{
			label: "non-fault error",
			err:   errors.New("some error"),
			key:   "key2",
			value: IntTagValue(42),
			expect: &FaultError{
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
			err: &FaultError{
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
			expect: &FaultError{
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
			if tc.err == nil || tc.expect == nil {
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
