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

func TestGatherObjectsRequiredObjectNames(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	_, req := GatherObjects(schema, []string{})

	assert.NotEmpty(t, req)

	expectedRequired := map[string]bool{
		"foo":              true,
		"foo > array_key":  true,
		"foo > map_key":    true,
		"foo > string_key": true,
		"map_key":          true,
	}

	assert.Equal(t, expectedRequired, req)
}

func TestGatherObjectsObjectsWrapper(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	o, _ := GatherObjects(schema, []string{})

	assert.NotEmpty(t, o)

	// not worth building up the whole object structure
	assert.Equal(t, 2, len(o))
	assert.Equal(t, 11, len(o["foo"].Properties))
}
