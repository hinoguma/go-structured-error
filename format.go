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
	if len(f.subErrors) > 0 {
		jsonStr += `,"sub_errors":[`
		for i, subErr := range f.subErrors {
			if subErr == nil {
				continue
			}
			var jf ErrorFormatter
			fe, ok := subErr.(interface{ JsonFormatter() ErrorFormatter })
			if ok {
				jf = fe.JsonFormatter()
			} else {
				jf = JsonFormatter{
					faultType: FaultTypeNone,
					err:       subErr,
				}
			}
			if i > 0 {
				jsonStr += `,`
			}
			jsonStr += jf.Format()
		}
		jsonStr += `]`
	}
	jsonStr += "}"
	return jsonStr
}

type TextFormatter struct {
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

func (f TextFormatter) Format() string {
	txt := "[" + "Type:" + f.faultType.StringWithDefaultNone() + "] "
	if f.err == nil {
		txt += "[Error:<no error>]"
	} else {
		txt += "[Error:" + f.err.Error() + "]"
	}
	if f.when != nil {
		txt += " [When:" + f.when.Format(time.RFC3339) + "]"
	}
	if f.requestId != "" {
		txt += " [RequestId:" + f.requestId + "]"
	}
	if len(f.tags.tags) > 0 {
		txt += "\n[Tags:\n"
		for _, value := range f.tags.tags {
			txt += " | " + value.String() + "\n"
		}
		txt += "]"
	}
	if len(f.stacktrace) > 0 {
		txt += "\n[StackTraces:\n"
		for _, item := range f.stacktrace {
			txt += " | " + item.String() + "\n"
		}
		txt += "]"
	}
	txt += "\n----end\n"

	if len(f.subErrors) > 0 {
		for _, subErr := range f.subErrors {
			if subErr == nil {
				continue
			}
			fe, ok := subErr.(interface{ TextFormatter() ErrorFormatter })
			var subFormatter TextFormatter
			if ok {
				subFormatter = fe.TextFormatter().(TextFormatter)
			} else {
				subFormatter = TextFormatter{
					faultType: FaultTypeNone,
					err:       subErr,
				}
			}
			txt += subFormatter.Format()
		}
	}
	return txt
}
