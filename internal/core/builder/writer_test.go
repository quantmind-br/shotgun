package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewFileWriter(t *testing.T) {
	writer := NewFileWriter()

	if writer == nil {
		t.Fatal("NewFileWriter() returned nil")
	}
}

func TestGenerateFilename(t *testing.T) {
	writer := NewFileWriter()
	timestamp := time.Date(2025, 9, 4, 14, 25, 30, 0, time.UTC)

	filename := writer.GenerateFilename(timestamp)
	expected := "shotgun_prompt_20250904_1425.md"

	if filename != expected {
		t.Errorf("GenerateFilename() = %s, want %s", filename, expected)
	}
}

func TestWritePromptFile_Success(t *testing.T) {
	writer := NewFileWriter()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "shotgun_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content := "# Test Prompt\n\nThis is a test prompt content."

	outputPath, err := writer.WritePromptFile(content, tempDir)
	if err != nil {
		t.Fatalf("WritePromptFile failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputPath)
	}

	// Check file content
	readContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", content, string(readContent))
	}

	// Check filename format
	basename := filepath.Base(outputPath)
	if !strings.HasPrefix(basename, "shotgun_prompt_") || !strings.HasSuffix(basename, ".md") {
		t.Errorf("Invalid filename format: %s", basename)
	}
}

func TestWritePromptFile_EmptyContent(t *testing.T) {
	writer := NewFileWriter()

	tempDir, err := os.MkdirTemp("", "shotgun_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = writer.WritePromptFile("", tempDir)
	if err == nil {
		t.Error("WritePromptFile should fail with empty content")
	}
}

func TestWritePromptFile_CurrentDirectory(t *testing.T) {
	writer := NewFileWriter()

	content := "# Test Prompt\n\nCurrent directory test."

	// Use empty path to test current directory
	outputPath, err := writer.WritePromptFile(content, "")
	if err != nil {
		t.Fatalf("WritePromptFile failed: %v", err)
	}

	// Clean up
	defer os.Remove(outputPath)

	// Check that file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputPath)
	}
}

func TestCheckCollisions_NoCollision(t *testing.T) {
	writer := NewFileWriter()

	tempDir, err := os.MkdirTemp("", "shotgun_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testPath := filepath.Join(tempDir, "test_file.md")

	result := writer.CheckCollisions(testPath)

	if result != testPath {
		t.Errorf("CheckCollisions should return original path when no collision. Expected: %s, Got: %s", testPath, result)
	}
}

func TestCheckCollisions_WithCollision(t *testing.T) {
	writer := NewFileWriter()

	tempDir, err := os.MkdirTemp("", "shotgun_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the original file to simulate collision
	originalPath := filepath.Join(tempDir, "test_file.md")
	err = os.WriteFile(originalPath, []byte("existing content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result := writer.CheckCollisions(originalPath)

	// Should return a path with _1 suffix
	expectedSuffix := "_1.md"
	if !strings.HasSuffix(result, expectedSuffix) {
		t.Errorf("CheckCollisions should add counter suffix. Expected suffix: %s, Got: %s", expectedSuffix, result)
	}

	if result == originalPath {
		t.Error("CheckCollisions should return different path when collision exists")
	}
}

func TestValidateWritePermissions_ValidDirectory(t *testing.T) {
	writer := NewFileWriter()

	tempDir, err := os.MkdirTemp("", "shotgun_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = writer.ValidateWritePermissions(tempDir)
	if err != nil {
		t.Errorf("ValidateWritePermissions failed for valid directory: %v", err)
	}
}

func TestValidateWritePermissions_NonexistentDirectory(t *testing.T) {
	writer := NewFileWriter()

	nonexistentDir := "/nonexistent/directory/path"

	err := writer.ValidateWritePermissions(nonexistentDir)
	if err == nil {
		t.Error("ValidateWritePermissions should fail for nonexistent directory")
	}
}

func TestValidateWritePermissions_FileInsteadOfDirectory(t *testing.T) {
	writer := NewFileWriter()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "shotgun_test_file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	err = writer.ValidateWritePermissions(tempFile.Name())
	if err == nil {
		t.Error("ValidateWritePermissions should fail when path is a file instead of directory")
	}
}

func TestWritePromptFile_Integration(t *testing.T) {
	writer := NewFileWriter()

	tempDir, err := os.MkdirTemp("", "shotgun_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content := "# Integration Test Prompt\n\nThis tests the full write workflow."

	// First write should succeed
	outputPath1, err := writer.WritePromptFile(content, tempDir)
	if err != nil {
		t.Fatalf("First WritePromptFile failed: %v", err)
	}

	// Second write should handle collision
	outputPath2, err := writer.WritePromptFile(content+" Modified", tempDir)
	if err != nil {
		t.Fatalf("Second WritePromptFile failed: %v", err)
	}

	// Paths should be different
	if outputPath1 == outputPath2 {
		t.Error("Second write should have different path due to collision handling")
	}

	// Both files should exist
	if _, err := os.Stat(outputPath1); os.IsNotExist(err) {
		t.Errorf("First output file should exist: %s", outputPath1)
	}

	if _, err := os.Stat(outputPath2); os.IsNotExist(err) {
		t.Errorf("Second output file should exist: %s", outputPath2)
	}
}
