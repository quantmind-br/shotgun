package app

import (
	"errors"
	"fmt"
)

// canAdvance checks if the current screen allows advancement to next screen
func (a *AppState) canAdvance() bool {
	switch a.CurrentScreen {
	case FileTreeScreen:
		// Requires at least one selected file
		selectedFiles := a.FileTree.GetSelectedFiles()
		return len(selectedFiles) > 0

	case TemplateScreen:
		// Requires a selected template
		return a.SelectedTemplate != nil

	case TaskScreen:
		// Requires valid task content using the TaskInputModel's validation
		return a.TaskInput.CanAdvance()

	case RulesScreen:
		// Rules are optional, always can advance
		return true

	case ConfirmScreen:
		// Final screen, can always "advance" (which means complete)
		return true

	default:
		return false
	}
}

// getValidationError returns the appropriate error message for validation failure
func (a *AppState) getValidationError() error {
	switch a.CurrentScreen {
	case FileTreeScreen:
		return errors.New("Please select at least one file before continuing")

	case TemplateScreen:
		return errors.New("Please select a template before continuing")

	case TaskScreen:
		return errors.New("Please enter a task description before continuing")

	case RulesScreen:
		// Should not happen since rules are optional
		return nil

	case ConfirmScreen:
		// Should not happen since this is the final screen
		return nil

	default:
		return errors.New("Unknown screen validation error")
	}
}

// validateScreenData performs comprehensive validation for a screen
func (a *AppState) validateScreenData(screen ScreenType) error {
	switch screen {
	case FileTreeScreen:
		selectedFiles := a.FileTree.GetSelectedFiles()
		if len(selectedFiles) == 0 {
			return errors.New("no files selected")
		}

		// Update shared state
		a.SelectedFiles = selectedFiles
		return nil

	case TemplateScreen:
		if a.SelectedTemplate == nil {
			return errors.New("no template selected")
		}
		return nil

	case TaskScreen:
		if !a.TaskInput.CanAdvance() {
			return errors.New("task description is required")
		}

		// Basic length validation - use the actual textarea content
		content := a.TaskInput.GetContent()
		if len(content) < 10 {
			return errors.New("task description should be at least 10 characters")
		}

		return nil

	case RulesScreen:
		// Rules are optional, but if provided, should be meaningful
		if a.RulesContent != "" && len(a.RulesContent) < 5 {
			return errors.New("rules should be at least 5 characters if provided")
		}
		return nil

	case ConfirmScreen:
		// Validate that all required data is present
		// Use current FileTree selections if available
		if len(a.SelectedFiles) == 0 && len(a.FileTree.GetSelectedFiles()) == 0 {
			return errors.New("no files selected")
		}
		if a.SelectedTemplate == nil {
			return errors.New("no template selected")
		}
		if a.TaskContent == "" {
			return errors.New("no task description provided")
		}
		return nil

	default:
		return fmt.Errorf("unknown screen: %v", screen)
	}
}

// getCurrentScreenProgress returns progress information for current screen
func (a *AppState) getCurrentScreenProgress() (current, total int) {
	total = 5                          // Total number of screens
	current = int(a.CurrentScreen) + 1 // Convert 0-based to 1-based
	return current, total
}

// getScreenTitle returns the display title for a screen
func (a *AppState) getScreenTitle(screen ScreenType) string {
	switch screen {
	case FileTreeScreen:
		return "Select Files"
	case TemplateScreen:
		return "Choose Template"
	case TaskScreen:
		return "Describe Task"
	case RulesScreen:
		return "Add Rules (Optional)"
	case ConfirmScreen:
		return "Review & Confirm"
	default:
		return "Unknown Screen"
	}
}

// isScreenComplete returns whether a screen has all required data
func (a *AppState) isScreenComplete(screen ScreenType) bool {
	return a.validateScreenData(screen) == nil
}
