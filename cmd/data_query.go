package cmd

import (
	"slices"
	"strings"

	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/spf13/cobra"
)

// dataQueryCmd is the parent for /open-api/v1/data-query/* subcommands.
//
// Subcommands:
//   query   POST /open-api/v1/data-query/query              sync NL query
//   stream  POST /open-api/v1/data-query/stream             SSE streaming variant
//   cancel  POST /open-api/v1/data-query/{traceId}/cancel   abort an in-flight query
var dataQueryCmd = &cobra.Command{
	Use:   "data-query",
	Short: "Natural-language data query commands",
	Long: `Query Ptengine analytics data via natural-language questions.

Subcommands:
  query    Submit a question and wait for the full result (sync)
  stream   Submit a question and consume Server-Sent Events as it progresses
  cancel   Abort an in-flight query by trace id

The data-query API translates a natural-language question (zh / en / ja) into
a deterministic SQL or atomic-template query against the profile's behavior
data. Returns rows + the SQL the backend ran.`,
}

func init() {
	rootCmd.AddCommand(dataQueryCmd)
}

// registerDataQueryFlags attaches the flag set shared by `query` and `stream`.
func registerDataQueryFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("question", "", "Natural language question (required unless --scenario is set)")
	f.String("profile-id", "", "Site profile ID (falls back to config file if omitted)")
	f.String("context", "", "Optional context: clarification reply, time anchor hints, ≤2000 chars")
	f.String("scenario", "", "Predefined scenario: "+strings.Join(api.ValidScenarios, " | "))
	f.String("params", "", "Scenario parameters as JSON object, e.g. '{\"userId\":\"...\"}'")
	f.String("request-id", "", "Custom trace id (UUID recommended); used as X-Request-Id and accepted by 'data-query cancel'")
}

// buildDataQueryRequest reads the shared flag set, applies config fallbacks,
// validates inputs, and returns the request body plus the optional X-Request-Id.
// On validation failure it renders stderr JSON via failValidation and returns
// the resulting *ExitError; callers should return that error as-is.
func buildDataQueryRequest(cmd *cobra.Command) (*api.DataQueryRequest, string, error) {
	if cfg.APIKey == "" {
		return nil, "", failValidation("API key is required",
			"Set via --api-key flag, PTENGINE_API_KEY env var, or 'ptengine-cli config set --api-key'.")
	}

	question, _ := cmd.Flags().GetString("question")
	profileID, _ := cmd.Flags().GetString("profile-id")
	contextStr, _ := cmd.Flags().GetString("context")
	scenario, _ := cmd.Flags().GetString("scenario")
	paramsStr, _ := cmd.Flags().GetString("params")
	requestID, _ := cmd.Flags().GetString("request-id")

	if profileID == "" {
		profileID = cfg.ProfileID
	}
	if profileID == "" {
		return nil, "", failValidation("profile-id is required",
			"Set via --profile-id flag or 'ptengine-cli config set --profile-id'.")
	}
	if question == "" && scenario == "" {
		return nil, "", failValidation("either --question or --scenario is required",
			"Free-form analytics: --question 'your question'. User-specific: --scenario user_overview --params '{\"userId\":\"...\"}'.")
	}
	if scenario != "" && !slices.Contains(api.ValidScenarios, scenario) {
		return nil, "", failValidation("invalid --scenario: "+scenario,
			"Valid values: "+strings.Join(api.ValidScenarios, ", ")+".")
	}
	if len(contextStr) > 2000 {
		return nil, "", failValidation("context exceeds 2000 characters",
			"Trim the context payload — the server enforces a 2000-char cap.")
	}

	req := &api.DataQueryRequest{
		Question:  question,
		ProfileID: profileID,
		Context:   contextStr,
		Scenario:  scenario,
	}
	if paramsStr != "" {
		params, err := parseJSONObject(paramsStr)
		if err != nil {
			return nil, "", failValidation("invalid --params JSON: "+err.Error(),
				"Pass a JSON object literal, e.g. --params '{\"userId\":\"abc\"}'.")
		}
		req.Params = params
	}
	return req, requestID, nil
}
