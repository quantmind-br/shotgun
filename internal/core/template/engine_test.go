package template

import (
	"testing"
)

func TestParseTemplateFromData_ValidTemplate(t *testing.T) {
	validTOML := `
id = "test-template"
name = "Test Template"
version = "1.0.0"
description = "A test template"
author = "Test Author"
tags = ["test", "demo"]
content = "Hello {{name}}!"

[variables.name]
name = "name"
type = "text"
required = true
placeholder = "Enter your name"
`

	template, err := parseTemplateFromData([]byte(validTOML))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if template.ID != "test-template" {
		t.Errorf("Expected ID 'test-template', got '%s'", template.ID)
	}

	if template.Name != "Test Template" {
		t.Errorf("Expected Name 'Test Template', got '%s'", template.Name)
	}

	if template.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", template.Version)
	}

	if template.Description != "A test template" {
		t.Errorf("Expected Description 'A test template', got '%s'", template.Description)
	}

	if template.Author != "Test Author" {
		t.Errorf("Expected Author 'Test Author', got '%s'", template.Author)
	}

	if len(template.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(template.Tags))
	}

	if template.Content != "Hello {{name}}!" {
		t.Errorf("Expected Content 'Hello {{name}}!', got '%s'", template.Content)
	}

	if len(template.Variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(template.Variables))
	}

	nameVar, exists := template.Variables["name"]
	if !exists {
		t.Error("Expected 'name' variable to exist")
	}

	if nameVar.Name != "name" {
		t.Errorf("Expected variable name 'name', got '%s'", nameVar.Name)
	}

	if nameVar.Type != "text" {
		t.Errorf("Expected variable type 'text', got '%s'", nameVar.Type)
	}

	if !nameVar.Required {
		t.Error("Expected variable to be required")
	}
}

func TestParseTemplateFromData_EmptyData(t *testing.T) {
	_, err := parseTemplateFromData([]byte{})
	if err == nil {
		t.Fatal("Expected error for empty data")
	}

	templateErr, ok := err.(*TemplateError)
	if !ok {
		t.Fatalf("Expected TemplateError, got %T", err)
	}

	if templateErr.Type != ErrorTypeParsing {
		t.Errorf("Expected parsing error, got %s", templateErr.Type)
	}
}

func TestParseTemplateFromData_OversizedFile(t *testing.T) {
	// Create data larger than 1MB
	largeData := make([]byte, 1024*1024+1)
	for i := range largeData {
		largeData[i] = 'a'
	}

	_, err := parseTemplateFromData(largeData)
	if err == nil {
		t.Fatal("Expected error for oversized file")
	}

	templateErr, ok := err.(*TemplateError)
	if !ok {
		t.Fatalf("Expected TemplateError, got %T", err)
	}

	if templateErr.Type != ErrorTypeContentSize {
		t.Errorf("Expected content size error, got %s", templateErr.Type)
	}
}

func TestParseTemplateFromData_MalformedTOML(t *testing.T) {
	malformedTOML := `
name = "Test Template"
version = 1.0.0"  # Missing opening quote
description = "A test template
content = "Hello {{name}}!"
`

	_, err := parseTemplateFromData([]byte(malformedTOML))
	if err == nil {
		t.Fatal("Expected error for malformed TOML")
	}

	templateErr, ok := err.(*TemplateError)
	if !ok {
		t.Fatalf("Expected TemplateError, got %T", err)
	}

	if templateErr.Type != ErrorTypeParsing {
		t.Errorf("Expected parsing error, got %s", templateErr.Type)
	}
}

func TestParseTemplateFromData_ValidationFailure(t *testing.T) {
	invalidTOML := `
name = "Test Template"
# Missing required version and description
content = "Hello {{name}}!"
`

	_, err := parseTemplateFromData([]byte(invalidTOML))
	if err == nil {
		t.Fatal("Expected validation error")
	}

	templateErr, ok := err.(*TemplateError)
	if !ok {
		t.Fatalf("Expected TemplateError, got %T", err)
	}

	if templateErr.Type != ErrorTypeValidation {
		t.Errorf("Expected validation error, got %s", templateErr.Type)
	}
}

func TestParseTemplateFromData_GenerateIDFromName(t *testing.T) {
	tomlWithoutID := `
name = "My Test Template"
version = "1.0.0"
description = "A test template"
content = "Hello {{name}}!"
`

	template, err := parseTemplateFromData([]byte(tomlWithoutID))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedID := "my-test-template"
	if template.ID != expectedID {
		t.Errorf("Expected generated ID '%s', got '%s'", expectedID, template.ID)
	}
}

func TestGenerateTemplateID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Template", "simple-template"},
		{"My_Special Template!", "my-special-template"},
		{"   Whitespace   Template   ", "whitespace-template"},
		{"", "unnamed-template"},
		{"123 Numbers", "123-numbers"},
		{"UPPERCASE", "uppercase"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := generateTemplateID(test.input)
			if result != test.expected {
				t.Errorf("For input '%s', expected '%s', got '%s'", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestParseTemplateFromData_ComplexVariables(t *testing.T) {
	complexTOML := `
id = "complex-template"
name = "Complex Template"
version = "2.0.0"
description = "A complex test template"
content = "{{message}} Priority: {{priority}}"

[variables.message]
name = "message"
type = "multiline"
required = true
placeholder = "Enter your message"
min_length = 10
max_length = 500

[variables.priority]
name = "priority"
type = "choice"
required = false
default = "medium"
options = ["low", "medium", "high", "critical"]
`

	template, err := parseTemplateFromData([]byte(complexTOML))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(template.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(template.Variables))
	}

	// Test message variable
	messageVar, exists := template.Variables["message"]
	if !exists {
		t.Error("Expected 'message' variable to exist")
	} else {
		if messageVar.Type != "multiline" {
			t.Errorf("Expected message type 'multiline', got '%s'", messageVar.Type)
		}
		if messageVar.MinLength != 10 {
			t.Errorf("Expected MinLength 10, got %d", messageVar.MinLength)
		}
		if messageVar.MaxLength != 500 {
			t.Errorf("Expected MaxLength 500, got %d", messageVar.MaxLength)
		}
	}

	// Test priority variable
	priorityVar, exists := template.Variables["priority"]
	if !exists {
		t.Error("Expected 'priority' variable to exist")
	} else {
		if priorityVar.Type != "choice" {
			t.Errorf("Expected priority type 'choice', got '%s'", priorityVar.Type)
		}
		if priorityVar.Default != "medium" {
			t.Errorf("Expected default 'medium', got '%s'", priorityVar.Default)
		}
		if len(priorityVar.Options) != 4 {
			t.Errorf("Expected 4 options, got %d", len(priorityVar.Options))
		}
	}
}