package errors

import (
	stderrors "errors"
	"fmt"
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
	return stderrors.New(text)
}

// compatibility functions for errors.Wrap in pkg/errors, cockroachdb/errors, etc.
// a lot of libraries make Wrap function to wrap errors with message
func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
