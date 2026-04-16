package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	APIKey    string `mapstructure:"api_key"`
	ProfileID string `mapstructure:"profile_id"`
	BaseURL   string `mapstructure:"base_url"`
	Output    string `mapstructure:"-"`
}

var configDir string

func init() {
	home, err := os.UserHomeDir()
	if err == nil {
		configDir = filepath.Join(home, ".config", "ptengine-cli")
	}
}

// ConfigDir returns the config directory path.
func ConfigDir() string {
	return configDir
}

// ConfigFilePath returns the full path to the config file.
func ConfigFilePath() string {
	return filepath.Join(configDir, "config.yaml")
}

// Load reads configuration from flags, env vars, and config file.
// Precedence: flag > env > config file > default.
func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("base_url", "https://xbackend.ptengine.com")

	// Config file
	if configDir != "" {
		v.AddConfigPath(configDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		_ = v.ReadInConfig() // ignore if not found
	}

	// Env vars
	v.SetEnvPrefix("")
	v.BindEnv("api_key", "PTENGINE_API_KEY")

	// Flag overrides
	if cmd != nil {
		if f := cmd.Flags().Lookup("api-key"); f != nil && f.Changed {
			v.Set("api_key", f.Value.String())
		}
		if f := cmd.Flags().Lookup("base-url"); f != nil && f.Changed {
			v.Set("base_url", f.Value.String())
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Output format from flag
	if cmd != nil {
		cfg.Output, _ = cmd.Flags().GetString("output")
	}
	if cfg.Output == "" {
		cfg.Output = "json"
	}

	return cfg, nil
}

// Save writes config values to the config file.
func Save(apiKey, profileID, baseURL string) error {
	if configDir == "" {
		return fmt.Errorf("unable to determine config directory")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(ConfigFilePath())
	v.SetConfigType("yaml")

	// Read existing config
	_ = v.ReadInConfig()

	if apiKey != "" {
		v.Set("api_key", apiKey)
	}
	if profileID != "" {
		v.Set("profile_id", profileID)
	}
	if baseURL != "" {
		v.Set("base_url", baseURL)
	}

	return v.WriteConfigAs(ConfigFilePath())
}
