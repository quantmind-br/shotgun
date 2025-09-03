package template

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Screen layout styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Align(lipgloss.Center).
			MarginTop(5)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("52")).
			Padding(1).
			MarginTop(2).
			MarginBottom(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center).
			MarginTop(1)

	mainStyle = lipgloss.NewStyle().
			Padding(1)
)

// View renders the template selection screen
func (m TemplateModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Loading..."
	}

	// Handle error state
	if m.err != nil {
		return m.renderError()
	}

	// Handle loading state
	if m.loading {
		return m.renderLoading()
	}

	// Handle empty state
	if len(m.templates) == 0 {
		return m.renderEmpty()
	}

	// Render normal screen
	return m.renderMain()
}

// renderError renders the error state
func (m TemplateModel) renderError() string {
	var content []string

	header := headerStyle.
		Width(m.width).
		Render("Template Selection - Error")

	content = append(content, header)

	errorMsg := errorStyle.
		Width(m.width - 4).
		Render("Error loading templates: " + m.err.Error())

	content = append(content, errorMsg)

	help := helpStyle.
		Width(m.width).
		Render("Press F2 to go back • Ctrl+R to retry • Ctrl+C to quit")

	content = append(content, help)

	return mainStyle.
		Width(m.width).
		Height(m.height).
		Render(strings.Join(content, "\n"))
}

// renderLoading renders the loading state
func (m TemplateModel) renderLoading() string {
	var content []string

	header := headerStyle.
		Width(m.width).
		Render("Template Selection")

	content = append(content, header)

	loading := loadingStyle.
		Width(m.width).
		Render("Loading templates...")

	content = append(content, loading)

	help := helpStyle.
		Width(m.width).
		Render("Press F2 to go back • Ctrl+C to quit")

	content = append(content, help)

	return mainStyle.
		Width(m.width).
		Height(m.height).
		Render(strings.Join(content, "\n"))
}

// renderEmpty renders the empty state (no templates found)
func (m TemplateModel) renderEmpty() string {
	var content []string

	header := headerStyle.
		Width(m.width).
		Render("Template Selection")

	content = append(content, header)

	emptyMsg := loadingStyle.
		Width(m.width).
		Render("No templates found")

	content = append(content, emptyMsg)

	help := helpStyle.
		Width(m.width).
		Render("Press F2 to go back • Ctrl+R to refresh • Ctrl+C to quit")

	content = append(content, help)

	return mainStyle.
		Width(m.width).
		Height(m.height).
		Render(strings.Join(content, "\n"))
}

// renderMain renders the main template selection view
func (m TemplateModel) renderMain() string {
	var content []string

	// Header
	header := headerStyle.
		Width(m.width).
		Render("Template Selection")

	content = append(content, header)

	// Main content area
	if m.showDetails && m.width > 80 {
		// Side-by-side layout: list + detail panel
		content = append(content, m.renderSideBySideLayout())
	} else {
		// Single column layout: just the list
		content = append(content, m.renderSingleColumnLayout())
	}

	// Help text
	helpText := m.buildHelpText()
	help := helpStyle.
		Width(m.width).
		Render(helpText)

	content = append(content, help)

	return mainStyle.
		Width(m.width).
		Height(m.height).
		Render(strings.Join(content, "\n"))
}

// renderSideBySideLayout renders list and detail panel side by side
func (m TemplateModel) renderSideBySideLayout() string {
	listWidth := m.width * 2 / 3
	detailWidth := m.width - listWidth - 2 // 2 for spacing

	listView := m.list.View()
	detailView := RenderDetailPanel(m.selected, detailWidth, m.height-6)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		detailView,
	)
}

// renderSingleColumnLayout renders just the list
func (m TemplateModel) renderSingleColumnLayout() string {
	return m.list.View()
}

// buildHelpText builds the help text based on current state
func (m TemplateModel) buildHelpText() string {
	var helpParts []string

	// Navigation
	helpParts = append(helpParts, "↑↓/jk: navigate")

	// Selection
	if m.selected != nil {
		helpParts = append(helpParts, "Enter/F3: select")
	}

	// Back
	helpParts = append(helpParts, "F2: back")

	// Refresh
	helpParts = append(helpParts, "Ctrl+R: refresh")

	// Detail panel toggle (only show if wide enough)
	if m.width > 80 {
		helpParts = append(helpParts, "Tab: toggle details")
	}

	// Quit
	helpParts = append(helpParts, "Ctrl+C: quit")

	return strings.Join(helpParts, " • ")
}
