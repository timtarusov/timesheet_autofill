package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version of the program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Timesheet Autofill -- version 0.0.1")
	},
}