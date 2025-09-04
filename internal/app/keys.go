package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/components/help"
)

// GlobalKeyHandler processes global navigation keys
func (a *AppState) GlobalKeyHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Skip global keys during text editing except ESC
	if a.isInInputMode() && msg.String() != "esc" {
		return nil, nil
	}

	switch msg.String() {
	case "f1":
		return a.showHelp()
	case "f2":
		if a.canGoToPreviousScreen() {
			return a.goToPreviousScreen()
		}
	case "f3":
		if a.canGoToNextScreen() {
			return a.goToNextScreen()
		}
	case "f4":
		return a.handleF4Action()
	case "f10":
		return a.handleF10Generate()
	case "esc":
		return a.showExitDialog()
	case "q", "ctrl+c":
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
	globalKeys := []string{"f1", "f2", "f3", "f4", "f10", "esc", "q", "ctrl+c"}
	for _, gk := range globalKeys {
		if key == gk {
			return true
		}
	}
	return false
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
		return len(a.SelectedFiles) > 0
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
