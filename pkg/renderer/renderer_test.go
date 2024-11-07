package renderer

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
)

func TestGetValueNone(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	v := GetValue(*schema.Properties["foo"].Properties["map_key"])

	assert.Equal(t, "", v)
}

func TestGetValueInteger(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	v := GetValue(*schema.Properties["foo"].Properties["max_key"])

	assert.Equal(t, "10", v)
}

func TestGetValueNumber(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	v := GetValue(*schema.Properties["foo"].Properties["range_float_key"])

	assert.Equal(t, "2.2", v)
}

func TestGetValueString(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	v := GetValue(*schema.Properties["foo"].Properties["one_of_strings"])

	assert.Equal(t, "one", v)
}

func TestGetValueArray(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	v := GetValue(*schema.Properties["foo"].Properties["array_key"])

	assert.Equal(t, "[]", v)
}

func TestGetValueBool(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	v := GetValue(*schema.Properties["foo"].Properties["bool"])

	assert.Equal(t, "true", v)
}
