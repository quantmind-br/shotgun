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
		keyStr := msg.String()

		// ABSOLUTELY CRITICAL: Intercept alt+c before ANYTHING else
		if keyStr == "alt+c" {
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
			// CRITICAL: Return immediately, bypass ALL textarea processing
			return m, tea.Batch(cmds...)
		}

		// Also try with alternative representations
		if len(keyStr) > 5 { // modifier + key
			if keyStr == "alt+c" {
				// Same logic as above
				if m.CanAdvance() {
					m.SetError(nil)
					cmds = append(cmds, func() tea.Msg {
						return TaskInputMsg{}
					})
				} else {
					m.SetError(errors.New("task description cannot be empty"))
				}
				return m, tea.Batch(cmds...)
			}
		}

		// Handle other control keys
		switch keyStr {
		case "ctrl+left":
			cmds = append(cmds, func() tea.Msg {
				return BackToTemplateMsg{}
			})
			return m, tea.Batch(cmds...)

		case "ctrl+c":
			content := m.textarea.Value()
			if content != "" {
				cmds = append(cmds, func() tea.Msg {
					return ClipboardCopyMsg{Text: content}
				})
			}
			return m, tea.Batch(cmds...)

		case "ctrl+v":
			cmds = append(cmds, func() tea.Msg {
				return ClipboardPasteMsg{Text: ""}
			})
			return m, tea.Batch(cmds...)

		default:
			// Pass everything else to textarea
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
			m.updateCounters()
		}

	default:
		// Update textarea with other messages
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
		m.updateCounters()
	}

	return m, tea.Batch(cmds...)
}

// (No enhanced key handler required on Bubble Tea v1.x)
