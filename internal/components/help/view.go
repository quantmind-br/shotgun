package help

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Help overlay styles
	helpOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2).
				Margin(1, 2)

	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center)

	helpContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	helpKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
)

// View renders the help overlay
func (h HelpModel) View() string {
	if !h.visible {
		return ""
	}

	// Create title based on current screen
	title := "Help - " + h.getScreenName()

	// Create the help content with proper styling
	content := h.viewport.View()

	// Style the content
	styledContent := helpContentStyle.Render(content)

	// Create the complete help overlay
	helpBox := lipgloss.JoinVertical(
		lipgloss.Center,
		helpTitleStyle.Render(title),
		"",
		styledContent,
	)

	// Apply the overlay styling
	overlay := helpOverlayStyle.
		Width(min(h.width-8, 60)).
		Height(min(h.height-6, 20)).
		Render(helpBox)

	// Center the overlay on screen
	return lipgloss.Place(
		h.width,
		h.height,
		lipgloss.Center,
		lipgloss.Center,
		overlay,
	)
}

// getScreenName returns a user-friendly screen name
func (h HelpModel) getScreenName() string {
	switch h.currentScreen {
	case FileTreeScreen:
		return "File Selection"
	case TemplateScreen:
		return "Template Selection"
	case TaskScreen:
		return "Task Input"
	case RulesScreen:
		return "Rules Input"
	case ConfirmScreen:
		return "Confirmation"
	case GenerateScreen:
		return "Generation"
	default:
		return "Unknown"
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
