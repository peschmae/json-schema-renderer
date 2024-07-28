/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "json-schema-renderer",
	Short: "Convert a json schema to an asciidoc file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		if err := validateInput(inputFile); err != nil {
			return err
		}

		output := strings.Trim(cmd.Flag("output").Value.String(), " ")
		format := strings.Trim(cmd.Flag("format").Value.String(), " ")
		title := strings.Trim(cmd.Flag("title").Value.String(), " ")

		return renderDoc(inputFile, output, format, title)
	},
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
}
