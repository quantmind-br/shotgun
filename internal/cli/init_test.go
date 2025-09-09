package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/diogopedro/shotgun/internal/cli/templates"
)

func TestNewInitCmd(t *testing.T) {
	cmd := NewInitCmd()

	if cmd == nil {
		t.Fatal("NewInitCmd() returned nil")
	}

	if cmd.Use != "init" {
		t.Errorf("expected Use = 'init', got %s", cmd.Use)
	}

	if cmd.Short != "Create a .shotgunignore file" {
		t.Errorf("expected Short description to contain 'Create a .shotgunignore file', got %s", cmd.Short)
	}

	// Check if force flag exists
	forceFlag := cmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("expected --force flag to exist")
	}

	// Check flag shorthand
	if forceFlag.Shorthand != "f" {
		t.Errorf("expected --force flag shorthand to be 'f', got %s", forceFlag.Shorthand)
	}
}

func TestCreateShotgunignore(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T, tempDir string)
		force       bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "create new file success",
			setupFunc:   func(t *testing.T, tempDir string) {},
			force:       false,
			expectError: false,
		},
		{
			name: "file exists without force fails",
			setupFunc: func(t *testing.T, tempDir string) {
				if err := os.WriteFile(filepath.Join(tempDir, ".shotgunignore"), []byte("existing"), 0644); err != nil {
					t.Fatal(err)
				}
			},
			force:       false,
			expectError: true,
			errorMsg:    "already exists",
		},
		{
			name: "file exists with force succeeds",
			setupFunc: func(t *testing.T, tempDir string) {
				if err := os.WriteFile(filepath.Join(tempDir, ".shotgunignore"), []byte("existing"), 0644); err != nil {
					t.Fatal(err)
				}
			},
			force:       true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir := t.TempDir()

			// Change to temp directory
			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.Chdir(originalWd) }()

			if err := os.Chdir(tempDir); err != nil {
				t.Fatal(err)
			}

			// Run setup function
			tt.setupFunc(t, tempDir)

			// Test CreateShotgunignore
			err = CreateShotgunignore(tt.force)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}

				// Verify file was created
				content, err := os.ReadFile(".shotgunignore")
				if err != nil {
					t.Errorf("failed to read created .shotgunignore file: %v", err)
				}

				// Verify content matches template
				expectedContent := templates.ShotgunignoreTemplate
				if string(content) != expectedContent {
					t.Error("created file content doesn't match template")
				}
			}
		})
	}
}

func TestValidateDirectory(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string // Returns directory to test
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid writable directory",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
		},
		{
			name: "read-only directory",
			setupFunc: func(t *testing.T) string {
				tempDir := t.TempDir()
				// Make directory read-only (this test may not work on all systems)
				if err := os.Chmod(tempDir, 0444); err != nil {
					t.Skip("Cannot create read-only directory on this system")
				}
				t.Cleanup(func() {
					_ = os.Chmod(tempDir, 0755) // Restore permissions for cleanup
				})
				return tempDir
			},
			expectError: true,
			errorMsg:    "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := tt.setupFunc(t)

			// Change to test directory
			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.Chdir(originalWd) }()

			if err := os.Chdir(testDir); err != nil {
				if tt.expectError && strings.Contains(err.Error(), "permission denied") {
					t.Skip("Cannot chdir to read-only directory on this system")
				}
				t.Fatal(err)
			}

			// Test ValidateDirectory
			err = ValidateDirectory()

			if tt.expectError {
				if err == nil {
					// On Windows, read-only directories may still be writable for the owner
					// Skip this test on Windows
					if runtime.GOOS == "windows" && tt.name == "read-only directory" {
						t.Skip("Read-only directory test not reliable on Windows")
					}
					t.Error("expected error but got nil")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGenerateTemplate(t *testing.T) {
	template := GenerateTemplate()

	if template == "" {
		t.Error("GenerateTemplate() returned empty string")
	}

	// Check that template matches the one from templates package
	expected := templates.ShotgunignoreTemplate
	if template != expected {
		t.Error("GenerateTemplate() doesn't match templates.ShotgunignoreTemplate")
	}

	// Check that template contains expected sections
	expectedSections := []string{
		"# Build artifacts",
		"# Dependencies",
		"# IDE and editor files",
		"# OS generated files",
		"# Go specific",
	}

	for _, section := range expectedSections {
		if !strings.Contains(template, section) {
			t.Errorf("template missing expected section: %s", section)
		}
	}

	// Check that template contains common patterns
	expectedPatterns := []string{
		"build/",
		"node_modules/",
		".DS_Store",
		"*.log",
		"coverage.txt",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(template, pattern) {
			t.Errorf("template missing expected pattern: %s", pattern)
		}
	}
}
