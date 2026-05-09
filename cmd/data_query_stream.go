package cmd

import (
	"os"

	"github.com/Kocoro-lab/ptengine-cli/internal/api"
	"github.com/spf13/cobra"
)

var dataQueryStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Submit a query and stream phase events as Server-Sent Events",
	Long: `Submit a query and stream the backend's progress events.

Each line written to stdout is one SSE frame in JSON form:
  {"event":"routed","data":{"phase":"routed","skill":"...","intent":"...","route_ms":1409,"trace_id":"..."}}
  {"event":"handler","data":{"phase":"handler","intent":"..."}}
  {"event":"final","data":{"phase":"final","code":200,"message":"OK","data":{...}}}
  {"event":"cancelled","data":{"phase":"cancelled","trace_id":"..."}}
  {"event":"error","data":{"phase":"error","message":"..."}}

The 'final' event payload uses the same envelope as 'data-query query'.
'routed' arrives within 1-2s and carries the trace_id for cancellation.

Use 'stream' when you want progress feedback to a user; use 'query' when you
just need the result and don't care about intermediate events.`,
	RunE: runDataQueryStream,
}

func init() {
	registerDataQueryFlags(dataQueryStreamCmd)
	dataQueryCmd.AddCommand(dataQueryStreamCmd)
}

func runDataQueryStream(cmd *cobra.Command, args []string) error {
	req, requestID, err := buildDataQueryRequest(cmd)
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)
	exitCode := client.DataQueryStream(req, requestID, os.Stdout)
	if exitCode != api.ExitOK {
		return &ExitError{Code: exitCode}
	}
	return nil
}
