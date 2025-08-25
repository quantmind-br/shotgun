# shotgun-cli-v3 Product Requirements Document (PRD)

## Goals and Background Context

### Goals
- Deliver a high-performance TUI application for generating standardized LLM prompts from templates
- Enable rapid, consistent prompt creation with project file context integration
- Provide intuitive keyboard-only navigation across all workflow steps
- Support both built-in and user-customized template workflows
- Achieve cross-platform compatibility (Windows, Linux, macOS) with excellent terminal support

### Background Context

The shotgun-cli addresses the critical need for standardized prompt generation in LLM-driven development workflows. Currently, developers manually copy-paste file contents and reconstruct context repeatedly, leading to formatting inconsistencies and time waste.

This TUI application, built with Go 1.22+ and Bubble Tea v2.0.0-beta.4, provides a wizard-driven interface that automatically maps project file structures, applies .gitignore/.shotgunignore rules, and generates comprehensive prompts through template-driven workflows. The solution targets development teams requiring reproducible, context-rich prompts for code analysis, debugging, planning, and project management tasks.

### Change Log
| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2024-03-15 | 1.0 | Initial PRD creation based on PLAN.md | John (PM) |

## Requirements

### Functional Requirements

**FR1**: The application must provide a file tree interface with hierarchical checkbox selection, where all items start selected by default and users deselect unwanted files/folders.

**FR2**: The system must automatically exclude files matching .gitignore and .shotgunignore patterns from the file tree display.

**FR3**: The application must support template selection from both embedded templates (analyze_bug, make_diff, make_plan, project_manager) and user-created custom templates.

**FR4**: The system must provide multiline text input editors with UTF-8 support for Task and Rules sections, including paste functionality and character counting.

**FR5**: The application must generate a complete file structure section showing both directory tree and file contents in tagged format.

**FR6**: The system must provide real-time size estimation of the final prompt with visual progress indicators before generation.

**FR7**: The application must support 100% keyboard navigation across all screens with F-key shortcuts and contextual controls.

**FR8**: The system must save generated prompts as timestamped Markdown files (shotgun_prompt_YYYYMMDD_HHMM.md) in the current directory.

### Non-Functional Requirements

**NFR1**: The application must start in under 2 seconds on standard hardware.

**NFR2**: File scanning of 1000+ files must complete in under 5 seconds with concurrent processing.

**NFR3**: The UI must maintain responsive frame rates (< 16ms frame time) during all operations.

**NFR4**: Memory usage must stay under 100MB for typical repositories.

**NFR5**: The application must provide graceful cross-platform compatibility across Windows PowerShell, Linux terminals, and macOS Terminal/iTerm2.

**NFR6**: The system must handle UTF-8 text input correctly including special characters (ç, á, ô, ñ, etc.).

**NFR7**: The application must recover gracefully from errors without crashes, providing meaningful error messages to users.

## User Interface Design Goals

### Overall UX Vision
Minimalist, monochrome terminal interface prioritizing speed and keyboard efficiency. The design follows a clean, distraction-free aesthetic with subtle accent colors (soft mint green #6ee7b7) for focus states and warm amber (#fbbf24) for highlights. Every interaction is optimized for professional developers who value precision and workflow efficiency over visual complexity.

### Key Interaction Paradigms
- **Keyboard-First Navigation**: 100% keyboard operation with global F-key shortcuts (F1-help, F2-back, F3-forward, ESC-exit)
- **Contextual Controls**: Screen-specific shortcuts that change based on current mode (navigation vs. editing)
- **Progressive Disclosure**: 5-screen wizard that maintains state while allowing free navigation between steps
- **Immediate Feedback**: Real-time validation, character counts, and visual indicators for all user actions

### Core Screens and Views
- **File Tree Selection** (1/5): Hierarchical file browser with checkbox states, initially all-selected with deselection workflow
- **Template Selection** (2/5): Vertical list selector showing both built-in and custom templates with metadata
- **Task Input Editor** (3/5): Multiline text editor with UTF-8 support, word wrap, and character counting
- **Rules Input Editor** (4/5): Optional multiline editor with skip functionality and visual "optional" indicators  
- **Confirmation Summary** (5/5): Complete review with size estimation, progress bars, and final generation step

### Accessibility: WCAG AA
Full keyboard navigation compliance, high contrast monochrome color scheme, clear focus indicators, and graceful degradation for limited terminals. Screen reader compatibility through semantic text structure and meaningful state descriptions.

### Branding
Early 1900s monochrome aesthetic with clean typography and minimal visual elements. Uses ASCII art borders, simple progress indicators, and restrained color palette to maintain professional, timeless appearance that works across all terminal environments.

### Target Device and Platforms: Cross-Platform
Optimized for professional terminal environments including Windows PowerShell, ConPTY, Linux bash/zsh, macOS Terminal/iTerm2, and modern terminals like WezTerm. Responsive design supporting minimum 80x24 terminal size with graceful scaling for larger displays.

## Technical Assumptions

### Repository Structure: Monorepo
Single repository containing the complete Go application with embedded templates, clear directory structure separating internal packages, and comprehensive testing suite. This supports the binary-only distribution model while maintaining development simplicity.

### Service Architecture
**Monolithic TUI Application** - Single-binary Go application using concurrent goroutines for file scanning, template processing, and UI operations. The architecture leverages Go's excellent concurrency primitives (channels, sync packages) within the Bubble Tea reactive framework for responsive user experience without the complexity of distributed services.

### Testing Requirements  
**Full Testing Pyramid** - Comprehensive testing strategy including:
- Unit tests with 90%+ coverage for core business logic
- Integration tests for TUI workflows using teatest helpers
- End-to-end tests with golden file validation
- Cross-platform compatibility tests in CI/CD
- Performance benchmarks for file scanning and UI responsiveness
- Fuzz testing for template parsing and file input handling

### Additional Technical Assumptions and Requests

**Language & Runtime**: Go 1.22+ for generics optimization, improved Windows compatibility, and enhanced performance characteristics. Native compilation eliminates runtime dependencies.

**TUI Framework**: Bubble Tea v2.0.0-beta.4 with Elm Architecture for predictable state management, enhanced keyboard handling, and improved Windows terminal support. Bubbles v0.21.0 for mature UI components (filepicker, textarea, list) with horizontal scrolling.

**Styling & Theming**: Lip Gloss v1.0.0 for sophisticated terminal styling with gradient support, flexible layouts, and consistent cross-platform rendering.

**Template Engine**: Go's native text/template package for security-safe template processing with custom function maps and validation.

**File Processing**: Concurrent file scanning using worker pools, channels, and goroutines. Integration with doublestar for glob pattern matching (.gitignore/.shotgunignore support).

**Configuration Management**: Viper for structured configuration with environment variable support and cross-platform config directories.

**Build & Distribution**: Single-binary compilation with embedded assets, cross-platform GitHub Actions CI/CD, and multi-architecture release artifacts.

**Dependencies**: Minimal external dependencies focusing on the Charm ecosystem (Bubble Tea, Bubbles, Lip Gloss) plus essential utilities for TOML parsing, file type detection, and text processing.

## Epic List

**Epic 1: Foundation & Core Infrastructure**
Establish project setup, Go module structure, core file scanning engine with concurrency, and basic TUI framework with a simple health-check interface to validate the technical stack.

**Epic 2: File Management & Navigation System**
Implement the complete file tree interface with hierarchical selection, .gitignore/.shotgunignore processing, binary detection, and keyboard navigation providing full file context gathering functionality.

**Epic 3: Template System & User Input**
Create template discovery engine for built-in and custom templates, implement multiline editors for Task/Rules input with UTF-8 support, and build the complete 5-screen wizard workflow.

**Epic 4: Prompt Generation & Output**
Develop the template processing engine, file content aggregation system, real-time size estimation, and final Markdown output generation with comprehensive error handling.

**Epic 5: Polish & Cross-Platform Optimization**
Refine UI styling with the monochrome theme, implement session management and history, optimize performance for large repositories, and ensure robust cross-platform terminal compatibility.

## Epic 1 - Foundation & Core Infrastructure

**Epic Goal:** Establish the technical foundation with Go project structure, Bubble Tea TUI framework, concurrent file scanning engine, and a basic interface that validates the complete technical stack while delivering initial functionality.

### Story 1.1: Project Setup and Go Module Initialization
As a developer,
I want a properly structured Go project with all dependencies configured,
so that I can begin development with a solid, maintainable foundation.

#### Acceptance Criteria
1. Go module initialized with go.mod specifying Go 1.22+ requirement
2. Complete directory structure created matching PLAN.md specifications (cmd/, internal/, templates/)
3. All core dependencies added to go.mod (Bubble Tea v2.0.0-beta.4, Bubbles v0.21.0, Lip Gloss v1.0.0, etc.)
4. Basic Makefile with build, test, clean targets
5. .gitignore configured for Go projects with common exclusions
6. Initial README.md with project overview and build instructions

### Story 1.2: Core File Scanner Engine with Concurrency
As a user,
I want the system to efficiently scan and catalog project files in the current directory,
so that I can work with projects of any size without performance issues.

#### Acceptance Criteria
1. FileScanner struct implemented with concurrent goroutine-based directory traversal
2. Worker pool pattern using runtime.NumCPU() for optimal performance
3. Channel-based file information streaming to prevent memory bottlenecks
4. Basic binary file detection using h2non/filetype library
5. File metadata collection (size, modification time, type) for all discovered files
6. Performance requirement: scan 1000+ files in under 5 seconds
7. Graceful error handling for permission denied and other file system errors

### Story 1.3: Basic Ignore Rules Processing
As a user,
I want the file scanner to respect .gitignore patterns,
so that irrelevant files are automatically excluded from consideration.

#### Acceptance Criteria
1. .gitignore file parsing using doublestar/v4 glob pattern matching
2. Ignore rule application during file scanning phase
3. Common patterns automatically excluded (node_modules/, .git/, dist/, build/)
4. Proper handling of nested .gitignore files and directory-specific rules
5. File exclusion logged at debug level for troubleshooting
6. Performance: ignore processing adds <10% to scan time overhead

### Story 1.4: Bubble Tea Application Foundation
As a developer,
I want the TUI application framework properly initialized,
so that all future UI screens can be built on a solid reactive architecture.

#### Acceptance Criteria
1. Main Bubble Tea application struct created following Elm Architecture
2. Initial model with state management for application lifecycle
3. Root command handler with proper initialization and cleanup
4. Global key bindings implemented (ESC to exit, F1 for help)
5. Terminal size detection and responsive layout foundation
6. Error recovery and graceful shutdown on panic or interrupt
7. TEA_DEBUG support for development debugging

### Story 1.5: Health Check Interface
As a user,
I want to see a simple interface that validates the technical stack,
so that I know the application is working correctly and can provide feedback.

#### Acceptance Criteria
1. Basic TUI screen displaying "shotgun-cli" title and version
2. File scan progress indicator showing discovered files count
3. Simple status display showing scan completion and total files found
4. "Press ESC to exit" instruction for user interaction
5. Real-time updates during file scanning operation
6. Basic styling using Lip Gloss with monochrome theme foundation
7. Cross-platform terminal compatibility validated (Windows/Linux/macOS)

## Epic 2 - File Management & Navigation System

**Epic Goal:** Implement the complete file tree interface with hierarchical checkbox selection, advanced ignore rule processing (.gitignore and .shotgunignore), binary file handling, and full keyboard navigation to provide comprehensive file context gathering functionality.

### Story 2.1: Hierarchical File Tree Widget
As a user,
I want to see my project files displayed in a tree structure with folder expansion/collapse,
so that I can understand the project organization and navigate efficiently.

#### Acceptance Criteria
1. Tree widget component built on Bubbles filepicker with hierarchical display
2. ASCII tree characters (├── └── │) for clean visual hierarchy
3. Folder expansion/collapse with → and ← arrow keys
4. Directory icons and file type indicators for visual distinction
5. Viewport scrolling for projects with many files using Bubbles viewport
6. Performance: smooth rendering for trees with 1000+ items
7. Keyboard navigation with ↑↓ keys maintaining current selection state

### Story 2.2: Checkbox Selection System with All-Selected Default
As a user,
I want all files to start selected with checkboxes so I can deselect unwanted items,
so that I can quickly exclude irrelevant files while including everything relevant by default.

#### Acceptance Criteria
1. All discovered files and folders display with checked checkboxes initially
2. Space bar toggles individual item selection state
3. Hierarchical selection: unchecking folder unchecks all contents automatically
4. Visual distinction between checked [✓], unchecked [ ], and partially selected [◐] states
5. Selection counter showing "X selected · Y excluded · Z ignored" in status bar
6. Ctrl+A select all, Ctrl+I invert selection keyboard shortcuts
7. Selection state persists during tree navigation and expansion/collapse operations

### Story 2.3: Enhanced Ignore Rules with .shotgunignore Support
As a user,
I want to use project-specific .shotgunignore rules in addition to .gitignore,
so that I can customize file exclusions for prompt generation without affecting Git.

#### Acceptance Criteria
1. .shotgunignore file detection and parsing in project root
2. Combined rule processing: .gitignore + .shotgunignore patterns
3. Rule precedence: .shotgunignore overrides .gitignore when conflicts exist
4. Pattern support: glob patterns, directory exclusions, negation with !
5. Visual indicators for ignored files (grayed out, "ignored" label)
6. Performance: rule processing adds <15% overhead to scanning
7. `shotgun init` command creates example .shotgunignore file

### Story 2.4: Binary File Detection and Handling
As a user,
I want binary files automatically identified and excluded from selection,
so that I don't accidentally include non-text files in my prompts.

#### Acceptance Criteria
1. Automatic binary detection using filetype library during scanning
2. Binary files visually distinct (🔒 icon, grayed appearance, non-selectable)
3. Common binary extensions pre-configured (.exe, .jpg, .png, .zip, .pdf, etc.)
4. Binary files excluded from selection count but visible in tree
5. Tooltip or status indication explaining why binary files cannot be selected
6. Override capability for edge cases where binary content is needed
7. Performance: detection adds <5% overhead to file scanning

### Story 2.5: Complete File Tree Screen Integration
As a user,
I want the complete file tree interface integrated as Screen 1 of the wizard,
so that I can select files and navigate to template selection seamlessly.

#### Acceptance Criteria
1. File tree screen integrated as first screen in 5-screen wizard flow
2. Header showing "File Selection [1/5]" with progress indicator
3. Status bar showing selection counts and keyboard shortcuts
4. F3 key advances to next screen when at least one file is selected
5. F1 key shows contextual help for file tree navigation
6. Screen state preservation when navigating back from subsequent screens
7. Validation: prevents advancement if no files are selected
8. Loading state during initial file scanning with progress indication

## Epic 3 - Template System & User Input

**Epic Goal:** Create comprehensive template discovery and management system for built-in and custom user templates, implement advanced multiline text editors with full UTF-8 support, and complete the 5-screen wizard workflow with seamless navigation and state management.

### Story 3.1: Template Discovery and Metadata System
As a user,
I want to see all available templates (built-in and custom) with rich metadata,
so that I can choose the most appropriate template for my current task.

#### Acceptance Criteria
1. Template discovery system scanning embedded templates and user config directories
2. Cross-platform config directory support (~/.config/shotgun-cli/templates, %APPDATA%/shotgun-cli/templates)
3. TOML metadata parsing for template name, version, description, author, and tags
4. Template validation ensuring required sections and variables are present
5. Unified template list showing both built-in and custom templates without visual distinction
6. Template precedence: user templates override built-in templates with same name
7. Error handling for malformed templates with user-friendly error messages

### Story 3.2: Template Selection Interface
As a user,
I want to select templates from a clean, navigable list interface,
so that I can quickly choose the right template and understand its purpose.

#### Acceptance Criteria
1. Vertical list interface using Bubbles list component with template metadata display
2. Template entries showing name, version, and description in structured format
3. Keyboard navigation with ↑↓ arrows and Enter/F3 for selection
4. Visual highlighting of currently selected template with focus indicators
5. Screen integration as "Template Selection [2/5]" in wizard flow
6. F2 navigation back to file tree with state preservation
7. Template preview or expanded description available via additional key press
8. Validation: prevents advancement without template selection

### Story 3.3: Advanced Multiline Text Editor for Task Input
As a user,
I want a sophisticated text editor for describing my task with full UTF-8 support,
so that I can provide detailed context in any language with special characters.

#### Acceptance Criteria
1. Multiline text editor built on Bubbles textarea v0.21.0 with horizontal scrolling
2. Full UTF-8 character support including accented characters (ç, á, ô, ñ, etc.)
3. Text input features: word wrap, line numbers, character/word count display
4. Clipboard integration (Ctrl+V paste, Ctrl+C copy) with proper UTF-8 handling
5. Editor modes: editing mode (normal text input) and navigation mode (F-key shortcuts)
6. Ctrl+Enter toggles between editing and navigation modes
7. Screen integration as "Task Description [3/5]" with state preservation
8. Validation: prevents advancement with empty task content

### Story 3.4: Optional Rules Input Editor
As a user,
I want an optional editor for additional rules and constraints,
so that I can provide specific guidance for the LLM when needed.

#### Acceptance Criteria
1. Optional multiline editor with same advanced features as task editor
2. Clear visual indication that field is optional with "optional" label
3. F4 key to skip this step entirely and advance to next screen
4. Auto-save functionality preserving content when navigating between screens
5. Screen integration as "Rules · optional [4/5]" with appropriate styling
6. Same UTF-8 and clipboard support as task editor
7. Graceful handling of empty content (treated as no additional rules)

### Story 3.5: Complete Wizard Flow Integration
As a user,
I want seamless navigation between all screens with preserved state,
so that I can review and modify any step without losing my work.

#### Acceptance Criteria
1. Complete 5-screen wizard flow: File Tree → Templates → Task → Rules → Confirmation
2. Global navigation: F2 (back), F3 (forward), F1 (help), ESC (exit with confirmation)
3. State preservation: all selections, text content, and UI states maintained during navigation
4. Progress indicator showing current screen (e.g., [3/5]) in all screen headers
5. Screen-specific validation preventing invalid forward navigation
6. Contextual help (F1) showing relevant shortcuts and instructions per screen
7. Consistent styling and theme across all screens using Lip Gloss monochrome palette
8. Error states and user feedback for invalid operations or missing required content

## Epic 4 - Prompt Generation & Output

**Epic Goal:** Develop the comprehensive template processing engine, concurrent file content aggregation system, real-time size estimation with progress visualization, and final Markdown output generation with robust error handling to deliver the complete MVP functionality.

### Story 4.1: Template Processing Engine
As a user,
I want the system to process templates with my inputs and selected files,
so that I can generate customized prompts based on the chosen template structure.

#### Acceptance Criteria
1. Template processing engine using Go's text/template package with custom function maps
2. Variable substitution system supporting all template placeholders (TASK, RULES, FILE_STRUCTURE, etc.)
3. Template function library including string manipulation (title, upper, lower, trim, wordCount, lineCount)
4. Conditional logic support in templates ({{if}}, {{range}}, {{with}}) for dynamic content
5. Template validation during processing with meaningful error messages for syntax issues
6. Safety measures preventing template injection or code execution vulnerabilities
7. Performance: template processing completes in <2 seconds for standard templates

### Story 4.2: Concurrent File Content Aggregation
As a user,
I want the system to efficiently read and organize selected file contents,
so that I can include comprehensive project context in my prompts.

#### Acceptance Criteria
1. Concurrent file reading using worker pool pattern with configurable concurrency
2. File content aggregation in structured format with file path headers
3. Directory tree generation using ASCII characters (├── └── │) for visual hierarchy
4. Content format: `<file path="RELATIVE/PATH">` blocks with actual file contents
5. Binary file exclusion with automatic detection and user notification
6. Large file handling with size limits and user warnings for oversized content
7. Performance: aggregate 100+ files in <10 seconds with progress indication
8. Error handling for file permission issues, missing files, or encoding problems

### Story 4.3: Real-Time Size Estimation System
As a user,
I want to see an accurate estimate of the final prompt size before generation,
so that I can make informed decisions about file inclusion and template complexity.

#### Acceptance Criteria
1. Real-time size calculation as user modifies selections or input content
2. Size estimation in human-readable units (KB, MB) with byte precision available
3. Progress bar visualization during size calculation with smooth animations
4. Warning thresholds for large prompts (>500KB warning, >1MB strong warning)
5. Breakdown showing contribution from different sections (template, files, inputs)
6. Performance impact estimation and recommendations for large file sets
7. Size recalculation triggers on file selection changes or input modifications

### Story 4.4: Confirmation Screen with Summary
As a user,
I want to review all my selections and see the complete summary before generation,
so that I can verify everything is correct before creating the final prompt.

#### Acceptance Criteria
1. Comprehensive summary screen showing template, selected files count, excluded items
2. Size estimate with visual progress bar and warning indicators
3. Generated filename preview with timestamp format (shotgun_prompt_YYYYMMDD_HHMM.md)
4. Final review of task description and rules (truncated preview if very long)
5. Clear action buttons: F2 (back for changes), F10 (generate), ESC (cancel)
6. Screen integration as "Confirm Generation [5/5]" completing the wizard
7. Warning messages for potential issues (very large size, no files selected, etc.)

### Story 4.5: Markdown Output Generation
As a user,
I want the system to generate a complete Markdown prompt file in the current directory,
so that I can immediately use the prompt with my preferred LLM tools.

#### Acceptance Criteria
1. Final prompt assembly combining template structure with processed content
2. Markdown file generation with timestamped filename in current working directory
3. File structure section with complete directory tree and file contents
4. Proper Markdown formatting with headers, code blocks, and structured sections
5. UTF-8 encoding support maintaining all special characters from inputs
6. Atomic file writing preventing partial files on interruption or errors
7. Success confirmation with file path display and file size information
8. Error handling for disk space issues, permission problems, or write failures
9. Generation progress indication during file writing process

## Epic 5 - Polish & Cross-Platform Optimization

**Epic Goal:** Refine the complete user interface with professional monochrome styling, implement session management and workflow history, optimize performance for large-scale repositories, and ensure robust cross-platform terminal compatibility for production deployment.

### Story 5.1: Professional UI Styling and Theme Implementation
As a user,
I want a polished, professional interface with consistent styling,
so that I have a pleasant, distraction-free experience that works beautifully across all terminals.

#### Acceptance Criteria
1. Complete monochrome theme implementation using Lip Gloss with specified color palette
2. Consistent styling across all screens: headers, borders, progress bars, text editors
3. Professional typography with proper spacing, alignment, and visual hierarchy
4. Accent color usage (mint green #6ee7b7) for focus states and highlights
5. Status indicators and progress bars with smooth animations and clear visual feedback
6. Terminal responsiveness handling various window sizes gracefully (minimum 80x24)
7. Theme consistency testing across major terminal emulators (PowerShell, bash, iTerm2, WezTerm)

### Story 5.2: Session Management and History System
As a user,
I want my previous sessions saved and easily restorable,
so that I can quickly repeat similar workflows or recover from interruptions.

#### Acceptance Criteria
1. Session persistence system saving user selections, template choices, and input content
2. Cross-platform session storage in appropriate config directories
3. Session history interface showing recent sessions with timestamps and summary info
4. Quick restore functionality allowing users to load and modify previous configurations
5. Session management: configurable history limit, cleanup of old sessions
6. Auto-save during workflow to prevent data loss on unexpected exits
7. Privacy controls: option to disable session saving or clear history

### Story 5.3: Performance Optimization for Large Repositories  
As a user,
I want the application to remain responsive even with very large codebases,
so that I can use it effectively on enterprise-scale projects.

#### Acceptance Criteria
1. Optimized file scanning with configurable concurrency and memory limits
2. Virtual scrolling in file tree for repositories with 10,000+ files
3. Lazy loading of file content with streaming processing to minimize memory usage
4. Background processing with progress indication for expensive operations
5. Configurable limits: max files to display, max file size to include, scan timeouts
6. Performance monitoring and user feedback for operations approaching limits
7. Graceful degradation when hitting system or configured limits
8. Memory usage optimization keeping total RAM usage under 200MB for large repos

### Story 5.4: Cross-Platform Terminal Compatibility
As a user,
I want the application to work flawlessly on all major operating systems and terminals,
so that I can use it consistently across different development environments.

#### Acceptance Criteria
1. Windows terminal compatibility: PowerShell, CMD, Windows Terminal, ConPTY support
2. Linux terminal support: bash, zsh, fish shells with various terminal emulators
3. macOS compatibility: Terminal.app, iTerm2, with proper key binding handling
4. Unicode and UTF-8 rendering consistency across all supported platforms
5. Keyboard shortcut handling accounting for OS-specific key combinations
6. Color support detection with graceful fallback for limited terminals
7. Integration testing on all target platforms with automated CI/CD validation

### Story 5.5: Advanced Configuration and CLI Integration
As a user,
I want comprehensive configuration options and CLI convenience features,
so that I can customize the application to fit my specific workflow needs.

#### Acceptance Criteria
1. Configuration system with YAML/TOML config files for user preferences
2. Command-line interface with Cobra framework supporting direct prompt generation
3. CLI flags for common operations: --template, --task, --output, --quick modes
4. Configuration options: default templates, file exclusions, UI preferences, performance tuning
5. Environment variable support for CI/CD integration and automation
6. `shotgun init` command for project setup with example .shotgunignore
7. `shotgun --help` with comprehensive usage documentation and examples
8. Integration with shell completion for enhanced CLI experience

## Checklist Results Report

*This section will contain the PM checklist validation results after executing the pm-checklist.*

## Next Steps

### UX Expert Prompt
Review this PRD to understand the monochrome TUI design requirements and create detailed interface mockups and interaction flows for the 5-screen wizard workflow.

### Architect Prompt
Use this PRD to design the technical architecture for the Go-based TUI application, focusing on concurrent file processing, Bubble Tea v2 integration, and cross-platform terminal compatibility.