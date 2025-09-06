package generate

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/diogopedro/shotgun/internal/core/builder"
)

// Update handles all messages for the generation screen
func (m GenerateModel) Update(msg tea.Msg) (GenerateModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateWindowSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.generating {
				// Cancel generation
				return m, CancelGenerationCmd()
			}
			// Exit generation screen
			return m, NavigateBackCmd()

		case "f1":
			// Return to file tree screen for restart
			return m, NavigateToFileTreeCmd()

		case "f2":
			if m.completed && !m.HasError() && m.outputFile != "" {
				// Open generated file in system default application
				return m, OpenFileCmd(m.outputFile)
			}

		case "f5":
			if m.HasError() {
				// Retry generation
				return m, RetryGenerationCmd()
			}

		case "s":
			// Toggle statistics display
			m.ToggleStats()
		}

	case builder.GenerationProgressMsg:
		// Update progress from generation
		if m.generating {
			// Convert progress to 0-1 range for progress bar
			m.progress.SetPercent(msg.Progress)
		}

	case builder.GenerationCompleteMsg:
		// Generation completed
		if msg.Error != nil {
			m.CompleteGeneration(nil, "", msg.Error)
		} else {
			// Write the generated prompt to file
			return m, WritePromptToFileCmd(msg.Result)
		}

	case FileWriteCompleteMsg:
		// File writing completed
		m.CompleteGeneration(msg.Result, msg.OutputFile, msg.Error)

	case GenerationCancelledMsg:
		// Generation was cancelled
		m.generating = false
		m.completed = false

	case StartGenerationMsg:
		// Start generation process
		m.StartGeneration()
		return m, StartGenerationCmd(msg.Config)
	}

	// Update spinner if generating
	if m.generating {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}

	// Update progress bar
	var progressCmd tea.Cmd
	progressModel, progressCmd := m.progress.Update(msg)
	m.progress = progressModel.(progress.Model)
	if progressCmd != nil {
		cmds = append(cmds, progressCmd)
	}

	// Update viewport
	var viewportCmd tea.Cmd
	m.viewport, viewportCmd = m.viewport.Update(msg)
	if viewportCmd != nil {
		cmds = append(cmds, viewportCmd)
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Message types for generation screen

// StartGenerationMsg begins the generation process
type StartGenerationMsg struct {
	Config builder.GenerationConfig
}

// FileWriteCompleteMsg indicates file writing has completed
type FileWriteCompleteMsg struct {
	Result     *builder.GeneratedPrompt
	OutputFile string
	Error      error
}

// GenerationCancelledMsg indicates generation was cancelled
type GenerationCancelledMsg struct{}

// NavigateBackMsg requests navigation back to previous screen
type NavigateBackMsg struct{}

// NavigateToFileTreeMsg requests navigation to file tree screen
type NavigateToFileTreeMsg struct{}

// OpenFileMsg requests opening a file in system default application
type OpenFileMsg struct {
	FilePath string
}

// RetryGenerationMsg requests retrying generation
type RetryGenerationMsg struct{}
