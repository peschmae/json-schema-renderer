package markdown

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type MarkdownRenderer struct{}

func (MarkdownRenderer) Header(title string, level int) string {
	return fmt.Sprintf("\n%s %s\n\n", strings.Repeat("#", level+1), title)
}

func (MarkdownRenderer) PropertyHeader(title string, level int) string {
	id := strings.ToLower(strings.ReplaceAll(title, " > ", "-"))

	return fmt.Sprintf("\n%s <a name=\"%s\"></a> Property: %s\n\n", strings.Repeat("#", level+1), id, title)
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
		return fmt.Sprintf("| %s | %s | %s |\n", schema.Title, strings.Join(schema.Types.ToStrings(), ", "), description)
	}

	id := strings.ToLower(schema.Title)
	if parent != "" {
		id = strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(schema.Title)
	}

	return fmt.Sprintf("| [%s](#%s) | %s | %s |\n", schema.Title, id, strings.Join(schema.Types.ToStrings(), ", "), description)

}
