package tomlreader

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// ModpackConfig represents the structure of the data we need from bcc-common.toml
type ModpackConfig struct {
	General struct {
		ModpackName    string `toml:"modpackName"`
		ModpackVersion string `toml:"modpackVersion"`
	} `toml:"general"`
}

// ReadConfig reads the TOML file at the specified path and returns the ModpackConfig
func ReadConfig(filepath string) (*ModpackConfig, error) {
	// Check if file exists and is readable
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", filepath)
		}
		return nil, fmt.Errorf("error accessing config file: %w", err)
	}

	var config ModpackConfig

	// Attempt to decode the TOML file
	meta, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TOML file: %w", err)
	}

	// Check for undecoded keys
	if len(meta.Undecoded()) > 0 {
		var keys []string
		for _, key := range meta.Undecoded() {
			keys = append(keys, key.String())
		}
		// Log warning about undecoded keys but don't fail
		fmt.Printf("Warning: Found undecoded keys in TOML file: %s\n", strings.Join(keys, ", "))
	}

	// Validate the config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// validateConfig performs validation on the ModpackConfig
func validateConfig(config *ModpackConfig) error {
	// Validate ModpackName
	if config.General.ModpackName == "" {
		return fmt.Errorf("modpackName is missing from config file")
	}
	if len(config.General.ModpackName) > 100 {
		return fmt.Errorf("modpackName is too long (max 100 characters)")
	}
	if strings.TrimSpace(config.General.ModpackName) != config.General.ModpackName {
		return fmt.Errorf("modpackName contains leading or trailing whitespace")
	}

	// Validate ModpackVersion
	if config.General.ModpackVersion == "" {
		return fmt.Errorf("modpackVersion is missing from config file")
	}

	// Check version format (should be in format like "1.0.0" or "v1.0.0")
	versionPattern := regexp.MustCompile(`^v?\d+(\.\d+)*(-[a-zA-Z0-9]+)?$`)
	if !versionPattern.MatchString(config.General.ModpackVersion) {
		return fmt.Errorf("invalid version format: %s (should be like '1.0.0' or 'v1.0.0')", config.General.ModpackVersion)
	}

	return nil
}
