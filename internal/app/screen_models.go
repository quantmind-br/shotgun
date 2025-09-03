package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Implement tea.Model interface for InputModel

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m InputModel) View() string {
	return m.content
}

// Implement tea.Model interface for TemplateModel

func (m TemplateModel) Init() tea.Cmd {
	return nil
}

func (m TemplateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TemplateModel) View() string {
	if len(m.items) == 0 {
		return "No templates available"
	}
	return "Templates loaded"
}

// Implement tea.Model interface for ConfirmModel

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConfirmModel) View() string {
	return m.summary
}