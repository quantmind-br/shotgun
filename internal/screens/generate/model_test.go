package generate

import (
	"errors"
	"testing"

	"github.com/diogopedro/shotgun/internal/core/builder"
)

func TestNewGenerateModel(t *testing.T) {
	model := NewGenerateModel()

	if model.generating {
		t.Error("New model should not be generating")
	}

	if model.completed {
		t.Error("New model should not be completed")
	}

	if model.generator == nil {
		t.Error("Generator should be initialized")
	}

	if model.fileWriter == nil {
		t.Error("FileWriter should be initialized")
	}

	if !model.showStats {
		t.Error("showStats should default to true")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	model := NewGenerateModel()

	width, height := 120, 40
	model.UpdateWindowSize(width, height)

	if model.width != width {
		t.Errorf("Width not updated. Expected: %d, Got: %d", width, model.width)
	}

	if model.height != height {
		t.Errorf("Height not updated. Expected: %d, Got: %d", height, model.height)
	}

	expectedViewportWidth := width - 4
	if model.viewport.Width != expectedViewportWidth {
		t.Errorf("Viewport width not updated correctly. Expected: %d, Got: %d", expectedViewportWidth, model.viewport.Width)
	}
}

func TestStartGeneration(t *testing.T) {
	model := NewGenerateModel()

	model.StartGeneration()

	if !model.generating {
		t.Error("Model should be generating after StartGeneration")
	}

	if model.completed {
		t.Error("Model should not be completed after StartGeneration")
	}

	if model.error != nil {
		t.Error("Error should be nil after StartGeneration")
	}

	if model.outputFile != "" {
		t.Error("Output file should be empty after StartGeneration")
	}
}

func TestCompleteGeneration_Success(t *testing.T) {
	model := NewGenerateModel()
	model.StartGeneration()

	result := &builder.GeneratedPrompt{
		Content:      "Test prompt content",
		TemplateSize: 100,
		FileCount:    3,
		TotalSize:    1500,
	}
	outputFile := "/path/to/output.md"

	model.CompleteGeneration(result, outputFile, nil)

	if model.generating {
		t.Error("Model should not be generating after completion")
	}

	if !model.completed {
		t.Error("Model should be completed after CompleteGeneration")
	}

	if model.error != nil {
		t.Error("Error should be nil for successful completion")
	}

	if model.outputFile != outputFile {
		t.Errorf("Output file mismatch. Expected: %s, Got: %s", outputFile, model.outputFile)
	}

	if model.templateSize != result.TemplateSize {
		t.Errorf("Template size mismatch. Expected: %d, Got: %d", result.TemplateSize, model.templateSize)
	}

	if model.fileCount != result.FileCount {
		t.Errorf("File count mismatch. Expected: %d, Got: %d", result.FileCount, model.fileCount)
	}

	if model.totalSize != result.TotalSize {
		t.Errorf("Total size mismatch. Expected: %d, Got: %d", result.TotalSize, model.totalSize)
	}
}

func TestCompleteGeneration_Error(t *testing.T) {
	model := NewGenerateModel()
	model.StartGeneration()

	testError := errors.New("generation failed")

	model.CompleteGeneration(nil, "", testError)

	if model.generating {
		t.Error("Model should not be generating after completion with error")
	}

	if !model.completed {
		t.Error("Model should be completed after CompleteGeneration even with error")
	}

	if model.error != testError {
		t.Errorf("Error mismatch. Expected: %v, Got: %v", testError, model.error)
	}

	if model.outputFile != "" {
		t.Error("Output file should be empty on error")
	}
}

func TestIsGenerating(t *testing.T) {
	model := NewGenerateModel()

	if model.IsGenerating() {
		t.Error("New model should not be generating")
	}

	model.StartGeneration()

	if !model.IsGenerating() {
		t.Error("Model should be generating after StartGeneration")
	}

	model.CompleteGeneration(nil, "", nil)

	if model.IsGenerating() {
		t.Error("Model should not be generating after completion")
	}
}

func TestIsCompleted(t *testing.T) {
	model := NewGenerateModel()

	if model.IsCompleted() {
		t.Error("New model should not be completed")
	}

	model.StartGeneration()

	if model.IsCompleted() {
		t.Error("Model should not be completed while generating")
	}

	model.CompleteGeneration(nil, "", nil)

	if !model.IsCompleted() {
		t.Error("Model should be completed after CompleteGeneration")
	}
}

func TestHasError(t *testing.T) {
	model := NewGenerateModel()

	if model.HasError() {
		t.Error("New model should not have error")
	}

	// Complete with success
	model.CompleteGeneration(nil, "", nil)

	if model.HasError() {
		t.Error("Model should not have error after successful completion")
	}

	// Complete with error
	testError := errors.New("test error")
	model.CompleteGeneration(nil, "", testError)

	if !model.HasError() {
		t.Error("Model should have error after error completion")
	}

	if model.GetError() != testError {
		t.Errorf("GetError() mismatch. Expected: %v, Got: %v", testError, model.GetError())
	}
}

func TestGetOutputFile(t *testing.T) {
	model := NewGenerateModel()

	if model.GetOutputFile() != "" {
		t.Error("New model should have empty output file")
	}

	testOutputFile := "/path/to/test/output.md"
	model.CompleteGeneration(nil, testOutputFile, nil)

	if model.GetOutputFile() != testOutputFile {
		t.Errorf("GetOutputFile() mismatch. Expected: %s, Got: %s", testOutputFile, model.GetOutputFile())
	}
}

func TestToggleStats(t *testing.T) {
	model := NewGenerateModel()

	initialShowStats := model.ShowingStats()

	model.ToggleStats()

	if model.ShowingStats() == initialShowStats {
		t.Error("ToggleStats should change ShowingStats value")
	}

	model.ToggleStats()

	if model.ShowingStats() != initialShowStats {
		t.Error("ToggleStats should restore original ShowingStats value")
	}
}

func TestShowingStats(t *testing.T) {
	model := NewGenerateModel()

	// Default should be true
	if !model.ShowingStats() {
		t.Error("ShowingStats should default to true")
	}

	model.ToggleStats()

	if model.ShowingStats() {
		t.Error("ShowingStats should be false after toggle")
	}
}
