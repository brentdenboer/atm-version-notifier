package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"atm10-version-notifier/config"
	"atm10-version-notifier/discord"
	"atm10-version-notifier/logger"
	"atm10-version-notifier/reference"
	"atm10-version-notifier/tomlreader"

	"github.com/fsnotify/fsnotify"
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
			modpackCfg.General.ModpackName,
			modpackCfg.General.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to save initial reference to %s: %w", cfg.ReferenceFilePath, err)
		}
		logger.Info("Initial reference file created successfully with version %s", modpackCfg.General.ModpackVersion)
		return nil
	}

	// Compare versions
	if ref.ModpackVersion != modpackCfg.General.ModpackVersion {
		logger.Info("Version change detected for %s! Old: %s, New: %s",
			modpackCfg.General.ModpackName,
			ref.ModpackVersion,
			modpackCfg.General.ModpackVersion)

		// Send Discord notification
		if err := discord.SendVersionUpdateNotification(
			cfg.DiscordWebhookURL,
			modpackCfg.General.ModpackName,
			ref.ModpackVersion,
			modpackCfg.General.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to send Discord notification: %w", err)
		}
		logger.Debug("Successfully sent Discord notification")

		// Update reference file with new version
		if err := reference.SaveReference(
			cfg.ReferenceFilePath,
			modpackCfg.General.ModpackName,
			modpackCfg.General.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to update reference file %s with new version: %w", cfg.ReferenceFilePath, err)
		}

		logger.Info("Successfully updated reference file with new version %s", modpackCfg.General.ModpackVersion)
	} else {
		logger.Debug("No version change detected. Current version: %s", modpackCfg.General.ModpackVersion)
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

	// Initial version check
	if err := checkVersionAndNotify(cfg); err != nil {
		logger.Error("Initial version check failed: %v", err)
	}

	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("Failed to create fsnotify watcher: %v", err)
	}
	defer watcher.Close()

	// Create a channel to handle errors
	done := make(chan bool)
	var lastEventTime time.Time

	// Start watching for changes in a goroutine
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					done <- true
					return
				}

				// Only process Write and Create events
				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					continue
				}

				// Skip if the event is not for our target file
				if filepath.Clean(event.Name) != filepath.Clean(cfg.ConfigFilePath) {
					continue
				}

				// Debounce events by checking the time since last event
				if time.Since(lastEventTime) < time.Duration(cfg.RecheckTimeoutSeconds)*time.Second {
					logger.Debug("Skipping event, too soon after last check")
					continue
				}

				logger.Info("Config file modification detected. Waiting %d seconds before rechecking...",
					cfg.RecheckTimeoutSeconds)

				// Wait for the specified timeout before rechecking
				time.Sleep(time.Duration(cfg.RecheckTimeoutSeconds) * time.Second)

				// Perform version check
				if err := checkVersionAndNotify(cfg); err != nil {
					logger.Error("Version check failed after file modification: %v", err)
				}

				lastEventTime = time.Now()

			case err, ok := <-watcher.Errors:
				if !ok {
					done <- true
					return
				}
				logger.Error("Error watching file: %v", err)
			}
		}
	}()

	// Add the config file to the watcher
	configDir := filepath.Dir(cfg.ConfigFilePath)
	logger.Info("Starting file monitoring for directory: %s", configDir)

	if err := watcher.Add(configDir); err != nil {
		logger.Fatal("Failed to add config directory to watcher: %v", err)
	}

	// Wait for done signal
	<-done
	logger.Info("Monitoring stopped")
}
