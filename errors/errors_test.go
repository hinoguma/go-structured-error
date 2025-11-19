package errors

import (
	"errors"
	"fmt"
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

func TestIs(t *testing.T) {
	stdErr1 := New("standard error")
	stdErr2 := New("another standard error")
	testCases := []struct {
		label  string
		err    error
		target error
		match  bool
	}{
		{
			label:  "same error",
			err:    stdErr1,
			target: stdErr1,
			match:  true,
		},
		{
			label:  "different errors",
			err:    stdErr1,
			target: stdErr2,
			match:  false,
		},
		{
			label:  "wrapped error matches target",
			err:    Wrap(stdErr1, "additional context"),
			target: stdErr1,
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
			target: &target,
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

}

func TestNew(t *testing.T) {

}
