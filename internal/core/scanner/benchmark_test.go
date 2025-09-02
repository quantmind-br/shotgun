package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/user/shotgun-cli/internal/models"
)

// createTestDirectory creates a directory structure with the specified number of files
func createTestDirectory(b *testing.B, numFiles int) string {
	tempDir, err := os.MkdirTemp("", "scanner_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	
	// Create nested directory structure
	for i := 0; i < numFiles; i++ {
		// Create files in nested directories to simulate real project structure
		subDir := filepath.Join(tempDir, fmt.Sprintf("dir%d", i%10))
		err := os.MkdirAll(subDir, 0755)
		if err != nil {
			b.Fatalf("Failed to create subdirectory: %v", err)
		}
		
		filename := filepath.Join(subDir, fmt.Sprintf("file%d.txt", i))
		content := []byte("This is test content for benchmarking file scanning performance")
		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}
	
	return tempDir
}

func BenchmarkConcurrentFileScanner_Small(b *testing.B) {
	tempDir := createTestDirectory(b, 100) // 100 files
	defer os.RemoveAll(tempDir)
	
	scanner := NewConcurrentFileScanner()
	ctx := context.Background()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		results, err := scanner.ScanDirectorySync(ctx, tempDir)
		if err != nil {
			b.Fatalf("Scan failed: %v", err)
		}
		if len(results) == 0 {
			b.Error("No results returned")
		}
	}
}

func BenchmarkConcurrentFileScanner_Medium(b *testing.B) {
	tempDir := createTestDirectory(b, 1000) // 1000 files
	defer os.RemoveAll(tempDir)
	
	scanner := NewConcurrentFileScanner()
	ctx := context.Background()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		results, err := scanner.ScanDirectorySync(ctx, tempDir)
		if err != nil {
			b.Fatalf("Scan failed: %v", err)
		}
		if len(results) == 0 {
			b.Error("No results returned")
		}
	}
}

func BenchmarkConcurrentFileScanner_Large(b *testing.B) {
	// Only run this benchmark with explicit flag to avoid slow tests
	if testing.Short() {
		b.Skip("Skipping large benchmark in short mode")
	}
	
	tempDir := createTestDirectory(b, 5000) // 5000 files
	defer os.RemoveAll(tempDir)
	
	scanner := NewConcurrentFileScanner()
	ctx := context.Background()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		results, err := scanner.ScanDirectorySync(ctx, tempDir)
		if err != nil {
			b.Fatalf("Scan failed: %v", err)
		}
		if len(results) == 0 {
			b.Error("No results returned")
		}
	}
}

func BenchmarkConcurrentFileScanner_WorkerCounts(b *testing.B) {
	tempDir := createTestDirectory(b, 1000)
	defer os.RemoveAll(tempDir)
	
	workerCounts := []int{1, 2, 4, 8, runtime.NumCPU(), runtime.NumCPU() * 2}
	
	for _, workerCount := range workerCounts {
		b.Run(fmt.Sprintf("workers_%d", workerCount), func(b *testing.B) {
			options := DefaultScanOptions()
			options.WorkerCount = workerCount
			scanner := NewConcurrentFileScannerWithOptions(options)
			ctx := context.Background()
			
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				results, err := scanner.ScanDirectorySync(ctx, tempDir)
				if err != nil {
					b.Fatalf("Scan failed: %v", err)
				}
				if len(results) == 0 {
					b.Error("No results returned")
				}
			}
		})
	}
}

func BenchmarkConcurrentFileScanner_BufferSizes(b *testing.B) {
	tempDir := createTestDirectory(b, 500)
	defer os.RemoveAll(tempDir)
	
	bufferSizes := []int{10, 50, 100, 500, 1000}
	
	for _, bufferSize := range bufferSizes {
		b.Run(fmt.Sprintf("buffer_%d", bufferSize), func(b *testing.B) {
			options := DefaultScanOptions()
			options.BufferSize = bufferSize
			scanner := NewConcurrentFileScannerWithOptions(options)
			ctx := context.Background()
			
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				results, err := scanner.ScanDirectorySync(ctx, tempDir)
				if err != nil {
					b.Fatalf("Scan failed: %v", err)
				}
				if len(results) == 0 {
					b.Error("No results returned")
				}
			}
		})
	}
}

func BenchmarkConcurrentFileScanner_BinaryDetection(b *testing.B) {
	tempDir := createTestDirectory(b, 500)
	defer os.RemoveAll(tempDir)
	
	// Add some binary files to the mix
	binaryContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	for i := 0; i < 50; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("binary%d.png", i))
		err := os.WriteFile(filename, binaryContent, 0644)
		if err != nil {
			b.Fatalf("Failed to create binary test file: %v", err)
		}
	}
	
	b.Run("with_binary_detection", func(b *testing.B) {
		options := DefaultScanOptions()
		options.DetectBinary = true
		scanner := NewConcurrentFileScannerWithOptions(options)
		ctx := context.Background()
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			results, err := scanner.ScanDirectorySync(ctx, tempDir)
			if err != nil {
				b.Fatalf("Scan failed: %v", err)
			}
			if len(results) == 0 {
				b.Error("No results returned")
			}
		}
	})
	
	b.Run("without_binary_detection", func(b *testing.B) {
		options := DefaultScanOptions()
		options.DetectBinary = false
		scanner := NewConcurrentFileScannerWithOptions(options)
		ctx := context.Background()
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			results, err := scanner.ScanDirectorySync(ctx, tempDir)
			if err != nil {
				b.Fatalf("Scan failed: %v", err)
			}
			if len(results) == 0 {
				b.Error("No results returned")
			}
		}
	})
}

func BenchmarkConcurrentFileScanner_StreamingVsSync(b *testing.B) {
	tempDir := createTestDirectory(b, 1000)
	defer os.RemoveAll(tempDir)
	
	scanner := NewConcurrentFileScanner()
	
	b.Run("sync_mode", func(b *testing.B) {
		ctx := context.Background()
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			results, err := scanner.ScanDirectorySync(ctx, tempDir)
			if err != nil {
				b.Fatalf("Scan failed: %v", err)
			}
			if len(results) == 0 {
				b.Error("No results returned")
			}
		}
	})
	
	b.Run("streaming_mode", func(b *testing.B) {
		ctx := context.Background()
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			resultChan, err := scanner.ScanDirectory(ctx, tempDir)
			if err != nil {
				b.Fatalf("Scan failed: %v", err)
			}
			
			count := 0
			for result := range resultChan {
				if result.Error == nil && result.FileNode != nil {
					count++
				}
			}
			
			if count == 0 {
				b.Error("No results returned")
			}
		}
	})
}

func BenchmarkWorkerPool_JobThroughput(b *testing.B) {
	processor := &MockJobProcessor{
		processFunc: func(ctx context.Context, job Job) []ScanResult {
			// Simulate minimal processing time
			return []ScanResult{{FileNode: &models.FileNode{Path: job.Path}}}
		},
	}
	
	workerCounts := []int{1, 2, 4, 8, runtime.NumCPU()}
	
	for _, workerCount := range workerCounts {
		b.Run(fmt.Sprintf("workers_%d", workerCount), func(b *testing.B) {
			wp := NewWorkerPool(workerCount, 100, processor)
			wp.Start()
			defer wp.Stop()
			
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				job := Job{Path: "benchmark_job", Depth: 0}
				wp.SubmitJob(job)
				
				// Consume the result
				<-wp.Results()
			}
		})
	}
}

// BenchmarkFileTypeDetection tests the performance of binary file detection
func BenchmarkFileTypeDetection(b *testing.B) {
	// Create test files of different types
	tempDir, err := os.MkdirTemp("", "filetype_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Text file
	textFile := filepath.Join(tempDir, "test.txt")
	textContent := []byte("This is a text file with some content for testing")
	err = os.WriteFile(textFile, textContent, 0644)
	if err != nil {
		b.Fatalf("Failed to create text file: %v", err)
	}
	
	// Binary file (PNG)
	binaryFile := filepath.Join(tempDir, "test.png")
	pngContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D}
	err = os.WriteFile(binaryFile, pngContent, 0644)
	if err != nil {
		b.Fatalf("Failed to create binary file: %v", err)
	}
	
	processor := &FileProcessor{options: DefaultScanOptions()}
	
	b.Run("text_file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isBinary := processor.isBinaryFile(textFile)
			if isBinary {
				b.Error("Text file detected as binary")
			}
		}
	})
	
	b.Run("binary_file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isBinary := processor.isBinaryFile(binaryFile)
			if !isBinary {
				b.Error("Binary file not detected as binary")
			}
		}
	})
}

