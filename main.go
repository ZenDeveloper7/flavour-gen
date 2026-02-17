package main

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "flavor-gen",
	Short: "Flavor Generator CLI - Build Android app flavors",
	Long:  `A tool to generate complete Android app flavor packages, including icons, keystores, and Gradle configuration.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigTypes([]string{"yaml", "json"})
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.flavor-gen")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("FLAVOR_GEN")
	viper.ReadInConfig()
}

func main() {
	Execute()
}