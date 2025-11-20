package errors

import "github.com/hinoguma/go-fault"

func ToFaultError(err error) *fault.FaultError {
	if err == nil {
		return fault.NewRawFaultError(nil)
	}
	fe, ok := err.(*fault.FaultError)
	if !ok {
		fe = fault.NewRawFaultError(err)
		return fe
	}
	return fe
}

func IsType(err error, t fault.ErrorType) bool {
	return fault.IsType(err, t)
}
