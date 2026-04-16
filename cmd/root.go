package cmd

import (
	"fmt"
	"os"

	"github.com/Kocoro-lab/ptengine-cli/internal/config"
	"github.com/spf13/cobra"
)

// ExitError carries a specific exit code from command handlers.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "ptengine-cli",
	Short: "Ptengine heatmap data CLI tool",
	Long: `CLI tool for querying Ptengine heatmap data via Open API.

API commands (query, filter-values) output JSON:
  stdout: {"success":true, "data":{...}, "rate_limit":{...}}
  stderr: {"success":false, "error":{"code":N, "message":"...", "hint":"..."}}

Exit codes: 0=ok, 1=validation, 2=auth, 3=param, 4=rate-limit, 5=server, 6=network

Config precedence: --flag > PTENGINE_API_KEY env > ~/.config/ptengine-cli/config.yaml`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load(cmd)
		if err != nil {
			return err
		}
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().String("api-key", "", "Ptengine API key (env: PTENGINE_API_KEY)")
	rootCmd.PersistentFlags().String("base-url", "https://xbackend.ptengine.com", "API base URL")
	rootCmd.PersistentFlags().StringP("output", "o", "json", "Output format [json|json-pretty|table]")
}

func Execute() error {
	err := rootCmd.Execute()
	if err == nil {
		return nil
	}
	// Commands already printed structured errors before returning ExitError
	if _, ok := err.(*ExitError); ok {
		return err
	}
	// Cobra framework errors (flag validation, etc.)
	fmt.Fprintln(os.Stderr, err)
	return &ExitError{Code: 1}
}
