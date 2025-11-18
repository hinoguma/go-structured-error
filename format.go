package fault

import "time"

type ErrorFormatter interface {
	Format() string
}

type JsonFormatter struct {
	// required
	faultType  FaultType
	err        error
	stacktrace StackTrace

	// optional
	when      *time.Time
	requestId string
	tags      Tags
	subErrors []error
}

func (f JsonFormatter) Format() string {
	jsonStr := "{"
	jsonStr += `"type":"` + f.faultType.StringWithDefaultNone() + `"`
	if f.err == nil {
		jsonStr += `,"message":""`
	} else {
		jsonStr += `,"message":"` + f.err.Error() + `"`
	}
	if f.when != nil {
		jsonStr += `,"when":"` + f.when.Format(time.RFC3339) + `"`
	}
	if f.requestId != "" {
		jsonStr += `,"request_id":"` + f.requestId + `"`
	}

	if len(f.tags.tags) > 0 {
		jsonStr += `,"tags":` + f.tags.JsonValueString()
	}
	if len(f.stacktrace) == 0 {
		jsonStr += `,"stacktrace":[]`
	} else {
		jsonStr += `,"stacktrace":` + f.stacktrace.JsonValueString()
	}
	jsonStr += "}"
	return jsonStr
}
