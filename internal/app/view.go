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
	screenContent := a.renderCurrentScreen()
	content.WriteString(screenContent)

	// Render exit dialog if showing (takes priority over help)
	if a.ShowingExit {
		return a.renderExitDialog()
	}

	// Render help overlay on top of current screen if showing
	if a.Help.IsVisible() {
		baseView := content.String()
		helpOverlay := a.Help.View()
		// Layer the help overlay on top of the base view
		return baseView + helpOverlay
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
	case GenerateScreen:
		return a.Generation.View()
	default:
		return "Unknown screen"
	}
}

// Screen renderers

func (a *AppState) renderTemplateScreen() string {
	return a.Template.View()
}

func (a *AppState) renderTaskScreen() string {
	return a.TaskInput.View()
}

func (a *AppState) renderRulesScreen() string {
	return a.RulesInput.View()
}

func (a *AppState) renderConfirmationScreen() string {
	return a.Confirmation.View()
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
	help := "Ctrl+H: Help • Ctrl+Left: Previous • Alt+C: Next • Ctrl+Q/ESC: Exit"
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
