/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var flatOutput string
var requiredOnly bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "json-schema-renderer",
	Short: "Convert a json schema to an asciidoc file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputFile string

		if len(args) > 0 {
			inputFile = args[0]
		} else {
			// check if there is somethinig to read on STDIN
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// read from stdin
				scanner := bufio.NewScanner(os.Stdin)

				var lines []string
				for {
					scanner.Scan()
					line := scanner.Text()
					if len(line) == 0 {
						break
					}
					lines = append(lines, line)
				}

				err := scanner.Err()
				if err != nil {
					log.Fatal(err)
				}

				// join lines with a linebreak to make it a valid yaml
				stdin := []byte(strings.Join(lines, "\n"))

				// create a temporary file
				tempFile, err := os.CreateTemp(os.TempDir(), "json-schema-renderer-")
				if err != nil {
					return err
				}
				defer os.Remove(tempFile.Name())

				// write the stdin to the temporary file
				if _, err := tempFile.Write(stdin); err != nil {
					return err
				}

				inputFile = tempFile.Name()

			}
		}

		if err := validateInputFile(inputFile); err != nil {
			return err
		}

		outputFile := strings.Trim(cmd.Flag("output").Value.String(), " ")
		format := strings.Trim(cmd.Flag("format").Value.String(), " ")
		title := strings.Trim(cmd.Flag("title").Value.String(), " ")
		depth, _ := cmd.Flags().GetInt("depth")
		flatObjects, _ := cmd.Flags().GetStringSlice("flat")
		headerOffset, _ := cmd.Flags().GetInt("header-level-offset")

		rendered, err := renderDoc(inputFile, format, title, depth, flatObjects, headerOffset)
		if err != nil {
			return err
		}

		if outputFile != "" {
			// write to file
			err := writeToFile(outputFile, rendered)
			if err != nil {
				return err
			}
			fmt.Printf("Output written to %s\n", outputFile)

		} else {
			fmt.Print(rendered)
		}

		return nil
	},
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "Output file")

	rootCmd.Flags().StringP("format", "f", "asciidoc", "Output format (asciidoc, markdown)")

	rootCmd.Flags().StringP("title", "t", "Root Schema", "Title of the document")
	rootCmd.Flags().Int("header-level-offset", 1, "Header level for the title. All subsequent headers will be offset accordingly")

	rootCmd.Flags().IntP("depth", "d", 0, "Depth of the schema to render")

	rootCmd.Flags().StringSlice("flat", []string{}, "Properties to always dump to json, and not recurse into, can be repeated multiple times, or comma separated. For Helm schemas, recommended to use: 'securityContext,resources,affinity,tolerations,nodeSelector'")

	rootCmd.Flags().StringVar(&flatOutput, "flat-type", "json", "How to dump the flat objects, either json or yaml")

	rootCmd.Flags().BoolVar(&requiredOnly, "required", false, "Only show required properties in the schema")
}
