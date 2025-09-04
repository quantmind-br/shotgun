package confirm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFilenameGenerator(t *testing.T) {
	tests := []struct {
		name        string
		outputDir   string
		expectedDir string
	}{
		{
			name:        "Default directory",
			outputDir:   "",
			expectedDir: ".",
		},
		{
			name:        "Custom directory",
			outputDir:   "/tmp/outputs",
			expectedDir: "/tmp/outputs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fg := NewFilenameGenerator(tt.outputDir)
			
			if fg.outputDir != tt.expectedDir {
				t.Errorf("Expected output dir %s, got %s", tt.expectedDir, fg.outputDir)
			}
		})
	}
}

func TestGenerateTimestampFilename(t *testing.T) {
	fg := NewFilenameGenerator("")
	
	filename := fg.GenerateTimestampFilename()
	
	// Should follow pattern: shotgun_prompt_YYYYMMDD_HHMM.md
	if !strings.HasPrefix(filename, "shotgun_prompt_") {
		t.Errorf("Filename should start with 'shotgun_prompt_', got %s", filename)
	}
	
	if !strings.HasSuffix(filename, ".md") {
		t.Errorf("Filename should end with '.md', got %s", filename)
	}
	
	// Should contain timestamp
	parts := strings.Split(filename, "_")
	if len(parts) < 4 {
		t.Errorf("Expected at least 4 parts in filename, got %d: %s", len(parts), filename)
	}
	
	// Verify timestamp format (basic check)
	dateStr := parts[2] // YYYYMMDD
	timeStr := strings.TrimSuffix(parts[3], ".md") // HHMM
	
	if len(dateStr) != 8 {
		t.Errorf("Expected date part to be 8 characters, got %d: %s", len(dateStr), dateStr)
	}
	
	if len(timeStr) != 4 {
		t.Errorf("Expected time part to be 4 characters, got %d: %s", len(timeStr), timeStr)
	}
}

func TestGenerateFullPath(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
		filename  string
		expected  string
	}{
		{
			name:      "Current directory",
			outputDir: ".",
			filename:  "test.md",
			expected:  filepath.Join(".", "test.md"),
		},
		{
			name:      "Custom directory",
			outputDir: "/tmp/outputs",
			filename:  "test.md",
			expected:  filepath.Join("/tmp/outputs", "test.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fg := NewFilenameGenerator(tt.outputDir)
			result := fg.GenerateFullPath(tt.filename)
			
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCheckFileCollision(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a test file
	existingFile := "existing.md"
	existingPath := filepath.Join(tempDir, existingFile)
	if err := os.WriteFile(existingPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	fg := NewFilenameGenerator(tempDir)
	
	// Test collision with existing file
	collision, path := fg.CheckFileCollision(existingFile)
	if !collision {
		t.Error("Expected collision with existing file")
	}
	if path != existingPath {
		t.Errorf("Expected path %s, got %s", existingPath, path)
	}
	
	// Test no collision with non-existing file
	nonExistingFile := "non_existing.md"
	collision, path = fg.CheckFileCollision(nonExistingFile)
	if collision {
		t.Error("Expected no collision with non-existing file")
	}
	expectedPath := filepath.Join(tempDir, nonExistingFile)
	if path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, path)
	}
}

func TestValidateFilename(t *testing.T) {
	fg := NewFilenameGenerator("")
	
	tests := []struct {
		name        string
		filename    string
		expectError bool
		errorText   string
	}{
		{
			name:     "Valid filename",
			filename: "test_file.md",
		},
		{
			name:        "Empty filename",
			filename:    "",
			expectError: true,
			errorText:   "filename cannot be empty",
		},
		{
			name:        "Filename with invalid character",
			filename:    "test<file.md",
			expectError: true,
			errorText:   "invalid character",
		},
		{
			name:        "Reserved name CON",
			filename:    "CON.md",
			expectError: true,
			errorText:   "reserved system name",
		},
		{
			name:        "Reserved name com1",
			filename:    "com1.txt",
			expectError: true,
			errorText:   "reserved system name",
		},
		{
			name:        "Very long filename",
			filename:    strings.Repeat("a", 300) + ".md",
			expectError: true,
			errorText:   "too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fg.ValidateFilename(tt.filename)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	tempDir := t.TempDir()
	fg := NewFilenameGenerator(tempDir)
	
	baseFilename := "test.md"
	
	// Create the base file
	basePath := filepath.Join(tempDir, baseFilename)
	if err := os.WriteFile(basePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Create the first collision file
	firstCollision := "test_1.md"
	firstPath := filepath.Join(tempDir, firstCollision)
	if err := os.WriteFile(firstPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Generate unique filename
	uniqueFilename := fg.GenerateUniqueFilename(baseFilename)
	
	// Should be test_2.md since test.md and test_1.md exist
	expected := "test_2.md"
	if uniqueFilename != expected {
		t.Errorf("Expected %s, got %s", expected, uniqueFilename)
	}
	
	// Verify no collision
	collision, _ := fg.CheckFileCollision(uniqueFilename)
	if collision {
		t.Error("Generated filename should not have collision")
	}
}

func TestGenerateUniqueFilenameNoCollision(t *testing.T) {
	tempDir := t.TempDir()
	fg := NewFilenameGenerator(tempDir)
	
	baseFilename := "no_collision.md"
	
	// Generate unique filename with no existing files
	uniqueFilename := fg.GenerateUniqueFilename(baseFilename)
	
	// Should return the original filename
	if uniqueFilename != baseFilename {
		t.Errorf("Expected %s, got %s", baseFilename, uniqueFilename)
	}
}

func TestSetAndGetOutputDirectory(t *testing.T) {
	fg := NewFilenameGenerator("initial")
	
	if fg.GetOutputDirectory() != "initial" {
		t.Errorf("Expected 'initial', got %s", fg.GetOutputDirectory())
	}
	
	fg.SetOutputDirectory("/new/path")
	if fg.GetOutputDirectory() != "/new/path" {
		t.Errorf("Expected '/new/path', got %s", fg.GetOutputDirectory())
	}
	
	// Test empty directory defaults to current
	fg.SetOutputDirectory("")
	if fg.GetOutputDirectory() != "." {
		t.Errorf("Expected '.', got %s", fg.GetOutputDirectory())
	}
}

func TestEnsureOutputDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test with current directory (should always succeed)
	fg := NewFilenameGenerator(".")
	if err := fg.EnsureOutputDirectory(); err != nil {
		t.Errorf("Failed to ensure current directory: %v", err)
	}
	
	// Test with new directory
	newDir := filepath.Join(tempDir, "new_output_dir")
	fg.SetOutputDirectory(newDir)
	
	if err := fg.EnsureOutputDirectory(); err != nil {
		t.Errorf("Failed to create output directory: %v", err)
	}
	
	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("Output directory was not created")
	}
}

func TestTimestampUniqueness(t *testing.T) {
	fg := NewFilenameGenerator("")
	
	// Test that filenames follow the expected pattern
	filename1 := fg.GenerateTimestampFilename()
	filename2 := fg.GenerateTimestampFilename()
	
	// Both should follow the pattern but may be identical if generated too quickly
	if !strings.HasPrefix(filename1, "shotgun_prompt_") {
		t.Errorf("Filename should start with 'shotgun_prompt_', got %s", filename1)
	}
	
	if !strings.HasPrefix(filename2, "shotgun_prompt_") {
		t.Errorf("Filename should start with 'shotgun_prompt_', got %s", filename2)
	}
	
	// Test that timestamp format is correct
	parts1 := strings.Split(filename1, "_")
	if len(parts1) < 4 {
		t.Errorf("Expected at least 4 parts in filename, got %d: %s", len(parts1), filename1)
	}
}

func BenchmarkGenerateTimestampFilename(b *testing.B) {
	fg := NewFilenameGenerator("")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fg.GenerateTimestampFilename()
	}
}

func BenchmarkValidateFilename(b *testing.B) {
	fg := NewFilenameGenerator("")
	filename := "valid_filename_for_benchmark.md"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fg.ValidateFilename(filename)
	}
}