package template

import (
	"context"
	"strings"
	"testing"

	"github.com/diogopedro/shotgun/internal/models"
)

func TestNewTemplateEngine(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		engine := NewTemplateEngine()
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("with custom options", func(t *testing.T) {
		engine := NewTemplateEngine(
			WithStrictMode(false),
			WithMaxSize(500),
			WithAllowedFunctions([]string{"upper", "lower"}),
		)
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})
}

func TestTemplateEngine_ProcessTemplate(t *testing.T) {
	ctx := context.Background()
	engine := NewTemplateEngine()

	tests := []struct {
		name      string
		template  *models.Template
		vars      map[string]interface{}
		expected  string
		shouldErr bool
	}{
		{
			name: "simple variable substitution",
			template: &models.Template{
				ID:      "simple",
				Content: "Hello {{.name}}!",
			},
			vars: map[string]interface{}{
				"name": "World",
			},
			expected: "Hello World!",
		},
		{
			name: "auto variables",
			template: &models.Template{
				ID:      "auto",
				Content: "Date: {{.CURRENT_DATE}}, Project: {{.PROJECT_NAME}}",
			},
			vars:     map[string]interface{}{},
			expected: "auto-variables-check", // Special case handled in test
		},
		{
			name: "template with functions",
			template: &models.Template{
				ID:      "functions",
				Content: "{{upper .message}} - {{lower .NAME}}",
			},
			vars: map[string]interface{}{
				"message": "hello",
				"NAME":    "WORLD",
			},
			expected: "HELLO - world",
		},
		{
			name: "conditional processing",
			template: &models.Template{
				ID:      "conditional",
				Content: "{{if .show}}Visible{{else}}Hidden{{end}}",
			},
			vars: map[string]interface{}{
				"show": true,
			},
			expected: "Visible",
		},
		{
			name: "conditional false",
			template: &models.Template{
				ID:      "conditional-false",
				Content: "{{if .show}}Visible{{else}}Hidden{{end}}",
			},
			vars: map[string]interface{}{
				"show": false,
			},
			expected: "Hidden",
		},
		{
			name: "complex template with nested conditions",
			template: &models.Template{
				ID:      "complex",
				Content: "{{if .user.admin}}Admin: {{upper .user.name}}{{else}}User: {{.user.name}}{{end}}",
			},
			vars: map[string]interface{}{
				"user": map[string]interface{}{
					"admin": true,
					"name":  "john",
				},
			},
			expected: "Admin: JOHN",
		},
		{
			name: "empty template content",
			template: &models.Template{
				ID:      "empty",
				Content: "",
			},
			vars:      map[string]interface{}{},
			shouldErr: true,
		},
		{
			name: "malformed template",
			template: &models.Template{
				ID:      "malformed",
				Content: "{{if .condition}}Missing end tag",
			},
			vars:      map[string]interface{}{},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ProcessTemplate(ctx, tt.template, tt.vars)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.name == "auto variables" {
				// Special case for auto variables - check that result contains expected parts
				if !strings.Contains(result, "Date: ") || !strings.Contains(result, "Project: ") {
					t.Errorf("Auto variables test failed: %s", result)
				}
			} else if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateEngine_ProcessTemplate_StrictMode(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		strict     bool
		template   *models.Template
		vars       map[string]interface{}
		shouldErr  bool
		errMessage string
	}{
		{
			name:   "strict mode with missing required variable",
			strict: true,
			template: &models.Template{
				ID:      "strict",
				Content: "Hello {{.name}}!",
				Variables: map[string]models.Variable{
					"name": {
						Name:     "name",
						Required: true,
					},
				},
			},
			vars:       map[string]interface{}{},
			shouldErr:  true,
			errMessage: "required variable 'name' is missing",
		},
		{
			name:   "non-strict mode with missing required variable",
			strict: false,
			template: &models.Template{
				ID:      "non-strict",
				Content: "Hello {{.name}}!",
				Variables: map[string]models.Variable{
					"name": {
						Name:     "name",
						Required: true,
					},
				},
			},
			vars:      map[string]interface{}{},
			shouldErr: false,
		},
		{
			name:   "default value substitution",
			strict: true,
			template: &models.Template{
				ID:      "default",
				Content: "Hello {{.name}}!",
				Variables: map[string]models.Variable{
					"name": {
						Name:    "name",
						Default: "Anonymous",
					},
				},
			},
			vars:      map[string]interface{}{},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewTemplateEngine(WithStrictMode(tt.strict))
			result, err := engine.ProcessTemplate(ctx, tt.template, tt.vars)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errMessage != "" && !strings.Contains(err.Error(), tt.errMessage) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.name == "default value substitution" && result != "Hello Anonymous!" {
				t.Errorf("Default value test failed. Expected 'Hello Anonymous!', got %q", result)
			}
		})
	}
}

func TestTemplateEngine_RegisterFunction(t *testing.T) {
	engine := NewTemplateEngine()

	// Test successful registration
	err := engine.RegisterFunction("reverse", func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test usage of custom function
	ctx := context.Background()
	template := &models.Template{
		ID:      "custom-func",
		Content: "{{reverse .word}}",
	}
	vars := map[string]interface{}{
		"word": "hello",
	}

	result, err := engine.ProcessTemplate(ctx, template, vars)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "olleh" {
		t.Errorf("Expected 'olleh', got %q", result)
	}

	// Test empty function name
	err = engine.RegisterFunction("", func() {})
	if err == nil {
		t.Error("Expected error for empty function name")
	}
}

func TestTemplateEngine_RegisterFunction_AllowedFunctions(t *testing.T) {
	engine := NewTemplateEngine(WithAllowedFunctions([]string{"upper", "custom"}))

	// Test allowed function
	err := engine.RegisterFunction("custom", func(s string) string { return s })
	if err != nil {
		t.Errorf("Expected no error for allowed function, got %v", err)
	}

	// Test disallowed function
	err = engine.RegisterFunction("disallowed", func(s string) string { return s })
	if err == nil {
		t.Error("Expected error for disallowed function")
	}
	if !strings.Contains(err.Error(), "not in allowed functions list") {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestTemplateEngine_ValidateTemplate(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "valid template",
			content:   "Hello {{.name}}!",
			shouldErr: false,
		},
		{
			name:      "valid conditional template",
			content:   "{{if .show}}Visible{{end}}",
			shouldErr: false,
		},
		{
			name:      "empty template",
			content:   "",
			shouldErr: true,
		},
		{
			name:      "malformed template",
			content:   "{{if .condition}}Missing end",
			shouldErr: true,
		},
		{
			name:      "invalid syntax",
			content:   "{{.name}",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.ValidateTemplate(tt.content)
			if tt.shouldErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestTemplateEngine_BuiltinFunctions(t *testing.T) {
	ctx := context.Background()
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		vars     map[string]interface{}
		expected string
	}{
		{
			name:     "upper function",
			template: "{{upper .text}}",
			vars:     map[string]interface{}{"text": "hello"},
			expected: "HELLO",
		},
		{
			name:     "lower function",
			template: "{{lower .text}}",
			vars:     map[string]interface{}{"text": "HELLO"},
			expected: "hello",
		},
		{
			name:     "trim function",
			template: "{{trim .text}}",
			vars:     map[string]interface{}{"text": "  hello  "},
			expected: "hello",
		},
		{
			name:     "replace function",
			template: "{{replace \"old\" \"new\" .text 1}}",
			vars:     map[string]interface{}{"text": "old text old"},
			expected: "new text old",
		},
		{
			name:     "replaceAll function",
			template: "{{replaceAll \"old\" \"new\" .text}}",
			vars:     map[string]interface{}{"text": "old text old"},
			expected: "new text new",
		},
		{
			name:     "contains function",
			template: "{{if contains .text \"hello\"}}Yes{{else}}No{{end}}",
			vars:     map[string]interface{}{"text": "hello world"},
			expected: "Yes",
		},
		{
			name:     "default function",
			template: "{{default \"fallback\" .missing}}",
			vars:     map[string]interface{}{},
			expected: "fallback",
		},
		{
			name:     "ternary function",
			template: "{{ternary .condition \"yes\" \"no\"}}",
			vars:     map[string]interface{}{"condition": true},
			expected: "yes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := &models.Template{
				ID:      tt.name,
				Content: tt.template,
			}

			result, err := engine.ProcessTemplate(ctx, tmpl, tt.vars)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateEngine_MaxSize(t *testing.T) {
	ctx := context.Background()
	engine := NewTemplateEngine(WithMaxSize(10)) // Very small limit

	template := &models.Template{
		ID:      "large",
		Content: "{{.text}}",
	}
	vars := map[string]interface{}{
		"text": "This is a very long text that exceeds the limit",
	}

	_, err := engine.ProcessTemplate(ctx, template, vars)
	if err == nil {
		t.Error("Expected size limit error")
	}
	if !strings.Contains(err.Error(), "exceeds size limit") {
		t.Errorf("Expected size limit error message, got %v", err)
	}
}

func TestTemplateEngine_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	engine := NewTemplateEngine()
	template := &models.Template{
		ID:      "cancelled",
		Content: "Hello {{.name}}!",
	}
	vars := map[string]interface{}{
		"name": "World",
	}

	_, err := engine.ProcessTemplate(ctx, template, vars)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestTemplateEngine_ComplexScenarios(t *testing.T) {
	ctx := context.Background()
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template *models.Template
		vars     map[string]interface{}
		expected string
	}{
		{
			name: "nested conditionals with functions",
			template: &models.Template{
				ID: "nested",
				Content: `{{if .user.admin}}Admin Panel: {{upper .user.name}}
{{if .user.permissions.write}}Can Write{{end}}{{else}}User: {{.user.name}}{{end}}`,
			},
			vars: map[string]interface{}{
				"user": map[string]interface{}{
					"admin": true,
					"name":  "john",
					"permissions": map[string]interface{}{
						"write": true,
					},
				},
			},
			expected: "Admin Panel: JOHN\nCan Write",
		},
		{
			name: "template with ranges and functions",
			template: &models.Template{
				ID:      "range",
				Content: "{{range .items}}{{upper .}}|{{end}}",
			},
			vars: map[string]interface{}{
				"items": []string{"hello", "world"},
			},
			expected: "HELLO|WORLD|",
		},
		{
			name: "all variable types",
			template: &models.Template{
				ID: "all-vars",
				Content: `Task: {{.TASK}}
Rules: {{.RULES}}
Files: {{.FILE_STRUCTURE}}
Date: {{.CURRENT_DATE}}
Project: {{.PROJECT_NAME}}`,
				Variables: map[string]models.Variable{
					"TASK": {Name: "TASK", Required: true},
					"RULES": {Name: "RULES", Default: "No rules"},
					"FILE_STRUCTURE": {Name: "FILE_STRUCTURE", Default: "No files"},
				},
			},
			vars: map[string]interface{}{
				"TASK":           "Test task",
				"FILE_STRUCTURE": "src/\n  main.go",
			},
			expected: "all-vars-check", // Special case handled in test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ProcessTemplate(ctx, tt.template, tt.vars)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Handle auto variables specially
			if tt.name == "all variable types" {
				if !strings.Contains(result, "Task: Test task") ||
					!strings.Contains(result, "Rules: No rules") ||
					!strings.Contains(result, "Files: src/\n  main.go") ||
					!strings.Contains(result, "Date: ") ||
					!strings.Contains(result, "Project: ") {
					t.Errorf("All variables test failed: %s", result)
				}
			} else if result != tt.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
			}
		})
	}
}

// Benchmark tests for performance verification
func BenchmarkTemplateEngine_ProcessTemplate(b *testing.B) {
	ctx := context.Background()
	engine := NewTemplateEngine()
	template := &models.Template{
		ID:      "bench",
		Content: "Hello {{.name}}! Today is {{.CURRENT_DATE}}. {{upper .message}}",
	}
	vars := map[string]interface{}{
		"name":    "World",
		"message": "this is a test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ProcessTemplate(ctx, template, vars)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkTemplateEngine_ComplexTemplate(b *testing.B) {
	ctx := context.Background()
	engine := NewTemplateEngine()
	template := &models.Template{
		ID: "complex-bench",
		Content: `{{range .users}}
User: {{upper .name}}
{{if .admin}}Admin: {{.permissions}}{{end}}
{{end}}`,
	}
	vars := map[string]interface{}{
		"users": []map[string]interface{}{
			{"name": "alice", "admin": true, "permissions": "read,write"},
			{"name": "bob", "admin": false},
			{"name": "charlie", "admin": true, "permissions": "read"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ProcessTemplate(ctx, template, vars)
		if err != nil {
			b.Fatalf("Complex benchmark failed: %v", err)
		}
	}
}