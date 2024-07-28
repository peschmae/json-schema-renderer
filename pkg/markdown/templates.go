package markdown

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type MarkdownRenderer struct{}

func (MarkdownRenderer) Header(title string, level int) string {
	return "\n" + strings.Repeat("#", level+1) + " " + title + "\n\n"
}

func (MarkdownRenderer) PropertyHeader(title string, level int) string {
	id := strings.ToLower(strings.ReplaceAll(title, " > ", "-"))

	return strings.Repeat("#", level+1) + " <a name=\"" + id + "\"></a>" + " Property: " + title + "\n\n"
}

func (MarkdownRenderer) TableHeader() string {

	return `| Name | Type | Description |
| :------ | :------: | :------------- |
`

}

func (MarkdownRenderer) TableFooter() string {
	return ""
}

func (MarkdownRenderer) PropertyRow(parent string, schema jsonschema.Schema) string {

	description := strings.ReplaceAll(schema.Description, "\n", "<br>")

	if schema.Types.String() != "[object]" {
		return "| " + schema.Title + " | " + strings.Join(schema.Types.ToStrings(), ", ") + " | " + description + " |\n"
	}

	id := strings.ToLower(schema.Title)
	if parent != "" {
		id = strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(schema.Title)
	}

	return "| [" + schema.Title + "](#" + id + ") | " + strings.Join(schema.Types.ToStrings(), ", ") + " | " + description + " |\n"

}
