package builder

import (
	"context"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/diogopedro/shotgun/internal/core/scanner"
)

// TreeFormat defines formatting options for tree visualization
type TreeFormat struct {
	UseUnicode bool // Use Unicode or ASCII tree characters
	ShowSizes  bool // Include file sizes in output
	ShowBinary bool // Show binary file placeholders
	IndentSize int  // Spaces per indentation level
}

// DefaultTreeFormat provides sensible defaults for tree formatting
var DefaultTreeFormat = TreeFormat{
	UseUnicode: false, // Use ASCII characters as specified in AC
	ShowSizes:  false,
	ShowBinary: true,
	IndentSize: 4,
}

// FileStructureBuilder generates tree-structured file content representations
type FileStructureBuilder struct {
	maxFileSize    int64
	maxConcurrency int
	treeFormat     TreeFormat
	binaryDetector *scanner.BinaryDetector
	sensitiveRegex []*regexp.Regexp
	mu             sync.RWMutex
}

// FileStructureBuilderInterface defines the contract for file structure assembly
type FileStructureBuilderInterface interface {
	GenerateStructure(ctx context.Context, files []string) (string, error)
	SetMaxFileSize(size int64) error
	SetMaxConcurrency(workers int) error
	SetTreeFormat(format TreeFormat) error
}

// Option is a functional option for configuring FileStructureBuilder
type Option func(*FileStructureBuilder)

// NewFileStructureBuilder creates a new FileStructureBuilder with functional options
func NewFileStructureBuilder(opts ...Option) *FileStructureBuilder {
	builder := &FileStructureBuilder{
		maxFileSize:    10 * 1024 * 1024, // 10MB default
		maxConcurrency: 10,               // 10 workers default
		treeFormat:     DefaultTreeFormat,
		binaryDetector: scanner.NewBinaryDetectorWithMaxSize(10 * 1024 * 1024), // 10MB for binary detection
		sensitiveRegex: initSensitivePatterns(),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

// initSensitivePatterns returns common sensitive file patterns
func initSensitivePatterns() []*regexp.Regexp {
	patterns := []string{
		`\.env$`,
		`\.env\..*$`,
		`.*\.key$`,
		`.*\.pem$`,
		`.*\.p12$`,
		`.*\.pfx$`,
		`.*\.jks$`,
		`.*\.keystore$`,
		`id_rsa$`,
		`id_ed25519$`,
		`\.ssh/.*$`,
		`secrets?\..*$`,
		`password.*\..*$`,
		`credentials?\..*$`,
		`config/.*\.conf$`,
		`\.aws/.*$`,
	}

	var regexes []*regexp.Regexp
	for _, pattern := range patterns {
		if regex, err := regexp.Compile("(?i)" + pattern); err == nil {
			regexes = append(regexes, regex)
		}
	}
	return regexes
}

// WithMaxFileSize sets the maximum file size for content reading
func WithMaxFileSize(size int64) Option {
	return func(b *FileStructureBuilder) {
		b.maxFileSize = size
	}
}

// WithMaxConcurrency sets the maximum number of concurrent workers
func WithMaxConcurrency(workers int) Option {
	return func(b *FileStructureBuilder) {
		b.maxConcurrency = workers
	}
}

// WithTreeFormat sets the tree formatting options
func WithTreeFormat(format TreeFormat) Option {
	return func(b *FileStructureBuilder) {
		b.treeFormat = format
	}
}

// SetMaxFileSize updates the maximum file size limit
func (b *FileStructureBuilder) SetMaxFileSize(size int64) error {
	if size <= 0 {
		return fmt.Errorf("file size must be positive, got %d", size)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maxFileSize = size
	return nil
}

// SetMaxConcurrency updates the maximum number of concurrent workers
func (b *FileStructureBuilder) SetMaxConcurrency(workers int) error {
	if workers <= 0 {
		return fmt.Errorf("concurrency must be positive, got %d", workers)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maxConcurrency = workers
	return nil
}

// SetTreeFormat updates the tree formatting options
func (b *FileStructureBuilder) SetTreeFormat(format TreeFormat) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.treeFormat = format
	return nil
}

// fileContent represents the result of reading a file
type fileContent struct {
	path    string
	content string
	err     error
}

// GenerateStructure creates a tree-structured representation with file contents
func (b *FileStructureBuilder) GenerateStructure(ctx context.Context, files []string) (string, error) {
	if len(files) == 0 {
		return "", nil
	}

	// Build tree structure from file paths
	tree := b.buildDirectoryTree(files)

	// Pre-load all file contents
	fileContents, err := b.readAllFilesConcurrently(ctx, files)
	if err != nil {
		return "", fmt.Errorf("failed to read file contents: %w", err)
	}

	// Generate tree visualization with content
	var result strings.Builder
	err = b.generateTreeWithContentMap(ctx, tree, "", true, fileContents, &result)
	if err != nil {
		return "", fmt.Errorf("failed to generate tree structure: %w", err)
	}

	return result.String(), nil
}

// readAllFilesConcurrently reads all files using a worker pool
func (b *FileStructureBuilder) readAllFilesConcurrently(ctx context.Context, files []string) (map[string]fileContent, error) {
	// Check context first
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Simplified version: read files sequentially instead of concurrently
	fileContents := make(map[string]fileContent)

	for _, filePath := range files {
		// Check context before each file
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Read file content using the existing readFileContent method
		content, err := b.readFileContent(ctx, filePath)
		if err != nil {
			// Skip files that can't be read instead of failing completely
			fileContents[filePath] = fileContent{
				path:    filePath,
				content: "",
				err:     err,
			}
			continue
		}

		fileContents[filePath] = fileContent{
			path:    filePath,
			content: content,
			err:     nil,
		}
	}

	return fileContents, nil
}

// fileReader is a worker that reads files from the channel
func (b *FileStructureBuilder) fileReader(ctx context.Context, paths <-chan string, results chan<- fileContent, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case path, ok := <-paths:
			if !ok {
				return
			}

			content, err := b.readFileContent(ctx, path)
			select {
			case results <- fileContent{path: path, content: content, err: err}:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// DirectoryNode represents a node in the directory tree
type DirectoryNode struct {
	Name        string
	Path        string
	IsDirectory bool
	Children    map[string]*DirectoryNode
	IsFile      bool
}

// buildDirectoryTree constructs a tree structure from file paths
func (b *FileStructureBuilder) buildDirectoryTree(files []string) *DirectoryNode {
	root := &DirectoryNode{
		Name:        "",
		Path:        "",
		IsDirectory: true,
		Children:    make(map[string]*DirectoryNode),
	}

	for _, file := range files {
		// Normalize path to use forward slashes and store original path
		normalizedFile := filepath.ToSlash(file)
		parts := strings.Split(normalizedFile, "/")
		current := root

		// Build directory structure
		for i, part := range parts {
			if part == "" {
				continue
			}

			if _, exists := current.Children[part]; !exists {
				isDir := i < len(parts)-1
				current.Children[part] = &DirectoryNode{
					Name:        part,
					Path:        file, // Store original file path for final files
					IsDirectory: isDir,
					Children:    make(map[string]*DirectoryNode),
					IsFile:      !isDir,
				}

				// For directories, don't set Path to the full file path
				if isDir {
					current.Children[part].Path = ""
				}
			}
			current = current.Children[part]
		}
	}

	return root
}

// generateTreeWithContentMap recursively generates tree visualization with pre-loaded file contents
func (b *FileStructureBuilder) generateTreeWithContentMap(ctx context.Context, node *DirectoryNode, prefix string, isLast bool, fileContents map[string]fileContent, result *strings.Builder) error {
	// Skip empty root node
	if node.Name == "" {
		// Sort children for consistent output
		children := make([]*DirectoryNode, 0, len(node.Children))
		for _, child := range node.Children {
			children = append(children, child)
		}
		sort.Slice(children, func(i, j int) bool {
			// Directories first, then files
			if children[i].IsDirectory != children[j].IsDirectory {
				return children[i].IsDirectory
			}
			return children[i].Name < children[j].Name
		})

		for i, child := range children {
			isChildLast := i == len(children)-1
			if err := b.generateTreeWithContentMap(ctx, child, "", isChildLast, fileContents, result); err != nil {
				return err
			}
		}
		return nil
	}

	// Generate tree characters
	treeChar := "├── "
	if isLast {
		treeChar = "└── "
	}

	// Write node line
	result.WriteString(prefix + treeChar + node.Name + "\n")

	// Handle files - add content from pre-loaded map
	if node.IsFile {
		if fileContent, exists := fileContents[node.Path]; exists {
			if fileContent.err != nil {
				result.WriteString(fmt.Sprintf("<file path=\"%s\">ERROR: %v</file>\n", node.Path, fileContent.err))
			} else {
				result.WriteString(fmt.Sprintf("<file path=\"%s\">%s</file>\n", node.Path, fileContent.content))
			}
		} else {
			result.WriteString(fmt.Sprintf("<file path=\"%s\">ERROR: File not found in content map</file>\n", node.Path))
		}
	}

	// Handle directories - process children
	if node.IsDirectory && len(node.Children) > 0 {
		// Sort children for consistent output
		children := make([]*DirectoryNode, 0, len(node.Children))
		for _, child := range node.Children {
			children = append(children, child)
		}
		sort.Slice(children, func(i, j int) bool {
			// Directories first, then files
			if children[i].IsDirectory != children[j].IsDirectory {
				return children[i].IsDirectory
			}
			return children[i].Name < children[j].Name
		})

		// Generate prefix for children
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}

		for i, child := range children {
			isChildLast := i == len(children)-1
			if err := b.generateTreeWithContentMap(ctx, child, childPrefix, isChildLast, fileContents, result); err != nil {
				return err
			}
		}
	}

	return nil
}

// readFileContent reads file content with size limits and binary detection
func (b *FileStructureBuilder) readFileContent(ctx context.Context, filePath string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	// Check file size limit
	b.mu.RLock()
	maxSize := b.maxFileSize
	b.mu.RUnlock()

	if info.Size() > maxSize {
		return fmt.Sprintf("File too large (%d bytes, limit %d bytes)", info.Size(), maxSize), nil
	}

	// Check if potentially sensitive file
	if b.isSensitiveFile(filePath) {
		return fmt.Sprintf("⚠️ Potentially sensitive file detected (%d bytes) - Use caution with file contents", info.Size()), nil
	}

	// Check if binary file
	if b.binaryDetector.IsBinary(filePath) {
		return fmt.Sprintf("Binary file (%d bytes)", info.Size()), nil
	}

	// Read file content
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Escape XML special characters for proper XML wrapping
	return html.EscapeString(string(content)), nil
}

// isSensitiveFile checks if a file path matches sensitive file patterns
func (b *FileStructureBuilder) isSensitiveFile(filePath string) bool {
	// Normalize path for consistent matching
	normalizedPath := filepath.ToSlash(filePath)

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, regex := range b.sensitiveRegex {
		if regex.MatchString(normalizedPath) {
			return true
		}
	}
	return false
}
