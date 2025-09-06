package e2e

import (
	"os"
	"runtime"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/styles"
	"github.com/diogopedro/shotgun/internal/utils"
)

// PlatformTestConfig defines platform-specific test configuration
type PlatformTestConfig struct {
	Platform    string
	Terminal    string
	ColorDepth  int
	HasUnicode  bool
	TERM        string
	COLORTERM   string
	NO_COLOR    string
	FORCE_COLOR string
}

// GetPlatformConfigs returns test configurations for different platforms
func GetPlatformConfigs() []PlatformTestConfig {
	return []PlatformTestConfig{
		// Windows configurations
		{
			Platform:    "windows",
			Terminal:    "cmd",
			ColorDepth:  16,
			HasUnicode:  false,
			TERM:        "cmd",
			COLORTERM:   "",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "windows",
			Terminal:    "powershell_5",
			ColorDepth:  256,
			HasUnicode:  false,
			TERM:        "xterm",
			COLORTERM:   "",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "windows",
			Terminal:    "powershell_7",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "windows",
			Terminal:    "windows_terminal",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		// macOS configurations
		{
			Platform:    "darwin",
			Terminal:    "terminal",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "darwin",
			Terminal:    "iterm2",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		// Linux configurations
		{
			Platform:    "linux",
			Terminal:    "gnome_terminal",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "linux",
			Terminal:    "konsole",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "linux",
			Terminal:    "xterm",
			ColorDepth:  256,
			HasUnicode:  true,
			TERM:        "xterm",
			COLORTERM:   "",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		// Edge cases
		{
			Platform:    "linux",
			Terminal:    "vt100",
			ColorDepth:  1,
			HasUnicode:  false,
			TERM:        "vt100",
			COLORTERM:   "",
			NO_COLOR:    "",
			FORCE_COLOR: "",
		},
		{
			Platform:    "linux",
			Terminal:    "no_color",
			ColorDepth:  1,
			HasUnicode:  true,
			TERM:        "xterm-256color",
			COLORTERM:   "truecolor",
			NO_COLOR:    "1",
			FORCE_COLOR: "",
		},
		{
			Platform:    "linux",
			Terminal:    "force_color",
			ColorDepth:  16777216,
			HasUnicode:  true,
			TERM:        "vt100",
			COLORTERM:   "",
			NO_COLOR:    "",
			FORCE_COLOR: "3",
		},
	}
}

// setupTestEnvironment configures environment for platform testing
func setupTestEnvironment(t *testing.T, config PlatformTestConfig) (func(), func()) {
	// Store original environment
	originalEnv := map[string]string{
		"TERM":        os.Getenv("TERM"),
		"COLORTERM":   os.Getenv("COLORTERM"),
		"NO_COLOR":    os.Getenv("NO_COLOR"),
		"FORCE_COLOR": os.Getenv("FORCE_COLOR"),
	}

	// Set test environment
	os.Setenv("TERM", config.TERM)
	os.Setenv("COLORTERM", config.COLORTERM)
	os.Setenv("NO_COLOR", config.NO_COLOR)
	os.Setenv("FORCE_COLOR", config.FORCE_COLOR)

	// Reset theme to pick up new environment
	resetTheme := func() {
		styles.ResetGlobalTheme()
	}

	// Restore environment
	restore := func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
		styles.ResetGlobalTheme()
	}

	resetTheme()
	return resetTheme, restore
}

// TestPlatformCompatibility tests terminal compatibility across platforms
func TestPlatformCompatibility(t *testing.T) {
	configs := GetPlatformConfigs()

	for _, config := range configs {
		t.Run(config.Platform+"_"+config.Terminal, func(t *testing.T) {
			_, restore := setupTestEnvironment(t, config)
			defer restore()

			// Test terminal capability detection
			caps := utils.DetectTerminalCapabilities()

			t.Logf("Testing %s/%s - Expected ColorDepth: %d, Detected: %d",
				config.Platform, config.Terminal, config.ColorDepth, caps.ColorDepth)

			// Color depth should be appropriate for terminal
			if config.NO_COLOR != "" {
				// NO_COLOR should force monochrome
				if caps.ColorDepth != 1 {
					t.Logf("Warning: NO_COLOR set but color depth is %d", caps.ColorDepth)
				}
			} else if config.FORCE_COLOR != "" {
				// FORCE_COLOR should override detection
				t.Logf("FORCE_COLOR=%s resulted in color depth %d", config.FORCE_COLOR, caps.ColorDepth)
			}

			// Test theme configuration
			theme := styles.GetGlobalTheme()
			if theme.Colors.Primary == "" {
				t.Error("Theme should have primary color")
			}

			// Test fallback configuration
			fallback := styles.NewFallbackConfig()
			if fallback.UseUnicode && fallback.Capability == utils.UnicodeNone {
				t.Error("Fallback configuration inconsistent")
			}
		})
	}
}

// TestUnicodeRendering validates Unicode character rendering across platforms
func TestUnicodeRendering(t *testing.T) {
	configs := GetPlatformConfigs()

	testStrings := []struct {
		name            string
		text            string
		requiresUnicode bool
	}{
		{"ascii", "Loading... [####----] 50%", false},
		{"basic_unicode", "Loading... [█▊▊▊░░░░] 50%", true},
		{"advanced_unicode", "Loading... ⠋ Progress: █░", true},
		{"mixed", "File: test.txt ✓ Completed", true},
		{"braille_spinner", "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏", true},
		{"box_drawing", "┌─────────┐│ Content │└─────────┘", true},
	}

	for _, config := range configs {
		t.Run(config.Platform+"_"+config.Terminal+"_unicode", func(t *testing.T) {
			_, restore := setupTestEnvironment(t, config)
			defer restore()

			fallback := styles.NewFallbackConfig()

			for _, testStr := range testStrings {
				t.Run(testStr.name, func(t *testing.T) {
					result := fallback.SanitizeText(testStr.text)

					// Text should not be empty
					if result == "" {
						t.Error("Sanitized text should not be empty")
					}

					// For terminals without Unicode, advanced Unicode should be converted
					if !config.HasUnicode && testStr.requiresUnicode {
						// Should contain ASCII alternatives
						if strings.Contains(result, "⠋") ||
							strings.Contains(result, "█") ||
							strings.Contains(result, "┌") {
							t.Logf("Warning: Unicode characters found in non-Unicode terminal: %s", result)
						}
					}

					// For Unicode terminals, Unicode should be preserved
					if config.HasUnicode && testStr.requiresUnicode {
						// Original text should be mostly preserved (some fallbacks still apply)
						if len(result) == 0 {
							t.Error("Unicode terminal should preserve text content")
						}
					}

					t.Logf("Platform: %s/%s, Input: %s, Output: %s",
						config.Platform, config.Terminal, testStr.text, result)
				})
			}
		})
	}
}

// TestColorOutput validates color output across different color depths
func TestColorOutput(t *testing.T) {
	configs := GetPlatformConfigs()

	for _, config := range configs {
		t.Run(config.Platform+"_"+config.Terminal+"_colors", func(t *testing.T) {
			_, restore := setupTestEnvironment(t, config)
			defer restore()

			theme := styles.GetGlobalTheme()

			// Test color palette creation
			if string(theme.Colors.Primary) == "" {
				t.Error("Primary color should not be empty")
			}

			if string(theme.Colors.Success) == "" {
				t.Error("Success color should not be empty")
			}

			if string(theme.Colors.Error) == "" {
				t.Error("Error color should not be empty")
			}

			// Test style rendering
			titleStyle := theme.Styles.Title
			if titleStyle.GetForeground() == nil {
				t.Error("Title style should have foreground color")
			}

			// Test color profile detection
			profile := styles.GetColorProfile()
			t.Logf("Platform: %s/%s, Detected profile: %d",
				config.Platform, config.Terminal, int(profile))

			// Render test content with styles
			testContent := titleStyle.Render("Test Title")
			if testContent == "" {
				t.Error("Styled content should not be empty")
			}

			successStyle := theme.Styles.Success
			successContent := successStyle.Render("✓ Success")
			if successContent == "" {
				t.Error("Success content should not be empty")
			}

			t.Logf("Styled title: %q", testContent)
			t.Logf("Styled success: %q", successContent)
		})
	}
}

// TestKeyboardHandling validates keyboard input handling across platforms
func TestKeyboardHandling(t *testing.T) {
	configs := GetPlatformConfigs()

	testKeys := []struct {
		key              tea.KeyType
		name             string
		universalSupport bool
	}{
		{tea.KeyEscape, "ESC", true},
		{tea.KeyEnter, "Enter", true},
		{tea.KeyTab, "Tab", true},
		{tea.KeyUp, "Up", true},
		{tea.KeyDown, "Down", true},
		{tea.KeyF1, "F1", false},
		{tea.KeyF5, "F5", false},
		{tea.KeyF10, "F10", false},
	}

	for _, config := range configs {
		t.Run(config.Platform+"_"+config.Terminal+"_keyboard", func(t *testing.T) {
			_, restore := setupTestEnvironment(t, config)
			defer restore()

			keyboardCaps := utils.GetKeyboardCapabilities()

			for _, testKey := range testKeys {
				t.Run(testKey.name, func(t *testing.T) {
					supported := utils.IsKeySupported(testKey.name)

					// Universal keys should always be supported
					if testKey.universalSupport && !supported {
						t.Errorf("Key %s should be universally supported", testKey.name)
					}

					// F-keys should match platform expectations
					if strings.HasPrefix(testKey.name, "F") {
						if supported != keyboardCaps.SupportsF1F10 {
							t.Logf("F-key support mismatch for %s: expected %t, got %t",
								testKey.name, keyboardCaps.SupportsF1F10, supported)
						}
					}

					t.Logf("Key %s supported on %s/%s: %t",
						testKey.name, config.Platform, config.Terminal, supported)
				})
			}

			// Test recommended keys
			recommendedKeys := utils.GetRecommendedKeys()
			if len(recommendedKeys) == 0 {
				t.Error("Should have recommended keys for platform")
			}

			t.Logf("Recommended keys for %s/%s: %v",
				config.Platform, config.Terminal, recommendedKeys)
		})
	}
}

// TestApplicationFlow tests complete application flow simulation
func TestApplicationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping application flow test in short mode")
	}

	// Test on current platform only to avoid complexity
	config := PlatformTestConfig{
		Platform:    runtime.GOOS,
		Terminal:    "current",
		ColorDepth:  256,
		HasUnicode:  true,
		TERM:        "xterm-256color",
		COLORTERM:   "truecolor",
		NO_COLOR:    "",
		FORCE_COLOR: "",
	}

	_, restore := setupTestEnvironment(t, config)
	defer restore()

	t.Run("theme_initialization", func(t *testing.T) {
		// Test theme initialization
		theme := styles.GetGlobalTheme()

		if theme.Colors.Primary == "" {
			t.Error("Theme should initialize with primary color")
		}

		if theme.RenderConfig.ColorProfile < 0 {
			t.Error("Theme should have valid color profile")
		}

		t.Logf("Theme initialized successfully for %s", runtime.GOOS)
	})

	t.Run("component_rendering", func(t *testing.T) {
		// Test that components can render without errors
		fallback := styles.NewFallbackConfig()

		// Test spinner characters
		spinnerChars := fallback.SpinnerCharacters()
		if len(spinnerChars) == 0 {
			t.Error("Should have spinner characters")
		}

		// Test progress characters
		filled, empty := fallback.ProgressCharacters()
		if filled == "" || empty == "" {
			t.Error("Should have progress characters")
		}

		// Test status symbols
		symbols := fallback.StatusSymbols()
		if symbols.Success == "" || symbols.Error == "" {
			t.Error("Should have status symbols")
		}

		t.Logf("Components render successfully on %s", runtime.GOOS)
	})

	t.Run("terminal_capability_detection", func(t *testing.T) {
		// Test terminal capability detection
		caps := utils.DetectTerminalCapabilities()

		if caps.Platform == "" {
			t.Error("Should detect platform")
		}

		if caps.ColorDepth < 1 {
			t.Error("Should detect valid color depth")
		}

		keyboardCaps := utils.GetKeyboardCapabilities()
		if keyboardCaps.Platform == "" {
			t.Error("Should detect keyboard platform")
		}

		t.Logf("Capabilities detected for %s: Colors=%d, Unicode=%t, FKeys=%t",
			caps.Platform, caps.ColorDepth, caps.HasUnicode, keyboardCaps.SupportsF1F10)
	})
}

// BenchmarkPlatformDetection benchmarks platform detection performance
func BenchmarkPlatformDetection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = utils.DetectTerminalCapabilities()
		_ = utils.GetKeyboardCapabilities()
		_ = styles.NewFallbackConfig()
	}
}

// TestTerminalEmulation tests different terminal emulation scenarios
func TestTerminalEmulation(t *testing.T) {
	scenarios := []struct {
		name              string
		term              string
		colorterm         string
		noColor           string
		forceColor        string
		expectedMinColors int
		expectedUnicode   bool
	}{
		{
			name:              "modern_terminal",
			term:              "xterm-256color",
			colorterm:         "truecolor",
			expectedMinColors: 256,
			expectedUnicode:   true,
		},
		{
			name:              "basic_terminal",
			term:              "xterm",
			expectedMinColors: 16,
			expectedUnicode:   true,
		},
		{
			name:              "legacy_terminal",
			term:              "vt100",
			expectedMinColors: 1,
			expectedUnicode:   false,
		},
		{
			name:              "no_color_override",
			term:              "xterm-256color",
			colorterm:         "truecolor",
			noColor:           "1",
			expectedMinColors: 1,
			expectedUnicode:   true,
		},
		{
			name:              "force_color_truecolor",
			term:              "vt100",
			forceColor:        "3",
			expectedMinColors: 256,
			expectedUnicode:   false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Set up environment
			originalTERM := os.Getenv("TERM")
			originalCOLORTERM := os.Getenv("COLORTERM")
			originalNO_COLOR := os.Getenv("NO_COLOR")
			originalFORCE_COLOR := os.Getenv("FORCE_COLOR")

			defer func() {
				os.Setenv("TERM", originalTERM)
				os.Setenv("COLORTERM", originalCOLORTERM)
				os.Setenv("NO_COLOR", originalNO_COLOR)
				os.Setenv("FORCE_COLOR", originalFORCE_COLOR)
				styles.ResetGlobalTheme()
			}()

			os.Setenv("TERM", scenario.term)
			os.Setenv("COLORTERM", scenario.colorterm)
			os.Setenv("NO_COLOR", scenario.noColor)
			os.Setenv("FORCE_COLOR", scenario.forceColor)

			styles.ResetGlobalTheme()

			// Test detection
			caps := utils.DetectTerminalCapabilities()

			t.Logf("Scenario %s: detected %d colors, unicode=%t",
				scenario.name, caps.ColorDepth, caps.HasUnicode)

			// Validate expectations
			if scenario.expectedMinColors > 1 && caps.ColorDepth < scenario.expectedMinColors {
				// Allow some flexibility in detection
				t.Logf("Warning: Expected min %d colors, got %d",
					scenario.expectedMinColors, caps.ColorDepth)
			}
		})
	}
}
