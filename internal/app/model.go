package app

import (
    "context"

    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"

	"github.com/diogopedro/shotgun/internal/components/help"
	"github.com/diogopedro/shotgun/internal/components/progress"
	"github.com/diogopedro/shotgun/internal/models"
	"github.com/diogopedro/shotgun/internal/screens/confirm"
	"github.com/diogopedro/shotgun/internal/screens/filetree"
	"github.com/diogopedro/shotgun/internal/screens/generate"
	"github.com/diogopedro/shotgun/internal/screens/input"
    "github.com/diogopedro/shotgun/internal/screens/template"
    tmplcore "github.com/diogopedro/shotgun/internal/core/template"
)

// ScreenType represents the different screens in the wizard
type ScreenType int

const (
	FileTreeScreen ScreenType = iota
	TemplateScreen
	TaskScreen
	RulesScreen
	ConfirmScreen
	GenerateScreen
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
	case GenerateScreen:
		return "Generate"
	default:
		return "Unknown"
	}
}

// Template represents a template selection
// Legacy template struct - remove once migration complete
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
// Legacy TemplateModel - remove once migration complete

// ConfirmModel represents confirmation screen model
type ConfirmModel = confirm.ConfirmModel

// AppState manages the overall application state
type AppState struct {
	// Current screen
	CurrentScreen ScreenType

	// Screen models
	FileTree     filetree.FileTreeModel
	Template     template.TemplateModel
	TaskInput    input.TaskInputModel
	RulesInput   input.RulesInputModel
	Confirmation ConfirmModel
	Generation   generate.GenerateModel

	// Progress indicator
	Progress progress.Model

	// Help overlay
	Help help.HelpModel

	// Shared data across screens
	SelectedFiles    []string
	SelectedTemplate *models.Template
	TaskContent      string
	RulesContent     string

	// UI state
	WindowSize tea.WindowSizeMsg
	Error      error

	// Dialog states
	ShowingHelp bool
	HelpContent string
	ShowingExit bool

	// Input mode tracking
	InputMode bool

    // Context for cancellation
    ctx    context.Context
    cancel context.CancelFunc

    // Services
    templateService tmplcore.TemplateService
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
		"Generate Prompt",
	}

    app := &AppState{
        CurrentScreen:    FileTreeScreen,
		SelectedFiles:    make([]string, 0),
		SelectedTemplate: nil,
		TaskContent:      "",
		RulesContent:     "",
		ShowingHelp:      false,
		HelpContent:      "",
		ShowingExit:      false,
        ctx:              ctx,
        cancel:           cancel,
    }

	// Initialize screen models with defaults
	app.FileTree = filetree.NewFileTreeModel()
	app.Template = template.NewTemplateModel()
	app.TaskInput = input.NewTaskInputModel()
	app.RulesInput = input.NewRulesInputModel()
	app.Confirmation = confirm.NewConfirmModel()
    app.Generation = generate.NewGenerateModel()

    // Initialize services
    app.templateService = tmplcore.NewTemplateService(nil)

	// Initialize progress indicator
	app.Progress = progress.NewModel(1, 6, screenTitles)

	// Initialize help overlay
	app.Help = help.NewHelpModel()
	app.Help.SetCurrentScreen(help.ScreenType(app.CurrentScreen))

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

	a.Template.UpdateSize(msg.Width, msg.Height)

	a.TaskInput.UpdateSize(msg.Width, msg.Height)

	a.RulesInput.UpdateSize(msg.Width, msg.Height)

	a.Confirmation.UpdateWindowSize(msg.Width, msg.Height)

	a.Generation.UpdateWindowSize(msg.Width, msg.Height)

	// Update help overlay size
	a.Help.UpdateSize(msg.Width, msg.Height)
}

// GetCurrentScreenModel returns the current active screen model
// Note: This is primarily for testing purposes
func (a *AppState) GetCurrentScreenModel() interface{} {
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
	case GenerateScreen:
		return &a.Generation
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

	// Update help overlay context
	a.Help.SetCurrentScreen(help.ScreenType(screen))

	// Load state into new screen
	a.loadScreenState(screen)
}

// saveCurrentScreenState saves the current screen's data to shared state
func (a *AppState) saveCurrentScreenState() {
	switch a.CurrentScreen {
	case FileTreeScreen:
		a.SelectedFiles = a.FileTree.GetSelectedFiles()
	case TemplateScreen:
		a.SelectedTemplate = a.Template.GetSelected()
	case TaskScreen:
		a.TaskContent = a.TaskInput.GetContent()
	case RulesScreen:
		a.RulesContent = a.RulesInput.GetContent()
	}
}

// loadScreenState loads shared data into the specified screen
func (a *AppState) loadScreenState(screen ScreenType) {
	switch screen {
	case FileTreeScreen:
		// FileTree state is managed internally
	case TemplateScreen:
		// Template selection is managed internally by the template screen
	case TaskScreen:
		a.TaskInput.SetContent(a.TaskContent)
	case RulesScreen:
		a.RulesInput.SetContent(a.RulesContent)
	case ConfirmScreen:
		// Build confirmation summary
		a.buildConfirmationSummary()
	case GenerateScreen:
		// Initialize generation screen with current app state
		a.initializeGenerationScreen()
	}
}

// buildConfirmationSummary creates a summary for the confirmation screen
func (a *AppState) buildConfirmationSummary() {
	// Set the confirmation data using the proper method
	a.Confirmation.SetData(a.SelectedTemplate, a.SelectedFiles, a.TaskContent, a.RulesContent)
}

// initializeGenerationScreen prepares the generation screen with current app state
func (a *AppState) initializeGenerationScreen() {
	// The generation screen will be initialized when generation starts
	// We don't need to do anything here as the screen handles its own initialization
}
