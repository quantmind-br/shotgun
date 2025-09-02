package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Ignorer handles .gitignore and .shotgunignore pattern matching
type Ignorer struct {
	patterns []string
	baseDir  string
}

// NewIgnorer creates a new Ignorer for the given directory
func NewIgnorer(baseDir string) (*Ignorer, error) {
	ignorer := &Ignorer{
		baseDir:  filepath.Clean(baseDir),
		patterns: make([]string, 0),
	}

	// Load .gitignore patterns first
	if err := ignorer.loadIgnoreFile(".gitignore"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .gitignore: %w", err)
	}

	// Load .shotgunignore patterns (these take priority)
	if err := ignorer.loadIgnoreFile(".shotgunignore"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .shotgunignore: %w", err)
	}

	return ignorer, nil
}

// loadIgnoreFile loads patterns from an ignore file
func (ig *Ignorer) loadIgnoreFile(filename string) error {
	ignoreFilePath := filepath.Join(ig.baseDir, filename)
	
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		ig.patterns = append(ig.patterns, line)
	}

	return scanner.Err()
}

// IsIgnored checks if a file path should be ignored based on loaded patterns
func (ig *Ignorer) IsIgnored(path string) bool {
	// Convert to relative path from base directory
	relPath, err := filepath.Rel(ig.baseDir, path)
	if err != nil {
		// If we can't make it relative, assume it's not ignored
		return false
	}

	// Normalize path separators for consistent matching
	relPath = filepath.ToSlash(relPath)

	ignored := false
	
	// First pass: check normal patterns
	for _, pattern := range ig.patterns {
		// Skip negation patterns in first pass
		if strings.HasPrefix(pattern, "!") {
			continue
		}

		if ig.matchPattern(relPath, pattern) {
			ignored = true
		}
	}

	// Second pass: check negation patterns (they override ignored status)
	for _, pattern := range ig.patterns {
		if strings.HasPrefix(pattern, "!") {
			negatePattern := strings.TrimPrefix(pattern, "!")
			if ig.matchPattern(relPath, negatePattern) {
				ignored = false // Explicitly not ignored
			}
		}
	}

	return ignored
}

// matchPattern checks if a path matches a gitignore-style pattern
func (ig *Ignorer) matchPattern(path, pattern string) bool {
	// Normalize pattern separators
	pattern = filepath.ToSlash(pattern)

	// Handle directory-only patterns (ending with /)
	if strings.HasSuffix(pattern, "/") {
		pattern = strings.TrimSuffix(pattern, "/")
		// For directory patterns, check if any parent directory matches
		dirs := strings.Split(path, "/")
		for i := range dirs {
			dirPath := strings.Join(dirs[:i+1], "/")
			if matched, _ := doublestar.Match(pattern, dirPath); matched {
				return true
			}
		}
		return false
	}

	// Handle patterns starting with / (absolute from root)
	if strings.HasPrefix(pattern, "/") {
		pattern = strings.TrimPrefix(pattern, "/")
		matched, _ := doublestar.Match(pattern, path)
		return matched
	}

	// Handle relative patterns - they can match at any level
	matched, _ := doublestar.Match(pattern, path)
	if matched {
		return true
	}

	// Also check if pattern matches any subdirectory
	matched, _ = doublestar.Match("**/"+pattern, path)
	return matched
}

// GetPatterns returns all loaded ignore patterns (useful for debugging)
func (ig *Ignorer) GetPatterns() []string {
	return append([]string(nil), ig.patterns...)
}