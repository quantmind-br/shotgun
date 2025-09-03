package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model interface
func (a *AppState) View() string {
	var content strings.Builder
	
	// Render progress indicator at top
	content.WriteString(a.Progress.View())
	content.WriteString("\n")
	
	// Show error if present
	if a.Error != nil {
		content.WriteString(a.renderError())
		content.WriteString("\n")
	}
	
	// Render current screen content
	content.WriteString(a.renderCurrentScreen())
	
	// Render help dialog if showing
	if a.ShowingHelp {
		return a.renderHelpDialog()
	}
	
	// Render exit dialog if showing
	if a.ShowingExit {
		return a.renderExitDialog()
	}
	
	return content.String()
}

// renderCurrentScreen renders the active screen
func (a *AppState) renderCurrentScreen() string {
	switch a.CurrentScreen {
	case FileTreeScreen:
		return a.FileTree.View()
	case TemplateScreen:
		return a.renderTemplateScreen()
	case TaskScreen:
		return a.renderTaskScreen()
	case RulesScreen:
		return a.renderRulesScreen()
	case ConfirmScreen:
		return a.renderConfirmationScreen()
	default:
		return "Unknown screen"
	}
}

// Screen renderers

func (a *AppState) renderTemplateScreen() string {
	var content strings.Builder
	
	content.WriteString(a.styleScreenTitle("Choose Template"))
	content.WriteString("\n\n")
	
	if len(a.Template.items) == 0 {
		content.WriteString("Loading templates...\n")
		return content.String()
	}
	
	// Render template list
	for i, template := range a.Template.items {
		var line strings.Builder
		
		// Selection indicator
		if i == a.Template.cursor {
			line.WriteString("→ ")
		} else {
			line.WriteString("  ")
		}
		
		// Template name
		name := template.Name
		if a.Template.selected == template {
			name = a.styleSelectedItem(name)
		} else if i == a.Template.cursor {
			name = a.styleCursorItem(name)
		}
		line.WriteString(name)
		
		// Template description
		if template.Description != "" {
			line.WriteString(" - ")
			line.WriteString(a.styleDescription(template.Description))
		}
		
		content.WriteString(line.String())
		content.WriteString("\n")
	}
	
	content.WriteString("\n")
	content.WriteString(a.renderNavigationHelp())
	
	return content.String()
}

func (a *AppState) renderTaskScreen() string {
	var content strings.Builder
	
	content.WriteString(a.styleScreenTitle("Describe Your Task"))
	content.WriteString("\n\n")
	
	content.WriteString("Enter a detailed description of what you want to accomplish:\n\n")
	
	// Render text input with cursor
	inputText := a.TaskInput.content
	if len(inputText) == 0 {
		inputText = a.stylePlaceholder("Type your task description here...")
	} else {
		// Insert cursor
		if a.TaskInput.cursor <= len(inputText) {
			inputText = inputText[:a.TaskInput.cursor] + a.styleCursor("│") + inputText[a.TaskInput.cursor:]
		}
	}
	
	content.WriteString(a.styleInputBox(inputText))
	content.WriteString("\n\n")
	
	content.WriteString("Examples:\n")
	content.WriteString("• Create a REST API for user management\n")
	content.WriteString("• Refactor this component to use hooks\n")
	content.WriteString("• Add error handling to the payment flow\n\n")
	
	content.WriteString(a.renderNavigationHelp())
	
	return content.String()
}

func (a *AppState) renderRulesScreen() string {
	var content strings.Builder
	
	content.WriteString(a.styleScreenTitle("Add Custom Rules (Optional)"))
	content.WriteString("\n\n")
	
	content.WriteString("Add any specific rules or constraints:\n\n")
	
	// Render text input with cursor
	inputText := a.RulesInput.content
	if len(inputText) == 0 {
		inputText = a.stylePlaceholder("Add custom rules or leave empty...")
	} else {
		// Insert cursor
		if a.RulesInput.cursor <= len(inputText) {
			inputText = inputText[:a.RulesInput.cursor] + a.styleCursor("│") + inputText[a.RulesInput.cursor:]
		}
	}
	
	content.WriteString(a.styleInputBox(inputText))
	content.WriteString("\n\n")
	
	content.WriteString("Examples:\n")
	content.WriteString("• Use TypeScript strict mode\n")
	content.WriteString("• Follow company coding standards\n")
	content.WriteString("• Include comprehensive error handling\n\n")
	
	content.WriteString(a.renderNavigationHelp())
	
	return content.String()
}

func (a *AppState) renderConfirmationScreen() string {
	var content strings.Builder
	
	content.WriteString(a.styleScreenTitle("Review & Confirm"))
	content.WriteString("\n\n")
	
	// Build summary
	summary := a.Confirmation.summary
	if summary == "" {
		a.buildConfirmationSummary()
		summary = a.Confirmation.summary
	}
	
	content.WriteString(a.styleConfirmationBox(summary))
	content.WriteString("\n\n")
	
	content.WriteString(a.styleAction("Press F3 to generate your prompt"))
	content.WriteString("\n")
	content.WriteString(a.renderNavigationHelp())
	
	return content.String()
}

// Dialog renderers

func (a *AppState) renderHelpDialog() string {
	if a.HelpContent == "" {
		a.HelpContent = a.getHelpContent()
	}
	
	// Create help dialog box
	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")). // Blue
		Padding(1, 2).
		Width(60).
		Render(a.HelpContent)
	
	// Center the dialog
	return lipgloss.Place(
		a.WindowSize.Width,
		a.WindowSize.Height,
		lipgloss.Center,
		lipgloss.Center,
		helpBox,
	)
}

func (a *AppState) renderExitDialog() string {
	exitText := "Are you sure you want to exit?\n\nPress 'y' to quit or 'n' to continue"
	
	exitBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")). // Red
		Padding(1, 2).
		Width(40).
		Render(exitText)
	
	// Center the dialog
	return lipgloss.Place(
		a.WindowSize.Width,
		a.WindowSize.Height,
		lipgloss.Center,
		lipgloss.Center,
		exitBox,
	)
}

// Helper renderers

func (a *AppState) renderError() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")). // Red
		Bold(true).
		Render("⚠ " + a.Error.Error())
}

func (a *AppState) renderNavigationHelp() string {
	help := "F1: Help • F2: Previous • F3: Next • ESC: Exit"
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Gray
		Render(help)
}

// Styling functions

func (a *AppState) styleScreenTitle(title string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14")). // Cyan
		Render(title)
}

func (a *AppState) styleSelectedItem(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")). // Green
		Render("✓ " + text)
}

func (a *AppState) styleCursorItem(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")). // Blue
		Render(text)
}

func (a *AppState) styleDescription(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Gray
		Render(text)
}

func (a *AppState) stylePlaceholder(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Gray
		Italic(true).
		Render(text)
}

func (a *AppState) styleCursor(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")). // Blue
		Bold(true).
		Render(text)
}

func (a *AppState) styleInputBox(text string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")). // Gray
		Padding(1, 2).
		Width(a.WindowSize.Width - 4).
		Render(text)
}

func (a *AppState) styleConfirmationBox(text string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10")). // Green
		Padding(1, 2).
		Width(a.WindowSize.Width - 4).
		Render(text)
}

func (a *AppState) styleAction(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")). // Green
		Render(text)
}