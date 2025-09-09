# PowerShell installation script for shotgun
Write-Host "Installing shotgun..." -ForegroundColor Green

try {
    # Use full path to go.exe
    $goExe = "C:\Program Files\Go\bin\go.exe"
    
    if (-not (Test-Path $goExe)) {
        throw "Go not found at $goExe. Please check Go installation."
    }
    
    Write-Host "Using Go from: $goExe" -ForegroundColor Cyan
    
    # Set environment variables directly from registry/known locations
    $env:GOPATH = "C:\Users\Diogo\go"
    $env:GOMODCACHE = "C:\Users\Diogo\go\pkg\mod"
    $env:GOCACHE = "C:\Users\Diogo\AppData\Local\go-build"
    $env:GOTMPDIR = "C:\Users\Diogo\AppData\Local\Temp"
    $env:TMP = "C:\Users\Diogo\AppData\Local\Temp"
    $env:TEMP = "C:\Users\Diogo\AppData\Local\Temp"
    
    Write-Host "GOPATH: $env:GOPATH" -ForegroundColor Cyan
    Write-Host "GOMODCACHE: $env:GOMODCACHE" -ForegroundColor Cyan
    
    # Run go install with full path and explicit environment
    Write-Host "Running go install..." -ForegroundColor Yellow
    
    $result = Start-Process -FilePath $goExe -ArgumentList @("install", "-ldflags=`"-s -w`"", "-buildvcs=false", "./cmd/shotgun") -Wait -PassThru -NoNewWindow
    
    if ($result.ExitCode -eq 0) {
        Write-Host "Installation complete!" -ForegroundColor Green
        Write-Host "Installed to: $env:GOPATH\bin\shotgun.exe" -ForegroundColor Green
        Write-Host "Verify with: shotgun version" -ForegroundColor Cyan
    } else {
        throw "go install failed with exit code $($result.ExitCode)"
    }
}
catch {
    Write-Host "Installation failed: $_" -ForegroundColor Red
    exit 1
}