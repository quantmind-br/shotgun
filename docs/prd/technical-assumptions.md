# Technical Assumptions

## Repository Structure: Monorepo
The application will be developed as a single Go repository containing all source code, templates, and documentation.

## Service Architecture
**Monolithic CLI Application** - Single binary executable with modular internal architecture using Go packages for separation of concerns. The application follows Elm Architecture patterns via Bubble Tea for state management.

## Testing Requirements
**Full Testing Pyramid** approach:
- Unit tests for all business logic components (90% coverage target)
- Integration tests for TUI flows using teatest framework
- E2E tests using expect scripts for complete user workflows
- Performance benchmarks for critical operations
- Cross-platform CI testing on Windows, Linux, and macOS

## Additional Technical Assumptions and Requests
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
