package filetree

import (
	"context"
	"testing"
	"time"

	"github.com/user/shotgun-cli/internal/models"
)

func TestNewFileTreeModel(t *testing.T) {
	model := NewFileTreeModel()

	if model.cursor != 0 {
		t.Errorf("Expected cursor = 0, got %d", model.cursor)
	}

	if len(model.items) != 0 {
		t.Errorf("Expected empty items, got %d items", len(model.items))
	}

	if model.selected == nil {
		t.Error("Expected selected map to be initialized")
	}

	if model.width != 80 {
		t.Errorf("Expected width = 80, got %d", model.width)
	}

	if model.height != 24 {
		t.Errorf("Expected height = 24, got %d", model.height)
	}
}

func TestLoadFileTree(t *testing.T) {
	model := NewFileTreeModel()

	// Create test nodes
	nodes := []*models.FileNode{
		{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsSelected:  true,
			IsBinary:    false,
			Size:        100,
			ModTime:     time.Now(),
		},
		{
			Path:        "/test/dir1",
			Name:        "dir1",
			IsDirectory: true,
			IsSelected:  true,
			IsExpanded:  true,
			Size:        0,
			ModTime:     time.Now(),
		},
	}

	model.LoadFileTree(nodes)

	if len(model.items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(model.items))
	}

	if model.cursor != 0 {
		t.Errorf("Expected cursor reset to 0, got %d", model.cursor)
	}

	// Check selected state
	if !model.selected["/test/file1.txt"] {
		t.Error("Expected file1.txt to be selected")
	}

	if !model.selected["/test/dir1"] {
		t.Error("Expected dir1 to be selected")
	}
}

func TestInitializeSelection(t *testing.T) {
	model := NewFileTreeModel()

	// Create test nodes with one binary file
	nodes := []*models.FileNode{
		{
			Path:        "/test/file1.txt",
			Name:        "file1.txt",
			IsDirectory: false,
			IsBinary:    false,
		},
		{
			Path:        "/test/binary.exe",
			Name:        "binary.exe",
			IsDirectory: false,
			IsBinary:    true,
		},
		{
			Path:        "/test/dir1",
			Name:        "dir1",
			IsDirectory: true,
			IsBinary:    false,
			Children: []*models.FileNode{
				{
					Path:        "/test/dir1/nested.txt",
					Name:        "nested.txt",
					IsDirectory: false,
					IsBinary:    false,
				},
			},
		},
	}

	model.initializeSelection(nodes, true)

	// Check normal file is selected
	if !nodes[0].IsSelected {
		t.Error("Expected normal file to be selected")
	}

	// Check binary file is not selected
	if nodes[1].IsSelected {
		t.Error("Expected binary file to not be selected")
	}

	// Check directory is selected
	if !nodes[2].IsSelected {
		t.Error("Expected directory to be selected")
	}

	// Check nested file is selected
	if !nodes[2].Children[0].IsSelected {
		t.Error("Expected nested file to be selected")
	}
}

func TestBuildTreeStructure_EmptyNodes(t *testing.T) {
	model := NewFileTreeModel()

	// Test with empty slice
	var emptyNodes []*models.FileNode
	result := model.buildTreeStructure(emptyNodes)

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty nodes, got %d nodes", len(result))
	}
}

func TestBuildTreeStructure_SingleFile(t *testing.T) {
	model := NewFileTreeModel()

	// Test with single file
	nodes := []*models.FileNode{
		{
			Path:        "/single/file.txt",
			Name:        "file.txt",
			IsDirectory: false,
			IsSelected:  true,
		},
	}

	result := model.buildTreeStructure(nodes)

	if len(result) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(result))
	}

	// The result should be a file.txt node with correct properties
	if result[0].Name != "file.txt" {
		t.Errorf("Expected node name 'file.txt', got '%s'", result[0].Name)
	}
}

func TestBuildTreeStructure(t *testing.T) {
	model := NewFileTreeModel()

	// Create flat list of nodes representing a directory structure
	flatNodes := []*models.FileNode{
		{Path: "/root", Name: "root", IsDirectory: true},
		{Path: "/root/file1.txt", Name: "file1.txt", IsDirectory: false},
		{Path: "/root/subdir", Name: "subdir", IsDirectory: true},
		{Path: "/root/subdir/file2.txt", Name: "file2.txt", IsDirectory: false},
	}

	result := model.buildTreeStructure(flatNodes)

	if len(result) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(result))
	}

	root := result[0]
	if root.Name != "root" {
		t.Errorf("Expected root name 'root', got '%s'", root.Name)
	}

	if len(root.Children) != 2 {
		t.Errorf("Expected root to have 2 children, got %d", len(root.Children))
	}

	// Check that children are properly linked
	var file1, subdir *models.FileNode
	for _, child := range root.Children {
		if child.Name == "file1.txt" {
			file1 = child
		} else if child.Name == "subdir" {
			subdir = child
		}
	}

	if file1 == nil {
		t.Error("Expected to find file1.txt as child of root")
	} else {
		// Check parent relationship for file1
		if file1.Parent != root {
			t.Error("Expected file1.txt parent to be root")
		}
	}

	if subdir == nil {
		t.Error("Expected to find subdir as child of root")
	} else {
		if len(subdir.Children) != 1 {
			t.Errorf("Expected subdir to have 1 child, got %d", len(subdir.Children))
		}

		if subdir.Children[0].Name != "file2.txt" {
			t.Errorf("Expected subdir child to be 'file2.txt', got '%s'", subdir.Children[0].Name)
		}

		// Check parent relationships
		if subdir.Children[0].Parent != subdir {
			t.Error("Expected file2.txt parent to be subdir")
		}
	}
}

func TestInit(t *testing.T) {
	model := NewFileTreeModel()
	cmd := model.Init()

	// Init should return nil command for FileTree
	if cmd != nil {
		t.Error("Expected Init to return nil command")
	}
}

func TestLoadFromScanner(t *testing.T) {
	model := NewFileTreeModel()
	ctx := context.Background()

	// Test with empty root path
	cmd := model.LoadFromScanner(ctx, "")

	// Should return a command
	if cmd == nil {
		t.Error("Expected LoadFromScanner to return a command")
	}

	// Execute the command to test it
	msg := cmd()

	// Should return either ScanCompleteMsg or ScanErrorMsg
	switch msg.(type) {
	case ScanCompleteMsg:
		// Success case
	case ScanErrorMsg:
		// Error case (expected for empty path)
	default:
		t.Errorf("Expected ScanCompleteMsg or ScanErrorMsg, got %T", msg)
	}
}

func TestLoadFromScannerStreaming(t *testing.T) {
	model := NewFileTreeModel()
	ctx := context.Background()

	// Test with empty root path (error case)
	cmd := model.LoadFromScannerStreaming(ctx, "")

	// Should return a command
	if cmd == nil {
		t.Error("Expected LoadFromScannerStreaming to return a command")
	}

	// Execute the command to test it
	msg := cmd()

	// Should return either ScanErrorMsg or ScanCompleteMsg (depending on implementation)
	switch msg.(type) {
	case ScanErrorMsg, ScanCompleteMsg:
		// Both are acceptable responses
	default:
		t.Errorf("Expected ScanErrorMsg or ScanCompleteMsg for empty path, got %T", msg)
	}

	// Test with valid root path
	cmd2 := model.LoadFromScannerStreaming(ctx, ".")
	if cmd2 == nil {
		t.Error("Expected LoadFromScannerStreaming to return a command for valid path")
	}
}

func TestSortTreeNodes(t *testing.T) {
	model := NewFileTreeModel()

	// Create unsorted nodes (mix of files and directories)
	nodes := []*models.FileNode{
		{Name: "zzz.txt", IsDirectory: false},
		{Name: "AAA", IsDirectory: true},
		{Name: "bbb.txt", IsDirectory: false},
		{Name: "CCC", IsDirectory: true},
		{Name: "aaa.txt", IsDirectory: false},
	}

	model.sortTreeNodes(nodes)

	// Expected order: directories first (AAA, CCC), then files (aaa.txt, bbb.txt, zzz.txt)
	expectedOrder := []string{"AAA", "CCC", "aaa.txt", "bbb.txt", "zzz.txt"}

	if len(nodes) != len(expectedOrder) {
		t.Fatalf("Expected %d nodes, got %d", len(expectedOrder), len(nodes))
	}

	for i, expected := range expectedOrder {
		if nodes[i].Name != expected {
			t.Errorf("Expected position %d to be '%s', got '%s'", i, expected, nodes[i].Name)
		}
	}
}
