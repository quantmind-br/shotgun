package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// App represents the main application
type App struct {
	state *AppState
}

// NewApplication creates a new application instance
func NewApplication() *App {
	return &App{
		state: NewApp(),
	}
}

// Run starts the application
func (app *App) Run() error {
	p := tea.NewProgram(app.state, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// GetState returns the current application state
func (app *App) GetState() *AppState {
	return app.state
}

// Shutdown performs cleanup operations
func (app *App) Shutdown() {
	if app.state != nil {
		app.state.Cleanup()
	}
}