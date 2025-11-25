package serrors

import (
	"reflect"
	"strings"
	"testing"
)

func assertEqualsStackTraceItem(t *testing.T, got, expected StackTraceItem, filterPrefix string) {
	// only check traces from this package
	// runtime and file system depends on environment
	if !strings.HasPrefix(got.Function, filterPrefix) {
		return
	}
	//if got.File != expected.File {
	//	t.Errorf("expected file %v, got %v", expected.File, got.File)
	//}
	//if got.Line != expected.Line {
	//	t.Errorf("expected line %v, got %v expected:%v got:%v", expected.Line, got.Line, expected, got)
	//}
	if got.Function != expected.Function {
		t.Errorf("expected function %v, got %v expected:%v got:%v", expected.Function, got.Function, expected, got)
	}
}

func assertEqualsStackTrace(t *testing.T, got, expected StackTrace, filterPrefix string) {
	if len(got) != len(expected) {
		t.Errorf("expected stack trace length %v, got %v expected:%v got :%v", len(expected), len(got), expected, got)
		return
	}
	for i := range got {
		assertEqualsStackTraceItem(t, got[i], expected[i], filterPrefix)
	}
}

func assertStructuredError(t *testing.T, got, expected *StructuredError) {
	if !reflect.DeepEqual(got.err, expected.err) {
		t.Errorf("expected err %v, got %v", expected.err, got.err)
	}
	if got.errorType != expected.errorType {
		t.Errorf("expected errorType %v, got %v", expected.errorType, got.errorType)
	}
	if got.when == nil || expected.when == nil {
		if got.when != expected.when {
			t.Errorf("expected when %v, got %v", expected.when, got.when)
		}
	} else {
		if !got.when.Equal(*expected.when) {
			t.Errorf("expected when %v, got %v", *expected.when, *got.when)
		}
	}
	if got.requestId != expected.requestId {
		t.Errorf("expected requestId %v, got %v", expected.requestId, got.requestId)
	}
	assertEqualsStackTrace(t, got.stacktrace, expected.stacktrace, "github.com/hinoguma/go-structured-error.")
	assertEqualsTags(t, got.tags, expected.tags)
	if len(got.subErrors) != len(expected.subErrors) {
		t.Errorf("expected subErrors length %v, got %v", len(expected.subErrors), len(got.subErrors))
	} else {
		for i := range got.subErrors {
			if got.subErrors[i] != expected.subErrors[i] {
				t.Errorf("expected subError %v, got %v", expected.subErrors[i], got.subErrors[i])
			}
		}
	}
}

func assertStructuredErrorWithErrorValue(t *testing.T, got, expected *StructuredError) {
	if got.err == nil || expected.err == nil {
		if got.err != expected.err {
			t.Errorf("expected err %v, got %v", expected.err, got.err)
		}
	} else {
		if got.err.Error() != expected.err.Error() {
			t.Errorf("expected err %v, got %v", expected.err, got.err)
		}
	}
	if got.errorType != expected.errorType {
		t.Errorf("expected errorType %v, got %v", expected.errorType, got.errorType)
	}
	if got.when == nil || expected.when == nil {
		if got.when != expected.when {
			t.Errorf("expected when %v, got %v", expected.when, got.when)
		}
	} else {
		if !got.when.Equal(*expected.when) {
			t.Errorf("expected when %v, got %v", *expected.when, *got.when)
		}
	}
	if got.requestId != expected.requestId {
		t.Errorf("expected requestId %v, got %v", expected.requestId, got.requestId)
	}
	assertEqualsStackTrace(t, got.stacktrace, expected.stacktrace, "github.com/hinoguma/go-structured-error.")
	assertEqualsTags(t, got.tags, expected.tags)
	if len(got.subErrors) != len(expected.subErrors) {
		t.Errorf("expected subErrors length %v, got %v", len(expected.subErrors), len(got.subErrors))
	} else {
		for i := range got.subErrors {
			gotSubErr := got.subErrors[i]
			expectedSubErr := expected.subErrors[i]
			if gotSubErr == nil || expectedSubErr == nil {
				if gotSubErr != expectedSubErr {
					t.Errorf("expected subError %v, got %v", expectedSubErr, gotSubErr)
				}
			} else {
				if gotSubErr.Error() != expectedSubErr.Error() {
					t.Errorf("expected subError %v, got %v", expectedSubErr, gotSubErr)
				}
			}
		}
	}
}

func assertEqualsStructuredWithoutStackTrace(t *testing.T, got, expected Structured) {
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
		unwrapGotFe, okGot := unwrapGot.(Structured)
		unwrapExpectedFe, okExpected := unwrapExpected.(Structured)
		if okGot && okExpected {
			assertEqualsStructuredWithoutStackTrace(t, unwrapGotFe, unwrapExpectedFe)
		} else {
			if !reflect.DeepEqual(unwrapGot, unwrapExpected) {
				t.Errorf("expected unwrapped error %v, got %v", unwrapExpected, unwrapGot)
			}
		}
	}
}
