package cmd

import "github.com/spf13/cobra"

var heatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "Heatmap data commands",
	Long:  "Query and explore Ptengine heatmap data. Use subcommands: query, filter-values, describe.",
}

func init() {
	rootCmd.AddCommand(heatmapCmd)
}
