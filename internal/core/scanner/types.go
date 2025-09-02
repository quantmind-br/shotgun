package scanner

import (
	"context"
	"time"

	"github.com/user/shotgun-cli/internal/models"
)

const (
	// MAX_FILE_SIZE_FOR_DETECTION limits binary detection to files under 1MB
	MAX_FILE_SIZE_FOR_DETECTION = 1024 * 1024
)

// ScannerInterface defines the interface for file scanning operations
type ScannerInterface interface {
	// ScanDirectory scans a directory and returns a stream of discovered files
	ScanDirectory(ctx context.Context, rootPath string) (<-chan ScanResult, error)
	
	// ScanDirectorySync scans a directory synchronously and returns all files at once
	ScanDirectorySync(ctx context.Context, rootPath string) ([]*models.FileNode, error)
}

// ScanResult represents the result of a file scan operation
type ScanResult struct {
	FileNode *models.FileNode
	Error    error
}

// ScanOptions configures scanning behavior
type ScanOptions struct {
	// MaxDepth limits directory traversal depth (0 = unlimited)
	MaxDepth int
	
	// FollowSymlinks determines whether to follow symbolic links
	FollowSymlinks bool
	
	// DetectBinary enables binary file detection
	DetectBinary bool
	
	// BufferSize sets the channel buffer size for streaming results
	BufferSize int
	
	// WorkerCount overrides the default worker pool size
	WorkerCount int
	
	// Timeout sets the maximum time for scanning operations
	Timeout time.Duration
}

// DefaultScanOptions returns sensible default scanning options
func DefaultScanOptions() ScanOptions {
	return ScanOptions{
		MaxDepth:       0, // Unlimited
		FollowSymlinks: false,
		DetectBinary:   true,
		BufferSize:     100,
		WorkerCount:    0, // Use runtime.NumCPU()
		Timeout:        5 * time.Minute,
	}
}

// FileInfo contains metadata about a discovered file
type FileInfo struct {
	Path        string
	Name        string
	IsDirectory bool
	IsBinary    bool
	Size        int64
	ModTime     time.Time
	Error       error
}