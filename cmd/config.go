package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  "Set and view persistent CLI configuration (API key, default profile ID, etc.).",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
