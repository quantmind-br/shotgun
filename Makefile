# Shotgun CLI Makefile

# Variables
BINARY_NAME=shotgun
BINARY_PATH=bin/$(BINARY_NAME)
CMD_PATH=./cmd/shotgun
BUILD_FLAGS=-ldflags="-s -w" -buildvcs=false

# Export Go environment variables for Windows compatibility
export GOPATH := $(shell go env GOPATH)
export GOMODCACHE := $(shell go env GOMODCACHE)
export GOCACHE := $(shell go env GOCACHE)

# OS detection
ifeq ($(OS),Windows_NT)
    BINARY_EXT=.exe
    # Use the same path as go install (GOPATH/bin)
    GOPATH_DIR=$(shell go env GOPATH)
    INSTALL_PATH=$(GOPATH_DIR)/bin
    MKDIR=mkdir -p
    CP=cp -f
    RM=rm -f
else
    BINARY_EXT=
    # On Unix, check if GOPATH/bin exists and is in PATH, otherwise use /usr/local/bin
    GOPATH_DIR=$(shell go env GOPATH)
    ifneq ($(GOPATH_DIR),)
        INSTALL_PATH=$(GOPATH_DIR)/bin
    else
        INSTALL_PATH=/usr/local/bin
    endif
    MKDIR=mkdir -p
    CP=install -m 755
    RM=rm -f
endif

BINARY_PATH_OS=bin/$(BINARY_NAME)$(BINARY_EXT)

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@$(MKDIR) bin
	go build $(BUILD_FLAGS) -o $(BINARY_PATH_OS) $(CMD_PATH)

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

# Install binary to system using PowerShell (works on Windows)
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME) using PowerShell script..."
	@powershell.exe -ExecutionPolicy Bypass -File install.ps1 || (echo "PowerShell failed, try: make install-local" && exit 1)

# Install using PowerShell script (Windows)
.PHONY: install-ps
install-ps:
	powershell.exe -ExecutionPolicy Bypass -File install.ps1

# Install using Bash script (Unix/Linux)  
.PHONY: install-sh
install-sh:
	bash ./install.sh

# Cross-platform install (chooses appropriate script)
.PHONY: install-go
install-go:
ifeq ($(OS),Windows_NT)
	$(MAKE) install-ps
else
	$(MAKE) install-sh
endif

# Install binary by building first (alternative method)
.PHONY: install-local
install-local: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@echo "(Same location as 'go install')"
	@$(MKDIR) $(INSTALL_PATH)
ifeq ($(OS),Windows_NT)
	@$(CP) $(BINARY_PATH_OS) $(INSTALL_PATH)/$(BINARY_NAME)$(BINARY_EXT)
	@echo "✓ Installation complete!"
	@echo ""
	@echo "Note: This is the standard Go binary location (GOPATH/bin)"
	@echo "If $(BINARY_NAME) command is not found, add to PATH:"
	@echo "  PowerShell: $$env:Path += ';$(INSTALL_PATH)'"
	@echo "  CMD: setx PATH \"%PATH%;$(INSTALL_PATH)\""
	@echo "  Git Bash: export PATH=$$PATH:$(INSTALL_PATH)"
else
	@$(CP) $(BINARY_PATH_OS) $(INSTALL_PATH)/
	@echo "✓ Installation complete!"
endif
	@echo ""
	@echo "Verify installation with: $(BINARY_NAME) version"

# Uninstall binary from system
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_PATH)..."
	@$(RM) $(INSTALL_PATH)/$(BINARY_NAME)$(BINARY_EXT)
	@echo "Uninstall complete!"

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
	@echo "  install      - Show installation instructions"
	@echo "  install-go   - Quick install using go install"
	@echo "  install-local- Install by building locally first"
	@echo "  uninstall    - Remove binary from GOPATH/bin"
	@echo "  help         - Show this help"