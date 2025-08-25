# Tech Stack

### Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| Primary Language | Go | 1.22+ | Core application development | Native compilation, excellent concurrency, cross-platform support, no runtime dependencies |
| TUI Framework | Bubble Tea | v2.0.0-beta.4 | Reactive terminal interface | Mature Elm Architecture implementation, excellent Windows support, active development |
| UI Components | Bubbles | v0.21.0 | Pre-built TUI components | Proven components (filepicker, textarea, list) with horizontal scrolling support |
| Terminal Styling | Lip Gloss | v1.0.0 | Advanced terminal styling | Sophisticated styling with gradients, layouts, consistent cross-platform rendering |
| Template Engine | text/template | Go stdlib | Template processing | Built-in security, custom function maps, no external dependencies |
| File Processing | doublestar | v4.6+ | Glob pattern matching | .gitignore/.shotgunignore pattern support with proper escaping |
| File Type Detection | filetype | v1.1+ | Binary file detection | Fast magic number detection for excluding binary files |
| Configuration | Viper | v1.18+ | Config management | Cross-platform config directories, environment variable support |
| CLI Framework | Cobra | v1.8+ | Command line interface | Standard Go CLI framework with shell completion support |
| Data Format | TOML | BurntSushi/toml v1.3+ | Template metadata | Human-readable config format, better than YAML for metadata |
| Build Tool | Go toolchain | 1.22+ | Native build system | Cross-compilation, module management, integrated testing |
| Task Runner | Makefile | GNU Make | Build automation | Simple, universal automation for common tasks |
| Version Control | Git | 2.40+ | Source control | Industry standard, required for .gitignore processing |
| CI/CD | GitHub Actions | Latest | Automated builds | Native GitHub integration, matrix builds for cross-platform |
| Testing Framework | testing | Go stdlib | Unit testing | Native Go testing with table-driven tests |
| TUI Testing | teatest | v0.1+ | TUI interaction testing | Bubble Tea testing helpers for UI workflows |
| Benchmarking | testing | Go stdlib | Performance testing | Built-in benchmarking for file scanning optimization |
| Linting | golangci-lint | v1.55+ | Code quality | Comprehensive linter suite with performance checks |
| Cross-Platform | Native Go | 1.22+ | Platform compatibility | Built-in cross-compilation without CGO dependencies |
