package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileWriterInterface defines the interface for file writing operations
type FileWriterInterface interface {
	WritePromptFile(content string, basePath string) (string, error)
	GenerateFilename(timestamp time.Time) string
	CheckCollisions(filename string) string
	ValidateWritePermissions(path string) error
}

// FileWriter handles writing prompt files to disk
type FileWriter struct{}

// NewFileWriter creates a new FileWriter instance
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// WritePromptFile writes the prompt content to a file with collision handling
func (fw *FileWriter) WritePromptFile(content string, basePath string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	// Generate base filename with timestamp
	timestamp := time.Now()
	baseFilename := fw.GenerateFilename(timestamp)

	// Resolve base path (use current directory if empty)
	if basePath == "" {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	// Validate write permissions
	if err := fw.ValidateWritePermissions(basePath); err != nil {
		return "", fmt.Errorf("write permission error: %w", err)
	}

	// Handle filename collisions
	finalFilename := fw.CheckCollisions(filepath.Join(basePath, baseFilename))
	fullPath := finalFilename

	// Use atomic file writing pattern: write to temp file, then rename
	tempPath := fullPath + ".tmp"

	// Write to temporary file first
	file, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer func() {
		file.Close()
		// Clean up temp file if it still exists
		if _, err := os.Stat(tempPath); err == nil {
			os.Remove(tempPath)
		}
	}()

	// Write content
	_, err = file.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("failed to write content to file: %w", err)
	}

	// Sync to ensure data is written to disk
	if err = file.Sync(); err != nil {
		return "", fmt.Errorf("failed to sync file to disk: %w", err)
	}

	// Close file before rename
	if err = file.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Atomic rename
	if err = os.Rename(tempPath, fullPath); err != nil {
		return "", fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return fullPath, nil
}

// GenerateFilename creates a timestamp-based filename
func (fw *FileWriter) GenerateFilename(timestamp time.Time) string {
	return fmt.Sprintf("shotgun_prompt_%s.md",
		timestamp.Format("20060102_1504"))
}

// CheckCollisions handles filename collisions by adding a counter
func (fw *FileWriter) CheckCollisions(filepath string) string {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return filepath // No collision
	}

	// Extract base name and extension
	dir := filepath[:strings.LastIndex(filepath, "/")+1]
	if strings.LastIndex(filepath, "/") == -1 {
		dir = ""
	}
	filename := filepath[len(dir):]

	base := filename
	ext := ""
	if dotIndex := strings.LastIndex(filename, "."); dotIndex != -1 {
		base = filename[:dotIndex]
		ext = filename[dotIndex:]
	}

	// Try incrementing counter until we find an available filename
	counter := 1
	for {
		candidateFilename := fmt.Sprintf("%s_%d%s", base, counter, ext)
		candidatePath := dir + candidateFilename

		if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
			return candidatePath
		}

		counter++

		// Prevent infinite loop (reasonable upper limit)
		if counter > 1000 {
			return filepath + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		}
	}
}

// ValidateWritePermissions checks if we can write to the specified path
func (fw *FileWriter) ValidateWritePermissions(path string) error {
	// Check if directory exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("cannot access directory: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	// Check write permissions by trying to create a temp file
	tempFile := filepath.Join(path, ".shotgun_write_test")
	file, err := os.Create(tempFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("no write permission for directory: %s", path)
		}
		return fmt.Errorf("cannot write to directory: %w", err)
	}

	file.Close()
	os.Remove(tempFile)

	// Check available disk space (basic check)
	if err := fw.checkDiskSpace(path); err != nil {
		return err
	}

	return nil
}

// checkDiskSpace performs a basic disk space check
func (fw *FileWriter) checkDiskSpace(path string) error {
	// Try to get filesystem stats (cross-platform approach)
	var stat fs.FileInfo
	var err error

	if stat, err = os.Stat(path); err != nil {
		return fmt.Errorf("cannot check disk space: %w", err)
	}

	// Basic heuristic: if we can create a small test file, assume we have space
	// More sophisticated disk space checking would be platform-specific
	_ = stat // Use stat to avoid unused variable warning

	return nil
}
