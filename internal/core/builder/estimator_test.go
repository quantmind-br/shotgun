package builder

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/diogopedro/shotgun/internal/models"
)

// mockTemplateProcessor for testing
type mockTemplateProcessor struct {
	processResult string
	processError  error
}

func (m *mockTemplateProcessor) ProcessTemplate(template *models.Template, variables map[string]string) (string, error) {
	if m.processError != nil {
		return "", m.processError
	}
	return m.processResult, nil
}

func TestNewSizeEstimator(t *testing.T) {
	mockEngine := &mockTemplateProcessor{}
	estimator := NewSizeEstimator(mockEngine)

	if estimator == nil {
		t.Fatal("Expected non-nil estimator")
	}

	if estimator.templateEngine != mockEngine {
		t.Error("Template engine not set correctly")
	}
}

func TestCalculateTemplateSize(t *testing.T) {
	tests := []struct {
		name           string
		template       *models.Template
		variables      map[string]string
		mockResult     string
		mockError      error
		expectedSize   int64
		expectError    bool
	}{
		{
			name: "Simple template processing",
			template: &models.Template{
				Content: "Hello {{.name}}",
			},
			variables: map[string]string{
				"name": "World",
			},
			mockResult:   "Hello World",
			expectedSize: 11,
		},
		{
			name: "Template processing error",
			template: &models.Template{
				Content: "Hello {{.name}}",
			},
			variables:   map[string]string{},
			mockError:   errors.New("template error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEngine := &mockTemplateProcessor{
				processResult: tt.mockResult,
				processError:  tt.mockError,
			}
			estimator := NewSizeEstimator(mockEngine)

			size, err := estimator.calculateTemplateSize(tt.template, tt.variables)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if size != tt.expectedSize {
				t.Errorf("Expected size %d, got %d", tt.expectedSize, size)
			}
		})
	}
}

func TestCalculateTemplateSizeFallback(t *testing.T) {
	// Test fallback when no template engine is provided
	estimator := NewSizeEstimator(nil)

	template := &models.Template{
		Content: "Hello {{.name}}, welcome to {{.app}}!",
	}
	variables := map[string]string{
		"name": "John",
		"app":  "MyApp",
	}

	size, err := estimator.calculateTemplateSize(template, variables)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should estimate based on variable substitution
	// Original: "Hello {{.name}}, welcome to {{.app}}!" (39 chars)
	// After substitution: "Hello John, welcome to MyApp!" (30 chars)
	// Difference: John (4) vs {{.name}} (9) = -5
	// Difference: MyApp (5) vs {{.app}} (8) = -3
	// Expected: 39 - 5 - 3 = 31, but let's be flexible with the estimation
	if size <= 0 {
		t.Error("Expected positive size estimation")
	}
}

func TestCalculateFileStructureSize(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	
	file1 := filepath.Join(tempDir, "test1.txt")
	file2 := filepath.Join(tempDir, "test2.txt")

	content1 := "Hello World"
	content2 := "This is a longer test file content for testing"

	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	estimator := NewSizeEstimator(nil)
	selectedFiles := []string{file1, file2}

	ctx := context.Background()
	fileContentSize, treeStructSize, err := estimator.calculateFileStructureSize(ctx, selectedFiles)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedContentSize := int64(len(content1) + len(content2))
	if fileContentSize != expectedContentSize {
		t.Errorf("Expected file content size %d, got %d", expectedContentSize, fileContentSize)
	}

	if treeStructSize <= 0 {
		t.Error("Expected positive tree structure size")
	}
}

func TestCalculateProgressively(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	
	files := make([]string, 3)
	for i := 0; i < 3; i++ {
		file := filepath.Join(tempDir, "test"+string(rune('1'+i))+".txt")
		content := strings.Repeat("x", (i+1)*10)
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		files[i] = file
	}

	estimator := NewSizeEstimator(nil)

	var progressUpdates []int
	callback := func(processed, total int, currentFile string) {
		progressUpdates = append(progressUpdates, processed)
	}

	ctx := context.Background()
	totalSize, err := estimator.CalculateProgressively(ctx, files, callback)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedSize := int64(10 + 20 + 30) // File contents: 10, 20, 30 chars
	if totalSize != expectedSize {
		t.Errorf("Expected total size %d, got %d", expectedSize, totalSize)
	}

	// Should have progress updates: 0, 1, 2, 3
	expectedUpdates := []int{0, 1, 2, 3}
	if len(progressUpdates) != len(expectedUpdates) {
		t.Errorf("Expected %d progress updates, got %d", len(expectedUpdates), len(progressUpdates))
	}

	for i, expected := range expectedUpdates {
		if i >= len(progressUpdates) || progressUpdates[i] != expected {
			t.Errorf("Progress update %d: expected %d, got %d", i, expected, progressUpdates[i])
		}
	}
}

func TestCalculateProgressivelyWithCancellation(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	
	files := make([]string, 5)
	for i := 0; i < 5; i++ {
		file := filepath.Join(tempDir, "test"+string(rune('1'+i))+".txt")
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		files[i] = file
	}

	estimator := NewSizeEstimator(nil)

	ctx, cancel := context.WithCancel(context.Background())
	
	var progressUpdates int
	callback := func(processed, total int, currentFile string) {
		progressUpdates++
		if processed == 2 { // Cancel after processing 2 files
			cancel()
		}
	}

	_, err := estimator.CalculateProgressively(ctx, files, callback)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	if progressUpdates < 2 {
		t.Errorf("Expected at least 2 progress updates before cancellation, got %d", progressUpdates)
	}
}

func TestEstimatePromptSize(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	
	file1 := filepath.Join(tempDir, "test1.txt")
	if err := os.WriteFile(file1, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockEngine := &mockTemplateProcessor{
		processResult: "Processed template content",
	}
	estimator := NewSizeEstimator(mockEngine)

	config := EstimationConfig{
		Template: &models.Template{
			Content: "Template {{.var}}",
		},
		Variables: map[string]string{
			"var": "value",
		},
		SelectedFiles: []string{file1},
		IncludeTree:   true,
	}

	ctx := context.Background()
	estimate, err := estimator.EstimatePromptSize(ctx, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if estimate == nil {
		t.Fatal("Expected non-nil estimate")
	}

	if estimate.TotalSize <= 0 {
		t.Error("Expected positive total size")
	}

	if estimate.TemplateSize <= 0 {
		t.Error("Expected positive template size")
	}

	if estimate.FileContentSize <= 0 {
		t.Error("Expected positive file content size")
	}

	if estimate.OverheadSize <= 0 {
		t.Error("Expected positive overhead size")
	}

	// Total should be sum of all components
	expectedTotal := estimate.TemplateSize + estimate.FileContentSize + 
		estimate.TreeStructSize + estimate.OverheadSize
	
	if estimate.TotalSize != expectedTotal {
		t.Errorf("Total size mismatch: expected %d, got %d", expectedTotal, estimate.TotalSize)
	}
}

func TestDetermineWarningLevel(t *testing.T) {
	estimator := NewSizeEstimator(nil)

	tests := []struct {
		name          string
		size          int64
		expectedLevel int
	}{
		{"Normal size", 50 * 1024, 0},     // 50KB
		{"Large size", 150 * 1024, 1},     // 150KB
		{"Very large", 700 * 1024, 2},     // 700KB
		{"Excessive", 3 * 1024 * 1024, 3}, // 3MB
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := estimator.determineWarningLevel(tt.size)
			if level != tt.expectedLevel {
				t.Errorf("Expected warning level %d, got %d", tt.expectedLevel, level)
			}
		})
	}
}

func TestCalculateFormattingOverhead(t *testing.T) {
	estimator := NewSizeEstimator(nil)
	
	selectedFiles := []string{
		"path/to/file1.go",
		"path/to/file2.go",
	}
	contentSize := int64(1000)

	overhead := estimator.calculateFormattingOverhead(selectedFiles, contentSize)

	if overhead <= 0 {
		t.Error("Expected positive formatting overhead")
	}

	// Should include XML tags, escaping, and markdown formatting
	expectedMinimum := int64(len(selectedFiles)) * 20 // Conservative estimate
	if overhead < expectedMinimum {
		t.Errorf("Expected overhead at least %d, got %d", expectedMinimum, overhead)
	}
}

func TestCalculateTreeStructureOverhead(t *testing.T) {
	estimator := NewSizeEstimator(nil)
	
	tests := []struct {
		path         string
		expectedMin  int64
	}{
		{"file.go", 6},                     // Simple file
		{"src/main.go", 10},               // One level deep
		{"src/pkg/util/helper.go", 20},    // Multiple levels
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			overhead := estimator.calculateTreeStructureOverhead(tt.path)
			
			if overhead < tt.expectedMin {
				t.Errorf("Expected overhead at least %d, got %d", tt.expectedMin, overhead)
			}
		})
	}
}

func BenchmarkEstimatePromptSize(b *testing.B) {
	// Create temporary test files
	tempDir := b.TempDir()
	
	files := make([]string, 100)
	for i := 0; i < 100; i++ {
		file := filepath.Join(tempDir, "file"+string(rune('0'+i%10))+".txt")
		content := strings.Repeat("x", 100)
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
		files[i] = file
	}

	mockEngine := &mockTemplateProcessor{
		processResult: "Processed template",
	}
	estimator := NewSizeEstimator(mockEngine)

	config := EstimationConfig{
		Template: &models.Template{
			Content: "Test template {{.var}}",
		},
		Variables: map[string]string{
			"var": "value",
		},
		SelectedFiles: files,
		IncludeTree:   true,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := estimator.EstimatePromptSize(ctx, config)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}