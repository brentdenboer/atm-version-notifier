package tomlreader

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// ModpackConfig represents the structure of the data we need from bcc-common.toml
type ModpackConfig struct {
	ModpackName    string `toml:"modpackName"`
	ModpackVersion string `toml:"modpackVersion"`
}

// ReadConfig reads the TOML file at the specified path and returns the ModpackConfig
func ReadConfig(filepath string) (*ModpackConfig, error) {
	var config ModpackConfig

	// Attempt to decode the TOML file
	_, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TOML file: %w", err)
	}

	// Validate that required fields are present
	if config.ModpackName == "" {
		return nil, fmt.Errorf("modpackName is missing from config file")
	}
	if config.ModpackVersion == "" {
		return nil, fmt.Errorf("modpackVersion is missing from config file")
	}

	return &config, nil
}
