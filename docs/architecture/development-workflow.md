# Development Workflow

## Local Development Setup

### Prerequisites
```bash
# Install Go 1.22+
brew install go  # macOS
# or
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz  # Linux

# Verify installation
go version

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### Initial Setup
```bash
# Clone repository
git clone https://github.com/user/shotgun-cli.git
cd shotgun-cli

# Download dependencies
go mod download

# Run tests
make test

# Build binary
make build
```

### Development Commands
```bash
# Start all services
make run

# Start frontend only
go run ./cmd/shotgun

# Start backend only
# N/A - monolithic application

# Run tests
make test           # All tests
make test-unit      # Unit tests only
make test-e2e       # E2E tests only
make test-cover     # With coverage
```

## Environment Configuration

### Required Environment Variables
```bash
# Frontend (.env.local)
# N/A - No frontend env vars needed

# Backend (.env)
DEBUG=true          # Enable debug logging
LOG_LEVEL=debug     # Log verbosity

# Shared
# N/A - No shared env vars needed
```
