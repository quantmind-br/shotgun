# Shotgun CLI Makefile

# Variables
BINARY_NAME=shotgun
BINARY_PATH=bin/$(BINARY_NAME)
CMD_PATH=./cmd/shotgun
BUILD_FLAGS=-ldflags="-s -w"

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@mkdir -p bin
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) $(CMD_PATH)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Cross-platform builds
.PHONY: build-all
build-all: build-linux build-windows build-darwin

.PHONY: build-linux
build-linux:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)

.PHONY: build-windows
build-windows:
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)

.PHONY: build-darwin
build-darwin:
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)

# Development targets
.PHONY: run
run:
	go run $(CMD_PATH)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: fmt vet

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Remove build artifacts"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  build-darwin - Build for macOS (Intel + Apple Silicon)"
	@echo "  run          - Run the application"
	@echo "  fmt          - Format Go code"
	@echo "  vet          - Run Go vet"
	@echo "  lint         - Run fmt and vet"
	@echo "  deps         - Install/update dependencies"
	@echo "  help         - Show this help"