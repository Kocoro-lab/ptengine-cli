package cmd

import (
	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)
var heatmapFilterValuesCmd = &cobra.Command{
	Use:   "filter-values",
	Short: "Get available values for a filter",
	Long: `Query the Ptengine API for available values of a given filter name.
Use this to discover what values can be passed in --filter for heatmap query.

Only dynamic filter names are supported (not fixed ones like deviceType/sourceType).
Run 'ptengine-cli heatmap describe' to see the full list of filter names.

Output: JSON envelope {"success":true, "data":[...]} on stdout.`,
	RunE: runHeatmapFilterValues,
}

func init() {
	f := heatmapFilterValuesCmd.Flags()
	f.String("profile-id", "", "Site profile ID, 8-char hex (falls back to config file if omitted)")
	f.String("name", "", "Filter name to query, e.g. country, browser, utmSource (required)")
	f.String("start-date", "", "Start date in YYYY-MM-DD format")
	f.String("end-date", "", "End date in YYYY-MM-DD format")
	f.String("search", "", "Fuzzy search keyword to narrow results")

	heatmapFilterValuesCmd.MarkFlagRequired("name")

	heatmapCmd.AddCommand(heatmapFilterValuesCmd)
}

func runHeatmapFilterValues(cmd *cobra.Command, args []string) error {
	if cfg.APIKey == "" {
		return failValidation("API key is required",
			"Set via --api-key flag, PTENGINE_API_KEY env var, or 'ptengine-cli config set --api-key'.")
	}

	profileID, _ := cmd.Flags().GetString("profile-id")
	name, _ := cmd.Flags().GetString("name")
	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	search, _ := cmd.Flags().GetString("search")

	if profileID == "" {
		profileID = cfg.ProfileID
	}
	if profileID == "" {
		return failValidation("profile-id is required",
			"Set via --profile-id flag or 'ptengine-cli config set --profile-id'.")
	}

	req := &api.FilterValuesRequest{
		ProfileID: profileID,
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Search:    search,
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)
	resp, exitCode := client.HeatmapFilterValues(req)
	output.PrintEnvelope(resp, cfg.Output)
	if exitCode != api.ExitOK {
		return &ExitError{Code: exitCode}
	}
	return nil
}
