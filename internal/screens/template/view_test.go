package template

import (
	"errors"
	"strings"
	"testing"

	"github.com/diogopedro/shotgun/internal/models"
)

func TestTemplateModel_View_LoadingState(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(80, 24)
	// Model starts in loading state

	view := model.View()

	if !strings.Contains(view, "Loading templates...") {
		t.Error("expected view to contain loading message")
	}

	if !strings.Contains(view, "Template Selection") {
		t.Error("expected view to contain screen title")
	}
}

func TestTemplateModel_View_ErrorState(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(80, 24)

	testError := errors.New("Failed to load templates")

	model.SetError(testError)

	view := model.View()

	if !strings.Contains(view, "Error loading templates") {
		t.Error("expected view to contain error message")
	}

	if !strings.Contains(view, "Failed to load templates") {
		t.Error("expected view to contain specific error details")
	}

	if !strings.Contains(view, "Ctrl+Left to go back") {
		t.Error("expected view to contain help text for error state")
	}
}

func TestTemplateModel_View_EmptyState(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(80, 24)

	// Set empty templates (not loading, no error)
	model.SetTemplates([]models.Template{})

	view := model.View()

	if !strings.Contains(view, "No templates found") {
		t.Error("expected view to contain empty state message")
	}

	if !strings.Contains(view, "Template Selection") {
		t.Error("expected view to contain screen title")
	}

	if !strings.Contains(view, "Ctrl+R to refresh") {
		t.Error("expected view to contain refresh help text")
	}
}

func TestTemplateModel_View_NormalState(t *testing.T) {
	model := setupModelWithTemplates()
	model.UpdateSize(80, 24)

	view := model.View()

	if !strings.Contains(view, "Template Selection") {
		t.Error("expected view to contain screen title")
	}

	// Should show help text for normal operation
	expectedHelpItems := []string{
		"navigate", // Navigation help
		"select",   // Selection help (when template is selected)
		"back",     // Back navigation
		"refresh",  // Refresh option
		"quit",     // Quit option
	}

	for _, helpItem := range expectedHelpItems {
		if !strings.Contains(view, helpItem) {
			t.Errorf("expected view to contain help text '%s'", helpItem)
		}
	}
}

func TestTemplateModel_View_ResponsiveLayout(t *testing.T) {
	model := setupModelWithTemplates()

	// Test narrow layout (single column)
	model.UpdateSize(60, 24)
	model.showDetails = true

	viewNarrow := model.View()

	// Test wide layout (side-by-side with details)
	model.UpdateSize(120, 24)
	model.showDetails = true

	viewWide := model.View()

	// Wide view should be different from narrow view when details are shown
	if viewWide == viewNarrow {
		t.Error("expected different layout for wide vs narrow screen")
	}

	// Test wide layout without details
	model.showDetails = false

	viewWideNoDetails := model.View()

	if viewWideNoDetails == viewWide {
		t.Error("expected different layout when details are hidden")
	}
}

func TestTemplateModel_View_ZeroSize(t *testing.T) {
	model := setupModelWithTemplates()

	// Test with zero or invalid size
	model.UpdateSize(0, 0)

	view := model.View()

	if view != "Loading..." {
		t.Errorf("expected 'Loading...' for zero size, got '%s'", view)
	}
}

func TestTemplateModel_View_HelpText(t *testing.T) {
	model := setupModelWithTemplates()

	// Test help text changes based on state
	tests := []struct {
		name        string
		width       int
		hasSelected bool
		showDetails bool
		expected    []string
		notExpected []string
	}{
		{
			name:        "narrow screen without details",
			width:       60,
			hasSelected: true,
			showDetails: false,
			expected:    []string{"navigate", "select", "back", "refresh", "quit"},
			notExpected: []string{"toggle details"},
		},
		{
			name:        "wide screen with details",
			width:       120,
			hasSelected: true,
			showDetails: true,
			expected:    []string{"navigate", "select", "back", "refresh", "quit", "toggle details"},
			notExpected: []string{},
		},
		{
			name:        "no template selected",
			width:       80,
			hasSelected: false,
			showDetails: true,
			expected:    []string{"navigate", "back", "refresh", "quit"},
			notExpected: []string{"select"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.hasSelected {
				// Clear selection
				model.selected = nil
			} else {
				// Ensure selection
				if len(model.templates) > 0 {
					model.selected = &model.templates[0]
				}
			}

			model.UpdateSize(tt.width, 24)
			model.showDetails = tt.showDetails

			view := model.View()

			for _, expected := range tt.expected {
				if !strings.Contains(view, expected) {
					t.Errorf("expected view to contain '%s'", expected)
				}
			}

			for _, notExpected := range tt.notExpected {
				if strings.Contains(view, notExpected) {
					t.Errorf("expected view to NOT contain '%s'", notExpected)
				}
			}
		})
	}
}

func TestRenderDetailPanel(t *testing.T) {
	tests := []struct {
		name     string
		template *models.Template
		width    int
		height   int
		expected []string
	}{
		{
			name:     "nil template",
			template: nil,
			width:    40,
			height:   20,
			expected: []string{"No template selected"},
		},
		{
			name: "complete template",
			template: &models.Template{
				ID:          "test1",
				Name:        "Test Template",
				Version:     "1.0.0",
				Description: "A comprehensive test template with all fields",
				Author:      "Test Author",
				Tags:        []string{"test", "template", "example"},
				Variables: map[string]models.Variable{
					"name": {
						Type:        "string",
						Placeholder: "Project name",
					},
					"version": {
						Type:        "string",
						Placeholder: "Project version",
					},
				},
			},
			width:  60,
			height: 25,
			expected: []string{
				"Test Template",
				"Version: 1.0.0",
				"Description:",
				"A comprehensive test template",
				"Author: Test Author",
				"Tags:",
				"test", "template", "example",
				"Variables:",
				"name (string)",
				"version (string)",
			},
		},
		{
			name: "minimal template",
			template: &models.Template{
				ID:   "minimal",
				Name: "Minimal Template",
			},
			width:    40,
			height:   15,
			expected: []string{"Minimal Template"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderDetailPanel(tt.template, tt.width, tt.height)

			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("expected detail panel to contain '%s', got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected []string
	}{
		{
			name:     "short text",
			text:     "Hello world",
			width:    20,
			expected: []string{"Hello world"},
		},
		{
			name:  "text that needs wrapping",
			text:  "This is a long text that needs to be wrapped across multiple lines",
			width: 20,
			expected: []string{
				"This is a long text",
				"that needs to be",
				"wrapped across",
				"multiple lines",
			},
		},
		{
			name:  "single long word",
			text:  "supercalifragilisticexpialidocious",
			width: 10,
			expected: []string{
				"supercalif",
				"ragilistic",
				"expialidoc",
				"ious",
			},
		},
		{
			name:     "empty text",
			text:     "",
			width:    10,
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapText(tt.text, tt.width)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d lines, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if i < len(result) && result[i] != expected {
					t.Errorf("line %d: expected '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}
