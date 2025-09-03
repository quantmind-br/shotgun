package template

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/user/shotgun-cli/internal/models"
)

// parseTemplateFromData parses TOML data into a Template struct
func parseTemplateFromData(data []byte) (*models.Template, error) {
	if len(data) == 0 {
		return nil, NewParsingError("", fmt.Errorf("template data is empty"))
	}

	// Check file size limit (1MB)
	const maxFileSize = 1024 * 1024 // 1MB
	if len(data) > maxFileSize {
		return nil, NewContentSizeError("", len(data))
	}

	var rawTemplate struct {
		ID          string                         `toml:"id"`
		Name        string                         `toml:"name"`
		Version     string                         `toml:"version"`
		Description string                         `toml:"description"`
		Author      string                         `toml:"author"`
		Tags        []string                       `toml:"tags"`
		Variables   map[string]tomlVariable        `toml:"variables"`
		Content     string                         `toml:"content"`
	}

	// Parse TOML data
	if err := toml.Unmarshal(data, &rawTemplate); err != nil {
		return nil, NewParsingError("", err)
	}

	// Create Template struct
	template := &models.Template{
		ID:          rawTemplate.ID,
		Name:        rawTemplate.Name,
		Version:     rawTemplate.Version,
		Description: rawTemplate.Description,
		Author:      rawTemplate.Author,
		Tags:        rawTemplate.Tags,
		Variables:   make(map[string]models.Variable),
		Content:     rawTemplate.Content,
	}

	// Convert variables
	for name, rawVar := range rawTemplate.Variables {
		variable := models.Variable{
			Name:        rawVar.Name,
			Type:        rawVar.Type,
			Required:    rawVar.Required,
			Default:     rawVar.Default,
			Placeholder: rawVar.Placeholder,
			MinLength:   rawVar.MinLength,
			MaxLength:   rawVar.MaxLength,
			Options:     rawVar.Options,
		}
		template.Variables[name] = variable
	}

	// Generate ID from name if not provided
	if template.ID == "" {
		template.ID = generateTemplateID(template.Name)
	}

	// Validate the parsed template
	if err := validateTemplate(template); err != nil {
		return nil, NewValidationError("", err.Error())
	}

	return template, nil
}

// tomlVariable represents a variable in TOML format
type tomlVariable struct {
	Name        string   `toml:"name"`
	Type        string   `toml:"type"`
	Required    bool     `toml:"required"`
	Default     string   `toml:"default"`
	Placeholder string   `toml:"placeholder"`
	MinLength   int      `toml:"min_length"`
	MaxLength   int      `toml:"max_length"`
	Options     []string `toml:"options"`
}

// generateTemplateID creates an ID from the template name
func generateTemplateID(name string) string {
	if name == "" {
		return "unnamed-template"
	}
	
	// Simple ID generation: lowercase, replace spaces with dashes
	id := ""
	lastWasDash := true // Start as true to avoid leading dashes
	
	for _, r := range name {
		switch {
		case r >= 'A' && r <= 'Z':
			id += string(r - 'A' + 'a')
			lastWasDash = false
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			id += string(r)
			lastWasDash = false
		case r == ' ', r == '_', r == '\t':
			if !lastWasDash && len(id) > 0 {
				id += "-"
				lastWasDash = true
			}
		}
	}
	
	// Remove trailing dash if any
	if len(id) > 0 && id[len(id)-1] == '-' {
		id = id[:len(id)-1]
	}
	
	if id == "" {
		return "unnamed-template"
	}
	
	return id
}