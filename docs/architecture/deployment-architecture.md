# Deployment Architecture

## Deployment Strategy

**Frontend Deployment:**
- **Platform:** User's local machine
- **Build Command:** `go build -o shotgun ./cmd/shotgun`
- **Output Directory:** `./dist/`
- **CDN/Edge:** N/A - Local binary

**Backend Deployment:**
- **Platform:** Same binary as frontend
- **Build Command:** Same as frontend
- **Deployment Method:** Binary distribution via GitHub Releases

## CI/CD Pipeline
```yaml
# .github/workflows/release.yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, windows, darwin]
        goarch: [amd64, arm64]
        exclude:
          - goos: windows
            goarch: arm64
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          go build -ldflags="-s -w" -o dist/shotgun-${{ matrix.goos }}-${{ matrix.goarch }} ./cmd/shotgun
      
      - name: Upload Release Assets
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
```

## Environments

| Environment | Frontend URL | Backend URL | Purpose |
|-------------|--------------|-------------|---------|
| Development | localhost (terminal) | localhost (same process) | Local development |
| Staging | N/A | N/A | Not applicable for CLI |
| Production | User's terminal | User's machine | Production use |
