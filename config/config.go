package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DiscordWebhookURL        string
	ReferenceFilePath        string
	ConfigFilePath           string
	FileCheckIntervalSeconds int // Interval between file checks in seconds
	RecheckTimeoutSeconds    int // Timeout before rechecking after an error
}

// New creates and validates a new Config instance
func New() (*Config, error) {
	fileCheckInterval, _ := strconv.Atoi(getEnvWithDefault("FILE_CHECK_INTERVAL_SECONDS", "3600"))
	recheckTimeout, _ := strconv.Atoi(getEnvWithDefault("RECHECK_TIMEOUT_SECONDS", "300"))

	cfg := &Config{
		DiscordWebhookURL:        os.Getenv("DISCORD_WEBHOOK_URL"),
		ReferenceFilePath:        getEnvWithDefault("REFERENCE_FILE_PATH", "/reference-data/version_reference.json"),
		ConfigFilePath:           getEnvWithDefault("CONFIG_FILE_PATH", "/data/config/bcc-common.toml"),
		FileCheckIntervalSeconds: fileCheckInterval,
		RecheckTimeoutSeconds:    recheckTimeout,
	}

	// Validate required environment variables
	if cfg.DiscordWebhookURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL environment variable is required")
	}

	// Validate intervals are positive
	if cfg.FileCheckIntervalSeconds <= 0 {
		return nil, fmt.Errorf("FILE_CHECK_INTERVAL_SECONDS must be positive")
	}
	if cfg.RecheckTimeoutSeconds <= 0 {
		return nil, fmt.Errorf("RECHECK_TIMEOUT_SECONDS must be positive")
	}

	return cfg, nil
}

// getEnvWithDefault returns the environment variable value or the default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
