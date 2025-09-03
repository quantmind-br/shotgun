package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// getUserTemplatesDir returns the user-specific templates directory
func getUserTemplatesDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
	case "darwin", "linux":
		baseDir = os.Getenv("HOME")
		if baseDir == "" {
			return "", fmt.Errorf("HOME environment variable not set")
		}
		baseDir = filepath.Join(baseDir, ".config")
	default:
		// Fallback for other Unix-like systems
		baseDir = os.Getenv("HOME")
		if baseDir == "" {
			return "", fmt.Errorf("HOME environment variable not set")
		}
		baseDir = filepath.Join(baseDir, ".config")
	}

	templatesDir := filepath.Join(baseDir, "shotgun-cli", "templates")

	// Clean the path to prevent any path traversal issues
	cleanPath := filepath.Clean(templatesDir)

	return cleanPath, nil
}

// ensureTemplateDir creates the template directory if it doesn't exist
func ensureTemplateDir(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Clean the path for security
	cleanPath := filepath.Clean(path)

	// Check if it already exists
	info, err := os.Stat(cleanPath)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", cleanPath)
		}
		return nil // Directory already exists
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory status: %w", err)
	}

	// Create directory with appropriate permissions
	if err := os.MkdirAll(cleanPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", cleanPath, err)
	}

	return nil
}

// GetUserTemplatesDir is a public wrapper for getUserTemplatesDir
func GetUserTemplatesDir() (string, error) {
	return getUserTemplatesDir()
}

// EnsureTemplateDir is a public wrapper for ensureTemplateDir
func EnsureTemplateDir(path string) error {
	return ensureTemplateDir(path)
}

// validatePathSafety performs additional path safety validation
func validatePathSafety(basePath, targetPath string) error {
	// Clean both paths
	cleanBase := filepath.Clean(basePath)
	cleanTarget := filepath.Clean(targetPath)

	// Convert to absolute paths
	absBase, err := filepath.Abs(cleanBase)
	if err != nil {
		return fmt.Errorf("failed to resolve base path: %w", err)
	}

	absTarget, err := filepath.Abs(cleanTarget)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	// Ensure target is within base directory
	relPath, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return fmt.Errorf("failed to determine relative path: %w", err)
	}

	// Check for path traversal attempts
	if len(relPath) > 0 && (relPath[0:1] == "." || filepath.IsAbs(relPath)) {
		return fmt.Errorf("path traversal detected: %s", targetPath)
	}

	return nil
}

// ValidatePathSafety is a public wrapper for validatePathSafety
func ValidatePathSafety(basePath, targetPath string) error {
	return validatePathSafety(basePath, targetPath)
}
