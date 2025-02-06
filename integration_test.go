package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"atm10-version-notifier/config"
	"atm10-version-notifier/discord"
)

func TestVersionMonitoringWorkflow(t *testing.T) {
	// Create temporary directories for test files
	tmpDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create config and reference directories
	configDir := filepath.Join(tmpDir, "config")
	referenceDir := filepath.Join(tmpDir, "reference-data")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	if err := os.MkdirAll(referenceDir, 0755); err != nil {
		t.Fatalf("Failed to create reference dir: %v", err)
	}

	// Set up test files
	configFile := filepath.Join(configDir, "bcc-common.toml")
	referenceFile := filepath.Join(referenceDir, "version_reference.json")

	// Create initial config file
	initialConfig := `[general]
modpackProjectID = 925200
modpackName = "All The Mods 10"
modpackVersion = "2.32"
useMetadata = false`

	if err := os.WriteFile(configFile, []byte(initialConfig), 0644); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	// Set up test environment variables
	os.Setenv("CONFIG_FILE_PATH", configFile)
	os.Setenv("REFERENCE_FILE_PATH", referenceFile)
	os.Setenv("DISCORD_WEBHOOK_URL", "https://discord.com/api/webhooks/1337145859212705840/NjG0-49RizfFqbmiILobROQWu1twuGt5pECvSlFYzW8aGG38WJHl-QGvg5qp_ErKU9PS")
	os.Setenv("RECHECK_TIMEOUT_SECONDS", "1") // Use shorter timeout for testing
	defer func() {
		os.Unsetenv("CONFIG_FILE_PATH")
		os.Unsetenv("REFERENCE_FILE_PATH")
		os.Unsetenv("DISCORD_WEBHOOK_URL")
		os.Unsetenv("RECHECK_TIMEOUT_SECONDS")
	}()

	// Store original notification function and restore it after test
	originalNotify := discord.SendVersionUpdateNotification
	defer func() {
		discord.SendVersionUpdateNotification = originalNotify
	}()

	// Create a mock Discord client for testing
	mockDiscordCalls := 0
	discord.SendVersionUpdateNotification = func(webhookURL, modpackName, oldVersion, newVersion string) error {
		mockDiscordCalls++
		if modpackName != "All The Mods 10" {
			t.Errorf("Unexpected modpack name: %s", modpackName)
		}
		if oldVersion != "2.32" || newVersion != "2.33" {
			t.Errorf("Unexpected version change: %s -> %s", oldVersion, newVersion)
		}
		return nil
	}

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Perform initial version check
	if err := checkVersionAndNotify(cfg); err != nil {
		t.Fatalf("Initial version check failed: %v", err)
	}

	// Verify that no Discord notification was sent for initial check
	if mockDiscordCalls != 0 {
		t.Errorf("Expected no Discord notifications for initial check, got %d", mockDiscordCalls)
	}

	// Update the config file with a new version
	updatedConfig := `[general]
modpackProjectID = 925200
modpackName = "All The Mods 10"
modpackVersion = "2.33"
useMetadata = false`

	// Wait a moment to ensure file system events are processed
	time.Sleep(time.Second)

	if err := os.WriteFile(configFile, []byte(updatedConfig), 0644); err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	// Wait for the recheck timeout
	time.Sleep(2 * time.Second)

	// Check version again
	if err := checkVersionAndNotify(cfg); err != nil {
		t.Fatalf("Version check after update failed: %v", err)
	}

	// Verify that a Discord notification was sent
	if mockDiscordCalls != 1 {
		t.Errorf("Expected 1 Discord notification after version change, got %d", mockDiscordCalls)
	}
}
