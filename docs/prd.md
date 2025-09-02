# Shotgun CLI Product Requirements Document (PRD)

## Goals and Background Context

### Goals
- Deliver a functional TUI application that standardizes LLM prompt generation for development teams
- Achieve <2 minute prompt generation time from project context to formatted output
- Support cross-platform deployment (Windows, macOS, Linux) with single binary distribution
- Provide built-in templates covering 80% of common developer LLM use cases
- Enable team-wide adoption through zero-configuration startup and intuitive keyboard navigation
- Establish foundation for community-driven template ecosystem post-MVP

### Background Context

Shotgun CLI addresses the growing friction between developers' increasing reliance on LLMs and the inefficient, manual processes currently required to provide context to these tools. As development teams integrate AI assistants into their workflows for code review, debugging, and planning, they face repetitive tasks of copying files, formatting prompts, and maintaining consistency across team members. This tool transforms prompt generation from a manual, error-prone process into a streamlined, reproducible workflow that maintains developers in their terminal environment.

The solution leverages Go's cross-platform capabilities and the Bubble Tea TUI framework to deliver a lightweight, performant tool that integrates seamlessly into existing development workflows. By focusing specifically on the developer use case with repository-aware features and development-focused templates, Shotgun CLI aims to become the standard utility for LLM prompt generation in software development.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-09-02 | v1.0 | Initial PRD creation from Project Brief | John (PM) |

## Requirements

### Functional

- FR1: The application SHALL provide a 5-screen TUI wizard for prompt generation (File Tree, Template Selection, Task Input, Rules Input, Confirmation)
- FR2: The file tree view SHALL display all files/folders in the current directory with checkbox selection, respecting .gitignore and .shotgunignore patterns
- FR3: The file tree SHALL support hierarchical selection where selecting/deselecting a folder affects all contained items
- FR4: The application SHALL automatically detect and exclude binary files from selection with visual indicators
- FR5: The template selection screen SHALL display both built-in and user-defined templates from ~/.config/shotgun-cli/templates
- FR6: The application SHALL include 4 built-in templates: analyze_bug, make_diff, make_plan, and project_manager
- FR7: The task input screen SHALL provide a multiline text editor supporting UTF-8 characters and clipboard operations
- FR8: The rules input screen SHALL provide an optional multiline text field with skip capability (F4)
- FR9: The confirmation screen SHALL display file count, estimated output size, and preview before generation
- FR10: The application SHALL generate output in Markdown format with embedded file contents using <file> tags
- FR11: The application SHALL support complete keyboard navigation using F-keys (F1-F10) and standard navigation keys
- FR12: The application SHALL save generated prompts with timestamp-based filenames to prevent overwrites
- FR13: The application SHALL support a 'shotgun init' command to create a .shotgunignore file
- FR14: The application SHALL process templates using Go's text/template engine with variable substitution
- FR15: The application SHALL display real-time progress indicators during file scanning and prompt generation

### Non Functional

- NFR1: The application SHALL start in less than 2 seconds on standard hardware
- NFR2: File scanning SHALL complete within 5 seconds for repositories with 1000+ files
- NFR3: The UI SHALL maintain responsive frame rates (<16ms) during all operations
- NFR4: Memory usage SHALL not exceed 100MB for typical repositories (<5000 files)
- NFR5: The application SHALL support terminals with minimum 80x24 character display
- NFR6: The application SHALL work on Windows (PowerShell/CMD), macOS (Terminal), and Linux terminals
- NFR7: The application SHALL gracefully degrade on terminals without full color support
- NFR8: The application SHALL handle UTF-8 encoded files and support international characters
- NFR9: The binary size SHALL not exceed 20MB for any platform
- NFR10: The application SHALL achieve 90% unit test coverage for core business logic
- NFR11: The application SHALL handle concurrent file operations without UI blocking
- NFR12: The application SHALL provide clear error messages and recovery options for all failures

## Technical Assumptions

### Repository Structure: Monorepo
The application will be developed as a single Go repository containing all source code, templates, and documentation.

### Service Architecture
**Monolithic CLI Application** - Single binary executable with modular internal architecture using Go packages for separation of concerns. The application follows Elm Architecture patterns via Bubble Tea for state management.

### Testing Requirements
**Full Testing Pyramid** approach:
- Unit tests for all business logic components (90% coverage target)
- Integration tests for TUI flows using teatest framework
- E2E tests using expect scripts for complete user workflows
- Performance benchmarks for critical operations
- Cross-platform CI testing on Windows, Linux, and macOS

### Additional Technical Assumptions and Requests
- **Language**: Go 1.22+ for improved Windows support and performance optimizations
- **TUI Framework**: Bubble Tea v2.0.0-beta.4 with Bubbles component library v0.21.0
- **Styling**: Lip Gloss v1.0.0 for terminal styling with minimalist monochrome palette
- **Template Engine**: Go's built-in text/template for secure template processing
- **Configuration**: TOML format for templates and user configuration files
- **Binary Distribution**: Single static binary per platform via GitHub Releases
- **Dependency Management**: Go modules with minimal external dependencies
- **File Detection**: Binary file detection using h2non/filetype library
- **Glob Patterns**: Support via bmatcuk/doublestar for .gitignore/.shotgunignore
- **Concurrency**: Goroutines with channels for file scanning and content reading
- **State Management**: Immutable state following Elm Architecture patterns
- **Error Handling**: Comprehensive error handling with panic recovery
- **No Network Access**: All operations performed locally without external API calls
- **No Code Execution**: Templates cannot execute arbitrary code for security

## Epic List

**Epic 1: Foundation & Core File Navigation** - Establish project structure, implement concurrent file scanner with .gitignore support, and deliver working file tree selection UI

**Epic 2: Template System & User Input** - Create template engine, implement template selection UI, and build multiline text input components for task and rules

**Epic 3: Prompt Generation & Output** - Implement prompt assembly pipeline, create confirmation screen with size estimation, and generate final Markdown output

**Epic 4: Polish & Cross-Platform Optimization** - Add keyboard navigation system, implement progress indicators, ensure cross-platform compatibility, and create distribution pipeline

## Epic 1: Foundation & Core File Navigation

Establish the foundational Go project structure with Bubble Tea TUI framework, implement the concurrent file scanning system with .gitignore support, and deliver a fully functional file tree selection interface that allows users to choose which files to include in their prompts.

### Story 1.1: Project Setup & Go Module Initialization

As a developer,  
I want the project properly structured with Go modules and dependencies,  
so that I can build and run the application consistently across platforms.

#### Acceptance Criteria
1. Go module initialized with proper module path (github.com/user/shotgun-cli)
2: Project follows standard Go layout with cmd/, internal/, and templates/ directories
3: All required dependencies added to go.mod (Bubble Tea, Bubbles, Lip Gloss, Cobra, TOML parser)
4: Makefile created with build, test, and clean targets
5: Basic main.go entry point that launches Bubble Tea application
6: README.md with project description and build instructions
7: Git repository initialized with appropriate .gitignore for Go projects

### Story 1.2: Core File Scanner Engine with Concurrency

As a user,  
I want the application to quickly scan my project directory,  
so that I can see all available files for selection.

#### Acceptance Criteria
1: FileScanner component implemented with concurrent directory traversal using goroutines
2: Support for .gitignore pattern matching using doublestar library
3: Support for .shotgunignore with identical pattern syntax
4: Binary file detection and automatic exclusion using filetype library
5: Channel-based file streaming for memory efficiency
6: Proper error handling for permission denied and invalid paths
7: Unit tests achieving 90% coverage for scanner logic

### Story 1.3: File Tree UI Component with Selection

As a user,  
I want to navigate and select files through an interactive tree view,  
so that I can choose which files to include in my prompt.

#### Acceptance Criteria
1: File tree component displays hierarchical directory structure
2: Checkboxes render for each file/folder (all checked by default)
3: Arrow keys navigate up/down through the tree
4: Space key toggles selection for current item
5: Right/Left arrows expand/collapse directories
6: Hierarchical selection works (deselecting folder deselects all contents)
7: Visual indicators for binary files (grayed out, unselectable)
8: File count status displays at bottom (X selected, Y excluded, Z ignored)

### Story 1.4: Keyboard Navigation & Screen Transitions

As a user,  
I want to navigate between screens using keyboard shortcuts,  
so that I can progress through the wizard efficiently.

#### Acceptance Criteria
1: F1 key displays contextual help for current screen
2: F2 navigates to previous screen (disabled on first screen)
3: F3 advances to next screen with validation
4: ESC key shows exit confirmation dialog
5: Screen state persists when navigating back and forth
6: Progress indicator shows current screen (1/5, 2/5, etc.)
7: Proper focus management when entering each screen

## Epic 2: Template System & User Input

Create the template discovery and selection system, implement the template parsing engine, and build the multiline text input components that allow users to provide task descriptions and optional rules for their prompts.

### Story 2.1: Template Discovery & Loading System

As a user,  
I want the application to find and load available templates,  
so that I can choose from built-in or custom templates.

#### Acceptance Criteria
1: Built-in templates embedded in binary from templates/ directory
2: User templates discovered from ~/.config/shotgun-cli/templates (or Windows equivalent)
3: TOML parser correctly reads template metadata (name, version, description, variables)
4: Template validation ensures required fields are present
5: Templates with parsing errors are skipped with warning logs
6: Both built-in and user templates appear in unified list
7: Unit tests cover template discovery and parsing logic

### Story 2.2: Template Selection UI Screen

As a user,  
I want to browse and select a template from available options,  
so that I can use the appropriate format for my use case.

#### Acceptance Criteria
1: List view displays all discovered templates with name and description
2: Version number shows inline with template name
3: Arrow keys navigate up/down through template list
4: Enter or F3 selects template and advances to next screen
5: Template metadata (author, tags) displays in detail panel
6: Visual distinction for currently selected template
7: F2 returns to file tree screen preserving selection

### Story 2.3: Task Input Screen with Multiline Editor

As a user,  
I want to describe my task in detail with proper formatting,  
so that the LLM receives clear context about what I need.

#### Acceptance Criteria
1: Multiline text editor using Bubbles textarea component
2: Support for copy/paste operations (Ctrl+C, Ctrl+V)
3: UTF-8 character support for international text
4: Line and character count displays in real-time
5: Ctrl+Enter finalizes input and advances
6: F3 advances only if content is non-empty
7: Text persists when navigating back via F2
8: Word wrap functions correctly for long lines

### Story 2.4: Rules Input Screen (Optional)

As a user,  
I want to optionally specify additional rules or constraints,  
so that I can guide the LLM's response style or requirements.

#### Acceptance Criteria
1: Multiline text editor similar to task input screen
2: Clear indication that this field is optional
3: F4 key skips this screen entirely
4: F3 advances regardless of content (can be empty)
5: Content persists when navigating between screens
6: Placeholder text suggests example rules
7: Same UTF-8 and clipboard support as task input

## Epic 3: Prompt Generation & Output

Implement the prompt assembly system that combines templates with file contents, create the confirmation screen with size estimation, and generate the final Markdown output file with proper formatting.

### Story 3.1: Template Processing Engine

As a developer,  
I want templates to be processed with variable substitution,  
so that dynamic content is properly inserted into the output.

#### Acceptance Criteria
1: Go text/template engine processes template content
2: Variable substitution works for all defined variables (TASK, RULES, FILE_STRUCTURE)
3: Automatic variables populated (CURRENT_DATE, PROJECT_NAME)
4: Conditional sections ({{if}}) process correctly
5: Template functions available (upper, lower, trim)
6: Error handling for missing variables or template errors
7: Unit tests cover various template scenarios

### Story 3.2: File Structure Assembly

As a user,  
I want selected files to be included with their content in the prompt,  
so that the LLM has full context for analysis.

#### Acceptance Criteria
1: Generate tree-like structure showing directory hierarchy
2: Use ASCII characters for tree visualization (├── └── │)
3: Read file contents asynchronously using goroutines
4: Wrap file contents in <file path="...">content</file> tags
5: Skip binary files with appropriate message
6: Handle large files gracefully (streaming read)
7: Preserve exact file content including whitespace

### Story 3.3: Confirmation Screen with Size Estimation

As a user,  
I want to review prompt details before generation,  
so that I can ensure the output will be appropriate.

#### Acceptance Criteria
1: Display summary of selections (template name, file count, excluded items)
2: Calculate and display estimated output size in KB/MB
3: Show output filename with timestamp
4: Progress bar displays during size calculation
5: Warning appears for very large outputs (>500KB)
6: F10 confirms and triggers generation
7: F2 allows returning to make adjustments

### Story 3.4: Prompt Generation & File Writing

As a user,  
I want the final prompt saved to a file,  
so that I can use it with my preferred LLM tool.

#### Acceptance Criteria
1: Combine template + variables + file structure into final output
2: Save to current directory with timestamp filename (shotgun_prompt_YYYYMMDD_HHMM.md)
3: Ensure no filename collisions with incrementing counter if needed
4: Display success message with full file path
5: Handle write errors gracefully with clear error message
6: Non-blocking generation using goroutines
7: Progress indicator during file writing for large outputs

## Epic 4: Polish & Cross-Platform Optimization

Add the complete keyboard navigation system with F-key shortcuts, implement visual progress indicators, ensure cross-platform compatibility, and create the distribution pipeline for releasing the application.

### Story 4.1: Global Keyboard Navigation System

As a user,  
I want consistent keyboard shortcuts across all screens,  
so that I can navigate efficiently without learning different commands.

#### Acceptance Criteria
1: F1-F10 keys properly mapped and handled globally
2: Help overlay (F1) shows context-sensitive shortcuts
3: Navigation between all 5 screens works smoothly
4: ESC key handling with confirmation dialog
5: Keyboard shortcuts don't conflict during text editing
6: Tab order logical for accessibility
7: All shortcuts documented in help screen

### Story 4.2: Progress Indicators & Loading States

As a user,  
I want visual feedback during long operations,  
so that I know the application is working.

#### Acceptance Criteria
1: Spinner displays during initial file scanning
2: Progress bar for file reading operations
3: Loading state for template discovery
4: Size calculation progress in confirmation screen
5: All progress indicators use consistent styling
6: Operations remain cancellable with ESC
7: Smooth animations without flickering

### Story 4.3: Cross-Platform Testing & Compatibility

As a developer,  
I want the application to work consistently across platforms,  
so that all users have the same experience.

#### Acceptance Criteria
1: Application runs on Windows PowerShell and CMD
2: Application runs on macOS Terminal and iTerm2
3: Application runs on common Linux terminals
4: Unicode characters display correctly on all platforms
5: Colors degrade gracefully on limited terminals
6: Keyboard shortcuts work across all environments
7: CI pipeline tests on Windows, macOS, and Linux

### Story 4.4: Init Command & Shotgunignore Support

As a user,  
I want to create a .shotgunignore file easily,  
so that I can customize which files are excluded.

#### Acceptance Criteria
1: 'shotgun init' command creates .shotgunignore file
2: Template .shotgunignore includes common patterns
3: File created only if it doesn't exist
4: Success message confirms file creation
5: Command available via Cobra CLI framework
6: Help text explains usage and purpose
7: Integration with main file scanner confirmed

### Story 4.5: Binary Distribution Pipeline

As a maintainer,  
I want automated builds for all platforms,  
so that users can easily download and use the application.

#### Acceptance Criteria
1: GitHub Actions workflow builds for Linux, macOS, Windows
2: Binaries created for amd64 and arm64 architectures where applicable
3: Version information embedded in binary
4: Artifacts uploaded to GitHub Releases
5: Checksums generated for verification
6: README includes installation instructions
7: Binary size optimized with appropriate build flags

## Checklist Results Report

_This section will be populated after PRD review and checklist execution._

## Next Steps

### UX Expert Prompt
Please review this PRD and create detailed UI/UX specifications for the Shotgun CLI TUI application, focusing on the 5-screen wizard flow, keyboard navigation patterns, and visual design within terminal constraints.

### Architect Prompt
Please review this PRD and create a comprehensive technical architecture document for Shotgun CLI, detailing the Go package structure, Bubble Tea component architecture, concurrency patterns, and implementation approach for each epic and story.