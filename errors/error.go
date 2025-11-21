package errors

import (
	stderrors "errors"
	"fmt"
	"github.com/hinoguma/go-fault"
)

// compatibility functions for errors.New
func New(text string) error {
	err := fault.NewRawFaultError((stderrors.New(text)))
	// set stack trace starting from caller of New
	_ = err.SetStackTraceWithSkipMaxDepth(2, fault.MaxStackTraceDepth)
	return err
}

// compatibility functions for errors.Is
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// compatibility functions for errors.As
func As(err error, target any) bool {
	return stderrors.As(err, target)
}

// compatibility functions for errors.Unwrap
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// compatibility functions for errors.Join
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}

// compatibility functions for errors.Wrap in pkg/errors, cockroachdb/errors, etc.
// a lot of libraries make Wrap function to wrap errors with message
// Wrap() clarifies return type is error interface due to compatibility.
// But Wrap() makes sure the returned error is always fault.Fault interface.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	fe := ToFault(err)
	if len(fe.StackTrace()) == 0 {
		_ = fe.WithStackTrace()
	}
	return fe.SetErr(fmt.Errorf("%s: %w", msg, fe.Unwrap()))
}

// Lift() is similar to Wrap() but now wrapping with message
// Lift() converts any error to fault.Fault
// if the error is already fault.Fault, it just adds stack trace if missing
func Lift(err error) error {
	if err == nil {
		return nil
	}
	fe := ToFault(err)
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

// This is original function in fault package
// IsType() checks whether the error is of the specified ErrorType
func IsType(err error, t fault.ErrorType) bool {
	return fault.IsType(err, t)
}
