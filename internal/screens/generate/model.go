package generate

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/diogopedro/shotgun/internal/core/builder"
)

// GenerateModel manages the prompt generation screen state
type GenerateModel struct {
	// Generation state
	generating    bool
	progress      progress.Model
	spinner       spinner.Model
	
	// Generation components
	generator     *builder.PromptGenerator
	fileWriter    *builder.FileWriter
	
	// Results
	completed     bool
	outputFile    string
	generatedSize int64
	error         error
	
	// UI state
	viewport      viewport.Model
	showStats     bool
	width         int
	height        int
	
	// Generation metadata
	templateSize  int64
	fileCount     int
	totalSize     int64
}

// NewGenerateModel creates a new GenerateModel instance
func NewGenerateModel() GenerateModel {
	p := progress.New(progress.WithDefaultGradient())
	s := spinner.New()
	s.Spinner = spinner.Dot
	
	return GenerateModel{
		generating:   false,
		progress:     p,
		spinner:      s,
		generator:    builder.NewPromptGenerator(),
		fileWriter:   builder.NewFileWriter(),
		completed:    false,
		showStats:    true,
		viewport:     viewport.New(80, 20),
	}
}

// UpdateWindowSize updates the model's window dimensions
func (m *GenerateModel) UpdateWindowSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 4  // Account for borders
	m.viewport.Height = height - 8 // Account for borders and progress
}

// StartGeneration begins the prompt generation process
func (m *GenerateModel) StartGeneration() {
	m.generating = true
	m.completed = false
	m.error = nil
	m.outputFile = ""
	m.generatedSize = 0
}

// CompleteGeneration marks the generation as complete with results
func (m *GenerateModel) CompleteGeneration(result *builder.GeneratedPrompt, outputFile string, err error) {
	m.generating = false
	m.completed = true
	m.error = err
	
	if result != nil {
		m.templateSize = result.TemplateSize
		m.fileCount = result.FileCount
		m.totalSize = result.TotalSize
		m.generatedSize = result.TotalSize
	}
	
	if outputFile != "" {
		m.outputFile = outputFile
	}
}

// IsGenerating returns whether generation is currently in progress
func (m *GenerateModel) IsGenerating() bool {
	return m.generating
}

// IsCompleted returns whether generation has completed
func (m *GenerateModel) IsCompleted() bool {
	return m.completed
}

// HasError returns whether generation completed with an error
func (m *GenerateModel) HasError() bool {
	return m.error != nil
}

// GetError returns the generation error, if any
func (m *GenerateModel) GetError() error {
	return m.error
}

// GetOutputFile returns the path to the generated file
func (m *GenerateModel) GetOutputFile() string {
	return m.outputFile
}

// ToggleStats toggles the display of generation statistics
func (m *GenerateModel) ToggleStats() {
	m.showStats = !m.showStats
}

// ShowingStats returns whether statistics are currently displayed
func (m *GenerateModel) ShowingStats() bool {
	return m.showStats
}