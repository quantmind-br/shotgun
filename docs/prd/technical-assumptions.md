# Technical Assumptions

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
