package filetree

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
)

func TestStartScanning(t *testing.T) {
	model := NewFileTreeModel()

	// Initially not scanning
	if model.IsScanning() {
		t.Error("New model should not be scanning")
	}

	// Start scanning
	cmd := model.StartScanning()

	if !model.IsScanning() {
		t.Error("Model should be scanning after StartScanning()")
	}

	if cmd == nil {
		t.Error("StartScanning() should return a command")
	}

	// Check state reset
	if model.scanError != nil {
		t.Error("scanError should be nil after StartScanning()")
	}

	if model.filesFound != 0 {
		t.Error("filesFound should be 0 after StartScanning()")
	}

	if model.currentDir != "" {
		t.Error("currentDir should be empty after StartScanning()")
	}
}

func TestStopScanning(t *testing.T) {
	model := NewFileTreeModel()

	// Start and then stop scanning
	model.StartScanning()
	model.StopScanning()

	if model.IsScanning() {
		t.Error("Model should not be scanning after StopScanning()")
	}
}

func TestScanningStateInView(t *testing.T) {
	model := NewFileTreeModel()

	// Start scanning
	model.StartScanning()

	// View should show scanning state
	view := model.View()
	if !strings.Contains(view, "Scanning project files...") {
		t.Error("View should show scanning message")
	}

	if !strings.Contains(view, "[Press ESC to cancel]") {
		t.Error("View should show cancel hint")
	}
}

func TestScanningProgress(t *testing.T) {
	model := NewFileTreeModel()

	// Start scanning
	model.StartScanning()

	// Update with progress
	model.filesFound = 50
	model.currentDir = "src"

	view := model.View()
	if !strings.Contains(view, "Found 50 files") {
		t.Error("View should show files found count")
	}

	if !strings.Contains(view, "(scanning src/)") {
		t.Error("View should show current directory")
	}
}

func TestErrorStateInView(t *testing.T) {
	model := NewFileTreeModel()

	// Set error state
	model.scanError = errors.New("permission denied")

	view := model.View()
	if !strings.Contains(view, "❌ Scan failed: permission denied") {
		t.Error("View should show error message")
	}

	if !strings.Contains(view, "Press 'r' to retry") {
		t.Error("View should show retry hint")
	}
}

func TestUpdateScanMessages(t *testing.T) {
	model := NewFileTreeModel()

	// Test ScanCompleteMsg
	nodes := createTestNodes()
	msg := ScanCompleteMsg{Nodes: nodes}

	newModel, _ := model.Update(msg)
	m := newModel.(FileTreeModel)

	if m.IsScanning() {
		t.Error("Scanning should stop after ScanCompleteMsg")
	}

	if m.filesFound != len(nodes) {
		t.Errorf("Expected filesFound to be %d, got %d", len(nodes), m.filesFound)
	}
}

func TestUpdateScanErrorMsg(t *testing.T) {
	model := NewFileTreeModel()
	model.StartScanning()

	// Test ScanErrorMsg
	testError := errors.New("test error")
	msg := ScanErrorMsg{Error: testError}

	newModel, _ := model.Update(msg)
	m := newModel.(FileTreeModel)

	if m.IsScanning() {
		t.Error("Scanning should stop after ScanErrorMsg")
	}

	if m.scanError != testError {
		t.Error("scanError should be set")
	}
}

func TestUpdateScanProgressMsg(t *testing.T) {
	model := NewFileTreeModel()
	model.StartScanning()

	// Test ScanProgressMsg
	msg := ScanProgressMsg{
		FilesFound: 25,
		CurrentDir: "lib",
	}

	newModel, _ := model.Update(msg)
	m := newModel.(FileTreeModel)

	if m.filesFound != 25 {
		t.Errorf("Expected filesFound to be 25, got %d", m.filesFound)
	}

	if m.currentDir != "lib" {
		t.Errorf("Expected currentDir to be 'lib', got '%s'", m.currentDir)
	}
}

func TestSpinnerUpdate(t *testing.T) {
	model := NewFileTreeModel()
	model.StartScanning()

	// Update with a generic message (spinner should update)
	newModel, cmd := model.Update(tea.Msg(nil))
	m := newModel.(FileTreeModel)

	if !m.IsScanning() {
		t.Error("Should still be scanning after spinner update")
	}

	// Command could be nil or a spinner command - both are valid
	_ = cmd
}

func TestESCCancellation(t *testing.T) {
	model := NewFileTreeModel()
	model.StartScanning()

	// Press ESC during scanning
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}, Alt: false}
	keyMsg.Type = tea.KeyEsc

	newModel, _ := model.Update(keyMsg)
	m := newModel.(FileTreeModel)

	if m.IsScanning() {
		t.Error("ESC should cancel scanning")
	}
}

func TestKeyBlockingDuringScanning(t *testing.T) {
	model := NewFileTreeModel()
	model.StartScanning()

	// Try to navigate during scanning
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}, Alt: false}

	initialCursor := model.cursor
	newModel, _ := model.Update(keyMsg)
	m := newModel.(FileTreeModel)

	if m.cursor != initialCursor {
		t.Error("Navigation keys should be blocked during scanning")
	}
}

func TestHelpBarDuringScanning(t *testing.T) {
	model := NewFileTreeModel()

	// Normal help
	help := model.helpBar()
	if !strings.Contains(help, "navigate") {
		t.Error("Normal help should mention navigation")
	}

	// Scanning help
	model.StartScanning()
	scanningHelp := model.helpBar()
	if !strings.Contains(scanningHelp, "ESC: cancel scanning") {
		t.Error("Scanning help should mention ESC cancellation")
	}

	if strings.Contains(scanningHelp, "navigate") {
		t.Error("Scanning help should not mention navigation")
	}
}

func TestLoadFromScannerAsync(t *testing.T) {
	model := NewFileTreeModel()
	ctx := context.Background()

	// This should return a command that will eventually send ScanCompleteMsg or ScanErrorMsg
	cmd := model.LoadFromScanner(ctx, ".")

	if cmd == nil {
		t.Error("LoadFromScanner should return a command")
	}

	// Execute the command to get the message
	msg := cmd()

	// Should get either ScanCompleteMsg or ScanErrorMsg
	switch msg.(type) {
	case ScanCompleteMsg:
		// Success case
	case ScanErrorMsg:
		// Error case (might happen in test environment)
	default:
		t.Errorf("LoadFromScanner should return ScanCompleteMsg or ScanErrorMsg, got %T", msg)
	}
}

// Helper function to create test nodes
func createTestNodes() []*models.FileNode {
	// Return mock nodes for testing - we need to import models
	return []*models.FileNode{
		{Path: "file1.go", Name: "file1.go", IsDirectory: false, IsBinary: false},
		{Path: "file2.txt", Name: "file2.txt", IsDirectory: false, IsBinary: false},
		{Path: "binary.exe", Name: "binary.exe", IsDirectory: false, IsBinary: true},
	}
}
