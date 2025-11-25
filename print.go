package serrors

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const indentation string = "    "
const JsonItemSeparator string = ","

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
	jsonStr := "{"
	jsonStr += BuildJsonStringOfType(f.errorType)
	jsonStr += JsonItemSeparator + BuildJsonStringOfMessage(f.err)

	if f.when != nil {
		jsonStr += JsonItemSeparator + BuildJsonStringOfWhen(*f.when, time.RFC3339)
	}
	if f.requestId != "" {
		jsonStr += JsonItemSeparator + BuildJsonStringOfRequestID(f.requestId)
	}

	if len(f.tags.tags) > 0 {
		jsonStr += JsonItemSeparator + BuildJsonStringOfTags(f.tags)
	}

	jsonStr += JsonItemSeparator + BuildJsonStringOfStackTrace(f.stacktrace)

	if len(f.subErrors) > 0 {
		jsonStr += JsonItemSeparator + BuildJsonStringOfSubErrors(f.subErrors)
	}
	jsonStr += "}"
	return jsonStr
}

func BuildJsonStringOfType(t ErrorType) string {
	return `"type":"` + t.StringWithDefaultNone() + `"`
}

func BuildJsonStringOfMessage(err error) string {
	if err == nil {
		return `"message":"` + NoErrStr + `"`
	}
	escaped, _ := json.Marshal(err.Error())
	return `"message":` + string(escaped)
}

func BuildJsonStringOfWhen(t time.Time, layout string) string {
	return fmt.Sprintf(`"when":"%s"`, t.Format(layout))
}

func BuildJsonStringOfRequestID(requestId string) string {
	escaped, _ := json.Marshal(requestId)
	return `"request_id":` + string(escaped)
}

func BuildJsonStringOfTags(tags Tags) string {
	return `"tags":` + tags.JsonValueString()
}

func BuildJsonStringOfStackTrace(stacktrace StackTrace) string {
	if len(stacktrace) == 0 {
		return `"stacktrace":[]`
	}
	return `"stacktrace":` + stacktrace.JsonValueString()
}

func BuildJsonStringOfSubErrors(subErrors []error) string {
	jsonStr := `"sub_errors":[`
	isFirst := true
	for _, subErr := range subErrors {
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
		if isFirst {
			isFirst = false
		} else {
			jsonStr += `,`
		}
		jsonStr += jf.Print()
	}
	jsonStr += `]`
	return jsonStr
}

type VerbosePrinter interface {
	Print() string
}

type ErrorVerbosePrinter struct {
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

func (f ErrorVerbosePrinter) Print() string {

	txt := ""
	txt += f.printSingle()

	if len(f.subErrors) > 0 {
		for i, subErr := range f.subErrors {
			if subErr == nil {
				continue
			}
			fe, ok := subErr.(interface{ VerbosePrinter() VerbosePrinter })
			var subFormatter ErrorVerbosePrinter
			if ok {
				subFormatter = fe.VerbosePrinter().(ErrorVerbosePrinter)
			} else {
				subFormatter = ErrorVerbosePrinter{
					errorType: ErrorTypeNone,
					err:       subErr,
				}
			}
			subFormatter.title = f.title + ".sub" + strconv.Itoa(i+1)
			txt += "\n" + subFormatter.Print()
		}
	}
	return txt
}

func (f ErrorVerbosePrinter) printSingle() string {
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
