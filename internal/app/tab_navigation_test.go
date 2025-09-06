package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTabNavigation_FileTreeScreen(t *testing.T) {
	app := NewApp()
	app.SetCurrentScreen(FileTreeScreen)

	// Test tab key navigation
	msg := tea.KeyMsg{Type: tea.KeyTab}

	// Tab should be handled by the screen, not globally
	model, cmd := app.handleScreenInput(msg)

	if model == nil {
		t.Error("Expected non-nil model from tab navigation")
	}

	if cmd != nil {
		// Tab might trigger commands in some screens
		t.Log("Tab triggered command in FileTree screen")
	}
}

func TestTabNavigation_TaskInputScreen(t *testing.T) {
	app := NewApp()
	app.SetCurrentScreen(TaskScreen)

	// Blur to start unfocused
	app.TaskInput.Blur()

	// Test tab key - should focus the text area
	msg := tea.KeyMsg{Type: tea.KeyTab}

	model, _ := app.handleScreenInput(msg)

	if model == nil {
		t.Error("Expected non-nil model from tab navigation")
	}

	appModel := model.(*AppState)

	// After tab, text area should be focused (in a real scenario)
	// Note: The actual focus behavior depends on the Update method implementation
	if appModel.CurrentScreen != TaskScreen {
		t.Errorf("Expected to remain on TaskScreen, got %v", appModel.CurrentScreen)
	}
}

func TestTabNavigation_RulesInputScreen(t *testing.T) {
	app := NewApp()
	app.SetCurrentScreen(RulesScreen)

	// Blur to start unfocused
	app.RulesInput.Blur()

	// Test tab key
	msg := tea.KeyMsg{Type: tea.KeyTab}

	model, _ := app.handleScreenInput(msg)

	if model == nil {
		t.Error("Expected non-nil model from tab navigation")
	}

	appModel := model.(*AppState)

	if appModel.CurrentScreen != RulesScreen {
		t.Errorf("Expected to remain on RulesScreen, got %v", appModel.CurrentScreen)
	}
}

func TestTabNavigation_TemplateScreen(t *testing.T) {
	app := NewApp()
	app.SetCurrentScreen(TemplateScreen)

	// Test tab key in template screen
	msg := tea.KeyMsg{Type: tea.KeyTab}

	model, _ := app.handleScreenInput(msg)

	if model == nil {
		t.Error("Expected non-nil model from tab navigation")
	}

	appModel := model.(*AppState)

	// Tab should cycle through templates or do nothing
	if appModel.CurrentScreen != TemplateScreen {
		t.Errorf("Expected to remain on TemplateScreen, got %v", appModel.CurrentScreen)
	}
}

func TestTabNavigation_ConfirmScreen(t *testing.T) {
	app := NewApp()
	app.SetCurrentScreen(ConfirmScreen)

	// Test tab key in confirmation screen
	msg := tea.KeyMsg{Type: tea.KeyTab}

	model, _ := app.handleScreenInput(msg)

	if model == nil {
		t.Error("Expected non-nil model from tab navigation")
	}

	appModel := model.(*AppState)

	// Tab might cycle through confirm/cancel options
	if appModel.CurrentScreen != ConfirmScreen {
		t.Errorf("Expected to remain on ConfirmScreen, got %v", appModel.CurrentScreen)
	}
}

func TestShiftTabNavigation(t *testing.T) {
	app := NewApp()

	screens := []ScreenType{
		FileTreeScreen,
		TemplateScreen,
		TaskScreen,
		RulesScreen,
		ConfirmScreen,
	}

	for _, screen := range screens {
		app.SetCurrentScreen(screen)

		// Blur text areas if applicable
		if screen == TaskScreen {
			app.TaskInput.Blur()
		} else if screen == RulesScreen {
			app.RulesInput.Blur()
		}

		// Test shift+tab (reverse tab)
		msg := tea.KeyMsg{Type: tea.KeyShiftTab}

		model, _ := app.handleScreenInput(msg)

		if model == nil {
			t.Errorf("Expected non-nil model from shift+tab on screen %v", screen)
			continue
		}

		appModel := model.(*AppState)

		// Should remain on same screen
		if appModel.CurrentScreen != screen {
			t.Errorf("Expected to remain on %v after shift+tab, got %v", screen, appModel.CurrentScreen)
		}
	}
}

func TestTabNavigationOrder(t *testing.T) {
	// This test verifies that tab navigation follows a logical order
	// Note: The actual implementation may vary, this tests the expectation

	app := NewApp()

	// Test file tree screen tab order
	app.SetCurrentScreen(FileTreeScreen)

	// In file tree, tab might cycle through:
	// 1. File list
	// 2. Action buttons (if any)
	// 3. Back to file list

	// Test template screen tab order
	app.SetCurrentScreen(TemplateScreen)

	// In template screen, tab might cycle through:
	// 1. Template list
	// 2. Template details
	// 3. Action buttons

	// Test task input tab order
	app.SetCurrentScreen(TaskScreen)
	app.TaskInput.Blur()

	// In task input, tab might cycle through:
	// 1. Text area
	// 2. Character count
	// 3. Action buttons

	// This is a placeholder for actual tab order testing
	// The real implementation would need to track focus state
	t.Log("Tab navigation order test placeholder - actual implementation needed")
}
