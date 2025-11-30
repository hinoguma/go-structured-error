package serrors

import (
	"errors"
	"testing"
	"time"
)

func TestErrorJsonPrinter_Print(t *testing.T) {
	testCases := []struct {
		label     string
		formatter ErrorJsonPrinter
		expected  string
	}{
		{
			label: "required fields",
			formatter: ErrorJsonPrinter{
				errorType: ErrorTypeNone,
				err:       errors.New("test error"),
				stacktrace: StackTrace{
					{
						File:     "example.go",
						Line:     10,
						Function: "main.exampleFunction",
					},
					{
						File:     "example.go",
						Line:     20,
						Function: "main.anotherFunction",
					},
				},
				when:      nil,
				requestId: "",
			},
			expected: `{"type":"none","message":"test error","stacktrace":[{"file":"example.go","line":10,"function":"main.exampleFunction"},{"file":"example.go","line":20,"function":"main.anotherFunction"}]}`,
		},
		{
			label: "with when and requestId",
			formatter: ErrorJsonPrinter{
				errorType:  ErrorType("testType"),
				err:        errors.New("another error"),
				stacktrace: make(StackTrace, 0),
				when: func() *time.Time {
					t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
					return &t
				}(),
				requestId: "req-12345",
			},
			expected: `{"type":"testType","message":"another error","when":"2024-01-01T12:00:00Z","request_id":"req-12345","stacktrace":[]}`,
		},
		{
			label: "tags",
			formatter: ErrorJsonPrinter{
				err:        errors.New("error with tags"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags: Tags{
					tags: []Tag{
						{Key: "session_id", Value: StringTagValue("sess-456")},
						{Key: "user_id", Value: StringTagValue("user-789")},
					},
					keyMap: map[string]int{
						"session_id": 0,
						"user_id":    1,
					},
				},
			},
			expected: `{"type":"none","message":"error with tags","tags":{"session_id":"sess-456","user_id":"user-789"},"stacktrace":[]}`,
		},
		{
			label: "empty tags",
			formatter: ErrorJsonPrinter{
				err:        errors.New("error with empty tags"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags: Tags{
					tags:   []Tag{},
					keyMap: map[string]int{},
				},
			},
			expected: `{"type":"none","message":"error with empty tags","stacktrace":[]}`,
		},
		{
			label: "sub errors",
			formatter: ErrorJsonPrinter{
				err:        errors.New("main error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				subErrors: []error{
					errors.New("sub error 1"),
					&StructuredError{
						err:       errors.New("sub error 2"),
						errorType: ErrorType("testType2"),
						stacktrace: StackTrace{
							{File: "sub_example.go", Line: 30, Function: "subFunction"},
						},
					},
				},
			},
			expected: `{"type":"none","message":"main error","stacktrace":[],"sub_errors":[{"type":"none","message":"sub error 1","stacktrace":[]},{"type":"testType2","message":"sub error 2","stacktrace":[{"file":"sub_example.go","line":30,"function":"subFunction"}]}]}`,
		},
		{
			label: "empty",
			formatter: ErrorJsonPrinter{
				err:        nil,
				stacktrace: nil,
				when:       nil,
				requestId:  "",
			},
			expected: `{"type":"none","message":"<no error>","stacktrace":[]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.formatter.Print()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestErrorVerbosePrinter_printSingle(t *testing.T) {
	testCases := []struct {
		label     string
		formatter ErrorVerbosePrinter
		expected  string
	}{
		{
			label: "required fields",
			formatter: ErrorVerbosePrinter{
				title:     "main_error",
				errorType: ErrorTypeNone,
				err:       errors.New("test error"),
				stacktrace: StackTrace{
					{
						File:     "example.go",
						Line:     10,
						Function: "main.exampleFunction",
					},
					{
						File:     "example.go",
						Line:     20,
						Function: "main.anotherFunction",
					},
				},
				when:      nil,
				requestId: "",
			},
			expected: `main_error:
    message: test error
    type: none
    stacktrace:
        main.exampleFunction() example.go:10
        main.anotherFunction() example.go:20`,
		},
		{
			label: "with when and requestId",
			formatter: ErrorVerbosePrinter{
				title:      "sub_error1",
				errorType:  ErrorType("testType"),
				err:        errors.New("another error"),
				stacktrace: make(StackTrace, 0),
				when: func() *time.Time {
					t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
					return &t
				}(),
				requestId: "req-12345",
			},
			expected: `sub_error1:
    message: another error
    type: testType
    when: 2024-01-01T12:00:00Z
    request_id: req-12345`,
		},
		{
			label: "tags",
			formatter: ErrorVerbosePrinter{
				title:      "sub_error2",
				err:        errors.New("error with tags"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags: Tags{
					tags: []Tag{
						{Key: "session_id", Value: StringTagValue("sess-456")},
						{Key: "user_id", Value: StringTagValue("user-789")},
					},
					keyMap: map[string]int{
						"session_id": 0,
						"user_id":    1,
					},
				},
			},
			expected: `sub_error2:
    message: error with tags
    type: none
    tags:
        session_id: sess-456
        user_id: user-789`,
		},
		{
			label: "empty tags",
			formatter: ErrorVerbosePrinter{
				title:      "sub_error3",
				err:        errors.New("error with empty tags"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags: Tags{
					tags:   []Tag{},
					keyMap: map[string]int{},
				},
			},
			expected: `sub_error3:
    message: error with empty tags
    type: none`,
		},
		{
			label: "empty",
			formatter: ErrorVerbosePrinter{
				title:      "main_error",
				err:        nil,
				stacktrace: nil,
				when:       nil,
				requestId:  "",
			},
			expected: `main_error:
    message: <no error>
    type: none`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.formatter.printSingle()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestErrorVerbosePrinter_Print(t *testing.T) {
	testCases := []struct {
		label     string
		formatter ErrorVerbosePrinter
		expected  string
	}{
		{
			label: "not sub errors",
			formatter: ErrorVerbosePrinter{
				title:     "main_error",
				errorType: ErrorTypeNone,
				err:       errors.New("test error"),
				stacktrace: StackTrace{
					{
						File:     "example.go",
						Line:     10,
						Function: "main.exampleFunction",
					},
				},
				when:      nil,
				requestId: "",
			},
			expected: `main_error:
    message: test error
    type: none
    stacktrace:
        main.exampleFunction() example.go:10`,
		},
		{
			label: "with sub errors",
			formatter: ErrorVerbosePrinter{
				title:      "main_error",
				err:        errors.New("main error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				subErrors: []error{
					errors.New("sub error 1"),
					&StructuredError{
						err:       errors.New("sub error 2"),
						errorType: ErrorType("testType2"),
						stacktrace: StackTrace{
							{File: "sub_example.go", Line: 30, Function: "subFunction"},
						},
						subErrors: []error{
							errors.New("nested sub error"),
							&StructuredError{
								err: errors.New("nested sub error 2"),
								stacktrace: StackTrace{
									{File: "sub_example.go", Line: 50, Function: "subFunction"},
								},
							},
						},
					},
				},
			},
			expected: `main_error:
    message: main error
    type: none
main_error.sub1:
    message: sub error 1
    type: none
main_error.sub2:
    message: sub error 2
    type: testType2
    stacktrace:
        subFunction() sub_example.go:30
main_error.sub2.sub1:
    message: nested sub error
    type: none
main_error.sub2.sub2:
    message: nested sub error 2
    type: none
    stacktrace:
        subFunction() sub_example.go:50`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.formatter.Print()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfMessage(t *testing.T) {
	testCases := []struct {
		label    string
		err      error
		expected string
	}{
		{
			label:    "normal error",
			err:      errors.New("sample error"),
			expected: `"message":"sample error"`,
		},
		{
			label:    "nil error",
			err:      nil,
			expected: `"message":"<no error>"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfMessage(tc.err)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfRequestID(t *testing.T) {
	testCases := []struct {
		label     string
		requestId string
		expected  string
	}{
		{
			label:     "non-empty request ID",
			requestId: "req-12345",
			expected:  `"request_id":"req-12345"`,
		},
		{
			label:     "empty request ID",
			requestId: "",
			expected:  `"request_id":""`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfRequestID(tc.requestId)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfWhen(t *testing.T) {
	testCases := []struct {
		label    string
		when     time.Time
		layout   string
		expected string
	}{
		{
			label:    "specific time",
			when:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			layout:   time.RFC3339,
			expected: `"when":"2024-01-01T12:00:00Z"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfWhen(tc.when, tc.layout)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfTags(t *testing.T) {
	testCases := []struct {
		label    string
		tags     Tags
		expected string
	}{
		{
			label: "multiple tags",
			tags: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			expected: `"tags":{"key1":"value1","key2":"value2"}`,
		},
		{
			label:    "no tags",
			tags:     Tags{tags: []Tag{}, keyMap: map[string]int{}},
			expected: `"tags":{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfTags(tc.tags)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfType(t *testing.T) {
	testCases := []struct {
		label    string
		errType  ErrorType
		expected string
	}{
		{
			label:    "defined error type",
			errType:  ErrorType("customType"),
			expected: `"type":"customType"`,
		},
		{
			label:    "none error type",
			errType:  ErrorTypeNone,
			expected: `"type":"none"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfType(tc.errType)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfStackTrace(t *testing.T) {
	testCases := []struct {
		label      string
		stacktrace StackTrace
		expected   string
	}{
		{
			label: "multiple stack frames",
			stacktrace: StackTrace{
				{File: "file1.go", Line: 10, Function: "func1"},
				{File: "file2.go", Line: 20, Function: "func2"},
			},
			expected: `"stacktrace":[{"file":"file1.go","line":10,"function":"func1"},{"file":"file2.go","line":20,"function":"func2"}]`,
		},
		{
			label:      "no stack frames",
			stacktrace: StackTrace{},
			expected:   `"stacktrace":[]`,
		},
		{
			label:      "nil stack trace",
			stacktrace: nil,
			expected:   `"stacktrace":[]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfStackTrace(tc.stacktrace)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestBuildJsonStringOfSubErrors(t *testing.T) {
	testCases := []struct {
		label     string
		subErrors []error
		expected  string
	}{
		{
			label: "multiple sub errors",
			subErrors: []error{
				errors.New("sub error 1"),
				&StructuredError{
					err:       errors.New("sub error 2"),
					errorType: ErrorType("type2"),
					stacktrace: StackTrace{
						{File: "file2.go", Line: 20, Function: "func2"},
					},
				},
			},
			expected: `"sub_errors":[{"type":"none","message":"sub error 1","stacktrace":[]},{"type":"type2","message":"sub error 2","stacktrace":[{"file":"file2.go","line":20,"function":"func2"}]}]`,
		},
		{
			label:     "no sub errors",
			subErrors: []error{},
			expected:  `"sub_errors":[]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := BuildJsonStringOfSubErrors(tc.subErrors)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
