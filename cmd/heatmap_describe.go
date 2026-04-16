package cmd

import (
	"fmt"

	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)

var heatmapDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show available query types, metrics, filters, and parameters",
	Long: `Display the API schema for heatmap queries. Useful for discovering
available parameters before constructing a query.

Without flags, shows the full schema. Use --query-type to filter
to a specific query type's metrics.`,
	RunE: runHeatmapDescribe,
}

func init() {
	heatmapDescribeCmd.Flags().String("query-type", "", "Show metrics for a specific query type [page_metrics|page_insight|block_metrics|element_metrics]")
	heatmapCmd.AddCommand(heatmapDescribeCmd)
}

func runHeatmapDescribe(cmd *cobra.Command, args []string) error {
	queryType, _ := cmd.Flags().GetString("query-type")

	if queryType != "" {
		metrics := api.GetMetricsForQueryType(queryType)
		if metrics == nil {
			cliErr := api.NewValidationError(
				fmt.Sprintf("unknown query type: %q", queryType),
				"Valid values: page_metrics, page_insight, block_metrics, element_metrics.",
			)
			exitCode := output.PrintError(cliErr, nil, cfg.Output)
			return &ExitError{Code: exitCode}
		}
		output.PrintJSON(map[string]interface{}{
			"query_type": queryType,
			"metrics":    metrics,
		}, cfg.Output)
		return nil
	}

	schema := api.GetSchema()
	output.PrintJSON(schema, cfg.Output)
	return nil
}
