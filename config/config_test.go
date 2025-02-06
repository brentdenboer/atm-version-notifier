package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Helper function to clean environment variables
	cleanEnv := func() {
		os.Unsetenv("DISCORD_WEBHOOK_URL")
		os.Unsetenv("REFERENCE_FILE_PATH")
		os.Unsetenv("CONFIG_FILE_PATH")
		os.Unsetenv("FILE_CHECK_INTERVAL_SECONDS")
		os.Unsetenv("RECHECK_TIMEOUT_SECONDS")
	}

	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config with all defaults except webhook",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL": "https://discord.webhook.url",
			},
			wantErr: false,
		},
		{
			name: "valid config with custom values",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":         "https://discord.webhook.url",
				"REFERENCE_FILE_PATH":         "/custom/reference.json",
				"CONFIG_FILE_PATH":            "/custom/config.toml",
				"FILE_CHECK_INTERVAL_SECONDS": "1800",
				"RECHECK_TIMEOUT_SECONDS":     "60",
			},
			wantErr: false,
		},
		{
			name:        "missing webhook URL",
			envVars:     map[string]string{},
			wantErr:     true,
			errContains: "DISCORD_WEBHOOK_URL environment variable is required",
		},
		{
			name: "invalid file check interval",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":         "https://discord.webhook.url",
				"FILE_CHECK_INTERVAL_SECONDS": "-1",
			},
			wantErr:     true,
			errContains: "must be positive",
		},
		{
			name: "file check interval too large",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":         "https://discord.webhook.url",
				"FILE_CHECK_INTERVAL_SECONDS": "86401",
			},
			wantErr:     true,
			errContains: "must not exceed 86400",
		},
		{
			name: "invalid recheck timeout",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":     "https://discord.webhook.url",
				"RECHECK_TIMEOUT_SECONDS": "-1",
			},
			wantErr:     true,
			errContains: "must be positive",
		},
		{
			name: "recheck timeout too large",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":     "https://discord.webhook.url",
				"RECHECK_TIMEOUT_SECONDS": "3601",
			},
			wantErr:     true,
			errContains: "must not exceed 3600",
		},
		{
			name: "invalid file check interval format",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":         "https://discord.webhook.url",
				"FILE_CHECK_INTERVAL_SECONDS": "not_a_number",
			},
			wantErr:     true,
			errContains: "invalid FILE_CHECK_INTERVAL_SECONDS value",
		},
		{
			name: "invalid recheck timeout format",
			envVars: map[string]string{
				"DISCORD_WEBHOOK_URL":     "https://discord.webhook.url",
				"RECHECK_TIMEOUT_SECONDS": "not_a_number",
			},
			wantErr:     true,
			errContains: "invalid RECHECK_TIMEOUT_SECONDS value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment before each test
			cleanEnv()

			// Set up environment for this test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Run the test
			cfg, err := New()
			if tt.wantErr {
				if err == nil {
					t.Error("New() expected error but got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("New() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
				return
			}

			// Verify the configuration values
			if cfg.DiscordWebhookURL != tt.envVars["DISCORD_WEBHOOK_URL"] {
				t.Errorf("New() DiscordWebhookURL = %v, want %v", cfg.DiscordWebhookURL, tt.envVars["DISCORD_WEBHOOK_URL"])
			}

			// Check that default values are set when not provided
			if tt.envVars["REFERENCE_FILE_PATH"] == "" && cfg.ReferenceFilePath != "/reference-data/version_reference.json" {
				t.Errorf("New() ReferenceFilePath = %v, want default value", cfg.ReferenceFilePath)
			}
			if tt.envVars["CONFIG_FILE_PATH"] == "" && cfg.ConfigFilePath != "/data/config/bcc-common.toml" {
				t.Errorf("New() ConfigFilePath = %v, want default value", cfg.ConfigFilePath)
			}
			if tt.envVars["FILE_CHECK_INTERVAL_SECONDS"] == "" && cfg.FileCheckIntervalSeconds != 3600 {
				t.Errorf("New() FileCheckIntervalSeconds = %v, want default value", cfg.FileCheckIntervalSeconds)
			}
			if tt.envVars["RECHECK_TIMEOUT_SECONDS"] == "" && cfg.RecheckTimeoutSeconds != 300 {
				t.Errorf("New() RecheckTimeoutSeconds = %v, want default value", cfg.RecheckTimeoutSeconds)
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return s != "" && substr != "" && len(s) >= len(substr) && s != substr && s[len(s)-1] != '/' && s != substr+"/" && s != "/"+substr
}
