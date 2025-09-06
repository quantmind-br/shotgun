package confirm

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewProgressManager(t *testing.T) {
	pm := NewProgressManager()

	if pm == nil {
		t.Fatal("NewProgressManager should not return nil")
	}

	// Test initial state
	state := pm.GetState()
	if state.Processed != 0 {
		t.Error("Initial processed count should be 0")
	}

	if state.Total != 0 {
		t.Error("Initial total count should be 0")
	}

	if state.Completed {
		t.Error("Initial state should not be completed")
	}
}

func TestStartProgress(t *testing.T) {
	pm := NewProgressManager()

	totalFiles := 50
	pm.StartProgress(totalFiles)

	state := pm.GetState()
	if state.Total != totalFiles {
		t.Errorf("Expected total to be %d, got %d", totalFiles, state.Total)
	}

	if state.Processed != 0 {
		t.Error("Processed count should start at 0")
	}

	if state.Percentage != 0.0 {
		t.Error("Percentage should start at 0.0")
	}

	if state.Completed {
		t.Error("Should not be completed initially")
	}
}

func TestUpdateProgress(t *testing.T) {
	pm := NewProgressManager()
	pm.StartProgress(10)

	// Update progress
	currentFile := "test.go"
	processed := 5

	cmd := pm.UpdateProgress(processed, currentFile)
	if cmd == nil {
		t.Error("UpdateProgress should return a command")
	}

	// Check state was updated
	state := pm.GetState()
	if state.Processed != processed {
		t.Errorf("Expected processed to be %d, got %d", processed, state.Processed)
	}

	if state.CurrentFile != currentFile {
		t.Errorf("Expected current file to be '%s', got '%s'", currentFile, state.CurrentFile)
	}

	expectedPercentage := 0.5 // 5/10
	if state.Percentage != expectedPercentage {
		t.Errorf("Expected percentage to be %.1f, got %.1f", expectedPercentage, state.Percentage)
	}
}

func TestCompleteProgress(t *testing.T) {
	pm := NewProgressManager()
	pm.StartProgress(10)

	// Complete progress
	cmd := pm.CompleteProgress()
	if cmd == nil {
		t.Error("CompleteProgress should return a command")
	}

	// Check state
	state := pm.GetState()
	if !state.Completed {
		t.Error("State should be completed")
	}

	if state.Percentage != 1.0 {
		t.Error("Percentage should be 100% when completed")
	}
}

func TestCancelProgress(t *testing.T) {
	pm := NewProgressManager()
	pm.StartProgress(10)

	// Cancel progress
	cmd := pm.CancelProgress()
	if cmd == nil {
		t.Error("CancelProgress should return a command")
	}

	// Context should be cancelled
	ctx := pm.GetContext()
	if ctx == nil {
		t.Error("Context should not be nil")
	}
}

func TestProgressManagerUpdate(t *testing.T) {
	pm := NewProgressManager()
	pm.StartProgress(5)

	// Test ProgressMsg
	progressMsg := ProgressMsg{
		Processed:   3,
		Total:       5,
		CurrentFile: "file.txt",
		Percentage:  0.6,
		Completed:   false,
	}

	cmd := pm.Update(progressMsg)
	// Command could be nil or not, both are valid
	_ = cmd

	// Check state was updated
	state := pm.GetState()
	if state.Processed != 3 {
		t.Errorf("Expected processed to be 3, got %d", state.Processed)
	}

	if state.CurrentFile != "file.txt" {
		t.Errorf("Expected current file to be 'file.txt', got '%s'", state.CurrentFile)
	}
}

func TestProgressManagerUpdateWithKeyMsg(t *testing.T) {
	pm := NewProgressManager()
	pm.StartProgress(10)

	// Test ESC key
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	cmd := pm.Update(escMsg)

	if cmd == nil {
		t.Error("ESC key should trigger cancel command")
	}

	// Test 'q' key
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	cmd = pm.Update(qMsg)

	if cmd == nil {
		t.Error("'q' key should trigger cancel command")
	}
}

func TestProgressManagerView(t *testing.T) {
	pm := NewProgressManager()

	// Test indeterminate progress (total = 0)
	view := pm.View()
	// The view contains styled content, so we just verify it's not empty
	// and that the spinner is working
	if view == "" {
		t.Error("Indeterminate progress view should not be empty")
	}

	// Test determinate progress
	pm.StartProgress(10)
	pm.UpdateProgress(3, "example.go")

	view = pm.View()
	// The view will contain our enhanced progress bar output
	// We can't easily test the exact output since it's styled,
	// but we can verify it's not empty
	if view == "" {
		t.Error("Determinate progress view should not be empty")
	}
}

func TestProgressManagerState(t *testing.T) {
	pm := NewProgressManager()

	// Test IsCompleted
	if pm.IsCompleted() {
		t.Error("Should not be completed initially")
	}

	// Test GetPercentage
	if pm.GetPercentage() != 0.0 {
		t.Error("Initial percentage should be 0.0")
	}

	// Start progress and update
	pm.StartProgress(4)
	pm.UpdateProgress(2, "test.go")

	if pm.GetPercentage() != 0.5 {
		t.Error("Percentage should be 0.5 after processing 2 of 4")
	}

	// Complete progress
	pm.CompleteProgress()

	if !pm.IsCompleted() {
		t.Error("Should be completed after CompleteProgress()")
	}

	if pm.GetPercentage() != 1.0 {
		t.Error("Percentage should be 1.0 when completed")
	}
}

func TestProgressManagerSetWidth(t *testing.T) {
	pm := NewProgressManager()

	// SetWidth should not panic
	pm.SetWidth(80)
	pm.SetWidth(20)
	pm.SetWidth(100)

	// We can't directly test the width was set since our progress component
	// doesn't expose it, but we can verify the method doesn't crash
}

func TestProgressManagerGetContext(t *testing.T) {
	pm := NewProgressManager()

	// Initially should return background context
	ctx := pm.GetContext()
	if ctx == nil {
		t.Error("GetContext should never return nil")
	}

	// After starting progress, should return the cancellable context
	pm.StartProgress(10)
	ctx2 := pm.GetContext()
	if ctx2 == nil {
		t.Error("GetContext should never return nil after StartProgress")
	}
}

func TestProgressTickCmd(t *testing.T) {
	cmd := ProgressTickCmd()
	if cmd == nil {
		t.Error("ProgressTickCmd should return a command")
	}

	// Execute the command to get a message
	msg := cmd()
	// Should return some kind of message (likely empty struct{})
	if msg == nil {
		t.Error("ProgressTickCmd should return a non-nil message")
	}
}

func TestSizeCalculationMessages(t *testing.T) {
	// Test SizeCalculationCompleteMsg
	completeMsg := SizeCalculationCompleteMsg{
		TotalSize: 1024,
		Breakdown: SizeBreakdown{
			TemplateSize:    200,
			FileContentSize: 600,
			TreeStructSize:  124,
			OverheadSize:    100,
		},
	}

	if completeMsg.TotalSize != 1024 {
		t.Error("SizeCalculationCompleteMsg should preserve TotalSize")
	}

	if completeMsg.Breakdown.TemplateSize != 200 {
		t.Error("SizeCalculationCompleteMsg should preserve breakdown")
	}

	// Test CancellationMsg
	cancelMsg := CancellationMsg{}
	_ = cancelMsg // Just verify it can be created
}

func TestProgressStateConversion(t *testing.T) {
	// Test that ProgressState can be converted to/from ProgressMsg
	originalState := ProgressState{
		Processed:   7,
		Total:       10,
		CurrentFile: "main.go",
		Percentage:  0.7,
		Completed:   false,
	}

	// Convert to ProgressMsg and back
	msg := ProgressMsg(originalState)
	convertedState := ProgressState(msg)

	if convertedState.Processed != originalState.Processed {
		t.Error("ProgressState conversion should preserve Processed")
	}

	if convertedState.Total != originalState.Total {
		t.Error("ProgressState conversion should preserve Total")
	}

	if convertedState.CurrentFile != originalState.CurrentFile {
		t.Error("ProgressState conversion should preserve CurrentFile")
	}

	if convertedState.Percentage != originalState.Percentage {
		t.Error("ProgressState conversion should preserve Percentage")
	}

	if convertedState.Completed != originalState.Completed {
		t.Error("ProgressState conversion should preserve Completed")
	}
}
