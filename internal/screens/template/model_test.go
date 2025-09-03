package template

import (
	"errors"
	"testing"

	"github.com/user/shotgun-cli/internal/models"
)

func TestNewTemplateModel(t *testing.T) {
	model := NewTemplateModel()

	if len(model.templates) != 0 {
		t.Errorf("expected empty templates slice, got %d items", len(model.templates))
	}

	if model.selected != nil {
		t.Error("expected selected template to be nil")
	}

	if !model.loading {
		t.Error("expected model to be in loading state initially")
	}

	if model.ready {
		t.Error("expected model to not be ready initially")
	}

	if !model.showDetails {
		t.Error("expected model to show details by default")
	}
}

func TestTemplateModel_SetTemplates(t *testing.T) {
	model := NewTemplateModel()

	templates := []models.Template{
		{
			ID:          "test1",
			Name:        "Test Template 1",
			Version:     "1.0.0",
			Description: "First test template",
			Author:      "Test Author",
			Tags:        []string{"test", "template"},
		},
		{
			ID:          "test2",
			Name:        "Test Template 2",
			Version:     "2.0.0",
			Description: "Second test template",
			Author:      "Another Author",
			Tags:        []string{"test", "example"},
		},
	}

	model.SetTemplates(templates)

	if len(model.templates) != 2 {
		t.Errorf("expected 2 templates, got %d", len(model.templates))
	}

	if model.loading {
		t.Error("expected model to not be loading after SetTemplates")
	}

	if !model.ready {
		t.Error("expected model to be ready after SetTemplates")
	}

	if model.selected == nil {
		t.Error("expected first template to be auto-selected")
	}

	if model.selected.ID != "test1" {
		t.Errorf("expected first template to be selected, got %s", model.selected.ID)
	}

	// Test list items were created
	if model.list.Items() == nil || len(model.list.Items()) != 2 {
		t.Error("expected list items to be created")
	}
}

func TestTemplateModel_CanAdvance(t *testing.T) {
	model := NewTemplateModel()

	// Initially should not be able to advance (no selection)
	if model.CanAdvance() {
		t.Error("expected CanAdvance to return false when no template selected")
	}

	// After setting templates, should be able to advance (auto-selects first)
	templates := []models.Template{
		{ID: "test1", Name: "Test Template", Description: "Test"},
	}
	model.SetTemplates(templates)

	if !model.CanAdvance() {
		t.Error("expected CanAdvance to return true when template is selected")
	}
}

func TestTemplateModel_UpdateSize(t *testing.T) {
	model := NewTemplateModel()

	width := 100
	height := 50

	model.UpdateSize(width, height)

	if model.width != width {
		t.Errorf("expected width %d, got %d", width, model.width)
	}

	if model.height != height {
		t.Errorf("expected height %d, got %d", height, model.height)
	}

	// Test that model dimensions are updated (list size is internal)
	if model.height != height {
		t.Errorf("expected model height %d, got %d", height, model.height)
	}

	// Test detail panel layout when wide enough and showing details
	if model.showDetails && width > 80 {
		expectedDetailWidth := width - (width * 2 / 3) - 2
		if model.viewport.Width != expectedDetailWidth {
			t.Errorf("expected viewport width %d, got %d", expectedDetailWidth, model.viewport.Width)
		}
	}
}

func TestTemplateModel_SetError(t *testing.T) {
	model := NewTemplateModel()

	testError := errors.New("Test error")

	model.SetError(testError)

	if model.err != testError {
		t.Error("expected error to be set")
	}

	if model.loading {
		t.Error("expected loading to be false after error")
	}
}

func TestTemplateModel_IsLoading(t *testing.T) {
	model := NewTemplateModel()

	// Initially should be loading
	if !model.IsLoading() {
		t.Error("expected model to be loading initially")
	}

	// After setting templates, should not be loading
	templates := []models.Template{
		{ID: "test1", Name: "Test Template", Description: "Test"},
	}
	model.SetTemplates(templates)

	if model.IsLoading() {
		t.Error("expected model to not be loading after SetTemplates")
	}

	// After error, should not be loading
	model = NewTemplateModel()
	model.SetError(nil)

	if model.IsLoading() {
		t.Error("expected model to not be loading after SetError")
	}
}

func TestTemplateItem_FilterValue(t *testing.T) {
	template := models.Template{
		ID:          "test1",
		Name:        "Test Template",
		Description: "A test template",
	}

	item := TemplateItem{Template: template}

	if item.FilterValue() != "Test Template" {
		t.Errorf("expected filter value 'Test Template', got '%s'", item.FilterValue())
	}
}

func TestTemplateItem_Title(t *testing.T) {
	tests := []struct {
		name     string
		template models.Template
		expected string
	}{
		{
			name: "template with version",
			template: models.Template{
				Name:    "Test Template",
				Version: "1.0.0",
			},
			expected: "Test Template v1.0.0",
		},
		{
			name: "template without version",
			template: models.Template{
				Name: "Test Template",
			},
			expected: "Test Template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := TemplateItem{Template: tt.template}

			if item.Title() != tt.expected {
				t.Errorf("expected title '%s', got '%s'", tt.expected, item.Title())
			}
		})
	}
}

func TestTemplateItem_Description(t *testing.T) {
	template := models.Template{
		Description: "A test template description",
	}

	item := TemplateItem{Template: template}

	if item.Description() != "A test template description" {
		t.Errorf("expected description 'A test template description', got '%s'", item.Description())
	}
}
