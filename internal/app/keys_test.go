package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
)

func TestGlobalKeyHandler_CtrlH(t *testing.T) {
	app := NewApp()

	msg := tea.KeyMsg{Type: tea.KeyCtrlH}

	model, _ := app.GlobalKeyHandler(msg)

	if model == nil {
		t.Error("Expected non-nil model from Ctrl+H handler")
		return
	}

	// Ctrl+H should toggle help
	appModel := model.(*AppState)
	if appModel.Error != nil {
		t.Errorf("Expected no error, got: %v", appModel.Error)
	}
	if !appModel.ShowingHelp {
		t.Error("Expected help to be shown after Ctrl+H")
	}
}

func TestGlobalKeyHandler_CtrlLeft(t *testing.T) {
	app := NewApp()

	// Start at Template screen (can go back to FileTree)
	app.SetCurrentScreen(TemplateScreen)

	msg := tea.KeyMsg{Type: tea.KeyCtrlLeft}

	model, _ := app.GlobalKeyHandler(msg)

	if model == nil {
		t.Error("Expected non-nil model from Ctrl+Left handler")
		return
	}

	appModel := model.(*AppState)
	if appModel.CurrentScreen != FileTreeScreen {
		t.Errorf("Expected FileTreeScreen after Ctrl+Left, got %v", appModel.CurrentScreen)
	}
}

func TestGlobalKeyHandler_CtrlLeftAtFirstScreen(t *testing.T) {
	app := NewApp()

	// Already at FileTree (first screen)
	app.SetCurrentScreen(FileTreeScreen)

	msg := tea.KeyMsg{Type: tea.KeyCtrlLeft}

	model, _ := app.GlobalKeyHandler(msg)

	if model == nil {
		t.Error("Expected non-nil model from Ctrl+Left handler")
		return
	}

	appModel := model.(*AppState)
	// Ctrl+Left at first screen should do nothing
	if appModel.CurrentScreen != FileTreeScreen {
		t.Errorf("Expected to stay at FileTreeScreen after Ctrl+Left, got %v", appModel.CurrentScreen)
	}
}

func TestGlobalKeyHandler_ESCExit(t *testing.T) {
	app := NewApp()

	msg := tea.KeyMsg{Type: tea.KeyEsc}

	model, cmd := app.GlobalKeyHandler(msg)

	if model == nil {
		t.Error("Expected non-nil model from ESC handler")
		return
	}

	// ESC now shows exit dialog instead of immediately quitting
	appModel := model.(*AppState)
	if !appModel.ShowingExit {
		t.Error("Expected ShowingExit to be true after ESC")
	}

	if cmd != nil {
		t.Error("Expected nil command from ESC handler (shows dialog instead)")
	}
}

func TestGlobalKeyHandler_CtrlQExit(t *testing.T) {
	app := NewApp()

	msg := tea.KeyMsg{Type: tea.KeyCtrlQ}

	model, cmd := app.GlobalKeyHandler(msg)

	if model == nil {
		t.Error("Expected non-nil model from Ctrl+Q handler")
		return
	}

	// Ctrl+Q now shows exit dialog instead of immediately quitting
	appModel := model.(*AppState)
	if !appModel.ShowingExit {
		t.Error("Expected ShowingExit to be true after Ctrl+Q")
	}

	if cmd != nil {
		t.Error("Expected nil command from Ctrl+Q handler (shows dialog instead)")
	}
}

func TestGlobalKeyHandler_CtrlCExit(t *testing.T) {
	app := NewApp()

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}

	model, cmd := app.GlobalKeyHandler(msg)

	if model == nil {
		t.Error("Expected non-nil model from Ctrl+C handler")
		return
	}

	// Ctrl+C now shows exit dialog instead of immediately quitting
	appModel := model.(*AppState)
	if !appModel.ShowingExit {
		t.Error("Expected ShowingExit to be true after Ctrl+C")
	}

	if cmd != nil {
		t.Error("Expected nil command from Ctrl+C handler (shows dialog instead)")
	}
}

func TestGlobalKeyHandler_UnknownKey(t *testing.T) {
	app := NewApp()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}

	model, cmd := app.GlobalKeyHandler(msg)

	// Unknown keys should return nil to let screen handle them
	if model != nil {
		t.Error("Expected nil model for unknown global key")
	}

	if cmd != nil {
		t.Error("Expected nil command for unknown global key")
	}
}

func TestGoToPreviousScreen_AllScreens(t *testing.T) {
	app := NewApp()

	tests := []struct {
		current  ScreenType
		expected ScreenType
	}{
		{TemplateScreen, FileTreeScreen},
		{TaskScreen, TemplateScreen},
		{RulesScreen, TaskScreen},
		{ConfirmScreen, RulesScreen},
		{FileTreeScreen, FileTreeScreen}, // Can't go back further
	}

	for _, tt := range tests {
		app.SetCurrentScreen(tt.current)

		model, _ := app.goToPreviousScreen()
		appModel := model.(*AppState)

		if appModel.CurrentScreen != tt.expected {
			t.Errorf("From %v, expected previous screen %v, got %v",
				tt.current, tt.expected, appModel.CurrentScreen)
		}
	}
}

func TestGoToNextScreen_AllScreens(t *testing.T) {
	app := NewApp()

	tests := []struct {
		current       ScreenType
		expected      ScreenType
		setupFunc     func(*AppState)
		shouldAdvance bool
	}{
		{
			FileTreeScreen,
			FileTreeScreen, // Won't advance because FileTree.GetSelectedFiles() returns empty
			func(a *AppState) {
				// Can't easily mock FileTree.GetSelectedFiles()
				// so this test will fail validation
				a.SelectedFiles = []string{"/test/file1.txt"}
			},
			false, // Will fail because FileTree.GetSelectedFiles() is empty
		},
		{
			TemplateScreen,
			TaskScreen,
			func(a *AppState) {
				a.SelectedTemplate = &models.Template{Name: "Test"}
			},
			true,
		},
		{
			TaskScreen,
			RulesScreen,
			func(a *AppState) {
				a.TaskContent = "Test task"
			},
			true,
		},
		{
			RulesScreen,
			ConfirmScreen,
			func(a *AppState) {
				// Rules are optional
			},
			true,
		},
	}

	for _, tt := range tests {
		app.SetCurrentScreen(tt.current)
		tt.setupFunc(app)

		model, _ := app.goToNextScreen()
		appModel := model.(*AppState)

		if tt.shouldAdvance {
			if appModel.CurrentScreen != tt.expected {
				t.Errorf("From %v, expected next screen %v, got %v",
					tt.current, tt.expected, appModel.CurrentScreen)
			}

			if appModel.Error != nil {
				t.Errorf("Expected no error advancing from %v, got: %v", tt.current, appModel.Error)
			}
		} else {
			// Should remain on current screen with error
			if appModel.CurrentScreen != tt.current {
				t.Errorf("From %v, expected to remain on current screen due to validation, got %v",
					tt.current, appModel.CurrentScreen)
			}

			// In the new implementation, validation failure doesn't set an error
			// It just doesn't advance the screen
		}
	}
}

func TestGoToNextScreen_FinalScreen(t *testing.T) {
	app := NewApp()
	app.SetCurrentScreen(ConfirmScreen)

	model, cmd := app.goToNextScreen()

	if model == nil {
		t.Error("Expected non-nil model from final screen navigation")
	}

	// Confirmation screen doesn't advance further
	if cmd != nil {
		t.Error("Expected nil command from confirmation screen navigation")
	}

	appModel := model.(*AppState)
	if appModel.CurrentScreen != ConfirmScreen {
		t.Errorf("Expected to remain on ConfirmScreen, got %v", appModel.CurrentScreen)
	}
}
