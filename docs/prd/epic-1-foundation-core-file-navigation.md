# Epic 1: Foundation & Core File Navigation

Establish the foundational Go project structure with Bubble Tea TUI framework, implement the concurrent file scanning system with .gitignore support, and deliver a fully functional file tree selection interface that allows users to choose which files to include in their prompts.

## Story 1.1: Project Setup & Go Module Initialization

As a developer,  
I want the project properly structured with Go modules and dependencies,  
so that I can build and run the application consistently across platforms.

### Acceptance Criteria
1. Go module initialized with proper module path (github.com/user/shotgun-cli)
2: Project follows standard Go layout with cmd/, internal/, and templates/ directories
3: All required dependencies added to go.mod (Bubble Tea, Bubbles, Lip Gloss, Cobra, TOML parser)
4: Makefile created with build, test, and clean targets
5: Basic main.go entry point that launches Bubble Tea application
6: README.md with project description and build instructions
7: Git repository initialized with appropriate .gitignore for Go projects

## Story 1.2: Core File Scanner Engine with Concurrency

As a user,  
I want the application to quickly scan my project directory,  
so that I can see all available files for selection.

### Acceptance Criteria
1: FileScanner component implemented with concurrent directory traversal using goroutines
2: Support for .gitignore pattern matching using doublestar library
3: Support for .shotgunignore with identical pattern syntax
4: Binary file detection and automatic exclusion using filetype library
5: Channel-based file streaming for memory efficiency
6: Proper error handling for permission denied and invalid paths
7: Unit tests achieving 90% coverage for scanner logic

## Story 1.3: File Tree UI Component with Selection

As a user,  
I want to navigate and select files through an interactive tree view,  
so that I can choose which files to include in my prompt.

### Acceptance Criteria
1: File tree component displays hierarchical directory structure
2: Checkboxes render for each file/folder (all checked by default)
3: Arrow keys navigate up/down through the tree
4: Space key toggles selection for current item
5: Right/Left arrows expand/collapse directories
6: Hierarchical selection works (deselecting folder deselects all contents)
7: Visual indicators for binary files (grayed out, unselectable)
8: File count status displays at bottom (X selected, Y excluded, Z ignored)

## Story 1.4: Keyboard Navigation & Screen Transitions

As a user,  
I want to navigate between screens using keyboard shortcuts,  
so that I can progress through the wizard efficiently.

### Acceptance Criteria
1: F1 key displays contextual help for current screen
2: F2 navigates to previous screen (disabled on first screen)
3: F3 advances to next screen with validation
4: ESC key shows exit confirmation dialog
5: Screen state persists when navigating back and forth
6: Progress indicator shows current screen (1/5, 2/5, etc.)
7: Proper focus management when entering each screen
