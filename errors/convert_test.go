package errors

import (
	"github.com/hinoguma/go-fault"
	"testing"
)

func TestToJsonString(t *testing.T) {
	// The format of Json follows FaultError's JsonString output
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
			err:      fault.NewRawFaultError(errStd),
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

func TestToFaultError(t *testing.T) {
	testCases := []struct {
		label string
		err   error
		fe    *fault.FaultError
	}{
		{
			label: "nil",
			err:   nil,
			fe:    fault.NewRawFaultError(nil),
		},
		{
			label: "standard error",
			err:   errStd,
			fe:    fault.NewRawFaultError(errStd),
		},
		{
			label: "fault error",
			err:   fault.NewRawFaultError(errStd),
			fe:    fault.NewRawFaultError(errStd),
		},
		{
			label: "fault error with custom type",
			err:   errC3,
			fe:    fault.NewRawFaultError(errC3),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			fe := ToFaultError(tc.err)
			assertEqualsFaultWithoutStackTrace(t, tc.fe, fe)
		})
	}
}

func TestToFault(t *testing.T) {
	testCases := []struct {
		label string
		err   error
		fe    fault.Fault
	}{
		{
			label: "nil",
			err:   nil,
			fe:    fault.NewRawFaultError(nil),
		},
		{
			label: "standard error",
			err:   errStd,
			fe:    fault.NewRawFaultError(errStd),
		},
		{
			label: "fault error",
			err:   fault.NewRawFaultError(errStd),
			fe:    fault.NewRawFaultError(errStd),
		},
		{
			label: "fault error with custom type",
			err:   errC3,
			fe:    errC3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			fe := ToFault(tc.err)
			assertEqualsFaultWithoutStackTrace(t, tc.fe, fe)
		})
	}
}
