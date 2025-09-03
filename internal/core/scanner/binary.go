package scanner

import (
	"os"

	"github.com/h2non/filetype"
)

// BinaryDetector handles binary file detection
type BinaryDetector struct {
	maxFileSize int64
}

// NewBinaryDetector creates a new binary detector with default settings
func NewBinaryDetector() *BinaryDetector {
	return &BinaryDetector{
		maxFileSize: MAX_FILE_SIZE_FOR_DETECTION,
	}
}

// NewBinaryDetectorWithMaxSize creates a new binary detector with custom max file size
func NewBinaryDetectorWithMaxSize(maxSize int64) *BinaryDetector {
	return &BinaryDetector{
		maxFileSize: maxSize,
	}
}

// IsBinary determines if a file is binary using filetype detection
func (bd *BinaryDetector) IsBinary(path string) bool {
	// Get file info first to check size
	info, err := os.Stat(path)
	if err != nil {
		// If we can't stat the file, assume it's not binary
		return false
	}

	// Skip detection for large files to avoid performance issues
	if info.Size() > bd.maxFileSize {
		return false
	}

	file, err := os.Open(path)
	if err != nil {
		// If we can't open the file, assume it's not binary
		return false
	}
	defer file.Close()

	// Read first 262 bytes for filetype detection
	buffer := make([]byte, 262)
	n, err := file.Read(buffer)
	if err != nil || n == 0 {
		return false
	}

	// Use filetype library to detect if it's a known binary type
	kind, err := filetype.Match(buffer[:n])
	if err != nil {
		// If filetype detection fails, fall back to simple heuristic
		return bd.containsNullBytes(buffer[:n])
	}

	// If filetype can identify a specific type, it's likely binary
	if kind != filetype.Unknown {
		return true
	}

	// For unknown types, use null byte detection
	return bd.containsNullBytes(buffer[:n])
}

// containsNullBytes checks if the buffer contains null bytes (common indicator of binary files)
func (bd *BinaryDetector) containsNullBytes(buffer []byte) bool {
	for _, b := range buffer {
		if b == 0 {
			return true
		}
	}
	return false
}

// GetMaxFileSize returns the maximum file size for binary detection
func (bd *BinaryDetector) GetMaxFileSize() int64 {
	return bd.maxFileSize
}
