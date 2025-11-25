package go_fault

func ToJsonString(err error) string {
	if err != nil {
		js, ok := err.(JsonStringer)
		if ok {
			return js.JsonString()
		}
		err, ok := err.(interface{ JsonPrinter() JsonPrinter })
		if ok {
			return err.JsonPrinter().Print()
		}
	}
	fe := NewRawStructuredError(err)
	return fe.JsonString()
}

func ToStructuredError(err error) *StructuredError {
	if err == nil {
		return NewRawStructuredError(err)
	}
	fe, ok := err.(*StructuredError)
	if !ok {
		fe = NewRawStructuredError(err)
		return fe
	}
	return fe
}

func ToStructured(err error) Structured {
	if err == nil {
		return NewRawStructuredError(err)
	}
	fe, ok := err.(Structured)
	if !ok {
		fe = NewRawStructuredError(err)
		return fe
	}
	return fe
}
