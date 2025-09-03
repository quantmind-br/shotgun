package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestTemplate_TOMLMarshaling(t *testing.T) {
	template := Template{
		ID:          "test-template",
		Name:        "Test Template",
		Version:     "1.0.0",
		Description: "A test template",
		Author:      "Test Author",
		Tags:        []string{"test", "example"},
		Variables: map[string]Variable{
			"task_description": {
				Name:        "task_description",
				Type:        "multiline",
				Required:    true,
				Placeholder: "Describe the task...",
			},
			"urgency": {
				Name:     "urgency",
				Type:     "choice",
				Required: false,
				Default:  "medium",
				Options:  []string{"low", "medium", "high", "critical"},
			},
		},
		Content: "# Task: {{task_description}}\n\nUrgency: {{urgency}}",
	}

	// Test TOML marshaling
	var buf strings.Builder
	err := toml.NewEncoder(&buf).Encode(template)
	if err != nil {
		t.Fatalf("Failed to encode TOML: %v", err)
	}

	// Test TOML unmarshaling
	var decoded Template
	err = toml.Unmarshal([]byte(buf.String()), &decoded)
	if err != nil {
		t.Fatalf("Failed to decode TOML: %v", err)
	}

	// Verify fields
	if decoded.ID != template.ID {
		t.Errorf("Expected ID %s, got %s", template.ID, decoded.ID)
	}
	if decoded.Name != template.Name {
		t.Errorf("Expected Name %s, got %s", template.Name, decoded.Name)
	}
	if len(decoded.Variables) != len(template.Variables) {
		t.Errorf("Expected %d variables, got %d", len(template.Variables), len(decoded.Variables))
	}
}

func TestTemplate_JSONMarshaling(t *testing.T) {
	template := Template{
		ID:          "json-test",
		Name:        "JSON Test Template",
		Version:     "2.0.0",
		Description: "A template for JSON testing",
		Author:      "JSON Tester",
		Tags:        []string{"json", "test"},
		Variables: map[string]Variable{
			"name": {
				Name:        "name",
				Type:        "text",
				Required:    true,
				Placeholder: "Enter your name",
				MinLength:   1,
				MaxLength:   100,
			},
		},
		Content: "Hello {{name}}!",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Test JSON unmarshaling
	var decoded Template
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify fields
	if decoded.ID != template.ID {
		t.Errorf("Expected ID %s, got %s", template.ID, decoded.ID)
	}
	if decoded.Variables["name"].MinLength != template.Variables["name"].MinLength {
		t.Errorf("Expected MinLength %d, got %d", template.Variables["name"].MinLength, decoded.Variables["name"].MinLength)
	}
}

func TestVariable_Types(t *testing.T) {
	tests := []struct {
		name     string
		variable Variable
		valid    bool
	}{
		{
			name: "valid text type",
			variable: Variable{
				Name:     "test",
				Type:     "text",
				Required: true,
			},
			valid: true,
		},
		{
			name: "valid choice type with options",
			variable: Variable{
				Name:     "choice_var",
				Type:     "choice",
				Required: false,
				Options:  []string{"option1", "option2"},
			},
			valid: true,
		},
		{
			name: "valid number type with constraints",
			variable: Variable{
				Name:      "number_var",
				Type:      "number",
				Required:  true,
				MinLength: 1,
				MaxLength: 10,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation that the variable can be created
			if tt.variable.Name == "" && tt.valid {
				t.Error("Valid variable should have a name")
			}

			// Check if type is in valid types
			isValidType := false
			for _, validType := range ValidVariableTypes {
				if tt.variable.Type == validType {
					isValidType = true
					break
				}
			}

			if tt.valid && !isValidType {
				t.Errorf("Type %s should be valid", tt.variable.Type)
			}
		})
	}
}

func TestTemplateSource_String(t *testing.T) {
	tests := []struct {
		source   TemplateSource
		expected string
	}{
		{TemplateSourceBuiltIn, "builtin"},
		{TemplateSourceUser, "user"},
		{TemplateSource(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.source.String(); got != tt.expected {
				t.Errorf("TemplateSource.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTemplateInfo(t *testing.T) {
	template := Template{
		ID:      "info-test",
		Name:    "Info Test",
		Version: "1.0.0",
	}

	info := TemplateInfo{
		Template: template,
		Source:   TemplateSourceBuiltIn,
		FilePath: "/embedded/test.toml",
	}

	if info.Template.ID != template.ID {
		t.Errorf("Expected template ID %s, got %s", template.ID, info.Template.ID)
	}

	if info.Source != TemplateSourceBuiltIn {
		t.Errorf("Expected source %v, got %v", TemplateSourceBuiltIn, info.Source)
	}

	if info.Source.String() != "builtin" {
		t.Errorf("Expected source string 'builtin', got %s", info.Source.String())
	}
}
