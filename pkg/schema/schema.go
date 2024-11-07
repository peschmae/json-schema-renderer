package util

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

var objects = make(map[string]jsonschema.Schema)
var requiredObjectNames = make(map[string]bool)

// Loop through all properties and gather all objects. Additionally a map of all requiredObjects is also returned
func GatherObjects(schema *jsonschema.Schema, flatObjects []string) (map[string]jsonschema.Schema, map[string]bool) {
	rootPropertyKeys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		rootPropertyKeys = append(rootPropertyKeys, k)
	}
	slices.Sort(rootPropertyKeys)

	for _, key := range rootPropertyKeys {
		gatherObjects("", key, schema.Properties[key], flatObjects)
	}

	return objects, requiredObjectNames
}

func gatherObjects(parentTitle, name string, schema *jsonschema.Schema, flatObjects []string) bool {

	if isRequired(schema) {
		requiredObjectNames[name] = true
	}

	// properties matching flatObjects will be dumped directly and shouldn't be added to the list of objects
	if slices.Contains(flatObjects, name) {
		return isRequired(schema)
	}

	if parentTitle != "" {
		name = strings.Join([]string{parentTitle, name}, " > ")
		if isRequired(schema) {
			requiredObjectNames[parentTitle] = true
		}
	}

	childRequired := isRequired(schema)
	// primitive types don't need to be nested
	if schema.Types.String() != "[object]" {
		return isRequired(schema)
	} else {
		objects[name] = *schema

		if len(schema.Required) > 0 {
			requiredObjectNames[name] = true
			childRequired = true
			for _, child := range schema.Required {
				requiredObjectNames[name+" > "+child] = true
			}
		}

		for key, sch := range schema.Properties {
			if sch.Types.String() == "[object]" {
				if gatherObjects(name, key, sch, flatObjects) {
					childRequired = true
				}
			} else {
				if isRequired(sch) {
					childRequired = true
					requiredObjectNames[name+" > "+key] = true
				}
			}
		}
	}

	if childRequired {
		requiredObjectNames[name] = true
	}

	return childRequired
}

func isRequired(schema *jsonschema.Schema) bool {
	if schema.MinProperties != nil && *schema.MinProperties > 0 {
		return true
	}
	if schema.MinLength != nil && *schema.MinLength > 0 {
		return true
	}
	if schema.MinItems != nil && *schema.MinItems > 0 {
		return true
	}
	if schema.MinContains != nil && *schema.MinContains > 0 {
		return true
	}

	return slices.Contains(schema.Required, schema.Title)
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
