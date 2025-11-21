package errors

import (
	"errors"
	"fmt"
	"github.com/hinoguma/go-fault"
	"reflect"
	"strconv"
	"testing"
)

type testCustomError struct {
	msg string
}

func (e testCustomError) Error() string {
	return e.msg
}

type testCustomError2 struct {
	code int
}

func (e *testCustomError2) Error() string {
	return "error code: " + strconv.Itoa(e.code)
}

const testErrorType3 fault.ErrorType = "testCustomError3"

type testCustomError3 struct {
	fault.FaultError
}

func newTestCustomError3() *testCustomError3 {
	err := &testCustomError3{
		FaultError: fault.FaultError{},
	}
	_ = err.SetType(testErrorType3)
	return err
}

func assertEqualsFaultWithoutStackTrace(t *testing.T, got, expected fault.Fault) {
	if got.Type() != expected.Type() {
		t.Errorf("expected fault type %v, got %v", expected.Type(), got.Type())
	}
	if (got.When() == nil) != (expected.When() == nil) {
		t.Errorf("expected when %v, got %v", expected.When(), got.When())
	} else if got.When() != nil && !got.When().Equal(*expected.When()) {
		t.Errorf("expected when %v, got %v", *expected.When(), *got.When())
	}
	if got.RequestID() != expected.RequestID() {
		t.Errorf("expected request ID %v, got %v", expected.RequestID(), got.RequestID())
	}
	unwrapGot := got.Unwrap()
	unwrapExpected := expected.Unwrap()
	if unwrapGot == nil || unwrapExpected == nil {
		if unwrapGot != unwrapExpected {
			t.Errorf("expected unwrapped error %v, got %v", unwrapExpected, unwrapGot)
		}
	} else {
		unwrapGotFe, okGot := unwrapGot.(fault.Fault)
		unwrapExpectedFe, okExpected := unwrapExpected.(fault.Fault)
		if okGot && okExpected {
			assertEqualsFaultWithoutStackTrace(t, unwrapGotFe, unwrapExpectedFe)
		} else {
			if !reflect.DeepEqual(unwrapGot, unwrapExpected) {
				t.Errorf("expected unwrapped error %v, got %v", unwrapExpected, unwrapGot)
			}
		}
	}
}

func TestIs(t *testing.T) {
	errStd1 := New("standard error")
	errStd2 := New("another standard error")
	testCases := []struct {
		label  string
		err    error
		target error
		match  bool
	}{
		{
			label:  "same error",
			err:    errStd1,
			target: errStd1,
			match:  true,
		},
		{
			label:  "different errors",
			err:    errStd1,
			target: errStd2,
			match:  false,
		},
		{
			label:  "wrapped error matches target",
			err:    Wrap(errStd1, "additional context"),
			target: errStd1,
			match:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			ok := Is(tc.err, tc.target)
			if ok != tc.match {
				t.Errorf("expected Is to return %v, got %v", tc.match, ok)
			}
		})
	}
}

func TestAs(t *testing.T) {
	var target testCustomError
	var target2 *testCustomError2
	testCases := []struct {
		label  string
		err    error
		target any
		match  bool
		want   any
	}{
		{
			label:  "custom error match",
			err:    testCustomError{msg: "custom error occurred"},
			target: &target,
			match:  true,
			want:   testCustomError{msg: "custom error occurred"},
		},
		{
			label:  "custom error2 match",
			err:    &testCustomError2{code: 404},
			target: &target2,
			match:  true,
			want:   &testCustomError2{code: 404},
		},
		{
			label:  "custom error no match",
			err:    &testCustomError2{code: 404},
			target: &target,
			match:  false,
			want:   nil,
		},
		{
			label:  "wrapped custom error match",
			err:    fmt.Errorf("wrapping: %w", testCustomError{msg: "wrapped custom error"}),
			target: &target,
			match:  true,
			want:   testCustomError{msg: "wrapped custom error"},
		},
		{
			label:  "wrapped custom error no match",
			err:    fmt.Errorf("wrapping: %w", &testCustomError2{code: 500}),
			target: &target,
			match:  false,
			want:   nil,
		},
		{
			label:  "nil error",
			err:    nil,
			target: &target,
			match:  false,
			want:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			match := errors.As(tc.err, tc.target)
			if match != tc.match {
				t.Fatalf("expected As to return %v, got %v", tc.match, match)
			}
			if match {
				rtarget := reflect.ValueOf(tc.target)
				got := rtarget.Elem().Interface()
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("expected target to be %v, got %v", tc.want, got)
				}
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errors.New("original error")
	wrappedErr := fmt.Errorf("wrapped: %w", err1)
	testCases := []struct {
		label    string
		err      error
		expected error
	}{
		{
			label:    "simple wrap",
			err:      fmt.Errorf("wrapped: %w", err1),
			expected: err1,
		},
		{
			label:    "no wrap",
			err:      errors.New("no wrap here"),
			expected: nil,
		},
		{
			label:    "nil error",
			err:      nil,
			expected: nil,
		},
		{
			label:    "double wrap",
			err:      fmt.Errorf("double wrapped: %w", wrappedErr),
			expected: wrappedErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := errors.Unwrap(tc.err)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	err1 := errors.New("error one")
	err2 := errors.New("error two")
	testCases := []struct {
		label    string
		target   []error
		expected []error
	}{
		{
			label:    "single error",
			target:   []error{err1},
			expected: []error{err1},
		},
		{
			label:    "two errors",
			target:   []error{err1, err2},
			expected: []error{err1, err2},
		},
		{
			label:    "errors with nil",
			target:   []error{err1, nil, err2, nil},
			expected: []error{err1, err2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := Join(tc.target...).(interface{ Unwrap() []error }).Unwrap()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
			if len(got) != cap(got) {
				t.Errorf("expected len==cap, got len=%v, cap=%v", len(got), cap(got))
			}
		})
	}
}

func TestJoinReturnsNil(t *testing.T) {
	if err := errors.Join(); err != nil {
		t.Errorf("errors.Join() = %v, want nil", err)
	}
	if err := errors.Join(nil); err != nil {
		t.Errorf("errors.Join(nil) = %v, want nil", err)
	}
	if err := errors.Join(nil, nil); err != nil {
		t.Errorf("errors.Join(nil, nil) = %v, want nil", err)
	}
}

func TestWrap(t *testing.T) {
	testCases := []struct {
		label    string
		err      error
		message  string
		expected error
	}{
		{
			label:   "wrap standard error",
			err:     errStd,
			message: "additional context",
			expected: fault.NewRawFaultError(
				fmt.Errorf("additional context: %w", errStd),
			),
		},
		{
			label:    "wrap nil error",
			err:      nil,
			message:  "should be nil",
			expected: nil,
		},
		{
			label:   "wrap fault error",
			err:     fault.New("Original fault error"),
			message: "wrapping fault",
			expected: func() error {
				fe := fault.New("Original fault error")
				return fe.SetErr(fmt.Errorf("wrapping fault: %w", fe.Unwrap()))
			}(),
		},
		{
			label:   "wrap custom error",
			err:     testCustomError{msg: "custom error occurred"},
			message: "wrapping custom error",
			expected: func() error {
				fe := fault.NewRawFaultError(testCustomError{msg: "custom error occurred"})
				_ = fe.WithStackTrace()
				return fe.SetErr(fmt.Errorf("wrapping custom error: %w", fe.Unwrap()))
			}(),
		},
		{
			label: "wrap custom error implementing error interface",
			err: func() error {
				err := &testCustomError3{}
				_ = err.SetErr(errStd)
				return err
			}(),
			message: "wrapping custom error3",
			expected: func() error {
				err := &testCustomError3{}
				_ = err.SetErr(fmt.Errorf("wrapping custom error3: %w", errStd))
				_ = err.WithStackTrace()
				return err
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := Wrap(tc.err, tc.message)
			if tc.expected == nil || got == nil {
				if tc.expected != got {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				return
			}
			expectedFe, expectedOk := tc.expected.(fault.Fault)
			fe, gotOk := got.(fault.Fault)
			if expectedOk != gotOk {
				t.Errorf("expected type fault.Fault: %v, got %v", expectedOk, gotOk)
			}
			if expectedOk && gotOk {
				assertEqualsFaultWithoutStackTrace(t, fe, expectedFe)
			} else {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				return
			}
			if len(fe.StackTrace()) == 0 {
				t.Errorf("expected stack trace to be set, but it was empty")
			}
		})
	}
}

func TestLift(t *testing.T) {
	testCases := []struct {
		label    string
		err      error
		expected error
	}{
		{
			label:    "wrap standard error",
			err:      errStd,
			expected: fault.NewRawFaultError(errStd),
		},
		{
			label:    "wrap nil error",
			err:      nil,
			expected: nil,
		},
		{
			label: "wrap fault error",
			err:   fault.New("Original fault error"),
			expected: func() error {
				fe := fault.New("Original fault error")
				return fe
			}(),
		},
		{
			label: "wrap custom error",
			err:   testCustomError{msg: "custom error occurred"},
			expected: func() error {
				fe := fault.NewRawFaultError(testCustomError{msg: "custom error occurred"})
				return fe
			}(),
		},
		{
			label: "wrap custom error implementing error interface",
			err: func() error {
				err := &testCustomError3{}
				_ = err.SetErr(errStd)
				return err
			}(),
			expected: func() error {
				err := &testCustomError3{}
				_ = err.SetErr(errStd)
				return err
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := Lift(tc.err)
			if tc.expected == nil || got == nil {
				if tc.expected != got {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				return
			}
			expectedFe, expectedOk := tc.expected.(fault.Fault)
			fe, gotOk := got.(fault.Fault)
			if expectedOk != gotOk {
				t.Errorf("expected type fault.Fault: %v, got %v", expectedOk, gotOk)
			}
			if expectedOk && gotOk {
				assertEqualsFaultWithoutStackTrace(t, fe, expectedFe)
			} else {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				return
			}
			if len(fe.StackTrace()) == 0 {
				t.Errorf("expected stack trace to be set, but it was empty")
			}
		})
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		label    string
		message  string
		expected error
	}{
		{
			label:   "basic error",
			message: "this is a test error",
			expected: func() error {
				fe := fault.NewRawFaultError(errors.New("this is a test error"))
				return fe
			}(),
		},
		{
			label:   "empty message",
			message: "",
			expected: func() error {
				fe := fault.NewRawFaultError(errors.New(""))
				return fe
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := New(tc.message)
			if tc.expected == nil || got == nil {
				if tc.expected != got {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				return
			}
			expectedFe, expectedOk := tc.expected.(fault.Fault)
			fe, gotOk := got.(fault.Fault)
			if expectedOk != gotOk {
				t.Errorf("expected type fault.Fault: %v, got %v", expectedOk, gotOk)
			}
			if expectedOk && gotOk {
				assertEqualsFaultWithoutStackTrace(t, fe, expectedFe)
			} else {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				return
			}
			if len(fe.StackTrace()) == 0 {
				t.Errorf("expected stack trace to be set, but it was empty")
			}
		})
	}
}

type causerTestError struct {
	msg   string
	cause error
}

func (e *causerTestError) Error() string {
	return e.msg
}

func (e *causerTestError) Cause() error {
	return e.cause
}

func TestCause(t *testing.T) {
	err1 := errors.New("root cause error")
	wrappedErr := fmt.Errorf("wrapped: %w", err1)
	faultNilErr := fault.NewRawFaultError(nil)
	joinedErr := errors.Join(err1, errors.New("another error"))
	causerErr := &causerTestError{
		msg:   "causer error",
		cause: err1,
	}
	testCases := []struct {
		label    string
		err      error
		expected error
	}{
		{
			label:    "simple wrapped error",
			err:      wrappedErr,
			expected: err1,
		},
		{
			label:    "non-wrapped error",
			err:      err1,
			expected: err1,
		},
		{
			label:    "nil error",
			err:      nil,
			expected: nil,
		},
		{
			label:    "fault error has error inside",
			err:      fault.NewRawFaultError(err1),
			expected: err1,
		},
		{
			label:    "fault error wrapping another error",
			err:      fault.NewRawFaultError(wrappedErr),
			expected: err1,
		},
		{
			label:    "fault error has nil inside",
			err:      faultNilErr,
			expected: faultNilErr,
		},
		{
			label:    "joined errors",
			err:      joinedErr,
			expected: joinedErr,
		},
		{
			label:    "wrapped joined errors",
			err:      fmt.Errorf("wrapping: %w", joinedErr),
			expected: joinedErr,
		},
		{
			label:    "causer error",
			err:      causerErr,
			expected: err1,
		},
		{
			label:    "wrapped causer error",
			err:      fmt.Errorf("wrapping: %w", causerErr),
			expected: err1,
		},
		{
			label: "complex wrapping with causer and fault",
			err: func() error {
				fe := fault.NewRawFaultError(causerErr)
				return fmt.Errorf("outer wrap: %w", fe)
			}(),
			expected: err1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := Cause(tc.err)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestIsType(t *testing.T) {
	// already tested in fault.fault_is_test.go
	testCases := []struct {
		label    string
		err      error
		errType  fault.ErrorType
		expected bool
	}{
		{
			label:    "nil error",
			err:      nil,
			errType:  fault.ErrorTypeNone,
			expected: false,
		},
		{
			label:    "fault error with matching type",
			err:      New("test error"),
			errType:  fault.ErrorTypeNone,
			expected: true,
		},
		{
			label:    "fault error with non-matching type",
			err:      New("test error"),
			errType:  "custom_type",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := IsType(tc.err, tc.errType)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
