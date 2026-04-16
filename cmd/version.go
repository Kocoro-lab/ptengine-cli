package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	appVersion = "dev"
	appCommit  = "none"
	appDate    = "unknown"
)

func SetVersionInfo(version, commit, date string) {
	appVersion = version
	appCommit = commit
	appDate = date
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Print version information.

Use --short to output only the version number (useful for scripting):
  ptengine-cli version --short   →  0.1.0`,
	Run: func(cmd *cobra.Command, args []string) {
		short, _ := cmd.Flags().GetBool("short")
		if short {
			fmt.Println(appVersion)
			return
		}

		outputFormat, _ := cmd.Flags().GetString("output")
		info := map[string]string{
			"version": appVersion,
			"commit":  appCommit,
			"date":    appDate,
		}

		switch outputFormat {
		case "json":
			data, _ := json.Marshal(info)
			fmt.Println(string(data))
		case "json-pretty":
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
		default:
			fmt.Printf("ptengine-cli %s (commit: %s, built: %s)\n", appVersion, appCommit, appDate)
		}
	},
}

func init() {
	versionCmd.Flags().Bool("short", false, "Print only the version number")
	rootCmd.AddCommand(versionCmd)
}
