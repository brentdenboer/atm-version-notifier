package reference

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadReference(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "reference_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test cases
	tests := []struct {
		name        string
		content     *ModpackReference
		setupFile   bool
		wantErr     bool
		errContains string
	}{
		{
			name: "valid reference",
			content: &ModpackReference{
				ModpackName:    "All The Mods 10",
				ModpackVersion: "2.32",
			},
			setupFile: true,
			wantErr:   false,
		},
		{
			name:      "non-existent file",
			setupFile: false,
			wantErr:   false, // ReadReference returns nil, nil for non-existent files
		},
		{
			name: "invalid json",
			content: &ModpackReference{
				ModpackName:    "",
				ModpackVersion: "",
			},
			setupFile:   true,
			wantErr:     true,
			errContains: "modpack name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, "reference.json")

			if tt.setupFile && tt.content != nil {
				data, err := json.MarshalIndent(tt.content, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal test data: %v", err)
				}
				if err := os.WriteFile(filePath, data, 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			} else if !tt.setupFile {
				// Ensure the file doesn't exist
				os.Remove(filePath)
			}

			got, err := ReadReference(filePath)
			if tt.wantErr {
				if err == nil {
					t.Error("ReadReference() expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ReadReference() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ReadReference() unexpected error: %v", err)
				return
			}

			if tt.setupFile {
				if got.ModpackName != tt.content.ModpackName {
					t.Errorf("ReadReference() ModpackName = %v, want %v", got.ModpackName, tt.content.ModpackName)
				}
				if got.ModpackVersion != tt.content.ModpackVersion {
					t.Errorf("ReadReference() ModpackVersion = %v, want %v", got.ModpackVersion, tt.content.ModpackVersion)
				}
			} else {
				if got != nil {
					t.Errorf("ReadReference() = %v, want nil for non-existent file", got)
				}
			}
		})
	}
}

func TestSaveReference(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "reference_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		modpackName    string
		modpackVersion string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "valid reference",
			modpackName:    "All The Mods 10",
			modpackVersion: "2.32",
			wantErr:        false,
		},
		{
			name:           "empty modpack name",
			modpackName:    "",
			modpackVersion: "2.32",
			wantErr:        true,
			errContains:    "modpack name cannot be empty",
		},
		{
			name:           "empty modpack version",
			modpackName:    "All The Mods 10",
			modpackVersion: "",
			wantErr:        true,
			errContains:    "modpack version cannot be empty",
		},
		{
			name:           "whitespace modpack name",
			modpackName:    "  ",
			modpackVersion: "2.32",
			wantErr:        true,
			errContains:    "modpack name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, "reference.json")

			err := SaveReference(filePath, tt.modpackName, tt.modpackVersion)
			if tt.wantErr {
				if err == nil {
					t.Error("SaveReference() expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SaveReference() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("SaveReference() unexpected error: %v", err)
				return
			}

			// Verify the saved file
			got, err := ReadReference(filePath)
			if err != nil {
				t.Errorf("Failed to read saved reference: %v", err)
				return
			}

			if got.ModpackName != tt.modpackName {
				t.Errorf("SaveReference() saved ModpackName = %v, want %v", got.ModpackName, tt.modpackName)
			}
			if got.ModpackVersion != tt.modpackVersion {
				t.Errorf("SaveReference() saved ModpackVersion = %v, want %v", got.ModpackVersion, tt.modpackVersion)
			}
		})
	}
}

func TestValidateReference(t *testing.T) {
	tests := []struct {
		name        string
		ref         *ModpackReference
		wantErr     bool
		errContains string
	}{
		{
			name: "valid reference",
			ref: &ModpackReference{
				ModpackName:    "All The Mods 10",
				ModpackVersion: "2.32",
			},
			wantErr: false,
		},
		{
			name:        "nil reference",
			ref:         nil,
			wantErr:     true,
			errContains: "reference cannot be nil",
		},
		{
			name: "empty modpack name",
			ref: &ModpackReference{
				ModpackName:    "",
				ModpackVersion: "2.32",
			},
			wantErr:     true,
			errContains: "modpack name cannot be empty",
		},
		{
			name: "empty modpack version",
			ref: &ModpackReference{
				ModpackName:    "All The Mods 10",
				ModpackVersion: "",
			},
			wantErr:     true,
			errContains: "modpack version cannot be empty",
		},
		{
			name: "too long modpack name",
			ref: &ModpackReference{
				ModpackName:    strings.Repeat("a", 101),
				ModpackVersion: "2.32",
			},
			wantErr:     true,
			errContains: "modpack name is too long",
		},
		{
			name: "too long modpack version",
			ref: &ModpackReference{
				ModpackName:    "All The Mods 10",
				ModpackVersion: strings.Repeat("1", 51),
			},
			wantErr:     true,
			errContains: "modpack version is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReference(tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Error("validateReference() expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateReference() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("validateReference() unexpected error: %v", err)
			}
		})
	}
}
