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
	// Parse intervals with proper error handling
	fileCheckInterval, err := strconv.Atoi(getEnvWithDefault("FILE_CHECK_INTERVAL_SECONDS", "3600"))
	if err != nil {
		return nil, fmt.Errorf("invalid FILE_CHECK_INTERVAL_SECONDS value: %w", err)
	}

	recheckTimeout, err := strconv.Atoi(getEnvWithDefault("RECHECK_TIMEOUT_SECONDS", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid RECHECK_TIMEOUT_SECONDS value: %w", err)
	}

	cfg := &Config{
		DiscordWebhookURL:        os.Getenv("DISCORD_WEBHOOK_URL"),
		ReferenceFilePath:        getEnvWithDefault("REFERENCE_FILE_PATH", "/reference-data/version_reference.json"),
		ConfigFilePath:           getEnvWithDefault("CONFIG_FILE_PATH", "/data/config/bcc-common.toml"),
		FileCheckIntervalSeconds: fileCheckInterval,
		RecheckTimeoutSeconds:    recheckTimeout,
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	// Validate required environment variables
	if c.DiscordWebhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL environment variable is required")
	}

	// Validate intervals are within reasonable bounds
	if c.FileCheckIntervalSeconds <= 0 {
		return fmt.Errorf("FILE_CHECK_INTERVAL_SECONDS must be positive")
	}
	if c.FileCheckIntervalSeconds > 86400 { // 24 hours
		return fmt.Errorf("FILE_CHECK_INTERVAL_SECONDS must not exceed 86400 (24 hours)")
	}

	if c.RecheckTimeoutSeconds <= 0 {
		return fmt.Errorf("RECHECK_TIMEOUT_SECONDS must be positive")
	}
	if c.RecheckTimeoutSeconds > 3600 { // 1 hour
		return fmt.Errorf("RECHECK_TIMEOUT_SECONDS must not exceed 3600 (1 hour)")
	}

	// Validate file paths
	if c.ReferenceFilePath == "" {
		return fmt.Errorf("REFERENCE_FILE_PATH must not be empty")
	}
	if c.ConfigFilePath == "" {
		return fmt.Errorf("CONFIG_FILE_PATH must not be empty")
	}

	return nil
}

// getEnvWithDefault returns the environment variable value or the default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
