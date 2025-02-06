package reference

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModpackReference represents the structure of our reference file
type ModpackReference struct {
	ModpackName    string `json:"modpackName"`
	ModpackVersion string `json:"modpackVersion"`
}

// ReadReference reads the current reference data from the file
func ReadReference(filePath string) (*ModpackReference, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, fmt.Errorf("empty reference file path")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read reference file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty reference file")
	}

	var ref ModpackReference
	if err := json.Unmarshal(data, &ref); err != nil {
		return nil, fmt.Errorf("parse reference file: %w", err)
	}

	if err := validateReference(&ref); err != nil {
		return nil, fmt.Errorf("invalid reference: %w", err)
	}

	return &ref, nil
}

// SaveReference writes the current modpack details to the reference file
func SaveReference(filePath string, modpackName, modpackVersion string) error {
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("empty reference file path")
	}
	if strings.TrimSpace(modpackName) == "" {
		return fmt.Errorf("empty modpack name")
	}
	if strings.TrimSpace(modpackVersion) == "" {
		return fmt.Errorf("empty modpack version")
	}

	ref := ModpackReference{
		ModpackName:    strings.TrimSpace(modpackName),
		ModpackVersion: strings.TrimSpace(modpackVersion),
	}

	if err := validateReference(&ref); err != nil {
		return fmt.Errorf("invalid reference: %w", err)
	}

	data, err := json.MarshalIndent(ref, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal reference data: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create reference directory: %w", err)
	}

	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}

	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("save reference file: %w", err)
	}

	return nil
}

// validateReference performs validation on the ModpackReference
func validateReference(ref *ModpackReference) error {
	if ref == nil {
		return fmt.Errorf("nil reference")
	}

	if strings.TrimSpace(ref.ModpackName) == "" {
		return fmt.Errorf("empty modpack name")
	}
	if len(ref.ModpackName) > 100 {
		return fmt.Errorf("modpack name too long (max 100 chars)")
	}

	if strings.TrimSpace(ref.ModpackVersion) == "" {
		return fmt.Errorf("empty modpack version")
	}
	if len(ref.ModpackVersion) > 50 {
		return fmt.Errorf("modpack version too long (max 50 chars)")
	}

	return nil
}
