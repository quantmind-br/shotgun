package models

// Template represents a prompt template with metadata and variables
type Template struct {
	ID          string              `toml:"id" json:"id"`
	Name        string              `toml:"name" json:"name"`
	Version     string              `toml:"version" json:"version"`
	Description string              `toml:"description" json:"description"`
	Author      string              `toml:"author" json:"author"`
	Tags        []string            `toml:"tags" json:"tags"`
	Variables   map[string]Variable `toml:"variables" json:"variables"`
	Content     string              `toml:"content" json:"content"`
}

// Variable represents a template variable with validation constraints
type Variable struct {
	Name        string   `toml:"name" json:"name"`
	Type        string   `toml:"type" json:"type"` // text, multiline, auto, choice, boolean, number
	Required    bool     `toml:"required" json:"required"`
	Default     string   `toml:"default" json:"default"`
	Placeholder string   `toml:"placeholder" json:"placeholder"`
	MinLength   int      `toml:"min_length,omitempty" json:"min_length,omitempty"`
	MaxLength   int      `toml:"max_length,omitempty" json:"max_length,omitempty"`
	Options     []string `toml:"options,omitempty" json:"options,omitempty"`
}

// TemplateSource indicates where a template originated from
type TemplateSource int

const (
	TemplateSourceBuiltIn TemplateSource = iota
	TemplateSourceUser
)

// String returns the string representation of TemplateSource
func (ts TemplateSource) String() string {
	switch ts {
	case TemplateSourceBuiltIn:
		return "builtin"
	case TemplateSourceUser:
		return "user"
	default:
		return "unknown"
	}
}

// TemplateInfo holds a template with its source information
type TemplateInfo struct {
	Template Template
	Source   TemplateSource
	FilePath string
}

// ValidVariableTypes contains all supported variable types
var ValidVariableTypes = []string{
	"text",
	"multiline",
	"auto",
	"choice",
	"boolean",
	"number",
}
