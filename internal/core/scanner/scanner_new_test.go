package scanner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "default constructor",
			opts: nil,
		},
		{
			name: "with workers option",
			opts: []Option{WithWorkers(4)},
		},
		{
			name: "with custom options",
			opts: []Option{WithOptions(ScanOptions{
				MaxDepth:     5,
				DetectBinary: false,
				BufferSize:   50,
			})},
		},
		{
			name:    "invalid worker count",
			opts:    []Option{WithWorkers(-1)},
			wantErr: true,
		},
		{
			name:    "zero worker count",
			opts:    []Option{WithWorkers(0)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner, err := New(tt.opts...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if scanner == nil {
				t.Error("scanner is nil")
			}
		})
	}
}

func TestScanner_DefaultConfiguration(t *testing.T) {
	scanner, err := New()
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	if scanner.workers != runtime.NumCPU() {
		t.Errorf("expected workers = %d, got %d", runtime.NumCPU(), scanner.workers)
	}

	expectedOptions := DefaultScanOptions()
	if scanner.options.BufferSize != expectedOptions.BufferSize {
		t.Errorf("expected BufferSize = %d, got %d", expectedOptions.BufferSize, scanner.options.BufferSize)
	}
}

func TestScanner_WithCustomOptions(t *testing.T) {
	customOptions := ScanOptions{
		MaxDepth:       10,
		FollowSymlinks: true,
		DetectBinary:   false,
		BufferSize:     200,
		WorkerCount:    8,
		Timeout:        1 * time.Minute,
	}

	scanner, err := New(WithOptions(customOptions))
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	if scanner.workers != customOptions.WorkerCount {
		t.Errorf("expected workers = %d, got %d", customOptions.WorkerCount, scanner.workers)
	}

	if scanner.options.MaxDepth != customOptions.MaxDepth {
		t.Errorf("expected MaxDepth = %d, got %d", customOptions.MaxDepth, scanner.options.MaxDepth)
	}
}

func TestScanner_ScanDirectory_InvalidPath(t *testing.T) {
	scanner, err := New()
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	ctx := context.Background()

	// Test non-existent directory
	_, err = scanner.ScanDirectory(ctx, "/non/existent/path")
	if err == nil {
		t.Error("expected error for non-existent path")
	}

	// Test file instead of directory
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = scanner.ScanDirectory(ctx, tempFile.Name())
	if err == nil {
		t.Error("expected error when scanning a file instead of directory")
	}
}

func TestScanner_ScanDirectory_ValidPath(t *testing.T) {
	// Create a temporary directory structure
	tempDir, cleanup := createNewTestDirectory(t)
	defer cleanup()

	scanner, err := New()
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	ctx := context.Background()
	resultChan, err := scanner.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("failed to scan directory: %v", err)
	}

	var results []ScanResult
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) == 0 {
		t.Error("expected at least one result")
	}

	// Verify we got some files
	foundFiles := false
	for _, result := range results {
		if result.Error != nil {
			t.Errorf("unexpected error in result: %v", result.Error)
		}
		if result.FileNode != nil && !result.FileNode.IsDirectory {
			foundFiles = true
		}
	}

	if !foundFiles {
		t.Error("expected to find at least one file")
	}
}

func TestScanner_ScanDirectorySync(t *testing.T) {
	// Create a temporary directory structure
	tempDir, cleanup := createNewTestDirectory(t)
	defer cleanup()

	scanner, err := New()
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	ctx := context.Background()
	nodes, err := scanner.ScanDirectorySync(ctx, tempDir)
	if err != nil {
		t.Fatalf("failed to scan directory synchronously: %v", err)
	}

	if len(nodes) == 0 {
		t.Error("expected at least one node")
	}

	// Verify we got the expected files
	foundTestFile := false
	for _, node := range nodes {
		if filepath.Base(node.Path) == "test.txt" {
			foundTestFile = true
		}
	}

	if !foundTestFile {
		t.Error("expected to find test.txt file")
	}
}

func TestScanner_ContextCancellation(t *testing.T) {
	tempDir, cleanup := createNewTestDirectory(t)
	defer cleanup()

	scanner, err := New()
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resultChan, err := scanner.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("failed to start scan: %v", err)
	}

	// Should complete quickly due to cancellation
	results := make([]ScanResult, 0)
	for result := range resultChan {
		results = append(results, result)
	}

	// With immediate cancellation, we might get no results or very few
	t.Logf("Got %d results with cancelled context", len(results))
}

// createNewTestDirectory creates a temporary directory with some test files
func createNewTestDirectory(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create some test files and directories
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	subFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatalf("failed to create sub file: %v", err)
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}
