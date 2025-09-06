package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/diogopedro/shotgun/internal/core/builder"
	"github.com/diogopedro/shotgun/internal/models"
)

// TemplateEngine interface defines template processing operations
type TemplateEngine interface {
	ProcessTemplate(ctx context.Context, tmpl *models.Template, vars map[string]interface{}) (string, error)
	ProcessTemplateWithFiles(ctx context.Context, tmpl *models.Template, vars map[string]interface{}, selectedFiles []string) (string, error)
	RegisterFunction(name string, fn interface{}) error
	ValidateTemplate(content string) error
}

// templateEngine implements the TemplateEngine interface
type templateEngine struct {
	funcMap           sync.Map // Custom template functions
	options           ProcessingOptions
	fileStructBuilder *builder.FileStructureBuilder
}

// ProcessingOptions configures template processing behavior
type ProcessingOptions struct {
	StrictMode   bool     // Fail on missing variables
	AllowedFuncs []string // Restrict available functions
	MaxSize      int64    // Maximum output size
}

// NewTemplateEngine creates a new template engine with optional configuration
func NewTemplateEngine(options ...func(*ProcessingOptions)) TemplateEngine {
	opts := ProcessingOptions{
		StrictMode: true,
		MaxSize:    1024 * 1024, // 1MB default
	}

	for _, opt := range options {
		opt(&opts)
	}

	engine := &templateEngine{
		options:           opts,
		fileStructBuilder: builder.NewFileStructureBuilder(),
	}

	// Register built-in functions
	engine.registerBuiltinFunctions()

	return engine
}

// ProcessTemplate executes template processing with variable substitution
func (e *templateEngine) ProcessTemplate(ctx context.Context, tmpl *models.Template, vars map[string]interface{}) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Validate template content
	if err := e.ValidateTemplate(tmpl.Content); err != nil {
		return "", fmt.Errorf("template validation failed: %w", err)
	}

	// Prepare variables with auto-populated values
	processedVars, err := e.prepareVariables(tmpl, vars)
	if err != nil {
		return "", fmt.Errorf("variable preparation failed: %w", err)
	}

	// Create function map
	funcMap := e.createFunctionMap()

	// Create and parse template
	goTemplate, err := template.New(tmpl.ID).Funcs(funcMap).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	// Execute template
	var result strings.Builder
	if err := goTemplate.Execute(&result, processedVars); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	// Check output size limit
	output := result.String()
	if e.options.MaxSize > 0 && int64(len(output)) > e.options.MaxSize {
		return "", fmt.Errorf("template output exceeds size limit of %d bytes", e.options.MaxSize)
	}

	return output, nil
}

// ProcessTemplateWithFiles executes template processing with selected files for FILE_STRUCTURE variable
func (e *templateEngine) ProcessTemplateWithFiles(ctx context.Context, tmpl *models.Template, vars map[string]interface{}, selectedFiles []string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Validate template content
	if err := e.ValidateTemplate(tmpl.Content); err != nil {
		return "", fmt.Errorf("template validation failed: %w", err)
	}

	// Prepare variables with auto-populated values including FILE_STRUCTURE
	processedVars, err := e.prepareVariablesWithFiles(tmpl, vars, selectedFiles)
	if err != nil {
		return "", fmt.Errorf("variable preparation failed: %w", err)
	}

	// Create function map
	funcMap := e.createFunctionMap()

	// Create and parse template
	goTemplate, err := template.New(tmpl.ID).Funcs(funcMap).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	// Execute template
	var result strings.Builder
	if err := goTemplate.Execute(&result, processedVars); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	// Check output size limit
	output := result.String()
	if e.options.MaxSize > 0 && int64(len(output)) > e.options.MaxSize {
		return "", fmt.Errorf("template output exceeds size limit of %d bytes", e.options.MaxSize)
	}

	return output, nil
}

// RegisterFunction adds a custom function to the template engine
func (e *templateEngine) RegisterFunction(name string, fn interface{}) error {
	if name == "" {
		return fmt.Errorf("function name cannot be empty")
	}

	// Check if function is allowed if restrictions are in place
	if len(e.options.AllowedFuncs) > 0 {
		allowed := false
		for _, allowedFunc := range e.options.AllowedFuncs {
			if allowedFunc == name {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("function '%s' is not in allowed functions list", name)
		}
	}

	e.funcMap.Store(name, fn)
	return nil
}

// ValidateTemplate checks if template syntax is valid
func (e *templateEngine) ValidateTemplate(content string) error {
	if content == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Create function map for validation
	funcMap := e.createFunctionMap()

	// Try to parse the template to validate syntax
	_, err := template.New("validation").Funcs(funcMap).Parse(content)
	if err != nil {
		return fmt.Errorf("invalid template syntax: %w", err)
	}

	return nil
}

// prepareVariables combines user variables with auto-generated ones
func (e *templateEngine) prepareVariables(tmpl *models.Template, userVars map[string]interface{}) (map[string]interface{}, error) {
	processedVars := make(map[string]interface{})

	// Copy user variables
	for k, v := range userVars {
		processedVars[k] = v
	}

	// Add auto-populated variables
	processedVars["CURRENT_DATE"] = time.Now().Format("2006-01-02")

	// Get project name from current directory
	if wd, err := os.Getwd(); err == nil {
		processedVars["PROJECT_NAME"] = filepath.Base(wd)
	}

	// Validate required variables in strict mode
	if e.options.StrictMode {
		for varName, variable := range tmpl.Variables {
			if variable.Required {
				if _, exists := processedVars[varName]; !exists {
					return nil, fmt.Errorf("required variable '%s' is missing", varName)
				}
			}
		}
	}

	// Apply defaults for missing variables
	for varName, variable := range tmpl.Variables {
		if _, exists := processedVars[varName]; !exists && variable.Default != "" {
			processedVars[varName] = variable.Default
		}
	}

	return processedVars, nil
}

// prepareVariablesWithFiles extends prepareVariables to include FILE_STRUCTURE generation
func (e *templateEngine) prepareVariablesWithFiles(tmpl *models.Template, userVars map[string]interface{}, selectedFiles []string) (map[string]interface{}, error) {
	processedVars := make(map[string]interface{})

	// Copy user variables
	for k, v := range userVars {
		processedVars[k] = v
	}

	// Add auto-populated variables
	processedVars["CURRENT_DATE"] = time.Now().Format("2006-01-02")

	// Get project name from current directory
	if wd, err := os.Getwd(); err == nil {
		processedVars["PROJECT_NAME"] = filepath.Base(wd)
	}

	// Generate FILE_STRUCTURE if selected files provided and not already set
	if _, exists := processedVars["FILE_STRUCTURE"]; !exists && len(selectedFiles) > 0 {
		ctx := context.Background() // Use background context for file structure generation
		fileStructure, err := e.fileStructBuilder.GenerateStructure(ctx, selectedFiles)
		if err != nil {
			return nil, fmt.Errorf("failed to generate file structure: %w", err)
		}
		processedVars["FILE_STRUCTURE"] = fileStructure
	}

	// Validate required variables in strict mode
	if e.options.StrictMode {
		for varName, variable := range tmpl.Variables {
			if variable.Required {
				if _, exists := processedVars[varName]; !exists {
					return nil, fmt.Errorf("required variable '%s' is missing", varName)
				}
			}
		}
	}

	// Apply defaults for missing variables
	for varName, variable := range tmpl.Variables {
		if _, exists := processedVars[varName]; !exists && variable.Default != "" {
			processedVars[varName] = variable.Default
		}
	}

	return processedVars, nil
}

// createFunctionMap creates the function map for template execution
func (e *templateEngine) createFunctionMap() template.FuncMap {
	funcMap := make(template.FuncMap)

	// Add all registered functions
	e.funcMap.Range(func(key, value interface{}) bool {
		name := key.(string)
		fn := value
		funcMap[name] = fn
		return true
	})

	return funcMap
}

// registerBuiltinFunctions registers the default template functions
func (e *templateEngine) registerBuiltinFunctions() {
	// String manipulation functions
	e.funcMap.Store("upper", strings.ToUpper)
	e.funcMap.Store("lower", strings.ToLower)
	e.funcMap.Store("trim", strings.TrimSpace)
	e.funcMap.Store("trimLeft", func(cutset, s string) string {
		return strings.TrimLeft(s, cutset)
	})
	e.funcMap.Store("trimRight", func(cutset, s string) string {
		return strings.TrimRight(s, cutset)
	})
	e.funcMap.Store("replace", func(old, new, s string, n int) string {
		return strings.Replace(s, old, new, n)
	})
	e.funcMap.Store("replaceAll", func(old, new, s string) string {
		return strings.ReplaceAll(s, old, new)
	})
	e.funcMap.Store("split", func(sep, s string) []string {
		return strings.Split(s, sep)
	})
	e.funcMap.Store("join", func(sep string, elems []string) string {
		return strings.Join(elems, sep)
	})
	e.funcMap.Store("contains", strings.Contains)
	e.funcMap.Store("hasPrefix", strings.HasPrefix)
	e.funcMap.Store("hasSuffix", strings.HasSuffix)

	// Utility functions
	e.funcMap.Store("default", func(defaultValue, value interface{}) interface{} {
		if value == nil || value == "" {
			return defaultValue
		}
		return value
	})

	// Conditional helper
	e.funcMap.Store("ternary", func(condition bool, trueValue, falseValue interface{}) interface{} {
		if condition {
			return trueValue
		}
		return falseValue
	})
}

// Option functions for engine configuration

// WithStrictMode enables or disables strict variable validation
func WithStrictMode(strict bool) func(*ProcessingOptions) {
	return func(opts *ProcessingOptions) {
		opts.StrictMode = strict
	}
}

// WithMaxSize sets the maximum output size
func WithMaxSize(maxSize int64) func(*ProcessingOptions) {
	return func(opts *ProcessingOptions) {
		opts.MaxSize = maxSize
	}
}

// WithAllowedFunctions restricts which functions can be used
func WithAllowedFunctions(funcs []string) func(*ProcessingOptions) {
	return func(opts *ProcessingOptions) {
		opts.AllowedFuncs = funcs
	}
}
