/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "json-schema-to-asciidoc",
	Short: "Convert a json schema to an asciidoc file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		if err := validateInput(inputFile); err != nil {
			return err
		}

		return convertToAsciiDoc(inputFile, cmd.Flag("output").Value.String())
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
}
