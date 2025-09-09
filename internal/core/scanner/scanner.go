package scanner

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/diogopedro/shotgun/internal/models"
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

