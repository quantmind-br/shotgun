package confirm

import (
	"testing"

	"github.com/diogopedro/shotgun/internal/models"
)

func TestNewConfirmModel(t *testing.T) {
	model := NewConfirmModel()

	if model.ready {
		t.Error("Expected new model to not be ready initially")
	}

	if model.calculating {
		t.Error("Expected new model to not be calculating initially")
	}

	if model.estimatedSize != 0 {
		t.Error("Expected initial estimated size to be 0")
	}

	if model.progress.Width != 40 {
		t.Errorf("Expected progress width to be 40, got %d", model.progress.Width)
	}
}

func TestUpdateWindowSize(t *testing.T) {
	model := NewConfirmModel()
	
	model.UpdateWindowSize(120, 40)

	if model.width != 120 {
		t.Errorf("Expected width 120, got %d", model.width)
	}

	if model.height != 40 {
		t.Errorf("Expected height 40, got %d", model.height)
	}

	expectedViewportWidth := 116 // 120 - 4 for borders
	if model.viewport.Width != expectedViewportWidth {
		t.Errorf("Expected viewport width %d, got %d", expectedViewportWidth, model.viewport.Width)
	}

	expectedViewportHeight := 28 // 40 - 12 for borders and headers
	if model.viewport.Height != expectedViewportHeight {
		t.Errorf("Expected viewport height %d, got %d", expectedViewportHeight, model.viewport.Height)
	}
}

func TestSetData(t *testing.T) {
	model := NewConfirmModel()
	
	template := &models.Template{
		Name:        "Test Template",
		Version:     "1.0",
		Description: "Test description",
		Content:     "Test content",
	}
	
	selectedFiles := []string{"file1.go", "file2.go"}
	taskContent := "Test task content"
	rulesContent := "Test rules content"

	model.SetData(template, selectedFiles, taskContent, rulesContent)

	if !model.ready {
		t.Error("Expected model to be ready after SetData")
	}

	if !model.IsReady() {
		t.Error("Expected IsReady() to return true")
	}

	if model.template != template {
		t.Error("Template not set correctly")
	}

	if len(model.selectedFiles) != 2 {
		t.Errorf("Expected 2 selected files, got %d", len(model.selectedFiles))
	}

	if model.taskContent != taskContent {
		t.Error("Task content not set correctly")
	}

	if model.rulesContent != rulesContent {
		t.Error("Rules content not set correctly")
	}
}

func TestSetEstimatedSize(t *testing.T) {
	model := NewConfirmModel()
	model.calculating = true

	breakdown := SizeBreakdown{
		TemplateSize:    1000,
		FileContentSize: 5000,
		TreeStructSize:  500,
		OverheadSize:    200,
	}

	model.SetEstimatedSize(6700, breakdown)

	if model.calculating {
		t.Error("Expected calculating to be false after SetEstimatedSize")
	}

	if model.estimatedSize != 6700 {
		t.Errorf("Expected estimated size 6700, got %d", model.estimatedSize)
	}

	if model.sizeBreakdown != breakdown {
		t.Error("Size breakdown not set correctly")
	}
}

func TestUpdateWarningLevel(t *testing.T) {
	tests := []struct {
		name         string
		size         int64
		expectedLevel WarningLevel
		expectWarning bool
	}{
		{
			name:         "Normal size",
			size:         50 * 1024, // 50KB
			expectedLevel: WarningNone,
			expectWarning: false,
		},
		{
			name:         "Large size",
			size:         150 * 1024, // 150KB
			expectedLevel: WarningLarge,
			expectWarning: true,
		},
		{
			name:         "Very large size",
			size:         700 * 1024, // 700KB
			expectedLevel: WarningVeryLarge,
			expectWarning: true,
		},
		{
			name:         "Excessive size",
			size:         3 * 1024 * 1024, // 3MB
			expectedLevel: WarningExcessive,
			expectWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewConfirmModel()
			breakdown := SizeBreakdown{}
			
			model.SetEstimatedSize(tt.size, breakdown)

			if model.warningLevel != tt.expectedLevel {
				t.Errorf("Expected warning level %d, got %d", tt.expectedLevel, model.warningLevel)
			}

			if model.showWarning != tt.expectWarning {
				t.Errorf("Expected showWarning %v, got %v", tt.expectWarning, model.showWarning)
			}
		})
	}
}

func TestStartCalculation(t *testing.T) {
	model := NewConfirmModel()
	model.estimatedSize = 1000
	model.showWarning = true

	model.StartCalculation()

	if !model.calculating {
		t.Error("Expected calculating to be true after StartCalculation")
	}

	if model.estimatedSize != 0 {
		t.Error("Expected estimated size to be reset to 0")
	}

	if model.showWarning {
		t.Error("Expected showWarning to be false after StartCalculation")
	}
}

func TestFilenameOperations(t *testing.T) {
	model := NewConfirmModel()
	filename := "test_output.md"

	model.SetOutputFilename(filename)

	if model.GetOutputFilename() != filename {
		t.Errorf("Expected filename %s, got %s", filename, model.GetOutputFilename())
	}
}

func TestIsCalculating(t *testing.T) {
	model := NewConfirmModel()

	if model.IsCalculating() {
		t.Error("Expected IsCalculating to be false initially")
	}

	model.StartCalculation()

	if !model.IsCalculating() {
		t.Error("Expected IsCalculating to be true after StartCalculation")
	}

	breakdown := SizeBreakdown{}
	model.SetEstimatedSize(1000, breakdown)

	if model.IsCalculating() {
		t.Error("Expected IsCalculating to be false after SetEstimatedSize")
	}
}