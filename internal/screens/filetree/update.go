package filetree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
)

// treeItem represents a flattened tree item for navigation
type treeItem struct {
	node  *models.FileNode
	depth int
}

// Update handles messages and updates the model state
func (m FileTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2 // Reserve space for status bar

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case ScanCompleteMsg:
		m.LoadFileTree(msg.Nodes)
		return m, nil

	case ScanErrorMsg:
		// Handle scan error (could add error display to UI)
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// handleKeyPress processes keyboard input
func (m FileTreeModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.moveCursorUp()
	case "down", "j":
		m.moveCursorDown()
	case "left", "h":
		m.collapseDirectory()
	case "right", "l":
		m.expandDirectory()
	case " ":
		m.toggleSelection()
	}

	m.updateViewport()
	return m, nil
}

// moveCursorUp moves cursor up one position
func (m *FileTreeModel) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// moveCursorDown moves cursor down one position
func (m *FileTreeModel) moveCursorDown() {
	flatItems := m.getFlattenedItems()
	if len(flatItems) > 0 && m.cursor < len(flatItems)-1 {
		m.cursor++
	}
}

// getFlattenedItems returns the current flattened tree items
func (m *FileTreeModel) getFlattenedItems() []treeItem {
	return m.flattenTreeInternal(m.items, 0)
}

// flattenTreeInternal converts tree to flat list respecting expand/collapse state
func (m *FileTreeModel) flattenTreeInternal(items []*models.FileNode, depth int) []treeItem {
	var result []treeItem

	for _, item := range items {
		result = append(result, treeItem{
			node:  item,
			depth: depth,
		})

		// Only show children if directory is expanded
		if item.IsDirectory && item.IsExpanded && len(item.Children) > 0 {
			childItems := m.flattenTreeInternal(item.Children, depth+1)
			result = append(result, childItems...)
		}
	}

	return result
}

// collapseDirectory collapses the current directory
func (m *FileTreeModel) collapseDirectory() {
	flatItems := m.getFlattenedItems()
	if len(flatItems) == 0 || m.cursor >= len(flatItems) {
		return
	}

	currentItem := flatItems[m.cursor]
	if currentItem.node.IsDirectory {
		currentItem.node.IsExpanded = false
	}
}

// expandDirectory expands the current directory
func (m *FileTreeModel) expandDirectory() {
	flatItems := m.getFlattenedItems()
	if len(flatItems) == 0 || m.cursor >= len(flatItems) {
		return
	}

	currentItem := flatItems[m.cursor]
	if currentItem.node.IsDirectory && len(currentItem.node.Children) > 0 {
		currentItem.node.IsExpanded = true
	}
}

// toggleSelection toggles selection for current item
func (m *FileTreeModel) toggleSelection() {
	flatItems := m.getFlattenedItems()
	if len(flatItems) == 0 || m.cursor >= len(flatItems) {
		return
	}

	currentItem := flatItems[m.cursor].node

	// Don't allow toggling binary files
	if currentItem.IsBinary {
		return
	}

	currentItem.IsSelected = !currentItem.IsSelected
	m.selected[currentItem.Path] = currentItem.IsSelected

	// Handle hierarchical selection for directories
	if currentItem.IsDirectory {
		if currentItem.IsSelected {
			// When selecting a directory, select all non-binary children
			m.selectChildren(currentItem, true)
		} else {
			// When deselecting a directory, deselect all children
			m.deselectChildren(currentItem)
		}
	}

	// Update parent selection state based on children
	if currentItem.Parent != nil {
		m.updateParentSelection(currentItem.Parent)
	}
}

// selectChildren recursively selects/deselects all non-binary children
func (m *FileTreeModel) selectChildren(node *models.FileNode, selected bool) {
	for _, child := range node.Children {
		if !child.IsBinary { // Don't select binary files
			child.IsSelected = selected
			m.selected[child.Path] = selected
		}
		if child.IsDirectory {
			m.selectChildren(child, selected)
		}
	}
}

// deselectChildren recursively deselects all children
func (m *FileTreeModel) deselectChildren(node *models.FileNode) {
	m.selectChildren(node, false)
}

// updateParentSelection updates parent selection based on children state
func (m *FileTreeModel) updateParentSelection(parent *models.FileNode) {
	if parent == nil || !parent.IsDirectory {
		return
	}

	// Count selected and total selectable children
	selectedChildren := 0
	selectableChildren := 0

	for _, child := range parent.Children {
		if !child.IsBinary { // Only count non-binary files
			selectableChildren++
			if child.IsSelected {
				selectedChildren++
			}
		}
	}

	// Update parent selection based on children
	if selectableChildren > 0 {
		// Parent is selected if all selectable children are selected
		parent.IsSelected = selectedChildren == selectableChildren
		m.selected[parent.Path] = parent.IsSelected
	}

	// Recursively update grandparent
	if parent.Parent != nil {
		m.updateParentSelection(parent.Parent)
	}
}

// updateViewport updates viewport content based on cursor position
func (m *FileTreeModel) updateViewport() {
	if len(m.items) == 0 {
		return
	}

	// Ensure cursor is visible in viewport
	if m.cursor < m.viewport.YOffset {
		m.viewport.YOffset = m.cursor
	} else if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.YOffset = m.cursor - m.viewport.Height + 1
	}
}

// GetSelectedFiles returns a slice of paths for all selected files
func (m *FileTreeModel) GetSelectedFiles() []string {
	var selectedFiles []string
	m.collectSelectedFiles(m.items, &selectedFiles)
	return selectedFiles
}

// collectSelectedFiles recursively collects paths of selected files
func (m *FileTreeModel) collectSelectedFiles(nodes []*models.FileNode, selected *[]string) {
	for _, node := range nodes {
		if !node.IsDirectory && node.IsSelected && !node.IsBinary {
			*selected = append(*selected, node.Path)
		}
		if node.IsDirectory && len(node.Children) > 0 {
			m.collectSelectedFiles(node.Children, selected)
		}
	}
}

// SetSize updates the width and height of the model and viewport
func (m *FileTreeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height - 2 // Reserve space for status bar
}
