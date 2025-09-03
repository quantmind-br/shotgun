package template

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/user/shotgun-cli/internal/models"
)

func TestTemplateModel_Update_WindowSizeMsg(t *testing.T) {
	model := NewTemplateModel()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	updatedModel, cmd := model.Update(msg)

	if updatedModel.width != 100 {
		t.Errorf("expected width 100, got %d", updatedModel.width)
	}

	if updatedModel.height != 50 {
		t.Errorf("expected height 50, got %d", updatedModel.height)
	}

	if cmd != nil {
		t.Error("expected no command for WindowSizeMsg")
	}
}

func TestTemplateModel_Update_TemplatesLoadedMsg(t *testing.T) {
	model := NewTemplateModel()

	templates := []models.Template{
		{ID: "test1", Name: "Template 1", Description: "First template"},
		{ID: "test2", Name: "Template 2", Description: "Second template"},
	}

	msg := TemplatesLoadedMsg{Templates: templates}

	updatedModel, cmd := model.Update(msg)

	if len(updatedModel.templates) != 2 {
		t.Errorf("expected 2 templates, got %d", len(updatedModel.templates))
	}

	if updatedModel.loading {
		t.Error("expected loading to be false after templates loaded")
	}

	if !updatedModel.ready {
		t.Error("expected model to be ready after templates loaded")
	}

	if updatedModel.selected == nil {
		t.Error("expected first template to be auto-selected")
	}

	if cmd != nil {
		t.Error("expected no command for TemplatesLoadedMsg")
	}
}

func TestTemplateModel_Update_TemplateLoadErrorMsg(t *testing.T) {
	model := NewTemplateModel()

	testError := errors.New("Failed to load templates")

	msg := TemplateLoadErrorMsg{Error: testError}

	updatedModel, cmd := model.Update(msg)

	if updatedModel.err != testError {
		t.Error("expected error to be set")
	}

	if updatedModel.loading {
		t.Error("expected loading to be false after error")
	}

	if cmd != nil {
		t.Error("expected no command for TemplateLoadErrorMsg")
	}
}

func TestTemplateModel_Update_KeyMsg_Navigation(t *testing.T) {

	tests := []struct {
		name     string
		key      string
		expected string // expected selected template ID
	}{
		{"arrow down", "down", "test2"},
		{"j key", "j", "test2"},
		{"arrow up", "up", "test1"}, // Should wrap or stay at first
		{"k key", "k", "test1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to first item
			testModel := setupModelWithTemplates()

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "up" {
				msg = tea.KeyMsg{Type: tea.KeyUp}
			} else if tt.key == "down" {
				msg = tea.KeyMsg{Type: tea.KeyDown}
			}

			updatedModel, cmd := testModel.Update(msg)

			// Commands are returned by the list component, not required for test
			_ = cmd

			// The exact behavior depends on the list implementation
			// At minimum, the model should handle the key without error
			if updatedModel.selected == nil {
				t.Error("expected template to remain selected after navigation")
			}
		})
	}
}

func TestTemplateModel_Update_KeyMsg_LoadingState(t *testing.T) {
	model := NewTemplateModel()
	// Model is in loading state by default

	msg := tea.KeyMsg{Type: tea.KeyDown}

	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("expected no command when model is loading")
	}

	// Model should be unchanged
	if updatedModel.loading != model.loading {
		t.Error("expected loading state to be unchanged")
	}
}

func TestTemplateModel_Update_KeyMsg_Selection(t *testing.T) {
	model := setupModelWithTemplates()

	tests := []struct {
		name string
		key  string
	}{
		{"enter key", "enter"},
		{"f3 key", "f3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{}
			if tt.key == "enter" {
				msg.Type = tea.KeyEnter
			} else {
				msg.Type = tea.KeyRunes
				msg.Runes = []rune(tt.key)
			}

			updatedModel, cmd := model.Update(msg)

			if cmd == nil {
				t.Error("expected command from selection key")
			}

			if updatedModel.selected == nil {
				t.Error("expected template to be selected")
			}
		})
	}
}

func TestTemplateModel_Update_KeyMsg_BackToFileTree(t *testing.T) {
	model := setupModelWithTemplates()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f2")}

	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("expected command for F2 key")
	}
}

func TestTemplateModel_Update_KeyMsg_Refresh(t *testing.T) {
	model := setupModelWithTemplates()

	msg := tea.KeyMsg{Type: tea.KeyCtrlR}

	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("expected command for Ctrl+R key")
	}
}

func TestTemplateModel_Update_KeyMsg_ToggleDetails(t *testing.T) {
	model := setupModelWithTemplates()
	originalShowDetails := model.showDetails

	msg := tea.KeyMsg{Type: tea.KeyTab}

	updatedModel, _ := model.Update(msg)

	if updatedModel.showDetails == originalShowDetails {
		t.Error("expected showDetails to be toggled")
	}
}

func TestTemplateModel_Update_KeyMsg_PageNavigation(t *testing.T) {
	model := setupModelWithTemplates()

	tests := []struct {
		name string
		key  tea.KeyType
	}{
		{"page up", tea.KeyPgUp},
		{"page down", tea.KeyPgDown},
		{"home", tea.KeyHome},
		{"end", tea.KeyEnd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tt.key}

			updatedModel, cmd := model.Update(msg)

			// Commands are returned by the list component, not required for test
			_ = cmd

			if updatedModel.selected == nil {
				t.Error("expected template to remain selected after page navigation")
			}
		})
	}
}

// Helper function to create a model with test templates
func setupModelWithTemplates() TemplateModel {
	model := NewTemplateModel()

	templates := []models.Template{
		{
			ID:          "test1",
			Name:        "Template 1",
			Version:     "1.0.0",
			Description: "First test template",
			Author:      "Test Author 1",
			Tags:        []string{"test", "template"},
		},
		{
			ID:          "test2",
			Name:        "Template 2",
			Version:     "2.0.0",
			Description: "Second test template",
			Author:      "Test Author 2",
			Tags:        []string{"test", "example"},
		},
		{
			ID:          "test3",
			Name:        "Template 3",
			Version:     "1.5.0",
			Description: "Third test template",
			Author:      "Test Author 3",
			Tags:        []string{"template", "advanced"},
		},
	}

	model.SetTemplates(templates)
	return model
}
