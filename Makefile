# Makefile for Shotgun CLI v3
# Build automation for Go TUI application

# Variables
BINARY_NAME=shotgun
CMD_DIR=cmd/shotgun
BUILD_DIR=dist
BINARY_UNIX=$(BUILD_DIR)/$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BUILD_DIR)/$(BINARY_NAME)_windows.exe
BINARY_DARWIN=$(BUILD_DIR)/$(BINARY_NAME)_darwin
GO_VERSION=1.22
LDFLAGS=-ldflags="-s -w"

# Detect OS for local builds
ifeq ($(OS),Windows_NT)
    BINARY_LOCAL=$(BUILD_DIR)/$(BINARY_NAME).exe
else
    BINARY_LOCAL=$(BUILD_DIR)/$(BINARY_NAME)
endif

.PHONY: help build test clean lint cross-compile install dev deps tidy check

# Default target
help: ## Show this help message
	@echo "Shotgun CLI v3 - Build Commands"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application for current platform
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BINARY_LOCAL) ./$(CMD_DIR)
	@echo "Built: $(BINARY_LOCAL)"

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...
	@echo "Running tests with race detector..."
	go test -race ./...

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	go clean
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@echo "Clean complete."

lint: ## Run code quality checks
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "Warning: golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Running basic vet and fmt checks..."; \
		go vet ./...; \
		go fmt ./...; \
	fi

# Cross-compilation targets
cross-compile: ## Build for all supported platforms
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_UNIX) ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_WINDOWS) ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DARWIN) ./$(CMD_DIR)
	@echo "Cross-compilation complete:"
	@echo "  Linux:   $(BINARY_UNIX)"
	@echo "  Windows: $(BINARY_WINDOWS)"
	@echo "  macOS:   $(BINARY_DARWIN)"

install: build ## Install the application locally
	@echo "Installing $(BINARY_NAME)..."
	go install ./$(CMD_DIR)
	@echo "Installation complete. $(BINARY_NAME) is now available in your PATH."

dev: deps build test ## Development build (deps, build, test)
	@echo "Development build complete."

deps: ## Download and verify dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	@echo "Dependencies updated."

tidy: ## Clean up go.mod and go.sum
	@echo "Tidying go modules..."
	go mod tidy
	@echo "Module cleanup complete."

check: lint test ## Run all checks (lint and test)
	@echo "All checks passed."

# Version and environment info
version: ## Show Go version and module info
	@echo "Go version: $(shell go version)"
	@echo "Module: $(shell head -1 go.mod)"
	@echo "Dependencies:"
	@go list -m all

# Benchmark targets
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Coverage targets
coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"