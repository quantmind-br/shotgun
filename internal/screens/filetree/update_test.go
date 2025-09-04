package filetree

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
)

func TestUpdate_WindowSizeMsg(t *testing.T) {
	model := NewFileTreeModel()

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(FileTreeModel)

	if m.width != 100 {
		t.Errorf("Expected width = 100, got %d", m.width)
	}

	if m.height != 30 {
		t.Errorf("Expected height = 30, got %d", m.height)
	}

	if m.viewport.Width != 100 {
		t.Errorf("Expected viewport width = 100, got %d", m.viewport.Width)
	}

	// Height should be reduced by 2 for status bar
	if m.viewport.Height != 28 {
		t.Errorf("Expected viewport height = 28, got %d", m.viewport.Height)
	}
}

func TestUpdate_ScanCompleteMsg(t *testing.T) {
	model := NewFileTreeModel()

	nodes := []*models.FileNode{
		{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  false,
		},
	}

	msg := ScanCompleteMsg{Nodes: nodes}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(FileTreeModel)

	if len(m.items) != 1 {
		t.Errorf("Expected 1 item after ScanCompleteMsg, got %d", len(m.items))
	}

	if m.items[0].Name != "file1.txt" {
		t.Errorf("Expected file1.txt, got %s", m.items[0].Name)
	}
}

func TestUpdate_ScanErrorMsg(t *testing.T) {
	model := NewFileTreeModel()

	msg := ScanErrorMsg{Error: nil}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(FileTreeModel)

	// Should not crash and model should remain unchanged
	if len(m.items) != 0 {
		t.Errorf("Expected no items after ScanErrorMsg, got %d", len(m.items))
	}
}

func TestHandleKeyPress_Navigation(t *testing.T) {
	model := NewFileTreeModel()

	// Setup test data
	nodes := createTestTreeNodes()
	model.LoadFileTree(nodes)

	tests := []struct {
		name           string
		key            string
		expectedCursor int
	}{
		{"down arrow", "down", 1},
		{"up arrow from position 1", "up", 0},
		{"j key (vim down)", "j", 1},
		{"k key (vim up)", "k", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cursor for each test
			if tt.name == "up arrow from position 1" || tt.name == "k key (vim up)" {
				model.cursor = 1
			} else {
				model.cursor = 0
			}

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "up" {
				msg = tea.KeyMsg{Type: tea.KeyUp}
			} else if tt.key == "down" {
				msg = tea.KeyMsg{Type: tea.KeyDown}
			}

			updatedModel, _ := model.handleKeyPress(msg)
			m := updatedModel.(FileTreeModel)

			if m.cursor != tt.expectedCursor {
				t.Errorf("Expected cursor = %d, got %d", tt.expectedCursor, m.cursor)
			}
		})
	}
}

func TestHandleKeyPress_ExpandCollapse(t *testing.T) {
	model := NewFileTreeModel()

	// Create a directory node
	nodes := []*models.FileNode{
		{
			Path:        "/test/dir1",
			Name:        "dir1",
			IsDirectory: true,
			IsExpanded:  true,
			Children: []*models.FileNode{
				{
					Path:        "/test/dir1/file1.txt",
					Name:        "file1.txt",
					IsDirectory: false,
				},
			},
		},
	}

	model.LoadFileTree(nodes)
	model.cursor = 0

	// Test collapse
	leftKey := tea.KeyMsg{Type: tea.KeyLeft}
	updatedModel, _ := model.handleKeyPress(leftKey)
	m := updatedModel.(FileTreeModel)

	flatItems := m.getFlattenedItems()
	if len(flatItems) > 0 && flatItems[0].node.IsExpanded {
		t.Error("Expected directory to be collapsed after left arrow")
	}

	// Test expand
	rightKey := tea.KeyMsg{Type: tea.KeyRight}
	updatedModel, _ = m.handleKeyPress(rightKey)
	m = updatedModel.(FileTreeModel)

	flatItems = m.getFlattenedItems()
	if len(flatItems) > 0 && !flatItems[0].node.IsExpanded {
		t.Error("Expected directory to be expanded after right arrow")
	}
}

func TestHandleKeyPress_ToggleSelection(t *testing.T) {
	model := NewFileTreeModel()

	// Create test nodes
	nodes := []*models.FileNode{
		{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  true,
			IsBinary:    false,
		},
		{
			Path:        "/test/binary.exe",
			Name:        "binary.exe",
			IsDirectory: false,
			IsSelected:  false,
			IsBinary:    true,
		},
	}

	model.LoadFileTree(nodes)

	// Test toggling regular file
	model.cursor = 0
	spaceKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	updatedModel, _ := model.handleKeyPress(spaceKey)
	m := updatedModel.(FileTreeModel)

	if m.items[0].IsSelected {
		t.Error("Expected file to be deselected after space key")
	}

	// Test toggling binary file (should not change)
	model.cursor = 1
	updatedModel, _ = model.handleKeyPress(spaceKey)
	m = updatedModel.(FileTreeModel)

	if m.items[1].IsSelected {
		t.Error("Expected binary file to remain deselected")
	}
}

func TestMoveCursor(t *testing.T) {
	model := NewFileTreeModel()

	// Setup test data
	nodes := createTestTreeNodes()
	model.LoadFileTree(nodes)

	// Test cursor bounds
	model.cursor = 0
	model.moveCursorUp()
	if model.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 when moving up from first position, got %d", model.cursor)
	}

	// Set cursor to last valid position
	flatItems := model.getFlattenedItems()
	model.cursor = len(flatItems) - 1
	model.moveCursorDown()
	if model.cursor != len(flatItems)-1 {
		t.Errorf("Expected cursor to stay at last position when moving down, got %d", model.cursor)
	}
}

func TestToggleSelection_Hierarchical(t *testing.T) {
	model := NewFileTreeModel()

	// Create directory with children
	childFile := &models.FileNode{
		Path:        "/test/dir1/file1.txt",
		Name:        "file1.txt",
		IsDirectory: false,
		IsSelected:  true,
		IsBinary:    false,
	}

	parentDir := &models.FileNode{
		Path:        "/test/dir1",
		Name:        "dir1",
		IsDirectory: true,
		IsSelected:  true,
		Children:    []*models.FileNode{childFile},
	}

	childFile.Parent = parentDir

	model.LoadFileTree([]*models.FileNode{parentDir})
	model.cursor = 0

	// Deselect parent directory
	model.toggleSelection()

	if parentDir.IsSelected {
		t.Error("Expected parent directory to be deselected")
	}

	if childFile.IsSelected {
		t.Error("Expected child file to be deselected when parent is deselected")
	}
}

func TestUpdateParentSelection(t *testing.T) {
	model := NewFileTreeModel()

	// Create a nested directory structure
	grandChild := &models.FileNode{
		Path:        "/test/parent/child/grandchild.txt",
		Name:        "grandchild.txt",
		IsDirectory: false,
		IsSelected:  true,
		IsBinary:    false,
	}

	child := &models.FileNode{
		Path:        "/test/parent/child",
		Name:        "child",
		IsDirectory: true,
		IsSelected:  true,
		Children:    []*models.FileNode{grandChild},
	}

	parent := &models.FileNode{
		Path:        "/test/parent",
		Name:        "parent",
		IsDirectory: true,
		IsSelected:  true,
		Children:    []*models.FileNode{child},
	}

	// Set up parent relationships
	grandChild.Parent = child
	child.Parent = parent

	model.LoadFileTree([]*models.FileNode{parent})

	// Deselect the grandchild
	grandChild.IsSelected = false

	// This should trigger parent selection update
	model.updateParentSelection(child)

	// Child should be deselected since all its children are deselected
	if child.IsSelected {
		t.Error("Expected child directory to be deselected when all children are deselected")
	}

	// Parent should be deselected since child is now deselected
	model.updateParentSelection(parent)
	if parent.IsSelected {
		t.Error("Expected parent directory to be deselected when child is deselected")
	}

	// Test with nil parent (edge case)
	model.updateParentSelection(nil)
	// Should not crash
}

func TestUpdateParentSelection_MixedChildren(t *testing.T) {
	model := NewFileTreeModel()

	// Create parent with mixed children (some selected, some not)
	child1 := &models.FileNode{
		Path:        "/test/parent/child1.txt",
		Name:        "child1.txt",
		IsDirectory: false,
		IsSelected:  true,
	}

	child2 := &models.FileNode{
		Path:        "/test/parent/child2.txt",
		Name:        "child2.txt",
		IsDirectory: false,
		IsSelected:  false,
	}

	parent := &models.FileNode{
		Path:        "/test/parent",
		Name:        "parent",
		IsDirectory: true,
		IsSelected:  true,
		Children:    []*models.FileNode{child1, child2},
	}

	child1.Parent = parent
	child2.Parent = parent

	// Don't use LoadFileTree as it auto-selects all files, instead manually set up
	model.items = []*models.FileNode{parent}

	// Manually set the desired selection states
	child1.IsSelected = true
	child2.IsSelected = false
	parent.IsSelected = true

	// Update parent selection with mixed children
	model.updateParentSelection(parent)

	// Parent should be deselected because not ALL children are selected
	if parent.IsSelected {
		t.Error("Expected parent to be deselected when children have mixed selection")
	}
}

func TestUpdateViewport(t *testing.T) {
	model := NewFileTreeModel()

	// Set up a model with some items and viewport dimensions
	model.height = 10
	model.viewport.Height = 8 // 2 lines reserved for status/help

	// Create enough items to require scrolling
	var nodes []*models.FileNode
	for i := 0; i < 20; i++ {
		nodes = append(nodes, &models.FileNode{
			Path:        fmt.Sprintf("/test/file%d.txt", i),
			Name:        fmt.Sprintf("file%d.txt", i),
			IsDirectory: false,
			IsSelected:  true,
		})
	}

	model.LoadFileTree(nodes)

	// Set cursor to bottom of viewport
	model.cursor = 15

	// Update viewport should adjust scroll position
	model.updateViewport()

	// Viewport should have content set
	content := model.viewport.View()
	if content == "" {
		t.Error("Expected viewport to have content after updateViewport")
	}
}

func TestExpandDirectory_EmptyDirectory(t *testing.T) {
	model := NewFileTreeModel()

	// Create an empty directory
	emptyDir := &models.FileNode{
		Path:        "/test/empty",
		Name:        "empty",
		IsDirectory: true,
		IsExpanded:  false,
		Children:    []*models.FileNode{}, // Empty directory
	}

	model.LoadFileTree([]*models.FileNode{emptyDir})
	model.cursor = 0

	// Try to expand empty directory
	model.expandDirectory()

	// Empty directory should remain collapsed (as per implementation - only expand if has children)
	if emptyDir.IsExpanded {
		t.Error("Expected empty directory to remain collapsed (no children to expand)")
	}
}

func TestCollapseDirectory_AlreadyCollapsed(t *testing.T) {
	model := NewFileTreeModel()

	// Create a collapsed directory
	collapsedDir := &models.FileNode{
		Path:        "/test/collapsed",
		Name:        "collapsed",
		IsDirectory: true,
		IsExpanded:  false, // Already collapsed
		Children: []*models.FileNode{
			{Name: "file.txt", IsDirectory: false},
		},
	}

	model.LoadFileTree([]*models.FileNode{collapsedDir})
	model.cursor = 0

	// Try to collapse already collapsed directory
	model.collapseDirectory()

	// Should remain collapsed
	if collapsedDir.IsExpanded {
		t.Error("Expected directory to remain collapsed")
	}
}

func TestToggleSelection_EdgeCases(t *testing.T) {
	model := NewFileTreeModel()

	// Test with no items
	model.items = []*models.FileNode{}
	model.cursor = 0

	// Should not crash
	model.toggleSelection()

	// Test with cursor out of bounds
	model.items = []*models.FileNode{
		{Name: "file.txt", IsDirectory: false, IsSelected: true},
	}
	model.cursor = 5 // Out of bounds

	// Should not crash
	model.toggleSelection()

	// File should remain selected (no change)
	if !model.items[0].IsSelected {
		t.Error("Expected file to remain selected when cursor out of bounds")
	}
}

// Helper function to create test tree nodes
func createTestTreeNodes() []*models.FileNode {
	return []*models.FileNode{
		{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
		{
			Path:        "/test/file2.txt",
			Name:        "file2.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
		{
			Path:        "/test/dir1",
			Name:        "dir1",
			IsDirectory: true,
			IsSelected:  true,
			IsExpanded:  true,
		},
	}
}
