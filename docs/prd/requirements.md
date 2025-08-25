# Requirements

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
