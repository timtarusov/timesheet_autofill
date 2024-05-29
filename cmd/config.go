package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config = &cobra.Command{
	Use:   "config",
	Short: "Print the current config",
	Long:  `Print the current config`,
	Run: func(_ *cobra.Command, _ []string) {
		// Do Stuff Here
		configPath := viper.ConfigFileUsed()
		fmt.Print(configPath)
	},
}

func init() {
	rootCmd.AddCommand(config)
}
