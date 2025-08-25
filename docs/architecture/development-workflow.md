# Development Workflow

### Local Development Setup

#### Prerequisites

```bash
# Go 1.22+ installation
curl -OL https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Development tools
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/vektra/mockery/v2@latest

# Optional: Air for hot reloading during development
go install github.com/cosmtrek/air@latest
```

#### Initial Setup

```bash
# Clone and setup project
git clone https://github.com/org/shotgun-cli-v3.git
cd shotgun-cli-v3

# Initialize Go module and download dependencies
go mod download
go mod verify

# Generate mocks for testing
make generate

# Build embedded templates
make embed-templates

# Run initial tests to verify setup
make test

# Build development binary
make build-dev
```

#### Development Commands

```bash
# Start development with hot reload (if using Air)
make dev

# Or standard development build and run
make build-dev && ./bin/shotgun-dev

# Run specific screen tests
make test-ui

# Run core logic tests only
make test-core

# Run full test suite with coverage
make test-coverage

# Lint and format code
make lint
make fmt

# Cross-platform build (all targets)
make build-all

# Run benchmarks
make bench

# Generate documentation
make docs
```

### Environment Configuration

#### Required Environment Variables

```bash
# Development (.env.local)
# Application configuration
SHOTGUN_LOG_LEVEL=debug
SHOTGUN_TEMPLATES_DIR=./templates/embedded
SHOTGUN_CONFIG_DIR=./.shotgun-dev

# Testing configuration
SHOTGUN_TEST_MODE=true
SHOTGUN_TEST_FIXTURES=./test/fixtures

# Build configuration
CGO_ENABLED=0
GOOS=linux  # or darwin, windows
GOARCH=amd64

# Optional: Enable Bubble Tea debug mode
TEA_DEBUG=1

# Optional: Performance profiling
SHOTGUN_ENABLE_PROFILING=true
SHOTGUN_PROFILE_OUTPUT=./profiles
```
