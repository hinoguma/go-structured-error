package fault

import (
	"fmt"
)

// todo: implemant TagValue types for map, slice , array, struct
type TagValue interface {
	String() string
	JsonValueString() string
}

type StringTagValue string

func (v StringTagValue) String() string {
	return string(v)
}

func (v StringTagValue) JsonValueString() string {
	return "\"" + string(v) + "\""
}

type IntTagValue int

func (v IntTagValue) String() string {
	return fmt.Sprintf("%d", v)
}

func (v IntTagValue) JsonValueString() string {
	return v.String()
}

type BoolTagValue bool

func (v BoolTagValue) String() string {
	if v {
		return "true"
	}
	return "false"
}

func (v BoolTagValue) JsonValueString() string {
	return v.String()
}

type FloatTagValue float64

func (v FloatTagValue) String() string {
	return fmt.Sprintf("%g", v)
}

func (v FloatTagValue) JsonValueString() string {
	return v.String()
}

type NilTagValue struct{}

func (v NilTagValue) String() string {
	return "null"
}

func (v NilTagValue) JsonValueString() string {
	return v.String()
}
