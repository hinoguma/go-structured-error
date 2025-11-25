package go_fault

import (
	"errors"
	"fmt"
)

// compatibility functions for errors.New
func New(text string) error {
	err := NewRawStructuredError((errors.New(text)))
	// set stack trace starting from caller of New
	_ = err.SetStackTraceWithSkipMaxDepth(2, MaxStackTraceDepth)
	return err
}

// compatibility functions for errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// compatibility functions for errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// compatibility functions for errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// compatibility functions for errors.Join
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// compatibility functions for errors.Wrap in pkg/errors, cockroachdb/errors, etc.
// a lot of libraries make Wrap function to wrap errors with message
// Wrap() clarifies return type is error interface due to compatibility.
// But Wrap() makes sure the returned error is always Structured interface.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	fe := ToStructured(err)
	if len(fe.StackTrace()) == 0 {
		_ = fe.WithStackTrace()
	}
	return fe.SetErr(fmt.Errorf("%s: %w", msg, fe.Unwrap()))
}

// Lift() is similar to Wrap() but now wrapping with message
// Lift() converts any error to Structured
// if the error is already Structured, it just adds stack trace if missing
func Lift(err error) error {
	if err == nil {
		return nil
	}
	fe := ToStructured(err)
	if len(fe.StackTrace()) == 0 {
		_ = fe.WithStackTrace()
	}
	return fe
}

type causer interface {
	Cause() error
}

// compatibility functions for errors.Cause in pkg/errors
// Cause() returns the underlying cause of the error
func Cause(err error) error {
	if err == nil {
		return nil
	}
	unwrapped := err
	var tmp error
	for {
		if c, ok := unwrapped.(causer); ok {
			tmp = c.Cause()
			if tmp == nil {
				return unwrapped
			}
			unwrapped = tmp
			continue
		}
		tmp = Unwrap(unwrapped)
		if tmp == nil {
			return unwrapped
		}
		unwrapped = tmp
	}
}

// IsTYpe() checks whether the given error or any of its wrapped errors is of the specified ErrorType.
// errors.Is() checks for error equality, but this function checks for error type.
func IsType(err error, t ErrorType) bool {
	if err == nil {
		return false
	}
	fe, ok := err.(HasType)
	if ok && fe.Type() == t {
		return true
	}

	switch x := err.(type) {
	case interface{ Unwrap() error }:
		return IsType(x.Unwrap(), t)
	case interface{ Unwrap() []error }:
		for _, subErr := range x.Unwrap() {
			if IsType(subErr, t) {
				return true
			}
		}
	}
	return false
}
