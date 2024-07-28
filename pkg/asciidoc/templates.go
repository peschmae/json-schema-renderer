package asciidoc

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Renderer interface {
	PropertyHeader(string, int) string
	TableHeader() string
	TableFooter() string
	PropertyRow(jsonschema.Schema) string
}

type AsciiDocRenderer struct{}

func (AsciiDocRenderer) PropertyHeader(title string, level int) string {
	return "\n" + strings.Repeat("=", level+1) + " " + title + "\n"
}

func (AsciiDocRenderer) TableHeader() string {

	return `[cols="1,1"]
|===
|Name |Type |Description

`

}

func (AsciiDocRenderer) TableFooter() string {
	return "|===\n"
}

func (AsciiDocRenderer) PropertyRow(schema jsonschema.Schema) string {
	return "|" + schema.Title + " |" + strings.Join(schema.Types.ToStrings(), ", ") + " |" + schema.Description + "\n\n"
}
