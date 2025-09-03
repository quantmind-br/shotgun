package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// GlobalKeyHandler processes global navigation keys
func (a *AppState) GlobalKeyHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "f1":
		return a.showHelp()
	case "f2":
		return a.goToPreviousScreen()
	case "f3":
		return a.goToNextScreen()
	case "esc":
		return a.showExitDialog()
	case "q", "ctrl+c":
		// Allow quit shortcut
		return a.showExitDialog()
	default:
		// Key not handled globally, let current screen handle it
		return nil, nil
	}
}

// showHelp displays contextual help for the current screen
func (a *AppState) showHelp() (tea.Model, tea.Cmd) {
	// Set help dialog state
	a.ShowingHelp = true
	a.HelpContent = a.getHelpContent()
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
	// Validate current screen before advancing
	if !a.canAdvance() {
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
		// At final screen, could trigger completion
		return a, tea.Quit
	}
	
	return a, nil
}

// showExitDialog shows exit confirmation dialog
func (a *AppState) showExitDialog() (tea.Model, tea.Cmd) {
	// For now, we'll implement a simple quit
	// Later we can add a confirmation dialog
	return a, tea.Quit
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
• F1: Show this help
• F2: Go to previous screen (disabled on first screen)
• F3: Go to next screen (requires selected files)
• ESC/q: Exit application

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
• F1: Show this help
• F2: Go to previous screen
• F3: Go to next screen (requires selected template)
• ESC/q: Exit application

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
• F1: Show this help
• F2: Go to previous screen
• F3: Go to next screen (requires task description)
• ESC/q: Exit application

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
• F1: Show this help
• F2: Go to previous screen
• F3: Go to next screen (rules are optional)
• ESC/q: Exit application

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
• F1: Show this help
• F2: Go to previous screen
• F3: Generate output
• ESC/q: Exit application

Action:
• Press F3 to generate the final output
• The application will create your prompt based on all inputs`

	default:
		return "Help not available for this screen."
	}
}

// IsGlobalKey checks if a key should be handled globally
func IsGlobalKey(key string) bool {
	globalKeys := []string{"f1", "f2", "f3", "esc", "q", "ctrl+c"}
	for _, gk := range globalKeys {
		if key == gk {
			return true
		}
	}
	return false
}