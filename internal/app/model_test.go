package app

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	
	if app == nil {
		t.Fatal("NewApp returned nil")
	}
	
	if app.CurrentScreen != FileTreeScreen {
		t.Errorf("Expected CurrentScreen = FileTreeScreen, got %v", app.CurrentScreen)
	}
	
	if app.SelectedFiles == nil {
		t.Error("Expected SelectedFiles to be initialized")
	}
	
	if len(app.SelectedFiles) != 0 {
		t.Errorf("Expected empty SelectedFiles, got %d items", len(app.SelectedFiles))
	}
	
	if app.SelectedTemplate != nil {
		t.Error("Expected SelectedTemplate to be nil initially")
	}
	
	if app.TaskContent != "" {
		t.Errorf("Expected empty TaskContent, got '%s'", app.TaskContent)
	}
	
	if app.RulesContent != "" {
		t.Errorf("Expected empty RulesContent, got '%s'", app.RulesContent)
	}
	
	if app.ctx == nil {
		t.Error("Expected context to be initialized")
	}
	
	if app.cancel == nil {
		t.Error("Expected cancel function to be initialized")
	}
	
	// Test progress indicator
	if app.Progress.GetCurrent() != 1 {
		t.Errorf("Expected progress current = 1, got %d", app.Progress.GetCurrent())
	}
	
	if app.Progress.GetTotal() != 5 {
		t.Errorf("Expected progress total = 5, got %d", app.Progress.GetTotal())
	}
}

func TestAppCleanup(t *testing.T) {
	app := NewApp()
	
	// Test that cleanup doesn't panic
	app.Cleanup()
	
	// Context should be cancelled after cleanup
	select {
	case <-app.ctx.Done():
		// Expected - context should be cancelled
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected context to be cancelled after cleanup")
	}
}

func TestScreenTypeString(t *testing.T) {
	tests := []struct {
		screen   ScreenType
		expected string
	}{
		{FileTreeScreen, "FileTree"},
		{TemplateScreen, "Template"},
		{TaskScreen, "TaskInput"},
		{RulesScreen, "RulesInput"},
		{ConfirmScreen, "Confirm"},
		{ScreenType(99), "Unknown"},
	}
	
	for _, tt := range tests {
		result := tt.screen.String()
		if result != tt.expected {
			t.Errorf("ScreenType(%d).String() = %q, want %q", tt.screen, result, tt.expected)
		}
	}
}

func TestSetCurrentScreen(t *testing.T) {
	app := NewApp()
	
	// Test navigation from FileTree to Template
	app.SetCurrentScreen(TemplateScreen)
	
	if app.CurrentScreen != TemplateScreen {
		t.Errorf("Expected CurrentScreen = TemplateScreen, got %v", app.CurrentScreen)
	}
	
	// Progress should be updated
	if app.Progress.GetCurrent() != 2 {
		t.Errorf("Expected progress current = 2, got %d", app.Progress.GetCurrent())
	}
	
	// Test navigation to Task screen
	app.SetCurrentScreen(TaskScreen)
	
	if app.CurrentScreen != TaskScreen {
		t.Errorf("Expected CurrentScreen = TaskScreen, got %v", app.CurrentScreen)
	}
	
	if app.Progress.GetCurrent() != 3 {
		t.Errorf("Expected progress current = 3, got %d", app.Progress.GetCurrent())
	}
}

func TestUpdateWindowSize(t *testing.T) {
	app := NewApp()
	
	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}
	
	app.UpdateWindowSize(msg)
	
	if app.WindowSize.Width != 100 {
		t.Errorf("Expected WindowSize.Width = 100, got %d", app.WindowSize.Width)
	}
	
	if app.WindowSize.Height != 30 {
		t.Errorf("Expected WindowSize.Height = 30, got %d", app.WindowSize.Height)
	}
	
	// Progress indicator should have updated width
	// Note: We can't easily test this without exposing internal state
}

func TestGetCurrentScreenModel(t *testing.T) {
	app := NewApp()
	
	// Test FileTree screen
	model := app.GetCurrentScreenModel()
	if model == nil {
		t.Error("Expected non-nil model for FileTree screen")
	}
	
	// Test Template screen
	app.SetCurrentScreen(TemplateScreen)
	model = app.GetCurrentScreenModel()
	if model == nil {
		t.Error("Expected non-nil model for Template screen")
	}
	
	// Test Task screen
	app.SetCurrentScreen(TaskScreen)
	model = app.GetCurrentScreenModel()
	if model == nil {
		t.Error("Expected non-nil model for Task screen")
	}
}

func TestSaveCurrentScreenState(t *testing.T) {
	app := NewApp()
	
	// Set some state in FileTree screen (simulate selected files)
	app.SelectedFiles = []string{"/test/file1.txt", "/test/file2.txt"}
	
	// Save state
	app.saveCurrentScreenState()
	
	// Switch to template screen and set template
	app.SetCurrentScreen(TemplateScreen)
	template := &Template{
		Name:        "Test Template",
		Path:        "/templates/test.md",
		Description: "Test template description",
	}
	app.Template.selected = template
	
	// Save template state
	app.saveCurrentScreenState()
	
	if app.SelectedTemplate == nil {
		t.Error("Expected SelectedTemplate to be saved")
	}
	
	if app.SelectedTemplate.Name != "Test Template" {
		t.Errorf("Expected template name 'Test Template', got '%s'", app.SelectedTemplate.Name)
	}
}

func TestLoadScreenState(t *testing.T) {
	app := NewApp()
	
	// Set up shared state
	app.SelectedFiles = []string{"/test/file1.txt"}
	app.SelectedTemplate = &Template{Name: "Test Template"}
	app.TaskContent = "Test task description"
	app.RulesContent = "Test rules"
	
	// Test loading Template screen state
	app.SetCurrentScreen(TemplateScreen)
	if app.Template.selected == nil {
		t.Error("Expected template to be loaded from shared state")
	}
	
	// Test loading Task screen state
	app.SetCurrentScreen(TaskScreen)
	if app.TaskInput.content != "Test task description" {
		t.Errorf("Expected task content to be loaded, got '%s'", app.TaskInput.content)
	}
	
	// Test loading Rules screen state
	app.SetCurrentScreen(RulesScreen)
	if app.RulesInput.content != "Test rules" {
		t.Errorf("Expected rules content to be loaded, got '%s'", app.RulesInput.content)
	}
}

func TestBuildConfirmationSummary(t *testing.T) {
	app := NewApp()
	
	// Set up test data
	app.SelectedFiles = []string{"/test/file1.txt", "/test/file2.txt"}
	app.SelectedTemplate = &Template{Name: "Test Template"}
	app.TaskContent = "Create a new feature"
	app.RulesContent = "Use TypeScript"
	
	app.buildConfirmationSummary()
	
	summary := app.Confirmation.summary
	
	if summary == "" {
		t.Error("Expected non-empty confirmation summary")
	}
	
	// Check that summary contains key information
	if !containsString(summary, "file1.txt") {
		t.Error("Expected summary to contain selected files")
	}
	
	if !containsString(summary, "Test Template") {
		t.Error("Expected summary to contain template name")
	}
	
	if !containsString(summary, "Create a new feature") {
		t.Error("Expected summary to contain task content")
	}
	
	if !containsString(summary, "Use TypeScript") {
		t.Error("Expected summary to contain rules content")
	}
}

func TestInit(t *testing.T) {
	app := NewApp()
	
	cmd := app.Init()
	
	// Init should return a batch command or nil
	// We don't enforce that it must return a command
	_ = cmd // Init may return nil, and that's okay
}

func TestContext(t *testing.T) {
	app := NewApp()
	
	ctx := app.Context()
	
	if ctx == nil {
		t.Error("Expected non-nil context")
	}
	
	// Context should not be cancelled initially
	select {
	case <-ctx.Done():
		t.Error("Expected context to not be cancelled initially")
	default:
		// Expected - context should not be done
	}
}

// Helper function
func containsString(text, substr string) bool {
	return len(text) >= len(substr) && 
		   text != substr && 
		   (text[:len(substr)] == substr || 
			text[len(text)-len(substr):] == substr || 
			findSubstring(text, substr))
}

func findSubstring(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}