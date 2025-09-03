package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/user/shotgun-cli/internal/components/progress"
	"github.com/user/shotgun-cli/internal/screens/filetree"
)

// ScreenType represents the different screens in the wizard
type ScreenType int

const (
	FileTreeScreen ScreenType = iota
	TemplateScreen
	TaskScreen
	RulesScreen
	ConfirmScreen
)

// String returns the string representation of ScreenType
func (s ScreenType) String() string {
	switch s {
	case FileTreeScreen:
		return "FileTree"
	case TemplateScreen:
		return "Template"
	case TaskScreen:
		return "TaskInput"
	case RulesScreen:
		return "RulesInput"
	case ConfirmScreen:
		return "Confirm"
	default:
		return "Unknown"
	}
}

// Template represents a template selection
type Template struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

// InputModel represents generic input screen models
type InputModel struct {
	viewport viewport.Model
	content  string
	cursor   int
	width    int
	height   int
}

// TemplateModel represents template selection screen model
type TemplateModel struct {
	items    []*Template
	cursor   int
	selected *Template
	viewport viewport.Model
	width    int
	height   int
}

// ConfirmModel represents confirmation screen model
type ConfirmModel struct {
	summary  string
	viewport viewport.Model
	width    int
	height   int
}

// AppState manages the overall application state
type AppState struct {
	// Current screen
	CurrentScreen ScreenType

	// Screen models
	FileTree     filetree.FileTreeModel
	Template     TemplateModel
	TaskInput    InputModel
	RulesInput   InputModel
	Confirmation ConfirmModel

	// Progress indicator
	Progress progress.Model

	// Shared data across screens
	SelectedFiles    []string
	SelectedTemplate *Template
	TaskContent      string
	RulesContent     string

	// UI state
	WindowSize tea.WindowSizeMsg
	Error      error

	// Dialog states
	ShowingHelp bool
	HelpContent string
	ShowingExit bool

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApp creates a new application state with default values
func NewApp() *AppState {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Define screen titles
	screenTitles := []string{
		"Select Files",
		"Choose Template", 
		"Describe Task",
		"Add Rules (Optional)",
		"Review & Confirm",
	}
	
	app := &AppState{
		CurrentScreen:    FileTreeScreen,
		SelectedFiles:    make([]string, 0),
		SelectedTemplate: nil,
		TaskContent:      "",
		RulesContent:     "",
		ShowingHelp:     false,
		HelpContent:     "",
		ShowingExit:     false,
		ctx:              ctx,
		cancel:           cancel,
	}

	// Initialize screen models with defaults
	app.FileTree = filetree.NewFileTreeModel()
	app.Template = TemplateModel{
		items:    make([]*Template, 0),
		cursor:   0,
		selected: nil,
	}
	app.TaskInput = InputModel{
		content: "",
		cursor:  0,
	}
	app.RulesInput = InputModel{
		content: "",
		cursor:  0,
	}
	app.Confirmation = ConfirmModel{
		summary: "",
	}
	
	// Initialize progress indicator
	app.Progress = progress.NewModel(1, 5, screenTitles)

	return app
}

// Cleanup releases resources
func (a *AppState) Cleanup() {
	if a.cancel != nil {
		a.cancel()
	}
}

// Context returns the application context
func (a *AppState) Context() context.Context {
	return a.ctx
}

// UpdateWindowSize updates the window size for all screen models
func (a *AppState) UpdateWindowSize(msg tea.WindowSizeMsg) {
	a.WindowSize = msg
	
	// Update viewport sizes for all models
	a.FileTree.SetSize(msg.Width, msg.Height)
	
	a.Template.width = msg.Width
	a.Template.height = msg.Height
	a.Template.viewport.Width = msg.Width
	a.Template.viewport.Height = msg.Height - 6 // Reserve space for header/footer
	
	a.TaskInput.width = msg.Width
	a.TaskInput.height = msg.Height
	a.TaskInput.viewport.Width = msg.Width
	a.TaskInput.viewport.Height = msg.Height - 6
	
	a.RulesInput.width = msg.Width
	a.RulesInput.height = msg.Height
	a.RulesInput.viewport.Width = msg.Width
	a.RulesInput.viewport.Height = msg.Height - 6
	
	a.Confirmation.width = msg.Width
	a.Confirmation.height = msg.Height
	a.Confirmation.viewport.Width = msg.Width
	a.Confirmation.viewport.Height = msg.Height - 6
}

// GetCurrentScreenModel returns the current active screen model
func (a *AppState) GetCurrentScreenModel() tea.Model {
	switch a.CurrentScreen {
	case FileTreeScreen:
		return &a.FileTree
	case TemplateScreen:
		return &a.Template
	case TaskScreen:
		return &a.TaskInput
	case RulesScreen:
		return &a.RulesInput
	case ConfirmScreen:
		return &a.Confirmation
	default:
		return &a.FileTree
	}
}

// SetCurrentScreen changes the current screen and updates shared state
func (a *AppState) SetCurrentScreen(screen ScreenType) {
	// Save current screen state before switching
	a.saveCurrentScreenState()
	
	a.CurrentScreen = screen
	
	// Update progress indicator
	a.Progress.SetCurrent(int(screen) + 1)
	
	// Load state into new screen
	a.loadScreenState(screen)
}

// saveCurrentScreenState saves the current screen's data to shared state
func (a *AppState) saveCurrentScreenState() {
	switch a.CurrentScreen {
	case FileTreeScreen:
		a.SelectedFiles = a.FileTree.GetSelectedFiles()
	case TemplateScreen:
		a.SelectedTemplate = a.Template.selected
	case TaskScreen:
		a.TaskContent = a.TaskInput.content
	case RulesScreen:
		a.RulesContent = a.RulesInput.content
	}
}

// loadScreenState loads shared data into the specified screen
func (a *AppState) loadScreenState(screen ScreenType) {
	switch screen {
	case FileTreeScreen:
		// FileTree state is managed internally
	case TemplateScreen:
		a.Template.selected = a.SelectedTemplate
	case TaskScreen:
		a.TaskInput.content = a.TaskContent
	case RulesScreen:
		a.RulesInput.content = a.RulesContent
	case ConfirmScreen:
		// Build confirmation summary
		a.buildConfirmationSummary()
	}
}

// buildConfirmationSummary creates a summary for the confirmation screen
func (a *AppState) buildConfirmationSummary() {
	summary := "Configuration Summary:\n\n"
	
	if len(a.SelectedFiles) > 0 {
		summary += "Selected Files:\n"
		for _, file := range a.SelectedFiles {
			summary += "  • " + file + "\n"
		}
		summary += "\n"
	}
	
	if a.SelectedTemplate != nil {
		summary += "Selected Template: " + a.SelectedTemplate.Name + "\n\n"
	}
	
	if a.TaskContent != "" {
		summary += "Task Description:\n" + a.TaskContent + "\n\n"
	}
	
	if a.RulesContent != "" {
		summary += "Custom Rules:\n" + a.RulesContent + "\n\n"
	}
	
	a.Confirmation.summary = summary
}