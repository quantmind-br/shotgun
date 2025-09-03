# Tech Stack and Dependencies

## Core Technologies
- **Language**: Go 1.18+ (targeting 1.22+ based on tech stack docs)
- **Frontend Framework**: Bubble Tea v1.1.0 (TUI framework)
- **UI Components**: Bubbles v0.20.0 (pre-built TUI components)
- **Styling**: Lip Gloss v1.0.0 (terminal styling)
- **CLI Framework**: Cobra v1.10.1

## Key Dependencies
- **File Type Detection**: github.com/h2non/filetype v1.1.3
- **Configuration**: TOML support via BurntSushi/toml v1.5.0
- **Concurrency**: Standard library sync and golang.org/x/sync
- **Terminal Support**: Various charmbracelet/* packages for cross-platform terminal handling

## Architecture Pattern
- **Elm Architecture**: Model-View-Update pattern
- **State Management**: Immutable state with Bubble Tea models
- **Concurrency**: Go channels and sync primitives
- **Testing**: Standard Go testing library with 90% coverage requirement

## Build System
- **Build Tool**: Go build with Makefile automation
- **Target Platforms**: Linux, macOS (Intel/ARM), Windows
- **Output**: Single binary executable