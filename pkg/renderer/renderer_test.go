package renderer

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
)

func TestHeaderLeveL(t *testing.T) {

	t.Run("default", func(t *testing.T) {
		HeaderOffset = 1

		assert.Equal(t, 1, headerLevel(0))
	})

	t.Run("two", func(t *testing.T) {
		HeaderOffset = 1
		assert.Equal(t, 2, headerLevel(1))
	})

	t.Run("set", func(t *testing.T) {
		HeaderOffset = 2

		assert.Equal(t, 2, headerLevel(0))
	})

	t.Run("max level", func(t *testing.T) {
		HeaderOffset = 1

		assert.Equal(t, 6, headerLevel(7))
	})

	t.Run("max < offset", func(t *testing.T) {
		HeaderOffset = 8

		assert.Equal(t, 6, headerLevel(1))
	})

}

func TestGetValue(t *testing.T) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile("../../testdata/test.schema.json")

	if err != nil {
		t.Fatal(err)
	}

	t.Run("None", func(t *testing.T) {

		v := getValue(*schema.Properties["foo"].Properties["map_key"])

		assert.Equal(t, "", v)
	})

	t.Run("integer", func(t *testing.T) {

		v := getValue(*schema.Properties["foo"].Properties["max_key"])

		assert.Equal(t, "10", v)
	})

	t.Run("number", func(t *testing.T) {
		v := getValue(*schema.Properties["foo"].Properties["range_float_key"])

		assert.Equal(t, "2.2", v)
	})

	t.Run("string", func(t *testing.T) {
		v := getValue(*schema.Properties["foo"].Properties["one_of_strings"])

		assert.Equal(t, "one", v)
	})

	t.Run("array", func(t *testing.T) {
		v := getValue(*schema.Properties["foo"].Properties["array_key"])

		assert.Equal(t, "[]\n", v)
	})

	t.Run("bool", func(t *testing.T) {
		v := getValue(*schema.Properties["foo"].Properties["bool"])

		assert.Equal(t, "true", v)
	})

}
