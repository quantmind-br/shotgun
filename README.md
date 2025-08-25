# Shotgun CLI v3

A fast, interactive terminal-based file scanner and template processor built with Go and Bubble Tea.

## Overview

Shotgun is a modern Terminal User Interface (TUI) application that provides an intuitive way to scan files, process templates, and generate code. Built with performance and usability in mind, it offers a rich interactive experience while maintaining the speed and efficiency of a command-line tool.

### Key Features

- **Interactive TUI**: Modern terminal interface powered by Bubble Tea v2
- **Fast File Scanning**: Efficiently scan large directories with intelligent filtering
- **Template Processing**: Support for Go templates with custom functions
- **Cross-Platform**: Native builds for Windows, Linux, and macOS
- **Git Integration**: Respects .gitignore and .shotgunignore patterns
- **Performance Optimized**: Handle 1000+ files in under 5 seconds

## Requirements

- Go 1.22 or later
- Terminal with 256 color support (recommended)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/your-org/shotgun-cli-v3.git
cd shotgun-cli-v3

# Build and install
make install
```

### Quick Build

```bash
# Build for current platform
make build

# Cross-compile for all platforms
make cross-compile
```

## Usage

### Basic Commands

```bash
# Show help
shotgun --help

# Launch interactive mode (coming in future stories)
shotgun

# Process templates (coming in future stories)
shotgun template --input ./templates --output ./generated
```

### Development Usage

```bash
# Development build with tests
make dev

# Run tests only
make test

# Run linter
make lint

# View all available commands
make help
```

## Project Structure

```
shotgun-cli-v3/
├── cmd/shotgun/           # Main application entry point
├── internal/              # Private application packages
│   ├── ui/               # TUI components and screens
│   ├── core/             # Business logic (scanner, template, etc.)
│   ├── services/         # Application services layer
│   ├── models/           # Shared data structures
│   ├── infrastructure/   # External dependencies abstraction
│   └── testutil/         # Testing utilities
├── templates/            # Template assets
│   ├── embedded/         # Built-in templates
│   └── examples/         # Example templates
├── configs/              # Configuration files
├── docs/                 # Documentation
├── test/                 # Integration and E2E tests
└── scripts/              # Build and development scripts
```

## Development Setup

### Prerequisites

- Go 1.22+
- Make (for build automation)
- Git

### Getting Started

1. **Clone and setup**:
   ```bash
   git clone https://github.com/your-org/shotgun-cli-v3.git
   cd shotgun-cli-v3
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Run development build**:
   ```bash
   make dev
   ```

4. **Run tests**:
   ```bash
   make test
   ```

### Code Quality

This project follows strict coding standards:

- **Error Handling**: All public functions return error as last parameter
- **Context Usage**: Long-running operations accept context.Context
- **Interface Design**: Program to interfaces for testability
- **Resource Cleanup**: Explicit cleanup using defer or context cancellation
- **Naming Conventions**: Follow Go best practices with specific TUI patterns

### Testing

```bash
# Run unit tests
make test

# Run benchmarks
make bench

# Generate coverage report
make coverage

# Run all quality checks
make check
```

## Technology Stack

| Component | Technology | Version | Purpose |
|-----------|------------|---------|---------|
| Language | Go | 1.22+ | Core application development |
| TUI Framework | Bubble Tea | v2.0.0-beta.4 | Reactive terminal interface |
| UI Components | Bubbles | v0.21.0 | Pre-built TUI components |
| Styling | Lip Gloss | v1.0.0 | Advanced terminal styling |
| CLI Framework | Cobra | v1.8+ | Command line interface |
| Config Management | Viper | v1.18+ | Configuration handling |
| Template Engine | text/template | Go stdlib | Template processing |
| File Matching | doublestar | v4.6+ | Glob pattern support |
| File Detection | filetype | v1.1+ | Binary file detection |

## Build Targets

| Target | Description |
|--------|-------------|
| `make build` | Build for current platform |
| `make test` | Run all tests |
| `make clean` | Remove build artifacts |
| `make lint` | Run code quality checks |
| `make cross-compile` | Build for all platforms |
| `make install` | Install locally |
| `make dev` | Development build (deps + build + test) |

## Contributing

We welcome contributions! Please follow these guidelines:

1. **Code Style**: Follow the established coding standards in `docs/architecture/coding-standards.md`
2. **Testing**: Include tests for new functionality
3. **Documentation**: Update documentation for user-facing changes
4. **Performance**: Ensure changes don't significantly impact scan performance

### Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make check` to verify quality
5. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs and request features on [GitHub Issues](https://github.com/your-org/shotgun-cli-v3/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/your-org/shotgun-cli-v3/discussions)
- **Documentation**: Full documentation in the [docs/](docs/) directory

## Roadmap

- ✅ **Story 1.1**: Project setup and Go module initialization
- 🔄 **Story 1.2**: Core file scanning engine (planned)
- 🔄 **Story 1.3**: TUI framework and basic screens (planned)
- 🔄 **Story 1.4**: Template processing system (planned)

---

*Built with ❤️ using Go and Bubble Tea*