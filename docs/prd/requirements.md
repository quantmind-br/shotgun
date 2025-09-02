# Requirements

## Functional

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

## Non Functional

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
