package input

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestTaskInputModel_ViewBasic(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	view := model.View()

	if view == "" {
		t.Error("expected non-empty view")
	}

	// Check for key components
	if !strings.Contains(view, "Task Description") {
		t.Error("expected view to contain header 'Task Description'")
	}

	if !strings.Contains(view, "Lines:") {
		t.Error("expected view to contain line count")
	}

	if !strings.Contains(view, "Characters:") {
		t.Error("expected view to contain character count")
	}

	if !strings.Contains(view, "Ctrl+Enter") {
		t.Error("expected view to contain keyboard shortcuts")
	}
}

func TestTaskInputModel_ViewWithContent(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)
	model.SetContent("This is a test task description with some content")

	view := model.View()

	// Check that counters are updated
	if !strings.Contains(view, "Lines: 1") {
		t.Error("expected view to show line count of 1")
	}

	// Check that character count is displayed (should be 49 characters)
	if !strings.Contains(view, "Characters: 49") {
		t.Error("expected view to show character count of 49")
	}
}

func TestTaskInputModel_ViewWithMultilineContent(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)
	model.SetContent("Line 1\nLine 2\nLine 3")

	view := model.View()

	// Check that line count is correct
	if !strings.Contains(view, "Lines: 3") {
		t.Error("expected view to show line count of 3")
	}

	// Character count should be 20 (including newlines)
	if !strings.Contains(view, "Characters: 20") {
		t.Error("expected view to show character count of 20")
	}
}

func TestTaskInputModel_ViewWithError(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)
	testError := errors.New("validation error")
	model.SetError(testError)

	view := model.View()

	// Should render error state
	if !strings.Contains(view, "Error:") {
		t.Error("expected view to contain error message")
	}

	if !strings.Contains(view, "validation error") {
		t.Error("expected view to contain specific error text")
	}
}

func TestTaskInputModel_ViewSmallWindow(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(0, 0)

	view := model.View()

	if view != "Loading..." {
		t.Errorf("expected 'Loading...' for zero dimensions, got %q", view)
	}
}

func TestTaskInputModel_ViewKeyboardShortcuts(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	view := model.View()

	shortcuts := []string{
		"Ctrl+Enter",
		"F3:",
		"F2:",
		"Ctrl+C:",
		"Ctrl+V:",
	}

	for _, shortcut := range shortcuts {
		if !strings.Contains(view, shortcut) {
			t.Errorf("expected view to contain keyboard shortcut %q", shortcut)
		}
	}
}

func TestTaskInputModel_ViewUTF8Content(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)
	model.SetContent("Hello 世界 🌍 café")

	view := model.View()

	// Should show correct character count for UTF-8
	expectedCharCount := len([]rune("Hello 世界 🌍 café"))

	// Check if character counting works correctly
	if !strings.Contains(view, "Characters:") {
		t.Error("expected view to contain character count")
	}

	// The exact number might vary due to how we display it,
	// but it should not be the byte count
	if model.charCount != expectedCharCount {
		t.Errorf("expected char count %d, got %d", expectedCharCount, model.charCount)
	}
}

func TestTaskInputModel_ViewConsistentLayout(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test with empty content
	viewEmpty := model.View()

	// Test with content
	model.SetContent("Some content")
	viewWithContent := model.View()

	// Both should have the same basic structure
	checkStructure := func(view, label string) {
		if !strings.Contains(view, "📝 Task Description") {
			t.Errorf("%s: expected header", label)
		}
		if !strings.Contains(view, "Lines:") {
			t.Errorf("%s: expected line count", label)
		}
		if !strings.Contains(view, "Characters:") {
			t.Errorf("%s: expected character count", label)
		}
		if !strings.Contains(view, "Ctrl+Enter") {
			t.Errorf("%s: expected keyboard shortcuts", label)
		}
	}

	checkStructure(viewEmpty, "Empty view")
	checkStructure(viewWithContent, "Content view")
}

func TestTaskInputModel_ViewResponsiveSize(t *testing.T) {
	model := NewTaskInputModel()

	// Test different window sizes
	sizes := []struct {
		width, height int
		name          string
	}{
		{40, 15, "small"},
		{80, 25, "medium"},
		{120, 40, "large"},
	}

	for _, size := range sizes {
		t.Run(size.name, func(t *testing.T) {
			model.UpdateSize(size.width, size.height)
			view := model.View()

			if view == "" {
				t.Errorf("expected non-empty view for %s window", size.name)
			}

			if view == "Loading..." {
				t.Errorf("unexpected loading state for %s window (%dx%d)", size.name, size.width, size.height)
			}
		})
	}
}

func TestTaskInputModel_ViewErrorState(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test with different error types
	testErrors := []error{
		errors.New("empty content"),
		errors.New("content too short"),
		errors.New("clipboard error"),
	}

	for _, testErr := range testErrors {
		model.SetError(testErr)
		view := model.View()

		if !strings.Contains(view, "Error:") {
			t.Errorf("expected error view to contain 'Error:', got %v", view)
		}

		if !strings.Contains(view, testErr.Error()) {
			t.Errorf("expected error view to contain error message %q", testErr.Error())
		}
	}

	// Test clearing error
	model.SetError(nil)
	view := model.View()

	if strings.Contains(view, "Error:") {
		t.Error("expected error to be cleared from view")
	}
}

func TestTaskInputModel_ViewInstructions(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	view := model.View()

	// Check for instructional text
	expectedInstructions := []string{
		"Describe your task in detail",
		"LLM will use this context",
	}

	for _, instruction := range expectedInstructions {
		if !strings.Contains(view, instruction) {
			t.Errorf("expected view to contain instruction %q", instruction)
		}
	}
}

// Additional view tests for missing coverage

func TestTaskInputModel_ViewTextareaFocus(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test focused vs unfocused textarea styling
	model.textarea.Focus()
	focusedView := model.View()

	model.textarea.Blur()
	blurredView := model.View()

	// Both views should render without error
	if focusedView == "" {
		t.Error("expected focused view to render")
	}
	if blurredView == "" {
		t.Error("expected blurred view to render")
	}

	// Views should be different when focus changes
	// (We can't test exact styling differences easily, but we can ensure both render)
}

func TestTaskInputModel_ViewLargeContent(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test with content that might affect layout
	largeContent := strings.Repeat("This is a long line that might cause wrapping issues. ", 50)
	model.SetContent(largeContent)

	view := model.View()

	// Should handle large content gracefully
	if view == "" {
		t.Error("expected view to render with large content")
	}

	// Character count should be displayed correctly
	expectedCharCount := len([]rune(largeContent))
	countText := fmt.Sprintf("Characters: %d", expectedCharCount)
	if !strings.Contains(view, countText) {
		t.Errorf("expected view to contain character count %q", countText)
	}
}

func TestTaskInputModel_ViewManyLines(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test with many lines
	manyLines := strings.Repeat("Line\n", 50)
	model.SetContent(manyLines)

	view := model.View()

	// Should handle many lines gracefully
	if view == "" {
		t.Error("expected view to render with many lines")
	}

	// Line count should be displayed correctly
	expectedLineCount := 51 // 50 lines + 1 empty line at end
	countText := fmt.Sprintf("Lines: %d", expectedLineCount)
	if !strings.Contains(view, countText) {
		t.Errorf("expected view to contain line count %q", countText)
	}
}

func TestTaskInputModel_ViewSpecialCharacters(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test with special characters that might affect rendering
	specialContent := "Special chars: \t\n\r\"'`<>&"
	model.SetContent(specialContent)

	view := model.View()

	// Should render without crashing
	if view == "" {
		t.Error("expected view to render with special characters")
	}
}

func TestTaskInputModel_ViewErrorStateEdgeCases(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test error view with various error conditions
	testErrors := []error{
		errors.New(""),
		errors.New("short error"),
		errors.New("very long error message that might cause display issues when rendered"),
		errors.New("error with UTF-8: éñ🌍"),
		errors.New("error\nwith\nnewlines"),
	}

	for i, err := range testErrors {
		t.Run(fmt.Sprintf("error_%d", i), func(t *testing.T) {
			model.SetError(err)
			view := model.View()

			// Should render error view
			if view == "" {
				t.Error("expected error view to render")
			}

			// Should contain error text
			if !strings.Contains(view, "Error:") {
				t.Error("expected error view to contain 'Error:'")
			}

			// Should contain the error message (unless it's empty)
			// Note: Long error messages may be truncated or wrapped, so we check for partial content
			if err.Error() != "" {
				errorText := err.Error()
				if len(errorText) > 50 {
					// For long messages, check for first part
					errorText = errorText[:50]
				}
				// For multiline messages, check for first line
				if strings.Contains(err.Error(), "\n") {
					errorText = strings.Split(err.Error(), "\n")[0]
				}
				if !strings.Contains(view, errorText) {
					t.Errorf("expected error view to contain part of error message %q", errorText)
				}
			}
		})
	}
}

func TestTaskInputModel_ViewDimensionHandling(t *testing.T) {
	model := NewTaskInputModel()

	// Test various dimensions that might affect rendering
	testCases := []struct {
		width, height int
		description   string
		expectLoading bool
	}{
		{0, 0, "zero dimensions", true},
		{-1, -1, "negative dimensions", true},
		{1, 1, "tiny dimensions", false},
		{40, 15, "small but valid dimensions", false},
		{200, 100, "large dimensions", false},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			model.UpdateSize(tc.width, tc.height)
			view := model.View()

			if tc.expectLoading {
				if view != "Loading..." {
					t.Errorf("expected 'Loading...' for %s, got %q", tc.description, view)
				}
			} else {
				if view == "Loading..." {
					t.Errorf("expected normal view for %s, got 'Loading...'", tc.description)
				}
				if view == "" {
					t.Errorf("expected non-empty view for %s", tc.description)
				}
			}
		})
	}
}

func TestTaskInputModel_ViewContentAndErrorInteraction(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test switching between error and normal states
	model.SetContent("Some content")
	normalView := model.View()

	// Switch to error state
	model.SetError(errors.New("test error"))
	errorView := model.View()

	// Clear error state
	model.SetError(nil)
	clearedView := model.View()

	// All views should render
	if normalView == "" {
		t.Error("expected normal view to render")
	}
	if errorView == "" {
		t.Error("expected error view to render")
	}
	if clearedView == "" {
		t.Error("expected cleared view to render")
	}

	// Error view should contain error text
	if !strings.Contains(errorView, "Error:") {
		t.Error("expected error view to contain 'Error:'")
	}

	// Normal and cleared views should not contain error text
	if strings.Contains(normalView, "Error:") {
		t.Error("expected normal view to not contain 'Error:'")
	}
	if strings.Contains(clearedView, "Error:") {
		t.Error("expected cleared view to not contain 'Error:'")
	}
}

func TestTaskInputModel_ViewCounterAccuracy(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test that counters in view match internal counters
	testContent := "Hello\nWorld\nwith UTF-8: éñ🌍"
	model.SetContent(testContent)

	view := model.View()

	// Check line count display
	expectedLineCount := 3
	lineCountText := fmt.Sprintf("Lines: %d", expectedLineCount)
	if !strings.Contains(view, lineCountText) {
		t.Errorf("expected view to contain %q", lineCountText)
	}

	// Check character count display
	expectedCharCount := len([]rune(testContent))
	charCountText := fmt.Sprintf("Characters: %d", expectedCharCount)
	if !strings.Contains(view, charCountText) {
		t.Errorf("expected view to contain %q", charCountText)
	}

	// Verify internal counters match
	if model.charCount != expectedCharCount {
		t.Errorf("internal character count mismatch: expected %d, got %d", expectedCharCount, model.charCount)
	}
	if model.lineCount != expectedLineCount {
		t.Errorf("internal line count mismatch: expected %d, got %d", expectedLineCount, model.lineCount)
	}
}
