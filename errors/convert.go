package errors

import "github.com/hinoguma/go-fault"

func NewConverter(err error) Converter {
	return Converter{err: err}
}

type Converter struct {
	err error
}

func (c Converter) JsonString() string {
	if c.err != nil {
		err, ok := c.err.(interface{ JsonFormatter() fault.ErrorFormatter })
		if ok {
			return err.JsonFormatter().Format()
		}
	}
	fe := fault.NewRawFaultError(c.err)
	return fe.JsonString()
}
