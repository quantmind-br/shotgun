package help

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ScreenType represents the different screens in the application
type ScreenType int

const (
	FileTreeScreen ScreenType = iota
	TemplateScreen
	TaskScreen
	RulesScreen
	ConfirmScreen
	GenerateScreen
)

// HelpItem represents a single help entry
type HelpItem struct {
	Key         string
	Description string
	Context     ScreenType
}

// HelpModel represents the help overlay component
type HelpModel struct {
	visible       bool
	content       []HelpItem
	viewport      viewport.Model
	currentScreen ScreenType
	width         int
	height        int
}

// NewHelpModel creates a new help model
func NewHelpModel() HelpModel {
	vp := viewport.New(0, 0)
	return HelpModel{
		visible:  false,
		viewport: vp,
	}
}

// SetVisible toggles the help overlay visibility
func (h *HelpModel) SetVisible(visible bool) {
	h.visible = visible
}

// IsVisible returns whether the help overlay is currently shown
func (h HelpModel) IsVisible() bool {
	return h.visible
}

// SetCurrentScreen updates the current screen context for help content
func (h *HelpModel) SetCurrentScreen(screen ScreenType) {
	h.currentScreen = screen
	h.updateContent()
}

// UpdateSize updates the help overlay dimensions
func (h *HelpModel) UpdateSize(width, height int) {
	h.width = width
	h.height = height

	// Set viewport size (leave space for borders and title)
	h.viewport.Width = width - 4
	h.viewport.Height = height - 8
}

// updateContent refreshes the help content based on current screen
func (h *HelpModel) updateContent() {
	h.content = GetHelpContent(h.currentScreen)
	h.viewport.SetContent(h.formatContent())
}

// formatContent formats help items for display
func (h HelpModel) formatContent() string {
	var content string

	// Global shortcuts section
	content += "Navigation:\n"
	for _, item := range h.content {
		if item.Context == h.currentScreen || item.Context == -1 { // -1 for global
			content += "  " + item.Key + "  " + item.Description + "\n"
		}
	}

	content += "\nGlobal:\n"
	content += "  F1        Show/hide this help\n"
	content += "  F2        Previous screen\n"
	content += "  F3        Next screen\n"
	if h.currentScreen == TaskScreen || h.currentScreen == RulesScreen {
		content += "  F4        Skip to next section\n"
	}
	if h.currentScreen == ConfirmScreen {
		content += "  F10       Generate prompt\n"
	}
	content += "  ESC       Exit application\n"

	content += "\nPress F1 to close help"

	return content
}

// Init implements tea.Model
func (h HelpModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (h HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if h.visible && msg.String() == "f1" {
			h.visible = false
			return h, nil
		}

		// Pass viewport navigation to viewport
		if h.visible {
			h.viewport, cmd = h.viewport.Update(msg)
		}
	}

	return h, cmd
}
