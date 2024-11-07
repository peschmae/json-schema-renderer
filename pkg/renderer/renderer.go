package renderer

import (
	"encoding/json"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

var HeaderOffset int

type Renderer interface {
	Header(string, int) string
	PropertyHeader(string, int) string
	TableHeader() string
	TableFooter() string
	PropertyRow(string, string, jsonschema.Schema, bool) string
	TextParagraph(string) string
}

// in JSON all values are basically strings, so they are converted before returned
func getValue(schema jsonschema.Schema) string {
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

func headerLevel(level int) int {
	return int(math.Min(6, float64(level+HeaderOffset)))
}

func RenderDocumentation(r Renderer, objects map[string]jsonschema.Schema, requiredObjectNames map[string]bool, requiredOnly bool, depth int, flatObjects []string) (string, error) {

	output := ""

	// print all schemas
	// sort keys
	objectsKeys := make([]string, 0, len(objects))
	for k := range objects {
		objectsKeys = append(objectsKeys, k)
	}
	slices.Sort(objectsKeys)

	for _, key := range objectsKeys {

		if !requiredOnly || requiredObjectNames[key] {

			if depth == 0 || strings.Count(key, ">") <= depth {
				output += r.PropertyHeader(key, strings.Count(key, ">")+1)

				if objects[key].Title != "" && objects[key].Description != "" {
					output += r.TextParagraph("**" + objects[key].Title + ":** " + objects[key].Description)
				} else if objects[key].Title != "" {
					output += r.TextParagraph(objects[key].Title)
				} else if objects[key].Description != "" {
					output += r.TextParagraph(objects[key].Description)
				}

				output += r.TableHeader()

				propertyKeys := make([]string, 0, len(objects[key].Properties))
				for k := range objects[key].Properties {
					propertyKeys = append(propertyKeys, k)
				}
				slices.Sort(propertyKeys)

				for _, s := range propertyKeys {

					dumpValue := strings.Count(key, ">") == depth
					if slices.Contains(flatObjects, s) {
						dumpValue = true
					} else if depth == 0 {
						// avoid dumping first level objects if depth is 0
						dumpValue = false
					}

					if requiredOnly && !requiredObjectNames[key+" > "+s] { // skip non required objects
						continue
					}

					output += r.PropertyRow(key, s, *objects[key].Properties[s], dumpValue)
				}

				output += r.TableFooter()
				output += "\n"
			}
		}
	}

	return output, nil
}
