package asciidoc

import (
	"fmt"
	"html"
	"math"
	"strings"

	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type AsciiDocRenderer struct {
	flatOutput string
}

func (a *AsciiDocRenderer) SetFlatOutput(output string) {
	a.flatOutput = output
}

func (AsciiDocRenderer) Header(title string, level int) string {
	return fmt.Sprintf("\n%s %s\n\n", strings.Repeat("=", int(math.Min(6, float64(level+1)))), title)
}

func (a AsciiDocRenderer) PropertyHeader(title string, level int) string {

	return fmt.Sprintf("\n[#%s]\n%s Property: %s\n\n", a.propertyId("", title), strings.Repeat("=", int(math.Min(6, float64(level+1)))), title)
}

func (AsciiDocRenderer) TableHeader() string {

	return `[cols="1,1,1a,1"]
|===
|Name |Type |Default |Description

`

}

func (AsciiDocRenderer) TableFooter() string {
	return "|===\n"
}

func (a AsciiDocRenderer) PropertyRow(parent, name string, schema jsonschema.Schema, maxDepth bool) string {

	// escape the description and replace | with {vbar} to avoid table row split
	descr := a.escapeText(schema.Description)

	if schema.Types.String() != "[object]" {

		return fmt.Sprintf("|%s |%s |``%s`` |%s\n", name, strings.Join(schema.Types.ToStrings(), ", "), renderer.GetValue(schema), descr)
	}

	// on maxDepth we dump the nested object, but don't link to it
	if maxDepth {
		return fmt.Sprintf("|%s |%s |%s |%s\n", name, strings.Join(schema.Types.ToStrings(), ", "), a.dumpPropertiesToValue(schema.Properties), descr)
	}

	return fmt.Sprintf("|%s |%s |%s |%s\n", a.link(a.propertyId(parent, name), name), strings.Join(schema.Types.ToStrings(), ", "), "", descr)
}

func (a AsciiDocRenderer) TextParagraph(text string) string {
	return a.escapeText(text) + "\n\n"
}

func (a AsciiDocRenderer) dumpPropertiesToValue(properties map[string]*jsonschema.Schema) string {

	var props string
	if a.flatOutput == "yaml" {
		props = renderer.DumpPropertiesToYaml(properties)
	} else {
		props = renderer.DumpPropertiesToJson(properties)
	}

	output := "[source," + a.flatOutput + "]\n----\n"
	output += string(props)
	output += "\n----\n"

	return output
}

// escape text to make it asciidoc compatible, replacing newlines and | with {vbar}
func (AsciiDocRenderer) escapeText(text string) string {
	// escape newline
	text = strings.ReplaceAll(text, "\n", " +\n")
	// escape | and replace with {vbar} to avoid table row split
	text = strings.ReplaceAll(text, "|", "{vbar}")

	return html.EscapeString(text)
}

func (AsciiDocRenderer) link(id, title string) string {
	return fmt.Sprintf("<<%s,%s>>", id, title)
}

func (AsciiDocRenderer) propertyId(parent, title string) string {
	if parent != "" {
		return strings.ToLower(strings.ReplaceAll(parent, " > ", "-")) + "-" + strings.ToLower(title)
	} else {
		return strings.ToLower(strings.ReplaceAll(title, " > ", "-"))
	}
}
