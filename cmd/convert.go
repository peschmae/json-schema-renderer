/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"slices"
	"strings"

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

func renderDoc(input, format, title string, depth int, flatObjects []string, headerOffset int) (string, error) {
	c := jsonschema.NewCompiler()
	schema, err := c.Compile(input)
	if err != nil {
		return "", err
	}

	objects, requiredObjectNames := util.GatherObjects(schema, flatObjects)

	var r renderer.Renderer
	if format == "markdown" {
		r = renderer.NewMarkdownRenderer(flatOutput, headerOffset)
	} else {
		r = renderer.NewAsciidocRenderer(flatOutput, headerOffset)
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
