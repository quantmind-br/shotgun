package confirm

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	// Initialize progress bar
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 40

	// Initialize spinner for indeterminate states
	s := spinner.New()
	s.Spinner = spinner.Dot
	
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
}

// UpdateProgress updates the current progress state
func (pm *ProgressManager) UpdateProgress(processed int, currentFile string) tea.Cmd {
	pm.state.Processed = processed
	pm.state.CurrentFile = currentFile
	
	if pm.state.Total > 0 {
		pm.state.Percentage = float64(processed) / float64(pm.state.Total)
	}
	
	return func() tea.Msg {
		return ProgressMsg(pm.state)
	}
}

// CompleteProgress marks calculation as complete
func (pm *ProgressManager) CompleteProgress() tea.Cmd {
	pm.state.Completed = true
	pm.state.Percentage = 1.0
	
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
	switch msg := msg.(type) {
	case ProgressMsg:
		pm.state = ProgressState(msg)
		if progressModel, progressCmd := pm.progress.Update(msg); progressCmd != nil {
			pm.progress = progressModel.(progress.Model)
			return progressCmd
		}
		
	case tea.KeyMsg:
		// Handle cancellation during progress
		if msg.String() == "esc" || msg.String() == "q" {
			return pm.CancelProgress()
		}
	}
	
	// Update spinner for indeterminate progress
	if spinnerModel, spinnerCmd := pm.spinner.Update(msg); spinnerCmd != nil {
		pm.spinner = spinnerModel
		return spinnerCmd
	}
	return nil
}

// View renders the progress display
func (pm *ProgressManager) View() string {
	if pm.state.Total == 0 {
		// Indeterminate progress with spinner
		return fmt.Sprintf("%s Calculating sizes...", pm.spinner.View())
	}
	
	// Determinate progress with progress bar
	progressView := pm.progress.ViewAs(pm.state.Percentage)
	
	status := fmt.Sprintf("Processed %d of %d files (%.0f%%)", 
		pm.state.Processed, pm.state.Total, pm.state.Percentage*100)
	
	if pm.state.CurrentFile != "" {
		currentFileDisplay := pm.state.CurrentFile
		if len(currentFileDisplay) > 50 {
			currentFileDisplay = "..." + currentFileDisplay[len(currentFileDisplay)-47:]
		}
		status += fmt.Sprintf("\nCurrent: %s", currentFileDisplay)
	}
	
	return fmt.Sprintf("%s\n%s", progressView, status)
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
	pm.progress.Width = width
}

// ProgressTickCmd returns a command that sends periodic progress updates
func ProgressTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		return struct{}{}
	})
}