package builder

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/diogopedro/shotgun/internal/models"
)

func TestNewPromptGenerator(t *testing.T) {
	generator := NewPromptGenerator()
	
	if generator == nil {
		t.Fatal("NewPromptGenerator() returned nil")
	}
	
	
	if generator.fileStructureBuilder == nil {
		t.Error("fileStructureBuilder is nil")
	}
}

func TestGeneratePrompt_Basic(t *testing.T) {
	generator := NewPromptGenerator()
	
	config := GenerationConfig{
		Template: &models.Template{
			ID:          "test",
			Name:        "Test Template",
			Version:     "1.0",
			Description: "Test template for unit tests",
			Content:     "Task: {{TASK}}\nRules: {{RULES}}\nFiles: {{SELECTED_FILES_COUNT}}",
		},
		Variables:     make(map[string]string),
		SelectedFiles: []string{"file1.txt", "file2.txt"},
		TaskContent:   "Implement feature X",
		RulesContent:  "Follow coding standards",
		OutputPath:    "",
	}
	
	ctx := context.Background()
	result, err := generator.GeneratePrompt(ctx, config)
	
	if err != nil {
		t.Fatalf("GeneratePrompt failed: %v", err)
	}
	
	if result == nil {
		t.Fatal("GeneratePrompt returned nil result")
	}
	
	// Check that variables were substituted
	if !strings.Contains(result.Content, "Implement feature X") {
		t.Errorf("Content should contain 'Implement feature X', got: %s", result.Content)
	}
	
	if result.FileCount != 2 {
		t.Errorf("FileCount mismatch. Expected 2, got: %d", result.FileCount)
	}
	
	if result.TemplateSize == 0 {
		t.Error("TemplateSize should not be 0")
	}
	
	if result.TotalSize == 0 {
		t.Error("TotalSize should not be 0")
	}
}

func TestGeneratePrompt_NilTemplate(t *testing.T) {
	generator := NewPromptGenerator()
	
	config := GenerationConfig{
		Template:      nil,
		Variables:     make(map[string]string),
		SelectedFiles: []string{"file1.txt"},
		TaskContent:   "Test task",
		RulesContent:  "Test rules",
		OutputPath:    "",
	}
	
	ctx := context.Background()
	result, err := generator.GeneratePrompt(ctx, config)
	
	if err == nil {
		t.Error("GeneratePrompt should fail with nil template")
	}
	
	if result != nil {
		t.Error("GeneratePrompt should return nil result on error")
	}
}

func TestGeneratePrompt_AutomaticVariables(t *testing.T) {
	generator := NewPromptGenerator()
	
	config := GenerationConfig{
		Template: &models.Template{
			ID:      "test",
			Name:    "Test Template",
			Version: "1.0",
			Content: "Date: {{CURRENT_DATE}}\nFile count: {{SELECTED_FILES_COUNT}}",
		},
		Variables:     make(map[string]string),
		SelectedFiles: []string{"file1.txt", "file2.txt", "file3.txt"},
		TaskContent:   "",
		RulesContent:  "",
		OutputPath:    "",
	}
	
	ctx := context.Background()
	result, err := generator.GeneratePrompt(ctx, config)
	
	if err != nil {
		t.Fatalf("GeneratePrompt failed: %v", err)
	}
	
	// Check that CURRENT_DATE was set
	today := time.Now().Format("2006-01-02")
	if !strings.Contains(result.Content, today) {
		t.Errorf("Content should contain today's date %s, got: %s", today, result.Content)
	}
	
	// Check that SELECTED_FILES_COUNT was set
	if !strings.Contains(result.Content, "File count: 3") {
		t.Errorf("Content should contain 'File count: 3', got: %s", result.Content)
	}
}

func TestGenerateAsync(t *testing.T) {
	generator := NewPromptGenerator()
	
	config := GenerationConfig{
		Template: &models.Template{
			ID:      "test",
			Name:    "Test Template",
			Version: "1.0",
			Content: "Async test: {{TASK}}",
		},
		Variables:     make(map[string]string),
		SelectedFiles: []string{},
		TaskContent:   "Async task",
		RulesContent:  "",
		OutputPath:    "",
	}
	
	callback := func(stage string, progress float64) {
		if progress < 0 || progress > 1 {
			t.Errorf("Invalid progress value: %f", progress)
		}
	}
	
	cmd := generator.GenerateAsync(config, callback)
	
	if cmd == nil {
		t.Fatal("GenerateAsync returned nil command")
	}
	
	// Execute the command
	msg := cmd()
	
	if msg == nil {
		t.Fatal("Command execution returned nil message")
	}
	
	// Check the result
	completeMsg, ok := msg.(GenerationCompleteMsg)
	if !ok {
		t.Fatalf("Expected GenerationCompleteMsg, got %T", msg)
	}
	
	if completeMsg.Error != nil {
		t.Errorf("Async generation failed: %v", completeMsg.Error)
	}
	
	if completeMsg.Result == nil {
		t.Error("Async generation returned nil result")
	}
}