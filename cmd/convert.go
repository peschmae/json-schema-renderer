/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"slices"
	"strings"

	"github.com/peschmae/json-schema-renderer/pkg/asciidoc"
	"github.com/peschmae/json-schema-renderer/pkg/markdown"
	"github.com/peschmae/json-schema-renderer/pkg/renderer"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var objects = make(map[string]jsonschema.Schema)
var requiredObjectNames = make(map[string]bool)

func validateInputFile(inputFile string) error {

	c := jsonschema.NewCompiler()
	_, err := c.Compile(inputFile)
	if err != nil {
		return err
	}
	// parse as json schema
	return nil
}

func renderDoc(input, format, title string, depth int, flatObjects []string, headerOffset int) (string, error) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile(input)
	if err != nil {
		return "", err
	}

	var r renderer.Renderer
	if format == "markdown" {
		r = markdown.NewRenderer(flatOutput, headerOffset)
	} else {
		r = asciidoc.NewRenderer(flatOutput, headerOffset)
	}

	output := ""

	// print schema
	output += r.Header(title, 0)
	output += r.TableHeader()
	rootPropertyKeys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		rootPropertyKeys = append(rootPropertyKeys, k)
	}
	slices.Sort(rootPropertyKeys)

	for _, key := range rootPropertyKeys {
		output += r.PropertyRow("", key, *schema.Properties[key], depth == 1)
		gatherObjects("", key, schema.Properties[key], flatObjects)
	}
	output += r.TableFooter()
	output += "\n"

	// print all schemas
	// sort keys
	objectsKeys := make([]string, 0, len(objects))
	for k := range objects {
		objectsKeys = append(objectsKeys, k)
	}
	slices.Sort(objectsKeys)

	for _, key := range objectsKeys {

		if !requiredOnly || requiredObjectNames[key] {

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
				slices.Sort(propertyKeys)

				for _, s := range propertyKeys {

					dumpValue := strings.Count(key, ">") == depth
					if slices.Contains(flatObjects, s) {
						dumpValue = true
					} else if depth == 0 {
						// avoid dumping first level objects if depth is 0
						dumpValue = false
					}

					if requiredOnly && !requiredObjectNames[key+" > "+s] { // skip non required objects
						continue
					}

					output += r.PropertyRow(key, s, *objects[key].Properties[s], dumpValue)
				}

				output += r.TableFooter()
				output += "\n"
			}
		}
	}

	return output, nil
}

func gatherObjects(parentTitle, name string, schema *jsonschema.Schema, flatObjects []string) bool {

	if isRequired(schema) {
		requiredObjectNames[name] = true
	}

	// properties matching flatObjects will be dumped directly and shouldn't be added to the list of objects
	if slices.Contains(flatObjects, name) {
		return isRequired(schema)
	}

	if parentTitle != "" {
		name = strings.Join([]string{parentTitle, name}, " > ")
		if isRequired(schema) {
			requiredObjectNames[parentTitle] = true
		}
	}

	childRequired := isRequired(schema)
	// primitive types don't need to be nested
	if schema.Types.String() != "[object]" {
		return isRequired(schema)
	} else {
		objects[name] = *schema

		if len(schema.Required) > 0 {
			requiredObjectNames[name] = true
			childRequired = true
			for _, child := range schema.Required {
				requiredObjectNames[name+" > "+child] = true
			}
		}

		for key, sch := range schema.Properties {
			if sch.Types.String() == "[object]" {
				if gatherObjects(name, key, sch, flatObjects) {
					childRequired = true
				}
			} else {
				if isRequired(sch) {
					childRequired = true
					requiredObjectNames[name+" > "+key] = true
				}
			}
		}
	}

	if childRequired {
		requiredObjectNames[name] = true
	}

	return childRequired
}

func isRequired(schema *jsonschema.Schema) bool {
	if schema.MinProperties != nil && *schema.MinProperties > 0 {
		return true
	}
	if schema.MinLength != nil && *schema.MinLength > 0 {
		return true
	}
	if schema.MinItems != nil && *schema.MinItems > 0 {
		return true
	}
	if schema.MinContains != nil && *schema.MinContains > 0 {
		return true
	}

	return slices.Contains(schema.Required, schema.Title)
}
