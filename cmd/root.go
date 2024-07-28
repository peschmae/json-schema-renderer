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
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateInput(cmd.Flag("input").Value.String()); err != nil {
			return err
		}

		return convertToAsciiDoc(cmd.Flag("input").Value.String())
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
	rootCmd.Flags().StringP("input", "i", "", "Input file")
	rootCmd.Flags().StringP("output", "o", "", "Output file")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
