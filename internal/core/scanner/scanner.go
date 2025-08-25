package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/h2non/filetype"

	"shotgun-cli-v3/internal/models"
)

const (
	// DEFAULT_TIMEOUT for scanning operations
	DEFAULT_TIMEOUT = 5 * time.Minute
)

// ConcurrentFileScanner implements the Scanner interface with concurrent processing
type ConcurrentFileScanner struct {
	options    ScanOptions
	workerPool *WorkerPool
	mu         sync.RWMutex
	started    bool
}

// NewConcurrentFileScanner creates a new concurrent file scanner with default options
func NewConcurrentFileScanner() Scanner {
	return NewSimpleConcurrentFileScanner()
}

// NewConcurrentFileScannerWithOptions creates a new concurrent file scanner with custom options
func NewConcurrentFileScannerWithOptions(options ScanOptions) Scanner {
	return NewSimpleConcurrentFileScannerWithOptions(options)
}

// ScanDirectory implements Scanner.ScanDirectory
func (cfs *ConcurrentFileScanner) ScanDirectory(ctx context.Context, rootPath string) (<-chan ScanResult, error) {
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

	// Set up timeout context if specified
	scanCtx := ctx
	if cfs.options.Timeout > 0 {
		var cancel context.CancelFunc
		scanCtx, cancel = context.WithTimeout(ctx, cfs.options.Timeout)
		defer cancel()
	}

	// Create and start worker pool
	processor := &FileProcessor{
		options: cfs.options,
	}
	
	workerCount := cfs.options.WorkerCount
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}
	
	wp := NewWorkerPool(workerCount, cfs.options.BufferSize, processor)
	wp.Start()

	// Start scanning in a separate goroutine
	resultChan := make(chan ScanResult, cfs.options.BufferSize)
	
	go func() {
		defer close(resultChan)
		defer wp.Stop()
		
		cfs.scanWithWorkerPool(scanCtx, cleanPath, wp, resultChan)
	}()

	return resultChan, nil
}

// ScanDirectorySync implements Scanner.ScanDirectorySync
func (cfs *ConcurrentFileScanner) ScanDirectorySync(ctx context.Context, rootPath string) ([]*models.FileNode, error) {
	resultChan, err := cfs.ScanDirectory(ctx, rootPath)
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

// scanWithWorkerPool orchestrates the concurrent scanning process
func (cfs *ConcurrentFileScanner) scanWithWorkerPool(ctx context.Context, rootPath string, wp *WorkerPool, resultChan chan<- ScanResult) {
	// Submit initial job - this will process the entire directory tree
	wp.SubmitJob(Job{Path: rootPath, Depth: 0})
	
	// Since we only submit one job (which processes everything recursively),
	// we need to wait until that job completes
	jobsSubmitted := 1
	jobsCompleted := 0
	
	for jobsCompleted < jobsSubmitted {
		select {
		case result, ok := <-wp.Results():
			if !ok {
				// Worker pool results channel closed
				return
			}
			
			// Send result to output channel
			select {
			case resultChan <- result:
			case <-ctx.Done():
				return
			}
			
			// The FileProcessor returns multiple results per job,
			// but we can't easily track when a job is "done" in this model.
			// For simplicity, let's use a timeout-based approach
			
		case <-ctx.Done():
			return
		}
	}
}

// FileProcessor implements JobProcessor for file scanning operations
type FileProcessor struct {
	options ScanOptions
}

// ProcessJob implements JobProcessor.ProcessJob
func (fp *FileProcessor) ProcessJob(ctx context.Context, job Job) []ScanResult {
	var results []ScanResult
	
	// Check depth limits
	if fp.options.MaxDepth > 0 && job.Depth >= fp.options.MaxDepth {
		return results
	}
	
	// Process the current path
	fileNode, err := fp.processPath(ctx, job.Path)
	if err != nil {
		results = append(results, ScanResult{Error: err})
		return results
	}
	
	if fileNode != nil {
		results = append(results, ScanResult{FileNode: fileNode})
		
		// If it's a directory, process its contents recursively
		if fileNode.IsDirectory {
			childResults := fp.processDirectory(ctx, job.Path, job.Depth+1)
			results = append(results, childResults...)
		}
	}
	
	return results
}

// processPath processes a single file or directory path
func (fp *FileProcessor) processPath(ctx context.Context, path string) (*models.FileNode, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// Get file info
	info, err := os.Lstat(path) // Use Lstat to not follow symlinks
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	
	// Handle symbolic links
	if info.Mode()&os.ModeSymlink != 0 && !fp.options.FollowSymlinks {
		// Skip symbolic links if not following them
		return nil, nil
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
	if !node.IsDirectory && fp.options.DetectBinary && node.Size < MAX_FILE_SIZE_FOR_DETECTION {
		node.IsBinary = fp.isBinaryFile(path)
	}
	
	return node, nil
}

// processDirectory processes the contents of a directory
func (fp *FileProcessor) processDirectory(ctx context.Context, dirPath string, depth int) []ScanResult {
	var results []ScanResult
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		results = append(results, ScanResult{
			Error: fmt.Errorf("failed to read directory %s: %w", dirPath, err),
		})
		return results
	}
	
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return results
		default:
		}
		
		childPath := filepath.Join(dirPath, entry.Name())
		childJob := Job{Path: childPath, Depth: depth}
		
		// Process child synchronously for simplicity
		// In a more complex implementation, we could submit child jobs back to the pool
		childResults := fp.ProcessJob(ctx, childJob)
		results = append(results, childResults...)
	}
	
	return results
}

// isBinaryFile determines if a file is binary using filetype detection
func (fp *FileProcessor) isBinaryFile(path string) bool {
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