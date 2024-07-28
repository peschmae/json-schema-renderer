/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/peschmae/json-schema-renderer/pkg/asciidoc"
	"github.com/peschmae/json-schema-renderer/pkg/markdown"
	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var objects = make(map[string]jsonschema.Schema)

func validateInput(input string) error {

	c := jsonschema.NewCompiler()
	_, err := c.Compile(input)
	if err != nil {
		return err
	}
	// parse as json schema
	return nil
}

func renderDoc(input, outFile, format, title string) error {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile(input)
	if err != nil {
		return err
	}

	var r renderer.Renderer
	if format == "markdown" {
		r = markdown.MarkdownRenderer{}
	} else {
		r = asciidoc.AsciiDocRenderer{}
	}

	output := ""

	// print schema
	output += r.Header(title, 0)
	output += r.TableHeader()
	for _, sch := range schema.Properties {
		output += r.PropertyRow("", *sch)
		gatherObjects("", sch)
	}
	output += r.TableFooter()
	output += "\n"

	// print all schemas
	// sort keys
	objectsKeys := make([]string, 0, len(objects))
	for k := range objects {
		objectsKeys = append(objectsKeys, k)
	}
	sort.Strings(objectsKeys)

	for _, key := range objectsKeys {
		output += r.PropertyHeader(key, strings.Count(key, ">"))
		output += r.TableHeader()
		propertyKeys := make([]string, 0, len(objects[key].Properties))
		for k := range objects[key].Properties {
			propertyKeys = append(propertyKeys, k)
		}
		sort.Strings(propertyKeys)

		for _, s := range propertyKeys {
			output += r.PropertyRow(key, *objects[key].Properties[s])
		}
		output += r.TableFooter()
		output += "\n"
	}

	if outFile != "" {
		// write to file
		err := writeToFile(outFile, output)
		if err != nil {
			return err
		}
		fmt.Printf("Output written to %s\n", outFile)

	} else {
		fmt.Print(output)
	}

	return nil
}

func writeToFile(outFile, output string) error {
	outFile = strings.Trim(outFile, " ")
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}

func gatherObjects(parentTitle string, schema *jsonschema.Schema) {
	name := schema.Title
	if parentTitle != "" {
		name = strings.Join([]string{parentTitle, schema.Title}, " > ")
	}

	if schema.Types.String() != "[object]" {
		return
	} else {
		objects[name] = *schema

		for _, sch := range schema.Properties {
			if sch.Types.String() == "[object]" {
				gatherObjects(name, sch)
			}
		}
	}
}
