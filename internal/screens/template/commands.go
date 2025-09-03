package template

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/user/shotgun-cli/internal/core/template"
)

// LoadTemplatesCmd creates a command to load templates from the service
func LoadTemplatesCmd(service template.TemplateService, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		templates, err := service.LoadAllTemplates(ctx)
		if err != nil {
			return TemplateLoadErrorMsg{Error: err}
		}
		return TemplatesLoadedMsg{Templates: templates}
	}
}

// RefreshTemplatesCmd creates a command to refresh templates from the service
func RefreshTemplatesCmd(service template.TemplateService, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		err := service.RefreshTemplates(ctx)
		if err != nil {
			return TemplateLoadErrorMsg{Error: err}
		}

		// After refresh, load all templates
		templates, err := service.LoadAllTemplates(ctx)
		if err != nil {
			return TemplateLoadErrorMsg{Error: err}
		}

		return TemplatesLoadedMsg{Templates: templates}
	}
}
