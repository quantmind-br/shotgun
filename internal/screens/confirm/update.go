package confirm

import (
	"context"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages for the confirmation screen
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateWindowSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "f10":
			// Confirm and trigger generation
			if !m.calculating && m.estimatedSize > 0 {
				return m, ConfirmGenerationCmd()
			}

		case "f2":
			// Return to rules input screen
			if !m.calculating {
				return m, NavigateToRulesCmd()
			}

		case "f1":
			// Return to file selection screen
			if !m.calculating {
				return m, NavigateToFileTreeCmd()
			}

		case "esc":
			// Cancel calculation if running, otherwise exit
			if m.calculating {
				return m, CancelSizeCalculationCmd()
			}
			return m, NavigateToExitCmd()

		case "up", "k":
			// Scroll viewport up
			m.viewport.LineUp(1)

		case "down", "j":
			// Scroll viewport down
			m.viewport.LineDown(1)

		case "pgup":
			// Page up in viewport
			m.viewport.HalfViewUp()

		case "pgdown":
			// Page down in viewport
			m.viewport.HalfViewDown()

		case "home":
			// Go to top
			m.viewport.GotoTop()

		case "end":
			// Go to bottom
			m.viewport.GotoBottom()
		}

	case ProgressMsg:
		// Update progress state
		progressState := ProgressState(msg)
		m.progress.SetPercent(progressState.Percentage)

	case SizeCalculationCompleteMsg:
		// Handle completed size calculation
		if msg.Error != nil {
			// Handle error - could add error state to model
			m.calculating = false
		} else {
			m.SetEstimatedSize(msg.TotalSize, msg.Breakdown)
		}

	case SizeCalculationStartMsg:
		// Start size calculation
		m.StartCalculation()
		cmds = append(cmds, StartSizeCalculationCmd(m.selectedFiles))

	case CancellationMsg:
		// Handle cancelled calculation
		m.calculating = false
		m.estimatedSize = 0

	case FilenameGeneratedMsg:
		// Handle generated filename
		m.SetOutputFilename(msg.Filename)
	}

	// Update progress bar
	if progressModel, progressCmd := m.progress.Update(msg); progressCmd != nil {
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, progressCmd)
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// Message types for screen navigation
type (
	NavigateToRulesMsg    struct{}
	NavigateToFileTreeMsg struct{}
	NavigateToExitMsg     struct{}
	ConfirmGenerationMsg  struct{}
	SizeCalculationStartMsg struct{}
	FilenameGeneratedMsg struct {
		Filename string
	}
)

// Navigation command functions
func NavigateToRulesCmd() tea.Cmd {
	return func() tea.Msg {
		return NavigateToRulesMsg{}
	}
}

func NavigateToFileTreeCmd() tea.Cmd {
	return func() tea.Msg {
		return NavigateToFileTreeMsg{}
	}
}

func NavigateToExitCmd() tea.Cmd {
	return func() tea.Msg {
		return NavigateToExitMsg{}
	}
}

func ConfirmGenerationCmd() tea.Cmd {
	return func() tea.Msg {
		return ConfirmGenerationMsg{}
	}
}

// Size calculation command functions
func StartSizeCalculationCmd(selectedFiles []string) tea.Cmd {
	return func() tea.Msg {
		// This would trigger the actual size calculation
		// For now, return a start message
		return SizeCalculationStartMsg{}
	}
}

func CancelSizeCalculationCmd() tea.Cmd {
	return func() tea.Msg {
		return CancellationMsg{}
	}
}

// CalculateSizeWithProgressCmd performs size calculation with progress updates  
func CalculateSizeWithProgressCmd(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual size calculation
		// For now, return a mock result
		return SizeCalculationCompleteMsg{
			TotalSize: 1024,
			Breakdown: SizeBreakdown{
				TemplateSize:    200,
				FileContentSize: 600,
				TreeStructSize:  124,
				OverheadSize:    100,
			},
		}
	}
}

// GenerateFilenameCmd generates a timestamped filename
func GenerateFilenameCmd() tea.Cmd {
	return func() tea.Msg {
		generator := NewFilenameGenerator("")
		filename := generator.GenerateTimestampFilename()
		
		return FilenameGeneratedMsg{
			Filename: filename,
		}
	}
}

// InitializeConfirmScreenCmd sets up the confirmation screen with data
func InitializeConfirmScreenCmd(template interface{}, selectedFiles []string, taskContent, rulesContent string) tea.Cmd {
	return func() tea.Msg {
		// This would be called when transitioning to confirm screen
		return tea.Batch(
			GenerateFilenameCmd(),
			func() tea.Msg {
				return SizeCalculationStartMsg{}
			},
		)()
	}
}