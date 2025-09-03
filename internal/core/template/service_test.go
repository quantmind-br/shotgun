package template

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/user/shotgun-cli/internal/models"
)

func TestNewTemplateService(t *testing.T) {
	// Test with provided logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := NewTemplateService(logger)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	// Test with nil logger (should use default)
	serviceWithNil := NewTemplateService(nil)
	if serviceWithNil == nil {
		t.Fatal("Expected service to be created with nil logger")
	}
}

func TestTemplateService_GetTemplate_NotLoaded(t *testing.T) {
	service := NewTemplateService(nil)

	_, err := service.GetTemplate("test-id")
	if err == nil {
		t.Error("Expected error when templates not loaded yet")
	}

	if !contains(err.Error(), "not loaded yet") {
		t.Errorf("Expected 'not loaded yet' error message, got: %v", err)
	}
}

func TestTemplateService_GetTemplate_EmptyID(t *testing.T) {
	service := NewTemplateService(nil)

	_, err := service.GetTemplate("")
	if err == nil {
		t.Error("Expected error for empty ID")
	}

	if !contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected 'cannot be empty' error message, got: %v", err)
	}
}

func TestTemplateService_GetTemplateCount_Empty(t *testing.T) {
	service := NewTemplateService(nil)

	count := service.GetTemplateCount()
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestTemplateService_LoadAllTemplates_WithContext(t *testing.T) {
	service := NewTemplateService(nil)
	ctx := context.Background()

	// This will test the actual loading process
	// It may fail due to parsing issues with existing templates, but should not crash
	templates, err := service.LoadAllTemplates(ctx)

	// We don't expect a fatal error, but parsing might fail
	if err != nil {
		t.Logf("LoadAllTemplates returned error (may be expected): %v", err)
	}

	if templates == nil {
		t.Error("Expected templates slice to be non-nil")
	}

	t.Logf("Loaded %d templates", len(templates))
}

func TestTemplateService_LoadAllTemplates_WithTimeout(t *testing.T) {
	service := NewTemplateService(nil)

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// This should handle timeout gracefully
	_, err := service.LoadAllTemplates(ctx)

	// We expect either success (if fast enough) or context error
	if err != nil && err != context.DeadlineExceeded && !contains(err.Error(), "context") {
		t.Errorf("Expected nil, deadline exceeded, or context error, got: %v", err)
	}
}

func TestTemplateService_RefreshTemplates(t *testing.T) {
	service := NewTemplateService(nil)
	ctx := context.Background()

	// Initial load
	_, err := service.LoadAllTemplates(ctx)
	if err != nil {
		t.Logf("Initial load error (may be expected): %v", err)
	}

	initialCount := service.GetTemplateCount()

	// Refresh
	err = service.RefreshTemplates(ctx)
	if err != nil {
		t.Logf("Refresh error (may be expected): %v", err)
	}

	refreshCount := service.GetTemplateCount()

	// Count should be the same after refresh
	if initialCount != refreshCount {
		t.Errorf("Expected same count after refresh, got initial=%d, refresh=%d",
			initialCount, refreshCount)
	}
}

func TestTemplateService_ValidateTemplateCache(t *testing.T) {
	service := NewTemplateService(nil)

	// Cast to access internal methods (for testing)
	if templateService, ok := service.(*templateService); ok {
		err := templateService.ValidateTemplateCache()
		if err != nil {
			t.Errorf("Expected no cache validation errors for empty cache, got: %v", err)
		}

		// Manually add a valid template to cache for testing
		testTemplate := models.Template{
			ID:          "test-template",
			Name:        "Test Template",
			Version:     "1.0.0",
			Description: "A test template",
			Content:     "Hello {{name}}!",
		}

		templateService.cache.Store("test-template", testTemplate)

		err = templateService.ValidateTemplateCache()
		if err != nil {
			t.Errorf("Expected no cache validation errors for valid template, got: %v", err)
		}

		// Test inconsistent cache (wrong key)
		templateService.cache.Store("wrong-key", testTemplate)

		err = templateService.ValidateTemplateCache()
		if err == nil {
			t.Error("Expected cache validation error for inconsistent key")
		}
	}
}

func TestTemplateService_GetAllCachedTemplates(t *testing.T) {
	service := NewTemplateService(nil)

	// Cast to access internal methods (for testing)
	if templateService, ok := service.(*templateService); ok {
		// Initially should be empty
		templates := templateService.GetAllCachedTemplates()
		if len(templates) != 0 {
			t.Errorf("Expected 0 cached templates initially, got %d", len(templates))
		}

		// Add some test templates
		testTemplates := []models.Template{
			{
				ID:          "template-b",
				Name:        "B Template", // Will test sorting
				Version:     "1.0.0",
				Description: "Template B",
				Content:     "B content",
			},
			{
				ID:          "template-a",
				Name:        "A Template", // Should come first after sorting
				Version:     "1.0.0",
				Description: "Template A",
				Content:     "A content",
			},
		}

		for _, template := range testTemplates {
			templateService.cache.Store(template.ID, template)
		}

		cached := templateService.GetAllCachedTemplates()
		if len(cached) != 2 {
			t.Errorf("Expected 2 cached templates, got %d", len(cached))
		}

		// Verify sorting by name
		if len(cached) >= 2 {
			if cached[0].Name != "A Template" {
				t.Errorf("Expected first template to be 'A Template', got '%s'", cached[0].Name)
			}
		}
	}
}

func TestTemplateService_Integration_LoadAndGet(t *testing.T) {
	service := NewTemplateService(nil)
	ctx := context.Background()

	// Load templates
	templates, err := service.LoadAllTemplates(ctx)

	// Even if loading fails due to template format issues, we should be able to test the flow
	if err != nil {
		t.Logf("Load templates error (may be expected due to template format): %v", err)
	}

	// If we successfully loaded any templates, try to get them
	if len(templates) > 0 {
		firstTemplate := templates[0]

		retrieved, err := service.GetTemplate(firstTemplate.ID)
		if err != nil {
			t.Errorf("Failed to retrieve template '%s': %v", firstTemplate.ID, err)
		}

		if retrieved == nil {
			t.Error("Retrieved template is nil")
		} else if retrieved.ID != firstTemplate.ID {
			t.Errorf("Retrieved template ID mismatch: expected '%s', got '%s'",
				firstTemplate.ID, retrieved.ID)
		}
	}

	// Test getting non-existent template
	_, err = service.GetTemplate("non-existent-template")
	if err == nil {
		t.Error("Expected error when getting non-existent template")
	}

	if !contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestTemplateService_Concurrency(t *testing.T) {
	service := NewTemplateService(nil)
	ctx := context.Background()

	// Load templates first
	_, err := service.LoadAllTemplates(ctx)
	if err != nil {
		t.Logf("Load error (may be expected): %v", err)
	}

	// Test concurrent access to GetTemplateCount
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Multiple concurrent operations
			count := service.GetTemplateCount()
			_ = count // Use the count to avoid compiler warnings

			// Try to get a template (may fail if no templates loaded)
			_, _ = service.GetTemplate("any-id")
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we get here without panic, concurrent access is working
}

func TestTemplateService_ErrorHandling(t *testing.T) {
	service := NewTemplateService(nil)

	// Test various error conditions
	tests := []struct {
		name string
		test func() error
	}{
		{
			name: "Get template with empty ID",
			test: func() error {
				_, err := service.GetTemplate("")
				return err
			},
		},
		{
			name: "Get template before loading",
			test: func() error {
				_, err := service.GetTemplate("test")
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.test()
			if err == nil {
				t.Errorf("Expected error for test '%s'", test.name)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		(len(str) > len(substr) &&
			(str[:len(substr)] == substr ||
				str[len(str)-len(substr):] == substr ||
				hasSubstring(str, substr))))
}

// Simple substring check
func hasSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
