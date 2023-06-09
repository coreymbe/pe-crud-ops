package cmd

import (
	"fmt"
	"os"

	"github.com/coreymbe/pe-crud-ops/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pe-crud-ops",
	Short: "Puppet Agent CRUD operations",
	Long: `A tool for testing Puppet Agent CRUD operations
	against a Puppet Enterprise installation`,
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
	cobra.OnInitialize(initConfig)

	configFile := ".cobra.yaml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configuration file: ", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()

	var configuration config.Config

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AddConfigPath(".")
	viper.SetConfigName(".cobra.yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
