package utils

import (
	"os"
	"runtime"
	"testing"
)

func TestGetPlatformInfo(t *testing.T) {
	info := GetPlatformInfo()

	if info.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, info.OS)
	}

	if info.Terminal == "" {
		t.Error("Expected non-empty terminal detection")
	}

	// Basic capability checks
	if info.Capabilities.Platform != runtime.GOOS {
		t.Errorf("Expected platform %s, got %s", runtime.GOOS, info.Capabilities.Platform)
	}

	if info.Capabilities.ColorDepth < 1 {
		t.Error("Expected color depth >= 1")
	}
}

func TestDetectTerminalCapabilities(t *testing.T) {
	caps := DetectTerminalCapabilities()

	// Test basic properties
	if caps.Platform != runtime.GOOS {
		t.Errorf("Expected platform %s, got %s", runtime.GOOS, caps.Platform)
	}

	// Color depth should be one of the valid values
	validColorDepths := []int{1, 8, 16, 256, 16777216}
	validDepth := false
	for _, depth := range validColorDepths {
		if caps.ColorDepth == depth {
			validDepth = true
			break
		}
	}
	if !validDepth {
		t.Errorf("Invalid color depth: %d", caps.ColorDepth)
	}

	// Unicode support should be boolean
	_ = caps.HasUnicode // Just verify it's accessible

	// F-key support should be boolean
	_ = caps.HasF1F10Keys // Just verify it's accessible

	// ESC support should be boolean
	_ = caps.ESCSupport // Just verify it's accessible
}

func TestIsWindowsCMD(t *testing.T) {
	originalComspec := os.Getenv("COMSPEC")
	originalPSModule := os.Getenv("PSModulePath")

	defer func() {
		os.Setenv("COMSPEC", originalComspec)
		os.Setenv("PSModulePath", originalPSModule)
	}()

	// Test non-Windows platform
	if runtime.GOOS != "windows" {
		if IsWindowsCMD() {
			t.Error("IsWindowsCMD should return false on non-Windows platforms")
		}
		return
	}

	// Test Windows CMD detection
	os.Setenv("COMSPEC", "C:\\Windows\\System32\\cmd.exe")
	os.Setenv("PSModulePath", "")

	if !IsWindowsCMD() {
		t.Error("Expected IsWindowsCMD to return true for CMD environment")
	}

	// Test PowerShell detection (should not be CMD)
	os.Setenv("PSModulePath", "C:\\Program Files\\WindowsPowerShell\\Modules")

	if IsWindowsCMD() {
		t.Error("Expected IsWindowsCMD to return false for PowerShell environment")
	}
}

func TestHasUnicodeSupport(t *testing.T) {
	// Test that function doesn't panic and returns a boolean
	result := HasUnicodeSupport()
	_ = result // Just verify it returns without error
}

func TestGetColorDepth(t *testing.T) {
	depth := GetColorDepth()

	if depth < 1 {
		t.Error("Color depth should be at least 1")
	}

	// Should be one of the standard depths
	validDepths := []int{1, 8, 16, 256, 16777216}
	valid := false
	for _, validDepth := range validDepths {
		if depth == validDepth {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("Invalid color depth: %d", depth)
	}
}

func TestDetectUnicodeSupport(t *testing.T) {
	originalLang := os.Getenv("LANG")
	defer os.Setenv("LANG", originalLang)

	// Test UTF-8 environment detection
	os.Setenv("LANG", "en_US.UTF-8")
	if !detectUnicodeSupport() {
		t.Error("Expected Unicode support with UTF-8 LANG environment")
	}

	// Test non-UTF-8 environment
	os.Setenv("LANG", "en_US.ASCII")
	// Note: This may still return true on modern systems due to other heuristics
	result := detectUnicodeSupport()
	_ = result // Just verify it doesn't panic
}

func TestDetectFKeySupport(t *testing.T) {
	result := detectFKeySupport()
	_ = result // Just verify it returns a boolean without panic
}

func TestGetTerminalMatrix(t *testing.T) {
	matrix := GetTerminalMatrix()

	if len(matrix) == 0 {
		t.Error("Terminal matrix should not be empty")
	}

	// Verify each entry has required fields
	for i, entry := range matrix {
		if entry.Platform == "" {
			t.Errorf("Matrix entry %d missing platform", i)
		}
		if entry.Terminal == "" {
			t.Errorf("Matrix entry %d missing terminal", i)
		}
		if entry.UnicodeLevel == "" {
			t.Errorf("Matrix entry %d missing unicode level", i)
		}
		if entry.ColorSupport < 1 {
			t.Errorf("Matrix entry %d has invalid color support: %d", i, entry.ColorSupport)
		}
	}

	// Verify we have entries for all major platforms
	platforms := make(map[string]bool)
	for _, entry := range matrix {
		platforms[entry.Platform] = true
	}

	expectedPlatforms := []string{"windows", "darwin", "linux"}
	for _, platform := range expectedPlatforms {
		if !platforms[platform] {
			t.Errorf("Missing platform in matrix: %s", platform)
		}
	}
}

func TestAdjustForWindows(t *testing.T) {
	// Only test on Windows or with mocked environment
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	baseCaps := TerminalCapabilities{
		HasUnicode:   true,
		ColorDepth:   16777216,
		HasF1F10Keys: true,
		Platform:     "windows",
		ESCSupport:   true,
	}

	// Test that adjustment doesn't panic
	adjusted := adjustForWindows(baseCaps)

	// Verify it returns valid capabilities
	if adjusted.Platform != "windows" {
		t.Error("Platform should remain windows")
	}

	if adjusted.ColorDepth < 1 {
		t.Error("Color depth should be at least 1")
	}
}

// Benchmark tests
func BenchmarkGetPlatformInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetPlatformInfo()
	}
}

func BenchmarkDetectTerminalCapabilities(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetectTerminalCapabilities()
	}
}

func BenchmarkIsWindowsCMD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsWindowsCMD()
	}
}
