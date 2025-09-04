package filetree

import (
	"strings"
	"testing"

	"github.com/diogopedro/shotgun/internal/models"
)

func TestView_EmptyModel(t *testing.T) {
	model := NewFileTreeModel()
	view := model.View()

	if !strings.Contains(view, "Loading file tree") {
		t.Error("Expected loading message for empty model")
	}
}

func TestView_WithData(t *testing.T) {
	model := NewFileTreeModel()

	// Create test data
	nodes := []*models.FileNode{
		{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  true,
			IsBinary:    false,
		},
	}

	model.LoadFileTree(nodes)
	view := model.View()

	// Should contain file name
	if !strings.Contains(view, "file1.txt") {
		t.Error("Expected view to contain file1.txt")
	}

	// Should contain status bar
	if !strings.Contains(view, "selected") {
		t.Error("Expected view to contain status information")
	}

	// Should contain help bar
	if !strings.Contains(view, "navigate") {
		t.Error("Expected view to contain help information")
	}
}

func TestRenderTreeItem_File(t *testing.T) {
	model := NewFileTreeModel()

	// Create file item
	item := treeItem{
		node: &models.FileNode{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  true,
			IsBinary:    false,
		},
		depth: 0,
	}

	rendered := model.renderTreeItem(item, false)

	// Should contain file icon
	if !strings.Contains(rendered, "📄") {
		t.Error("Expected file icon in rendered item")
	}

	// Should contain selected checkbox
	if !strings.Contains(rendered, "✅") {
		t.Error("Expected selected checkbox for selected file")
	}

	// Should contain file name
	if !strings.Contains(rendered, "file1.txt") {
		t.Error("Expected file name in rendered item")
	}
}

func TestRenderTreeItem_Directory(t *testing.T) {
	model := NewFileTreeModel()

	// Create expanded directory item
	item := treeItem{
		node: &models.FileNode{
			Path:        "/test/dir1",
			Name:        "dir1",
			IsDirectory: true,
			IsSelected:  true,
			IsExpanded:  true,
			Children: []*models.FileNode{
				{Name: "child.txt", IsDirectory: false},
			},
		},
		depth: 0,
	}

	rendered := model.renderTreeItem(item, false)

	// Should contain directory icon
	if !strings.Contains(rendered, "📁") {
		t.Error("Expected directory icon in rendered item")
	}

	// Should contain expanded indicator
	if !strings.Contains(rendered, "▼") {
		t.Error("Expected expanded indicator for expanded directory")
	}

	// Should contain directory name
	if !strings.Contains(rendered, "dir1") {
		t.Error("Expected directory name in rendered item")
	}
}

func TestRenderTreeItem_BinaryFile(t *testing.T) {
	model := NewFileTreeModel()

	// Create binary file item
	item := treeItem{
		node: &models.FileNode{
			Path:        "/test/binary.exe",
			Name:        "binary.exe",
			IsDirectory: false,
			IsSelected:  false,
			IsBinary:    true,
		},
		depth: 0,
	}

	rendered := model.renderTreeItem(item, false)

	// Should contain binary indicator
	if !strings.Contains(rendered, "⚫") {
		t.Error("Expected binary indicator for binary file")
	}

	// Should contain file name
	if !strings.Contains(rendered, "binary.exe") {
		t.Error("Expected binary file name in rendered item")
	}
}

func TestRenderTreeItem_NestedItem(t *testing.T) {
	model := NewFileTreeModel()

	// Create nested item (depth > 0)
	item := treeItem{
		node: &models.FileNode{
			Path:        "/test/dir1/nested.txt",
			Name:        "nested.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
		depth: 2, // Two levels deep
	}

	rendered := model.renderTreeItem(item, false)

	// Should have proper indentation (2 levels = 4 spaces)
	if !strings.HasPrefix(rendered, "    ") {
		t.Error("Expected proper indentation for nested item")
	}
}

func TestRenderTreeItem_SelectedHighlight(t *testing.T) {
	model := NewFileTreeModel()

	item := treeItem{
		node: &models.FileNode{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
		depth: 0,
	}

	// Test without highlight
	renderedNormal := model.renderTreeItem(item, false)

	// Test with highlight (cursor on item)
	renderedHighlighted := model.renderTreeItem(item, true)

	// Both should contain the file name and basic structure
	if !strings.Contains(renderedNormal, "file1.txt") {
		t.Error("Expected normal render to contain file name")
	}

	if !strings.Contains(renderedHighlighted, "file1.txt") {
		t.Error("Expected highlighted render to contain file name")
	}

	// Note: In testing environment, lipgloss styling may not be applied,
	// so we just verify the function doesn't crash and returns content
}

func TestStatusBar(t *testing.T) {
	model := NewFileTreeModel()

	// Create test nodes with different states (only files, no directories to avoid selection logic complexity)
	nodes := []*models.FileNode{
		{
			Path:        "/test/selected.txt",
			Name:        "selected.txt",
			IsDirectory: false,
			IsSelected:  true,
			IsBinary:    false,
			IsIgnored:   false,
		},
		{
			Path:        "/test/binary.exe",
			Name:        "binary.exe",
			IsDirectory: false,
			IsSelected:  false,
			IsBinary:    true,
			IsIgnored:   false,
		},
		{
			Path:        "/test/ignored.log",
			Name:        "ignored.log",
			IsDirectory: false,
			IsSelected:  false,
			IsBinary:    false,
			IsIgnored:   true,
		},
		{
			Path:        "/test/unselected.txt",
			Name:        "unselected.txt",
			IsDirectory: false,
			IsSelected:  false,
			IsBinary:    false,
			IsIgnored:   false,
		},
	}

	// Don't use LoadFileTree as it auto-selects files; instead set manually
	model.items = nodes

	// Manually set the exact selection states we want for testing
	nodes[0].IsSelected = true  // selected.txt
	nodes[1].IsSelected = false // binary.exe (binary, should be excluded)
	nodes[2].IsSelected = false // ignored.log (ignored)
	nodes[3].IsSelected = false // unselected.txt (unselected, should be excluded)

	statusBar := model.statusBar()

	// Should show 1 selected (using the actual format with emoji)
	if !strings.Contains(statusBar, "✅ 1 selected") {
		t.Errorf("Expected status bar to show '✅ 1 selected', got: %s", statusBar)
	}

	// Should show 2 excluded (1 binary + 1 unselected)
	if !strings.Contains(statusBar, "⚫ 2 excluded") {
		t.Errorf("Expected status bar to show '⚫ 2 excluded', got: %s", statusBar)
	}

	// Should show 1 ignored
	if !strings.Contains(statusBar, "🚫 1 ignored") {
		t.Errorf("Expected status bar to show '🚫 1 ignored', got: %s", statusBar)
	}

	// Should show 3 total (excluding ignored files)
	if !strings.Contains(statusBar, "📄 3 total") {
		t.Errorf("Expected status bar to show '📄 3 total', got: %s", statusBar)
	}
}

func TestCalculateCounts(t *testing.T) {
	model := NewFileTreeModel()

	// Create test tree with directory and files
	childFile1 := &models.FileNode{
		Path:        "/test/dir1/file1.txt",
		Name:        "file1.txt",
		IsDirectory: false,
		IsSelected:  true,
		IsBinary:    false,
		IsIgnored:   false,
	}

	childFile2 := &models.FileNode{
		Path:        "/test/dir1/binary.exe",
		Name:        "binary.exe",
		IsDirectory: false,
		IsSelected:  false,
		IsBinary:    true,
		IsIgnored:   false,
	}

	parentDir := &models.FileNode{
		Path:        "/test/dir1",
		Name:        "dir1",
		IsDirectory: true,
		IsSelected:  true,
		Children:    []*models.FileNode{childFile1, childFile2},
	}

	rootFile := &models.FileNode{
		Path:        "/test/ignored.log",
		Name:        "ignored.log",
		IsDirectory: false,
		IsSelected:  false,
		IsBinary:    false,
		IsIgnored:   true,
	}

	model.LoadFileTree([]*models.FileNode{parentDir, rootFile})

	selected, excluded, ignored := model.calculateCounts()

	// Should count files only, not directories
	if selected != 1 {
		t.Errorf("Expected 1 selected file, got %d", selected)
	}

	if excluded != 1 {
		t.Errorf("Expected 1 excluded file (binary), got %d", excluded)
	}

	if ignored != 1 {
		t.Errorf("Expected 1 ignored file, got %d", ignored)
	}
}

func TestHelpBar(t *testing.T) {
	model := NewFileTreeModel()
	help := model.helpBar()

	// Should contain key navigation hints
	expectedKeys := []string{"↑/↓", "k/j", "←/→", "h/l", "space", "F3", "q"}

	for _, key := range expectedKeys {
		if !strings.Contains(help, key) {
			t.Errorf("Expected help bar to contain key hint: %s", key)
		}
	}
}

func TestRenderTreeItem_EdgeCases(t *testing.T) {
	model := NewFileTreeModel()

	// Test directory with no children (empty directory)
	emptyDir := treeItem{
		node: &models.FileNode{
			Path:        "/test/empty",
			Name:        "empty",
			IsDirectory: true,
			IsExpanded:  true,
			Children:    []*models.FileNode{}, // No children
		},
		depth: 0,
	}

	rendered := model.renderTreeItem(emptyDir, false)

	// Should contain directory icon and name
	if !strings.Contains(rendered, "📁") {
		t.Error("Expected directory icon for empty directory")
	}
	if !strings.Contains(rendered, "empty") {
		t.Error("Expected directory name in rendered item")
	}

	// Test collapsed directory with children
	collapsedDir := treeItem{
		node: &models.FileNode{
			Path:        "/test/collapsed",
			Name:        "collapsed",
			IsDirectory: true,
			IsExpanded:  false,
			Children: []*models.FileNode{
				{Name: "child.txt", IsDirectory: false},
			},
		},
		depth: 0,
	}

	rendered = model.renderTreeItem(collapsedDir, false)

	// Should contain collapsed indicator
	if !strings.Contains(rendered, "▶") {
		t.Error("Expected collapsed indicator for collapsed directory")
	}

	// Test deeply nested item (high depth)
	deepItem := treeItem{
		node: &models.FileNode{
			Path:        "/very/deep/nested/file.txt",
			Name:        "file.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
		depth: 5, // Very deep nesting
	}

	rendered = model.renderTreeItem(deepItem, false)

	// Should have proper indentation (5 levels = 10 spaces)
	if !strings.HasPrefix(rendered, "          ") {
		t.Error("Expected proper indentation for deeply nested item")
	}
}

func TestStatusBar_EmptyTree(t *testing.T) {
	model := NewFileTreeModel()

	// Test with empty tree
	statusBar := model.statusBar()

	// Should show all zeros
	if !strings.Contains(statusBar, "✅ 0 selected") {
		t.Errorf("Expected status bar to show '✅ 0 selected', got: %s", statusBar)
	}
	if !strings.Contains(statusBar, "⚫ 0 excluded") {
		t.Errorf("Expected status bar to show '⚫ 0 excluded', got: %s", statusBar)
	}
	if !strings.Contains(statusBar, "🚫 0 ignored") {
		t.Errorf("Expected status bar to show '🚫 0 ignored', got: %s", statusBar)
	}
	if !strings.Contains(statusBar, "📄 0 total") {
		t.Errorf("Expected status bar to show '📄 0 total', got: %s", statusBar)
	}
}

func TestStatusBar_WithWidth(t *testing.T) {
	model := NewFileTreeModel()
	model.width = 100 // Set specific width

	// Add some test data
	nodes := []*models.FileNode{
		{
			Path:        "/test/file.txt",
			Name:        "file.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
	}

	model.LoadFileTree(nodes)
	statusBar := model.statusBar()

	// Status bar should include padding for width
	// (This is hard to test exactly due to lipgloss styling, but we can verify it runs)
	if len(statusBar) == 0 {
		t.Error("Expected non-empty status bar")
	}
}

func TestFlattenTree(t *testing.T) {
	model := NewFileTreeModel()

	// Create test tree structure
	childFile := &models.FileNode{
		Path:        "/test/dir1/child.txt",
		Name:        "child.txt",
		IsDirectory: false,
		IsExpanded:  false,
	}

	parentDir := &models.FileNode{
		Path:        "/test/dir1",
		Name:        "dir1",
		IsDirectory: true,
		IsExpanded:  true,
		Children:    []*models.FileNode{childFile},
	}

	rootFile := &models.FileNode{
		Path:        "/test/root.txt",
		Name:        "root.txt",
		IsDirectory: false,
		IsExpanded:  false,
	}

	nodes := []*models.FileNode{parentDir, rootFile}
	model.LoadFileTree(nodes)

	flatItems := model.flattenTree(nodes, 0)

	// Should have parent dir, child file, and root file = 3 items
	if len(flatItems) != 3 {
		t.Errorf("Expected 3 flattened items, got %d", len(flatItems))
	}

	// Check ordering and depths
	if flatItems[0].node.Name != "dir1" || flatItems[0].depth != 0 {
		t.Error("Expected first item to be dir1 at depth 0")
	}

	if flatItems[1].node.Name != "child.txt" || flatItems[1].depth != 1 {
		t.Error("Expected second item to be child.txt at depth 1")
	}

	if flatItems[2].node.Name != "root.txt" || flatItems[2].depth != 0 {
		t.Error("Expected third item to be root.txt at depth 0")
	}
}
