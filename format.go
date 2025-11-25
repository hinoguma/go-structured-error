package go_fault

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type ErrorFormatter interface {
	Format() string
}

const NoErrStr string = "<no error>"
const indentation string = "    "

type JsonPrinter interface {
	Print() string
}

type ErrorJsonPrinter struct {
	// required
	errorType  ErrorType
	err        error
	stacktrace StackTrace

	// optional
	when      *time.Time
	requestId string
	tags      Tags
	subErrors []error
}

func (f ErrorJsonPrinter) Print() string {
	jsonStr := `{"type":"` + f.errorType.StringWithDefaultNone() + `"`
	if f.err == nil {
		jsonStr += `,"message":"` + NoErrStr + `"`
	} else {
		escaped, ok := json.Marshal(f.err.Error())
		if ok != nil {
			escaped = []byte(`"` + f.err.Error() + `"`)
		}
		jsonStr += `,"message":` + string(escaped)
	}
	if f.when != nil {
		jsonStr += `,"when":"` + f.when.Format(time.RFC3339) + `"`
	}
	if f.requestId != "" {
		escaped, ok := json.Marshal(f.requestId)
		if ok != nil {
			escaped = []byte(`"` + f.requestId + `"`)
		}
		jsonStr += `,"request_id":` + string(escaped)
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
			var jf JsonPrinter
			fe, ok := subErr.(HasJsonPrinter)
			if ok {
				jf = fe.JsonPrinter()
			} else {
				jf = ErrorJsonPrinter{
					errorType: ErrorTypeNone,
					err:       subErr,
				}
			}
			if i > 0 {
				jsonStr += `,`
			}
			jsonStr += jf.Print()
		}
		jsonStr += `]`
	}
	jsonStr += "}"
	return jsonStr
}

type VerboseFormatter struct {
	// required
	title      string
	errorType  ErrorType
	err        error
	stacktrace StackTrace

	// optional
	when      *time.Time
	requestId string
	tags      Tags
	subErrors []error
}

func (f VerboseFormatter) Format() string {

	txt := ""
	txt += f.formatMain()

	if len(f.subErrors) > 0 {
		for i, subErr := range f.subErrors {
			if subErr == nil {
				continue
			}
			fe, ok := subErr.(interface{ VerboseFormatter() ErrorFormatter })
			var subFormatter VerboseFormatter
			if ok {
				subFormatter = fe.VerboseFormatter().(VerboseFormatter)
			} else {
				subFormatter = VerboseFormatter{
					errorType: ErrorTypeNone,
					err:       subErr,
				}
			}
			subFormatter.title = f.title + ".sub" + strconv.Itoa(i+1)
			txt += "\n" + subFormatter.Format()
		}
	}
	return txt
}

func (f VerboseFormatter) formatMain() string {
	txt := ""
	if f.err == nil {
		txt += "\n" + "message: " + NoErrStr
	} else {
		txt += "\n" + "message: " + f.err.Error()
	}
	txt += "\n" + "type: " + f.errorType.StringWithDefaultNone()

	if f.when != nil {
		txt += "\n" + "when: " + f.when.Format(time.RFC3339)
	}
	if f.requestId != "" {
		txt += "\n" + "request_id: " + f.requestId
	}

	// tags
	if len(f.tags.tags) > 0 {
		txt += "\n" + "tags:"
		for _, tag := range f.tags.tags {
			txt += "\n" + indentation + tag.Key + ": " + tag.Value.String()
		}
	}

	if len(f.stacktrace) > 0 {
		txt += "\n" + "stacktrace:"
		for _, frame := range f.stacktrace {
			txt += "\n" + indentation + frame.String()
		}
	}

	txt = strings.ReplaceAll(txt, "\n", "\n"+indentation)
	txt = f.title + ":" + txt
	return txt
}
