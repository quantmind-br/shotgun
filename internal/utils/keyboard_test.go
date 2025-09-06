package utils

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestGetKeyboardCapabilities(t *testing.T) {
	caps := GetKeyboardCapabilities()

	// Should have valid platform
	if caps.Platform == "" {
		t.Error("Platform should not be empty")
	}

	// Should have valid terminal
	if caps.Terminal == "" {
		t.Error("Terminal should not be empty")
	}

	// Should have recommended keys
	if len(caps.RecommendedKeys) == 0 {
		t.Error("Should have recommended keys")
	}

	// Should have known issues (even if empty)
	_ = caps.KnownIssues // Just verify it's accessible

	// Should have tested terminals
	if len(caps.TestedTerminals) == 0 {
		t.Error("Should have tested terminals info")
	}
}

func TestIsKeySupported(t *testing.T) {
	tests := []struct {
		key      string
		expected bool // We can't predict exact support, but these should not panic
	}{
		{"ESC", true},
		{"Escape", true},
		{"Enter", true},
		{"Tab", true},
		{"Up", true},
		{"Down", true},
		{"Left", true},
		{"Right", true},
		{"Ctrl+C", true},
	}

	for _, test := range tests {
		result := IsKeySupported(test.key)
		// We mainly test that it doesn't panic and returns a boolean
		_ = result
	}
}

func TestGetRecommendedKeys(t *testing.T) {
	keys := GetRecommendedKeys()

	if len(keys) == 0 {
		t.Error("Should return recommended keys")
	}

	// Should contain ESC as it's universally supported
	found := false
	for _, key := range keys {
		if strings.Contains(strings.ToUpper(key), "ESC") {
			found = true
			break
		}
	}
	if !found {
		t.Log("Warning: ESC not found in recommended keys:", keys)
	}
}

func TestGetKeyMappingInfo(t *testing.T) {
	info := GetKeyMappingInfo()

	if info == "" {
		t.Error("Key mapping info should not be empty")
	}

	// Should contain platform information
	if !strings.Contains(info, "Platform:") {
		t.Error("Info should contain platform information")
	}

	// Should contain support information
	if !strings.Contains(info, "Support:") {
		t.Error("Info should contain support information")
	}
}

func TestTestKeyboardInput(t *testing.T) {
	// Create test key messages
	tests := []tea.KeyMsg{
		{Type: tea.KeyEscape},
		{Type: tea.KeyEnter},
		{Type: tea.KeyTab},
		{Type: tea.KeyF1},
		{Type: tea.KeyF10},
	}

	for _, keyMsg := range tests {
		result := TestKeyboardInput(keyMsg)

		if result == "" {
			t.Errorf("TestKeyboardInput should return non-empty result for key %v", keyMsg)
		}

		// Should contain key information
		if !strings.Contains(result, "Key:") {
			t.Errorf("Result should contain key information: %s", result)
		}

		// Should contain supported information
		if !strings.Contains(result, "Supported:") {
			t.Errorf("Result should contain supported information: %s", result)
		}
	}
}

func TestKeyboardTestMatrix(t *testing.T) {
	matrices := KeyboardTestMatrix()

	if len(matrices) == 0 {
		t.Error("Should return test matrices")
	}

	// Should have entries for all major platforms
	platforms := make(map[string]bool)
	terminals := make(map[string]bool)

	for _, matrix := range matrices {
		platforms[matrix.Platform] = true
		terminals[matrix.Terminal] = true

		// Each matrix should have valid data
		if matrix.Platform == "" {
			t.Error("Matrix entry should have platform")
		}

		if matrix.Terminal == "" {
			t.Error("Matrix entry should have terminal")
		}

		if len(matrix.RecommendedKeys) == 0 {
			t.Errorf("Matrix entry for %s/%s should have recommended keys",
				matrix.Platform, matrix.Terminal)
		}
	}

	// Should have all major platforms
	expectedPlatforms := []string{"windows", "darwin", "linux"}
	for _, platform := range expectedPlatforms {
		if !platforms[platform] {
			t.Errorf("Missing platform in test matrix: %s", platform)
		}
	}

	// Should have key terminals
	expectedTerminals := []string{"cmd", "powershell_7", "terminal", "iterm2", "gnome_terminal"}
	for _, terminal := range expectedTerminals {
		if !terminals[terminal] {
			t.Errorf("Missing terminal in test matrix: %s", terminal)
		}
	}
}

func TestConfigureWindowsKeyboard(t *testing.T) {
	baseCaps := KeyboardCapabilities{
		Platform:     "windows",
		FKeyMappings: make(map[string]string),
	}

	terminals := []string{"cmd", "powershell_5", "powershell_7", "windows_terminal"}

	for _, terminal := range terminals {
		caps := configureWindowsKeyboard(baseCaps, terminal)

		// Should have known issues documented
		if len(caps.KnownIssues) == 0 {
			t.Errorf("Windows terminal %s should have documented known issues", terminal)
		}

		// Should have tested terminals
		if len(caps.TestedTerminals) == 0 {
			t.Errorf("Windows terminal %s should have tested terminals info", terminal)
		}

		// Should have recommended keys
		if len(caps.RecommendedKeys) == 0 {
			t.Errorf("Windows terminal %s should have recommended keys", terminal)
		}

		// Check specific terminal behaviors
		switch terminal {
		case "cmd", "powershell_5":
			if caps.SupportsF1F10 {
				t.Errorf("Terminal %s should not support F1-F10", terminal)
			}
		case "powershell_7", "windows_terminal":
			if !caps.SupportsF1F10 {
				t.Errorf("Terminal %s should support F1-F10", terminal)
			}
			if len(caps.FKeyMappings) == 0 {
				t.Errorf("Terminal %s should have F-key mappings", terminal)
			}
		}
	}
}

func TestConfigureMacOSKeyboard(t *testing.T) {
	baseCaps := KeyboardCapabilities{
		Platform:     "darwin",
		FKeyMappings: make(map[string]string),
	}

	terminals := []string{"terminal", "iterm2"}

	for _, terminal := range terminals {
		caps := configureMacOSKeyboard(baseCaps, terminal)

		// macOS terminals generally support F-keys
		if !caps.SupportsF1F10 {
			t.Errorf("macOS terminal %s should support F1-F10", terminal)
		}

		// Should have F-key mappings
		if len(caps.FKeyMappings) == 0 {
			t.Errorf("macOS terminal %s should have F-key mappings", terminal)
		}

		// Should have recommended keys
		if len(caps.RecommendedKeys) == 0 {
			t.Errorf("macOS terminal %s should have recommended keys", terminal)
		}
	}
}

func TestConfigureLinuxKeyboard(t *testing.T) {
	baseCaps := KeyboardCapabilities{
		Platform:     "linux",
		FKeyMappings: make(map[string]string),
	}

	terminals := []string{"gnome_terminal", "konsole", "xterm"}

	for _, terminal := range terminals {
		caps := configureLinuxKeyboard(baseCaps, terminal)

		// Linux terminals generally support F-keys
		if !caps.SupportsF1F10 {
			t.Errorf("Linux terminal %s should support F1-F10", terminal)
		}

		// Should have F-key mappings
		if len(caps.FKeyMappings) == 0 {
			t.Errorf("Linux terminal %s should have F-key mappings", terminal)
		}

		// Should have recommended keys
		if len(caps.RecommendedKeys) == 0 {
			t.Errorf("Linux terminal %s should have recommended keys", terminal)
		}
	}
}

func TestConfigureUnknownKeyboard(t *testing.T) {
	caps := configureUnknownKeyboard(KeyboardCapabilities{
		Platform:     "unknown",
		FKeyMappings: make(map[string]string),
	})

	// Unknown platforms should be conservative
	if caps.SupportsF1F10 {
		t.Error("Unknown platform should not claim F1-F10 support")
	}

	// Should still have basic keys recommended
	if len(caps.RecommendedKeys) == 0 {
		t.Error("Unknown platform should still have basic recommended keys")
	}

	// Should document uncertainty
	if len(caps.KnownIssues) == 0 {
		t.Error("Unknown platform should document uncertainty in known issues")
	}
}

func TestFKeyMappings(t *testing.T) {
	// Test that F-key mappings are valid ANSI escape sequences
	matrices := KeyboardTestMatrix()

	for _, matrix := range matrices {
		for key, mapping := range matrix.FKeyMappings {
			// Should be a valid F-key
			if !strings.HasPrefix(key, "F") {
				t.Errorf("Invalid F-key name: %s", key)
			}

			// Should be a valid escape sequence
			if !strings.HasPrefix(mapping, "\\x1b") {
				t.Errorf("Invalid escape sequence for %s: %s", key, mapping)
			}
		}
	}
}

// Benchmark tests
func BenchmarkGetKeyboardCapabilities(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKeyboardCapabilities()
	}
}

func BenchmarkIsKeySupported(b *testing.B) {
	keys := []string{"ESC", "F1", "F5", "Enter", "Tab"}
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			IsKeySupported(key)
		}
	}
}

func BenchmarkKeyboardTestMatrix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KeyboardTestMatrix()
	}
}
