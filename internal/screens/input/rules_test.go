package input

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewRulesInputModel(t *testing.T) {
	model := NewRulesInputModel()

	// Test initial state
	if model.content != "" {
		t.Errorf("Expected empty initial content, got '%s'", model.content)
	}

	if model.charCount != 0 {
		t.Errorf("Expected initial character count of 0, got %d", model.charCount)
	}

	if model.lineCount != 1 {
		t.Errorf("Expected initial line count of 1, got %d", model.lineCount)
	}

	if model.ready != false {
		t.Errorf("Expected initial ready state to be false")
	}

	if model.err != nil {
		t.Errorf("Expected no initial error, got %v", model.err)
	}

	// Test textarea configuration
	if model.textarea.CharLimit != 5000 {
		t.Errorf("Expected character limit of 5000, got %d", model.textarea.CharLimit)
	}

	if !model.textarea.Focused() {
		t.Error("Expected textarea to be focused initially")
	}
}

func TestRulesInputModel_SetContent(t *testing.T) {
	model := NewRulesInputModel()

	testContent := "Use TypeScript strict mode\nInclude error handling"
	model.SetContent(testContent)

	if model.GetContent() != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, model.GetContent())
	}

	// Test that counters are updated
	expectedCharCount := len([]rune(testContent))
	if model.charCount != expectedCharCount {
		t.Errorf("Expected character count %d, got %d", expectedCharCount, model.charCount)
	}

	expectedLineCount := 2 // Two lines due to \n
	if model.lineCount != expectedLineCount {
		t.Errorf("Expected line count %d, got %d", expectedLineCount, model.lineCount)
	}
}

func TestRulesInputModel_UTF8Support(t *testing.T) {
	model := NewRulesInputModel()

	// Test UTF-8 characters
	utf8Content := "使用 TypeScript 模式 🚀\n日本語 テスト"
	model.SetContent(utf8Content)

	// Character count should be based on runes, not bytes
	expectedCharCount := len([]rune(utf8Content))
	if model.charCount != expectedCharCount {
		t.Errorf("Expected UTF-8 character count %d, got %d", expectedCharCount, model.charCount)
	}

	expectedLineCount := 2
	if model.lineCount != expectedLineCount {
		t.Errorf("Expected line count %d, got %d", expectedLineCount, model.lineCount)
	}
}

func TestRulesInputModel_CanAdvance(t *testing.T) {
	model := NewRulesInputModel()

	// Rules input should always allow advancing since it's optional
	if !model.CanAdvance() {
		t.Error("Expected CanAdvance to be true for empty content (optional field)")
	}

	// Test with content
	model.SetContent("Some rules")
	if !model.CanAdvance() {
		t.Error("Expected CanAdvance to be true with content")
	}
}

func TestRulesInputModel_ErrorHandling(t *testing.T) {
	model := NewRulesInputModel()

	testError := errors.New("test error")
	model.SetError(testError)

	if model.GetError() != testError {
		t.Errorf("Expected error '%v', got '%v'", testError, model.GetError())
	}

	// Clear error
	model.SetError(nil)
	if model.GetError() != nil {
		t.Errorf("Expected no error after clearing, got '%v'", model.GetError())
	}
}

func TestRulesInputModel_ReadyState(t *testing.T) {
	model := NewRulesInputModel()

	// Initial state
	if model.IsReady() {
		t.Error("Expected initial ready state to be false")
	}

	// Set ready
	model.SetReady(true)
	if !model.IsReady() {
		t.Error("Expected ready state to be true after setting")
	}

	// Unset ready
	model.SetReady(false)
	if model.IsReady() {
		t.Error("Expected ready state to be false after unsetting")
	}
}

func TestRulesInputModel_UpdateSize(t *testing.T) {
	model := NewRulesInputModel()

	width, height := 100, 30
	model.UpdateSize(width, height)

	if model.width != width {
		t.Errorf("Expected width %d, got %d", width, model.width)
	}

	if model.height != height {
		t.Errorf("Expected height %d, got %d", height, model.height)
	}

	// Test textarea dimensions are updated (check they're reasonable, exact values may differ due to internal calculations)
	if model.textarea.Width() < width/2 {
		t.Errorf("Expected textarea width to be reasonable relative to window width %d, got %d", width, model.textarea.Width())
	}

	if model.textarea.Height() < height/2 {
		t.Errorf("Expected textarea height to be reasonable relative to window height %d, got %d", height, model.textarea.Height())
	}
}

func TestRulesInputModel_KeyboardNavigation(t *testing.T) {
	model := NewRulesInputModel()

	// Test F3 key (advance)
	keyMsg := tea.KeyMsg{Type: tea.KeyF3}
	updatedModel, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command to be returned for F3 key")
	}

	// The model should be returned (not modified for this key)
	_ = updatedModel

	// Test F4 key (skip)
	keyMsg = tea.KeyMsg{Type: tea.KeyF4}
	updatedModel, cmd = model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command to be returned for F4 key")
	}

	// Test F2 key (back)
	keyMsg = tea.KeyMsg{Type: tea.KeyF2}
	updatedModel, cmd = model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command to be returned for F2 key")
	}
}

func TestRulesInputModel_ClipboardOperations(t *testing.T) {
	model := NewRulesInputModel()

	// Add some content first
	model.SetContent("Test content")

	// Test Ctrl+C (copy)
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command to be returned for Ctrl+C")
	}

	_ = updatedModel

	// Test Ctrl+V (paste)
	keyMsg = tea.KeyMsg{Type: tea.KeyCtrlV}
	updatedModel, cmd = model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command to be returned for Ctrl+V")
	}

	// Test clipboard paste message
	pasteMsg := ClipboardPasteMsg{Text: " pasted text"}
	updatedModel, cmd = model.Update(pasteMsg)

	// Content should be updated
	if !containsSubstring(updatedModel.GetContent(), "pasted text") {
		t.Error("Expected pasted text to be added to content")
	}
}

func TestRulesInputModel_MessageHandling(t *testing.T) {
	model := NewRulesInputModel()

	// Test RulesContentUpdatedMsg
	contentMsg := RulesContentUpdatedMsg{Content: "Updated rules content"}
	updatedModel, _ := model.Update(contentMsg)

	if updatedModel.GetContent() != "Updated rules content" {
		t.Errorf("Expected content to be updated to 'Updated rules content', got '%s'", updatedModel.GetContent())
	}

	// Test ClipboardErrorMsg
	errorMsg := ClipboardErrorMsg{Error: errors.New("clipboard error")}
	updatedModel, _ = model.Update(errorMsg)

	if updatedModel.GetError() == nil {
		t.Error("Expected error to be set from ClipboardErrorMsg")
	}

	if updatedModel.GetError().Error() != "clipboard error" {
		t.Errorf("Expected error message 'clipboard error', got '%s'", updatedModel.GetError().Error())
	}
}

func TestRulesInputModel_WindowSizeMsg(t *testing.T) {
	model := NewRulesInputModel()

	// Test window size message
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := model.Update(sizeMsg)

	if updatedModel.width != 120 {
		t.Errorf("Expected width to be updated to 120, got %d", updatedModel.width)
	}

	if updatedModel.height != 40 {
		t.Errorf("Expected height to be updated to 40, got %d", updatedModel.height)
	}
}

func TestRulesInputModel_View(t *testing.T) {
	model := NewRulesInputModel()
	model.UpdateSize(80, 25)

	view := model.View()

	// Should not be empty
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain "Optional" indication
	if !containsSubstring(view, "Optional") {
		t.Error("Expected view to contain 'Optional' indication")
	}

	// Test error view
	model.SetError(errors.New("test error"))
	errorView := model.View()

	if !containsSubstring(errorView, "Error") {
		t.Error("Expected error view to contain 'Error'")
	}

	if !containsSubstring(errorView, "test error") {
		t.Error("Expected error view to contain the error message")
	}
}

func TestRulesInputModel_CounterUpdates(t *testing.T) {
	model := NewRulesInputModel()

	// Test empty content
	model.SetContent("")
	if model.charCount != 0 {
		t.Errorf("Expected character count 0 for empty content, got %d", model.charCount)
	}
	if model.lineCount != 1 {
		t.Errorf("Expected line count 1 for empty content, got %d", model.lineCount)
	}

	// Test single line
	model.SetContent("Single line")
	expectedChars := len([]rune("Single line"))
	if model.charCount != expectedChars {
		t.Errorf("Expected character count %d, got %d", expectedChars, model.charCount)
	}
	if model.lineCount != 1 {
		t.Errorf("Expected line count 1 for single line, got %d", model.lineCount)
	}

	// Test multiple lines
	model.SetContent("Line 1\nLine 2\nLine 3")
	expectedChars = len([]rune("Line 1\nLine 2\nLine 3"))
	if model.charCount != expectedChars {
		t.Errorf("Expected character count %d, got %d", expectedChars, model.charCount)
	}
	if model.lineCount != 3 {
		t.Errorf("Expected line count 3 for three lines, got %d", model.lineCount)
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(text, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(text) < len(substr) {
		return false
	}

	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Test message types exist and can be instantiated
func TestMessageTypes(t *testing.T) {
	// Test that our custom message types can be created and used
	rulesMsg := RulesInputMsg{}
	backMsg := BackToTaskMsg{}
	contentMsg := RulesContentUpdatedMsg{Content: "test"}
	skipMsg := SkipRulesMsg{}

	// Verify types exist and can be instantiated
	_ = rulesMsg
	_ = backMsg
	_ = contentMsg
	_ = skipMsg
}
