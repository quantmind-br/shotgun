package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styling for the generation screen
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F56")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFBD2E"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(60)

	progressBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(0, 1).
				Width(50)
)

// View renders the generation screen
func (m GenerateModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content strings.Builder

	// Title
	title := titleStyle.Render("Generating Prompt")
	content.WriteString(title)
	content.WriteString("\n\n")

	if m.generating {
		content.WriteString(m.renderGeneratingView())
	} else if m.completed {
		if m.HasError() {
			content.WriteString(m.renderErrorView())
		} else {
			content.WriteString(m.renderSuccessView())
		}
	} else {
		content.WriteString(m.renderIdleView())
	}

	// Add footer with key bindings
	content.WriteString("\n\n")
	content.WriteString(m.renderFooter())

	return content.String()
}

// renderGeneratingView renders the view while generation is in progress
func (m GenerateModel) renderGeneratingView() string {
	var content strings.Builder

	// Spinner and current status
	content.WriteString(m.spinner.View())
	content.WriteString(" Processing template...\n\n")

	// Progress bar
	progressBox := progressBoxStyle.Render(m.progress.View())
	content.WriteString(progressBox)
	content.WriteString("\n\n")

	// Generation info
	if m.templateSize > 0 {
		content.WriteString(fmt.Sprintf("Template size: %s\n", formatBytes(m.templateSize)))
	}
	if m.fileCount > 0 {
		content.WriteString(fmt.Sprintf("Files selected: %d\n", m.fileCount))
	}
	if m.totalSize > 0 {
		content.WriteString(fmt.Sprintf("Estimated size: %s\n", formatBytes(m.totalSize)))
	}

	return boxStyle.Render(content.String())
}

// renderSuccessView renders the view when generation completes successfully
func (m GenerateModel) renderSuccessView() string {
	var content strings.Builder

	// Success message
	content.WriteString(successStyle.Render("✓ Prompt successfully generated!"))
	content.WriteString("\n\n")

	// File information
	if m.outputFile != "" {
		fileName := filepath.Base(m.outputFile)
		content.WriteString(fmt.Sprintf("File: %s\n", infoStyle.Render(fileName)))

		// Show truncated path if it's long
		dir := filepath.Dir(m.outputFile)
		if len(dir) > 40 {
			dir = "..." + dir[len(dir)-37:]
		}
		content.WriteString(fmt.Sprintf("Path: %s\n", infoStyle.Render(dir)))

		if m.generatedSize > 0 {
			content.WriteString(fmt.Sprintf("Size: %s", infoStyle.Render(formatBytes(m.generatedSize))))
			if m.totalSize > 0 && m.totalSize != m.generatedSize {
				content.WriteString(fmt.Sprintf(" (estimate was %s)", formatBytes(m.totalSize)))
			}
			content.WriteString("\n")
		}
	}

	// Generation statistics
	if m.showStats && (m.templateSize > 0 || m.fileCount > 0) {
		content.WriteString("\n")
		if m.templateSize > 0 {
			content.WriteString(fmt.Sprintf("Template processed: %s\n", formatBytes(m.templateSize)))
		}
		if m.fileCount > 0 {
			fileSize := m.generatedSize - m.templateSize
			if fileSize > 0 {
				content.WriteString(fmt.Sprintf("Files included: %d files, %s\n", m.fileCount, formatBytes(fileSize)))
			} else {
				content.WriteString(fmt.Sprintf("Files included: %d files\n", m.fileCount))
			}
		}
	}

	return boxStyle.Render(content.String())
}

// renderErrorView renders the view when generation fails
func (m GenerateModel) renderErrorView() string {
	var content strings.Builder

	// Error message
	content.WriteString(errorStyle.Render("✗ Failed to generate prompt"))
	content.WriteString("\n\n")

	// Error details
	if m.error != nil {
		errorMsg := m.error.Error()
		if len(errorMsg) > 80 {
			errorMsg = errorMsg[:77] + "..."
		}
		content.WriteString(fmt.Sprintf("Error: %s\n", errorStyle.Render(errorMsg)))
		content.WriteString("\n")
	}

	// Suggestions
	content.WriteString("Suggestions:\n")
	content.WriteString("• Check file permissions and disk space\n")
	content.WriteString("• Verify template and file selections\n")
	content.WriteString("• Try reducing the number of selected files\n")

	return boxStyle.Render(content.String())
}

// renderIdleView renders the view when not generating
func (m GenerateModel) renderIdleView() string {
	return boxStyle.Render("Ready to generate prompt...")
}

// renderFooter renders the key binding footer
func (m GenerateModel) renderFooter() string {
	var keys []string

	if m.generating {
		keys = append(keys, "Esc: Cancel")
	} else if m.completed {
		if m.HasError() {
			keys = append(keys, "F5: Retry", "F1: Start over", "Esc: Exit")
		} else {
			if m.outputFile != "" {
				keys = append(keys, "F2: Open file")
			}
			keys = append(keys, "F1: Start over", "Esc: Exit")
		}
		if m.showStats {
			keys = append(keys, "s: Hide stats")
		} else {
			keys = append(keys, "s: Show stats")
		}
	} else {
		keys = append(keys, "F1: Start over", "Esc: Exit")
	}

	return infoStyle.Render(strings.Join(keys, "  "))
}

// formatBytes converts bytes to human readable format
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

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
