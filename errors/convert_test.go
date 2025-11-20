package errors

import (
	"github.com/hinoguma/go-fault"
	"reflect"
	"testing"
)

func TestNewConverter(t *testing.T) {
	testCases := []struct {
		label    string
		err      error
		expected Converter
	}{
		{
			label:    "NewConverter creates a non-nil Converter",
			err:      errStd,
			expected: Converter{err: errStd},
		},
		{
			label:    "NewConverter with nil error",
			err:      nil,
			expected: Converter{err: nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			converter := NewConverter(tc.err)
			if !reflect.DeepEqual(converter, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, converter)
			}
		})
	}
}

func TestConverter_JsonString(t *testing.T) {
	// The format of Json follows FaultError's JsonString output
	// test cases of each props are implemented in fault package tests
	testCases := []struct {
		label     string
		converter Converter
		expected  string
	}{
		{
			label:     "nil",
			converter: Converter{err: nil},
			expected:  `{"type":"none","message":"","stacktrace":[]}`,
		},
		{
			label:     "fault error",
			converter: Converter{err: fault.NewRawFaultError(errStd)},
			expected:  `{"type":"none","message":"standard error","stacktrace":[]}`,
		},
		{
			label:     "standard error",
			converter: Converter{err: errStd},
			expected:  `{"type":"none","message":"standard error","stacktrace":[]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			jsonStr := tc.converter.JsonString()
			if jsonStr != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, jsonStr)
			}
		})
	}

}
