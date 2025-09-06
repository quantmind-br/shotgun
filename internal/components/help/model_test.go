package help

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewHelpModel(t *testing.T) {
	model := NewHelpModel()

	if model.IsVisible() {
		t.Error("Expected help to be hidden initially")
	}

	if model.currentScreen != 0 {
		t.Errorf("Expected current screen to be 0, got %d", model.currentScreen)
	}
}

func TestHelpModel_SetVisible(t *testing.T) {
	model := NewHelpModel()

	// Test showing help
	model.SetVisible(true)
	if !model.IsVisible() {
		t.Error("Expected help to be visible after SetVisible(true)")
	}

	// Test hiding help
	model.SetVisible(false)
	if model.IsVisible() {
		t.Error("Expected help to be hidden after SetVisible(false)")
	}
}

func TestHelpModel_SetCurrentScreen(t *testing.T) {
	model := NewHelpModel()

	screens := []ScreenType{
		FileTreeScreen,
		TemplateScreen,
		TaskScreen,
		RulesScreen,
		ConfirmScreen,
		GenerateScreen,
	}

	for _, screen := range screens {
		model.SetCurrentScreen(screen)
		if model.currentScreen != screen {
			t.Errorf("Expected current screen to be %v, got %v", screen, model.currentScreen)
		}
	}
}

func TestHelpModel_UpdateSize(t *testing.T) {
	model := NewHelpModel()

	width, height := 100, 50
	model.UpdateSize(width, height)

	if model.width != width {
		t.Errorf("Expected width to be %d, got %d", width, model.width)
	}

	if model.height != height {
		t.Errorf("Expected height to be %d, got %d", height, model.height)
	}

	// Viewport should be adjusted for borders
	expectedVPWidth := width - 4
	expectedVPHeight := height - 8

	if model.viewport.Width != expectedVPWidth {
		t.Errorf("Expected viewport width to be %d, got %d", expectedVPWidth, model.viewport.Width)
	}

	if model.viewport.Height != expectedVPHeight {
		t.Errorf("Expected viewport height to be %d, got %d", expectedVPHeight, model.viewport.Height)
	}
}

func TestHelpModel_Update_F1Toggle(t *testing.T) {
	model := NewHelpModel()
	model.SetVisible(true)

	// Test F1 key toggles visibility
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f1")}
	updatedModel, _ := model.Update(msg)

	if updatedModel.IsVisible() {
		t.Error("Expected help to be hidden after F1 key")
	}
}

func TestHelpModel_formatContent(t *testing.T) {
	model := NewHelpModel()
	model.SetCurrentScreen(TaskScreen)

	content := model.formatContent()

	// Check that content includes expected sections
	if !strings.Contains(content, "Navigation:") {
		t.Error("Expected help content to contain 'Navigation:' section")
	}

	if !strings.Contains(content, "Global:") {
		t.Error("Expected help content to contain 'Global:' section")
	}

	if !strings.Contains(content, "F1") {
		t.Error("Expected help content to mention F1 key")
	}

	if !strings.Contains(content, "F2") {
		t.Error("Expected help content to mention F2 key")
	}

	if !strings.Contains(content, "F3") {
		t.Error("Expected help content to mention F3 key")
	}

	// Task and Rules screens should show F4
	if !strings.Contains(content, "F4") {
		t.Error("Expected help content for TaskScreen to mention F4 key")
	}

	// Test confirmation screen shows F10
	model.SetCurrentScreen(ConfirmScreen)
	content = model.formatContent()

	if !strings.Contains(content, "F10") {
		t.Error("Expected help content for ConfirmScreen to mention F10 key")
	}
}

func TestHelpModel_Init(t *testing.T) {
	model := NewHelpModel()
	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil command")
	}
}

func TestGetHelpContent(t *testing.T) {
	screens := []ScreenType{
		FileTreeScreen,
		TemplateScreen,
		TaskScreen,
		RulesScreen,
		ConfirmScreen,
		GenerateScreen,
	}

	for _, screen := range screens {
		content := GetHelpContent(screen)

		if len(content) == 0 {
			t.Errorf("Expected non-empty help content for screen %v", screen)
		}

		// All screens should have at least some help items
		hasContent := false
		for _, item := range content {
			if item.Context == screen || item.Context == -1 {
				hasContent = true
				break
			}
		}

		if !hasContent && screen != GenerateScreen {
			t.Errorf("Expected screen %v to have specific help content", screen)
		}
	}
}
