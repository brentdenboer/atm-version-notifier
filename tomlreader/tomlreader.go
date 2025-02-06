package tomlreader

import (
	"fmt"
	"regexp"
	"strings"

	"atm10-version-notifier/logger"

	"github.com/BurntSushi/toml"
)

// ModpackConfig represents modpack configuration data from a TOML file
type ModpackConfig struct {
	General struct {
		ModpackName    string `toml:"modpackName"`
		ModpackVersion string `toml:"modpackVersion"`
	} `toml:"general"`
}

// ReadConfig reads and validates a TOML configuration file
func ReadConfig(filepath string) (*ModpackConfig, error) {
	var config ModpackConfig

	meta, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		var keys []string
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		logger.Warn("Found undecoded keys in TOML config: %s", strings.Join(keys, ", "))
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

var versionRegex = regexp.MustCompile(`^v?\d+(\.\d+)*(-[a-zA-Z0-9]+)?$`)

func validateConfig(config *ModpackConfig) error {
	// Validate ModpackName
	name := strings.TrimSpace(config.General.ModpackName)
	switch {
	case name == "":
		return fmt.Errorf("missing modpackName")
	case len(name) > 100:
		return fmt.Errorf("modpackName too long (100 character max)")
	case name != config.General.ModpackName:
		return fmt.Errorf("modpackName contains leading/trailing whitespace")
	}

	// Validate ModpackVersion
	version := config.General.ModpackVersion
	if version == "" {
		return fmt.Errorf("missing modpackVersion")
	}
	if !versionRegex.MatchString(version) {
		return fmt.Errorf("invalid version format: %q (expected format: 1.0.0 or v1.0.0)", version)
	}

	return nil
}
