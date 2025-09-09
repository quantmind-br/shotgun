#!/bin/bash
# Installation script for Windows

echo "Installing shotgun..."
export GOPATH="$(go env GOPATH)"
export GOMODCACHE="$(go env GOMODCACHE)" 
export GOCACHE="$(go env GOCACHE)"

go install -ldflags="-s -w" -buildvcs=false ./cmd/shotgun

if [ $? -eq 0 ]; then
    echo "✓ Installation complete!"
    echo "Installed to: $(go env GOPATH)/bin/shotgun.exe"
    echo "Verify with: shotgun version"
else
    echo "❌ Installation failed"
    exit 1
fi