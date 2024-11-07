package schema

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
)

func TestIsRequiredMinItems(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRequired(schema.Properties["foo"].Properties["array_key"]))
}

func TestIsRequiredMinProperties(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRequired(schema.Properties["foo"].Properties["map_key"]))
}

func TestIsRequiredMinLength(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRequired(schema.Properties["foo"].Properties["string_key"]))
}

func TestIsRequiredFalse(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, isRequired(schema.Properties["foo"].Properties["max_key"]))
}
