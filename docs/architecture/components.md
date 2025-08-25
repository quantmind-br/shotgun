# Components

### TUI Application Shell

**Responsibility:** Main application orchestration, Bubble Tea program lifecycle, and global state management

**Key Interfaces:**
- `tea.Model` implementation for root application state
- Global key binding handlers (ESC, F1-F10)
- Screen navigation and wizard flow control

**Dependencies:** All UI screen components, StateManager

**Technology Stack:** Bubble Tea v2 Elm Architecture, Lip Gloss styling, Cobra CLI integration

### File Tree Component

**Responsibility:** Hierarchical file browser with checkbox selection, tree expansion/collapse, and real-time file discovery

**Key Interfaces:**
- `tea.Model` for file tree state and rendering
- Selection change events and validation
- Tree navigation and keyboard shortcuts

**Dependencies:** FileScanner, IgnoreProcessor, FileNode data models

**Technology Stack:** Bubbles filepicker extended with custom checkbox logic, concurrent file scanning

### Template Selector Component

**Responsibility:** Template discovery, metadata display, and user selection interface

**Key Interfaces:**
- Template list rendering with rich metadata
- Selection events and template preview
- Built-in vs custom template differentiation

**Dependencies:** TemplateEngine, Template data models

**Technology Stack:** Bubbles list component with custom item rendering, TOML metadata parsing

### Multiline Text Editor Component

**Responsibility:** Advanced text input with UTF-8 support, clipboard integration, and editing modes

**Key Interfaces:**
- Text editing events and content validation
- Mode switching (edit vs navigation)
- Character/word counting and status display

**Dependencies:** None (self-contained UI component)

**Technology Stack:** Bubbles textarea v0.21.0 with UTF-8 handling, clipboard integration

### File Scanner Engine

**Responsibility:** Concurrent directory traversal, binary detection, and ignore rule processing

**Key Interfaces:**
- `ScanDirectory(path string) (*FileNode, error)`
- `ApplyIgnoreRules(root *FileNode, rules []string) error`
- Progress reporting through channels

**Dependencies:** IgnoreProcessor, doublestar glob matching, filetype detection

**Technology Stack:** Go worker pools, channels, doublestar v4, filetype library

### Template Processing Engine

**Responsibility:** Template compilation, variable substitution, and content generation

**Key Interfaces:**
- `ProcessTemplate(template *Template, data *TemplateData) (string, error)`
- Template validation and custom function registration
- Error handling and debugging information

**Dependencies:** Go text/template, custom function library

**Technology Stack:** Native Go text/template with security-safe processing

### Ignore Rules Processor

**Responsibility:** .gitignore/.shotgunignore parsing, pattern matching, and file exclusion logic

**Key Interfaces:**
- `LoadRules(gitignorePath, shotgunignorePath string) error`
- `ShouldIgnore(path string, isDir bool) bool`
- Rule precedence and conflict resolution

**Dependencies:** doublestar pattern matching

**Technology Stack:** doublestar v4 glob patterns, custom precedence logic

### Output Generator

**Responsibility:** Final prompt assembly, Markdown formatting, and file writing with atomic operations

**Key Interfaces:**
- `GeneratePrompt(generation *PromptGeneration) error`
- Size estimation and progress reporting
- File system operations with error recovery

**Dependencies:** Template processing results, selected file contents

**Technology Stack:** Go text/template, Markdown formatting, atomic file operations

### Session Manager

**Responsibility:** Application state persistence, session history, and restoration capabilities

**Key Interfaces:**
- `SaveSession(state *ApplicationState) error`
- `LoadSession(sessionID string) (*ApplicationState, error)`
- Cross-platform config directory management

**Dependencies:** ApplicationState serialization, Viper configuration

**Technology Stack:** JSON serialization, Viper config management, cross-platform file paths

### Configuration Manager

**Responsibility:** User preferences, environment variables, and cross-platform configuration

**Key Interfaces:**
- Configuration loading and validation
- Default value management
- Environment variable override support

**Dependencies:** Viper configuration library

**Technology Stack:** Viper v1.18+, TOML/YAML config formats
