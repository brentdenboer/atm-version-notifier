package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"atm10-version-notifier/config"
	"atm10-version-notifier/discord"
	"atm10-version-notifier/logger"
	"atm10-version-notifier/reference"
	"atm10-version-notifier/tomlreader"

	"github.com/fsnotify/fsnotify"
)

const (
	defaultLogFile       = "app.log"
	configReloadDebounce = 2 * time.Second // Prevents rapid successive reloads
)

func checkVersionAndNotify(cfg *config.Config) error {
	modpackCfg, err := tomlreader.ReadConfig(cfg.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to read modpack config: %v", err)
	}

	ref, err := reference.ReadReference(cfg.ReferenceFilePath)
	if err != nil {
		return fmt.Errorf("failed to read reference file: %v", err)
	}

	if ref == nil {
		logger.Info("Creating initial reference file at %s", cfg.ReferenceFilePath)
		if err := reference.SaveReference(
			cfg.ReferenceFilePath,
			modpackCfg.General.ModpackName,
			modpackCfg.General.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to create initial reference: %v", err)
		}
		logger.Info("Initial reference created with version %s", modpackCfg.General.ModpackVersion)
		return nil
	}

	if ref.ModpackVersion != modpackCfg.General.ModpackVersion {
		logger.Info("Version change detected: %s -> %s", ref.ModpackVersion, modpackCfg.General.ModpackVersion)

		if err := discord.SendVersionUpdateNotification(
			cfg.DiscordWebhookURL,
			modpackCfg.General.ModpackName,
			ref.ModpackVersion,
			modpackCfg.General.ModpackVersion,
		); err != nil {
			return fmt.Errorf("discord notification failed: %v", err)
		}

		if err := reference.SaveReference(
			cfg.ReferenceFilePath,
			modpackCfg.General.ModpackName,
			modpackCfg.General.ModpackVersion,
		); err != nil {
			return fmt.Errorf("failed to update reference: %v", err)
		}
		logger.Info("Reference file updated to version %s", modpackCfg.General.ModpackVersion)
	}

	logger.Debug("Current version: %s", modpackCfg.General.ModpackVersion)
	return nil
}

func setupLogger() io.Closer {
	logFile, err := os.OpenFile(defaultLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal("failed to open log file: %v", err)
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "debug":
		logger.SetLevel(logger.DEBUG)
	case "warn":
		logger.SetLevel(logger.WARNING)
	case "error":
		logger.SetLevel(logger.ERROR)
	default:
		logger.SetLevel(logger.INFO)
	}
	return logFile
}

func main() {
	closer := setupLogger()
	defer func() {
		logger.Info("Shutting down")
		if closer != nil {
			closer.Close()
		}
		logger.Sync()
	}()

	logger.Info("ATM10 Version Monitor starting")

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("configuration error: %v", err)
	}

	logger.Info("Monitoring config file: %s", cfg.ConfigFilePath)
	logger.Debug("Using reference file: %s", cfg.ReferenceFilePath)

	if err := checkVersionAndNotify(cfg); err != nil {
		logger.Fatal("initial version check failed: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("failed to initialize file watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	defer close(done)
	go watchLoop(watcher, cfg, done)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	configDir := filepath.Dir(cfg.ConfigFilePath)
	if err := watcher.Add(configDir); err != nil {
		logger.Fatal("failed to watch config directory %q: %v", configDir, err)
	}
	logger.Info("Watching directory: %s", configDir)

	select {
	case <-done:
	case sig := <-sigChan:
		logger.Info("Received signal: %v. Shutting down...", sig)
		close(done)
	}
	logger.Info("Exiting main function")
}

func watchLoop(watcher *fsnotify.Watcher, cfg *config.Config, done chan bool) {
	var (
		debounceTimer *time.Timer
	)

	for {
		select {
		case <-done:
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if filepath.Clean(event.Name) != filepath.Clean(cfg.ConfigFilePath) {
				continue
			}

			if debounceTimer != nil {
				debounceTimer.Stop()
			}

			debounceTimer = time.AfterFunc(configReloadDebounce, func() {
				logger.Info("Config file changed, rechecking...")
				if err := checkVersionAndNotify(cfg); err != nil {
					logger.Error("version check failed: %v", err)
				}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Error("file watcher error: %v", err)
		}
	}
}
