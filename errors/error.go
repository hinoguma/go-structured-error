package errors

import (
	stderrors "errors"
	"fmt"
	"github.com/hinoguma/go-fault"
)

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

// compatibility functions for errors.New
func New(text string) error {
	err := fault.NewRawFaultError((stderrors.New(text)))
	// set stack trace starting from caller of New
	err.SetStackTraceWithSkipMaxDepth(2, fault.GetMaxDepthStackTrace())
	return err
}

// compatibility functions for errors.Wrap in pkg/errors, cockroachdb/errors, etc.
// a lot of libraries make Wrap function to wrap errors with message
// Wrap() always return fault.Fault with stack trace
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	fe, ok := err.(fault.Fault)
	if !ok {
		fe = fault.NewRawFaultError(err)
	}
	if len(fe.StackTrace()) == 0 {
		fe.WithStackTrace()
	}
	return fe.SetErr(fmt.Errorf("%s: %w", msg, fe.Unwrap()))
}
