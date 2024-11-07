package renderer

import (
	"encoding/json"
	"strconv"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Renderer interface {
	Header(string, int) string
	PropertyHeader(string, int) string
	TableHeader() string
	TableFooter() string
	PropertyRow(string, string, jsonschema.Schema, bool) string
	TextParagraph(string) string
	HeaderLevel(int) int
}

// in JSON all values are basically strings, so they are converted before returned
func GetValue(schema jsonschema.Schema) string {
	if schema.Types.String() != "[object]" {

		if schema.Default == nil {
			return ""
		}

		var value string
		switch schema.Types.String() {
		case "[string]":
			value = (*schema.Default).(string)
		case "[number]":
			value = (*schema.Default).(json.Number).String()
		case "[integer]":
			value = (*schema.Default).(json.Number).String()
		case "[boolean]":
			value = strconv.FormatBool((*schema.Default).(bool))
		case "[array]":
			if schema.Default != nil {
				b, _ := json.Marshal((*schema.Default).([]interface{}))
				value = string(b)
			} else if schema.Items != nil {
				s := (schema.Items).(*jsonschema.Schema)
				if len(s.AnyOf) > 0 {
					value = (schema.Items).(*jsonschema.Schema).AnyOf[0].Types.String()
				} else {
					value = "[]"
				}

			} else {
				value = "[]"
			}

		default:
			value = "unknown"
		}

		return value
	}

	return ""
}
