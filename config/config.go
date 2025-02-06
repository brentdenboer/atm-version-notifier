package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultReferenceFile = "/reference-data/version_reference.json"
	defaultConfigFile    = "/data/config/bcc-common.toml"
)

type Config struct {
	DiscordWebhookURL string `env:"DISCORD_WEBHOOK_URL,required"`
	ReferenceFilePath string `env:"REFERENCE_FILE_PATH"`
	ConfigFilePath    string `env:"CONFIG_FILE_PATH"`
}

// New creates and validates a new Config instance
func New() (*Config, error) {
	cfg := &Config{
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		ReferenceFilePath: getEnvWithDefault("REFERENCE_FILE_PATH", defaultReferenceFile),
		ConfigFilePath:    getEnvWithDefault("CONFIG_FILE_PATH", defaultConfigFile),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return cfg, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	// Validate required environment variables
	if c.DiscordWebhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL environment variable is required")
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

// getEnvWithDefault returns environment variable or default value
func getEnvWithDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// maskWebhookURL returns a masked version of the webhook URL for secure logging
func maskWebhookURL(url string) string {
	if url == "" {
		return ""
	}

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "[MASKED]"
	}

	// Keep the domain but mask the path/token
	domain := parts[2] // parts[0] is "https:", parts[1] is "", parts[2] is domain
	return fmt.Sprintf("https://%s/[MASKED]", domain)
}

// String implements Stringer interface for safe logging
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{DiscordWebhookURL: %q, ReferenceFilePath: %q, ConfigFilePath: %q}",
		maskWebhookURL(c.DiscordWebhookURL),
		c.ReferenceFilePath,
		c.ConfigFilePath,
	)
}
