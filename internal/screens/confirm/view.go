package confirm

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Style definitions
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Margin(1, 0)

	summaryStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("36")).
			Padding(1, 2).
			Margin(0, 0, 1, 0)

	sizeStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("33")).
			Padding(1, 2).
			Margin(0, 0, 1, 0)

	warningStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2).
			Margin(0, 0, 1, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Margin(1, 0, 0, 0)

	// Warning level colors
	normalColor    = lipgloss.Color("82")  // Green
	largeColor     = lipgloss.Color("214") // Yellow  
	veryLargeColor = lipgloss.Color("208") // Orange
	excessiveColor = lipgloss.Color("196") // Red
)

// View renders the confirmation screen
func (m ConfirmModel) View() string {
	if !m.IsReady() {
		return "Loading confirmation screen..."
	}

	var sections []string

	// Title
	sections = append(sections, titleStyle.Render("Confirmation"))

	// Selection summary section
	sections = append(sections, m.renderSelectionSummary())

	// Size estimation section
	if m.calculating {
		sections = append(sections, m.renderSizeCalculation())
	} else {
		sections = append(sections, m.renderSizeEstimation())
	}

	// Warning section (if applicable)
	if m.showWarning {
		sections = append(sections, m.renderWarningSection())
	}

	// Navigation help
	sections = append(sections, m.renderNavigationHelp())

	content := strings.Join(sections, "\n")

	// Wrap in viewport if content is scrollable
	if len(content) > m.viewport.Height*80 { // Estimate line wrapping
		m.viewport.SetContent(content)
		return m.viewport.View()
	}

	return content
}

// renderSelectionSummary renders the template and file selection summary
func (m ConfirmModel) renderSelectionSummary() string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Template & Selection Summary"))
	content.WriteString("\n\n")

	// Template information
	templateInfo := fmt.Sprintf("Template: \"%s\" v%s", 
		m.template.Name, m.template.Version)
	if m.template.Description != "" {
		templateInfo += fmt.Sprintf("\nDescription: %s", m.template.Description)
	}
	content.WriteString(templateInfo)
	content.WriteString("\n\n")

	// File selection summary
	totalSelected := len(m.selectedFiles)
	filesSummary := fmt.Sprintf("Files: %d selected", totalSelected)
	
	// Add excluded file count if available (would need to be passed from AppState)
	// For now, just show selected count
	content.WriteString(filesSummary)
	content.WriteString("\n\n")

	// User inputs summary
	if m.taskContent != "" {
		content.WriteString("Task: ")
		taskPreview := m.taskContent
		if len(taskPreview) > 80 {
			taskPreview = taskPreview[:77] + "..."
		}
		content.WriteString(fmt.Sprintf("\"%s\"", taskPreview))
		content.WriteString("\n")
	}

	if m.rulesContent != "" {
		content.WriteString("Rules: ")
		rulesPreview := m.rulesContent
		if len(rulesPreview) > 80 {
			rulesPreview = rulesPreview[:77] + "..."
		}
		content.WriteString(fmt.Sprintf("\"%s\"", rulesPreview))
		content.WriteString("\n")
	}

	// Show sample of selected files (first few)
	if len(m.selectedFiles) > 0 {
		content.WriteString("\nSelected Files (sample):\n")
		maxFiles := 5
		if len(m.selectedFiles) < maxFiles {
			maxFiles = len(m.selectedFiles)
		}
		
		for i := 0; i < maxFiles; i++ {
			fileName := filepath.Base(m.selectedFiles[i])
			content.WriteString(fmt.Sprintf("  • %s\n", fileName))
		}
		
		if len(m.selectedFiles) > maxFiles {
			content.WriteString(fmt.Sprintf("  ... and %d more files\n", len(m.selectedFiles)-maxFiles))
		}
	}

	return summaryStyle.Render(content.String())
}

// renderSizeCalculation renders the progress bar during calculation
func (m ConfirmModel) renderSizeCalculation() string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Output Size Estimation"))
	content.WriteString("\n\n")

	content.WriteString("Calculating size estimation...\n\n")
	content.WriteString(m.progress.View())
	content.WriteString("\n")

	return sizeStyle.Render(content.String())
}

// renderSizeEstimation renders the final size estimation results
func (m ConfirmModel) renderSizeEstimation() string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Output Size Estimation"))
	content.WriteString("\n\n")

	// Size breakdown
	content.WriteString(fmt.Sprintf("Template + Variables: %s\n", formatBytes(m.sizeBreakdown.TemplateSize)))
	content.WriteString(fmt.Sprintf("File Contents: %s\n", formatBytes(m.sizeBreakdown.FileContentSize)))
	content.WriteString(fmt.Sprintf("Structure Overhead: %s\n", formatBytes(m.sizeBreakdown.TreeStructSize)))
	content.WriteString(fmt.Sprintf("Formatting Overhead: %s\n", formatBytes(m.sizeBreakdown.OverheadSize)))
	content.WriteString("\n")

	// Total size with color coding
	totalSizeStr := fmt.Sprintf("Total Estimated Size: %s", formatBytes(m.estimatedSize))
	switch m.warningLevel {
	case WarningLarge:
		totalSizeStr = lipgloss.NewStyle().Foreground(largeColor).Render(totalSizeStr)
	case WarningVeryLarge:
		totalSizeStr = lipgloss.NewStyle().Foreground(veryLargeColor).Render(totalSizeStr)
	case WarningExcessive:
		totalSizeStr = lipgloss.NewStyle().Foreground(excessiveColor).Render(totalSizeStr)
	default:
		totalSizeStr = lipgloss.NewStyle().Foreground(normalColor).Render(totalSizeStr)
	}
	content.WriteString(totalSizeStr)
	content.WriteString("\n\n")

	// Output filename
	if m.outputFilename != "" {
		content.WriteString(fmt.Sprintf("Output File: %s", m.outputFilename))
	}

	return sizeStyle.Render(content.String())
}

// renderWarningSection renders size warnings if applicable
func (m ConfirmModel) renderWarningSection() string {
	var content strings.Builder
	var warningIcon, warningText string

	switch m.warningLevel {
	case WarningLarge:
		warningIcon = "⚠️"
		warningText = "Large Output Notice"
	case WarningVeryLarge:
		warningIcon = "⚠️"
		warningText = "Large Output Warning"
	case WarningExcessive:
		warningIcon = "🚨"
		warningText = "Very Large Output Warning"
	}

	content.WriteString(fmt.Sprintf("%s  %s", warningIcon, 
		lipgloss.NewStyle().Bold(true).Render(warningText)))
	content.WriteString("\n\n")

	switch m.warningLevel {
	case WarningLarge:
		content.WriteString("This prompt will be moderately large. ")
		content.WriteString("Consider reviewing file selection if performance is a concern.")
	case WarningVeryLarge:
		content.WriteString(fmt.Sprintf("This prompt will be very large (%s).\n", formatBytes(m.estimatedSize)))
		content.WriteString("Consider excluding some files or content.\n")
		content.WriteString("Large prompts may impact LLM performance.")
	case WarningExcessive:
		content.WriteString(fmt.Sprintf("This prompt will be extremely large (%s)!\n", formatBytes(m.estimatedSize)))
		content.WriteString("This may exceed LLM token limits or cause performance issues.\n")
		content.WriteString("Consider significantly reducing file selection.")
	}

	return warningStyle.Render(content.String())
}

// renderNavigationHelp renders keyboard navigation instructions
func (m ConfirmModel) renderNavigationHelp() string {
	help := []string{
		"F10: Confirm and generate prompt",
		"F2: Return to rules input", 
		"F1: Return to file selection",
		"Esc: Exit application",
	}

	if m.viewport.TotalLineCount() > m.viewport.Height {
		help = append(help, "↑/↓: Scroll content")
	}

	return helpStyle.Render(strings.Join(help, " • "))
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}