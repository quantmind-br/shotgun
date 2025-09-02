package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewConcurrentFileScanner(t *testing.T) {
	scanner := NewSimpleConcurrentFileScanner()
	if scanner == nil {
		t.Fatal("NewConcurrentFileScanner returned nil")
	}
	
	// Verify it's the concrete type we expect
	cfs, ok := scanner.(*SimpleConcurrentFileScanner)
	if !ok {
		t.Fatal("NewConcurrentFileScanner did not return *SimpleConcurrentFileScanner")
	}
	
	// Verify default options
	if cfs.options.BufferSize == 0 {
		t.Error("Expected default buffer size to be set")
	}
	if !cfs.options.DetectBinary {
		t.Error("Expected binary detection to be enabled by default")
	}
}

func TestNewSimpleConcurrentFileScannerWithOptions(t *testing.T) {
	options := ScanOptions{
		MaxDepth:     5,
		BufferSize:   50,
		DetectBinary: false,
	}
	
	scanner := NewSimpleConcurrentFileScannerWithOptions(options)
	cfs, ok := scanner.(*SimpleConcurrentFileScanner)
	if !ok {
		t.Fatal("NewSimpleConcurrentFileScannerWithOptions did not return *SimpleConcurrentFileScanner")
	}
	
	if cfs.options.MaxDepth != 5 {
		t.Errorf("Expected MaxDepth=5, got %d", cfs.options.MaxDepth)
	}
	if cfs.options.BufferSize != 50 {
		t.Errorf("Expected BufferSize=50, got %d", cfs.options.BufferSize)
	}
	if cfs.options.DetectBinary {
		t.Error("Expected DetectBinary=false")
	}
}

func TestConcurrentFileScanner_ScanDirectorySync(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test files and directories
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"subdir/file3.txt",
		"subdir/file4.json",
		"subdir/nested/file5.txt",
	}
	
	for _, testFile := range testFiles {
		fullPath := filepath.Join(tempDir, testFile)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", testFile, err)
		}
		
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFile, err)
		}
	}
	
	tests := []struct {
		name        string
		options     ScanOptions
		expectFiles int // Minimum expected files (actual may be more due to directories)
		wantErr     bool
	}{
		{
			name:        "scan with default options",
			options:     DefaultScanOptions(),
			expectFiles: 5,
			wantErr:     false,
		},
		{
			name: "scan with depth limit",
			options: ScanOptions{
				MaxDepth:     2,
				DetectBinary: true,
				BufferSize:   10,
			},
			expectFiles: 3, // Only files at depth 0 and 1
			wantErr:     false,
		},
		{
			name: "scan with binary detection disabled",
			options: ScanOptions{
				DetectBinary: false,
				BufferSize:   10,
			},
			expectFiles: 5,
			wantErr:     false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewSimpleConcurrentFileScannerWithOptions(tt.options)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			results, err := scanner.ScanDirectorySync(ctx, tempDir)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(results) < tt.expectFiles {
				t.Errorf("Expected at least %d files, got %d", tt.expectFiles, len(results))
			}
			
			// Verify all results have required fields
			for _, result := range results {
				if result.Path == "" {
					t.Error("Found result with empty Path")
				}
				if result.Name == "" {
					t.Error("Found result with empty Name")
				}
			}
		})
	}
}

func TestConcurrentFileScanner_ScanDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "subdir/file3.txt"}
	for _, testFile := range testFiles {
		fullPath := filepath.Join(tempDir, testFile)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", testFile, err)
		}
		
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFile, err)
		}
	}
	
	scanner := NewSimpleConcurrentFileScanner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	resultChan, err := scanner.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}
	
	var results []ScanResult
	var errors []error
	
	for result := range resultChan {
		results = append(results, result)
		if result.Error != nil {
			errors = append(errors, result.Error)
		}
	}
	
	if len(errors) > 0 {
		t.Errorf("Got %d errors during scanning: %v", len(errors), errors[0])
	}
	
	if len(results) == 0 {
		t.Error("Expected some scan results, got none")
	}
	
	// Count file nodes (excluding error results)
	fileNodeCount := 0
	for _, result := range results {
		if result.FileNode != nil {
			fileNodeCount++
		}
	}
	
	if fileNodeCount < 3 {
		t.Errorf("Expected at least 3 file nodes, got %d", fileNodeCount)
	}
}

func TestConcurrentFileScanner_InvalidDirectory(t *testing.T) {
	scanner := NewSimpleConcurrentFileScanner()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test with non-existent directory
	_, err := scanner.ScanDirectory(ctx, "/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
	
	// Test with a file instead of directory
	tempFile, err := os.CreateTemp("", "notadir")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	_, err = scanner.ScanDirectory(ctx, tempFile.Name())
	if err == nil {
		t.Error("Expected error when scanning a file as directory")
	}
}

func TestConcurrentFileScanner_ContextCancellation(t *testing.T) {
	// Create a temporary directory with some files
	tempDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a few test files
	for i := 0; i < 5; i++ {
		filename := filepath.Join(tempDir, "file"+string(rune('0'+i))+".txt")
		err = os.WriteFile(filename, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}
	
	scanner := NewSimpleConcurrentFileScanner()
	ctx, cancel := context.WithCancel(context.Background())
	
	resultChan, err := scanner.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}
	
	// Cancel context immediately
	cancel()
	
	// Collect results - should be cancelled quickly
	var results []ScanResult
	timeout := time.After(2 * time.Second)
	
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				// Channel closed, scanning finished
				return
			}
			results = append(results, result)
			
		case <-timeout:
			t.Fatal("Scanning did not respect context cancellation within timeout")
		}
	}
}

func TestDefaultScanOptions(t *testing.T) {
	options := DefaultScanOptions()
	
	if options.MaxDepth != 0 {
		t.Errorf("Expected MaxDepth=0 (unlimited), got %d", options.MaxDepth)
	}
	if options.FollowSymlinks {
		t.Error("Expected FollowSymlinks=false")
	}
	if !options.DetectBinary {
		t.Error("Expected DetectBinary=true")
	}
	if options.BufferSize <= 0 {
		t.Errorf("Expected BufferSize>0, got %d", options.BufferSize)
	}
	if options.WorkerCount != 0 {
		t.Errorf("Expected WorkerCount=0 (auto), got %d", options.WorkerCount)
	}
	if options.Timeout <= 0 {
		t.Errorf("Expected Timeout>0, got %v", options.Timeout)
	}
}