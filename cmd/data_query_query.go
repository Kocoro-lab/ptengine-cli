package cmd

import (
	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)

var dataQueryQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Submit a natural-language data query (sync)",
	Long: `Submit a natural-language analytics question and wait for the full response.

Output: JSON envelope {"success":true, "data":{...}, "rate_limit":{...}} on stdout.
Errors: JSON envelope {"success":false, "error":{...}} on stderr.

Examples:
  ptengine-cli data-query query --question "昨天的 UV"
  ptengine-cli data-query query --question "中国 vs 海外的支付完成率" --request-id $(uuidgen)
  ptengine-cli data-query query --scenario user_overview --params '{"userId":"abc..."}'

Tip: pass --request-id (UUID) so you can later cancel via 'data-query cancel <id>'
without waiting for the response. Without --request-id the server generates one
that only arrives in the final response body.`,
	RunE: runDataQueryQuery,
}

func init() {
	registerDataQueryFlags(dataQueryQueryCmd)
	dataQueryCmd.AddCommand(dataQueryQueryCmd)
}

func runDataQueryQuery(cmd *cobra.Command, args []string) error {
	req, requestID, err := buildDataQueryRequest(cmd)
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)
	resp, exitCode := client.DataQueryQuery(req, requestID)
	output.PrintEnvelope(resp, cfg.Output)
	if exitCode != api.ExitOK {
		return &ExitError{Code: exitCode}
	}
	return nil
}
