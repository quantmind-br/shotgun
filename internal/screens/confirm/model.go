package confirm

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/diogopedro/shotgun/internal/components/progress"
	"github.com/diogopedro/shotgun/internal/models"
)

// WarningLevel represents the severity of size warnings
type WarningLevel int

const (
	WarningNone WarningLevel = iota
	WarningLarge
	WarningVeryLarge
	WarningExcessive
)

// ConfirmModel holds the state for the confirmation screen
type ConfirmModel struct {
	// Summary data from AppState
	template      *models.Template
	selectedFiles []string
	taskContent   string
	rulesContent  string

	// Size estimation state
	estimatedSize int64
	calculating   bool
	progress      progress.Model
	progressMgr   *ProgressManager
	sizeBreakdown SizeBreakdown

	// Output configuration
	outputFilename string
	outputPath     string

	// UI state
	viewport     viewport.Model
	showWarning  bool
	warningLevel WarningLevel
	width        int
	height       int

	// Navigation state
	ready bool

	// Key mappings
	keyMap KeyMap
}

// SizeBreakdown provides detailed size information
type SizeBreakdown struct {
	TemplateSize    int64
	FileContentSize int64
	TreeStructSize  int64
	OverheadSize    int64
}

// NewConfirmModel creates a new confirmation screen model
func NewConfirmModel() ConfirmModel {
	// Initialize our enhanced progress bar
	p := progress.NewFileProgressModel(100) // Default to 100 files, will be updated
	p.SetWidth(40)

	// Initialize viewport for scrollable content
	vp := viewport.New(78, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return ConfirmModel{
		progress:    p,
		progressMgr: NewProgressManager(),
		viewport:    vp,
		ready:       false,
		keyMap:      DefaultKeyMap(),
	}
}

// UpdateWindowSize updates the model dimensions for responsive display
func (m *ConfirmModel) UpdateWindowSize(width, height int) {
	m.width = width
	m.height = height

	// Update viewport size with padding for borders and headers
	viewportWidth := width - 4    // Account for borders
	viewportHeight := height - 12 // Account for borders, headers, and navigation

	if viewportWidth < 40 {
		viewportWidth = 40
	}
	if viewportHeight < 10 {
		viewportHeight = 10
	}

	m.viewport.Width = viewportWidth
	m.viewport.Height = viewportHeight

	// Update progress bar width
	progressWidth := viewportWidth - 20
	if progressWidth < 20 {
		progressWidth = 20
	}
	m.progress.SetWidth(progressWidth)

	// Update progress manager width too
	if m.progressMgr != nil {
		m.progressMgr.SetWidth(progressWidth)
	}
}

// SetData populates the model with data from AppState
func (m *ConfirmModel) SetData(template *models.Template, selectedFiles []string, taskContent, rulesContent string) {
	m.template = template
	m.selectedFiles = selectedFiles
	m.taskContent = taskContent
	m.rulesContent = rulesContent
	m.ready = true
}

// IsReady returns whether the model has been populated with data
func (m *ConfirmModel) IsReady() bool {
	return m.ready && m.template != nil
}

// GetEstimatedSize returns the current estimated size
func (m *ConfirmModel) GetEstimatedSize() int64 {
	return m.estimatedSize
}

// SetEstimatedSize updates the estimated size and breakdown
func (m *ConfirmModel) SetEstimatedSize(size int64, breakdown SizeBreakdown) {
	m.estimatedSize = size
	m.sizeBreakdown = breakdown
	m.calculating = false

	// Determine warning level based on size
	m.updateWarningLevel()
}

// updateWarningLevel sets warning level based on estimated size
func (m *ConfirmModel) updateWarningLevel() {
	const (
		largeSizeThreshold     = 100 * 1024  // 100KB
		veryLargeSizeThreshold = 500 * 1024  // 500KB
		excessiveSizeThreshold = 2048 * 1024 // 2MB
	)

	switch {
	case m.estimatedSize >= excessiveSizeThreshold:
		m.warningLevel = WarningExcessive
		m.showWarning = true
	case m.estimatedSize >= veryLargeSizeThreshold:
		m.warningLevel = WarningVeryLarge
		m.showWarning = true
	case m.estimatedSize >= largeSizeThreshold:
		m.warningLevel = WarningLarge
		m.showWarning = true
	default:
		m.warningLevel = WarningNone
		m.showWarning = false
	}
}

// StartCalculation marks the model as calculating size
func (m *ConfirmModel) StartCalculation() {
	m.calculating = true
	m.estimatedSize = 0
	m.showWarning = false
}

// IsCalculating returns whether size calculation is in progress
func (m *ConfirmModel) IsCalculating() bool {
	return m.calculating
}

// SetOutputFilename sets the generated output filename
func (m *ConfirmModel) SetOutputFilename(filename string) {
	m.outputFilename = filename
}

// GetOutputFilename returns the current output filename
func (m *ConfirmModel) GetOutputFilename() string {
	return m.outputFilename
}
