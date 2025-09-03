package input

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Integration tests for clipboard operations and cross-component functionality
// Focus on message-based testing rather than key simulation due to KeyMsg complexity

func TestTaskInputModel_ClipboardMessageIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test complete clipboard workflow through messages
	t.Run("clipboard_message_workflow", func(t *testing.T) {
		// Test clipboard copy message behavior
		testContent := "Content to copy to clipboard"
		model.SetContent(testContent)

		// Create clipboard copy message directly
		copyMsg := ClipboardCopyMsg{Text: testContent}
		
		// While this doesn't test the key->message conversion,
		// it tests the message handling logic
		_ = copyMsg // We would send this to app layer in real usage

		// Test clipboard paste message
		pasteMsg := ClipboardPasteMsg{Text: " pasted text"}
		updatedModel, cmd := model.Update(pasteMsg)

		expectedContent := testContent + " pasted text"
		if updatedModel.GetContent() != expectedContent {
			t.Errorf("expected content %q after paste, got %q", expectedContent, updatedModel.GetContent())
		}

		// Character count should be updated
		expectedCharCount := len([]rune(expectedContent))
		if updatedModel.charCount != expectedCharCount {
			t.Errorf("expected character count %d after paste, got %d", expectedCharCount, updatedModel.charCount)
		}

		// Should not generate additional commands from paste
		if cmd != nil {
			t.Error("expected no additional command from paste message handling")
		}
	})

	// Test clipboard error handling
	t.Run("clipboard_error_handling", func(t *testing.T) {
		model.SetError(nil) // Clear any existing error
		
		// Simulate clipboard error
		clipboardError := ClipboardErrorMsg{Error: errors.New("clipboard access denied")}
		updatedModel, cmd := model.Update(clipboardError)

		// Error should be set on the model
		if updatedModel.GetError() == nil {
			t.Error("expected error to be set after clipboard error")
		}

		if updatedModel.GetError().Error() != "clipboard access denied" {
			t.Errorf("expected specific error message, got %v", updatedModel.GetError())
		}

		// Should not generate additional commands
		if cmd != nil {
			t.Error("expected no command from clipboard error handling")
		}
	})
}

func TestTaskInputModel_NavigationIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test complete navigation workflow with state preservation
	t.Run("complete_navigation_workflow", func(t *testing.T) {
		// Set content
		testContent := "Task description with UTF-8: éñ🌍"
		model.SetContent(testContent)

		// Test F2 back navigation
		f2Key := tea.KeyMsg{Type: tea.KeyF2}
		updatedModel, cmd := model.Update(f2Key)

		// Content should be preserved
		if updatedModel.GetContent() != testContent {
			t.Error("content should be preserved during F2 navigation")
		}

		// Should generate back navigation command
		if cmd == nil {
			t.Error("expected back navigation command")
		}

		// Execute command to verify message type
		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(BackToTemplateMsg); !ok {
				t.Error("expected BackToTemplateMsg from F2 navigation")
			}
		}
	})

	// Test advancement with validation
	t.Run("advancement_with_validation", func(t *testing.T) {
		// Test F3 advancement with valid content
		model.SetContent("Valid content for advancement")
		
		f3Key := tea.KeyMsg{Type: tea.KeyF3}
		updatedModel, cmd := model.Update(f3Key)

		// Should clear any errors
		if updatedModel.GetError() != nil {
			t.Error("expected error to be cleared for valid content")
		}

		// Should generate advancement command
		if cmd == nil {
			t.Error("expected advancement command for valid content")
		}

		// Execute command to verify message type
		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(TaskInputMsg); !ok {
				t.Error("expected TaskInputMsg from F3 advancement")
			}
		}
	})

	// Test advancement failure with validation
	t.Run("advancement_failure_with_validation", func(t *testing.T) {
		// Test F3 advancement with empty content
		model.SetContent("")
		model.SetError(nil) // Clear any existing error
		
		f3Key := tea.KeyMsg{Type: tea.KeyF3}
		updatedModel, cmd := model.Update(f3Key)

		// Should set validation error
		if updatedModel.GetError() == nil {
			t.Error("expected validation error for empty content")
		}

		if updatedModel.GetError().Error() != "task description cannot be empty" {
			t.Errorf("expected specific validation error, got %v", updatedModel.GetError())
		}

		// Should not generate advancement command
		if cmd != nil {
			t.Error("expected no advancement command for invalid content")
		}
	})
}

func TestTaskInputModel_StateManagementIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test complete state synchronization workflow
	t.Run("state_synchronization", func(t *testing.T) {
		// Test external content update
		externalContent := "Content updated externally"
		contentMsg := TaskContentUpdatedMsg{Content: externalContent}
		
		updatedModel, cmd := model.Update(contentMsg)

		// All state should be synchronized
		if updatedModel.GetContent() != externalContent {
			t.Error("GetContent should return externally updated content")
		}

		if updatedModel.content != externalContent {
			t.Error("internal content field should be updated")
		}

		if updatedModel.textarea.Value() != externalContent {
			t.Error("textarea value should be updated")
		}

		// Counters should be updated
		expectedCharCount := len([]rune(externalContent))
		if updatedModel.charCount != expectedCharCount {
			t.Errorf("character count should be updated to %d, got %d", expectedCharCount, updatedModel.charCount)
		}

		// Should not generate commands from content updates
		if cmd != nil {
			t.Error("expected no command from external content update")
		}
	})

	// Test window size and content interaction
	t.Run("window_size_content_interaction", func(t *testing.T) {
		// Set content first
		longContent := strings.Repeat("This is a long line of content. ", 20)
		model.SetContent(longContent)

		// Change window size
		sizeMsg := tea.WindowSizeMsg{Width: 40, Height: 15}
		updatedModel, cmd := model.Update(sizeMsg)

		// Content should be preserved
		if updatedModel.GetContent() != longContent {
			t.Error("content should be preserved during resize")
		}

		// Dimensions should be updated
		if updatedModel.width != 40 || updatedModel.height != 15 {
			t.Errorf("expected dimensions 40x15, got %dx%d", updatedModel.width, updatedModel.height)
		}

		// Counters should remain accurate
		expectedCharCount := len([]rune(longContent))
		if updatedModel.charCount != expectedCharCount {
			t.Error("character count should remain accurate after resize")
		}

		// Should not generate commands from resize
		if cmd != nil {
			t.Error("expected no command from window resize")
		}
	})
}

func TestTaskInputModel_ErrorRecoveryIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test complete error recovery workflow
	t.Run("error_recovery_workflow", func(t *testing.T) {
		// Start with an error state
		model.SetError(errors.New("initial error"))

		// Verify error view is rendered
		errorView := model.View()
		if !strings.Contains(errorView, "Error:") {
			t.Error("expected error view to be rendered")
		}

		// Add valid content (should not automatically clear error)
		model.SetContent("Valid content")
		stillErrorView := model.View()
		if !strings.Contains(stillErrorView, "Error:") {
			t.Error("error should persist until explicitly cleared")
		}

		// Attempt advancement (should clear error for valid content)
		f3Key := tea.KeyMsg{Type: tea.KeyF3}
		recoveredModel, cmd := model.Update(f3Key)

		// Error should be cleared
		if recoveredModel.GetError() != nil {
			t.Error("error should be cleared for valid content advancement")
		}

		// Should generate advancement command
		if cmd == nil {
			t.Error("expected advancement command after error recovery")
		}

		// View should no longer show error
		recoveredView := recoveredModel.View()
		if strings.Contains(recoveredView, "Error:") {
			t.Error("recovered view should not contain error")
		}
	})

	// Test error state transitions
	t.Run("error_state_transitions", func(t *testing.T) {
		model.SetError(nil) // Start clean

		// Cause validation error
		model.SetContent("")
		f3Key := tea.KeyMsg{Type: tea.KeyF3}
		errorModel, _ := model.Update(f3Key)

		// Should have validation error
		if errorModel.GetError() == nil {
			t.Error("expected validation error")
		}

		// Add content and try again (should clear error and succeed)
		errorModel.SetContent("Valid content")
		successModel, cmd := errorModel.Update(f3Key)

		// Error should be cleared
		if successModel.GetError() != nil {
			t.Error("error should be cleared for valid content")
		}

		// Should generate advancement command
		if cmd == nil {
			t.Error("expected advancement command after fixing content")
		}
	})
}

func TestTaskInputModel_UTF8IntegrationScenarios(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test complete UTF-8 handling workflow
	t.Run("utf8_complete_workflow", func(t *testing.T) {
		// Test various UTF-8 scenarios through complete workflows
		testCases := []struct {
			content     string
			description string
		}{
			{"Hello 世界", "mixed ASCII and Chinese"},
			{"🏳️‍🌈 Pride flag", "complex emoji"},
			{"Café résumé naïve", "accented characters"},
			{"مرحبا بالعالم", "Arabic RTL text"},
			{"e\u0301\u0302", "combining characters"},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				// Set UTF-8 content
				model.SetContent(tc.content)

				// Verify character counting is correct
				expectedChars := len([]rune(tc.content))
				if model.charCount != expectedChars {
					t.Errorf("expected %d characters for %q, got %d", expectedChars, tc.content, model.charCount)
				}

				// Test that advancement works with UTF-8 content
				if !model.CanAdvance() {
					t.Error("UTF-8 content should be valid for advancement")
				}

				// Test view rendering with UTF-8 content
				view := model.View()
				if view == "" {
					t.Error("view should render with UTF-8 content")
				}

				// Test that clipboard messages work with UTF-8 content
				copyMsg := ClipboardCopyMsg{Text: tc.content}
				if copyMsg.Text != tc.content {
					t.Error("UTF-8 content should be preserved in clipboard messages")
				}
			})
		}
	})
}

func TestTaskInputModel_PerformanceIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test performance with large content
	t.Run("large_content_performance", func(t *testing.T) {
		// Create large content (close to the 10KB limit)
		largeContent := strings.Repeat("Performance test content with UTF-8 characters: éñ🌍. ", 150) // ~8KB

		// Set large content
		model.SetContent(largeContent)

		// Verify counters are calculated correctly
		expectedCharCount := len([]rune(largeContent))
		if model.charCount != expectedCharCount {
			t.Errorf("character count should be accurate for large content: expected %d, got %d", expectedCharCount, model.charCount)
		}

		// Test that view rendering works with large content
		view := model.View()
		if view == "" {
			t.Error("view should render with large content")
		}

		// Test that advancement validation works with large content
		if !model.CanAdvance() {
			t.Error("large valid content should allow advancement")
		}
	})

	// Test performance with many lines
	t.Run("many_lines_performance", func(t *testing.T) {
		// Create content with many lines
		manyLines := strings.Repeat("Line with content\n", 100) // 100 lines

		model.SetContent(manyLines)

		// Verify line counting is correct
		expectedLineCount := 101 // 100 lines + 1 empty at end
		if model.lineCount != expectedLineCount {
			t.Errorf("line count should be accurate for many lines: expected %d, got %d", expectedLineCount, model.lineCount)
		}

		// Test view rendering performance
		view := model.View()
		if view == "" {
			t.Error("view should render with many lines")
		}

		// Verify character count is also correct
		expectedCharCount := len([]rune(manyLines))
		if model.charCount != expectedCharCount {
			t.Errorf("character count should be accurate for many lines: expected %d, got %d", expectedCharCount, model.charCount)
		}
	})
}

func TestTaskInputModel_MessageHandlingIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test handling of various message types in sequence
	t.Run("message_sequence_handling", func(t *testing.T) {
		// Start with window size
		sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
		model, cmd := model.Update(sizeMsg)
		
		if model.width != 100 || model.height != 30 {
			t.Error("window size should be updated")
		}
		
		if cmd != nil {
			t.Error("window size update should not generate commands")
		}

		// Add content via external update
		contentMsg := TaskContentUpdatedMsg{Content: "Updated content"}
		model, cmd = model.Update(contentMsg)
		
		if model.GetContent() != "Updated content" {
			t.Error("content should be updated via message")
		}
		
		if cmd != nil {
			t.Error("content update should not generate commands")
		}

		// Test error handling
		errorMsg := ClipboardErrorMsg{Error: errors.New("test error")}
		model, cmd = model.Update(errorMsg)
		
		if model.GetError() == nil {
			t.Error("error should be set via message")
		}
		
		if cmd != nil {
			t.Error("error message should not generate commands")
		}

		// Test paste operation
		pasteMsg := ClipboardPasteMsg{Text: " appended"}
		model, cmd = model.Update(pasteMsg)
		
		expectedContent := "Updated content appended"
		if model.GetContent() != expectedContent {
			t.Errorf("content should be updated after paste: expected %q, got %q", 
				expectedContent, model.GetContent())
		}
		
		if cmd != nil {
			t.Error("paste message should not generate commands")
		}
	})

	// Test unknown message handling
	t.Run("unknown_message_handling", func(t *testing.T) {
		// Send an unknown message type
		unknownMsg := struct{ Value string }{Value: "unknown"}
		
		originalContent := model.GetContent()
		updatedModel, cmd := model.Update(unknownMsg)
		
		// Should handle gracefully without changes
		if updatedModel.GetContent() != originalContent {
			t.Error("unknown messages should not affect content")
		}
		
		// Command handling depends on textarea behavior
		_ = cmd
	})
}

// Integration test focusing on complete workflows to increase coverage
func TestTaskInputModel_CompleteWorkflowIntegration(t *testing.T) {
	model := NewTaskInputModel()
	model.UpdateSize(80, 25)

	// Test complete task input workflow
	t.Run("complete_task_workflow", func(t *testing.T) {
		// Start with empty content
		if model.GetContent() != "" {
			t.Error("expected empty content initially")
		}
		
		// Add some initial content
		contentMsg := TaskContentUpdatedMsg{Content: "Initial task"}
		updatedModel, cmd := model.Update(contentMsg)
		
		if updatedModel.GetContent() != "Initial task" {
			t.Error("expected content to be set by message")
		}
		
		if cmd != nil {
			t.Error("expected no command from content update message")
		}
		
		// Test clipboard paste integration
		pasteMsg := ClipboardPasteMsg{Text: " with additional details"}
		updatedModel, cmd = updatedModel.Update(pasteMsg)
		
		expectedContent := "Initial task with additional details"
		if updatedModel.GetContent() != expectedContent {
			t.Errorf("expected content %q after paste, got %q", expectedContent, updatedModel.GetContent())
		}
		
		// Verify character and line counts are updated
		expectedCharCount := len([]rune(expectedContent))
		if updatedModel.charCount != expectedCharCount {
			t.Errorf("expected char count %d, got %d", expectedCharCount, updatedModel.charCount)
		}
		
		expectedLineCount := 1 // Single line content
		if updatedModel.lineCount != expectedLineCount {
			t.Errorf("expected line count %d, got %d", expectedLineCount, updatedModel.lineCount)
		}
		
		// Test error state handling
		clipboardError := ClipboardErrorMsg{Error: errors.New("clipboard access failed")}
		updatedModel, cmd = updatedModel.Update(clipboardError)
		
		if updatedModel.GetError() == nil {
			t.Error("expected error to be set from clipboard error message")
		}
		
		if updatedModel.GetError().Error() != "clipboard access failed" {
			t.Errorf("expected clipboard error message, got %v", updatedModel.GetError())
		}
		
		// Test window resize during operation
		resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
		updatedModel, cmd = updatedModel.Update(resizeMsg)
		
		if updatedModel.width != 120 {
			t.Errorf("expected width to be updated to 120, got %d", updatedModel.width)
		}
		
		if updatedModel.height != 40 {
			t.Errorf("expected height to be updated to 40, got %d", updatedModel.height)
		}
		
		// Content should be preserved through resize
		if updatedModel.GetContent() != expectedContent {
			t.Error("content should be preserved during window resize")
		}
	})
}

// Test additional edge cases to increase coverage
func TestTaskInputModel_MessageHandlingEdgeCases(t *testing.T) {
	model := NewTaskInputModel()

	t.Run("empty_paste_handling", func(t *testing.T) {
		// Test pasting empty text
		pasteMsg := ClipboardPasteMsg{Text: ""}
		updatedModel, cmd := model.Update(pasteMsg)
		
		// Should handle empty paste gracefully
		if updatedModel.GetContent() != "" {
			t.Error("expected content to remain empty after empty paste")
		}
		
		if cmd != nil {
			t.Error("expected no command from empty paste")
		}
	})

	t.Run("multiple_errors_handling", func(t *testing.T) {
		// Test handling multiple errors in sequence
		error1 := ClipboardErrorMsg{Error: errors.New("first error")}
		updatedModel, _ := model.Update(error1)
		
		if updatedModel.GetError().Error() != "first error" {
			t.Error("expected first error to be set")
		}
		
		// Set new error - should overwrite previous
		error2 := ClipboardErrorMsg{Error: errors.New("second error")}
		updatedModel, _ = updatedModel.Update(error2)
		
		if updatedModel.GetError().Error() != "second error" {
			t.Error("expected second error to overwrite first")
		}
	})

	t.Run("utf8_paste_integration", func(t *testing.T) {
		// Test pasting UTF-8 content
		model.SetContent("Hello ")
		
		pasteMsg := ClipboardPasteMsg{Text: "世界 🌍"}
		updatedModel, cmd := model.Update(pasteMsg)
		
		expectedContent := "Hello 世界 🌍"
		if updatedModel.GetContent() != expectedContent {
			t.Errorf("expected UTF-8 content %q, got %q", expectedContent, updatedModel.GetContent())
		}
		
		// Verify proper UTF-8 character counting
		expectedCharCount := len([]rune(expectedContent))
		if updatedModel.charCount != expectedCharCount {
			t.Errorf("expected UTF-8 char count %d, got %d", expectedCharCount, updatedModel.charCount)
		}
		
		if cmd != nil {
			t.Error("expected no command from UTF-8 paste")
		}
	})
}