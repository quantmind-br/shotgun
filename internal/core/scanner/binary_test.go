package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewBinaryDetector(t *testing.T) {
	detector := NewBinaryDetector()
	if detector == nil {
		t.Error("detector is nil")
	}

	if detector.maxFileSize != MAX_FILE_SIZE_FOR_DETECTION {
		t.Errorf("expected maxFileSize = %d, got %d", MAX_FILE_SIZE_FOR_DETECTION, detector.maxFileSize)
	}
}

func TestNewBinaryDetectorWithMaxSize(t *testing.T) {
	customSize := int64(512 * 1024) // 512KB
	detector := NewBinaryDetectorWithMaxSize(customSize)

	if detector.maxFileSize != customSize {
		t.Errorf("expected maxFileSize = %d, got %d", customSize, detector.maxFileSize)
	}
}

func TestBinaryDetector_IsBinary(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "binary_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	detector := NewBinaryDetector()

	// Create test files
	textFile := filepath.Join(tempDir, "text.txt")
	if err := os.WriteFile(textFile, []byte("This is plain text content"), 0644); err != nil {
		t.Fatalf("failed to create text file: %v", err)
	}

	// Create a file with null bytes (binary indicator)
	binaryFile := filepath.Join(tempDir, "binary.bin")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE}
	if err := os.WriteFile(binaryFile, binaryContent, 0644); err != nil {
		t.Fatalf("failed to create binary file: %v", err)
	}

	// Create an empty file
	emptyFile := filepath.Join(tempDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	// Create a JSON file (should be text)
	jsonFile := filepath.Join(tempDir, "data.json")
	jsonContent := `{"name": "test", "value": 123}`
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to create JSON file: %v", err)
	}

	tests := []struct {
		name     string
		file     string
		expected bool
	}{
		{
			name:     "plain text file",
			file:     textFile,
			expected: false,
		},
		{
			name:     "binary file with null bytes",
			file:     binaryFile,
			expected: true,
		},
		{
			name:     "empty file",
			file:     emptyFile,
			expected: false,
		},
		{
			name:     "JSON file",
			file:     jsonFile,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.IsBinary(tt.file)
			if result != tt.expected {
				t.Errorf("IsBinary(%s) = %v, expected %v", tt.file, result, tt.expected)
			}
		})
	}
}

func TestBinaryDetector_IsBinary_NonExistentFile(t *testing.T) {
	detector := NewBinaryDetector()

	// Should return false for non-existent files
	result := detector.IsBinary("/non/existent/file.txt")
	if result != false {
		t.Errorf("expected false for non-existent file, got %v", result)
	}
}

func TestBinaryDetector_IsBinary_LargeFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "binary_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a detector with very small max size
	detector := NewBinaryDetectorWithMaxSize(10)

	// Create a file larger than max size
	largeFile := filepath.Join(tempDir, "large.txt")
	largeContent := make([]byte, 100) // 100 bytes > 10 bytes max
	for i := range largeContent {
		largeContent[i] = byte('a' + (i % 26))
	}
	if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
		t.Fatalf("failed to create large file: %v", err)
	}

	// Should return false due to size limit
	result := detector.IsBinary(largeFile)
	if result != false {
		t.Errorf("expected false for large file, got %v", result)
	}
}

func TestBinaryDetector_ContainsNullBytes(t *testing.T) {
	detector := NewBinaryDetector()

	tests := []struct {
		name     string
		buffer   []byte
		expected bool
	}{
		{
			name:     "no null bytes",
			buffer:   []byte("hello world"),
			expected: false,
		},
		{
			name:     "with null byte at start",
			buffer:   []byte{0x00, 0x48, 0x65, 0x6C, 0x6C, 0x6F},
			expected: true,
		},
		{
			name:     "with null byte in middle",
			buffer:   []byte{0x48, 0x65, 0x00, 0x6C, 0x6C, 0x6F},
			expected: true,
		},
		{
			name:     "with null byte at end",
			buffer:   []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00},
			expected: true,
		},
		{
			name:     "empty buffer",
			buffer:   []byte{},
			expected: false,
		},
		{
			name:     "all null bytes",
			buffer:   []byte{0x00, 0x00, 0x00},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.containsNullBytes(tt.buffer)
			if result != tt.expected {
				t.Errorf("containsNullBytes(%v) = %v, expected %v", tt.buffer, result, tt.expected)
			}
		})
	}
}

func TestBinaryDetector_GetMaxFileSize(t *testing.T) {
	customSize := int64(256 * 1024)
	detector := NewBinaryDetectorWithMaxSize(customSize)

	if detector.GetMaxFileSize() != customSize {
		t.Errorf("expected GetMaxFileSize() = %d, got %d", customSize, detector.GetMaxFileSize())
	}
}

func TestBinaryDetector_IsBinary_PermissionDenied(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "binary_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	detector := NewBinaryDetector()

	// Create a file and make it unreadable
	unreadableFile := filepath.Join(tempDir, "unreadable.txt")
	if err := os.WriteFile(unreadableFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// Change permissions to make it unreadable
	if err := os.Chmod(unreadableFile, 0000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}
	defer os.Chmod(unreadableFile, 0644) // Restore for cleanup

	// Should return false when can't read file
	result := detector.IsBinary(unreadableFile)
	if result != false {
		t.Errorf("expected false for unreadable file, got %v", result)
	}
}
