package template

import (
	"fmt"
	"strings"

	"github.com/diogopedro/shotgun/internal/models"
)

// validateTemplate validates a template struct for required fields and consistency
func validateTemplate(template *models.Template) error {
	if template == nil {
		return fmt.Errorf("template is nil")
	}

	// Validate required fields
	if err := validateRequiredFields(template); err != nil {
		return err
	}

	// Validate variables
	if err := validateVariables(template.Variables); err != nil {
		return err
	}

	// Validate content
	if err := validateContent(template.Content); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields ensures all required fields are present and valid
func validateRequiredFields(template *models.Template) error {
	if strings.TrimSpace(template.Name) == "" {
		return fmt.Errorf("template name is required")
	}

	if strings.TrimSpace(template.Version) == "" {
		return fmt.Errorf("template version is required")
	}

	if strings.TrimSpace(template.Description) == "" {
		return fmt.Errorf("template description is required")
	}

	if strings.TrimSpace(template.Content) == "" {
		return fmt.Errorf("template content is required")
	}

	if strings.TrimSpace(template.ID) == "" {
		return fmt.Errorf("template ID is required")
	}

	// Validate version format (basic semver check)
	if !isValidVersion(template.Version) {
		return fmt.Errorf("invalid version format: %s (expected semver like 1.0.0)", template.Version)
	}

	return nil
}

// validateVariables validates all variables in the template
func validateVariables(variables map[string]models.Variable) error {
	for name, variable := range variables {
		if err := validateVariable(name, variable); err != nil {
			return fmt.Errorf("variable '%s': %w", name, err)
		}
	}
	return nil
}

// validateVariable validates a single variable
func validateVariable(name string, variable models.Variable) error {
	// Validate variable name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if strings.TrimSpace(variable.Name) == "" {
		return fmt.Errorf("variable Name field is required")
	}

	// Variable name in map should match Variable.Name field
	if name != variable.Name {
		return fmt.Errorf("variable map key '%s' does not match Variable.Name '%s'", name, variable.Name)
	}

	// Validate variable type
	if !isValidVariableType(variable.Type) {
		return fmt.Errorf("invalid variable type: %s (valid types: %v)",
			variable.Type, models.ValidVariableTypes)
	}

	// Type-specific validation
	switch variable.Type {
	case "choice":
		if len(variable.Options) == 0 {
			return fmt.Errorf("choice variable must have options")
		}
		// Validate default value is in options if provided
		if variable.Default != "" {
			found := false
			for _, option := range variable.Options {
				if option == variable.Default {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("default value '%s' is not in options %v",
					variable.Default, variable.Options)
			}
		}
	case "boolean":
		// Validate boolean default
		if variable.Default != "" {
			if variable.Default != "true" && variable.Default != "false" {
				return fmt.Errorf("boolean variable default must be 'true' or 'false', got: %s",
					variable.Default)
			}
		}
	case "number":
		// Could add number validation here if needed
	}

	// Validate length constraints
	if variable.MinLength < 0 {
		return fmt.Errorf("MinLength cannot be negative: %d", variable.MinLength)
	}

	if variable.MaxLength > 0 && variable.MaxLength < variable.MinLength {
		return fmt.Errorf("MaxLength (%d) cannot be less than MinLength (%d)",
			variable.MaxLength, variable.MinLength)
	}

	// Validate default value meets length constraints
	if variable.Default != "" {
		if len(variable.Default) < variable.MinLength {
			return fmt.Errorf("default value length (%d) is less than MinLength (%d)",
				len(variable.Default), variable.MinLength)
		}
		if variable.MaxLength > 0 && len(variable.Default) > variable.MaxLength {
			return fmt.Errorf("default value length (%d) exceeds MaxLength (%d)",
				len(variable.Default), variable.MaxLength)
		}
	}

	return nil
}

// validateContent performs basic content validation
func validateContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Check for basic template syntax (variable placeholders)
	// This is a simple check for {{ }} patterns
	if strings.Contains(content, "{{") && !strings.Contains(content, "}}") {
		return fmt.Errorf("template content has malformed variable placeholders (unclosed {{)")
	}

	if strings.Contains(content, "}}") && !strings.Contains(content, "{{") {
		return fmt.Errorf("template content has malformed variable placeholders (unopened }})")
	}

	return nil
}

// isValidVariableType checks if the variable type is valid
func isValidVariableType(varType string) bool {
	for _, validType := range models.ValidVariableTypes {
		if varType == validType {
			return true
		}
	}
	return false
}

// isValidVersion performs basic semantic version validation
func isValidVersion(version string) bool {
	if version == "" {
		return false
	}

	// Basic semver pattern: X.Y.Z (can be extended)
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}

	// Each part should be numeric (basic check)
	for _, part := range parts {
		if part == "" {
			return false
		}
		// Check if all characters are digits
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
		}
	}

	return true
}
