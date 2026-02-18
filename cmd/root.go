package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flavor-gen",
	Short: "Flavor Generator CLI - Build Android app flavors",
	Long:  `A tool to generate complete Android app flavor packages, including icons, keystores, and Gradle configuration.`,
}

func Execute() {
	// Enable color output
	color.NoColor = false
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
}
