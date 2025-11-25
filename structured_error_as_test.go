package serrors

import (
	"errors"
	"testing"
)

func TestStructuredError_As(t *testing.T) {

	ne := errors.New("go standard error")
	ca := newTestCustomNonStructuredError("custom non-fault error", 300)
	cb := newTestCustomStructuredError1(100)

	/**
	Normal E = NE = error
	Custom A = CA = custom error not implementing ExtendError
	Custom B = CB = custom error implementing ExtendError
	Fault E  = FE = StructuredError

	Scenarios:
	1. FE has NE as FE -> true
	2. FE has NE as CA -> false
	3. FE has NE as CB -> false

	4. FE has CA as CA -> true
	5. FE has CA as CB -> false

	6. FE has CB as CB -> true
	7. FE has CB as CA -> false

	8. CB has NE as CB -> true
	9. CB has CA as CA -> true
	10. CB has CA as FE -> false
	11. CB has FE as FE -> true

	*/

	testCases := []struct {
		label        string
		scenarioFunc func(t *testing.T)
	}{
		{
			label: "1. FE has NE as FE -> true",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: ne}
				var target *StructuredError
				expectedOk := true

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "2. FE has NE as CA -> false",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: ne}
				var target *testCustomNonStructuredError
				expectedOk := false

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "3. FE has NE as CB -> false",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: ne}
				var target *testCustomStructuredError1
				expectedOk := false

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "4. FE has CA as CA -> true",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: ca}
				var target *testCustomNonStructuredError
				expectedOk := true

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "5. FE has CA as CB -> false",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: ca}
				var target *testCustomStructuredError1
				expectedOk := false

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "6. FE has CB as CB -> true",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: cb}
				var target *testCustomStructuredError1
				expectedOk := true

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "7. FE has CB as CA -> false",
			scenarioFunc: func(t *testing.T) {
				// prepare
				fe := &StructuredError{err: cb}
				var target *testCustomNonStructuredError
				expectedOk := false

				// execute
				ok := errors.As(fe, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "8. CB has NE as CB -> true",
			scenarioFunc: func(t *testing.T) {
				// prepare
				err := newTestCustomStructuredError1(200)
				_ = err.SetErr(ne)
				var target *testCustomStructuredError1
				expectedOk := true

				// execute
				ok := errors.As(err, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "9. CB has CA as CA -> true",
			scenarioFunc: func(t *testing.T) {
				// prepare
				err := newTestCustomStructuredError1(200)
				_ = err.SetErr(ca)
				var target *testCustomNonStructuredError
				expectedOk := true

				// execute
				ok := errors.As(err, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "10. CB has CA as FE -> false",
			scenarioFunc: func(t *testing.T) {
				// prepare
				err := newTestCustomStructuredError1(200)
				_ = err.SetErr(ca)
				var target *StructuredError
				expectedOk := false

				// execute
				ok := errors.As(err, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
		{
			label: "11. CB has FE as FE -> true",
			scenarioFunc: func(t *testing.T) {
				// prepare
				err := newTestCustomStructuredError1(200)
				_ = err.SetErr(&StructuredError{err: ne})
				var target *StructuredError
				expectedOk := true

				// execute
				ok := errors.As(err, &target)

				// verify
				if ok != expectedOk {
					t.Errorf("expected As to return %v, got %v", expectedOk, ok)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			tc.scenarioFunc(t)
		})
	}
}
