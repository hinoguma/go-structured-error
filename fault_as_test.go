package fault

import (
	"errors"
	"testing"
)

func TestFaultError_As(t *testing.T) {

	ne := errors.New("go standard error")
	ca := newTestCustomNonFaultError("custom non-fault error", 300)
	cb := newTestCustomFaultError1(100)

	/**
	Normal E = NE = error
	Custom A = CA = custom error not implementing ExtendError
	Custom B = CB = custom error implementing ExtendError
	Fault E  = FE = FaultError

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
				fe := &FaultError{err: ne}
				var target *FaultError
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
				fe := &FaultError{err: ne}
				var target *testCustomNonFaultError
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
				fe := &FaultError{err: ne}
				var target *testCustomFaultError1
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
				fe := &FaultError{err: ca}
				var target *testCustomNonFaultError
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
				fe := &FaultError{err: ca}
				var target *testCustomFaultError1
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
				fe := &FaultError{err: cb}
				var target *testCustomFaultError1
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
				fe := &FaultError{err: cb}
				var target *testCustomNonFaultError
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
				err := newTestCustomFaultError1(200)
				err.SetErr(ne)
				var target *testCustomFaultError1
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
				err := newTestCustomFaultError1(200)
				err.SetErr(ca)
				var target *testCustomNonFaultError
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
				err := newTestCustomFaultError1(200)
				err.SetErr(ca)
				var target *FaultError
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
				err := newTestCustomFaultError1(200)
				err.SetErr(&FaultError{err: ne})
				var target *FaultError
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
