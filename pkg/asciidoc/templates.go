package asciidoc

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type AsciiDocRenderer struct{}

func (AsciiDocRenderer) Header(title string, level int) string {
	return "\n" + strings.Repeat("=", level+1) + " " + title + "\n\n"
}

func (AsciiDocRenderer) PropertyHeader(title string, level int) string {
	id := strings.ToLower(strings.ReplaceAll(title, " > ", "-"))

	return "\n[#" + id + "]\n" + strings.Repeat("=", level+1) + " Property: " + title + "\n\n"
}

func (AsciiDocRenderer) TableHeader() string {

	return `[cols="1,1,1"]
|===
|Name |Type |Description

`

}

func (AsciiDocRenderer) TableFooter() string {
	return "|===\n"
}

func (AsciiDocRenderer) PropertyRow(parent string, schema jsonschema.Schema) string {

	if schema.Types.String() != "[object]" {
		return "|" + schema.Title + " |" + strings.Join(schema.Types.ToStrings(), ", ") + " |" + schema.Description + "\n\n"
	}

	id := strings.ToLower(schema.Title)
	if parent != "" {
		id = strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(schema.Title)
	}

	return "|<<" + id + "," + schema.Title + ">> |" + strings.Join(schema.Types.ToStrings(), ", ") + " |" + schema.Description + "\n\n"
}
