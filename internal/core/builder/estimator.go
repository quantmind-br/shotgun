package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/diogopedro/shotgun/internal/models"
)

// SizeEstimator provides size estimation functionality
type SizeEstimator struct {
	templateEngine TemplateProcessor
}

// EstimationConfig holds configuration for size estimation
type EstimationConfig struct {
	Template      *models.Template
	Variables     map[string]string
	SelectedFiles []string
	IncludeTree   bool
}

// SizeEstimate contains detailed size breakdown
type SizeEstimate struct {
	TotalSize       int64
	TemplateSize    int64
	FileContentSize int64
	TreeStructSize  int64
	OverheadSize    int64
	WarningLevel    int
}

// ProgressCallback is called during progressive calculation
type ProgressCallback func(processed, total int, currentFile string)

// TemplateProcessor interface for template processing
type TemplateProcessor interface {
	ProcessTemplate(template *models.Template, variables map[string]string) (string, error)
}

// NewSizeEstimator creates a new size estimator
func NewSizeEstimator(templateEngine TemplateProcessor) *SizeEstimator {
	return &SizeEstimator{
		templateEngine: templateEngine,
	}
}

// EstimatePromptSize calculates the estimated total size of the prompt output
func (e *SizeEstimator) EstimatePromptSize(ctx context.Context, config EstimationConfig) (*SizeEstimate, error) {
	estimate := &SizeEstimate{}

	// Calculate template size after variable substitution
	templateSize, err := e.calculateTemplateSize(config.Template, config.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate template size: %w", err)
	}
	estimate.TemplateSize = templateSize

	// Calculate file content size
	fileContentSize, treeStructSize, err := e.calculateFileStructureSize(ctx, config.SelectedFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file structure size: %w", err)
	}
	estimate.FileContentSize = fileContentSize
	estimate.TreeStructSize = treeStructSize

	// Calculate XML and formatting overhead
	estimate.OverheadSize = e.calculateFormattingOverhead(config.SelectedFiles, fileContentSize)

	// Calculate total size
	estimate.TotalSize = estimate.TemplateSize + estimate.FileContentSize + 
		estimate.TreeStructSize + estimate.OverheadSize

	// Set warning level
	estimate.WarningLevel = e.determineWarningLevel(estimate.TotalSize)

	return estimate, nil
}

// CalculateProgressively calculates size with progress callbacks
func (e *SizeEstimator) CalculateProgressively(ctx context.Context, files []string, callback ProgressCallback) (int64, error) {
	var totalSize int64
	
	for i, filePath := range files {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		// Notify progress
		if callback != nil {
			callback(i, len(files), filePath)
		}

		// Get file size
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			// Skip files that can't be accessed
			continue
		}

		if !fileInfo.IsDir() {
			totalSize += fileInfo.Size()
		}
	}

	// Final progress notification
	if callback != nil {
		callback(len(files), len(files), "")
	}

	return totalSize, nil
}

// calculateTemplateSize processes template with variables and calculates size
func (e *SizeEstimator) calculateTemplateSize(template *models.Template, variables map[string]string) (int64, error) {
	if e.templateEngine == nil {
		// Fallback estimation without template processing
		baseSize := int64(len(template.Content))
		
		// Estimate variable expansion
		for key, value := range variables {
			placeholder := "{{." + key + "}}"
			occurrences := int64(strings.Count(template.Content, placeholder))
			expansionDiff := int64(len(value)) - int64(len(placeholder))
			baseSize += occurrences * expansionDiff
		}
		
		return baseSize, nil
	}

	// Process template to get exact size
	processed, err := e.templateEngine.ProcessTemplate(template, variables)
	if err != nil {
		return 0, fmt.Errorf("failed to process template: %w", err)
	}

	return int64(len(processed)), nil
}

// calculateFileStructureSize calculates total size of files and tree structure
func (e *SizeEstimator) calculateFileStructureSize(ctx context.Context, selectedFiles []string) (int64, int64, error) {
	var fileContentSize int64
	var treeStructSize int64

	for _, filePath := range selectedFiles {
		select {
		case <-ctx.Done():
			return 0, 0, ctx.Err()
		default:
		}

		// Get file size
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue // Skip inaccessible files
		}

		if !fileInfo.IsDir() {
			fileContentSize += fileInfo.Size()

			// Calculate tree structure overhead for this file
			treeStructSize += e.calculateTreeStructureOverhead(filePath)
		}
	}

	return fileContentSize, treeStructSize, nil
}

// calculateTreeStructureOverhead estimates ASCII tree character overhead
func (e *SizeEstimator) calculateTreeStructureOverhead(filePath string) int64 {
	// Count directory levels for tree structure
	levels := int64(strings.Count(filePath, string(filepath.Separator)))
	
	// Estimate tree characters: ├── └── │ plus spaces
	// Each level adds approximately 4-8 characters
	treeChars := levels * 6
	
	// Add file path length in tree display
	pathDisplay := int64(len(filePath))
	
	return treeChars + pathDisplay
}

// calculateFormattingOverhead estimates XML tags and escaping overhead
func (e *SizeEstimator) calculateFormattingOverhead(selectedFiles []string, contentSize int64) int64 {
	var overhead int64

	// XML tag overhead per file: <file path="...">content</file>
	for _, filePath := range selectedFiles {
		// Opening tag: <file path="filepath">
		openTag := int64(len(`<file path="`) + len(filePath) + len(`">`))
		
		// Closing tag: </file>
		closeTag := int64(len(`</file>`))
		
		overhead += openTag + closeTag
	}

	// XML escaping expansion factor (estimated 5% increase)
	escapingOverhead := contentSize / 20

	// Additional markdown formatting (headers, code blocks)
	markdownOverhead := int64(len(selectedFiles)) * 50 // Estimated per file

	return overhead + escapingOverhead + markdownOverhead
}

// determineWarningLevel returns warning level based on size
func (e *SizeEstimator) determineWarningLevel(totalSize int64) int {
	const (
		largeSizeThreshold     = 100 * 1024  // 100KB
		veryLargeSizeThreshold = 500 * 1024  // 500KB
		excessiveSizeThreshold = 2048 * 1024 // 2MB
	)

	switch {
	case totalSize >= excessiveSizeThreshold:
		return 3 // Excessive
	case totalSize >= veryLargeSizeThreshold:
		return 2 // Very Large
	case totalSize >= largeSizeThreshold:
		return 1 // Large
	default:
		return 0 // Normal
	}
}