/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"slices"

	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	util "github.com/peschmae/json-schema-renderer/pkg/schema"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func validateInputFile(inputFile string) error {

	c := jsonschema.NewCompiler()
	_, err := c.Compile(inputFile)
	if err != nil {
		return err
	}
	// parse as json schema
	return nil
}

func renderDoc(input, format, title string, depth int, flatObjects []string) (string, error) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile(input)
	if err != nil {
		return "", err
	}

	objects, requiredObjectNames := util.GatherObjects(schema, flatObjects)

	var r renderer.Renderer
	if format == "markdown" {
		r = renderer.NewMarkdownRenderer(flatOutput)
	} else {
		r = renderer.NewAsciidocRenderer(flatOutput)
	}

	output := ""

	// print schema root
	output += r.Header(title, 0)
	output += r.TableHeader()
	rootPropertyKeys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		rootPropertyKeys = append(rootPropertyKeys, k)
	}
	slices.Sort(rootPropertyKeys)

	for _, key := range rootPropertyKeys {
		output += r.PropertyRow("", key, *schema.Properties[key], depth == 1)
	}
	output += r.TableFooter()
	output += "\n"

	// print lower level objects
	out, err := renderer.RenderDocumentation(r, objects, requiredObjectNames, false, depth, flatObjects)
	if err != nil {
		return "", err
	}
	output += out

	return output, nil
}
