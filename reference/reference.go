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
	// Validate file path
	if strings.TrimSpace(filePath) == "" {
		return nil, fmt.Errorf("reference file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Return nil to indicate file doesn't exist
		}
		return nil, fmt.Errorf("error accessing reference file: %w", err)
	}

	// Read file with explicit permissions check
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read reference file: %w", err)
	}

	// Handle empty file case
	if len(data) == 0 {
		return nil, fmt.Errorf("reference file is empty")
	}

	var ref ModpackReference
	if err := json.Unmarshal(data, &ref); err != nil {
		return nil, fmt.Errorf("failed to parse reference file: %w", err)
	}

	// Validate reference data
	if err := validateReference(&ref); err != nil {
		return nil, fmt.Errorf("invalid reference data: %w", err)
	}

	return &ref, nil
}

// SaveReference writes the current modpack details to the reference file
func SaveReference(filePath string, modpackName, modpackVersion string) error {
	// Validate input parameters
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("reference file path cannot be empty")
	}
	if strings.TrimSpace(modpackName) == "" {
		return fmt.Errorf("modpack name cannot be empty")
	}
	if strings.TrimSpace(modpackVersion) == "" {
		return fmt.Errorf("modpack version cannot be empty")
	}

	ref := ModpackReference{
		ModpackName:    strings.TrimSpace(modpackName),
		ModpackVersion: strings.TrimSpace(modpackVersion),
	}

	// Validate reference data
	if err := validateReference(&ref); err != nil {
		return fmt.Errorf("invalid reference data: %w", err)
	}

	// Marshal with indentation for better readability
	data, err := json.MarshalIndent(ref, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reference data: %w", err)
	}

	// Ensure the directory exists with proper permissions
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create reference directory: %w", err)
	}

	// Write file with explicit permissions
	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary reference file: %w", err)
	}

	// Atomic rename for safer file writing
	if err := os.Rename(tempFile, filePath); err != nil {
		// Clean up temporary file if rename fails
		os.Remove(tempFile)
		return fmt.Errorf("failed to save reference file: %w", err)
	}

	return nil
}

// validateReference performs validation on the ModpackReference
func validateReference(ref *ModpackReference) error {
	if ref == nil {
		return fmt.Errorf("reference cannot be nil")
	}

	// Validate ModpackName
	if strings.TrimSpace(ref.ModpackName) == "" {
		return fmt.Errorf("modpack name cannot be empty")
	}
	if len(ref.ModpackName) > 100 {
		return fmt.Errorf("modpack name is too long (max 100 characters)")
	}

	// Validate ModpackVersion
	if strings.TrimSpace(ref.ModpackVersion) == "" {
		return fmt.Errorf("modpack version cannot be empty")
	}
	if len(ref.ModpackVersion) > 50 {
		return fmt.Errorf("modpack version is too long (max 50 characters)")
	}

	return nil
}
