package reference

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ModpackReference represents the structure of our reference file
type ModpackReference struct {
	ModpackName    string `json:"modpackName"`
	ModpackVersion string `json:"modpackVersion"`
}

// ReadReference reads the current reference data from the file
func ReadReference(filePath string) (*ModpackReference, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil // Return nil to indicate file doesn't exist
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read reference file: %w", err)
	}

	var ref ModpackReference
	if err := json.Unmarshal(data, &ref); err != nil {
		return nil, fmt.Errorf("failed to parse reference file: %w", err)
	}

	return &ref, nil
}

// SaveReference writes the current modpack details to the reference file
func SaveReference(filePath string, modpackName, modpackVersion string) error {
	ref := ModpackReference{
		ModpackName:    modpackName,
		ModpackVersion: modpackVersion,
	}

	data, err := json.MarshalIndent(ref, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reference data: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create reference directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write reference file: %w", err)
	}

	return nil
}
