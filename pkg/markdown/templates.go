package markdown

import (
	"fmt"
	"math"
	"strings"

	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func NewRenderer(flatOutput string, headerOffset int) renderer.Renderer {
	return &MarkdownRenderer{flatOutput: flatOutput, headerOffset: headerOffset}
}

type MarkdownRenderer struct {
	flatOutput   string
	headerOffset int
}

func (m *MarkdownRenderer) HeaderLevel(level int) int {
	return int(math.Min(6, float64(level+m.headerOffset)))
}

func (m *MarkdownRenderer) Header(title string, level int) string {
	return fmt.Sprintf("\n%s %s\n\n", strings.Repeat("#", m.HeaderLevel(level)), title)
}

func (m *MarkdownRenderer) PropertyHeader(title string, level int) string {
	id := strings.ToLower(strings.ReplaceAll(title, " > ", "-"))

	return fmt.Sprintf("\n%s <a name=\"%s\"></a> Property: %s\n\n", strings.Repeat("#", m.HeaderLevel(level)), id, title)
}

func (MarkdownRenderer) TableHeader() string {

	return `| Name | Type | Default | Description |
| :------ | :------: | :------------- | :------------- |
`

}

func (MarkdownRenderer) TableFooter() string {
	return ""
}

func (m *MarkdownRenderer) PropertyRow(parent, name string, schema jsonschema.Schema, maxDepth bool) string {

	description := strings.ReplaceAll(schema.Description, "\n", "<br>")

	if schema.Types.String() != "[object]" {
		return fmt.Sprintf("| %s | %s | `%s` | %s |\n", name, strings.Join(schema.Types.ToStrings(), ", "), renderer.GetValue(schema), description)
	}

	id := strings.ToLower(name)
	if parent != "" {
		id = strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(name)
	}

	if maxDepth {
		return fmt.Sprintf("| %s | %s | %s | %s |\n", name, strings.Join(schema.Types.ToStrings(), ", "), m.dumpPropertiesToValue(schema.Properties), description)
	}

	return fmt.Sprintf("| [%s](#%s) | %s | %s | %s |\n", name, id, strings.Join(schema.Types.ToStrings(), ", "), "", description)

}

func (MarkdownRenderer) TextParagraph(text string) string {
	return strings.ReplaceAll(text, "\n", "  \n") + "\n\n"
}

func (m *MarkdownRenderer) dumpPropertiesToValue(properties map[string]*jsonschema.Schema) string {

	var props string
	if m.flatOutput == "yaml" {
		props = renderer.DumpPropertiesToYaml(properties)
	} else {
		props = renderer.DumpPropertiesToJson(properties)
	}

	props = strings.ReplaceAll(string(props), " ", "&nbsp;")
	props = strings.ReplaceAll(string(props), "\n", "<br>")

	output := "<pre>"
	output += props
	output += "</pre>"

	return output
}
