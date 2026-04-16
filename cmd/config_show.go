package cmd

import (
	"strings"

	"github.com/Kocoro-lab/ptengine-cli/internal/config"
	"github.com/Kocoro-lab/ptengine-cli/internal/output"
	"github.com/spf13/cobra"
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current effective configuration (merged from flag, env, and config file).",
	Run:   runConfigShow,
}

func init() {
	configCmd.AddCommand(configShowCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) {
	// Mask the API key — only show last 4 chars
	maskedKey := ""
	if cfg.APIKey != "" {
		if len(cfg.APIKey) > 4 {
			maskedKey = strings.Repeat("*", len(cfg.APIKey)-4) + cfg.APIKey[len(cfg.APIKey)-4:]
		} else {
			maskedKey = "****"
		}
	}

	info := map[string]interface{}{
		"api_key":     maskedKey,
		"profile_id":  cfg.ProfileID,
		"base_url":    cfg.BaseURL,
		"output":      cfg.Output,
		"config_file": config.ConfigFilePath(),
	}

	output.PrintJSON(info, cfg.Output)
}
