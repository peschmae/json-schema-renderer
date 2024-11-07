package schema

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
)

func TestDumpPropertiesToJson(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	json := DumpPropertiesToJson(schema.Properties)

	expected := `{
 "foo": {
  "array_key": [],
  "map_key": {},
  "max_key": 10,
  "min_key": 10,
  "one_of_integers": 1,
  "one_of_mixed": "one",
  "one_of_strings": "one",
  "range_float_key": 2.2,
  "range_key": 10,
  "string_key": ""
 }
}`

	assert.Equal(t, expected, json)
}

func TestDumpPropertiesToYaml(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	json := DumpPropertiesToYaml(schema.Properties)

	expected := `foo:
    array_key: []
    map_key: {}
    max_key: "10"
    min_key: "10"
    one_of_integers: "1"
    one_of_mixed: one
    one_of_strings: one
    range_float_key: "2.2"
    range_key: "10"
    string_key: ""
`

	assert.Equal(t, expected, json)
}
