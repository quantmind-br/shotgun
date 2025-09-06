package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/diogopedro/shotgun/internal/cli/templates"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	var force bool

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create a .shotgunignore file",
		Long: `Create a .shotgunignore file in the current directory.

The .shotgunignore file contains patterns for files and directories 
that should be excluded when scanning project files. This is useful 
for excluding build artifacts, dependencies, and other non-source files.

Examples:
  shotgun init              # Create .shotgunignore if it doesn't exist
  shotgun init --force      # Create .shotgunignore, overwriting if exists`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return CreateShotgunignore(force)
		},
	}

	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing .shotgunignore file")

	return initCmd
}

// CreateShotgunignore creates a .shotgunignore file in the current directory
func CreateShotgunignore(force bool) error {
	if err := ValidateDirectory(); err != nil {
		return err
	}

	filename := ".shotgunignore"
	filepath := filepath.Join(".", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); err == nil && !force {
		return fmt.Errorf(".shotgunignore file already exists (use --force to overwrite)")
	}

	// Generate template content
	content := GenerateTemplate()

	// Write file
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create .shotgunignore file: %w", err)
	}

	if force {
		fmt.Println("✓ .shotgunignore file created (overwritten)")
	} else {
		fmt.Println("✓ .shotgunignore file created successfully")
	}

	return nil
}

// ValidateDirectory ensures we can write to the current directory
func ValidateDirectory() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check if we can write to the directory
	tempFile := filepath.Join(wd, ".shotgun-temp-test")
	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("cannot write to current directory: %w", err)
	}
	os.Remove(tempFile) // Clean up

	return nil
}

// GenerateTemplate generates the default .shotgunignore content
func GenerateTemplate() string {
	return templates.ShotgunignoreTemplate
}
