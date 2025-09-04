package input

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// TaskInputModel represents the task input screen state
type TaskInputModel struct {
	// Text editing components
	textarea textarea.Model
	viewport viewport.Model

	// Content state
	content string

	// Counters
	charCount int
	lineCount int

	// Layout dimensions
	width  int
	height int

	// Screen state
	ready bool
	err   error

	// Key mappings
	keyMap KeyMap
}

// NewTaskInputModel creates a new task input screen model
func NewTaskInputModel() TaskInputModel {
	ta := textarea.New()
	ta.Placeholder = "Describe your task in detail..."
	ta.Focus()
	ta.CharLimit = 10000 // Reasonable limit for task descriptions
	ta.SetWidth(80)
	ta.SetHeight(10)

	vp := viewport.New(80, 15)

	return TaskInputModel{
		textarea:  ta,
		viewport:  vp,
		content:   "",
		charCount: 0,
		lineCount: 1,
		width:     80,
		height:    25,
		ready:     false,
		err:       nil,
		keyMap:    DefaultKeyMap(),
	}
}

// Init initializes the task input model
func (m TaskInputModel) Init() tea.Cmd {
	return textarea.Blink
}

// UpdateSize updates the model dimensions for responsive layout
func (m *TaskInputModel) UpdateSize(width, height int) {
	m.width = width
	m.height = height

	// Calculate textarea dimensions based on available space
	// Leave room for counters, borders, and instructions
	textareaWidth := width - 4  // Account for borders and padding
	textareaHeight := height - 8 // Account for counters and instructions

	if textareaWidth < 40 {
		textareaWidth = 40
	}
	if textareaHeight < 5 {
		textareaHeight = 5
	}

	m.textarea.SetWidth(textareaWidth)
	m.textarea.SetHeight(textareaHeight)

	// Update viewport size
	m.viewport.Width = width - 2
	m.viewport.Height = height - 6
}

// SetContent sets the textarea content and updates counters
func (m *TaskInputModel) SetContent(content string) {
	m.content = content
	m.textarea.SetValue(content)
	m.updateCounters()
}

// GetContent returns the current textarea content
func (m TaskInputModel) GetContent() string {
	return m.textarea.Value()
}

// CanAdvance returns true if the content is valid for advancement
func (m TaskInputModel) CanAdvance() bool {
	content := strings.TrimSpace(m.textarea.Value())
	return len(content) > 0
}

// SetError sets an error state on the model
func (m *TaskInputModel) SetError(err error) {
	m.err = err
}

// GetError returns the current error state
func (m TaskInputModel) GetError() error {
	return m.err
}

// IsReady returns true if the model is ready for interaction
func (m TaskInputModel) IsReady() bool {
	return m.ready
}

// SetReady sets the ready state
func (m *TaskInputModel) SetReady(ready bool) {
	m.ready = ready
}

// updateCounters updates the character and line counters
func (m *TaskInputModel) updateCounters() {
	content := m.textarea.Value()
	m.content = content
	m.charCount = len([]rune(content)) // Use rune count for proper UTF-8 character counting
	m.lineCount = len(strings.Split(content, "\n"))
	if m.lineCount == 0 {
		m.lineCount = 1
	}
}

// Focused returns true if the textarea is currently focused for text input
func (t TaskInputModel) Focused() bool {
	return t.textarea.Focused()
}

// Blur removes focus from the textarea
func (t *TaskInputModel) Blur() {
	t.textarea.Blur()
}
