package main

import (
	"go/build"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestProjectStructure(t *testing.T) {
	// Required directories based on story requirements
	requiredDirs := []string{
		"cmd/shotgun",
		"internal/ui/app",
		"internal/ui/screens",
		"internal/ui/components", 
		"internal/ui/styles",
		"internal/ui/messages",
		"internal/core/scanner",
		"internal/core/template",
		"internal/core/ignore",
		"internal/core/output",
		"internal/services",
		"internal/models",
		"internal/infrastructure/filesystem",
		"internal/infrastructure/storage",
		"internal/infrastructure/platform",
		"internal/testutil/fixtures",
		"internal/testutil/mocks",
		"internal/testutil/helpers",
		"templates/embedded",
		"templates/examples",
		"configs",
		"scripts",
		"test/integration",
		"test/e2e",
		".github/workflows",
	}
	
	for _, dir := range requiredDirs {
		t.Run("directory "+dir+" exists", func(t *testing.T) {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				t.Errorf("Required directory %s does not exist", dir)
			}
		})
	}
}

func TestRequiredFiles(t *testing.T) {
	requiredFiles := map[string]string{
		"go.mod":                "Go module file",
		"cmd/shotgun/main.go":   "Main application entry point", 
		"Makefile":              "Build automation",
		".gitignore":            "Git exclusions",
		"README.md":             "Project documentation",
	}
	
	for file, description := range requiredFiles {
		t.Run(description+" exists", func(t *testing.T) {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("Required file %s does not exist", file)
			}
		})
	}
}

func TestGoModuleConfiguration(t *testing.T) {
	t.Run("go.mod has correct module name", func(t *testing.T) {
		content, err := os.ReadFile("go.mod")
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}
		
		if !strings.Contains(string(content), "module shotgun-cli-v3") {
			t.Error("go.mod should contain 'module shotgun-cli-v3'")
		}
	})
	
	t.Run("go.mod specifies Go 1.22+", func(t *testing.T) {
		content, err := os.ReadFile("go.mod")
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}
		
		// Match go version pattern
		goVersionRegex := regexp.MustCompile(`go (\d+)\.(\d+)`)
		matches := goVersionRegex.FindStringSubmatch(string(content))
		
		if len(matches) < 3 {
			t.Fatal("Could not find Go version in go.mod")
		}
		
		major := matches[1]
		minor := matches[2]
		
		if major != "1" {
			t.Errorf("Expected Go 1.x, got %s.x", major)
		}
		
		// Convert minor version to int for comparison
		minorInt := 0
		for _, r := range minor {
			if r >= '0' && r <= '9' {
				minorInt = minorInt*10 + int(r-'0')
			}
		}
		
		if minorInt < 22 {
			t.Errorf("Expected Go 1.22+, got 1.%d", minorInt)
		}
	})
}

func TestCoreDependencies(t *testing.T) {
	// Required core dependencies from the story requirements
	requiredDeps := []string{
		"github.com/charmbracelet/bubbletea/v2",
		"github.com/charmbracelet/bubbles",
		"github.com/charmbracelet/lipgloss", 
		"github.com/bmatcuk/doublestar/v4",
		"github.com/h2non/filetype",
		"github.com/spf13/viper",
		"github.com/spf13/cobra",
		"github.com/BurntSushi/toml",
	}
	
	content, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	
	modContent := string(content)
	
	for _, dep := range requiredDeps {
		t.Run("dependency "+dep+" is present", func(t *testing.T) {
			if !strings.Contains(modContent, dep) {
				t.Errorf("Required dependency %s not found in go.mod", dep)
			}
		})
	}
}

func TestMakefileTargets(t *testing.T) {
	content, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}
	
	makefileContent := string(content)
	requiredTargets := []string{
		"build:",
		"test:", 
		"clean:",
		"lint:",
		"cross-compile:",
	}
	
	for _, target := range requiredTargets {
		t.Run("Makefile has "+target+" target", func(t *testing.T) {
			if !strings.Contains(makefileContent, target) {
				t.Errorf("Required Makefile target %s not found", target)
			}
		})
	}
}

func TestGitignoreConfiguration(t *testing.T) {
	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}
	
	gitignoreContent := string(content)
	requiredPatterns := []string{
		"*.exe",
		"*.test", 
		"*.out",
		"vendor/",
		"dist/",
		".DS_Store",
		"Thumbs.db",
		".vscode/",
		".idea/",
	}
	
	for _, pattern := range requiredPatterns {
		t.Run(".gitignore contains "+pattern, func(t *testing.T) {
			if !strings.Contains(gitignoreContent, pattern) {
				t.Errorf("Required .gitignore pattern %s not found", pattern)
			}
		})
	}
}

func TestReadmeContent(t *testing.T) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}
	
	readmeContent := string(content)
	requiredSections := []string{
		"# Shotgun CLI v3",
		"## Overview",
		"## Installation", 
		"## Usage",
		"## Development Setup",
		"## Technology Stack",
	}
	
	for _, section := range requiredSections {
		t.Run("README contains "+section, func(t *testing.T) {
			if !strings.Contains(readmeContent, section) {
				t.Errorf("Required README section %s not found", section)
			}
		})
	}
}

func TestGoEnvironment(t *testing.T) {
	t.Run("Go version meets requirements", func(t *testing.T) {
		goVersion := build.Default.ReleaseTags
		if len(goVersion) == 0 {
			t.Skip("Cannot determine Go version")
		}
		
		// Check if we have a recent enough version
		hasGo122 := false
		for _, tag := range goVersion {
			if strings.HasPrefix(tag, "go1.22") || 
			   strings.HasPrefix(tag, "go1.23") ||
			   strings.HasPrefix(tag, "go1.24") {
				hasGo122 = true
				break
			}
		}
		
		if !hasGo122 {
			t.Error("Go 1.22+ is required")
		}
	})
}