package input

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RulesInputModel represents the state for the optional rules input screen
type RulesInputModel struct {
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

// NewRulesInputModel creates a new RulesInputModel with default configuration
func NewRulesInputModel() RulesInputModel {
	ta := textarea.New()
	ta.Placeholder = "Optional: Add rules or constraints to guide the AI response...\nExample: 'Write in a conversational tone' or 'Keep responses under 500 words'"
	ta.Focus()
	ta.CharLimit = 5000 // Reasonable limit for rules
	ta.SetWidth(80)
	ta.SetHeight(8) // Smaller than task input since rules are typically shorter

	// Enable word wrap and proper text handling
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(true)

	// Configure for multiline editing with UTF-8 support
	// (The textarea is already properly configured by default)

	vp := viewport.New(80, 15)

	return RulesInputModel{
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

// Init initializes the RulesInputModel
func (m RulesInputModel) Init() RulesInputModel {
	return m
}

// UpdateSize updates the model's dimensions for responsive layout
func (m *RulesInputModel) UpdateSize(width, height int) {
	// Guard against invalid dimensions
	if width <= 0 || height <= 0 {
		return
	}

	m.width = width
	m.height = height

	// Update textarea dimensions with padding for UI elements
	textareaWidth := width - 4
	textareaHeight := height - 10 // Leave room for header, footer, counters

	// Ensure minimum viable dimensions
	if textareaWidth < 20 {
		textareaWidth = 20
	}
	if textareaHeight < 3 {
		textareaHeight = 3
	}

	m.textarea.SetWidth(textareaWidth)
	m.textarea.SetHeight(textareaHeight)
	m.viewport.Width = textareaWidth
	m.viewport.Height = textareaHeight + 2
}

// SetContent sets the content and updates counters
func (m *RulesInputModel) SetContent(content string) {
	m.content = content
	m.textarea.SetValue(content)
	m.updateCounters()
}

// GetContent returns the current content
func (m RulesInputModel) GetContent() string {
	return m.content
}

// CanAdvance returns true since rules input is optional (no validation required)
func (m RulesInputModel) CanAdvance() bool {
	return true // Rules are optional, can always advance
}

// SetError sets an error state
func (m *RulesInputModel) SetError(err error) {
	m.err = err
}

// GetError returns the current error
func (m RulesInputModel) GetError() error {
	return m.err
}

// IsReady returns the ready state
func (m RulesInputModel) IsReady() bool {
	return m.ready
}

// SetReady sets the ready state
func (m *RulesInputModel) SetReady(ready bool) {
	m.ready = ready
}

// Message types for Rules Input

// RulesInputMsg represents a message to advance from rules input screen
type RulesInputMsg struct{}

// BackToTaskMsg represents a message to return to task screen
type BackToTaskMsg struct{}

// RulesContentUpdatedMsg represents a message when rules content is updated
type RulesContentUpdatedMsg struct {
	Content string
}

// SkipRulesMsg represents a message to skip rules input entirely
type SkipRulesMsg struct{}

// Update handles messages and updates the RulesInputModel state
func (m RulesInputModel) Update(msg tea.Msg) (RulesInputModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update model dimensions for responsive layout
		m.UpdateSize(msg.Width, msg.Height)

	case RulesContentUpdatedMsg:
		// Handle external content updates (e.g., from AppState)
		m.SetContent(msg.Content)

	case ClipboardPasteMsg:
		// Handle clipboard paste result
		current := m.textarea.Value()

		// For now, append pasted text - cursor positioning will be handled by textarea
		newContent := current + msg.Text

		m.textarea.SetValue(newContent)
		m.updateCounters()

	case ClipboardErrorMsg:
		// Handle clipboard operation errors
		m.SetError(msg.Error)

	case tea.KeyMsg:
		switch msg.String() {
		case "f4":
			// F4 skips rules screen entirely
			cmds = append(cmds, func() tea.Msg {
				return SkipRulesMsg{}
			})

		case "f3":
			// F3 advances regardless of content (rules are optional)
			// Clear any previous errors
			m.SetError(nil)
			// Send message to advance to confirmation screen
			cmds = append(cmds, func() tea.Msg {
				return RulesInputMsg{}
			})

		case "f2":
			// Return to task screen with state preservation
			cmds = append(cmds, func() tea.Msg {
				return BackToTaskMsg{}
			})

		case "ctrl+c":
			// Copy selected text to clipboard
			// Note: For now, we'll copy all content since selection API is not available
			content := m.textarea.Value()
			if content != "" {
				cmds = append(cmds, func() tea.Msg {
					return ClipboardCopyMsg{Text: content}
				})
			}

		case "ctrl+v":
			// Initiate clipboard paste operation
			cmds = append(cmds, func() tea.Msg {
				// This would be handled by the app layer to perform actual clipboard access
				return ClipboardPasteMsg{Text: ""} // Placeholder - actual text comes from clipboard
			})

		default:
			// Let the textarea handle other keys (typing, cursor movement, etc.)
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)

			// Update counters after text changes
			m.updateCounters()
		}

	default:
		// Update textarea with other messages
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)

		// Update counters after potential text changes
		m.updateCounters()
	}

	return m, tea.Batch(cmds...)
}

// View renders the RulesInputModel screen
func (m RulesInputModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Loading..."
	}

	// Handle error state
	if m.err != nil {
		return m.renderError()
	}

	// Render normal screen
	return m.renderMain()
}

// renderMain renders the main rules input screen
func (m RulesInputModel) renderMain() string {
	var sections []string

	// Header with "Optional" indication
	header := headerStyle.Width(m.width).Render("📋 Rules & Constraints (Optional)")
	sections = append(sections, header)

	// Instructions
	instruction := instructionStyle.Width(m.width).Render(
		"Optional: Add rules or constraints to guide the AI's response style or requirements. This field can be left empty.",
	)
	sections = append(sections, instruction)

	// Text area with proper styling based on focus
	textareaContent := m.textarea.View()
	if m.textarea.Focused() {
		textareaContent = focusedTextareaStyle.Width(m.width - 4).Render(textareaContent)
	} else {
		textareaContent = textareaStyle.Width(m.width - 4).Render(textareaContent)
	}
	sections = append(sections, textareaContent)

	// Character and line counts
	countText := fmt.Sprintf("Lines: %d | Characters: %d", m.lineCount, m.charCount)
	countDisplay := countStyle.Width(m.width - 4).Render(countText)
	sections = append(sections, countDisplay)

	// Help text for keyboard shortcuts
	helpText := []string{
		"F3: Continue to confirmation",
		"F4: Skip this screen entirely",
		"F2: Back to Task Description",
		"Ctrl+C: Copy • Ctrl+V: Paste",
	}
	help := helpStyle.Width(m.width - 4).Render(strings.Join(helpText, " • "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderError renders the error state
func (m RulesInputModel) renderError() string {
	var sections []string

	// Header
	header := headerStyle.Width(m.width).Render("📋 Rules & Constraints (Optional) - Error")
	sections = append(sections, header)

	// Error message
	errorMsg := errorStyle.Width(m.width - 4).Render(fmt.Sprintf("Error: %s", m.err.Error()))
	sections = append(sections, errorMsg)

	// Instructions
	instruction := instructionStyle.Width(m.width).Render("Please try again or press F2 to go back.")
	sections = append(sections, instruction)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// updateCounters updates character and line counters
func (m *RulesInputModel) updateCounters() {
	// Update content from textarea
	m.content = m.textarea.Value()

	// UTF-8 aware character counting
	m.charCount = len([]rune(m.content))

	// Line counting
	if m.content == "" {
		m.lineCount = 1
	} else {
		m.lineCount = 1
		for _, char := range m.content {
			if char == '\n' {
				m.lineCount++
			}
		}
	}
}

// Focused returns true if the textarea is currently focused for text input
func (r RulesInputModel) Focused() bool {
	return r.textarea.Focused()
}

// Blur removes focus from the textarea
func (r *RulesInputModel) Blur() {
	r.textarea.Blur()
}
