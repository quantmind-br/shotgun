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

// TemplateModel implementation moved to internal/screens/template/

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
