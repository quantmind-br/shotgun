package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/user/shotgun-cli/internal/models"
)

const (
	// DEFAULT_TIMEOUT for scanning operations
	DEFAULT_TIMEOUT = 5 * time.Minute
)

// Scanner implements the core file scanning functionality with concurrency
type Scanner struct {
	ignorer  *Ignorer
	detector *BinaryDetector
	workers  int
	options  ScanOptions
}

// Option defines functional options for Scanner configuration
type Option func(*Scanner) error

// New creates a new Scanner with functional options pattern
func New(opts ...Option) (*Scanner, error) {
	s := &Scanner{
		workers: runtime.NumCPU(),
		options: DefaultScanOptions(),
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return s, nil
}

// WithWorkers sets the number of worker goroutines
func WithWorkers(workers int) Option {
	return func(s *Scanner) error {
		if workers <= 0 {
			return fmt.Errorf("worker count must be positive, got %d", workers)
		}
		s.workers = workers
		s.options.WorkerCount = workers
		return nil
	}
}

// WithIgnorer sets a custom ignorer
func WithIgnorer(ignorer *Ignorer) Option {
	return func(s *Scanner) error {
		s.ignorer = ignorer
		return nil
	}
}

// WithBinaryDetector sets a custom binary detector
func WithBinaryDetector(detector *BinaryDetector) Option {
	return func(s *Scanner) error {
		s.detector = detector
		return nil
	}
}

// WithOptions sets scan options
func WithOptions(options ScanOptions) Option {
	return func(s *Scanner) error {
		s.options = options
		if options.WorkerCount > 0 {
			s.workers = options.WorkerCount
		}
		return nil
	}
}

// ScanDirectory implements the main scanning functionality
func (s *Scanner) ScanDirectory(ctx context.Context, rootPath string) (<-chan ScanResult, error) {
	// For now, delegate to the working SimpleConcurrentFileScanner implementation
	// while maintaining the new functional options interface
	simpleScanner := NewSimpleConcurrentFileScannerWithOptions(s.options)
	return simpleScanner.ScanDirectory(ctx, rootPath)
}

// ScanDirectorySync scans synchronously and returns all files at once
func (s *Scanner) ScanDirectorySync(ctx context.Context, rootPath string) ([]*models.FileNode, error) {
	resultChan, err := s.ScanDirectory(ctx, rootPath)
	if err != nil {
		return nil, err
	}

	var results []*models.FileNode
	var scanErrors []error

	for result := range resultChan {
		if result.Error != nil {
			scanErrors = append(scanErrors, result.Error)
		} else if result.FileNode != nil {
			results = append(results, result.FileNode)
		}
	}

	// Return partial results with aggregated errors
	if len(scanErrors) > 0 {
		return results, fmt.Errorf("encountered %d errors during scanning: %v", len(scanErrors), scanErrors[0])
	}

	return results, nil
}

// scanDirectoryRecursive performs the actual directory scanning
func (s *Scanner) scanDirectoryRecursive(ctx context.Context, rootPath string, resultChan chan<- ScanResult) {
	// Use a simple approach: walk the directory and process files directly
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			// Send error to result channel and continue
			select {
			case resultChan <- ScanResult{Error: fmt.Errorf("failed to access %s: %w", path, err)}:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil // Continue walking
		}

		// Check if path should be ignored
		if s.ignorer != nil && s.ignorer.IsIgnored(path) {
			if d.IsDir() {
				return filepath.SkipDir // Skip entire directory
			}
			return nil // Skip this file
		}

		// Process the path
		s.processPath(ctx, path, resultChan)
		return nil
	})

	if err != nil && err != context.Canceled {
		select {
		case resultChan <- ScanResult{Error: fmt.Errorf("directory walk failed: %w", err)}:
		case <-ctx.Done():
		}
	}
}

// processPath processes a single file or directory path
func (s *Scanner) processPath(ctx context.Context, path string, resultChan chan<- ScanResult) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	// Get file info
	info, err := os.Lstat(path) // Use Lstat to not follow symlinks
	if err != nil {
		select {
		case resultChan <- ScanResult{Error: fmt.Errorf("failed to stat %s: %w", path, err)}:
		case <-ctx.Done():
		}
		return
	}

	// Handle symbolic links
	if info.Mode()&os.ModeSymlink != 0 && !s.options.FollowSymlinks {
		// Skip symbolic links if not following them
		return
	}

	// Create FileNode
	node := &models.FileNode{
		Path:        path,
		Name:        info.Name(),
		IsDirectory: info.IsDir(),
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		IsIgnored:   s.ignorer != nil && s.ignorer.IsIgnored(path),
	}

	// Detect binary files for regular files
	if !node.IsDirectory && s.options.DetectBinary && s.detector != nil {
		node.IsBinary = s.detector.IsBinary(path)
	}

	// Send result
	select {
	case resultChan <- ScanResult{FileNode: node}:
	case <-ctx.Done():
	}
}
