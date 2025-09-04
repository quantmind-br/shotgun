package template

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/diogopedro/shotgun/internal/models"
)

var (
	// Styles for template list items
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			MarginRight(1)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(1).
				MarginRight(1).
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99"))

	detailPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1).
				MarginLeft(1)

	detailTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				Underline(true)

	detailFieldStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	detailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("159"))

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("227")).
			Background(lipgloss.Color("57")).
			Padding(0, 1).
			MarginRight(1)
)

// TemplateDelegate implements list.ItemDelegate for custom template rendering
type TemplateDelegate struct {
	selected bool
}

// Height returns the height of each list item
func (d TemplateDelegate) Height() int {
	return 3
}

// Spacing returns the spacing between list items
func (d TemplateDelegate) Spacing() int {
	return 1
}

// Update handles delegate updates
func (d TemplateDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render renders a template list item
func (d TemplateDelegate) Render(w int, m list.Model, index int, item list.Item) string {
	templateItem, ok := item.(TemplateItem)
	if !ok {
		return ""
	}

	template := templateItem.Template
	isSelected := index == m.Index()

	// Build the item display
	var lines []string

	// Line 1: Name and Version
	nameVersion := template.Name
	if template.Version != "" {
		nameVersion += " " + versionStyle.Render("v"+template.Version)
	}

	if isSelected {
		lines = append(lines, selectedItemStyle.Render(titleStyle.Render(nameVersion)))
	} else {
		lines = append(lines, itemStyle.Render(titleStyle.Render(nameVersion)))
	}

	// Line 2: Description
	description := template.Description
	if len(description) > w-4 {
		description = description[:w-7] + "..."
	}

	if isSelected {
		lines = append(lines, selectedItemStyle.Render(descStyle.Render(description)))
	} else {
		lines = append(lines, itemStyle.Render(descStyle.Render(description)))
	}

	// Line 3: Author and tags preview
	authorInfo := ""
	if template.Author != "" {
		authorInfo = "by " + template.Author
	}

	tagPreview := ""
	if len(template.Tags) > 0 {
		// Show first 2 tags
		tagsToShow := template.Tags
		if len(tagsToShow) > 2 {
			tagsToShow = tagsToShow[:2]
		}
		tagPreview = " [" + strings.Join(tagsToShow, ", ")
		if len(template.Tags) > 2 {
			tagPreview += fmt.Sprintf(", +%d more", len(template.Tags)-2)
		}
		tagPreview += "]"
	}

	metaInfo := authorInfo + tagPreview
	if len(metaInfo) > w-4 {
		metaInfo = metaInfo[:w-7] + "..."
	}

	if isSelected {
		lines = append(lines, selectedItemStyle.Render(descStyle.Render(metaInfo)))
	} else {
		lines = append(lines, itemStyle.Render(descStyle.Render(metaInfo)))
	}

	return strings.Join(lines, "\n")
}

// RenderDetailPanel renders the detail panel for the selected template
func RenderDetailPanel(template *models.Template, width, height int) string {
	if template == nil {
		return detailPanelStyle.
			Width(width - 2).
			Height(height - 2).
			Render("No template selected")
	}

	var details []string

	// Title
	details = append(details, detailTitleStyle.Render(template.Name))

	if template.Version != "" {
		details = append(details, detailFieldStyle.Render("Version: ")+detailValueStyle.Render(template.Version))
	}

	// Description
	if template.Description != "" {
		details = append(details, "")
		details = append(details, detailFieldStyle.Render("Description:"))

		// Wrap description to fit panel width
		descLines := wrapText(template.Description, width-6)
		for _, line := range descLines {
			details = append(details, detailValueStyle.Render("  "+line))
		}
	}

	// Author
	if template.Author != "" {
		details = append(details, "")
		details = append(details, detailFieldStyle.Render("Author: ")+detailValueStyle.Render(template.Author))
	}

	// Tags
	if len(template.Tags) > 0 {
		details = append(details, "")
		details = append(details, detailFieldStyle.Render("Tags:"))

		// Render tags with styling
		var tagLines []string
		currentLine := "  "
		for i, tag := range template.Tags {
			tagRendered := tagStyle.Render(tag)
			if len(currentLine+tagRendered) > width-6 {
				if currentLine != "  " {
					tagLines = append(tagLines, currentLine)
					currentLine = "  "
				}
			}
			currentLine += tagRendered
			if i < len(template.Tags)-1 {
				currentLine += " "
			}
		}
		if currentLine != "  " {
			tagLines = append(tagLines, currentLine)
		}

		details = append(details, tagLines...)
	}

	// Variables
	if len(template.Variables) > 0 {
		details = append(details, "")
		details = append(details, detailFieldStyle.Render("Variables:"))

		for name, variable := range template.Variables {
			varLine := fmt.Sprintf("  %s (%s)", name, variable.Type)
			if variable.Placeholder != "" {
				varLine += " - " + variable.Placeholder
			}

			// Wrap if too long
			if len(varLine) > width-6 {
				varLine = varLine[:width-9] + "..."
			}

			details = append(details, detailValueStyle.Render(varLine))
		}
	}

	content := strings.Join(details, "\n")

	return detailPanelStyle.
		Width(width - 2).
		Height(height - 2).
		Render(content)
}

// wrapText wraps text to specified width
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine+" "+word) > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Word is longer than width, split it
				for len(word) > width {
					lines = append(lines, word[:width])
					word = word[width:]
				}
				currentLine = word
			}
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
