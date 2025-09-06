package templates

import (
	"strings"
	"testing"
)

func TestShotgunignoreTemplate(t *testing.T) {
	if ShotgunignoreTemplate == "" {
		t.Error("ShotgunignoreTemplate is empty")
	}

	// Check that template contains expected header
	if !strings.Contains(ShotgunignoreTemplate, "# Shotgun ignore patterns") {
		t.Error("template missing expected header")
	}

	// Check that template contains expected sections
	expectedSections := []string{
		"# Build artifacts",
		"# Dependencies",
		"# IDE and editor files",
		"# OS generated files",
		"# Logs and temporary files",
		"# Runtime and environment files",
		"# Version control",
		"# Package managers",
		"# Go specific",
		"# Language specific artifacts",
	}

	for _, section := range expectedSections {
		if !strings.Contains(ShotgunignoreTemplate, section) {
			t.Errorf("template missing expected section: %s", section)
		}
	}

	// Check that template contains key patterns
	expectedPatterns := []string{
		"build/",
		"dist/",
		"node_modules/",
		"vendor/",
		".vscode/",
		".idea/",
		".DS_Store",
		"Thumbs.db",
		"*.log",
		".env",
		".git/",
		"go.sum",
		"*.test",
		"__pycache__/",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(ShotgunignoreTemplate, pattern) {
			t.Errorf("template missing expected pattern: %s", pattern)
		}
	}
}

func TestShotgunignoreCategories(t *testing.T) {
	if len(ShotgunignoreCategories) == 0 {
		t.Error("ShotgunignoreCategories is empty")
	}

	// Check that we have expected number of categories (should be 10 based on template)
	expectedMinCategories := 10
	if len(ShotgunignoreCategories) < expectedMinCategories {
		t.Errorf("expected at least %d categories, got %d", expectedMinCategories, len(ShotgunignoreCategories))
	}

	// Verify each category has required fields
	for i, category := range ShotgunignoreCategories {
		if category.Name == "" {
			t.Errorf("category %d has empty Name", i)
		}

		if category.Description == "" {
			t.Errorf("category %d (%s) has empty Description", i, category.Name)
		}

		if len(category.Patterns) == 0 {
			t.Errorf("category %d (%s) has no patterns", i, category.Name)
		}

		// Verify patterns are non-empty strings
		for j, pattern := range category.Patterns {
			if pattern == "" {
				t.Errorf("category %d (%s) has empty pattern at index %d", i, category.Name, j)
			}
		}
	}
}

func TestCategoryNames(t *testing.T) {
	expectedNames := []string{
		"Build artifacts",
		"Dependencies",
		"IDE and editor files",
		"OS generated files",
		"Logs and temporary files",
		"Runtime and environment files",
		"Version control",
		"Package managers",
		"Go specific",
		"Language specific artifacts",
	}

	actualNames := make([]string, len(ShotgunignoreCategories))
	for i, category := range ShotgunignoreCategories {
		actualNames[i] = category.Name
	}

	for _, expected := range expectedNames {
		found := false
		for _, actual := range actualNames {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected category '%s' not found in categories", expected)
		}
	}
}

func TestSpecificCategoryPatterns(t *testing.T) {
	// Test Build artifacts category
	buildCategory := findCategory("Build artifacts")
	if buildCategory == nil {
		t.Fatal("Build artifacts category not found")
	}

	expectedBuildPatterns := []string{"build/", "dist/", "target/", "*.exe", "*.dll"}
	for _, expected := range expectedBuildPatterns {
		if !containsPattern(buildCategory.Patterns, expected) {
			t.Errorf("Build artifacts category missing pattern: %s", expected)
		}
	}

	// Test Go specific category
	goCategory := findCategory("Go specific")
	if goCategory == nil {
		t.Fatal("Go specific category not found")
	}

	expectedGoPatterns := []string{"go.sum", "*.test", "*.out", "coverage.txt"}
	for _, expected := range expectedGoPatterns {
		if !containsPattern(goCategory.Patterns, expected) {
			t.Errorf("Go specific category missing pattern: %s", expected)
		}
	}
}

// Helper functions

func findCategory(name string) *IgnoreCategory {
	for _, category := range ShotgunignoreCategories {
		if category.Name == name {
			return &category
		}
	}
	return nil
}

func containsPattern(patterns []string, pattern string) bool {
	for _, p := range patterns {
		if p == pattern {
			return true
		}
	}
	return false
}
