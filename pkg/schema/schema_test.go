package schema

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
)

func TestIsRequired(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}
	t.Run("minItems", func(t *testing.T) {
		assert.True(t, isRequired(schema.Properties["foo"].Properties["array_key"]))
	})

	t.Run("minProperties", func(t *testing.T) {
		assert.True(t, isRequired(schema.Properties["foo"].Properties["map_key"]))
	})

	t.Run("minLength", func(t *testing.T) {
		assert.True(t, isRequired(schema.Properties["foo"].Properties["string_key"]))
	})

	t.Run("negative test", func(t *testing.T) {
		assert.False(t, isRequired(schema.Properties["foo"].Properties["max_key"]))
	})
}

func TestGatherObjects(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	obj, req := GatherObjects(schema, []string{})

	t.Run("requiredObjectNames", func(t *testing.T) {
		assert.NotEmpty(t, req)

		expectedRequired := map[string]bool{
			"foo":              true,
			"foo > array_key":  true,
			"foo > map_key":    true,
			"foo > string_key": true,
			"map_key":          true,
		}

		assert.Equal(t, expectedRequired, req)
	})

	t.Run("objectsWrapper", func(t *testing.T) {
		assert.NotEmpty(t, obj)

		// not worth building up the whole object structure
		assert.Equal(t, 2, len(obj))
		assert.Equal(t, 11, len(obj["foo"].Properties))
	})
}
