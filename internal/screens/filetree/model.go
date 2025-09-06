package filetree

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/components/spinner"
	"github.com/diogopedro/shotgun/internal/core/scanner"
	"github.com/diogopedro/shotgun/internal/models"
)

// Message types for scanner integration
type ScanCompleteMsg struct {
	Nodes []*models.FileNode
}

type ScanErrorMsg struct {
	Error error
}

type ScanProgressMsg struct {
	FilesFound int
	CurrentDir string
}

// FileTreeModel represents the state of the file tree screen
type FileTreeModel struct {
	items    []*models.FileNode
	cursor   int
	selected map[string]bool
	viewport viewport.Model
	width    int
	height   int
	keyMap   KeyMap
	// Loading state fields
	scanning   bool
	spinner    spinner.Model
	scanError  error
	filesFound int
	currentDir string
}

// NewFileTreeModel creates a new FileTreeModel with defaults
func NewFileTreeModel() FileTreeModel {
	vp := viewport.New(80, 20)
	vp.YPosition = 0

	return FileTreeModel{
		items:    make([]*models.FileNode, 0),
		cursor:   0,
		selected: make(map[string]bool),
		viewport: vp,
		width:    80,
		height:   24,
		keyMap:   DefaultKeyMap(),
		scanning: false,
		spinner:  spinner.New(spinner.SpinnerDots),
	}
}

// LoadFileTree initializes the model with FileNode data and sets all files as selected by default
func (m *FileTreeModel) LoadFileTree(nodes []*models.FileNode) {
	m.items = nodes
	m.cursor = 0
	m.selected = make(map[string]bool)

	// Initialize all files as selected by default (IsSelected: true)
	m.initializeSelection(nodes, true)
}

// initializeSelection recursively sets initial selection state
func (m *FileTreeModel) initializeSelection(nodes []*models.FileNode, isSelected bool) {
	for _, node := range nodes {
		if !node.IsBinary { // Don't select binary files
			node.IsSelected = isSelected
			m.selected[node.Path] = isSelected
		} else {
			node.IsSelected = false
			m.selected[node.Path] = false
		}

		// Recursively handle children
		if node.IsDirectory && len(node.Children) > 0 {
			m.initializeSelection(node.Children, isSelected)
		}
	}
}

// Init implements the Bubble Tea Model interface
func (m FileTreeModel) Init() tea.Cmd {
	return nil
}

// StartScanning initiates the scanning process
func (m *FileTreeModel) StartScanning() tea.Cmd {
	m.scanning = true
	m.scanError = nil
	m.filesFound = 0
	m.currentDir = ""
	return m.spinner.Start()
}

// StopScanning stops the scanning process
func (m *FileTreeModel) StopScanning() {
	m.scanning = false
	m.spinner.Stop()
}

// IsScanning returns whether scanning is in progress
func (m FileTreeModel) IsScanning() bool {
	return m.scanning
}

// LoadFromScanner loads file tree data from the scanner service with loading state
func (m *FileTreeModel) LoadFromScanner(ctx context.Context, rootPath string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// Create scanner with default options
		scannerInstance, err := scanner.New()
		if err != nil {
			return ScanErrorMsg{Error: err}
		}

		// Scan directory synchronously for simplicity
		nodes, err := scannerInstance.ScanDirectorySync(ctx, rootPath)
		if err != nil {
			return ScanErrorMsg{Error: err}
		}

		// Convert flat list to tree structure
		treeNodes := m.buildTreeStructure(nodes)

		return ScanCompleteMsg{Nodes: treeNodes}
	})
}

// LoadFromScannerStreaming loads file tree data using streaming scanner with progress updates
func (m *FileTreeModel) LoadFromScannerStreaming(ctx context.Context, rootPath string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// Create scanner with default options
		scannerInstance, err := scanner.New()
		if err != nil {
			return ScanErrorMsg{Error: err}
		}

		// Start streaming scan
		resultChan, err := scannerInstance.ScanDirectory(ctx, rootPath)
		if err != nil {
			return ScanErrorMsg{Error: err}
		}

		var nodes []*models.FileNode
		filesFound := 0

		for result := range resultChan {
			if result.Error != nil {
				return ScanErrorMsg{Error: result.Error}
			}
			if result.FileNode != nil {
				nodes = append(nodes, result.FileNode)
				filesFound++

				// Send progress updates periodically (every 10 files)
				if filesFound%10 == 0 {
					// Note: In a real streaming implementation, we'd send these as separate messages
					// For now, we'll just count them and show at the end
				}
			}
		}

		// Convert flat list to tree structure
		treeNodes := m.buildTreeStructure(nodes)

		return ScanCompleteMsg{Nodes: treeNodes}
	})
}

// LoadFromScannerWithProgress loads file tree with enhanced progress tracking
func (m *FileTreeModel) LoadFromScannerWithProgress(ctx context.Context, rootPath string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// Create scanner with default options
		scannerInstance, err := scanner.New()
		if err != nil {
			return ScanErrorMsg{Error: err}
		}

		// For synchronous scanning, we'll simulate progress
		// In a real implementation, this would use the streaming scanner
		nodes, err := scannerInstance.ScanDirectorySync(ctx, rootPath)
		if err != nil {
			return ScanErrorMsg{Error: err}
		}

		// Convert flat list to tree structure
		treeNodes := m.buildTreeStructure(nodes)

		return ScanCompleteMsg{Nodes: treeNodes}
	})
}

// buildTreeStructure converts flat list of FileNode to hierarchical tree
func (m *FileTreeModel) buildTreeStructure(flatNodes []*models.FileNode) []*models.FileNode {
	if len(flatNodes) == 0 {
		return nil
	}

	// Create maps for quick lookups
	nodeMap := make(map[string]*models.FileNode)
	pathToParentMap := make(map[string]string)

	// First pass: populate maps and prepare nodes
	for _, node := range flatNodes {
		// Clean the path
		cleanPath := filepath.Clean(node.Path)
		node.Path = cleanPath
		nodeMap[cleanPath] = node

		// Initialize expanded state for directories
		if node.IsDirectory {
			node.IsExpanded = true // Start expanded by default
		}

		// Find parent path
		parentPath := filepath.Dir(cleanPath)
		if parentPath != cleanPath && parentPath != "." {
			pathToParentMap[cleanPath] = parentPath
		}
	}

	// Second pass: build parent-child relationships
	var rootNodes []*models.FileNode

	for _, node := range flatNodes {
		parentPath, hasParent := pathToParentMap[node.Path]

		if hasParent {
			// This node has a parent
			if parent, exists := nodeMap[parentPath]; exists {
				// Set parent relationship
				node.Parent = parent
				// Add to parent's children
				parent.Children = append(parent.Children, node)
			} else {
				// Parent doesn't exist in our scan results, treat as root
				rootNodes = append(rootNodes, node)
			}
		} else {
			// This is a root node
			rootNodes = append(rootNodes, node)
		}
	}

	// Sort children in each directory
	m.sortTreeNodes(rootNodes)
	for _, node := range flatNodes {
		if node.IsDirectory {
			m.sortTreeNodes(node.Children)
		}
	}

	return rootNodes
}

// sortTreeNodes sorts nodes with directories first, then files, alphabetically
func (m *FileTreeModel) sortTreeNodes(nodes []*models.FileNode) {
	sort.Slice(nodes, func(i, j int) bool {
		// Directories come before files
		if nodes[i].IsDirectory != nodes[j].IsDirectory {
			return nodes[i].IsDirectory
		}
		// Within same type, sort alphabetically
		return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
	})
}
