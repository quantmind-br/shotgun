package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDebugScanning(t *testing.T) {
	// Create a simple test directory
	tempDir, err := os.MkdirTemp("", "debug_scanner")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Logf("Created test directory: %s", tempDir)
	t.Logf("Created test file: %s", testFile)

	// Test with SimpleConcurrentFileScanner directly
	scanner := NewSimpleConcurrentFileScanner()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := scanner.ScanDirectorySync(ctx, tempDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	t.Logf("Found %d results", len(results))
	for i, result := range results {
		t.Logf("Result %d: Path=%s, Name=%s, IsDir=%t, Size=%d",
			i, result.Path, result.Name, result.IsDirectory, result.Size)
	}

	if len(results) == 0 {
		t.Error("Expected to find at least the directory and file")
	}
}
