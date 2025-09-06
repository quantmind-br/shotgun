package spinner

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		style SpinnerStyle
	}{
		{"dots style", SpinnerDots},
		{"line style", SpinnerLine},
		{"circle style", SpinnerCircle},
		{"unknown style defaults to dots", SpinnerStyle("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New(tt.style)

			if model.loading {
				t.Error("New spinner should not be loading initially")
			}

			if model.message != "" {
				t.Error("New spinner should have empty message initially")
			}

			if model.minDuration != MinSpinnerDuration {
				t.Errorf("Expected minDuration %v, got %v", MinSpinnerDuration, model.minDuration)
			}
		})
	}
}

func TestStartStop(t *testing.T) {
	model := New(SpinnerDots)

	// Initially not loading
	if model.IsLoading() {
		t.Error("New spinner should not be loading")
	}

	// Start loading
	cmd := model.Start()
	if !model.IsLoading() {
		t.Error("Spinner should be loading after Start()")
	}

	if cmd == nil {
		t.Error("Start() should return a tick command")
	}

	// Stop loading
	model.Stop()
	if model.IsLoading() {
		t.Error("Spinner should not be loading after Stop()")
	}
}

func TestSetMessage(t *testing.T) {
	model := New(SpinnerDots)

	testMessage := "Loading files..."
	model.SetMessage(testMessage)

	if model.message != testMessage {
		t.Errorf("Expected message '%s', got '%s'", testMessage, model.message)
	}
}

func TestView(t *testing.T) {
	model := New(SpinnerDots)

	// Not loading - should return empty
	view := model.View()
	if view != "" {
		t.Error("View should return empty string when not loading")
	}

	// Start loading
	model.Start()
	view = model.View()
	if view == "" {
		t.Error("View should return non-empty string when loading")
	}

	// With message
	model.SetMessage("Processing...")
	view = model.View()
	if !strings.Contains(view, "Processing...") {
		t.Error("View should contain the message when set")
	}

	// Stop loading
	model.Stop()
	view = model.View()
	if view != "" {
		t.Error("View should return empty string after stopping")
	}
}

func TestViewWithCancel(t *testing.T) {
	model := New(SpinnerDots)

	// Not loading - should return empty
	view := model.ViewWithCancel()
	if view != "" {
		t.Error("ViewWithCancel should return empty string when not loading")
	}

	// Start loading
	model.Start()
	model.SetMessage("Loading...")
	view = model.ViewWithCancel()

	if !strings.Contains(view, "Loading...") {
		t.Error("ViewWithCancel should contain the message")
	}

	if !strings.Contains(view, "[Press ESC to cancel]") {
		t.Error("ViewWithCancel should contain cancellation hint")
	}
}

func TestUpdate(t *testing.T) {
	model := New(SpinnerDots)

	// Not loading - should return no command
	newModel, cmd := model.Update(nil)
	if cmd != nil {
		t.Error("Update should return nil command when not loading")
	}

	// Start loading
	model.Start()

	// Update with spinner tick
	newModel, cmd = model.Update(tea.Msg(nil))
	if newModel.loading != model.loading {
		t.Error("Update should preserve loading state")
	}
}

func TestShouldHide(t *testing.T) {
	model := New(SpinnerDots)

	// Start and immediately stop
	model.Start()
	model.Stop()

	// Should not hide immediately (anti-flicker)
	if model.ShouldHide() {
		t.Error("Should not hide immediately after stopping (anti-flicker)")
	}

	// Simulate time passage
	model.started = time.Now().Add(-MinSpinnerDuration - time.Millisecond)
	if !model.ShouldHide() {
		t.Error("Should hide after minimum duration has passed")
	}
}

func TestLoadingTracker(t *testing.T) {
	tracker := NewLoadingTracker()

	// Should not hide immediately
	if tracker.ShouldHide() {
		t.Error("LoadingTracker should not hide immediately")
	}

	if tracker.HasShownMinimum() {
		t.Error("LoadingTracker should not have shown minimum initially")
	}

	// Simulate time passage
	tracker.startTime = time.Now().Add(-MinSpinnerDuration - time.Millisecond)

	if !tracker.ShouldHide() {
		t.Error("LoadingTracker should hide after minimum duration")
	}

	if !tracker.HasShownMinimum() {
		t.Error("LoadingTracker should have shown minimum after ShouldHide() call")
	}
}

func TestMinSpinnerDuration(t *testing.T) {
	expectedDuration := 500 * time.Millisecond
	if MinSpinnerDuration != expectedDuration {
		t.Errorf("Expected MinSpinnerDuration to be %v, got %v", expectedDuration, MinSpinnerDuration)
	}
}

// Benchmark tests for performance
func BenchmarkSpinnerView(b *testing.B) {
	model := New(SpinnerDots)
	model.Start()
	model.SetMessage("Loading files...")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkSpinnerUpdate(b *testing.B) {
	model := New(SpinnerDots)
	model.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model, _ = model.Update(nil)
	}
}
