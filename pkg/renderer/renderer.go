package renderer

import (
	"encoding/json"
	"strconv"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

type Renderer interface {
	SetFlatOutput(string)
	Header(string, int) string
	PropertyHeader(string, int) string
	TableHeader() string
	TableFooter() string
	PropertyRow(string, string, jsonschema.Schema, bool) string
	TextParagraph(string) string
}

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

func DumpPropertiesToJson(properties map[string]*jsonschema.Schema) string {

	props := dumpPropertiesToMap(properties)

	b, _ := json.MarshalIndent(props, "", " ")
	return string(b)
}

func DumpPropertiesToYaml(properties map[string]*jsonschema.Schema) string {

	props := dumpPropertiesToMap(properties)

	b, _ := yaml.Marshal(props)
	return string(b)
}

func dumpPropertiesToMap(properties map[string]*jsonschema.Schema) map[string]interface{} {

	props := map[string]interface{}{}
	for k, v := range properties {
		if v.Types.String() == "[object]" {
			props[k] = dumpPropertiesToMap(v.Properties)
		} else {
			props[k] = v.Default
		}
	}

	return props
}
