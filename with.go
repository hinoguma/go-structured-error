package serrors

import (
	"time"
)

func With(err error, options ...WithFunc) error {
	if err == nil {
		err = ToStructuredError(err)
	}
	for _, opt := range options {
		err = opt(err)
	}
	return err
}

type WithFunc func(err error) error

func WithRequestID(id string) WithFunc {
	return func(err error) error {
		fe := ToStructuredError(err)
		if fe == nil {
			return nil
		}
		_ = fe.SetRequestID(id)
		return fe
	}
}

func WithWhen(t time.Time) WithFunc {
	return func(err error) error {
		fe := ToStructuredError(err)
		if fe == nil {
			return nil
		}
		_ = fe.SetWhen(t)
		return fe
	}
}

func WithType(t ErrorType) WithFunc {
	return func(err error) error {
		fe := ToStructuredError(err)
		if fe == nil {
			return nil
		}
		_ = fe.SetType(t)
		return fe
	}
}

func WithTagSafe(key string, value TagValue) WithFunc {
	return func(err error) error {
		fe := ToStructuredError(err)
		if fe == nil {
			return nil
		}
		_ = fe.AddTagSafe(key, value)
		return fe
	}
}
