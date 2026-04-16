package cmd

import (
	"fmt"

	"github.com/Kocoro-lab/ptengine-cli/internal/config"
	"github.com/spf13/cobra"
)

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long: `Save configuration to ~/.config/ptengine-cli/config.yaml.

Uses the global --api-key and --base-url flags, plus a local --profile-id flag.

Examples:
  ptengine-cli config set --api-key pt-xxxxx
  ptengine-cli config set --profile-id 566d12f9
  ptengine-cli config set --api-key pt-xxxxx --profile-id 566d12f9`,
	RunE: runConfigSet,
}

func init() {
	// Only --profile-id is a local flag; --api-key and --base-url come from root persistent flags.
	configSetCmd.Flags().String("profile-id", "", "Default site profile ID")
	configCmd.AddCommand(configSetCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	// Read values only if the flag was explicitly provided by the user.
	var apiKey, profileID, baseURL string

	if f := cmd.Flags().Lookup("api-key"); f != nil && f.Changed {
		apiKey = f.Value.String()
	}
	if f := cmd.Flags().Lookup("base-url"); f != nil && f.Changed {
		baseURL = f.Value.String()
	}
	if f := cmd.Flags().Lookup("profile-id"); f != nil && f.Changed {
		profileID = f.Value.String()
	}

	if apiKey == "" && profileID == "" && baseURL == "" {
		return fmt.Errorf("at least one of --api-key, --profile-id, or --base-url must be provided")
	}

	if err := config.Save(apiKey, profileID, baseURL); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration saved to %s\n", config.ConfigFilePath())
	return nil
}
