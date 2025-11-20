package fault

import (
	"errors"
	"testing"
	"time"
)

func TestJsonFormatter_Format(t *testing.T) {
	testCases := []struct {
		label     string
		formatter JsonFormatter
		expected  string
	}{
		{
			label: "required fields",
			formatter: JsonFormatter{
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
			formatter: JsonFormatter{
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
			formatter: JsonFormatter{
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
			formatter: JsonFormatter{
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
			formatter: JsonFormatter{
				err:        errors.New("main error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				subErrors: []error{
					errors.New("sub error 1"),
					&FaultError{
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
			formatter: JsonFormatter{
				err:        nil,
				stacktrace: nil,
				when:       nil,
				requestId:  "",
			},
			expected: `{"type":"none","message":"","stacktrace":[]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.formatter.Format()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestTextFormatter_Format(t *testing.T) {
	testCases := []struct {
		label     string
		formatter TextFormatter
		expected  string
	}{
		{
			label: "required fields",
			formatter: TextFormatter{
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
			expected: `[Type:none] [Error:test error]
[StackTraces:
 | example.go:10 main.exampleFunction
 | example.go:20 main.anotherFunction
]
----end
`,
		},
		{
			label: "with when and requestId",
			formatter: TextFormatter{
				errorType:  ErrorType("testType"),
				err:        errors.New("another error"),
				stacktrace: make(StackTrace, 0),
				when: func() *time.Time {
					t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
					return &t
				}(),
				requestId: "req-12345",
			},
			expected: `[Type:testType] [Error:another error] [When:2024-01-01T12:00:00Z] [RequestId:req-12345]
----end
`,
		},
		{
			label: "tags",
			formatter: TextFormatter{
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
			expected: `[Type:none] [Error:error with tags]
[Tags:
 | key:session_id value:sess-456
 | key:user_id value:user-789
]
----end
`,
		},
		{
			label: "empty tags",
			formatter: TextFormatter{
				err:        errors.New("error with empty tags"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				tags: Tags{
					tags:   []Tag{},
					keyMap: map[string]int{},
				},
			},
			expected: `[Type:none] [Error:error with empty tags]
----end
`,
		},
		{
			label: "sub errors",
			formatter: TextFormatter{
				err:        errors.New("main error"),
				stacktrace: make(StackTrace, 0),
				when:       nil,
				requestId:  "",
				subErrors: []error{
					errors.New("sub error 1"),
					&FaultError{
						err:       errors.New("sub error 2"),
						errorType: ErrorType("testType2"),
						stacktrace: StackTrace{
							{File: "sub_example.go", Line: 30, Function: "subFunction"},
						},
					},
				},
			},
			expected: `[Type:none] [Error:main error]
----end
[Type:none] [Error:sub error 1]
----end
[Type:testType2] [Error:sub error 2]
[StackTraces:
 | sub_example.go:30 subFunction
]
----end
`,
		},
		{
			label: "empty",
			formatter: TextFormatter{
				err:        nil,
				stacktrace: nil,
				when:       nil,
				requestId:  "",
			},
			expected: `[Type:none] [Error:<no error>]
----end
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.formatter.Format()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
