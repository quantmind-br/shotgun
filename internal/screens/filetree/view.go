package filetree

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopedro/shotgun/internal/models"
)

var (
	// Styles for the file tree
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true)

	binaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Strikethrough(true)

	directoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("248")).
			Italic(true).
			Padding(0, 1)
)

// View renders the file tree screen
func (m FileTreeModel) View() string {
	// Show loading state with spinner during scanning
	if m.scanning {
		return m.renderScanningState()
	}

	// Show error state if scan failed
	if m.scanError != nil {
		return m.renderErrorState()
	}

	if len(m.items) == 0 {
		return "Loading file tree...\n" + m.statusBar()
	}

	var content strings.Builder
	flatItems := m.flattenTree(m.items, 0)

	for i, item := range flatItems {
		line := m.renderTreeItem(item, i == m.cursor)
		content.WriteString(line + "\n")
	}

	m.viewport.SetContent(content.String())

	return m.viewport.View() + "\n" + m.statusBar() + "\n" + m.helpBar()
}

// renderScanningState shows the spinner and progress during file scanning
func (m FileTreeModel) renderScanningState() string {
	var message string
	if m.filesFound > 0 {
		message = fmt.Sprintf("Found %d files", m.filesFound)
		if m.currentDir != "" {
			message += fmt.Sprintf(" (scanning %s/)", m.currentDir)
		}
	} else {
		message = "Scanning project files..."
	}

	m.spinner.SetMessage(message)
	spinnerView := m.spinner.ViewWithCancel()

	return spinnerView + "\n\n" + m.renderScanningStatusBar()
}

// renderErrorState shows scan error information
func (m FileTreeModel) renderErrorState() string {
	errorMsg := fmt.Sprintf("❌ Scan failed: %s", m.scanError.Error())
	retryHint := "\nPress 'r' to retry scanning or 'q' to quit"

	return errorMsg + retryHint
}

// renderScanningStatusBar shows status during scanning
func (m FileTreeModel) renderScanningStatusBar() string {
	var status string
	if m.filesFound > 0 {
		status = fmt.Sprintf("📄 %d files discovered", m.filesFound)
	} else {
		status = "🔍 Discovering files..."
	}

	// Add padding to fill width if needed
	if m.width > 0 {
		padding := m.width - len(status)
		if padding > 0 {
			status += strings.Repeat(" ", padding)
		}
	}

	return statusStyle.Render(status)
}

// flattenTree converts the tree structure to a flat list for rendering
func (m FileTreeModel) flattenTree(items []*models.FileNode, depth int) []treeItem {
	return m.flattenTreeInternal(items, depth)
}

// flattenTreeInternal is now defined in update.go to avoid duplication

// renderTreeItem renders a single tree item
func (m FileTreeModel) renderTreeItem(item treeItem, isSelected bool) string {
	// Create tree structure with proper indentation
	var treeStructure strings.Builder

	// Add indentation for nested levels
	for i := 0; i < item.depth; i++ {
		treeStructure.WriteString("  ")
	}

	// Add expand/collapse indicator for directories
	var expandIndicator string
	if item.node.IsDirectory {
		if len(item.node.Children) > 0 {
			if item.node.IsExpanded {
				expandIndicator = "▼ " // Expanded directory
			} else {
				expandIndicator = "▶ " // Collapsed directory
			}
		} else {
			expandIndicator = "  " // Empty directory
		}
	} else {
		expandIndicator = "  " // Files get space alignment
	}

	// Checkbox with different states
	var checkbox string
	if item.node.IsBinary {
		checkbox = "⚫ " // Unselectable binary file
	} else if item.node.IsSelected {
		checkbox = "✅ "
	} else {
		checkbox = "⬜ "
	}

	// File/directory icon
	var icon string
	if item.node.IsDirectory {
		icon = "📁 "
	} else {
		icon = "📄 "
	}

	// Apply styling to file name
	name := item.node.Name
	if item.node.IsDirectory {
		name = directoryStyle.Render(name)
	} else if item.node.IsBinary {
		name = binaryStyle.Render(name)
	}

	// Combine all parts
	line := treeStructure.String() + expandIndicator + checkbox + icon + name

	// Highlight current cursor position
	if isSelected {
		line = selectedStyle.Render(line)
	}

	return line
}

// statusBar renders the status bar with file counts
func (m FileTreeModel) statusBar() string {
	selected, excluded, ignored := m.calculateCounts()
	total := selected + excluded

	// Create individual sections with icons
	selectedText := fmt.Sprintf("✅ %d selected", selected)
	excludedText := fmt.Sprintf("⚫ %d excluded", excluded)
	ignoredText := fmt.Sprintf("🚫 %d ignored", ignored)
	totalText := fmt.Sprintf("📄 %d total", total)

	// Combine sections with separators
	status := fmt.Sprintf("%s  │  %s  │  %s  │  %s",
		selectedText, excludedText, ignoredText, totalText)

	// Add padding to fill width if needed
	if m.width > 0 {
		padding := m.width - len(status)
		if padding > 0 {
			status += strings.Repeat(" ", padding)
		}
	}

	return statusStyle.Render(status)
}

// helpBar renders the help/navigation bar
func (m FileTreeModel) helpBar() string {
	var help string
	if m.scanning {
		help = "ESC: cancel scanning │ q: quit"
	} else {
		help = "↑/↓ or k/j: navigate │ ←/→ or h/l: expand/collapse │ space: toggle │ F3: continue │ q: quit"
	}
	return helpStyle.Render(help)
}

// calculateCounts calculates file count statistics
func (m FileTreeModel) calculateCounts() (selected, excluded, ignored int) {
	m.countNodes(m.items, &selected, &excluded, &ignored)
	return
}

// countNodes recursively counts nodes by type
func (m FileTreeModel) countNodes(nodes []*models.FileNode, selected, excluded, ignored *int) {
	for _, node := range nodes {
		if node.IsDirectory {
			// Recursively count children for directories
			m.countNodes(node.Children, selected, excluded, ignored)
		} else {
			// Count files only
			if node.IsIgnored {
				*ignored++
			} else if node.IsBinary {
				*excluded++
			} else if node.IsSelected {
				*selected++
			} else {
				*excluded++
			}
		}
	}
}
