package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeText(t *testing.T) {

	renderer := NewAsciidocRenderer("").(*AsciiDocRenderer)

	t.Run("newline", func(t *testing.T) {
		assert.Equal(t, " +\n", renderer.escapeText("\n"))
	})

	t.Run("pipe", func(t *testing.T) {
		assert.Equal(t, "{vbar}", renderer.escapeText("|"))
	})

}

func TestPropertyId(t *testing.T) {
	renderer := NewAsciidocRenderer("").(*AsciiDocRenderer)

	t.Run("first level", func(t *testing.T) {
		assert.Equal(t, "root", renderer.propertyId("", "root"))
	})

	t.Run("second level", func(t *testing.T) {
		assert.Equal(t, "parent-child", renderer.propertyId("parent", "child"))
	})

	t.Run("nested", func(t *testing.T) {
		assert.Equal(t, "root-parent-child", renderer.propertyId("root > parent", "child"))
	})

	t.Run("case-insensitive", func(t *testing.T) {
		assert.Equal(t, "root-parent-child", renderer.propertyId("Root > paRent", "chiLd"))
	})
}
