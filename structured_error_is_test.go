package serrors

import (
	"errors"
	"fmt"
	"testing"
)

const testErrorType1 ErrorType = "testCustom1"

func newTestCustomStructuredError1(code int) *testCustomStructuredError1 {
	return &testCustomStructuredError1{
		StructuredError: StructuredError{
			errorType: testErrorType1,
		},
		code: code,
	}
}

func newTestCustomNonStructuredError(message string, code int) *testCustomNonStructuredError {
	return &testCustomNonStructuredError{
		message: message,
		code:    code,
	}
}

type testCustomStructuredError1 struct {
	StructuredError
	code int
}

type testCustomNonStructuredError struct {
	message string
	code    int
}

func (e *testCustomNonStructuredError) Error() string {
	return fmt.Sprintf("non-fault error: code:%d message:%s", e.code, e.message)
}

func TestStructuredError_Is(t *testing.T) {
	ne := errors.New("go standard error")
	ca := newTestCustomNonStructuredError("custom non-fault error", 300)
	cb := newTestCustomStructuredError1(100)

	testCases := []struct {
		label    string
		err      error
		target   error
		expected bool
	}{
		/**
		Normal E = NE = error
		Custom A = CA = custom error not implementing StructuredError
		Custom B = CB = custom error implementing StructuredError
		Fault  E = FE = StructuredError

		Scenarios:
		1. FE has NE is FE has same NE -> true
		2. FE has NE is FE has dif NE  -> false

		3. FE has CA is FE has same CA -> true
		4. FE has CA is FE has dif CA -> false

		5. CB has NE is CB has same NE -> true
		6. CB has NE is CB has dif NE -> false
		7. CB has NE us is FE has same NE -> false
		8. FE has nil target is nil -> false

		9. FE has NE is NE -> true
		*/
		{
			label:    "1. FE has NE is FE has same NE -> true",
			err:      &StructuredError{err: ne},
			target:   &StructuredError{err: ne},
			expected: true,
		},
		{
			label:    "2. FE has NE is FE has dif NE  -> false",
			err:      &StructuredError{err: ne},
			target:   &StructuredError{err: errors.New("another go std error")},
			expected: false,
		},
		{
			label:    "3. FE has CA is FE has same CA -> true",
			err:      &StructuredError{err: ca},
			target:   &StructuredError{err: ca},
			expected: true,
		},
		{
			label:    "4. FE has CA is FE has dif CA -> false",
			err:      &StructuredError{err: cb},
			target:   &StructuredError{err: newTestCustomNonStructuredError("custom non-fault error", 300)},
			expected: false,
		},
		{
			label:    "5. CB has NE is CB has same NE -> true",
			err:      newTestCustomStructuredError1(100).SetErr(ne),
			target:   newTestCustomStructuredError1(100).SetErr(ne),
			expected: true,
		},
		{
			label:    "6. CB has NE is CB has dif NE -> false",
			err:      newTestCustomStructuredError1(100).SetErr(ne),
			target:   newTestCustomStructuredError1(100).SetErr(errors.New("another go std error")),
			expected: false,
		},
		{
			label:    "7. CB has NE is FE has same NE -> false",
			err:      newTestCustomStructuredError1(100).SetErr(ne),
			target:   &StructuredError{err: ne},
			expected: false,
		},
		{
			label:    "8. FE has nil target is nil -> false",
			err:      &StructuredError{err: nil},
			target:   nil,
			expected: false,
		},
		{
			label:    "9. NE is FE has same NE -> true",
			err:      &StructuredError{err: ne},
			target:   ne,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			result := errors.Is(tc.err, tc.target)
			if result != tc.expected {
				t.Errorf("expected Is result %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsType(t *testing.T) {

	testCases := []struct {
		label    string
		err      error
		target   ErrorType
		expected bool
	}{
		{
			label:    "Initial FE is None -> true",
			err:      &StructuredError{},
			target:   ErrorTypeNone,
			expected: true,
		},
		{
			label:    "Initial FE is testCustom1 -> true",
			err:      &StructuredError{},
			target:   testErrorType1,
			expected: false,
		},
		{
			label:    "FE with type testCustom1 -> true",
			err:      &StructuredError{errorType: testErrorType1},
			target:   testErrorType1,
			expected: true,
		},
		{
			label:    "testCustom1 is testCustom1 -> true",
			err:      newTestCustomStructuredError1(100),
			target:   testErrorType1,
			expected: true,
		},
		{
			label:    "testCustom1 is None -> false",
			err:      newTestCustomStructuredError1(100),
			target:   ErrorTypeNone,
			expected: false,
		},
		{
			label:    "non-Fault error is None -> false",
			err:      errors.New("standard go error"),
			target:   ErrorTypeNone,
			expected: false,
		},
		{
			label:    "warped StructuredError with type testCustom1 is None -> false",
			err:      fmt.Errorf("wrapping fault error: %w", newTestCustomStructuredError1(100)),
			target:   ErrorTypeNone,
			expected: false,
		},
		{
			label:    "warped StructuredError with type testCustom1 is testCustom1 -> true",
			err:      fmt.Errorf("wrapping fault error: %w", newTestCustomStructuredError1(100)),
			target:   testErrorType1,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			result := IsType(tc.err, tc.target)
			if result != tc.expected {
				t.Errorf("expected IsType result %v, got %v", tc.expected, result)
			}
		})
	}
}
