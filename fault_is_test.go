package fault

import (
	"errors"
	"fmt"
	"testing"
)

const testErrorType1 FaultType = "testCustom1"

func newTestCustomFaultError1(code int) *testCustomFaultError1 {
	return &testCustomFaultError1{
		FaultError: FaultError{
			faultType: testErrorType1,
		},
		code: code,
	}
}

func newTestCustomNonFaultError(message string, code int) *testCustomNonFaultError {
	return &testCustomNonFaultError{
		message: message,
		code:    code,
	}
}

type testCustomFaultError1 struct {
	FaultError
	code int
}

type testCustomNonFaultError struct {
	message string
	code    int
}

func (e *testCustomNonFaultError) Error() string {
	return fmt.Sprintf("non-fault error: code:%d message:%s", e.code, e.message)
}

func TestFaultError_Is(t *testing.T) {
	ne := errors.New("go standard error")
	ca := newTestCustomNonFaultError("custom non-fault error", 300)
	cb := newTestCustomFaultError1(100)

	testCases := []struct {
		label    string
		err      error
		target   error
		expected bool
	}{
		/**
		Normal E = NE = error
		Custom A = CA = custom error not implementing FaultError
		Custom B = CB = custom error implementing FaultError
		Fault  E = FE = FaultError

		Scenarios:
		1. FE has NE is FE has same NE -> true
		2. FE has NE is FE has dif NE  -> false

		3. FE has CA is FE has same CA -> true
		4. FE has CA is FE has dif CA -> false

		5. CB has NE is CB has same NE -> true
		6. CB has NE is CB has dif NE -> false
		7. CB has NE us is FE has same NE -> false
		8. FE has nil target is nil -> false
		*/
		{
			label:    "1. FE has NE is FE has same NE -> true",
			err:      &FaultError{err: ne},
			target:   &FaultError{err: ne},
			expected: true,
		},
		{
			label:    "2. FE has NE is FE has dif NE  -> false",
			err:      &FaultError{err: ne},
			target:   &FaultError{err: errors.New("another go std error")},
			expected: false,
		},
		{
			label:    "3. FE has CA is FE has same CA -> true",
			err:      &FaultError{err: ca},
			target:   &FaultError{err: ca},
			expected: true,
		},
		{
			label:    "4. FE has CA is FE has dif CA -> false",
			err:      &FaultError{err: cb},
			target:   &FaultError{err: newTestCustomNonFaultError("custom non-fault error", 300)},
			expected: false,
		},
		{
			label:    "5. CB has NE is CB has same NE -> true",
			err:      newTestCustomFaultError1(100).SetErr(ne),
			target:   newTestCustomFaultError1(100).SetErr(ne),
			expected: true,
		},
		{
			label:    "6. CB has NE is CB has dif NE -> false",
			err:      newTestCustomFaultError1(100).SetErr(ne),
			target:   newTestCustomFaultError1(100).SetErr(errors.New("another go std error")),
			expected: false,
		},
		{
			label:    "7. CB has NE is FE has same NE -> false",
			err:      newTestCustomFaultError1(100).SetErr(ne),
			target:   &FaultError{err: ne},
			expected: false,
		},
		{
			label:    "8. FE has nil target is nil -> false",
			err:      &FaultError{err: nil},
			target:   nil,
			expected: false,
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
