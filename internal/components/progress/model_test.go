package progress

import (
	"strings"
	"testing"
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
	
	// Invalid current step (too low)
	model.SetCurrent(0)
	if model.current != 3 {
		t.Errorf("Expected current to remain 3, got %d", model.current)
	}
	
	// Invalid current step (too high)
	model.SetCurrent(6)
	if model.current != 3 {
		t.Errorf("Expected current to remain 3, got %d", model.current)
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
		{1, 5, 0.0},   // First step
		{3, 5, 50.0},  // Middle step
		{5, 5, 100.0}, // Last step
		{1, 1, 100.0}, // Single step
		{2, 3, 50.0},  // Two of three
	}
	
	for _, tt := range tests {
		model := NewModel(tt.current, tt.total, []string{})
		result := model.GetProgressPercent()
		
		if result != tt.expected {
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
	
	// Should contain step counter
	if !strings.Contains(progress, "Step 2 of 4") {
		t.Error("Expected progress to contain step counter")
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
		{1, 5, false}, // 0% progress (step 1 of 5 = 0/4 progress)
		{3, 5, false}, // 50% progress (step 3 of 5 = 2/4 progress)
		{5, 5, true},  // 100% progress (step 5 of 5 = 4/4 progress)
		{1, 1, true},  // 100% progress (single step)
	}
	
	for _, tt := range tests {
		model := NewModel(tt.current, tt.total, []string{})
		model.SetWidth(50)
		
		progressBar := model.renderProgressBar()
		
		if tt.expectFull {
			if !strings.Contains(progressBar, "100%") {
				t.Errorf("Expected 100%% for step %d of %d, got: %s", 
					tt.current, tt.total, progressBar)
			}
		} else {
			if strings.Contains(progressBar, "100%") {
				t.Errorf("Expected less than 100%% for step %d of %d, got: %s", 
					tt.current, tt.total, progressBar)
			}
		}
	}
}