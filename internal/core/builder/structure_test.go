package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFileStructureBuilder(t *testing.T) {
	builder := NewFileStructureBuilder()
	if builder == nil {
		t.Error("Expected non-nil builder")
	}

	if builder.maxFileSize != 10*1024*1024 {
		t.Errorf("Expected default max file size 10MB, got %d", builder.maxFileSize)
	}

	if builder.maxConcurrency != 10 {
		t.Errorf("Expected default max concurrency 10, got %d", builder.maxConcurrency)
	}

	if builder.binaryDetector == nil {
		t.Error("Expected non-nil binary detector")
	}

	if len(builder.sensitiveRegex) == 0 {
		t.Error("Expected sensitive patterns to be initialized")
	}
}

func TestFileStructureBuilder_SetMaxFileSize(t *testing.T) {
	builder := NewFileStructureBuilder()

	// Test valid size
	err := builder.SetMaxFileSize(1024)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	builder.mu.RLock()
	size := builder.maxFileSize
	builder.mu.RUnlock()

	if size != 1024 {
		t.Errorf("Expected max file size 1024, got %d", size)
	}

	// Test invalid size
	err = builder.SetMaxFileSize(0)
	if err == nil {
		t.Error("Expected error for zero file size")
	}

	err = builder.SetMaxFileSize(-1)
	if err == nil {
		t.Error("Expected error for negative file size")
	}
}

func TestFileStructureBuilder_SetMaxConcurrency(t *testing.T) {
	builder := NewFileStructureBuilder()

	// Test valid concurrency
	err := builder.SetMaxConcurrency(5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	builder.mu.RLock()
	concurrency := builder.maxConcurrency
	builder.mu.RUnlock()

	if concurrency != 5 {
		t.Errorf("Expected max concurrency 5, got %d", concurrency)
	}

	// Test invalid concurrency
	err = builder.SetMaxConcurrency(0)
	if err == nil {
		t.Error("Expected error for zero concurrency")
	}

	err = builder.SetMaxConcurrency(-1)
	if err == nil {
		t.Error("Expected error for negative concurrency")
	}
}

func TestFileStructureBuilder_SetTreeFormat(t *testing.T) {
	builder := NewFileStructureBuilder()

	format := TreeFormat{
		UseUnicode: true,
		ShowSizes:  true,
		ShowBinary: false,
		IndentSize: 2,
	}

	err := builder.SetTreeFormat(format)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	builder.mu.RLock()
	actualFormat := builder.treeFormat
	builder.mu.RUnlock()

	if actualFormat.UseUnicode != true ||
		actualFormat.ShowSizes != true ||
		actualFormat.ShowBinary != false ||
		actualFormat.IndentSize != 2 {
		t.Errorf("Tree format not set correctly: %+v", actualFormat)
	}
}

func TestFileStructureBuilder_GenerateStructure_EmptyFiles(t *testing.T) {
	builder := NewFileStructureBuilder()
	ctx := context.Background()

	result, err := builder.GenerateStructure(ctx, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result, got: %s", result)
	}
}

func TestFileStructureBuilder_buildDirectoryTree(t *testing.T) {
	builder := NewFileStructureBuilder()

	files := []string{
		"src/main.go",
		"src/utils/helper.go",
		"docs/README.md",
		"test/main_test.go",
	}

	tree := builder.buildDirectoryTree(files)

	if tree == nil {
		t.Error("Expected non-nil tree")
	}

	// Check root has expected children
	if len(tree.Children) != 3 {
		t.Errorf("Expected 3 root children, got %d", len(tree.Children))
	}

	// Check src directory
	src, exists := tree.Children["src"]
	if !exists {
		t.Error("Expected 'src' directory")
	}
	if !src.IsDirectory {
		t.Error("Expected 'src' to be a directory")
	}
	if len(src.Children) != 2 {
		t.Errorf("Expected 2 children in src, got %d", len(src.Children))
	}

	// Check src/main.go file
	mainGo, exists := src.Children["main.go"]
	if !exists {
		t.Error("Expected 'main.go' file in src")
	}
	if mainGo.IsDirectory {
		t.Error("Expected 'main.go' to be a file")
	}
	if !mainGo.IsFile {
		t.Error("Expected 'main.go' to have IsFile=true")
	}

	// Check src/utils directory
	utils, exists := src.Children["utils"]
	if !exists {
		t.Error("Expected 'utils' directory in src")
	}
	if !utils.IsDirectory {
		t.Error("Expected 'utils' to be a directory")
	}
}

func setupTestFiles(t *testing.T) (string, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "file_structure_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test file structure
	files := map[string]string{
		"simple.txt":        "Hello, World!",
		"config.json":       `{"name": "test", "version": "1.0"}`,
		"src/main.go":       "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}",
		"src/utils/math.go": "package utils\n\nfunc Add(a, b int) int {\n\treturn a + b\n}",
		"docs/README.md":    "# Test Project\n\nThis is a test.",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tempDir, relPath)
		dir := filepath.Dir(fullPath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tempDir)
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			os.RemoveAll(tempDir)
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestFileStructureBuilder_GenerateStructure_WithFiles(t *testing.T) {
	tempDir, cleanup := setupTestFiles(t)
	defer cleanup()

	builder := NewFileStructureBuilder()
	ctx := context.Background()

	files := []string{
		filepath.Join(tempDir, "simple.txt"),
		filepath.Join(tempDir, "src", "main.go"),
	}

	result, err := builder.GenerateStructure(ctx, files)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Check for tree structure
	if !strings.Contains(result, "├──") && !strings.Contains(result, "└──") {
		t.Error("Expected tree characters in result")
	}

	// Check for file content wrapping
	if !strings.Contains(result, "<file path=") || !strings.Contains(result, "</file>") {
		t.Error("Expected XML file wrapping in result")
	}

	// Check for actual file content
	if !strings.Contains(result, "Hello, World!") {
		t.Error("Expected file content 'Hello, World!' in result")
	}

	if !strings.Contains(result, "package main") {
		t.Error("Expected Go code content in result")
	}
}

func TestFileStructureBuilder_GenerateStructure_ContextCancellation(t *testing.T) {
	tempDir, cleanup := setupTestFiles(t)
	defer cleanup()

	builder := NewFileStructureBuilder()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	files := []string{
		filepath.Join(tempDir, "simple.txt"),
	}

	result, err := builder.GenerateStructure(ctx, files)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got: %s", result)
	}
}

func TestFileStructureBuilder_GenerateStructure_LargeFiles(t *testing.T) {
	// Create temporary file that exceeds default limit
	tempDir, err := os.MkdirTemp("", "large_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	builder := NewFileStructureBuilder()
	// Set small limit for testing
	_ = builder.SetMaxFileSize(10)

	// Create large file content (> 10 bytes)
	largeContent := strings.Repeat("A", 20)
	largeFile := filepath.Join(tempDir, "large.txt")
	if err := os.WriteFile(largeFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to write large file: %v", err)
	}

	ctx := context.Background()
	result, err := builder.GenerateStructure(ctx, []string{largeFile})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should contain size limit message
	if !strings.Contains(result, "File too large") {
		t.Error("Expected 'File too large' message for oversized file")
	}
}

func TestFileStructureBuilder_readFileContent_XMLEscaping(t *testing.T) {
	// Create file with XML special characters
	tempDir, err := os.MkdirTemp("", "xml_escape_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	builder := NewFileStructureBuilder()
	ctx := context.Background()

	// Content with XML special characters
	xmlContent := `<tag attr="value">Content & more</tag>`
	xmlFile := filepath.Join(tempDir, "xml.txt")
	if err := os.WriteFile(xmlFile, []byte(xmlContent), 0644); err != nil {
		t.Fatalf("Failed to write XML file: %v", err)
	}

	content, err := builder.readFileContent(ctx, xmlFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that XML characters are escaped
	if strings.Contains(content, "<tag") && !strings.Contains(content, "&lt;tag") {
		t.Error("XML < characters should be escaped")
	}

	if strings.Contains(content, "& more") && !strings.Contains(content, "&amp; more") {
		t.Error("XML & characters should be escaped")
	}

	// Content should be properly escaped
	if !strings.Contains(content, "&lt;tag") {
		t.Errorf("Expected escaped '<' character, got: %s", content)
	}

	if !strings.Contains(content, "&amp; more") {
		t.Errorf("Expected escaped '&' character, got: %s", content)
	}
}

func TestFileStructureBuilder_ConcurrentFileReading(t *testing.T) {
	tempDir, cleanup := setupTestFiles(t)
	defer cleanup()

	builder := NewFileStructureBuilder()
	_ = builder.SetMaxConcurrency(2) // Test with limited concurrency
	ctx := context.Background()

	// Get all test files
	files := []string{
		filepath.Join(tempDir, "simple.txt"),
		filepath.Join(tempDir, "config.json"),
		filepath.Join(tempDir, "src", "main.go"),
		filepath.Join(tempDir, "src", "utils", "math.go"),
		filepath.Join(tempDir, "docs", "README.md"),
	}

	result, err := builder.GenerateStructure(ctx, files)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that all files are included
	for _, file := range files {
		basename := filepath.Base(file)
		if !strings.Contains(result, basename) {
			t.Errorf("Expected file %s to be in result", basename)
		}
	}

	// Check that all content is included
	expectedContents := []string{
		"Hello, World!",
		"test",
		"package main",
		"package utils",
		"Test Project",
	}

	for _, content := range expectedContents {
		if !strings.Contains(result, content) {
			t.Errorf("Expected content '%s' to be in result", content)
		}
	}
}

func TestFileStructureBuilder_SensitiveFileDetection(t *testing.T) {
	builder := NewFileStructureBuilder()

	tests := []struct {
		name         string
		fileName     string
		shouldDetect bool
		description  string
	}{
		{"env file", ".env", true, "Environment file"},
		{"env local", ".env.local", true, "Local environment file"},
		{"private key", "id_rsa", true, "SSH private key"},
		{"pem cert", "certificate.pem", true, "PEM certificate"},
		{"keystore", "app.keystore", true, "Java keystore"},
		{"aws config", ".aws/credentials", true, "AWS credentials"},
		{"secrets", "secrets.json", true, "Secrets file"},
		{"password file", "password.txt", true, "Password file"},
		{"normal file", "README.md", false, "Normal markdown file"},
		{"source file", "main.go", false, "Source code file"},
		{"config file", "app.yaml", false, "Normal config file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detected := builder.isSensitiveFile(tt.fileName)
			if detected != tt.shouldDetect {
				t.Errorf("isSensitiveFile(%s) = %v; expected %v (%s)",
					tt.fileName, detected, tt.shouldDetect, tt.description)
			}
		})
	}
}

func TestFileStructureBuilder_SensitiveFileWarning(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sensitive_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	builder := NewFileStructureBuilder()
	ctx := context.Background()

	// Create a sensitive file
	envFile := filepath.Join(tempDir, ".env")
	envContent := "API_KEY=secret123\nDB_PASSWORD=password456"
	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create env file: %v", err)
	}

	result, err := builder.GenerateStructure(ctx, []string{envFile})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should contain sensitive file warning
	if !strings.Contains(result, "⚠️ Potentially sensitive file detected") {
		t.Error("Expected sensitive file warning for .env file")
	}

	// Should not contain actual sensitive content
	if strings.Contains(result, "API_KEY=secret123") {
		t.Error("Should not include actual sensitive file content")
	}
}

func TestFileStructureBuilder_initSensitivePatterns(t *testing.T) {
	patterns := initSensitivePatterns()

	if len(patterns) == 0 {
		t.Error("Expected at least one sensitive pattern")
	}

	// Test a few key patterns
	testCases := []string{
		".env",
		"id_rsa",
		"certificate.pem",
		"config/app.conf",
	}

	for _, testCase := range testCases {
		found := false
		for _, pattern := range patterns {
			if pattern.MatchString(testCase) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pattern not matching expected file: %s", testCase)
		}
	}
}

func TestFileStructureBuilder_ErrorHandling(t *testing.T) {
	builder := NewFileStructureBuilder()
	ctx := context.Background()

	// Test with non-existent file
	nonExistentFile := "/path/that/does/not/exist.txt"
	result, err := builder.GenerateStructure(ctx, []string{nonExistentFile})
	if err != nil {
		t.Errorf("GenerateStructure should handle non-existent files gracefully, got error: %v", err)
	}

	// Should contain error message in result
	if !strings.Contains(result, "ERROR:") {
		t.Error("Expected error message for non-existent file")
	}
}

func TestFileStructureBuilder_ContextCancellationEdgeCases(t *testing.T) {
	tempDir, cleanup := setupTestFiles(t)
	defer cleanup()

	builder := NewFileStructureBuilder()

	// Test context already cancelled
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	files := []string{filepath.Join(tempDir, "simple.txt")}
	result, err := builder.GenerateStructure(cancelledCtx, files)

	if err == nil {
		t.Error("Expected error when context is already cancelled")
	}
	if result != "" {
		t.Error("Expected empty result when context is cancelled")
	}
}

func TestFileStructureBuilder_ThreadSafety(t *testing.T) {
	builder := NewFileStructureBuilder()

	// Test concurrent configuration changes
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Concurrent configuration updates
			_ = builder.SetMaxFileSize(int64(1024 * (id + 1)))
			_ = builder.SetMaxConcurrency(id + 1)
			_ = builder.SetTreeFormat(TreeFormat{
				UseUnicode: id%2 == 0,
				ShowSizes:  id%3 == 0,
			})
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify builder is still functional
	if builder.maxFileSize <= 0 {
		t.Error("Invalid maxFileSize after concurrent updates")
	}
	if builder.maxConcurrency <= 0 {
		t.Error("Invalid maxConcurrency after concurrent updates")
	}
}

func TestFileStructureBuilder_FunctionalOptions(t *testing.T) {
	// Test all functional options
	builder := NewFileStructureBuilder(
		WithMaxFileSize(5*1024*1024),
		WithMaxConcurrency(5),
		WithTreeFormat(TreeFormat{
			UseUnicode: true,
			ShowSizes:  true,
			ShowBinary: false,
			IndentSize: 2,
		}),
	)

	if builder.maxFileSize != 5*1024*1024 {
		t.Errorf("Expected maxFileSize 5MB, got %d", builder.maxFileSize)
	}

	if builder.maxConcurrency != 5 {
		t.Errorf("Expected maxConcurrency 5, got %d", builder.maxConcurrency)
	}

	if !builder.treeFormat.UseUnicode {
		t.Error("Expected UseUnicode to be true")
	}

	if !builder.treeFormat.ShowSizes {
		t.Error("Expected ShowSizes to be true")
	}
}

func TestFileStructureBuilder_EdgeCaseValidation(t *testing.T) {
	builder := NewFileStructureBuilder()

	// Test invalid parameters
	err := builder.SetMaxFileSize(-1)
	if err == nil {
		t.Error("Expected error for negative file size")
	}

	err = builder.SetMaxConcurrency(-1)
	if err == nil {
		t.Error("Expected error for negative concurrency")
	}

	err = builder.SetMaxConcurrency(0)
	if err == nil {
		t.Error("Expected error for zero concurrency")
	}
}

// Benchmark tests for performance verification
func BenchmarkFileStructureBuilder_GenerateStructure(b *testing.B) {
	tempDir, cleanup := setupTestFiles(nil)
	defer cleanup()

	builder := NewFileStructureBuilder()
	ctx := context.Background()

	files := []string{
		filepath.Join(tempDir, "simple.txt"),
		filepath.Join(tempDir, "src", "main.go"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.GenerateStructure(ctx, files)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkFileStructureBuilder_LargeFileSet(b *testing.B) {
	// Create larger set of test files
	tempDir, err := os.MkdirTemp("", "large_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple files
	var files []string
	for i := 0; i < 50; i++ {
		fileName := filepath.Join(tempDir, fmt.Sprintf("file_%d.txt", i))
		content := fmt.Sprintf("File %d content\nLine 2\nLine 3", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
		files = append(files, fileName)
	}

	builder := NewFileStructureBuilder()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.GenerateStructure(ctx, files)
		if err != nil {
			b.Fatalf("Large benchmark failed: %v", err)
		}
	}
}
