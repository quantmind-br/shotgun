package confirm

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
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
        case "ctrl+enter", "enter":
            // Confirm and trigger generation
            if !m.calculating && m.estimatedSize > 0 {
                return m, ConfirmGenerationCmd()
            }

        case "ctrl+left":
            // Return to rules input screen
            if !m.calculating {
                return m, NavigateToRulesCmd()
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
		// Update progress state via the ProgressManager
		if m.calculating && m.progressMgr != nil {
			progressCmd := m.progressMgr.UpdateProgress(msg.Processed, msg.CurrentFile)
			if progressCmd != nil {
				cmds = append(cmds, progressCmd)
			}
		}

	case SizeCalculationCompleteMsg:
		// Handle completed size calculation
		if msg.Error != nil {
			// Handle error - could add error state to model
			m.calculating = false
		} else {
			m.SetEstimatedSize(msg.TotalSize, msg.Breakdown)
			// Complete progress manager
			if m.progressMgr != nil {
				completeCmd := m.progressMgr.CompleteProgress()
				if completeCmd != nil {
					cmds = append(cmds, completeCmd)
				}
			}
		}

	case SizeCalculationStartMsg:
		// Start size calculation
		m.StartCalculation()
		// Start progress tracking with estimated file count
		m.progressMgr.StartProgress(len(m.selectedFiles) + 3) // files + template + task + rules
		ctx := m.progressMgr.GetContext()
		cmds = append(cmds, CalculateSizeWithProgressCmd(ctx, m.selectedFiles, m.template, m.taskContent, m.rulesContent))

	case CancellationMsg:
		// Handle cancelled calculation
		m.calculating = false
		m.estimatedSize = 0
		// Cancel progress manager if active
		if m.progressMgr != nil {
			cancelCmd := m.progressMgr.CancelProgress()
			if cancelCmd != nil {
				cmds = append(cmds, cancelCmd)
			}
		}

	case FilenameGeneratedMsg:
		// Handle generated filename
		m.SetOutputFilename(msg.Filename)
	}

	// Update progress bar
	if progressModel, progressCmd := m.progress.Update(msg); progressCmd != nil {
		m.progress = progressModel
		cmds = append(cmds, progressCmd)
	}

	// Update progress manager
	if m.progressMgr != nil {
		if progressCmd := m.progressMgr.Update(msg); progressCmd != nil {
			cmds = append(cmds, progressCmd)
		}
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
	NavigateToRulesMsg      struct{}
	NavigateToFileTreeMsg   struct{}
	NavigateToExitMsg       struct{}
	ConfirmGenerationMsg    struct{}
	SizeCalculationStartMsg struct{}
	FilenameGeneratedMsg    struct {
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
func CalculateSizeWithProgressCmd(ctx context.Context, selectedFiles []string, template *models.Template, taskContent, rulesContent string) tea.Cmd {
	return tea.Sequence(
		// Start progress indicator
		func() tea.Msg {
			return ProgressMsg{
				Processed:   0,
				Total:       len(selectedFiles) + 3, // files + template + task + rules
				CurrentFile: "Starting calculation...",
				Percentage:  0.0,
				Completed:   false,
			}
		},
		// Perform calculation with progress updates
		func() tea.Msg {
			return calculateSizeWithProgress(ctx, selectedFiles, template, taskContent, rulesContent)
		},
	)
}

// calculateSizeWithProgress performs the actual size calculation with progress updates
func calculateSizeWithProgress(ctx context.Context, selectedFiles []string, template *models.Template, taskContent, rulesContent string) tea.Msg {
	var breakdown SizeBreakdown
	processed := 0
	total := len(selectedFiles) + 3 // files + template + task + rules

	// Check for cancellation
	select {
	case <-ctx.Done():
		return CancellationMsg{}
	default:
	}

	// Calculate template size
	if template != nil {
		breakdown.TemplateSize = int64(len(template.Content) + len(template.Name) + len(template.Description))
		processed++
	}

	// Calculate task content size
	breakdown.TreeStructSize = int64(len(taskContent))
	processed++

	// Calculate rules content size
	breakdown.OverheadSize = int64(len(rulesContent))
	processed++

	// Calculate file content sizes
	for i, filePath := range selectedFiles {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return CancellationMsg{}
		default:
		}

		// Send progress update
		currentFile := fmt.Sprintf("Processing %s", filePath)
		percentage := float64(processed+i) / float64(total)

		// In a real implementation, you would read the file here
		// For now, we'll estimate based on filename
		estimatedSize := int64(len(filePath) * 100) // Mock estimation
		breakdown.FileContentSize += estimatedSize

		processed++

		// Send progress update (would be sent via channel in real implementation)
		if i%5 == 0 || i == len(selectedFiles)-1 { // Update every 5 files or on last file
			_ = ProgressMsg{
				Processed:   processed,
				Total:       total,
				CurrentFile: currentFile,
				Percentage:  percentage,
				Completed:   false,
			}
		}
	}

	// Final calculation complete
	totalSize := breakdown.TemplateSize + breakdown.FileContentSize + breakdown.TreeStructSize + breakdown.OverheadSize

	return SizeCalculationCompleteMsg{
		TotalSize: totalSize,
		Breakdown: breakdown,
		Error:     nil,
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
