package cmd

import (
	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)

var dataQueryCancelCmd = &cobra.Command{
	Use:   "cancel <trace-id>",
	Short: "Cancel an in-flight data-query by trace id",
	Long: `Cancel an in-flight data-query.

The trace id must be the X-Request-Id you passed to 'data-query query --request-id'
(or to '... stream --request-id'), or the traceId you read from a streaming
'routed' event.

Output (cancelling): {"success":true, "data":{"traceId":"...", "status":"cancelling"}}
Output (not_found): {"success":false, "data":{"traceId":"...", "status":"not_found"}} — exit 0
  The query already finished or the id is unknown; not a failure from the caller's
  perspective. Inspect the original query response instead.`,
	Args: cobra.ExactArgs(1),
	RunE: runDataQueryCancel,
}

func init() {
	dataQueryCmd.AddCommand(dataQueryCancelCmd)
}

func runDataQueryCancel(cmd *cobra.Command, args []string) error {
	if cfg.APIKey == "" {
		return failValidation("API key is required", "")
	}
	traceID := args[0]
	client := api.NewClient(cfg.BaseURL, cfg.APIKey)
	resp, exitCode := client.DataQueryCancel(traceID)
	output.PrintEnvelope(resp, cfg.Output)
	if exitCode != api.ExitOK {
		return &ExitError{Code: exitCode}
	}
	return nil
}
