package confirm

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/components/common"
	"github.com/diogopedro/shotgun/internal/components/progress"
	"github.com/diogopedro/shotgun/internal/components/spinner"
)

// ProgressState represents the current state of size calculation
type ProgressState struct {
	Processed   int
	Total       int
	CurrentFile string
	Percentage  float64
	Completed   bool
}

// ProgressMsg is sent when progress updates
type ProgressMsg ProgressState

// SizeCalculationCompleteMsg is sent when calculation is done
type SizeCalculationCompleteMsg struct {
	TotalSize int64
	Breakdown SizeBreakdown
	Error     error
}

// CancellationMsg is sent when user cancels calculation
type CancellationMsg struct{}

// ProgressManager manages progress tracking for size calculation
type ProgressManager struct {
	progress progress.Model
	spinner  spinner.Model
	state    ProgressState
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewProgressManager creates a new progress manager
func NewProgressManager() *ProgressManager {
	// Initialize our enhanced progress bar
	p := progress.NewFileProgressModel(100) // Default to 100 files, will be updated
	p.SetWidth(40)

	// Initialize spinner for indeterminate states
	s := spinner.New(spinner.SpinnerDots)

	return &ProgressManager{
		progress: p,
		spinner:  s,
	}
}

// StartProgress initializes progress tracking
func (pm *ProgressManager) StartProgress(total int) {
	pm.ctx, pm.cancel = context.WithCancel(context.Background())
	pm.state = ProgressState{
		Processed:  0,
		Total:      total,
		Percentage: 0.0,
		Completed:  false,
	}

	// Update progress bar with total
	pm.progress.SetFileCount(0, total)
}

// UpdateProgress updates the current progress state
func (pm *ProgressManager) UpdateProgress(processed int, currentFile string) tea.Cmd {
	pm.state.Processed = processed
	pm.state.CurrentFile = currentFile

	if pm.state.Total > 0 {
		pm.state.Percentage = float64(processed) / float64(pm.state.Total)
	}

	// Update progress bar
	pm.progress.SetFileCount(processed, pm.state.Total)
	pm.progress.SetMessage(fmt.Sprintf("Processing: %s", currentFile))

	return func() tea.Msg {
		return ProgressMsg(pm.state)
	}
}

// CompleteProgress marks calculation as complete
func (pm *ProgressManager) CompleteProgress() tea.Cmd {
	pm.state.Completed = true
	pm.state.Percentage = 1.0

	// Complete progress bar
	pm.progress.SetFileCount(pm.state.Total, pm.state.Total)
	pm.progress.SetMessage("Calculation complete")

	return func() tea.Msg {
		return ProgressMsg(pm.state)
	}
}

// CancelProgress cancels the ongoing calculation
func (pm *ProgressManager) CancelProgress() tea.Cmd {
	if pm.cancel != nil {
		pm.cancel()
	}

	return func() tea.Msg {
		return CancellationMsg{}
	}
}

// GetContext returns the cancellation context
func (pm *ProgressManager) GetContext() context.Context {
	if pm.ctx == nil {
		pm.ctx = context.Background()
	}
	return pm.ctx
}

// Update handles progress-related messages
func (pm *ProgressManager) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ProgressMsg:
		pm.state = ProgressState(msg)

	case tea.KeyMsg:
		// Handle cancellation during progress
		if msg.String() == "esc" || msg.String() == "q" {
			return pm.CancelProgress()
		}
	}

	// Update progress bar
	if progressModel, progressCmd := pm.progress.Update(msg); progressCmd != nil {
		pm.progress = progressModel
		cmds = append(cmds, progressCmd)
	}

	// Update spinner for indeterminate progress
	if spinnerModel, spinnerCmd := pm.spinner.Update(msg); spinnerCmd != nil {
		pm.spinner = spinnerModel
		cmds = append(cmds, spinnerCmd)
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}

// View renders the progress display
func (pm *ProgressManager) View() string {
	if pm.state.Total == 0 {
		// Indeterminate progress with spinner
		pm.spinner.SetMessage("Calculating sizes...")
		pm.spinner.Start()
		return pm.spinner.ViewWithCancel()
	}

	// Determinate progress with our enhanced progress bar
	return pm.progress.View()
}

// GetState returns the current progress state
func (pm *ProgressManager) GetState() ProgressState {
	return pm.state
}

// IsCompleted returns whether calculation is completed
func (pm *ProgressManager) IsCompleted() bool {
	return pm.state.Completed
}

// GetPercentage returns the current completion percentage
func (pm *ProgressManager) GetPercentage() float64 {
	return pm.state.Percentage
}

// SetWidth updates the progress bar width
func (pm *ProgressManager) SetWidth(width int) {
	pm.progress.SetWidth(width)
}

// ProgressTickCmd returns a command that sends periodic progress updates
func ProgressTickCmd() tea.Cmd {
	return tea.Tick(common.ProgressUpdateRate, func(time.Time) tea.Msg {
		return struct{}{}
	})
}
