package confirm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FilenameGenerator handles output filename generation
type FilenameGenerator struct {
	outputDir string
}

// NewFilenameGenerator creates a new filename generator
func NewFilenameGenerator(outputDir string) *FilenameGenerator {
	if outputDir == "" {
		outputDir = "."
	}
	return &FilenameGenerator{
		outputDir: outputDir,
	}
}

// GenerateTimestampFilename generates a timestamp-based filename
func (fg *FilenameGenerator) GenerateTimestampFilename() string {
	timestamp := time.Now().Format("20060102_1504")
	return fmt.Sprintf("shotgun_prompt_%s.md", timestamp)
}

// GenerateFullPath generates the full output path
func (fg *FilenameGenerator) GenerateFullPath(filename string) string {
	return filepath.Join(fg.outputDir, filename)
}

// CheckFileCollision checks if a file already exists and returns collision info
func (fg *FilenameGenerator) CheckFileCollision(filename string) (bool, string) {
	fullPath := fg.GenerateFullPath(filename)

	if _, err := os.Stat(fullPath); err == nil {
		return true, fullPath
	}

	return false, fullPath
}

// ValidateFilename checks if filename is safe and valid
func (fg *FilenameGenerator) ValidateFilename(filename string) error {
	// Check for empty filename
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Check for invalid characters
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	// Check for reserved names on Windows
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	baseFilename := strings.ToUpper(strings.TrimSuffix(filename, filepath.Ext(filename)))
	for _, reserved := range reservedNames {
		if baseFilename == reserved {
			return fmt.Errorf("filename cannot be a reserved system name: %s", reserved)
		}
	}

	// Check path length (reasonable limit)
	if len(filename) > 255 {
		return fmt.Errorf("filename is too long (max 255 characters)")
	}

	return nil
}

// EnsureOutputDirectory creates the output directory if it doesn't exist
func (fg *FilenameGenerator) EnsureOutputDirectory() error {
	if fg.outputDir == "." {
		return nil // Current directory always exists
	}

	return os.MkdirAll(fg.outputDir, 0755)
}

// GenerateUniqueFilename generates a unique filename by appending a counter if needed
func (fg *FilenameGenerator) GenerateUniqueFilename(baseFilename string) string {
	filename := baseFilename
	counter := 1

	for {
		collision, _ := fg.CheckFileCollision(filename)
		if !collision {
			return filename
		}

		// Extract extension and base name
		ext := filepath.Ext(baseFilename)
		baseName := strings.TrimSuffix(baseFilename, ext)

		// Append counter
		filename = fmt.Sprintf("%s_%d%s", baseName, counter, ext)
		counter++

		// Prevent infinite loop
		if counter > 1000 {
			break
		}
	}

	return filename
}

// SetOutputDirectory updates the output directory
func (fg *FilenameGenerator) SetOutputDirectory(dir string) {
	if dir == "" {
		fg.outputDir = "."
	} else {
		fg.outputDir = dir
	}
}

// GetOutputDirectory returns the current output directory
func (fg *FilenameGenerator) GetOutputDirectory() string {
	return fg.outputDir
}
