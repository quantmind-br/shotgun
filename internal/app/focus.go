package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// FocusState represents the focus state for a screen
type FocusState struct {
	CursorPosition int
	ScrollOffset   int
	SelectionState map[string]bool
	TextContent    string
}

// saveFocusState saves the current screen's focus state
func (a *AppState) saveFocusState() *FocusState {
	switch a.CurrentScreen {
	case FileTreeScreen:
		return &FocusState{
			CursorPosition: a.getFocusPosition(&a.FileTree),
			ScrollOffset:   a.getScrollOffset(&a.FileTree),
		}
	case TemplateScreen:
		return &FocusState{
			CursorPosition: 0, // Template list manages its own cursor internally
		}
	case TaskScreen:
		return &FocusState{
			CursorPosition: a.TaskInput.cursor,
			TextContent:    a.TaskInput.content,
		}
	case RulesScreen:
		return &FocusState{
			CursorPosition: a.RulesInput.cursor,
			TextContent:    a.RulesInput.content,
		}
	case ConfirmScreen:
		return &FocusState{
			ScrollOffset: a.getConfirmScrollOffset(),
		}
	default:
		return &FocusState{}
	}
}

// restoreFocusState restores focus state for the current screen
func (a *AppState) restoreFocusState(state *FocusState) tea.Cmd {
	if state == nil {
		return a.initializeScreenFocus()
	}

	switch a.CurrentScreen {
	case FileTreeScreen:
		return a.restoreFileTreeFocus(state)
	case TemplateScreen:
		return a.restoreTemplateFocus(state)
	case TaskScreen:
		return a.restoreTaskInputFocus(state)
	case RulesScreen:
		return a.restoreRulesInputFocus(state)
	case ConfirmScreen:
		return a.restoreConfirmationFocus(state)
	default:
		return nil
	}
}

// initializeScreenFocus sets up initial focus for a screen
func (a *AppState) initializeScreenFocus() tea.Cmd {
	switch a.CurrentScreen {
	case FileTreeScreen:
		// FileTree manages its own focus internally
		return nil
	case TemplateScreen:
		// Template screen manages its own focus internally
		return nil
	case TaskScreen:
		// Focus on text input, cursor at end of existing content
		a.TaskInput.cursor = len(a.TaskInput.content)
		return nil
	case RulesScreen:
		// Focus on text input, cursor at end of existing content
		a.RulesInput.cursor = len(a.RulesInput.content)
		return nil
	case ConfirmScreen:
		// Scroll to top of confirmation screen
		return nil
	default:
		return nil
	}
}

// Focus restoration helpers for each screen type

func (a *AppState) restoreFileTreeFocus(state *FocusState) tea.Cmd {
	// FileTree has its own cursor management
	// We can't directly set cursor but this provides the interface
	return nil
}

func (a *AppState) restoreTemplateFocus(state *FocusState) tea.Cmd {
	// Template list manages its own cursor internally
	// No direct cursor access needed
	return nil
}

func (a *AppState) restoreTaskInputFocus(state *FocusState) tea.Cmd {
	// Restore text content if provided
	if state.TextContent != "" {
		a.TaskInput.content = state.TextContent
		a.TaskContent = state.TextContent
	}

	// Restore cursor position
	if state.CursorPosition >= 0 && state.CursorPosition <= len(a.TaskInput.content) {
		a.TaskInput.cursor = state.CursorPosition
	} else {
		a.TaskInput.cursor = len(a.TaskInput.content)
	}

	return nil
}

func (a *AppState) restoreRulesInputFocus(state *FocusState) tea.Cmd {
	// Restore text content if provided
	if state.TextContent != "" {
		a.RulesInput.content = state.TextContent
		a.RulesContent = state.TextContent
	}

	// Restore cursor position
	if state.CursorPosition >= 0 && state.CursorPosition <= len(a.RulesInput.content) {
		a.RulesInput.cursor = state.CursorPosition
	} else {
		a.RulesInput.cursor = len(a.RulesInput.content)
	}

	return nil
}

func (a *AppState) restoreConfirmationFocus(state *FocusState) tea.Cmd {
	// Restore scroll position if supported
	return nil
}

// Helper methods to get focus information

func (a *AppState) getFocusPosition(model interface{}) int {
	// This would need to be implemented based on the specific model
	// For now, return 0
	return 0
}

func (a *AppState) getScrollOffset(model interface{}) int {
	// This would need to be implemented based on the specific model
	// For now, return 0
	return 0
}

func (a *AppState) getConfirmScrollOffset() int {
	// This would need to be implemented based on the confirmation model
	// For now, return 0
	return 0
}

// Screen initialization commands for proper focus setup

// InitScreenCmd returns initialization command for current screen
func (a *AppState) InitScreenCmd() tea.Cmd {
	return a.initializeScreenFocus()
}

// CleanupScreenCmd performs cleanup when leaving a screen
func (a *AppState) CleanupScreenCmd() tea.Cmd {
	switch a.CurrentScreen {
	case FileTreeScreen:
		// FileTree cleanup if needed
		return nil
	case TemplateScreen:
		// Template cleanup if needed
		return nil
	case TaskScreen:
		// Save task content before leaving
		a.TaskContent = a.TaskInput.content
		return nil
	case RulesScreen:
		// Save rules content before leaving
		a.RulesContent = a.RulesInput.content
		return nil
	case ConfirmScreen:
		// Confirmation cleanup if needed
		return nil
	default:
		return nil
	}
}

// IsFocused returns whether the current screen has focus
func (a *AppState) IsFocused() bool {
	// If showing dialogs, main screen doesn't have focus
	if a.ShowingHelp || a.ShowingExit {
		return false
	}
	return true
}

// CanReceiveInput returns whether current screen can receive keyboard input
func (a *AppState) CanReceiveInput() bool {
	if !a.IsFocused() {
		return false
	}

	switch a.CurrentScreen {
	case TaskScreen, RulesScreen:
		// Text input screens can always receive input
		return true
	case FileTreeScreen, TemplateScreen, ConfirmScreen:
		// Navigation screens can receive input
		return true
	default:
		return false
	}
}
