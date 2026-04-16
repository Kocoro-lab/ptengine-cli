package cmd

import "github.com/spf13/cobra"

var heatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "Heatmap data commands",
	Long: `Query and explore Ptengine heatmap data.

Subcommands:
  query          Query heatmap metrics (requires API key)
  filter-values  List available values for a filter (requires API key)
  describe       Show available query types, metrics, and filter names (no API key needed)`,
}

func init() {
	rootCmd.AddCommand(heatmapCmd)
}
