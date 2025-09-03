package input

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

// TaskInputMsg represents a message to advance to the task input screen
type TaskInputMsg struct{}

// BackToTemplateMsg represents a message to return to template screen
type BackToTemplateMsg struct{}

// TaskContentUpdatedMsg represents a message when task content is updated
type TaskContentUpdatedMsg struct {
	Content string
}

// ClipboardCopyMsg represents a clipboard copy operation
type ClipboardCopyMsg struct {
	Text string
}

// ClipboardPasteMsg represents a clipboard paste operation  
type ClipboardPasteMsg struct {
	Text string
}

// ClipboardErrorMsg represents a clipboard operation error
type ClipboardErrorMsg struct {
	Error error
}

// Update handles messages for the task input model
func (m TaskInputModel) Update(msg tea.Msg) (TaskInputModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update model dimensions for responsive layout
		m.UpdateSize(msg.Width, msg.Height)

	case TaskContentUpdatedMsg:
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
		case "ctrl+enter":
			// Validate content and advance if valid
			if m.CanAdvance() {
				// Clear any previous errors
				m.SetError(nil)
				// Send message to advance to next screen
				cmds = append(cmds, func() tea.Msg {
					return TaskInputMsg{}
				})
			} else {
				// Set validation error
				m.SetError(errors.New("task description cannot be empty"))
			}

		case "f3":
			// F3 advances only if content is non-empty
			if m.CanAdvance() {
				// Clear any previous errors
				m.SetError(nil)
				// Send message to advance to next screen
				cmds = append(cmds, func() tea.Msg {
					return TaskInputMsg{}
				})
			} else {
				// Set validation error - content is empty
				m.SetError(errors.New("task description cannot be empty"))
			}

		case "f2":
			// Return to template screen with state preservation
			cmds = append(cmds, func() tea.Msg {
				return BackToTemplateMsg{}
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