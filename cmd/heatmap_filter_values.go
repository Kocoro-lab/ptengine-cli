package cmd

import (
	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)

var heatmapFilterValuesCmd = &cobra.Command{
	Use:   "filter-values",
	Short: "Get available filter values",
	Long: `Query available values for a specific filter type.

Dynamic filter names: os, osVersion, browser, browserVersion, screenResolution,
deviceBrand, country, region, searchEngine, socialNetwork, socialUrl, aiName,
referralSource, referralUrl, campaignUrl, utmCampaign, utmSource, utmMedium,
utmTerm, utmContent, combinedPages, originalPages, conversionName, eventName,
customDimension, eventVariable.

Use 'ptengine-cli heatmap describe' to see all filter names.`,
	RunE: runHeatmapFilterValues,
}

func init() {
	f := heatmapFilterValuesCmd.Flags()
	f.String("profile-id", "", "Site profile ID (8-char hex)")
	f.String("name", "", "Filter name to query values for (required)")
	f.String("start-date", "", "Start date YYYY-MM-DD (optional)")
	f.String("end-date", "", "End date YYYY-MM-DD (optional)")
	f.String("search", "", "Fuzzy search within values (optional)")

	heatmapFilterValuesCmd.MarkFlagRequired("name")

	heatmapCmd.AddCommand(heatmapFilterValuesCmd)
}

func runHeatmapFilterValues(cmd *cobra.Command, args []string) error {
	if cfg.APIKey == "" {
		cliErr := api.NewValidationError(
			"API key is required",
			"Set via --api-key flag, PTENGINE_API_KEY env var, or 'ptengine-cli config set --api-key'.",
		)
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
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
		cliErr := api.NewValidationError(
			"profile-id is required",
			"Set via --profile-id flag or 'ptengine-cli config set --profile-id'.",
		)
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
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

	if resp.Success {
		output.PrintSuccess(resp, cfg.Output)
		return nil
	}

	output.PrintError(resp.Error, resp.RateLimit, cfg.Output)
	return &ExitError{Code: exitCode}
}
