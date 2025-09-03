package template

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages for the template selection screen
func (m TemplateModel) Update(msg tea.Msg) (TemplateModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update model dimensions
		m.UpdateSize(msg.Width, msg.Height)

	case TemplatesLoadedMsg:
		// Handle templates loaded from service
		m.SetTemplates(msg.Templates)

	case TemplateLoadErrorMsg:
		// Handle template loading error
		m.SetError(msg.Error)

	case tea.KeyMsg:
		// Don't handle keys if still loading
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			// Move cursor up
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			// Update selected template
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if templateItem, ok := selectedItem.(TemplateItem); ok {
					m.selected = &templateItem.Template
				}
			}

		case "down", "j":
			// Move cursor down
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			// Update selected template
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if templateItem, ok := selectedItem.(TemplateItem); ok {
					m.selected = &templateItem.Template
				}
			}

		case "enter", "f3":
			// Select template and advance to next screen
			if m.selected != nil {
				cmds = append(cmds, func() tea.Msg {
					return TemplateSelectedMsg{Template: m.selected}
				})
			}

		case "f2":
			// Return to file tree screen
			cmds = append(cmds, func() tea.Msg {
				return BackToFileTreeMsg{}
			})

		case "pgup":
			// Page up
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			// Update selected template
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if templateItem, ok := selectedItem.(TemplateItem); ok {
					m.selected = &templateItem.Template
				}
			}

		case "pgdown":
			// Page down
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			// Update selected template
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if templateItem, ok := selectedItem.(TemplateItem); ok {
					m.selected = &templateItem.Template
				}
			}

		case "home":
			// Go to first item
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			// Update selected template
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if templateItem, ok := selectedItem.(TemplateItem); ok {
					m.selected = &templateItem.Template
				}
			}

		case "end":
			// Go to last item
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			// Update selected template
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if templateItem, ok := selectedItem.(TemplateItem); ok {
					m.selected = &templateItem.Template
				}
			}

		case "ctrl+r":
			// Refresh templates
			cmds = append(cmds, func() tea.Msg {
				return RefreshTemplatesMsg{}
			})

		case "tab":
			// Toggle detail panel visibility
			m.showDetails = !m.showDetails
			m.UpdateSize(m.width, m.height) // Recalculate layout

		default:
			// Let the list handle other keys
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		// Update list with other messages
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Additional messages for template screen
type BackToFileTreeMsg struct{}

type RefreshTemplatesMsg struct{}
