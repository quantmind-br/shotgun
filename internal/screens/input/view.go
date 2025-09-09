package input

import (
	"fmt"
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

	textareaStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1)

	focusedTextareaStyle = textareaStyle.Copy().
				BorderForeground(lipgloss.Color("205"))

	countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Right).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("52")).
			Padding(1).
			MarginTop(1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Left).
			MarginTop(1)

	instructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99")).
				Italic(true).
				MarginBottom(1)
)

// View renders the task input screen
func (m TaskInputModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Loading..."
	}

	// Handle error state
	if m.err != nil {
		return m.renderError()
	}

	// Render normal screen
	return m.renderMain()
}

// renderMain renders the main task input screen
func (m TaskInputModel) renderMain() string {
	var sections []string

	// Header
	header := headerStyle.Width(m.width).Render("📝 Task Description")
	sections = append(sections, header)

	// Instructions
	instruction := instructionStyle.Width(m.width).Render(
		"Describe your task in detail. The LLM will use this context to understand what you need.",
	)
	sections = append(sections, instruction)

	// Text area with proper styling based on focus
	textareaContent := m.textarea.View()
	if m.textarea.Focused() {
		textareaContent = focusedTextareaStyle.Width(m.width - 4).Render(textareaContent)
	} else {
		textareaContent = textareaStyle.Width(m.width - 4).Render(textareaContent)
	}
	sections = append(sections, textareaContent)

	// Character and line counts
	countText := fmt.Sprintf("Lines: %d | Characters: %d", m.lineCount, m.charCount)
	countDisplay := countStyle.Width(m.width - 4).Render(countText)
	sections = append(sections, countDisplay)

	// Help text for keyboard shortcuts
	helpText := []string{
		"Alt+C: Continue (if content is not empty)",
		"Ctrl+Left: Back to Template Selection",
		"Ctrl+C: Copy selected text",
		"Ctrl+V: Paste from clipboard",
	}
	help := helpStyle.Width(m.width - 4).Render(strings.Join(helpText, " • "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderError renders the error state
func (m TaskInputModel) renderError() string {
	var sections []string

	// Header
	header := headerStyle.Width(m.width).Render("📝 Task Description")
	sections = append(sections, header)

	// Error message
	errorMsg := fmt.Sprintf("Error: %s", m.err.Error())
	errorDisplay := errorStyle.Width(m.width - 4).Render(errorMsg)
	sections = append(sections, errorDisplay)

	// Instructions
	instruction := instructionStyle.Width(m.width).Render(
		"Please provide a task description to continue.",
	)
	sections = append(sections, instruction)

	// Text area with error styling
	textareaContent := m.textarea.View()
	textareaContent = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")). // Red border for error state
		Padding(1).
		Width(m.width - 4).
		Render(textareaContent)
	sections = append(sections, textareaContent)

	// Character and line counts
	countText := fmt.Sprintf("Lines: %d | Characters: %d", m.lineCount, m.charCount)
	countDisplay := countStyle.Width(m.width - 4).Render(countText)
	sections = append(sections, countDisplay)

	// Help text
	helpText := []string{
		"Alt+C: Continue (if content is not empty)",
		"Ctrl+Left: Back to Template Selection",
		"Ctrl+C: Copy selected text",
		"Ctrl+V: Paste from clipboard",
	}
	help := helpStyle.Width(m.width - 4).Render(strings.Join(helpText, " • "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
