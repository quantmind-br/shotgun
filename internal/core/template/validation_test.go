package template

import (
	"testing"

	"github.com/user/shotgun-cli/internal/models"
)

func TestValidateTemplate_ValidTemplate(t *testing.T) {
	template := &models.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Version:     "1.0.0",
		Description: "A valid test template",
		Content:     "Hello {{name}}!",
		Variables: map[string]models.Variable{
			"name": {
				Name:     "name",
				Type:     "text",
				Required: true,
			},
		},
	}

	err := validateTemplate(template)
	if err != nil {
		t.Errorf("Expected no error for valid template, got: %v", err)
	}
}

func TestValidateTemplate_NilTemplate(t *testing.T) {
	err := validateTemplate(nil)
	if err == nil {
		t.Error("Expected error for nil template")
	}
}

func TestValidateRequiredFields_MissingName(t *testing.T) {
	template := &models.Template{
		ID:          "test-template",
		Name:        "", // Missing name
		Version:     "1.0.0",
		Description: "A test template",
		Content:     "Hello world!",
	}

	err := validateRequiredFields(template)
	if err == nil {
		t.Error("Expected error for missing name")
	}
}

func TestValidateRequiredFields_MissingVersion(t *testing.T) {
	template := &models.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Version:     "", // Missing version
		Description: "A test template",
		Content:     "Hello world!",
	}

	err := validateRequiredFields(template)
	if err == nil {
		t.Error("Expected error for missing version")
	}
}

func TestValidateRequiredFields_MissingDescription(t *testing.T) {
	template := &models.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Version:     "1.0.0",
		Description: "", // Missing description
		Content:     "Hello world!",
	}

	err := validateRequiredFields(template)
	if err == nil {
		t.Error("Expected error for missing description")
	}
}

func TestValidateRequiredFields_MissingContent(t *testing.T) {
	template := &models.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Version:     "1.0.0",
		Description: "A test template",
		Content:     "", // Missing content
	}

	err := validateRequiredFields(template)
	if err == nil {
		t.Error("Expected error for missing content")
	}
}

func TestValidateRequiredFields_MissingID(t *testing.T) {
	template := &models.Template{
		ID:          "", // Missing ID
		Name:        "Test Template",
		Version:     "1.0.0",
		Description: "A test template",
		Content:     "Hello world!",
	}

	err := validateRequiredFields(template)
	if err == nil {
		t.Error("Expected error for missing ID")
	}
}

func TestIsValidVersion(t *testing.T) {
	tests := []struct {
		version string
		valid   bool
	}{
		{"1.0.0", true},
		{"0.1.0", true},
		{"10.20.30", true},
		{"1.0", false},     // Not 3 parts
		{"1.0.0.0", false}, // Too many parts
		{"1.0.a", false},   // Non-numeric
		{"", false},        // Empty
		{"1..0", false},    // Empty part
		{"1.0.", false},    // Trailing dot
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			result := isValidVersion(test.version)
			if result != test.valid {
				t.Errorf("For version '%s', expected %v, got %v",
					test.version, test.valid, result)
			}
		})
	}
}

func TestValidateVariable_ValidVariable(t *testing.T) {
	variable := models.Variable{
		Name:        "test_var",
		Type:        "text",
		Required:    true,
		Placeholder: "Enter text",
	}

	err := validateVariable("test_var", variable)
	if err != nil {
		t.Errorf("Expected no error for valid variable, got: %v", err)
	}
}

func TestValidateVariable_EmptyName(t *testing.T) {
	variable := models.Variable{
		Name: "",
		Type: "text",
	}

	err := validateVariable("", variable)
	if err == nil {
		t.Error("Expected error for empty variable name")
	}
}

func TestValidateVariable_NameMismatch(t *testing.T) {
	variable := models.Variable{
		Name: "different_name",
		Type: "text",
	}

	err := validateVariable("test_var", variable)
	if err == nil {
		t.Error("Expected error for name mismatch")
	}
}

func TestIsValidVariableType(t *testing.T) {
	// First check what valid types are defined
	validTypes := models.ValidVariableTypes
	if len(validTypes) == 0 {
		t.Skip("No valid variable types defined in models")
	}

	// Test valid types
	for _, validType := range validTypes {
		if !isValidVariableType(validType) {
			t.Errorf("Expected '%s' to be valid", validType)
		}
	}

	// Test invalid types
	invalidTypes := []string{"invalid", "unknown", ""}
	for _, invalidType := range invalidTypes {
		if isValidVariableType(invalidType) {
			t.Errorf("Expected '%s' to be invalid", invalidType)
		}
	}
}

func TestValidateVariable_ChoiceType(t *testing.T) {
	// Valid choice variable
	validChoice := models.Variable{
		Name:    "priority",
		Type:    "choice",
		Options: []string{"low", "medium", "high"},
		Default: "medium",
	}

	err := validateVariable("priority", validChoice)
	if err != nil {
		t.Errorf("Expected no error for valid choice variable, got: %v", err)
	}

	// Choice without options
	invalidChoice := models.Variable{
		Name:    "priority",
		Type:    "choice",
		Options: []string{}, // Empty options
	}

	err = validateVariable("priority", invalidChoice)
	if err == nil {
		t.Error("Expected error for choice variable without options")
	}

	// Choice with invalid default
	invalidDefault := models.Variable{
		Name:    "priority",
		Type:    "choice",
		Options: []string{"low", "medium", "high"},
		Default: "invalid", // Not in options
	}

	err = validateVariable("priority", invalidDefault)
	if err == nil {
		t.Error("Expected error for choice variable with invalid default")
	}
}

func TestValidateVariable_BooleanType(t *testing.T) {
	// Valid boolean defaults
	validDefaults := []string{"", "true", "false"}

	for _, defaultVal := range validDefaults {
		boolean := models.Variable{
			Name:    "enabled",
			Type:    "boolean",
			Default: defaultVal,
		}

		err := validateVariable("enabled", boolean)
		if err != nil {
			t.Errorf("Expected no error for boolean with default '%s', got: %v",
				defaultVal, err)
		}
	}

	// Invalid boolean default
	invalidBoolean := models.Variable{
		Name:    "enabled",
		Type:    "boolean",
		Default: "yes", // Invalid boolean
	}

	err := validateVariable("enabled", invalidBoolean)
	if err == nil {
		t.Error("Expected error for boolean variable with invalid default")
	}
}

func TestValidateVariable_LengthConstraints(t *testing.T) {
	// Valid length constraints
	validVar := models.Variable{
		Name:      "text_field",
		Type:      "text",
		MinLength: 5,
		MaxLength: 100,
		Default:   "valid text",
	}

	err := validateVariable("text_field", validVar)
	if err != nil {
		t.Errorf("Expected no error for valid length constraints, got: %v", err)
	}

	// Negative MinLength
	negativeMin := models.Variable{
		Name:      "text_field",
		Type:      "text",
		MinLength: -1,
	}

	err = validateVariable("text_field", negativeMin)
	if err == nil {
		t.Error("Expected error for negative MinLength")
	}

	// MaxLength less than MinLength
	invalidRange := models.Variable{
		Name:      "text_field",
		Type:      "text",
		MinLength: 100,
		MaxLength: 50,
	}

	err = validateVariable("text_field", invalidRange)
	if err == nil {
		t.Error("Expected error for MaxLength < MinLength")
	}

	// Default too short
	tooShort := models.Variable{
		Name:      "text_field",
		Type:      "text",
		MinLength: 10,
		Default:   "short",
	}

	err = validateVariable("text_field", tooShort)
	if err == nil {
		t.Error("Expected error for default value too short")
	}

	// Default too long
	tooLong := models.Variable{
		Name:      "text_field",
		Type:      "text",
		MaxLength: 5,
		Default:   "this is too long",
	}

	err = validateVariable("text_field", tooLong)
	if err == nil {
		t.Error("Expected error for default value too long")
	}
}

func TestValidateContent(t *testing.T) {
	// Valid content
	validContent := []string{
		"Simple text",
		"Text with {{variable}}",
		"Multiple {{var1}} and {{var2}} placeholders",
	}

	for _, content := range validContent {
		err := validateContent(content)
		if err != nil {
			t.Errorf("Expected no error for valid content '%s', got: %v",
				content, err)
		}
	}

	// Invalid content
	invalidContent := []string{
		"",                    // Empty
		"   ",                 // Only whitespace
		"Unclosed {{variable", // Missing closing
		"Unopened variable}}", // Missing opening
	}

	for _, content := range invalidContent {
		err := validateContent(content)
		if err == nil {
			t.Errorf("Expected error for invalid content '%s'", content)
		}
	}
}

func TestValidateVariables_MultipleVariables(t *testing.T) {
	variables := map[string]models.Variable{
		"name": {
			Name:     "name",
			Type:     "text",
			Required: true,
		},
		"priority": {
			Name:    "priority",
			Type:    "choice",
			Options: []string{"low", "medium", "high"},
			Default: "medium",
		},
		"enabled": {
			Name: "enabled",
			Type: "boolean",
		},
	}

	err := validateVariables(variables)
	if err != nil {
		t.Errorf("Expected no error for valid variables, got: %v", err)
	}

	// Add invalid variable
	variables["invalid"] = models.Variable{
		Name: "invalid",
		Type: "unknown_type",
	}

	err = validateVariables(variables)
	if err == nil {
		t.Error("Expected error for invalid variable type")
	}
}
