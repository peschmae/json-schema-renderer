/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/peschmae/json-schema-to-asciidoc/pkg/asciidoc"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var objects = make(map[string]jsonschema.Schema)

var output string

func validateInput(input string) error {

	c := jsonschema.NewCompiler()
	_, err := c.Compile(input)
	if err != nil {
		return err
	}
	// parse as json schema
	return nil
}

func convertToAsciiDoc(input string) error {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile(input)
	if err != nil {
		return err
	}

	renderer := asciidoc.AsciiDocRenderer{}

	// print schema
	output = renderer.PropertyHeader("Root Schema", 1)
	output += renderer.TableHeader()
	for _, sch := range schema.Properties {
		renderer.PropertyRow(*sch)
		gatherObjects("", sch)
	}
	output += renderer.TableFooter()
	output += "\n"

	// print all schemas
	// sort keys
	keys := make([]string, 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		output += renderer.PropertyHeader(fmt.Sprintf("Schema: %s \n", key), 1)
		output += renderer.TableHeader()
		for _, s := range objects[key].Properties {
			output += renderer.PropertyRow(*s)
		}
		output += renderer.TableFooter()
		output += "\n"
	}

	fmt.Print(output)

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
				objects[name] = *sch
				gatherObjects(name, sch)
			}
		}
	}
}
