package serrors

import (
	"errors"
	"reflect"
	"testing"
)

func TestToJsonString(t *testing.T) {
	// The format of Json follows StructuredError's JsonString output
	// test cases of each props are implemented in fault package tests
	testCases := []struct {
		label    string
		err      error
		expected string
	}{
		{
			label:    "nil",
			err:      nil,
			expected: `{"type":"none","message":"<no error>","stacktrace":[]}`,
		},
		{
			label:    "fault error",
			err:      NewRawStructuredError(errStd),
			expected: `{"type":"none","message":"standard error","stacktrace":[]}`,
		},
		{
			label:    "standard error",
			err:      errStd,
			expected: `{"type":"none","message":"standard error","stacktrace":[]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			jsonStr := ToJsonString(tc.err)
			if jsonStr != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, jsonStr)
			}
		})
	}
}

func TestToStructuredError(t *testing.T) {
	testCases := []struct {
		label string
		err   error
		fe    *StructuredError
	}{
		{
			label: "nil",
			err:   nil,
			fe:    NewRawStructuredError(nil),
		},
		{
			label: "standard error",
			err:   errStd,
			fe:    NewRawStructuredError(errStd),
		},
		{
			label: "fault error",
			err:   NewRawStructuredError(errStd),
			fe:    NewRawStructuredError(errStd),
		},
		{
			label: "fault error with custom type",
			err:   errC3,
			fe:    NewRawStructuredError(errC3),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			fe := ToStructuredError(tc.err)
			assertEqualsStructuredWithoutStackTrace(t, tc.fe, fe)
		})
	}
}

func TestToStructured(t *testing.T) {
	stdErr := errors.New("standard error")
	testCases := []struct {
		label    string
		err      error
		expected SError
	}{
		{
			label: "nil error",
			err:   nil,
			expected: &StructuredError{
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
			label: "convert standard error",
			err:   errors.New("standard error"),
			expected: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        stdErr,
				stacktrace: make(StackTrace, 0),
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label: "already a StructuredError",
			err: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        stdErr,
				stacktrace: make(StackTrace, 0),
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
			expected: &StructuredError{
				errorType:  ErrorTypeNone,
				err:        stdErr,
				stacktrace: make(StackTrace, 0),
				tags:       NewTags(),
				subErrors:  make([]error, 0),
			},
		},
		{
			label:    "already a SError interface",
			err:      &testCustomStructuredError1{},
			expected: &testCustomStructuredError1{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := ToStructured(tc.err)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected SError %v, got %v", tc.expected, got)
			}
		})
	}
}
