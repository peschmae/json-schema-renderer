package schema

import (
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

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

func dumpPropertiesToMap(properties map[string]*jsonschema.Schema) map[string]any {

	props := make(map[string]any)
	for k, v := range properties {
		if v.Types.String() == "[object]" {
			props[k] = dumpPropertiesToMap(v.Properties)
		} else {
			props[k] = v.Default
		}
	}

	return props
}
