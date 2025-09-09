package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestShotgunignoreIntegration verifies end-to-end integration between
// the init command and the file scanner's ignore functionality
func TestShotgunignoreIntegration(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Create test files structure
	testFiles := []string{
		"src/main.go",
		"build/app.exe",
		"node_modules/package/index.js",
		"dist/bundle.js",
		"logs/app.log",
		".DS_Store",
		"coverage.txt",
		"important.txt",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create .shotgunignore file with typical patterns
	shotgunignoreContent := `# Build artifacts
build/
dist/

# Dependencies
node_modules/

# Logs
*.log

# OS files
.DS_Store

# Go specific
coverage.txt
`

	shotgunignorePath := filepath.Join(tempDir, ".shotgunignore")
	if err := os.WriteFile(shotgunignorePath, []byte(shotgunignoreContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize scanner
	scanner := NewSimpleConcurrentFileScanner()

	// Scan directory
	ctx := context.Background()
	resultChan, err := scanner.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Collect results
	var foundFiles []string
	for result := range resultChan {
		if result.Error != nil {
			t.Logf("Scan error: %v", result.Error)
			continue
		}

		// Skip directories, only count files
		if result.FileNode.IsDirectory {
			continue
		}

		// Convert to relative path for easier testing
		relPath, _ := filepath.Rel(tempDir, result.FileNode.Path)
		relPath = filepath.ToSlash(relPath) // Normalize path separators
		foundFiles = append(foundFiles, relPath)
		t.Logf("Found file: %s", relPath)
	}

	// Verify results
	expectedFound := []string{
		"src/main.go",
		"important.txt",
		".shotgunignore", // The ignore file itself should be found
	}

	expectedIgnored := []string{
		"build/app.exe",
		"node_modules/package/index.js",
		"dist/bundle.js",
		"logs/app.log",
		".DS_Store",
		"coverage.txt",
	}

	// Check that expected files are found
	for _, expected := range expectedFound {
		found := false
		for _, actual := range foundFiles {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s was not found", expected)
		}
	}

	// Check that ignored files are NOT found
	for _, ignored := range expectedIgnored {
		for _, actual := range foundFiles {
			if actual == ignored {
				t.Errorf("Expected ignored file %s was found", ignored)
			}
		}
	}

	t.Logf("Found %d files (expected around %d)", len(foundFiles), len(expectedFound))
}

// TestShotgunignoreWithGitignore verifies that .shotgunignore patterns
// are properly merged with .gitignore patterns
func TestShotgunignoreWithGitignore(t *testing.T) {
	tempDir := t.TempDir()

	// Create .gitignore with some patterns
	gitignoreContent := `# Git patterns
*.tmp
debug/
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .shotgunignore with additional patterns
	shotgunignoreContent := `# Shotgun patterns
*.log
build/
`
	shotgunignorePath := filepath.Join(tempDir, ".shotgunignore")
	if err := os.WriteFile(shotgunignorePath, []byte(shotgunignoreContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create test files
	testFiles := []string{
		"main.go",
		"test.tmp",      // Should be ignored by .gitignore
		"app.log",       // Should be ignored by .shotgunignore
		"debug/out.txt", // Should be ignored by .gitignore
		"build/app.exe", // Should be ignored by .shotgunignore
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test ignorer directly
	ignorer, err := NewIgnorer(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"regular file", filepath.Join(tempDir, "main.go"), false},
		{"gitignore tmp", filepath.Join(tempDir, "test.tmp"), true},
		{"shotgunignore log", filepath.Join(tempDir, "app.log"), true},
		{"gitignore debug dir", filepath.Join(tempDir, "debug/out.txt"), true},
		{"shotgunignore build dir", filepath.Join(tempDir, "build/app.exe"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ignorer.IsIgnored(tt.path)
			if result != tt.expected {
				t.Errorf("IsIgnored(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}
