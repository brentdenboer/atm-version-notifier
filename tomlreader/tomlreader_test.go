package tomlreader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "tomlreader_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test cases
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
		wantConfig  *ModpackConfig
	}{
		{
			name: "valid config",
			content: `[general]
modpackName = "All The Mods 10"
modpackVersion = "2.32"`,
			wantErr: false,
			wantConfig: &ModpackConfig{
				General: struct {
					ModpackName    string `toml:"modpackName"`
					ModpackVersion string `toml:"modpackVersion"`
				}{
					ModpackName:    "All The Mods 10",
					ModpackVersion: "2.32",
				},
			},
		},
		{
			name: "missing modpack name",
			content: `[general]
modpackVersion = "2.32"`,
			wantErr:     true,
			errContains: "modpackName is missing",
		},
		{
			name: "missing modpack version",
			content: `[general]
modpackName = "All The Mods 10"`,
			wantErr:     true,
			errContains: "modpackVersion is missing",
		},
		{
			name: "invalid version format",
			content: `[general]
modpackName = "All The Mods 10"
modpackVersion = "invalid.version"`,
			wantErr:     true,
			errContains: "invalid version format",
		},
		{
			name: "whitespace in modpack name",
			content: `[general]
modpackName = " All The Mods 10 "
modpackVersion = "2.32"`,
			wantErr:     true,
			errContains: "contains leading or trailing whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for this test case
			tmpFile := filepath.Join(tmpDir, "test.toml")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Run the test
			got, err := ReadConfig(tmpFile)
			if tt.wantErr {
				if err == nil {
					t.Error("ReadConfig() expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ReadConfig() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ReadConfig() unexpected error: %v", err)
				return
			}

			// Compare the results
			if got.General.ModpackName != tt.wantConfig.General.ModpackName {
				t.Errorf("ReadConfig() ModpackName = %v, want %v", got.General.ModpackName, tt.wantConfig.General.ModpackName)
			}
			if got.General.ModpackVersion != tt.wantConfig.General.ModpackVersion {
				t.Errorf("ReadConfig() ModpackVersion = %v, want %v", got.General.ModpackVersion, tt.wantConfig.General.ModpackVersion)
			}
		})
	}
}

func TestReadConfig_FileErrors(t *testing.T) {
	tests := []struct {
		name        string
		filepath    string
		errContains string
	}{
		{
			name:        "non-existent file",
			filepath:    "nonexistent.toml",
			errContains: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadConfig(tt.filepath)
			if err == nil {
				t.Error("ReadConfig() expected error but got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ReadConfig() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
