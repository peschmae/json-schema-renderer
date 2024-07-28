package asciidoc

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type AsciiDocRenderer struct{}

func (AsciiDocRenderer) Header(title string, level int) string {
	return fmt.Sprintf("\n%s %s\n\n", strings.Repeat("=", level+1), title)
}

func (AsciiDocRenderer) PropertyHeader(title string, level int) string {
	id := strings.ToLower(strings.ReplaceAll(title, " > ", "-"))

	return fmt.Sprintf("\n[#%s]\n%s Property: %s\n\n", id, strings.Repeat("=", level+1), title)
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

		return fmt.Sprintf("|%s |%s |%s\n", schema.Title, strings.Join(schema.Types.ToStrings(), ", "), schema.Description)
	}

	id := strings.ToLower(schema.Title)
	if parent != "" {
		id = strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(schema.Title)
	}

	return fmt.Sprintf("|<<%s,%s>> |%s |%s\n", id, schema.Title, strings.Join(schema.Types.ToStrings(), ", "), schema.Description)
}
