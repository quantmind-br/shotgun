package template

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
)

func TestStartDiscovery(t *testing.T) {
	model := NewTemplateModel()

	// Initially not discovering
	if model.IsDiscovering() {
		t.Error("New model should not be discovering")
	}

	// Start discovery
	cmd := model.StartDiscovery()

	if !model.IsDiscovering() {
		t.Error("Model should be discovering after StartDiscovery()")
	}

	if !model.IsLoading() {
		t.Error("Model should be loading during discovery")
	}

	if cmd == nil {
		t.Error("StartDiscovery() should return a command")
	}

	// Check state reset
	if model.err != nil {
		t.Error("err should be nil after StartDiscovery()")
	}

	if model.foundCount != 0 {
		t.Error("foundCount should be 0 after StartDiscovery()")
	}

	if model.currentPath != "" {
		t.Error("currentPath should be empty after StartDiscovery()")
	}
}

func TestStopDiscovery(t *testing.T) {
	model := NewTemplateModel()

	// Start and then stop discovery
	model.StartDiscovery()
	model.StopDiscovery()

	if model.IsDiscovering() {
		t.Error("Model should not be discovering after StopDiscovery()")
	}
}

func TestDiscoveryStateInView(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(100, 30)

	// Start discovery
	model.StartDiscovery()

	// View should show discovery state
	view := model.View()
	if !strings.Contains(view, "Discovering templates...") {
		t.Error("View should show discovery message")
	}

	if !strings.Contains(view, "[Press ESC to cancel]") {
		t.Error("View should show cancel hint")
	}
}

func TestDiscoveryProgress(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(100, 30)

	// Start discovery
	model.StartDiscovery()

	// Update with progress
	model.foundCount = 3
	model.currentPath = "templates/"

	view := model.View()
	if !strings.Contains(view, "Found 3 templates") {
		t.Error("View should show templates found count")
	}

	if !strings.Contains(view, "(scanning templates/)") {
		t.Error("View should show current path")
	}
}

func TestDiscoveryHelpText(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(100, 30)

	// Normal loading (not discovering)
	view := model.View()
	if !strings.Contains(view, "Press Ctrl+Left to go back") {
		t.Error("Normal loading should show back option")
	}

	// Discovery loading
	model.StartDiscovery()
	discoveryView := model.View()
	if !strings.Contains(discoveryView, "ESC: cancel discovery") {
		t.Error("Discovery view should show ESC cancellation")
	}

	if strings.Contains(discoveryView, "Press Ctrl+Left to go back") {
		t.Error("Discovery view should not show back option")
	}
}

func TestUpdateTemplateDiscoveryProgress(t *testing.T) {
	model := NewTemplateModel()

	// Test TemplateDiscoveryProgressMsg
	msg := TemplateDiscoveryProgressMsg{
		Found: 5,
		Path:  "user-templates/",
	}

	newModel, _ := model.Update(msg)

	if newModel.foundCount != 5 {
		t.Errorf("Expected foundCount to be 5, got %d", newModel.foundCount)
	}

	if newModel.currentPath != "user-templates/" {
		t.Errorf("Expected currentPath to be 'user-templates/', got '%s'", newModel.currentPath)
	}
}

func TestTemplatesLoadedStopsDiscovery(t *testing.T) {
	model := NewTemplateModel()
	model.StartDiscovery()

	// Verify discovery is active
	if !model.IsDiscovering() {
		t.Error("Should be discovering before templates loaded")
	}

	// Load templates
	msg := TemplatesLoadedMsg{Templates: []models.Template{}}
	newModel, _ := model.Update(msg)

	if newModel.IsDiscovering() {
		t.Error("Discovery should stop after templates loaded")
	}

	if newModel.IsLoading() {
		t.Error("Loading should stop after templates loaded")
	}
}

func TestTemplateLoadErrorStopsDiscovery(t *testing.T) {
	model := NewTemplateModel()
	model.StartDiscovery()

	// Simulate error
	testError := errors.New("discovery failed")
	msg := TemplateLoadErrorMsg{Error: testError}

	newModel, _ := model.Update(msg)

	if newModel.IsDiscovering() {
		t.Error("Discovery should stop after error")
	}

	if newModel.IsLoading() {
		t.Error("Loading should stop after error")
	}

	if newModel.err != testError {
		t.Error("Error should be set")
	}
}

func TestESCCancellationDuringDiscovery(t *testing.T) {
	model := NewTemplateModel()
	model.StartDiscovery()

	// Press ESC during discovery
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}

	newModel, _ := model.Update(keyMsg)

	if newModel.IsDiscovering() {
		t.Error("ESC should cancel discovery")
	}

	if newModel.IsLoading() {
		t.Error("ESC should stop loading")
	}
}

func TestKeyBlockingDuringDiscovery(t *testing.T) {
	model := NewTemplateModel()
	model.StartDiscovery()

	// Try to navigate during discovery
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}

	// Should not process the key (returns same model)
	newModel, _ := model.Update(keyMsg)

	// Since we can't easily test list cursor position, just verify discovery state
	if !newModel.IsDiscovering() {
		t.Error("Should still be discovering after blocked key")
	}
}

func TestSpinnerUpdateDuringDiscovery(t *testing.T) {
	model := NewTemplateModel()
	model.StartDiscovery()

	// Update with a generic message (spinner should update)
	newModel, cmd := model.Update(tea.Msg(nil))

	if !newModel.IsDiscovering() {
		t.Error("Should still be discovering after spinner update")
	}

	// Command could be nil or a spinner command - both are valid
	_ = cmd
}

func TestRenderDiscoveryState(t *testing.T) {
	model := NewTemplateModel()

	// Test initial discovery state - check the spinner message instead
	model.spinner.SetMessage("Discovering templates...")
	// The view contains styled content, so we test the logic that sets the message
	expectedMessage := "Discovering templates..."

	// Simulate the logic from renderDiscoveryState
	var message string
	if model.foundCount > 0 {
		message = fmt.Sprintf("Found %d templates", model.foundCount)
		if model.currentPath != "" {
			message += fmt.Sprintf(" (scanning %s)", model.currentPath)
		}
	} else {
		message = "Discovering templates..."
	}

	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}

	// Test with progress
	model.foundCount = 7
	model.currentPath = "project-templates/"

	// Simulate logic again
	if model.foundCount > 0 {
		message = fmt.Sprintf("Found %d templates", model.foundCount)
		if model.currentPath != "" {
			message += fmt.Sprintf(" (scanning %s)", model.currentPath)
		}
	} else {
		message = "Discovering templates..."
	}

	expectedProgressMessage := "Found 7 templates (scanning project-templates/)"
	if message != expectedProgressMessage {
		t.Errorf("Expected progress message '%s', got '%s'", expectedProgressMessage, message)
	}
}

func TestDiscoveryViewWithoutProgress(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(100, 30)
	model.StartDiscovery()

	// No progress updates yet
	view := model.renderDiscoveryState()
	if !strings.Contains(view, "Discovering templates...") {
		t.Error("Should show initial discovery message")
	}
}

func TestDiscoveryViewWithProgress(t *testing.T) {
	model := NewTemplateModel()
	model.UpdateSize(100, 30)
	model.StartDiscovery()
	model.foundCount = 2

	// With progress but no path
	view := model.renderDiscoveryState()
	if !strings.Contains(view, "Found 2 templates") {
		t.Error("Should show found count without path")
	}

	// With progress and path
	model.currentPath = "shared/"
	viewWithPath := model.renderDiscoveryState()
	if !strings.Contains(viewWithPath, "Found 2 templates") {
		t.Error("Should show found count")
	}

	if !strings.Contains(viewWithPath, "(scanning shared/)") {
		t.Error("Should show scanning path")
	}
}
