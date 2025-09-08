package app

import (
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/diogopedro/shotgun/internal/components/help"
)

// GlobalKeyHandler processes global navigation keys
func (a *AppState) GlobalKeyHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    key := normalizeKey(msg)

    // In input mode, allow only safe global combos: ESC and Ctrl+Q
    if a.isInInputMode() && key != "esc" && key != "ctrl+q" && key != "ctrl+c" {
        return nil, nil
    }

    switch key {
    case "ctrl+h":
        return a.showHelp()
    case "ctrl+left":
        if a.canGoToPreviousScreen() {
            return a.goToPreviousScreen()
        }
    case "ctrl+enter":
        if a.canGoToNextScreen() {
            return a.goToNextScreen()
        }
    case "ctrl+q", "esc", "ctrl+c":
        return a.showExitDialog()
    default:
        // Key not handled globally, let current screen handle it
        return nil, nil
    }

    // Return current state if navigation not allowed
    return a, nil
}

// showHelp displays contextual help for the current screen
func (a *AppState) showHelp() (tea.Model, tea.Cmd) {
	// Toggle help visibility
	a.Help.SetVisible(!a.Help.IsVisible())
	a.Help.SetCurrentScreen(help.ScreenType(a.CurrentScreen))
	a.ShowingHelp = a.Help.IsVisible()
	a.Error = nil // Clear any errors

	return a, nil
}

// goToPreviousScreen navigates to the previous screen if allowed
func (a *AppState) goToPreviousScreen() (tea.Model, tea.Cmd) {
	// Check if we can go back
	if a.CurrentScreen == FileTreeScreen {
		// Already at first screen
		return a, nil
	}

	// Navigate to previous screen
	switch a.CurrentScreen {
	case TemplateScreen:
		a.SetCurrentScreen(FileTreeScreen)
	case TaskScreen:
		a.SetCurrentScreen(TemplateScreen)
	case RulesScreen:
		a.SetCurrentScreen(TaskScreen)
	case ConfirmScreen:
		a.SetCurrentScreen(RulesScreen)
	}

	return a, nil
}

// goToNextScreen navigates to the next screen with validation
func (a *AppState) goToNextScreen() (tea.Model, tea.Cmd) {
	// Use the new validation method
	if !a.canGoToNextScreen() {
		// Set error message for validation failure
		a.Error = a.getValidationError()
		return a, nil
	}

	// Clear any previous errors
	a.Error = nil

	// Navigate to next screen
	switch a.CurrentScreen {
	case FileTreeScreen:
		a.SetCurrentScreen(TemplateScreen)
	case TemplateScreen:
		a.SetCurrentScreen(TaskScreen)
	case TaskScreen:
		a.SetCurrentScreen(RulesScreen)
	case RulesScreen:
		a.SetCurrentScreen(ConfirmScreen)
	case ConfirmScreen:
		// Use F10 for generation from confirmation screen
		return a, nil
	}

	return a, nil
}

// showExitDialog shows exit confirmation dialog
func (a *AppState) showExitDialog() (tea.Model, tea.Cmd) {
	// Show exit confirmation dialog
	a.ShowingExit = true
	return a, nil
}

// getHelpContent returns contextual help content for current screen
func (a *AppState) getHelpContent() string {
    switch a.CurrentScreen {
    case FileTreeScreen:
        return `File Tree Navigation Help:

Navigation:
• ↑/↓ or k/j: Move cursor up/down
• ←/→ or h/l: Collapse/expand directories  
• Space: Toggle file/directory selection

Global Keys:
• Ctrl+H: Show this help
• Ctrl+Left: Go to previous screen (disabled on first screen)
• Ctrl+Enter: Go to next screen (requires selected files)
• Ctrl+Q/ESC: Exit application

File Selection:
• Select files and directories you want to include
• Binary files are automatically excluded
• At least one file must be selected to continue`

    case TemplateScreen:
        return `Template Selection Help:

Navigation:
• ↑/↓ or k/j: Move cursor up/down
• Enter/Space: Select template

Global Keys:
• Ctrl+H: Show this help
• Ctrl+Left: Go to previous screen
• Ctrl+Enter: Go to next screen (requires selected template)
• Ctrl+Q/ESC: Exit application

Template Types:
• Choose a template that matches your task
• Templates provide structured prompts for AI generation`

    case TaskScreen:
        return `Task Input Help:

Input:
• Type your task description
• Be specific about what you want to accomplish
• Include any special requirements or constraints

Global Keys:
• Ctrl+H: Show this help
• Ctrl+Left: Go to previous screen
• Ctrl+Enter: Go to next screen (requires task description)
• Ctrl+Q/ESC: Exit application

Tips:
• Describe your goal clearly
• Mention specific technologies or patterns if relevant`

    case RulesScreen:
        return `Rules Input Help:

Input:
• Add custom rules or constraints (optional)
• Specify coding standards or preferences
• Add any special requirements

Global Keys:
• Ctrl+H: Show this help
• Ctrl+Left: Go to previous screen
• Ctrl+Enter: Go to next screen (rules are optional)
• Ctrl+Q/ESC: Exit application

Examples:
• "Use TypeScript strict mode"
• "Follow REST API conventions"
• "Include comprehensive error handling"`

    case ConfirmScreen:
        return `Confirmation Help:

Review:
• Review all your selections and inputs
• Check that everything is correct
• Make changes using F2 to go back if needed

Global Keys:
• Ctrl+H: Show this help
• Ctrl+Left: Go to previous screen
• Ctrl+Enter: Generate output
• Ctrl+Q/ESC: Exit application

Action:
• Press Ctrl+Enter to generate the final output
• The application will create your prompt based on all inputs`

	default:
		return "Help not available for this screen."
	}
}

// IsGlobalKey checks if a key should be handled globally
func IsGlobalKey(key string) bool {
    // Note: ctrl+enter is handled per-screen (e.g., FileTree, Confirm)
    globalKeys := []string{"ctrl+h", "ctrl+left", "ctrl+q", "esc", "ctrl+c"}
    for _, gk := range globalKeys {
        if key == gk {
            return true
        }
    }
    return false
}

// isFunctionKeyMsg returns true if msg.Type is one of the function keys we handle
func isFunctionKeyMsg(msg tea.KeyMsg) bool { return false }

// normalizeKey maps tea.KeyMsg to a stable string used by handlers
func normalizeKey(msg tea.KeyMsg) string {
    switch msg.Type {
    case tea.KeyEsc:
        return "esc"
    }
    return strings.ToLower(msg.String())
}

// isInInputMode checks if user is currently editing text
func (a *AppState) isInInputMode() bool {
	switch a.CurrentScreen {
	case TaskScreen:
		return a.TaskInput.Focused()
	case RulesScreen:
		return a.RulesInput.Focused()
	default:
		return false
	}
}

// canGoToPreviousScreen validates if backward navigation is allowed
func (a *AppState) canGoToPreviousScreen() bool {
	switch a.CurrentScreen {
	case FileTreeScreen:
		return false // First screen
	case TemplateScreen:
		return true
	case TaskScreen:
		return true
	case RulesScreen:
		return true
	case ConfirmScreen:
		return true
	case GenerateScreen:
		return true
	default:
		return false
	}
}

// canGoToNextScreen validates if forward navigation is allowed
func (a *AppState) canGoToNextScreen() bool {
    switch a.CurrentScreen {
    case FileTreeScreen:
        // Derive directly from the FileTree model to reflect live selections
        return len(a.FileTree.GetSelectedFiles()) > 0
	case TemplateScreen:
		return a.SelectedTemplate != nil
	case TaskScreen:
		return len(a.TaskContent) > 0
	case RulesScreen:
		return true // Rules are optional
	case ConfirmScreen:
		return false // Use F10 for generation
	case GenerateScreen:
		return false // Last screen
	default:
		return false
	}
}

// handleF4Action provides skip functionality for input screens
func (a *AppState) handleF4Action() (tea.Model, tea.Cmd) {
	switch a.CurrentScreen {
	case TaskScreen:
		// Skip to rules input
		a.SetCurrentScreen(RulesScreen)
		return a, nil
	case RulesScreen:
		// Skip rules, go to confirmation
		a.SetCurrentScreen(ConfirmScreen)
		return a, nil
	default:
		return a, nil
	}
}

// handleF10Generate triggers prompt generation from confirmation screen
func (a *AppState) handleF10Generate() (tea.Model, tea.Cmd) {
	if a.CurrentScreen == ConfirmScreen && a.canGenerate() {
		cmd := a.generatePrompt()
		return a, cmd
	}
	return a, nil
}

// canGenerate validates all required data is present for generation
func (a *AppState) canGenerate() bool {
	return len(a.SelectedFiles) > 0 && a.SelectedTemplate != nil && len(a.TaskContent) > 0
}
