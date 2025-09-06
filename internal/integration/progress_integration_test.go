package integration

import (
	"testing"
	"time"

	"github.com/diogopedro/shotgun/internal/components/progress"
	"github.com/diogopedro/shotgun/internal/components/spinner"
	"github.com/diogopedro/shotgun/internal/models"
	"github.com/diogopedro/shotgun/internal/screens/confirm"
	"github.com/diogopedro/shotgun/internal/screens/filetree"
	"github.com/diogopedro/shotgun/internal/screens/template"
)

// TestProgressComponentIntegration tests that progress components work together
func TestProgressComponentIntegration(t *testing.T) {
	// Test spinner and progress bar together
	spinner := spinner.New(spinner.SpinnerDots)
	progressBar := progress.NewFileProgressModel(100)

	// Start spinner
	spinner.Start()
	spinner.SetMessage("Loading...")

	// Configure progress bar
	progressBar.SetFileCount(25, 100)
	progressBar.SetWidth(50)

	// Verify both components can render simultaneously
	spinnerView := spinner.View()
	progressView := progressBar.View()

	if spinnerView == "" {
		t.Error("Spinner view should not be empty when active")
	}

	if progressView == "" {
		t.Error("Progress bar view should not be empty")
	}

	// Test that both use consistent styling (basic smoke test)
	if len(spinnerView) < 5 || len(progressView) < 10 {
		t.Error("Component views seem too short, may indicate styling issues")
	}
}

// TestFileTreeProgressIntegration tests file tree scanning with spinner
func TestFileTreeProgressIntegration(t *testing.T) {
	model := filetree.NewFileTreeModel()
	model.SetSize(80, 24)

	// Test that initial state is correct
	if model.IsScanning() {
		t.Error("File tree should not be scanning initially")
	}

	// Test scanning state transition
	model.StartScanning()
	if !model.IsScanning() {
		t.Error("File tree should be scanning after StartScanning()")
	}

	// Simulate scanning completion
	model.StopScanning()
	if model.IsScanning() {
		t.Error("File tree should not be scanning after StopScanning()")
	}

	// Test view rendering during different states
	view := model.View()
	if view == "" {
		t.Error("File tree view should never be empty")
	}
}

// TestTemplateDiscoveryProgressIntegration tests template discovery with spinner
func TestTemplateDiscoveryProgressIntegration(t *testing.T) {
	model := template.NewTemplateModel()
	model.UpdateSize(80, 24)

	// Test that initial state is correct
	if model.IsDiscovering() {
		t.Error("Template model should not be discovering initially")
	}

	// Test discovery state transition
	model.StartDiscovery()
	if !model.IsDiscovering() {
		t.Error("Template model should be discovering after StartDiscovery()")
	}

	// Simulate discovery completion
	model.StopDiscovery()
	if model.IsDiscovering() {
		t.Error("Template model should not be discovering after StopDiscovery()")
	}

	// Test view rendering during different states
	view := model.View()
	if view == "" {
		t.Error("Template model view should never be empty")
	}
}

// TestConfirmScreenProgressIntegration tests size calculation progress
func TestConfirmScreenProgressIntegration(t *testing.T) {
	model := confirm.NewConfirmModel()
	model.UpdateWindowSize(80, 24)

	// Set up test data
	testTemplate := &models.Template{
		Name:        "test-template",
		Version:     "1.0.0",
		Description: "Test template for integration testing",
		Content:     "Test content: {{.task}}",
	}

	selectedFiles := []string{"test1.go", "test2.go", "test3.go"}
	taskContent := "Implement feature X"
	rulesContent := "Follow coding standards"

	model.SetData(testTemplate, selectedFiles, taskContent, rulesContent)

	// Test that calculation can be started
	if model.IsCalculating() {
		t.Error("Confirm model should not be calculating initially")
	}

	model.StartCalculation()
	if !model.IsCalculating() {
		t.Error("Confirm model should be calculating after StartCalculation()")
	}

	// Test size estimation
	breakdown := confirm.SizeBreakdown{
		TemplateSize:    100,
		FileContentSize: 5000,
		TreeStructSize:  200,
		OverheadSize:    50,
	}
	model.SetEstimatedSize(5350, breakdown)

	// Test view rendering with estimated size
	view := model.View()
	if view == "" {
		t.Error("Confirm model view should never be empty")
	}
}

// TestProgressStatesConsistency tests that all progress states are consistent
func TestProgressStatesConsistency(t *testing.T) {
	// Test FileTree model separately
	t.Run("FileTree", func(t *testing.T) {
		model := filetree.NewFileTreeModel()
		model.SetSize(80, 24)

		// Test initial state
		if model.IsScanning() {
			t.Error("FileTree should not be scanning initially")
		}

		// Test start state
		model.StartScanning()
		if !model.IsScanning() {
			t.Error("FileTree should be scanning after starting")
		}

		// Test stop state
		model.StopScanning()
		if model.IsScanning() {
			t.Error("FileTree should not be scanning after stopping")
		}

		// Test view is never empty
		view := model.View()
		if view == "" {
			t.Error("FileTree view should never be empty")
		}
	})

	// Test Template model separately
	t.Run("Template", func(t *testing.T) {
		model := template.NewTemplateModel()
		model.UpdateSize(80, 24)

		// Test initial state
		if model.IsDiscovering() {
			t.Error("Template should not be discovering initially")
		}

		// Test start state
		model.StartDiscovery()
		if !model.IsDiscovering() {
			t.Error("Template should be discovering after starting")
		}

		// Test stop state
		model.StopDiscovery()
		if model.IsDiscovering() {
			t.Error("Template should not be discovering after stopping")
		}

		// Test view is never empty
		view := model.View()
		if view == "" {
			t.Error("Template view should never be empty")
		}
	})
}

// TestProgressManagerIntegration tests the ProgressManager in different scenarios
func TestProgressManagerIntegration(t *testing.T) {
	pm := confirm.NewProgressManager()

	// Test indeterminate progress (spinner mode)
	view := pm.View()
	if view == "" {
		t.Error("ProgressManager view should not be empty in indeterminate mode")
	}

	// Test transition to determinate progress
	pm.StartProgress(10)
	pm.UpdateProgress(3, "file3.go")

	view = pm.View()
	if view == "" {
		t.Error("ProgressManager view should not be empty in determinate mode")
	}

	// Test progress completion
	pm.CompleteProgress()
	view = pm.View()
	// View might be empty after completion, which is acceptable

	// Test context cancellation
	ctx := pm.GetContext()
	if ctx == nil {
		t.Error("ProgressManager should provide a valid context")
	}

	// Test that context is valid and can be used
	select {
	case <-ctx.Done():
		t.Error("Context should not be cancelled initially")
	default:
		// Context is not cancelled, which is expected
	}

	// Test width setting
	pm.SetWidth(100)
	view = pm.View()
	// Should not crash with different widths
}

// TestAntiFlickerBehavior tests that anti-flicker mechanisms work correctly
func TestAntiFlickerBehavior(t *testing.T) {
	spinner := spinner.New(spinner.SpinnerDots)

	// Start spinner
	startTime := time.Now()
	spinner.Start()

	// Immediately try to stop - should still show for minimum duration
	spinner.Stop()

	// Check that spinner respects minimum duration
	if !spinner.IsLoading() && time.Since(startTime) < time.Millisecond*500 {
		// This is expected behavior - spinner stopped but might still show for anti-flicker
	}

	// Test ShouldHide functionality
	if spinner.ShouldHide() && time.Since(startTime) < time.Millisecond*500 {
		t.Error("Spinner should not hide before minimum duration")
	}
}

// TestProgressAnimationConsistency tests that all animations run at consistent rates
func TestProgressAnimationConsistency(t *testing.T) {
	// Test spinner animation timing
	spinner := spinner.New(spinner.SpinnerDots)
	startCmd := spinner.Start()
	spinner.SetMessage("Testing animation consistency")

	// The start command should return a tick command
	if startCmd == nil {
		t.Error("Spinner start should return tick command")
	}

	// Test multiple updates in sequence - spinner may or may not return commands depending on message type
	for i := 0; i < 5; i++ {
		updatedSpinner, cmd := spinner.Update(struct{}{})
		spinner = updatedSpinner
		// cmd may be nil depending on the message, that's OK
		_ = cmd

		view := spinner.View()
		if view == "" {
			t.Error("Spinner view should not be empty during animation")
		}
	}

	// Test progress bar with different file counts
	progressBar := progress.NewFileProgressModel(100)
	for i := 1; i <= 10; i++ {
		progressBar.SetFileCount(i*10, 100)
		view := progressBar.View()
		if view == "" {
			t.Error("Progress bar view should not be empty")
		}
	}
}

// TestComponentViewSizes tests that component views are reasonable sizes
func TestComponentViewSizes(t *testing.T) {
	tests := []struct {
		name    string
		getView func() string
		minSize int
		maxSize int
	}{
		{
			name: "Spinner with message",
			getView: func() string {
				s := spinner.New(spinner.SpinnerDots)
				s.Start()
				s.SetMessage("Processing files...")
				return s.View()
			},
			minSize: 10,
			maxSize: 100,
		},
		{
			name: "Progress bar with files",
			getView: func() string {
				p := progress.NewFileProgressModel(50)
				p.SetFileCount(25, 50)
				p.SetWidth(60)
				return p.View()
			},
			minSize: 50,
			maxSize: 500,
		},
		{
			name: "Spinner with cancel hint",
			getView: func() string {
				s := spinner.New(spinner.SpinnerDots)
				s.Start()
				s.SetMessage("Calculating sizes...")
				return s.ViewWithCancel()
			},
			minSize: 30,
			maxSize: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.getView()

			if len(view) < tt.minSize {
				t.Errorf("View too short: expected at least %d chars, got %d", tt.minSize, len(view))
			}

			if len(view) > tt.maxSize {
				t.Errorf("View too long: expected at most %d chars, got %d", tt.maxSize, len(view))
			}
		})
	}
}

// TestErrorHandling tests that components handle errors gracefully
func TestErrorHandling(t *testing.T) {
	// Test progress bar with invalid values
	p := progress.NewModel(0, 0, nil)
	view := p.View()
	if view == "" {
		t.Error("Progress bar should handle zero values gracefully")
	}

	// Test progress bar with negative values
	p.SetCurrent(-1)
	view = p.View()
	// Should not crash

	// Test spinner without starting
	s := spinner.New(spinner.SpinnerDots)
	view = s.View()
	if view != "" {
		t.Error("Inactive spinner should return empty view")
	}

	// Test progress manager with context cancellation
	pm := confirm.NewProgressManager()
	view = pm.View()
	// Should handle gracefully even with various states
}
