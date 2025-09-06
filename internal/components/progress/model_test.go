package progress

import (
	"strings"
	"testing"
	"time"
)

func TestNewModel(t *testing.T) {
	titles := []string{"Step 1", "Step 2", "Step 3"}
	model := NewModel(2, 3, titles)

	if model.current != 2 {
		t.Errorf("Expected current = 2, got %d", model.current)
	}

	if model.total != 3 {
		t.Errorf("Expected total = 3, got %d", model.total)
	}

	if model.width != 80 {
		t.Errorf("Expected default width = 80, got %d", model.width)
	}

	if len(model.titles) != 3 {
		t.Errorf("Expected 3 titles, got %d", len(model.titles))
	}

	if model.titles[0] != "Step 1" {
		t.Errorf("Expected first title 'Step 1', got '%s'", model.titles[0])
	}

	if !model.showETA {
		t.Error("Expected showETA to be true by default")
	}
}

func TestNewFileProgressModel(t *testing.T) {
	model := NewFileProgressModel(100)

	if model.totalFiles != 100 {
		t.Errorf("Expected totalFiles to be 100, got %d", model.totalFiles)
	}
	if model.fileCount != 0 {
		t.Errorf("Expected fileCount to be 0, got %d", model.fileCount)
	}
	if model.current != 0 {
		t.Errorf("Expected current to be 0, got %d", model.current)
	}
	if !model.showETA {
		t.Error("Expected showETA to be true")
	}
}

func TestNewBytesProgressModel(t *testing.T) {
	totalBytes := int64(1024 * 1024) // 1MB
	model := NewBytesProgressModel(totalBytes)

	if model.totalBytes != totalBytes {
		t.Errorf("Expected totalBytes to be %d, got %d", totalBytes, model.totalBytes)
	}
	if model.bytesRead != 0 {
		t.Errorf("Expected bytesRead to be 0, got %d", model.bytesRead)
	}
	if model.total != 100 {
		t.Errorf("Expected total to be 100 (percentage), got %d", model.total)
	}
}

func TestSetWidth(t *testing.T) {
	model := NewModel(1, 3, []string{})
	model.SetWidth(120)

	if model.width != 120 {
		t.Errorf("Expected width = 120, got %d", model.width)
	}
}

func TestSetCurrent(t *testing.T) {
	model := NewModel(1, 5, []string{})

	// Valid current step
	model.SetCurrent(3)
	if model.current != 3 {
		t.Errorf("Expected current = 3, got %d", model.current)
	}

	// Check percentage is calculated
	expectedPercentage := 60.0 // 3/5 * 100
	if model.percentage != expectedPercentage {
		t.Errorf("Expected percentage to be %.1f, got %.1f", expectedPercentage, model.percentage)
	}

	// Test allowing 0 now
	model.SetCurrent(0)
	if model.current != 0 {
		t.Errorf("Expected current = 0, got %d", model.current)
	}

	// Invalid current step (too high)
	model.SetCurrent(6)
	if model.current != 0 {
		t.Errorf("Expected current to remain 0, got %d", model.current)
	}

	// Valid boundary values
	model.SetCurrent(1)
	if model.current != 1 {
		t.Errorf("Expected current = 1, got %d", model.current)
	}

	model.SetCurrent(5)
	if model.current != 5 {
		t.Errorf("Expected current = 5, got %d", model.current)
	}
}

func TestSetBytes(t *testing.T) {
	model := NewBytesProgressModel(1000)

	model.SetBytes(250, 1000)
	if model.bytesRead != 250 {
		t.Errorf("Expected bytesRead to be 250, got %d", model.bytesRead)
	}

	expectedPercentage := 25.0
	if model.percentage != expectedPercentage {
		t.Errorf("Expected percentage to be %.1f, got %.1f", expectedPercentage, model.percentage)
	}
}

func TestSetFileCount(t *testing.T) {
	model := NewFileProgressModel(10)

	model.SetFileCount(3, 10)
	if model.fileCount != 3 {
		t.Errorf("Expected fileCount to be 3, got %d", model.fileCount)
	}

	expectedPercentage := 30.0
	if model.percentage != expectedPercentage {
		t.Errorf("Expected percentage to be %.1f, got %.1f", expectedPercentage, model.percentage)
	}
}

func TestSetMessage(t *testing.T) {
	model := NewModel(1, 3, nil)
	message := "Processing files..."

	model.SetMessage(message)
	if model.message != message {
		t.Errorf("Expected message to be '%s', got '%s'", message, model.message)
	}
}

func TestIncrementFile(t *testing.T) {
	model := NewFileProgressModel(5)

	model.IncrementFile()
	if model.fileCount != 1 {
		t.Errorf("Expected fileCount to be 1, got %d", model.fileCount)
	}

	model.IncrementFile()
	if model.fileCount != 2 {
		t.Errorf("Expected fileCount to be 2, got %d", model.fileCount)
	}

	expectedPercentage := 40.0 // 2/5 * 100
	if model.percentage != expectedPercentage {
		t.Errorf("Expected percentage to be %.1f, got %.1f", expectedPercentage, model.percentage)
	}
}

func TestAddBytes(t *testing.T) {
	model := NewBytesProgressModel(1000)

	model.AddBytes(100)
	if model.bytesRead != 100 {
		t.Errorf("Expected bytesRead to be 100, got %d", model.bytesRead)
	}

	model.AddBytes(150)
	if model.bytesRead != 250 {
		t.Errorf("Expected bytesRead to be 250, got %d", model.bytesRead)
	}
}

func TestGetETA(t *testing.T) {
	model := NewModel(0, 100, nil)
	model.startTime = time.Now().Add(-10 * time.Second)
	model.SetCurrent(25) // 25% complete after 10 seconds

	eta := model.GetETA()
	// Should be approximately 30 seconds (75% remaining at current rate)
	expectedRange := 25 * time.Second
	if eta < expectedRange || eta > 35*time.Second {
		t.Errorf("Expected ETA to be around 30s, got %v", eta)
	}

	// Test edge cases
	model.percentage = 0
	eta = model.GetETA()
	if eta != 0 {
		t.Error("ETA should be 0 for 0% progress")
	}

	model.percentage = 100
	eta = model.GetETA()
	if eta != 0 {
		t.Error("ETA should be 0 for 100% progress")
	}
}

func TestViewWithMessage(t *testing.T) {
	model := NewModel(1, 3, nil)
	customMessage := "Custom progress message"

	view := model.ViewWithMessage(customMessage)
	if !strings.Contains(view, customMessage) {
		t.Error("ViewWithMessage should contain the custom message")
	}

	// Original message should be unchanged
	if model.message == customMessage {
		t.Error("ViewWithMessage should not modify the original message")
	}
}

func TestViewCompact(t *testing.T) {
	// Test with file progress
	model := NewFileProgressModel(10)
	model.SetFileCount(3, 10)

	view := model.ViewCompact()
	if !strings.Contains(view, "3/10 files") {
		t.Error("ViewCompact should show file count")
	}

	if !strings.Contains(view, "30.0%") {
		t.Error("ViewCompact should show percentage")
	}

	// Test with percentage only
	model2 := NewModel(5, 10, nil)
	model2.SetCurrent(5) // Set to get 50% progress
	view2 := model2.ViewCompact()
	if !strings.Contains(view2, "50.0%") {
		t.Error("ViewCompact should show percentage for non-file progress")
	}
}

func TestGetCurrent(t *testing.T) {
	model := NewModel(3, 5, []string{})

	if model.GetCurrent() != 3 {
		t.Errorf("Expected GetCurrent() = 3, got %d", model.GetCurrent())
	}
}

func TestGetTotal(t *testing.T) {
	model := NewModel(2, 7, []string{})

	if model.GetTotal() != 7 {
		t.Errorf("Expected GetTotal() = 7, got %d", model.GetTotal())
	}
}

func TestIsComplete(t *testing.T) {
	model := NewModel(3, 5, []string{})

	// Not complete
	if model.IsComplete() {
		t.Error("Expected IsComplete() = false for step 3 of 5")
	}

	// Complete
	model.SetCurrent(5)
	if !model.IsComplete() {
		t.Error("Expected IsComplete() = true for step 5 of 5")
	}

	// Edge case - single step
	singleModel := NewModel(1, 1, []string{})
	if !singleModel.IsComplete() {
		t.Error("Expected IsComplete() = true for step 1 of 1")
	}
}

func TestGetProgressPercent(t *testing.T) {
	tests := []struct {
		current  int
		total    int
		expected float64
	}{
		{1, 5, 20.0},  // First step: 1/5 = 20%
		{3, 5, 60.0},  // Middle step: 3/5 = 60%
		{5, 5, 100.0}, // Last step: 5/5 = 100%
		{1, 1, 100.0}, // Single step: 1/1 = 100%
		{2, 3, 66.7},  // Two of three: 2/3 = 66.7%
	}

	for _, tt := range tests {
		model := NewModel(tt.current, tt.total, []string{})
		result := model.GetProgressPercent()

		if result < tt.expected-0.1 || result > tt.expected+0.1 {
			t.Errorf("GetProgressPercent() for step %d of %d = %.1f, want %.1f",
				tt.current, tt.total, result, tt.expected)
		}
	}
}

func TestView(t *testing.T) {
	titles := []string{"Select Files", "Choose Template", "Enter Task"}
	model := NewModel(2, 3, titles)

	view := model.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain step counter
	if !strings.Contains(view, "Step 2 of 3") {
		t.Error("Expected view to contain step counter")
	}

	// Should contain current step title
	if !strings.Contains(view, "Choose Template") {
		t.Error("Expected view to contain current step title")
	}

	// Should contain progress elements (may be styled, so check for basic characters)
	if !strings.Contains(view, "[") || !strings.Contains(view, "]") {
		t.Error("Expected view to contain progress bar brackets")
	}
}

func TestRenderProgressBar(t *testing.T) {
	model := NewModel(2, 5, []string{})
	model.SetWidth(50) // Set a specific width for consistent testing

	progressBar := model.renderProgressBar()

	if progressBar == "" {
		t.Error("Expected non-empty progress bar")
	}

	// Should contain brackets
	if !strings.Contains(progressBar, "[") || !strings.Contains(progressBar, "]") {
		t.Error("Expected progress bar to contain brackets")
	}

	// Should contain percentage
	if !strings.Contains(progressBar, "%") {
		t.Error("Expected progress bar to contain percentage")
	}
}

func TestRenderProgressBar_SmallWidth(t *testing.T) {
	model := NewModel(1, 3, []string{})
	model.SetWidth(10) // Very small width

	progressBar := model.renderProgressBar()

	// Should return empty string for very small width
	if progressBar != "" {
		t.Errorf("Expected empty progress bar for small width, got '%s'", progressBar)
	}
}

func TestRenderStepIndicators(t *testing.T) {
	model := NewModel(3, 5, []string{})

	indicators := model.renderStepIndicators()

	if indicators == "" {
		t.Error("Expected non-empty step indicators")
	}

	// Should contain arrows for multi-step
	if !strings.Contains(indicators, "→") {
		t.Error("Expected step indicators to contain arrows")
	}

	// Should contain numbers or checkmarks
	if !strings.Contains(indicators, "3") { // Current step should be visible
		t.Error("Expected step indicators to contain current step number")
	}
}

func TestRenderStepIndicators_SingleStep(t *testing.T) {
	model := NewModel(1, 1, []string{})

	indicators := model.renderStepIndicators()

	if indicators == "" {
		t.Error("Expected non-empty step indicators for single step")
	}

	// Should not contain arrows for single step
	if strings.Contains(indicators, "→") {
		t.Error("Expected no arrows for single step indicator")
	}
}

func TestRenderProgress_WithoutTitles(t *testing.T) {
	model := NewModel(2, 4, []string{}) // No titles

	progress := model.renderProgress()

	if progress == "" {
		t.Error("Expected non-empty progress render")
	}

	// Should contain progress bar, but may not have step counter with new enhanced logic
	if !strings.Contains(progress, "[") {
		t.Error("Expected progress to contain progress bar")
	}

	// Should not crash with empty titles
}

func TestRenderProgress_WithTitles(t *testing.T) {
	titles := []string{"First", "Second", "Third"}
	model := NewModel(2, 3, titles)

	progress := model.renderProgress()

	if progress == "" {
		t.Error("Expected non-empty progress render")
	}

	// Should contain step counter
	if !strings.Contains(progress, "Step 2 of 3") {
		t.Error("Expected progress to contain step counter")
	}

	// Should contain current title
	if !strings.Contains(progress, "Second") {
		t.Error("Expected progress to contain current step title")
	}
}

func TestRenderProgress_TitleOutOfBounds(t *testing.T) {
	titles := []string{"First", "Second"} // Only 2 titles
	model := NewModel(3, 4, titles)       // But asking for step 3

	progress := model.renderProgress()

	if progress == "" {
		t.Error("Expected non-empty progress render")
	}

	// Should contain step counter
	if !strings.Contains(progress, "Step 3 of 4") {
		t.Error("Expected progress to contain step counter")
	}

	// Should not crash with out-of-bounds title access
	// (No specific title should be shown)
}

func TestProgressBarCalculations(t *testing.T) {
	tests := []struct {
		current    int
		total      int
		expectFull bool
	}{
		{1, 5, false}, // 20% progress (step 1 of 5)
		{3, 5, false}, // 60% progress (step 3 of 5)
		{5, 5, true},  // 100% progress (step 5 of 5)
		{1, 1, true},  // 100% progress (single step)
	}

	for _, tt := range tests {
		model := NewModel(tt.current, tt.total, []string{})
		model.SetWidth(50)

		progressBar := model.renderProgressBar()

		if tt.expectFull {
			if !strings.Contains(progressBar, "100.0%") {
				t.Errorf("Expected 100.0%% for step %d of %d, got: %s",
					tt.current, tt.total, progressBar)
			}
		} else {
			if strings.Contains(progressBar, "100.0%") {
				t.Errorf("Expected less than 100.0%% for step %d of %d, got: %s",
					tt.current, tt.total, progressBar)
			}
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		result := formatBytes(tt.input)
		if result != tt.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m30s"},
		{3661 * time.Second, "1h01m"},
	}

	for _, tt := range tests {
		result := formatDuration(tt.input)
		if result != tt.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestUpdate(t *testing.T) {
	model := NewModel(1, 5, nil)

	newModel, cmd := model.Update(nil)
	if cmd != nil {
		t.Error("Update should return nil command for basic progress bar")
	}

	if newModel.current != model.current {
		t.Error("Update should preserve model state")
	}
}

// Benchmark tests
func BenchmarkProgressView(b *testing.B) {
	model := NewFileProgressModel(1000)
	model.SetFileCount(500, 1000)
	model.SetMessage("Processing files...")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkProgressUpdate(b *testing.B) {
	model := NewFileProgressModel(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model, _ = model.Update(nil)
	}
}
