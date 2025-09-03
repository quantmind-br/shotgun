package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewIgnorer(t *testing.T) {
	tempDir, cleanup := createTestDirWithIgnoreFiles(t)
	defer cleanup()

	ignorer, err := NewIgnorer(tempDir)
	if err != nil {
		t.Fatalf("failed to create ignorer: %v", err)
	}

	if ignorer == nil {
		t.Error("ignorer is nil")
	}

	if ignorer.baseDir != tempDir {
		t.Errorf("expected baseDir = %s, got %s", tempDir, ignorer.baseDir)
	}

	// Should have loaded patterns from both .gitignore and .shotgunignore
	patterns := ignorer.GetPatterns()
	if len(patterns) == 0 {
		t.Error("expected some patterns to be loaded")
	}
}

func TestIgnorer_IsIgnored(t *testing.T) {
	tempDir, cleanup := createTestDirWithIgnoreFiles(t)
	defer cleanup()

	ignorer, err := NewIgnorer(tempDir)
	if err != nil {
		t.Fatalf("failed to create ignorer: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "ignored file from gitignore",
			path:     filepath.Join(tempDir, "build", "output.bin"),
			expected: true,
		},
		{
			name:     "ignored directory from gitignore",
			path:     filepath.Join(tempDir, "node_modules", "package"),
			expected: true,
		},
		{
			name:     "ignored file from shotgunignore",
			path:     filepath.Join(tempDir, "temp.log"),
			expected: true,
		},
		{
			name:     "not ignored file",
			path:     filepath.Join(tempDir, "src", "main.go"),
			expected: false,
		},
		{
			name:     "negated pattern",
			path:     filepath.Join(tempDir, "important.log"),
			expected: false, // Should not be ignored due to ! pattern
		},
		{
			name:     "absolute pattern",
			path:     filepath.Join(tempDir, "root-only.txt"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ignorer.IsIgnored(tt.path)
			if result != tt.expected {
				t.Errorf("IsIgnored(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIgnorer_MatchPattern(t *testing.T) {
	ignorer := &Ignorer{baseDir: "/test"}

	tests := []struct {
		name     string
		path     string
		pattern  string
		expected bool
	}{
		{
			name:     "simple wildcard",
			path:     "file.txt",
			pattern:  "*.txt",
			expected: true,
		},
		{
			name:     "double star pattern",
			path:     "dir/subdir/file.js",
			pattern:  "**/*.js",
			expected: true,
		},
		{
			name:     "directory pattern",
			path:     "node_modules/package",
			pattern:  "node_modules/",
			expected: true,
		},
		{
			name:     "absolute pattern from root",
			path:     "root-file.txt",
			pattern:  "/root-file.txt",
			expected: true,
		},
		{
			name:     "absolute pattern not matching subdirectory",
			path:     "sub/root-file.txt",
			pattern:  "/root-file.txt",
			expected: false,
		},
		{
			name:     "no match",
			path:     "file.go",
			pattern:  "*.js",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ignorer.matchPattern(tt.path, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchPattern(%s, %s) = %v, expected %v", tt.path, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestIgnorer_LoadIgnoreFile_NonExistent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignorer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ignorer, err := NewIgnorer(tempDir)
	if err != nil {
		t.Fatalf("failed to create ignorer for empty dir: %v", err)
	}

	// Should not error even if ignore files don't exist
	patterns := ignorer.GetPatterns()
	if len(patterns) != 0 {
		t.Errorf("expected no patterns, got %d", len(patterns))
	}
}

func TestIgnorer_CommentAndEmptyLines(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignorer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .gitignore with comments and empty lines
	gitignoreContent := `
# This is a comment
*.tmp

# Another comment

*.log
	# Indented comment
build/
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("failed to write .gitignore: %v", err)
	}

	ignorer, err := NewIgnorer(tempDir)
	if err != nil {
		t.Fatalf("failed to create ignorer: %v", err)
	}

	patterns := ignorer.GetPatterns()

	// Should only have 3 patterns (*.tmp, *.log, build/)
	expectedPatterns := 3
	if len(patterns) != expectedPatterns {
		t.Errorf("expected %d patterns, got %d: %v", expectedPatterns, len(patterns), patterns)
	}

	// Verify no comments are included
	for _, pattern := range patterns {
		if pattern[0] == '#' {
			t.Errorf("comment found in patterns: %s", pattern)
		}
	}
}

func TestIgnorer_NegationPatterns(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignorer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .gitignore with negation patterns
	gitignoreContent := `*.log
!important.log
!critical/*.log`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("failed to write .gitignore: %v", err)
	}

	ignorer, err := NewIgnorer(tempDir)
	if err != nil {
		t.Fatalf("failed to create ignorer: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "normal log file ignored",
			path:     filepath.Join(tempDir, "app.log"),
			expected: true,
		},
		{
			name:     "important.log not ignored due to negation",
			path:     filepath.Join(tempDir, "important.log"),
			expected: false,
		},
		{
			name:     "critical log not ignored due to negation",
			path:     filepath.Join(tempDir, "critical", "system.log"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ignorer.IsIgnored(tt.path)
			if result != tt.expected {
				t.Errorf("IsIgnored(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

// createTestDirWithIgnoreFiles creates a temp directory with .gitignore and .shotgunignore files
func createTestDirWithIgnoreFiles(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "ignorer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create .gitignore
	gitignoreContent := `# Build directories
build/
dist/
node_modules/

# Log files
*.log
!important.log

# Binary files
*.bin
*.exe

# Absolute path pattern
/root-only.txt`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("failed to write .gitignore: %v", err)
	}

	// Create .shotgunignore (takes priority)
	shotgunignoreContent := `# Temporary files
*.tmp
temp.log

# Cache
.cache/`
	shotgunignorePath := filepath.Join(tempDir, ".shotgunignore")
	if err := os.WriteFile(shotgunignorePath, []byte(shotgunignoreContent), 0644); err != nil {
		t.Fatalf("failed to write .shotgunignore: %v", err)
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}
