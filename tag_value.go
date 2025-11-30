package serrors

import (
	"encoding/json"
	"fmt"
)

// TagValue.String is a methof to convert the tag value to a string for verbose output
// TagValue.JsonValueString is a method to convert the tag value to a json string safely
type TagValue interface {
	String() string
	JsonValueString() string
}

type StringTagValue string

func (v StringTagValue) String() string {
	return string(v)
}

func (v StringTagValue) JsonValueString() string {
	// escape the string for JSON
	// line breaks and special characters will be escaped
	s, _ := json.Marshal(v.String())
	return string(s)
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
