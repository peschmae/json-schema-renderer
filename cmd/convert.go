/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/peschmae/json-schema-renderer/pkg/asciidoc"
	"github.com/peschmae/json-schema-renderer/pkg/markdown"
	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var objects = make(map[string]jsonschema.Schema)

func validateInputFile(inputFile string) error {

	c := jsonschema.NewCompiler()
	_, err := c.Compile(inputFile)
	if err != nil {
		return err
	}
	// parse as json schema
	return nil
}

func renderDoc(input, outFile, format, title string, depth int, flatObjects []string, flatOutput string) error {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile(input)
	if err != nil {
		return err
	}

	var r renderer.Renderer
	if format == "markdown" {
		r = &markdown.MarkdownRenderer{}
	} else {
		r = &asciidoc.AsciiDocRenderer{}
	}
	r.SetFlatOutput(flatOutput)

	output := ""

	// print schema
	output += r.Header(title, 0)
	output += r.TableHeader()
	for key, sch := range schema.Properties {
		output += r.PropertyRow("", key, *sch, depth == 1)
		gatherObjects("", key, sch, flatObjects)
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

		if depth == 0 || strings.Count(key, ">") <= depth {
			output += r.PropertyHeader(key, strings.Count(key, ">")+1)

			if objects[key].Title != "" && objects[key].Description != "" {
				output += r.TextParagraph("**" + objects[key].Title + ":** " + objects[key].Description)
			} else if objects[key].Title != "" {
				output += r.TextParagraph(objects[key].Title)
			} else if objects[key].Description != "" {
				output += r.TextParagraph(objects[key].Description)
			}

			output += r.TableHeader()

			propertyKeys := make([]string, 0, len(objects[key].Properties))
			for k := range objects[key].Properties {
				propertyKeys = append(propertyKeys, k)
			}
			sort.Strings(propertyKeys)

			for _, s := range propertyKeys {

				dumpValue := strings.Count(key, ">") == depth
				if slices.Contains(flatObjects, s) {
					dumpValue = true
				} else if depth == 0 {
					// avoid dumping first level objects if depth is 0
					dumpValue = false
				}

				output += r.PropertyRow(key, s, *objects[key].Properties[s], dumpValue)
			}

			output += r.TableFooter()
			output += "\n"
		}
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

func gatherObjects(parentTitle, name string, schema *jsonschema.Schema, flatObjects []string) {

	// properties matching flatObjects will be dumped directly and shouldn't be added to the list of objects
	if slices.Contains(flatObjects, name) {
		return
	}

	if parentTitle != "" {
		name = strings.Join([]string{parentTitle, name}, " > ")
	}

	// primitive types don't need to be nested
	if schema.Types.String() != "[object]" {
		return
	} else {
		objects[name] = *schema

		for key, sch := range schema.Properties {
			if sch.Types.String() == "[object]" {
				gatherObjects(name, key, sch, flatObjects)
			}
		}
	}
}
