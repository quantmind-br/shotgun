# Essential Development Commands

## Build Commands
- `make build` - Build the application for current platform
- `make build-all` - Build for all platforms (Linux, Windows, macOS)
- `go build -o bin/shotgun ./cmd/shotgun` - Direct Go build

## Testing Commands
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report (generates coverage.html)
- `go test -v ./...` - Run tests with verbose output

## Code Quality Commands
- `make lint` - Run formatting and vetting (combines fmt + vet)
- `make fmt` - Format Go code with gofmt
- `make vet` - Run go vet static analysis
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run static analysis

## Development Commands
- `make run` - Run the application directly
- `go run ./cmd/shotgun` - Run without building
- `make deps` - Install/update dependencies (go mod tidy + download)
- `make clean` - Remove build artifacts

## Windows-Specific Commands
Since this is Windows environment:
- Use `dir` instead of `ls`
- Use `type` instead of `cat`
- Use `findstr` instead of `grep`
- Git commands work normally

## Task Completion Workflow
After completing any development task:
1. `make lint` - Ensure code quality
2. `make test` - Verify all tests pass
3. `make build` - Ensure clean build
4. Test the binary manually if needed