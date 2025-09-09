# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Building
```bash
# Build for current platform
make build
# or directly:
go build -ldflags="-s -w" -o bin/shotgun ./cmd/shotgun

# Build for all platforms
make build-all

# Build for specific platforms
make build-linux    # Linux AMD64
make build-windows  # Windows AMD64  
make build-darwin   # macOS Intel + Apple Silicon
```

### Testing
```bash
# Run all tests
make test
# or:
go test -v ./...

# Run tests with coverage
make test-coverage
# or:
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run tests for a specific package
go test -v ./internal/core/scanner
go test -v ./internal/screens/filetree

# Run a single test
go test -v -run TestSpecificFunction ./internal/core/scanner
```

### Code Quality
```bash
# Format code
make fmt
# or:
go fmt ./...

# Run go vet
make vet
# or:
go vet ./...

# Run both fmt and vet
make lint
```

### Running in Development
```bash
# Run the TUI application
make run
# or:
go run ./cmd/shotgun

# Run with specific command
go run ./cmd/shotgun init
go run ./cmd/shotgun version --json
```

### Dependencies
```bash
# Update and tidy dependencies
make deps
# or:
go mod tidy
go mod download
```

### Installation
```bash
# Install to GOPATH/bin using make (recommended on Windows)
make install

# Alternative: use go install directly
go install -ldflags="-s -w" -buildvcs=false ./cmd/shotgun

# Build and install locally (fallback method)
make install-local
```

## Architecture Overview

### Core Design Pattern: Elm Architecture
The application follows the Elm Architecture (Model-View-Update) pattern using the Bubble Tea framework:
- **Model**: Immutable application state stored in `AppState` and screen-specific models
- **View**: Pure rendering functions that display the UI based on the model
- **Update**: Message-based state transitions that handle user input and system events

### Project Structure

#### Main Entry Points
- `cmd/shotgun/main.go` - CLI entry point using Cobra for command routing
- `cmd/shotgun/program_options_*.go` - Platform-specific Bubble Tea configuration

#### Core Components

**`internal/app/`** - Main TUI application controller
- `model.go` - Central `AppState` struct managing screen navigation and state
- `update.go` - Main update loop handling global messages and screen transitions
- `view.go` - Top-level view composition
- `keys.go` - Global keyboard handling and normalization

**`internal/screens/`** - Individual TUI screens following Model-View-Update pattern
- `filetree/` - File selection screen with directory scanning
- `template/` - Template selection and management
- `input/` - User input for task names and rules
- `confirm/` - Confirmation screen showing generation summary
- `generate/` - Progress screen during file generation

**`internal/core/`** - Business logic layer
- `scanner/` - File system scanning with `.shotgunignore` support
  - Binary file detection
  - Concurrent scanning with worker pools
  - Ignore pattern matching (gitignore-style)
- `template/` - Template engine and management
  - TOML-based template definitions
  - Template discovery from multiple sources
  - Built-in template functions
  - Variable substitution engine
- `builder/` - Project structure generation
  - File tree construction from scanned files
  - Progress tracking and estimation
  - Concurrent file writing

**`internal/cli/`** - CLI-specific commands (non-TUI)
- `init.go` - Init command for creating `.shotgunignore` files
- `templates/` - Embedded default templates

**`internal/components/`** - Reusable UI components
- `help/` - Help dialog system
- `spinner/` - Loading indicators
- `progress/` - Progress bars

### Key Architectural Decisions

1. **Screen State Management**: Each screen maintains its own state, with `AppState` coordinating transitions and preserving state between screens.

2. **Concurrent Operations**: File scanning and generation use worker pools for performance, with proper context cancellation.

3. **Template System**: Templates are discovered from multiple sources (built-in, user directory) and parsed from TOML files with a custom template engine.

4. **Cross-Platform Support**: Platform-specific code is isolated in build-tagged files (e.g., `program_options_windows.go`).

5. **Keyboard Navigation**: Centralized key normalization to handle cross-platform differences, with Ctrl-based shortcuts for consistency.

### Screen Flow
1. **FileTree Screen** → User selects files/directories to include
2. **Template Screen** → User selects a project template  
3. **Input Screen** → User provides task name and optional rules
4. **Confirm Screen** → Shows summary and output filename
5. **Generate Screen** → Displays progress during generation

### Testing Strategy
- Unit tests alongside implementation files (`*_test.go`)
- Integration tests in `internal/integration/`
- End-to-end tests in `e2e/`
- Platform-specific tests with build tags
- Minimum 60% coverage target for core packages

### Important Patterns
- Use `tea.Batch()` for combining multiple commands
- Implement `tea.Model` interface for all screens
- Use `context.Context` for cancellation in long-running operations
- Normalize keyboard input through `normalizeKey()` function
- Handle Windows-specific terminal requirements through build tags