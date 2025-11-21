package errors

import "github.com/hinoguma/go-fault"

func ToJsonString(err error) string {
	if err != nil {
		js, ok := err.(fault.JsonStringer)
		if ok {
			return js.JsonString()
		}
		err, ok := err.(interface{ JsonFormatter() fault.ErrorFormatter })
		if ok {
			return err.JsonFormatter().Format()
		}
	}
	fe := fault.NewRawFaultError(err)
	return fe.JsonString()
}

func ToFaultError(err error) *fault.FaultError {
	if err == nil {
		return fault.NewRawFaultError(err)
	}
	fe, ok := err.(*fault.FaultError)
	if !ok {
		fe = fault.NewRawFaultError(err)
		return fe
	}
	return fe
}

func ToFault(err error) fault.Fault {
	if err == nil {
		return fault.NewRawFaultError(err)
	}
	fe, ok := err.(fault.Fault)
	if !ok {
		fe = fault.NewRawFaultError(err)
		return fe
	}
	return fe
}
