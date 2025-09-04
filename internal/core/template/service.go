package template

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/diogopedro/shotgun/internal/models"
)

// TemplateService interface defines the template loading service contract
type TemplateService interface {
	LoadAllTemplates(ctx context.Context) ([]models.Template, error)
	GetTemplate(id string) (*models.Template, error)
	RefreshTemplates(ctx context.Context) error
	GetTemplateCount() int
}

// templateService implements TemplateService
type templateService struct {
	discovery *DiscoveryService
	cache     sync.Map // map[string]models.Template
	logger    *slog.Logger
	mu        sync.RWMutex
	loaded    bool
}

// NewTemplateService creates a new template service
func NewTemplateService(logger *slog.Logger) TemplateService {
	if logger == nil {
		logger = slog.Default()
	}

	return &templateService{
		discovery: NewDiscoveryService(logger),
		logger:    logger,
		loaded:    false,
	}
}

// LoadAllTemplates loads all templates from all sources and returns them sorted
func (s *templateService) LoadAllTemplates(ctx context.Context) ([]models.Template, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Discover templates
	templateInfos, err := s.discovery.DiscoverAllTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Clear existing cache
	s.cache.Range(func(key, value interface{}) bool {
		s.cache.Delete(key)
		return true
	})

	// Convert to templates and cache them
	var templates []models.Template
	var errors []error

	for _, info := range templateInfos {
		// Store in cache
		s.cache.Store(info.Template.ID, info.Template)
		templates = append(templates, info.Template)

		s.logger.Debug("Template loaded and cached",
			"id", info.Template.ID,
			"name", info.Template.Name,
			"source", info.Source,
			"path", info.FilePath)
	}

	// Sort templates by name for consistent ordering
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	s.loaded = true

	// Log aggregated errors if any
	if len(errors) > 0 {
		s.logger.Warn("Some templates failed to load",
			"error_count", len(errors),
			"successful_count", len(templates))
		for i, err := range errors {
			s.logger.Warn("Template loading error", "index", i, "error", err)
		}
	}

	s.logger.Info("Template loading completed",
		"total_templates", len(templates),
		"error_count", len(errors))

	return templates, nil
}

// GetTemplate retrieves a specific template by ID from cache
func (s *templateService) GetTemplate(id string) (*models.Template, error) {
	if id == "" {
		return nil, fmt.Errorf("template ID cannot be empty")
	}

	// Check if templates have been loaded
	s.mu.RLock()
	loaded := s.loaded
	s.mu.RUnlock()

	if !loaded {
		return nil, fmt.Errorf("templates not loaded yet, call LoadAllTemplates first")
	}

	// Try to get from cache
	if value, ok := s.cache.Load(id); ok {
		if template, ok := value.(models.Template); ok {
			return &template, nil
		}
		s.logger.Error("Invalid template type in cache", "id", id)
		return nil, fmt.Errorf("invalid template data in cache for ID: %s", id)
	}

	return nil, fmt.Errorf("template not found: %s", id)
}

// RefreshTemplates clears the cache and reloads all templates
func (s *templateService) RefreshTemplates(ctx context.Context) error {
	s.logger.Info("Refreshing template cache")

	// Clear the loaded flag
	s.mu.Lock()
	s.loaded = false
	s.mu.Unlock()

	// Reload templates
	_, err := s.LoadAllTemplates(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh templates: %w", err)
	}

	s.logger.Info("Template cache refreshed successfully")
	return nil
}

// GetTemplateCount returns the number of cached templates
func (s *templateService) GetTemplateCount() int {
	count := 0
	s.cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// GetAllCachedTemplates returns all cached templates (for testing/debugging)
func (s *templateService) GetAllCachedTemplates() []models.Template {
	var templates []models.Template

	s.cache.Range(func(key, value interface{}) bool {
		if template, ok := value.(models.Template); ok {
			templates = append(templates, template)
		}
		return true
	})

	// Sort by name for consistency
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates
}

// ValidateTemplateCache performs integrity checks on the cache
func (s *templateService) ValidateTemplateCache() error {
	var errors []error

	s.cache.Range(func(key, value interface{}) bool {
		// Validate key type
		id, ok := key.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("invalid cache key type: %T", key))
			return true
		}

		// Validate value type
		template, ok := value.(models.Template)
		if !ok {
			errors = append(errors, fmt.Errorf("invalid cache value type for key %s: %T", id, value))
			return true
		}

		// Validate ID consistency
		if template.ID != id {
			errors = append(errors, fmt.Errorf("cache key '%s' does not match template ID '%s'", id, template.ID))
		}

		return true
	})

	if len(errors) > 0 {
		return fmt.Errorf("cache validation failed with %d errors: %v", len(errors), errors)
	}

	return nil
}
