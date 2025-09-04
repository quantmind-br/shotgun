package builder

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/models"
)

// GenerationConfig contains all the configuration needed for prompt generation
type GenerationConfig struct {
	Template      *models.Template
	Variables     map[string]string
	SelectedFiles []string
	TaskContent   string
	RulesContent  string
	OutputPath    string
}

// GeneratedPrompt contains the result of prompt generation with metadata
type GeneratedPrompt struct {
	Content       string
	TemplateSize  int64
	FileCount     int
	TotalSize     int64
	GeneratedAt   time.Time
}

// GenerationProgressCallback is called during async generation to report progress
type GenerationProgressCallback func(stage string, progress float64)

// GenerationStage represents different phases of prompt generation
type GenerationStage string

const (
	StageProcessingTemplate  GenerationStage = "Processing template"
	StageLoadingFiles       GenerationStage = "Loading file contents"
	StageAssemblingStructure GenerationStage = "Assembling file structure"
	StageWritingOutput      GenerationStage = "Writing output file"
	StageComplete           GenerationStage = "Generation complete"
)

// GenerationProgressMsg is sent during async generation
type GenerationProgressMsg struct {
	Stage    GenerationStage
	Progress float64
	Message  string
}

// GenerationCompleteMsg is sent when generation is complete
type GenerationCompleteMsg struct {
	Result *GeneratedPrompt
	Error  error
}

// GenerationCancelMsg is sent to cancel generation
type GenerationCancelMsg struct{}

// PromptGeneratorInterface defines the interface for prompt generation
type PromptGeneratorInterface interface {
	GeneratePrompt(ctx context.Context, config GenerationConfig) (*GeneratedPrompt, error)
	GenerateAsync(config GenerationConfig, callback GenerationProgressCallback) tea.Cmd
}

// PromptGenerator handles the generation of final prompts
type PromptGenerator struct {
	fileStructureBuilder *FileStructureBuilder
}

// NewPromptGenerator creates a new PromptGenerator instance
func NewPromptGenerator() *PromptGenerator {
	return &PromptGenerator{
		fileStructureBuilder: NewFileStructureBuilder(),
	}
}

// GeneratePrompt combines template, variables, and file structure into final prompt
func (pg *PromptGenerator) GeneratePrompt(ctx context.Context, config GenerationConfig) (*GeneratedPrompt, error) {
	if config.Template == nil {
		return nil, fmt.Errorf("template is required for generation")
	}

	startTime := time.Now()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Step 1: Prepare variables for template processing
	variables := make(map[string]string)
	
	// Copy provided variables
	for k, v := range config.Variables {
		variables[k] = v
	}
	
	// Add required template variables
	variables["TASK"] = config.TaskContent
	variables["RULES"] = config.RulesContent
	
	// Add automatic variables
	variables["CURRENT_DATE"] = time.Now().Format("2006-01-02")
	variables["SELECTED_FILES_COUNT"] = fmt.Sprintf("%d", len(config.SelectedFiles))
	
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Step 2: Generate file structure if files are selected
	var fileStructure string
	if len(config.SelectedFiles) > 0 {
		structure, err := pg.fileStructureBuilder.GenerateStructure(ctx, config.SelectedFiles)
		if err != nil {
			return nil, fmt.Errorf("failed to generate file structure: %w", err)
		}
		fileStructure = structure
	}
	
	variables["FILE_STRUCTURE"] = fileStructure

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Step 3: Process template with simple variable substitution
	processedTemplate := config.Template.Content
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		processedTemplate = strings.ReplaceAll(processedTemplate, placeholder, value)
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Step 4: Use processed template as final content
	finalContent := processedTemplate

	// Step 5: Calculate metadata
	templateSize := int64(len(config.Template.Content))
	totalSize := int64(len(finalContent))
	
	result := &GeneratedPrompt{
		Content:       finalContent,
		TemplateSize:  templateSize,
		FileCount:     len(config.SelectedFiles),
		TotalSize:     totalSize,
		GeneratedAt:   startTime,
	}

	return result, nil
}

// GenerateAsync performs prompt generation asynchronously with progress updates
func (pg *PromptGenerator) GenerateAsync(config GenerationConfig, callback GenerationProgressCallback) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Create a channel to handle cancellation
		done := make(chan struct{})
		var result *GeneratedPrompt
		var err error

		go func() {
			defer close(done)
			
			// Report progress through callback if provided
			if callback != nil {
				callback(string(StageProcessingTemplate), 0.0)
			}

			// Step 1: Template processing (25% progress)
			if callback != nil {
				callback(string(StageProcessingTemplate), 0.25)
			}

			// Step 2: File loading (50% progress)  
			if callback != nil {
				callback(string(StageLoadingFiles), 0.50)
			}

			// Step 3: Structure assembly (75% progress)
			if callback != nil {
				callback(string(StageAssemblingStructure), 0.75)
			}

			// Generate the prompt
			result, err = pg.GeneratePrompt(ctx, config)

			// Step 4: Complete (100% progress)
			if callback != nil {
				if err != nil {
					callback("Generation failed", 0.75)
				} else {
					callback(string(StageComplete), 1.0)
				}
			}
		}()

		// Wait for completion or cancellation
		select {
		case <-done:
			return GenerationCompleteMsg{
				Result: result,
				Error:  err,
			}
		case <-ctx.Done():
			return GenerationCompleteMsg{
				Result: nil,
				Error:  ctx.Err(),
			}
		}
	})
}

