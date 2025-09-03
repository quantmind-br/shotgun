package template

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed templates
var builtinTemplatesFS embed.FS

// listBuiltinTemplates returns a list of all embedded template file paths
func listBuiltinTemplates() ([]string, error) {
	var templatePaths []string

	err := fs.WalkDir(builtinTemplatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-TOML files
		if d.IsDir() || !isTomlFile(path) {
			return nil
		}

		templatePaths = append(templatePaths, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list built-in templates: %w", err)
	}

	return templatePaths, nil
}

// loadBuiltinTemplate loads a specific built-in template by path
func loadBuiltinTemplate(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("template path cannot be empty")
	}

	data, err := fs.ReadFile(builtinTemplatesFS, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read built-in template '%s': %w", path, err)
	}

	return data, nil
}

// isTomlFile checks if a file has .toml extension
func isTomlFile(filename string) bool {
	return len(filename) > 5 && filename[len(filename)-5:] == ".toml"
}

// validateBuiltinTemplatesFS validates that the embedded filesystem is accessible
func validateBuiltinTemplatesFS() error {
	// Try to open the root directory
	entries, err := fs.ReadDir(builtinTemplatesFS, ".")
	if err != nil {
		return fmt.Errorf("failed to access built-in templates filesystem: %w", err)
	}

	// Check if we have any entries
	if len(entries) == 0 {
		return fmt.Errorf("built-in templates filesystem is empty")
	}

	return nil
}
