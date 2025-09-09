package input

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Integration tests that verify the interaction between
// different input models and the main application flow

func TestTaskInputModel_NavigationIntegration(t *testing.T) {
	// This tests the integration of task input with navigation
	model := NewTaskInputModel()

	// Initialize with window size
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, _ = model.Update(windowMsg)

	// Test complete navigation workflow
	t.Run("complete_navigation_workflow", func(t *testing.T) {
		// Set content first
		testContent := "Test task content"
		model.SetContent(testContent)

		// Test Ctrl+Left back navigation
		ctrlLeftKey := tea.KeyMsg{Type: tea.KeyCtrlLeft}
		updatedModel, cmd := model.Update(ctrlLeftKey)

		// Content should be preserved
		if updatedModel.GetContent() != testContent {
			t.Error("content should be preserved during Ctrl+Left navigation")
		}

		// Should generate back navigation command
		if cmd == nil {
			t.Error("expected back navigation command")
		}

		// Execute command to verify message type
		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(BackToTemplateMsg); !ok {
				t.Error("expected BackToTemplateMsg from Ctrl+Left navigation")
			}
		}
	})

	// Test advancement with validation
	t.Run("advancement_with_validation", func(t *testing.T) {
		// Test Alt+C advancement with valid content
		model.SetContent("Valid content for advancement")

		altCKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c"), Alt: true}
		updatedModel, cmd := model.Update(altCKey)

		// Should clear any errors
		if updatedModel.GetError() != nil {
			t.Error("expected error to be cleared for valid content")
		}

		// Should generate advancement command
		if cmd == nil {
			t.Error("expected advancement command for valid content")
		}
	})

	// Test advancement failure with validation
	t.Run("advancement_failure_with_validation", func(t *testing.T) {
		// Test Alt+C advancement with empty content (should fail validation)
		model.SetContent("")

		altCKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c"), Alt: true}
		updatedModel, _ := model.Update(altCKey)

		// Should set error for empty content
		if updatedModel.GetError() == nil {
			t.Error("expected validation error for empty content")
		}
	})
}

func TestRulesInputModel_NavigationIntegration(t *testing.T) {
	// This tests the integration of rules input with navigation
	model := NewRulesInputModel()

	// Initialize with window size
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, _ = model.Update(windowMsg)

	t.Run("back_navigation", func(t *testing.T) {
		// Test back navigation
		ctrlLeftKey := tea.KeyMsg{Type: tea.KeyCtrlLeft}
		_, cmd := model.Update(ctrlLeftKey)

		if cmd == nil {
			t.Error("expected back navigation command")
		}

		// Verify message type
		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(BackToTaskMsg); !ok {
				t.Error("expected BackToTaskMsg from Ctrl+Left navigation")
			}
		}
	})

	t.Run("forward_navigation", func(t *testing.T) {
		// Test forward navigation (rules are optional)
		altCKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c"), Alt: true}
		_, cmd := model.Update(altCKey)

		if cmd == nil {
			t.Error("expected forward navigation command")
		}
	})

	t.Run("skip_navigation", func(t *testing.T) {
		// Test skip navigation
		altSKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s"), Alt: true}
		_, cmd := model.Update(altSKey)

		if cmd == nil {
			t.Error("expected skip navigation command")
		}
	})
}

func TestIntegration_StatePreservation(t *testing.T) {
	// Test that state is preserved across navigation
	taskModel := NewTaskInputModel()
	rulesModel := NewRulesInputModel()

	// Set content in models
	taskContent := "Task content to preserve"
	rulesContent := "Rules content to preserve"

	taskModel.SetContent(taskContent)
	rulesModel.SetContent(rulesContent)

	// Verify content is preserved
	if taskModel.GetContent() != taskContent {
		t.Error("task content not preserved")
	}

	if rulesModel.GetContent() != rulesContent {
		t.Error("rules content not preserved")
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	// Test error handling across models
	taskModel := NewTaskInputModel()

	// Set an error
	testError := errors.New("Test error message")
	taskModel.SetError(testError)

	// Verify error is set
	if taskModel.GetError() == nil {
		t.Error("error not set")
	}

	if taskModel.GetError().Error() != testError.Error() {
		t.Errorf("expected error '%s', got '%v'", testError, taskModel.GetError())
	}
}
