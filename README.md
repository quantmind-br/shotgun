# 🔫 Shotgun CLI

An AI-powered terminal-based prompt engineering tool that helps you craft, manage, and optimize AI prompts efficiently.

## Features

- **Interactive Terminal UI**: Clean, keyboard-driven interface built with Bubble Tea
- **Prompt Management**: Organize and version control your AI prompts
- **Template System**: Reusable prompt templates for common use cases
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/user/shotgun-cli.git
cd shotgun-cli

# Build the application
make build

# Run the application
./bin/shotgun
```

### Using Go Install

```bash
go install github.com/user/shotgun-cli/cmd/shotgun@latest
```

## Usage

Launch the interactive TUI:

```bash
shotgun
```

### Keyboard Shortcuts

- `q` or `Ctrl+C` - Quit the application
- `?` - Show help (coming soon)

## Development

### Prerequisites

- Go 1.18 or higher
- Make (optional, for build automation)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run with coverage
make test-coverage

# Clean build artifacts
make clean
```

### Project Structure

```
shotgun-cli/
├── cmd/shotgun/          # Main application entry point
│   └── main.go
├── internal/             # Private application packages
│   ├── app/             # Main app controller
│   ├── components/      # Reusable UI components
│   ├── core/           # Business logic
│   │   └── scanner/    # File scanning utilities
│   ├── models/         # Data models
│   └── screens/        # TUI screens
├── templates/          # Built-in prompt templates
├── go.mod             # Go module definition
├── go.sum            # Dependency checksums
├── Makefile          # Build automation
└── README.md         # Project documentation
```

### Dependencies

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[TOML](https://github.com/BurntSushi/toml)** - Configuration parsing

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and use `gofmt`
- Write tests for new functionality
- Maintain minimum 90% test coverage for core packages
- Use meaningful commit messages

## Architecture

This project follows the Elm Architecture pattern with:

- **Model**: Immutable application state
- **View**: Pure functions that render the UI
- **Update**: State transitions based on messages

The TUI is built using the Bubble Tea framework, providing a reactive, component-based architecture similar to modern web frameworks but for terminal applications.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Prompt template management
- [ ] AI provider integrations
- [ ] Prompt versioning and history
- [ ] Export/import functionality
- [ ] Plugin system
- [ ] Configuration management

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/user/shotgun-cli/issues) on GitHub.