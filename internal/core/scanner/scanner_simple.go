package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/h2non/filetype"

	"shotgun-cli-v3/internal/models"
)

// SimpleConcurrentFileScanner is a simpler implementation using goroutines and channels
type SimpleConcurrentFileScanner struct {
	options ScanOptions
}

// NewSimpleConcurrentFileScanner creates a new simple concurrent file scanner
func NewSimpleConcurrentFileScanner() Scanner {
	return NewSimpleConcurrentFileScannerWithOptions(DefaultScanOptions())
}

// NewSimpleConcurrentFileScannerWithOptions creates a new scanner with custom options
func NewSimpleConcurrentFileScannerWithOptions(options ScanOptions) Scanner {
	return &SimpleConcurrentFileScanner{
		options: options,
	}
}

// ScanDirectory implements Scanner.ScanDirectory using a simpler concurrent approach
func (scfs *SimpleConcurrentFileScanner) ScanDirectory(ctx context.Context, rootPath string) (<-chan ScanResult, error) {
	// Validate and clean the root path
	cleanPath := filepath.Clean(rootPath)
	if !filepath.IsAbs(cleanPath) {
		var err error
		cleanPath, err = filepath.Abs(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve absolute path for %s: %w", rootPath, err)
		}
	}

	// Verify the directory exists
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat directory %s: %w", cleanPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", cleanPath)
	}

	// Create result channel
	resultChan := make(chan ScanResult, scfs.options.BufferSize)

	// Determine worker count
	workerCount := scfs.options.WorkerCount
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}

	// Start scanning in a separate goroutine
	go func() {
		defer close(resultChan)
		
		// Set up timeout context if specified
		scanCtx := ctx
		if scfs.options.Timeout > 0 {
			var cancel context.CancelFunc
			scanCtx, cancel = context.WithTimeout(ctx, scfs.options.Timeout)
			defer cancel()
		}
		
		scfs.scanConcurrent(scanCtx, cleanPath, workerCount, resultChan)
	}()

	return resultChan, nil
}

// ScanDirectorySync implements Scanner.ScanDirectorySync
func (scfs *SimpleConcurrentFileScanner) ScanDirectorySync(ctx context.Context, rootPath string) ([]*models.FileNode, error) {
	resultChan, err := scfs.ScanDirectory(ctx, rootPath)
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

// scanConcurrent performs concurrent directory scanning
func (scfs *SimpleConcurrentFileScanner) scanConcurrent(ctx context.Context, rootPath string, workerCount int, resultChan chan<- ScanResult) {
	// Channel for paths to be processed
	pathChan := make(chan string, scfs.options.BufferSize)
	
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go scfs.worker(ctx, pathChan, resultChan, &wg)
	}

	// Submit the root path to start scanning
	go func() {
		defer close(pathChan)
		scfs.discoverPaths(ctx, rootPath, 0, pathChan)
	}()

	// Wait for all workers to complete
	wg.Wait()
}

// discoverPaths recursively discovers all paths and submits them for processing
func (scfs *SimpleConcurrentFileScanner) discoverPaths(ctx context.Context, path string, depth int, pathChan chan<- string) {
	// Check depth limits
	if scfs.options.MaxDepth > 0 && depth >= scfs.options.MaxDepth {
		return
	}

	select {
	case <-ctx.Done():
		return
	case pathChan <- path:
		// Path submitted successfully
	}

	// If this is a directory, discover its children
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return
		default:
		}

		childPath := filepath.Join(path, entry.Name())
		
		// Submit all paths (files and directories) for processing
		scfs.discoverPaths(ctx, childPath, depth+1, pathChan)
	}
}

// worker processes paths from pathChan and sends results to resultChan
func (scfs *SimpleConcurrentFileScanner) worker(ctx context.Context, pathChan <-chan string, resultChan chan<- ScanResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-pathChan:
			if !ok {
				return
			}

			// Process the path
			result := scfs.processPath(ctx, path)
			
			// Send result only if it's not empty (FileNode != nil or Error != nil)
			if result.FileNode != nil || result.Error != nil {
				select {
				case <-ctx.Done():
					return
				case resultChan <- result:
				}
			}
		}
	}
}

// processPath processes a single file or directory path
func (scfs *SimpleConcurrentFileScanner) processPath(ctx context.Context, path string) ScanResult {
	select {
	case <-ctx.Done():
		return ScanResult{Error: ctx.Err()}
	default:
	}

	// Get file info
	info, err := os.Lstat(path) // Use Lstat to not follow symlinks
	if err != nil {
		return ScanResult{Error: fmt.Errorf("failed to stat %s: %w", path, err)}
	}

	// Handle symbolic links
	if info.Mode()&os.ModeSymlink != 0 && !scfs.options.FollowSymlinks {
		// Skip symbolic links if not following them - return empty result (no error)
		return ScanResult{}
	}

	// Create FileNode
	node := &models.FileNode{
		Path:        path,
		Name:        info.Name(),
		IsDirectory: info.IsDir(),
		Size:        info.Size(),
		ModTime:     info.ModTime(),
	}

	// Detect binary files for regular files
	if !node.IsDirectory && scfs.options.DetectBinary && node.Size < MAX_FILE_SIZE_FOR_DETECTION {
		node.IsBinary = scfs.isBinaryFile(path)
	}

	return ScanResult{FileNode: node}
}

// isBinaryFile determines if a file is binary using filetype detection
func (scfs *SimpleConcurrentFileScanner) isBinaryFile(path string) bool {
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
	_, err = filetype.Match(buffer[:n])
	return err == nil // If filetype can identify it, it's likely binary
}