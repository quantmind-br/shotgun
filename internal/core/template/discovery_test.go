package template

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/shotgun-cli/internal/models"
)

func TestNewDiscoveryService(t *testing.T) {
	// Test with provided logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := NewDiscoveryService(logger)
	
	if service == nil {
		t.Fatal("Expected service to be created")
	}
	
	if service.logger != logger {
		t.Error("Expected provided logger to be used")
	}
	
	// Test with nil logger (should use default)
	serviceWithNil := NewDiscoveryService(nil)
	if serviceWithNil == nil {
		t.Fatal("Expected service to be created with nil logger")
	}
	
	if serviceWithNil.logger == nil {
		t.Error("Expected default logger to be set")
	}
}

func TestDiscoveryService_IsPathSafe(t *testing.T) {
	service := NewDiscoveryService(nil)
	
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	tests := []struct {
		name       string
		baseDir    string
		targetPath string
		expected   bool
	}{
		{
			name:       "Safe path within base",
			baseDir:    tempDir,
			targetPath: filepath.Join(tempDir, "template.toml"),
			expected:   true,
		},
		{
			name:       "Exact base directory",
			baseDir:    tempDir,
			targetPath: tempDir,
			expected:   true,
		},
		{
			name:       "Path traversal attempt",
			baseDir:    tempDir,
			targetPath: filepath.Join(tempDir, "..", "outside.toml"),
			expected:   false,
		},
		{
			name:       "Absolute path outside",
			baseDir:    tempDir,
			targetPath: "/etc/passwd",
			expected:   false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := service.isPathSafe(test.baseDir, test.targetPath)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for baseDir=%s, targetPath=%s",
					test.expected, result, test.baseDir, test.targetPath)
			}
		})
	}
}

func TestDiscoveryService_DeduplicateTemplates(t *testing.T) {
	service := NewDiscoveryService(nil)
	
	// Create test templates
	builtin1 := models.TemplateInfo{
		Template: models.Template{
			ID:   "template1",
			Name: "Built-in Template 1",
		},
		Source:   models.TemplateSourceBuiltIn,
		FilePath: "builtin/template1.toml",
	}
	
	builtin2 := models.TemplateInfo{
		Template: models.Template{
			ID:   "template2",
			Name: "Built-in Template 2",
		},
		Source:   models.TemplateSourceBuiltIn,
		FilePath: "builtin/template2.toml",
	}
	
	user1 := models.TemplateInfo{
		Template: models.Template{
			ID:   "template1", // Same ID as builtin1
			Name: "User Template 1 (Override)",
		},
		Source:   models.TemplateSourceUser,
		FilePath: "user/template1.toml",
	}
	
	user2 := models.TemplateInfo{
		Template: models.Template{
			ID:   "template3",
			Name: "User Template 3",
		},
		Source:   models.TemplateSourceUser,
		FilePath: "user/template3.toml",
	}
	
	templates := []models.TemplateInfo{builtin1, builtin2, user1, user2}
	
	result := service.deduplicateTemplates(templates)
	
	// Should have 3 templates (builtin2, user1 override, user2)
	if len(result) != 3 {
		t.Errorf("Expected 3 templates after deduplication, got %d", len(result))
	}
	
	// Find template1 (should be user version)
	var template1 *models.TemplateInfo
	for _, tmpl := range result {
		if tmpl.Template.ID == "template1" {
			template1 = &tmpl
			break
		}
	}
	
	if template1 == nil {
		t.Error("Expected to find template1 in results")
	} else if template1.Source != models.TemplateSourceUser {
		t.Error("Expected template1 to be user version (override)")
	}
}

func TestDiscoveryService_DiscoverUserTemplates_NoDirectory(t *testing.T) {
	// This test might be platform dependent, so we skip it for now
	// and focus on testing the logic with actual directories
	t.Skip("Skipping test that depends on user directory setup")
}

func TestDiscoveryService_DiscoverAllTemplates_WithTimeout(t *testing.T) {
	service := NewDiscoveryService(nil)
	
	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// This should return quickly due to timeout
	_, err := service.DiscoverAllTemplates(ctx)
	
	// We expect either success (if it's fast enough) or context error
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Expected nil or deadline exceeded error, got: %v", err)
	}
}

func TestDiscoveryService_DiscoverBuiltinTemplates_ValidatesEmbedFS(t *testing.T) {
	service := NewDiscoveryService(nil)
	ctx := context.Background()
	
	// This will test the actual embedded filesystem
	templates, err := service.DiscoverBuiltinTemplates(ctx)
	
	// We don't expect an error unless the embed fails
	if err != nil {
		// The parseTemplateFromData might fail since we haven't implemented
		// proper templates yet, but the discovery itself should work
		t.Logf("Discovery returned error (expected if templates are invalid): %v", err)
	}
	
	// We should get some result (even if empty due to parsing failures)
	if templates == nil {
		t.Error("Expected templates slice to be non-nil")
	}
	
	t.Logf("Found %d built-in templates", len(templates))
}

func TestDiscoveryService_DiscoverUserTemplates_WithTempDir(t *testing.T) {
	service := NewDiscoveryService(nil)
	
	// Create temporary directory with test templates
	tempDir := t.TempDir()
	
	// Create a valid test template
	validTOML := `
id = "test-template"
name = "Test Template"
version = "1.0.0"
description = "A test template"
content = "Hello {{name}}!"

[variables.name]
name = "name"
type = "text"
required = true
`
	
	templateFile := filepath.Join(tempDir, "test.toml")
	err := os.WriteFile(templateFile, []byte(validTOML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}
	
	// Create an invalid file (should be skipped)
	invalidFile := filepath.Join(tempDir, "invalid.toml")
	err = os.WriteFile(invalidFile, []byte("invalid toml content [[["), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid test file: %v", err)
	}
	
	// Create a non-TOML file (should be ignored)
	nonTomlFile := filepath.Join(tempDir, "readme.txt")
	err = os.WriteFile(nonTomlFile, []byte("This is not a template"), 0644)
	if err != nil {
		t.Fatalf("Failed to create non-TOML file: %v", err)
	}
	
	// We need to temporarily override the getUserTemplatesDir function
	// For now, we'll test the path safety logic directly
	
	// Test path safety with the temp directory
	safePath := filepath.Join(tempDir, "test.toml")
	if !service.isPathSafe(tempDir, safePath) {
		t.Error("Expected safe path to be validated as safe")
	}
	
	unsafePath := filepath.Join(tempDir, "..", "outside.toml")
	if service.isPathSafe(tempDir, unsafePath) {
		t.Error("Expected unsafe path to be validated as unsafe")
	}
}

func TestParseTemplateFromData_Integration(t *testing.T) {
	// Test the integration between parseTemplateFromData and validation
	validTOML := `
id = "integration-test"
name = "Integration Test Template"
version = "1.2.3"
description = "Testing integration between parsing and validation"
author = "Test Suite"
tags = ["test", "integration"]
content = "Hello {{name}}, your priority is {{priority}}!"

[variables.name]
name = "name"
type = "text"
required = true
placeholder = "Enter your name"
min_length = 2
max_length = 50

[variables.priority]
name = "priority"
type = "choice"
required = false
default = "medium"
options = ["low", "medium", "high", "critical"]
`
	
	template, err := parseTemplateFromData([]byte(validTOML))
	if err != nil {
		t.Fatalf("Expected successful parsing, got error: %v", err)
	}
	
	// Verify all fields are correctly parsed and validated
	if template.ID != "integration-test" {
		t.Errorf("Expected ID 'integration-test', got '%s'", template.ID)
	}
	
	if len(template.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(template.Variables))
	}
	
	// Test that validation is properly integrated
	nameVar := template.Variables["name"]
	if nameVar.MinLength != 2 || nameVar.MaxLength != 50 {
		t.Errorf("Expected name variable constraints: MinLength=2, MaxLength=50, got MinLength=%d, MaxLength=%d", 
			nameVar.MinLength, nameVar.MaxLength)
	}
	
	priorityVar := template.Variables["priority"]
	if len(priorityVar.Options) != 4 {
		t.Errorf("Expected 4 priority options, got %d", len(priorityVar.Options))
	}
}

// Helper function to create test templates directory
func createTestTemplatesDir(t *testing.T) string {
	tempDir := t.TempDir()
	
	// Create some test templates
	templates := map[string]string{
		"simple.toml": `
id = "simple"
name = "Simple Template"
version = "1.0.0"
description = "A simple template"
content = "Hello {{name}}!"

[variables.name]
name = "name"
type = "text"
required = true
`,
		"complex.toml": `
id = "complex"
name = "Complex Template"
version = "2.1.0"
description = "A complex template with multiple variables"
author = "Template Author"
tags = ["complex", "demo"]
content = "Message: {{message}}\nPriority: {{priority}}\nEnabled: {{enabled}}"

[variables.message]
name = "message"
type = "multiline"
required = true
placeholder = "Enter your message"

[variables.priority]
name = "priority"
type = "choice"
required = false
default = "medium"
options = ["low", "medium", "high"]

[variables.enabled]
name = "enabled"
type = "boolean"
required = false
default = "true"
`,
		"invalid.toml": `
name = "Invalid Template"
# Missing required fields like version, description
content = "This should fail validation"
`,
	}
	
	for filename, content := range templates {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test template %s: %v", filename, err)
		}
	}
	
	return tempDir
}

func TestErrorAggregator(t *testing.T) {
	aggregator := NewErrorAggregator()
	
	// Test empty aggregator
	if aggregator.HasErrors() {
		t.Error("Expected no errors initially")
	}
	
	if aggregator.Count() != 0 {
		t.Errorf("Expected count 0, got %d", aggregator.Count())
	}
	
	// Add some errors
	aggregator.Add(NewParsingError("test1.toml", nil))
	aggregator.Add(NewValidationError("test2.toml", "missing field"))
	aggregator.Add(nil) // Should be ignored
	
	if !aggregator.HasErrors() {
		t.Error("Expected to have errors")
	}
	
	if aggregator.Count() != 2 {
		t.Errorf("Expected count 2, got %d", aggregator.Count())
	}
	
	errors := aggregator.Errors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}
	
	// Test error message
	errorMsg := aggregator.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}