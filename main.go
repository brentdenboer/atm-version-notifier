package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"atm10-version-notifier/config"
	"atm10-version-notifier/discord"
	"atm10-version-notifier/logger"
	"atm10-version-notifier/reference"
	"atm10-version-notifier/tomlreader"
)

// checkVersionAndNotify reads the current modpack version and compares it with the reference.
// It sends a notification if the version has changed.
func checkVersionAndNotify(cfg *config.Config) error {
	// Read the modpack configuration
	modpackCfg, err := tomlreader.ReadConfig(cfg.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to read modpack config from %s: %w", cfg.ConfigFilePath, err)
	}

	// Read the reference file
	ref, err := reference.ReadReference(cfg.ReferenceFilePath)
	if err != nil {
		return fmt.Errorf("error reading reference file %s: %w", cfg.ReferenceFilePath, err)
	}

	// If reference file doesn't exist, create it with current version
	if ref == nil {
		logger.Info("No existing reference file found at %s. Creating initial reference...", cfg.ReferenceFilePath)
		if err := reference.SaveReference(
			cfg.ReferenceFilePath,
			modpackCfg.ModpackName,
			modpackCfg.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to save initial reference to %s: %w", cfg.ReferenceFilePath, err)
		}
		logger.Info("Initial reference file created successfully with version %s", modpackCfg.ModpackVersion)
		return nil
	}

	// Compare versions
	if ref.ModpackVersion != modpackCfg.ModpackVersion {
		logger.Info("Version change detected for %s! Old: %s, New: %s",
			modpackCfg.ModpackName,
			ref.ModpackVersion,
			modpackCfg.ModpackVersion)

		// Send Discord notification
		if err := discord.SendVersionUpdateNotification(
			cfg.DiscordWebhookURL,
			modpackCfg.ModpackName,
			ref.ModpackVersion,
			modpackCfg.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to send Discord notification: %w", err)
		}
		logger.Debug("Successfully sent Discord notification")

		// Update reference file with new version
		if err := reference.SaveReference(
			cfg.ReferenceFilePath,
			modpackCfg.ModpackName,
			modpackCfg.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to update reference file %s with new version: %w", cfg.ReferenceFilePath, err)
		}

		logger.Info("Successfully updated reference file with new version %s", modpackCfg.ModpackVersion)
	} else {
		logger.Debug("No version change detected. Current version: %s", modpackCfg.ModpackVersion)
	}

	return nil
}

func main() {
	// Set up logging to both file and stdout
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Configure logger to write to both stdout and file
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Set log level based on environment variable (default to INFO)
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		logger.SetLevel(logger.DEBUG)
	case "INFO":
		logger.SetLevel(logger.INFO)
	case "WARN":
		logger.SetLevel(logger.WARNING)
	case "ERROR":
		logger.SetLevel(logger.ERROR)
	default:
		logger.SetLevel(logger.INFO)
	}

	logger.Info("ATM10 Modpack Version Monitor started")

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	logger.Info("Configuration loaded successfully. Using config file: %s", cfg.ConfigFilePath)
	logger.Debug("File check interval: %d seconds, Recheck timeout: %d seconds",
		cfg.FileCheckIntervalSeconds,
		cfg.RecheckTimeoutSeconds)

	// Initial version check
	if err := checkVersionAndNotify(cfg); err != nil {
		logger.Error("Initial version check failed: %v", err)
	}

	// Store the last modification time
	lastModTime := time.Now()

	// Start the monitoring loop
	ticker := time.NewTicker(time.Duration(cfg.FileCheckIntervalSeconds) * time.Second)
	defer ticker.Stop()

	logger.Info("Starting file monitoring loop. Checking every %d seconds", cfg.FileCheckIntervalSeconds)

	for range ticker.C {
		// Check if the config file exists and get its modification time
		fileInfo, err := os.Stat(cfg.ConfigFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Warn("Config file %s does not exist", cfg.ConfigFilePath)
			} else {
				logger.Error("Error checking config file %s: %v", cfg.ConfigFilePath, err)
			}
			continue
		}

		// Check if the file has been modified
		if fileInfo.ModTime().After(lastModTime) {
			logger.Info("Config file modification detected at %s. Waiting %d seconds before rechecking...",
				fileInfo.ModTime().Format(time.RFC3339),
				cfg.RecheckTimeoutSeconds)

			// Wait for the specified timeout before rechecking
			time.Sleep(time.Duration(cfg.RecheckTimeoutSeconds) * time.Second)

			// Perform version check
			if err := checkVersionAndNotify(cfg); err != nil {
				logger.Error("Version check failed after file modification: %v", err)
			}

			// Update the last modification time
			lastModTime = fileInfo.ModTime()
		}
	}
}
