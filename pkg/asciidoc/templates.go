package asciidoc

import (
	"fmt"
	"html"
	"strings"

	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type AsciiDocRenderer struct{}

func (AsciiDocRenderer) Header(title string, level int) string {
	return fmt.Sprintf("\n%s %s\n\n", strings.Repeat("=", level+1), title)
}

func (a AsciiDocRenderer) PropertyHeader(title string, level int) string {

	return fmt.Sprintf("\n[#%s]\n%s Property: %s\n\n", a.propertyId("", title), strings.Repeat("=", level+1), title)
}

func (AsciiDocRenderer) TableHeader() string {

	return `[cols="1,1,1,1"]
|===
|Name |Type |Default |Description

`

}

func (AsciiDocRenderer) TableFooter() string {
	return "|===\n"
}

func (a AsciiDocRenderer) PropertyRow(parent string, schema jsonschema.Schema) string {

	descr := html.EscapeString(schema.Description)

	if schema.Types.String() != "[object]" {

		return fmt.Sprintf("|%s |%s |%s |%s\n", schema.Title, strings.Join(schema.Types.ToStrings(), ", "), renderer.GetValue(schema), descr)
	}

	// we don't show the default value for objects
	return fmt.Sprintf("|%s |%s |%s |%s\n", a.link(a.propertyId(parent, schema.Title), schema.Title), strings.Join(schema.Types.ToStrings(), ", "), "", descr)
}

func (AsciiDocRenderer) link(id, title string) string {
	return fmt.Sprintf("<<%s,%s>>", id, title)
}

func (AsciiDocRenderer) propertyId(parent, title string) string {
	if parent != "" {
		return strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(title)
	} else {
		return strings.ToLower(title)
	}
}
