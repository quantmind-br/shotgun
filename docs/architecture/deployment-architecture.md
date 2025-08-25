# Deployment Architecture

### Deployment Strategy

**Frontend Deployment:**
- **Platform:** GitHub Releases with cross-compiled binaries
- **Build Command:** `make build-all` (builds for Windows, Linux, macOS)
- **Output Directory:** `./dist/` with platform-specific binaries
- **CDN/Edge:** GitHub's global CDN for release asset distribution

**Backend Deployment:**
- **Platform:** Single binary deployment (no separate backend)
- **Build Command:** `go build -ldflags="-s -w" -o shotgun cmd/shotgun/main.go`
- **Deployment Method:** Direct binary execution, no server deployment required

### CI/CD Pipeline

```yaml
name: Build and Release

on:
  push:
    branches: [main]
    tags: ['v*']
  pull_request:
    branches: [main]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
          
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html
          
      - name: Run linting
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          
      - name: Upload coverage reports
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out

  build:
    name: Build Cross-Platform Binaries
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/')
    
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
        exclude:
          - goos: windows
            goarch: arm64
            
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
          
      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: 0
        run: |
          # Set binary extension for Windows
          EXT=""
          if [ "$GOOS" = "windows" ]; then EXT=".exe"; fi
          
          # Build with version information
          VERSION=${GITHUB_REF#refs/tags/}
          LDFLAGS="-s -w -X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
          
          # Create output directory
          mkdir -p dist
          
          # Build binary
          go build -ldflags="$LDFLAGS" -o "dist/shotgun-$GOOS-$GOARCH$EXT" cmd/shotgun/main.go
          
          # Create archive
          if [ "$GOOS" = "windows" ]; then
            zip -j "dist/shotgun-$GOOS-$GOARCH.zip" "dist/shotgun-$GOOS-$GOARCH$EXT"
          else
            tar -czf "dist/shotgun-$GOOS-$GOARCH.tar.gz" -C dist "shotgun-$GOOS-$GOARCH"
          fi
          
      - name: Upload build artifacts
        uses: actions/upload-artifact@v3
        with:
          name: shotgun-${{ matrix.goos }}-${{ matrix.goarch }}
          path: dist/shotgun-*

  release:
    name: Create Release
    needs: build
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/')
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Download all artifacts
        uses: actions/download-artifact@v3
        with:
          path: ./dist
          
      - name: Generate release notes
        run: |
          # Extract version from tag
          VERSION=${GITHUB_REF#refs/tags/}
          
          # Generate changelog from commits since last tag
          LAST_TAG=$(git describe --tags --abbrev=0 HEAD~1 2>/dev/null || echo "")
          if [ -n "$LAST_TAG" ]; then
            echo "## Changes since $LAST_TAG" > CHANGELOG.md
            git log --pretty=format:"- %s (%an)" $LAST_TAG..HEAD >> CHANGELOG.md
          else
            echo "## Initial Release" > CHANGELOG.md
          fi
          
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            dist/**/*.tar.gz
            dist/**/*.zip
          body_path: CHANGELOG.md
          draft: false
          prerelease: ${{ contains(github.ref, '-') }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Gosec Security Scanner
        uses: securecodewarrior/github-action-gosec@master
        with:
          args: './...'
          
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
```

### Environments

| Environment | Frontend URL | Backend URL | Purpose |
|-------------|-------------|-------------|---------|
| Development | Local binary execution | N/A - Single binary | Local development and testing |
| Staging | GitHub Actions artifacts | N/A - Single binary | Pre-release testing and validation |
| Production | GitHub Releases | N/A - Single binary | Live distribution to end users |
