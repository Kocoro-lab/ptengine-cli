package cmd

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)

var heatmapQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query heatmap data",
	Long: `Query Ptengine heatmap data. Returns JSON with page/block/element metrics.

Query types:
  page_metrics     Page-level aggregate metrics (pv, uv, bounceRate, etc.)
  page_insight     Same metrics grouped by a dimension (requires --fun-name)
  block_metrics    Per-block metrics (impression, dropoff, etc.; page must have blocks configured)
  element_metrics  Per-element metrics (click, impression, etc.; page must have elements configured)

Output: JSON envelope {"success":true, "data":{...}, "rate_limit":{...}} on stdout.
Errors: JSON envelope {"success":false, "error":{...}} on stderr.

Run 'ptengine-cli heatmap describe' to list all available metrics, filters, and funName values.`,
	RunE: runHeatmapQuery,
}

func init() {
	f := heatmapQueryCmd.Flags()
	f.String("query-type", "", "Query type: page_metrics | page_insight | block_metrics | element_metrics (required)")
	f.String("profile-id", "", "Site profile ID, 8-char hex (falls back to config file if omitted)")
	f.String("url", "", "Target page URL to query (required)")
	f.String("start-date", "", "Start date in YYYY-MM-DD format (required)")
	f.String("end-date", "", "End date in YYYY-MM-DD format (required)")
	f.String("device-type", "ALL", "Device filter: ALL | PC | MOBILE | TABLET")
	f.StringSlice("metrics", nil, "Metrics to return, comma-separated (omit for all). Run 'heatmap describe --query-type <type>' to list")
	f.StringSlice("conversion-names", nil, "Conversion goal names, comma-separated (supports fuzzy match)")
	f.StringArray("filter", nil, "Filter condition, repeatable. Format: 'name include|exclude val1,val2'")
	f.String("filter-json", "", `Filters as JSON array, e.g. '[{"name":"country","op":"include","value":["Japan"]}]'`)
	f.String("fun-name", "", "Grouping dimension (required for page_insight): terminalType | sourceType | visitType | aiName | utmCampaign | utmSource | utmMedium | utmTerm | utmContent | week | day")

	heatmapQueryCmd.MarkFlagRequired("query-type")
	heatmapQueryCmd.MarkFlagRequired("url")
	heatmapQueryCmd.MarkFlagRequired("start-date")
	heatmapQueryCmd.MarkFlagRequired("end-date")

	heatmapCmd.AddCommand(heatmapQueryCmd)
}

func runHeatmapQuery(cmd *cobra.Command, args []string) error {
	if cfg.APIKey == "" {
		cliErr := api.NewValidationError(
			"API key is required",
			"Set via --api-key flag, PTENGINE_API_KEY env var, or 'ptengine-cli config set --api-key'.",
		)
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
	}

	queryType, _ := cmd.Flags().GetString("query-type")
	profileID, _ := cmd.Flags().GetString("profile-id")
	pageURL, _ := cmd.Flags().GetString("url")
	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	deviceType, _ := cmd.Flags().GetString("device-type")
	metrics, _ := cmd.Flags().GetStringSlice("metrics")
	conversionNames, _ := cmd.Flags().GetStringSlice("conversion-names")
	filterStrs, _ := cmd.Flags().GetStringArray("filter")
	filterJSON, _ := cmd.Flags().GetString("filter-json")
	funName, _ := cmd.Flags().GetString("fun-name")

	// Use profile-id from config if not provided via flag
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

	if !slices.Contains(api.ValidQueryTypes, queryType) {
		cliErr := api.NewValidationError(
			fmt.Sprintf("invalid query-type: %q", queryType),
			"Valid values: page_metrics, page_insight, block_metrics, element_metrics.",
		)
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
	}

	if !slices.Contains(api.ValidDeviceTypes, deviceType) {
		cliErr := api.NewValidationError(
			fmt.Sprintf("invalid device-type: %q", deviceType),
			"Valid values: ALL, PC, MOBILE, TABLET.",
		)
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
	}

	if queryType == "page_insight" && funName == "" {
		cliErr := api.NewValidationError(
			"--fun-name is required when query-type is page_insight",
			fmt.Sprintf("Valid values: %s.", strings.Join(api.ValidFunNames, ", ")),
		)
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
	}

	filters, err := parseFilters(filterStrs, filterJSON)
	if err != nil {
		cliErr := api.NewValidationError(err.Error(), "Filter format: 'name include|exclude val1,val2'. Or use --filter-json for raw JSON array.")
		exitCode := output.PrintError(cliErr, nil, cfg.Output)
		return &ExitError{Code: exitCode}
	}

	req := &api.HeatmapQueryRequest{
		QueryType:       queryType,
		ProfileID:       profileID,
		URL:             pageURL,
		StartDate:       startDate,
		EndDate:         endDate,
		DeviceType:      deviceType,
		Metrics:         metrics,
		ConversionNames: conversionNames,
		Filters:         filters,
		FunName:         funName,
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)
	resp, exitCode := client.HeatmapQuery(req)

	if resp.Success {
		output.PrintSuccess(resp, cfg.Output)
		return nil
	}

	output.PrintError(resp.Error, resp.RateLimit, cfg.Output)
	return &ExitError{Code: exitCode}
}

// parseFilters handles both --filter string syntax and --filter-json.
func parseFilters(filterStrs []string, filterJSON string) ([]api.Filter, error) {
	var filters []api.Filter

	if filterJSON != "" {
		if err := json.Unmarshal([]byte(filterJSON), &filters); err != nil {
			return nil, fmt.Errorf("invalid --filter-json: %w", err)
		}
	}

	for _, f := range filterStrs {
		// Expected: "name include|exclude val1,val2,..."
		parts := strings.SplitN(f, " ", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid filter: %q — expected format: 'name include|exclude val1,val2'", f)
		}

		name, op := parts[0], parts[1]
		if op != "include" && op != "exclude" {
			return nil, fmt.Errorf("invalid filter op: %q — must be 'include' or 'exclude'", op)
		}

		values := strings.Split(parts[2], ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}

		filters = append(filters, api.Filter{Name: name, Op: op, Value: values})
	}

	return filters, nil
}
