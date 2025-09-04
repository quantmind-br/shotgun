package template

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/diogopedro/shotgun/internal/core/config"
	"github.com/diogopedro/shotgun/internal/models"
)

// DiscoveryService handles template discovery from various sources
type DiscoveryService struct {
	logger *slog.Logger
}

// NewDiscoveryService creates a new template discovery service
func NewDiscoveryService(logger *slog.Logger) *DiscoveryService {
	if logger == nil {
		logger = slog.Default()
	}
	return &DiscoveryService{
		logger: logger,
	}
}

// DiscoverAllTemplates finds all templates from built-in and user sources
func (d *DiscoveryService) DiscoverAllTemplates(ctx context.Context) ([]models.TemplateInfo, error) {
	var allTemplates []models.TemplateInfo

	// Discover built-in templates
	builtinTemplates, err := d.DiscoverBuiltinTemplates(ctx)
	if err != nil {
		d.logger.Warn("Failed to discover built-in templates", "error", err)
	} else {
		allTemplates = append(allTemplates, builtinTemplates...)
	}

	// Discover user templates
	userTemplates, err := d.DiscoverUserTemplates(ctx)
	if err != nil {
		d.logger.Warn("Failed to discover user templates", "error", err)
	} else {
		allTemplates = append(allTemplates, userTemplates...)
	}

	// Deduplicate templates (user templates override built-in ones by ID)
	deduplicated := d.deduplicateTemplates(allTemplates)

	d.logger.Debug("Template discovery completed",
		"total_found", len(allTemplates),
		"after_deduplication", len(deduplicated),
		"builtin_count", len(builtinTemplates),
		"user_count", len(userTemplates))

	return deduplicated, nil
}

// DiscoverBuiltinTemplates finds all embedded templates
func (d *DiscoveryService) DiscoverBuiltinTemplates(ctx context.Context) ([]models.TemplateInfo, error) {
	var templates []models.TemplateInfo

	// Walk through embedded filesystem
	err := fs.WalkDir(builtinTemplatesFS, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non-TOML files
		if !strings.HasSuffix(strings.ToLower(dirEntry.Name()), ".toml") {
			return nil
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read template file
		data, readErr := fs.ReadFile(builtinTemplatesFS, path)
		if readErr != nil {
			d.logger.Warn("Failed to read built-in template file",
				"path", path,
				"error", readErr)
			return nil // Continue processing other files
		}

		// Parse template (will be implemented in engine.go)
		template, parseErr := parseTemplateFromData(data)
		if parseErr != nil {
			d.logger.Warn("Failed to parse built-in template",
				"path", path,
				"error", parseErr)
			return nil // Continue processing other files
		}

		// Add to results
		templates = append(templates, models.TemplateInfo{
			Template: *template,
			Source:   models.TemplateSourceBuiltIn,
			FilePath: path,
		})

		d.logger.Debug("Discovered built-in template",
			"id", template.ID,
			"name", template.Name,
			"path", path)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk built-in templates: %w", err)
	}

	return templates, nil
}

// DiscoverUserTemplates finds all user templates from config directory
func (d *DiscoveryService) DiscoverUserTemplates(ctx context.Context) ([]models.TemplateInfo, error) {
	var templates []models.TemplateInfo

	// Get user templates directory
	userDir, err := config.GetUserTemplatesDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user templates directory: %w", err)
	}

	// Check if directory exists
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		d.logger.Debug("User templates directory does not exist", "path", userDir)
		return templates, nil // Return empty slice, not an error
	}

	// Walk through user templates directory
	err = filepath.WalkDir(userDir, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			// Log but don't fail completely for permission errors
			d.logger.Warn("Error accessing user template path",
				"path", path,
				"error", err)
			return nil
		}

		// Skip non-TOML files
		if !strings.HasSuffix(strings.ToLower(dirEntry.Name()), ".toml") {
			return nil
		}

		// Validate path safety (prevent directory traversal)
		if !d.isPathSafe(userDir, path) {
			d.logger.Warn("Unsafe template path detected, skipping",
				"path", path,
				"base", userDir)
			return nil
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read template file
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			d.logger.Warn("Failed to read user template file",
				"path", path,
				"error", readErr)
			return nil // Continue processing other files
		}

		// Parse template
		template, parseErr := parseTemplateFromData(data)
		if parseErr != nil {
			d.logger.Warn("Failed to parse user template",
				"path", path,
				"error", parseErr)
			return nil // Continue processing other files
		}

		// Add to results
		templates = append(templates, models.TemplateInfo{
			Template: *template,
			Source:   models.TemplateSourceUser,
			FilePath: path,
		})

		d.logger.Debug("Discovered user template",
			"id", template.ID,
			"name", template.Name,
			"path", path)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk user templates directory: %w", err)
	}

	return templates, nil
}

// deduplicateTemplates removes duplicate templates, preferring user templates over built-in ones
func (d *DiscoveryService) deduplicateTemplates(templates []models.TemplateInfo) []models.TemplateInfo {
	templateMap := make(map[string]models.TemplateInfo)

	for _, template := range templates {
		existing, exists := templateMap[template.Template.ID]

		if !exists {
			// New template, add it
			templateMap[template.Template.ID] = template
		} else if template.Source == models.TemplateSourceUser && existing.Source == models.TemplateSourceBuiltIn {
			// User template overrides built-in template
			d.logger.Debug("User template overriding built-in template",
				"id", template.Template.ID,
				"user_path", template.FilePath,
				"builtin_path", existing.FilePath)
			templateMap[template.Template.ID] = template
		} else if template.Source == existing.Source {
			// Same source type - log warning about duplicate
			d.logger.Warn("Duplicate template found, using first occurrence",
				"id", template.Template.ID,
				"first_path", existing.FilePath,
				"duplicate_path", template.FilePath)
		}
		// Built-in template trying to override user template - ignore
	}

	// Convert map back to slice
	result := make([]models.TemplateInfo, 0, len(templateMap))
	for _, template := range templateMap {
		result = append(result, template)
	}

	return result
}

// isPathSafe checks if a path is within the allowed base directory
func (d *DiscoveryService) isPathSafe(baseDir, targetPath string) bool {
	// Clean and resolve paths
	cleanBase, err := filepath.Abs(filepath.Clean(baseDir))
	if err != nil {
		d.logger.Error("Failed to resolve base directory", "path", baseDir, "error", err)
		return false
	}

	cleanTarget, err := filepath.Abs(filepath.Clean(targetPath))
	if err != nil {
		d.logger.Error("Failed to resolve target path", "path", targetPath, "error", err)
		return false
	}

	// Check if target path starts with base directory
	return strings.HasPrefix(cleanTarget, cleanBase+string(filepath.Separator)) ||
		cleanTarget == cleanBase
}

// parseTemplateFromData is implemented in engine.go
