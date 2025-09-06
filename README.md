# 🔫 Shotgun CLI

A fast project scaffolding tool that helps you quickly generate project structures from templates, with an intuitive terminal interface and flexible file ignore management.

## Features

- **Interactive Terminal UI**: Clean, keyboard-driven interface built with Bubble Tea
- **Project Scaffolding**: Fast project structure generation from templates
- **File Ignore Management**: Easy creation and management of .shotgunignore files
- **Cross-Platform**: Works on Linux, macOS, and Windows

## 📥 Installation

### Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [GitHub Releases](https://github.com/diogopedro/shotgun/releases) page.

#### Linux / macOS / FreeBSD

```bash
# Download and extract (replace with your platform/architecture)
wget https://github.com/diogopedro/shotgun/releases/latest/download/shotgun-linux-amd64.tar.gz
tar -xzf shotgun-linux-amd64.tar.gz

# Verify checksum (recommended)
sha256sum -c shotgun-linux-amd64.sha256

# Make executable and install
chmod +x shotgun-linux-amd64
sudo mv shotgun-linux-amd64 /usr/local/bin/shotgun

# Verify installation
shotgun version
```

#### Windows (PowerShell)

```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/diogopedro/shotgun/releases/latest/download/shotgun-windows-amd64.zip" -OutFile "shotgun-windows-amd64.zip"
Expand-Archive -Path "shotgun-windows-amd64.zip" -DestinationPath "."

# Verify checksum (recommended)
$expected = Get-Content shotgun-windows-amd64.exe.sha256
$actual = Get-FileHash shotgun-windows-amd64.exe -Algorithm SHA256
if ($expected.Split(" ")[0] -eq $actual.Hash.ToLower()) { Write-Host "✅ Checksum verified" } else { Write-Host "❌ Checksum failed" }

# Add to PATH or run directly
./shotgun-windows-amd64.exe version
```

### Platform Support Matrix

| Platform | Architecture | Binary | Archive | Status |
|----------|-------------|---------|---------|--------|
| **Linux** | x86_64 | `shotgun-linux-amd64` | `shotgun-linux-amd64.tar.gz` | ✅ Supported |
| **Linux** | ARM64 | `shotgun-linux-arm64` | `shotgun-linux-arm64.tar.gz` | ✅ Supported |
| **macOS** | Intel | `shotgun-darwin-amd64` | `shotgun-darwin-amd64.tar.gz` | ✅ Supported |
| **macOS** | Apple Silicon | `shotgun-darwin-arm64` | `shotgun-darwin-arm64.tar.gz` | ✅ Supported |
| **Windows** | x86_64 | `shotgun-windows-amd64.exe` | `shotgun-windows-amd64.zip` | ✅ Supported |
| **Windows** | ARM64 | `shotgun-windows-arm64.exe` | `shotgun-windows-arm64.zip` | ✅ Supported |
| **FreeBSD** | x86_64 | `shotgun-freebsd-amd64` | `shotgun-freebsd-amd64.tar.gz` | ✅ Supported |

### Alternative Installation Methods

#### Using Go Install

```bash
go install github.com/diogopedro/shotgun/cmd/shotgun@latest
```

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/diogopedro/shotgun.git
cd shotgun

# Build the application (requires Go 1.22+)
go build -o shotgun ./cmd/shotgun

# Run the application
./shotgun
```

### Terminal Compatibility

Shotgun automatically detects and adapts to your terminal capabilities:

- 🎨 **True Color** support (16.7 million colors)
- 🌈 **256 Color** terminal support
- 🔧 **16 Color** terminal fallback
- ⚫ **Monochrome** terminal fallback
- 🔤 **Unicode** character support with ASCII fallbacks
- ⌨️ **Cross-platform** keyboard handling

### Environment Variables

Configure terminal behavior:

- `FORCE_COLOR=0|1|2|3` - Force specific color support level
- `NO_COLOR=1` - Disable all colors
- `TERM` - Automatically detected for terminal capabilities
- `COLORTERM=truecolor` - Enhanced color support detection
- `SHOTGUN_CONFIG_DIR` - Custom configuration directory

### 🔐 Security & Verification

For security, always verify downloaded binaries:

#### Linux/macOS/FreeBSD
```bash
# Verify checksum matches
sha256sum -c shotgun-linux-amd64.sha256
```

#### Windows (PowerShell)
```powershell
# Verify checksum matches
$expected = (Get-Content shotgun-windows-amd64.exe.sha256).Split(" ")[0]
$actual = (Get-FileHash shotgun-windows-amd64.exe -Algorithm SHA256).Hash.ToLower()
if ($expected -eq $actual) { Write-Host "✅ Verified" } else { Write-Host "❌ Failed" }
```

## Usage

### Interactive TUI Mode (Default)

Launch the interactive terminal interface:

```bash
shotgun
```

### CLI Commands

Shotgun also provides direct CLI commands for quick operations:

#### Initialize .shotgunignore file

Create a .shotgunignore file in your project to customize which files are excluded during scanning:

```bash
# Create .shotgunignore with default patterns
shotgun init

# Force overwrite existing .shotgunignore
shotgun init --force

# Show help for init command
shotgun init --help
```

#### Version Information

```bash
# Human-readable format
shotgun version

# JSON format (for scripts/automation)
shotgun version --json
```

Example output:
```bash
$ shotgun version
Shotgun v1.0.0
Build time: 2024-12-06_10:30:15
Git commit: abc123def456
Platform: linux/amd64
Go version: go1.22.0

$ shotgun version --json
{
  "version": "v1.0.0",
  "build_time": "2024-12-06_10:30:15",
  "git_commit": "abc123def456", 
  "platform": "linux/amd64",
  "go_version": "go1.22.0"
}
```

#### Help

```bash
shotgun --help
```

### .shotgunignore File Format

The .shotgunignore file uses the same syntax as .gitignore files and supports:

- **Wildcard patterns**: `*.log`, `temp*`
- **Directory exclusions**: `node_modules/`, `build/`  
- **Negation patterns**: `!important.log`
- **Comments**: Lines starting with `#`

Example .shotgunignore:
```gitignore
# Build artifacts
build/
dist/
*.exe

# Dependencies
node_modules/
vendor/

# OS files
.DS_Store
Thumbs.db

# Logs (but keep important ones)
*.log
!important.log
```

### Keyboard Shortcuts (TUI Mode)

- `q` or `Ctrl+C` - Quit the application
- `?` - Show help (coming soon)

## Development

### Prerequisites

- Go 1.22 or higher (1.23+ recommended)
- Make (optional, for build automation)
- Git

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
│   ├── main.go          # CLI entry point with Cobra commands
│   └── main_test.go     # CLI tests
├── internal/             # Private application packages
│   ├── app/             # Main TUI application controller
│   ├── cli/             # CLI commands and logic
│   │   ├── init.go      # Init command implementation
│   │   ├── init_test.go # Init command tests
│   │   └── templates/   # Template definitions
│   ├── components/      # Reusable UI components
│   ├── core/           # Business logic
│   │   └── scanner/    # File scanning with .shotgunignore support
│   ├── models/         # Data models
│   └── screens/        # TUI screens
├── e2e/                # End-to-end integration tests
├── templates/          # Built-in project templates
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
- Maintain minimum 60% test coverage for core packages
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

### Completed ✅
- [x] CLI framework with Cobra integration
- [x] .shotgunignore file creation and management
- [x] File scanner integration with ignore patterns
- [x] Comprehensive test coverage (>90%)

### In Progress 🚧
- [ ] Project template management
- [ ] Enhanced file scanning capabilities

### Planned 📋
- [ ] Custom template creation
- [ ] Export/import functionality  
- [ ] Plugin system
- [ ] Configuration management
- [ ] Template versioning and history

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/diogopedro/shotgun/issues) on GitHub.

## 🔐 Security

All release binaries are built using GitHub Actions with:
- Static linking (CGO disabled) for maximum compatibility
- Build flags `-ldflags="-s -w"` for optimized binaries  
- SHA256 checksums for integrity verification
- Reproducible builds from tagged releases

For security issues, please email security@diogopedro.dev instead of opening a public issue.