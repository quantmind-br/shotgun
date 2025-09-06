package styles

import (
	"strings"
	"testing"

	"github.com/diogopedro/shotgun/internal/utils"
)

func TestNewFallbackConfig(t *testing.T) {
	config := NewFallbackConfig()

	// Should have valid spinner and progress styles
	if config.SpinnerStyle < 0 || config.SpinnerStyle > 1 {
		t.Error("Invalid spinner style")
	}

	if config.ProgressStyle < 0 || config.ProgressStyle > 2 {
		t.Error("Invalid progress style")
	}

	// Capability should be valid
	if config.Capability != utils.UnicodeNone &&
		config.Capability != utils.UnicodeBasic &&
		config.Capability != utils.UnicodeFull {
		t.Error("Invalid Unicode capability")
	}
}

func TestSpinnerCharacters(t *testing.T) {
	tests := []struct {
		style    SpinnerStyle
		expected []string
	}{
		{SpinnerStyleUnicode, []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}},
		{SpinnerStyleASCII, []string{"|", "/", "-", "\\"}},
	}

	for _, test := range tests {
		config := FallbackConfig{SpinnerStyle: test.style}
		result := config.SpinnerCharacters()

		if len(result) != len(test.expected) {
			t.Errorf("SpinnerCharacters(%v) length = %d, expected %d",
				test.style, len(result), len(test.expected))
			continue
		}

		for i, expected := range test.expected {
			if i < len(result) && result[i] != expected {
				t.Errorf("SpinnerCharacters(%v)[%d] = %q, expected %q",
					test.style, i, result[i], expected)
			}
		}
	}
}

func TestProgressCharacters(t *testing.T) {
	tests := []struct {
		style          ProgressStyle
		expectedFilled string
		expectedEmpty  string
	}{
		{ProgressStyleUnicode, "█", "░"},
		{ProgressStyleASCII, "#", "-"},
		{ProgressStyleMinimal, "=", " "},
	}

	for _, test := range tests {
		config := FallbackConfig{ProgressStyle: test.style}
		filled, empty := config.ProgressCharacters()

		if filled != test.expectedFilled {
			t.Errorf("ProgressCharacters(%v) filled = %q, expected %q",
				test.style, filled, test.expectedFilled)
		}

		if empty != test.expectedEmpty {
			t.Errorf("ProgressCharacters(%v) empty = %q, expected %q",
				test.style, empty, test.expectedEmpty)
		}
	}
}

func TestBorderCharacters(t *testing.T) {
	// Test Unicode border characters
	unicodeConfig := FallbackConfig{
		UseUnicode: true,
		Capability: utils.UnicodeFull,
	}

	borders := unicodeConfig.BorderCharacters()
	expectedUnicode := []string{"┌", "┐", "└", "┘", "─", "│"}
	actualUnicode := []string{
		borders.TopLeft, borders.TopRight, borders.BottomLeft,
		borders.BottomRight, borders.Horizontal, borders.Vertical,
	}

	for i, expected := range expectedUnicode {
		if actualUnicode[i] != expected {
			t.Errorf("Unicode border character %d = %q, expected %q",
				i, actualUnicode[i], expected)
		}
	}

	// Test ASCII border characters
	asciiConfig := FallbackConfig{
		UseUnicode: false,
		Capability: utils.UnicodeNone,
	}

	asciiBorders := asciiConfig.BorderCharacters()
	expectedASCII := []string{"+", "+", "+", "+", "-", "|"}
	actualASCII := []string{
		asciiBorders.TopLeft, asciiBorders.TopRight, asciiBorders.BottomLeft,
		asciiBorders.BottomRight, asciiBorders.Horizontal, asciiBorders.Vertical,
	}

	for i, expected := range expectedASCII {
		if actualASCII[i] != expected {
			t.Errorf("ASCII border character %d = %q, expected %q",
				i, actualASCII[i], expected)
		}
	}
}

func TestStatusSymbols(t *testing.T) {
	// Test Unicode status symbols
	unicodeConfig := FallbackConfig{
		UseUnicode: true,
		Capability: utils.UnicodeFull,
	}

	symbols := unicodeConfig.StatusSymbols()
	if symbols.Success != "✓" {
		t.Errorf("Unicode success symbol = %q, expected ✓", symbols.Success)
	}
	if symbols.Error != "✗" {
		t.Errorf("Unicode error symbol = %q, expected ✗", symbols.Error)
	}

	// Test ASCII status symbols
	asciiConfig := FallbackConfig{
		UseUnicode: false,
		Capability: utils.UnicodeNone,
	}

	asciiSymbols := asciiConfig.StatusSymbols()
	if asciiSymbols.Success != "+" {
		t.Errorf("ASCII success symbol = %q, expected +", asciiSymbols.Success)
	}
	if asciiSymbols.Error != "x" {
		t.Errorf("ASCII error symbol = %q, expected x", asciiSymbols.Error)
	}
}

func TestSanitizeText(t *testing.T) {
	config := FallbackConfig{Capability: utils.UnicodeNone}

	input := "Loading... ⠋"
	result := config.SanitizeText(input)

	// Should convert Unicode to ASCII
	if strings.Contains(result, "⠋") {
		t.Error("SanitizeText should convert Unicode characters in UnicodeNone mode")
	}

	if !strings.Contains(result, "Loading...") {
		t.Error("SanitizeText should preserve ASCII text")
	}
}

func TestGetFallbackSpinnerFrames(t *testing.T) {
	frames := GetFallbackSpinnerFrames()

	// Should have both Unicode and ASCII frames
	if _, exists := frames[SpinnerStyleUnicode]; !exists {
		t.Error("Should have Unicode spinner frames")
	}

	if _, exists := frames[SpinnerStyleASCII]; !exists {
		t.Error("Should have ASCII spinner frames")
	}

	// Unicode frames should contain Braille characters
	unicodeFrames := frames[SpinnerStyleUnicode]
	if len(unicodeFrames) == 0 {
		t.Error("Unicode frames should not be empty")
	}

	// ASCII frames should contain basic ASCII characters
	asciiFrames := frames[SpinnerStyleASCII]
	if len(asciiFrames) == 0 {
		t.Error("ASCII frames should not be empty")
	}

	for _, frame := range asciiFrames {
		if len(frame) != 1 {
			t.Errorf("ASCII frame should be single character, got %q", frame)
		}
		if frame[0] > 127 {
			t.Errorf("ASCII frame should be ASCII character, got %q", frame)
		}
	}
}

func TestGetFallbackProgressChars(t *testing.T) {
	chars := GetFallbackProgressChars()

	// Should have all progress styles
	expectedStyles := []ProgressStyle{
		ProgressStyleUnicode,
		ProgressStyleASCII,
		ProgressStyleMinimal,
	}

	for _, style := range expectedStyles {
		if _, exists := chars[style]; !exists {
			t.Errorf("Should have progress chars for style %v", style)
		}
	}

	// Check specific characters
	if chars[ProgressStyleUnicode][0] != "█" {
		t.Error("Unicode progress filled char should be █")
	}

	if chars[ProgressStyleASCII][0] != "#" {
		t.Error("ASCII progress filled char should be #")
	}
}

func TestTestTerminalFeatures(t *testing.T) {
	features := TestTerminalFeatures()

	// Should have valid platform
	if features.Platform == "" {
		t.Error("Platform should not be empty")
	}

	// Should have valid terminal
	if features.Terminal == "" {
		t.Error("Terminal should not be empty")
	}

	// Color depth should be positive
	if features.ColorDepth < 1 {
		t.Error("Color depth should be at least 1")
	}

	// Test string representation
	str := features.String()
	if !strings.Contains(str, "Platform:") {
		t.Error("String representation should contain Platform")
	}

	if !strings.Contains(str, "Terminal:") {
		t.Error("String representation should contain Terminal")
	}
}

func TestTerminalFeaturesString(t *testing.T) {
	tests := []struct {
		features TerminalFeatures
		contains []string
	}{
		{
			TerminalFeatures{
				Platform:            "linux",
				Terminal:            "gnome_terminal",
				SupportsUnicode:     true,
				SupportsFullUnicode: true,
				ColorDepth:          256,
				FKeys:               true,
				ESCSupport:          true,
			},
			[]string{"linux", "gnome_terminal", "Full", "256-color", "Yes"},
		},
		{
			TerminalFeatures{
				Platform:            "windows",
				Terminal:            "cmd",
				SupportsUnicode:     false,
				SupportsFullUnicode: false,
				ColorDepth:          16,
				FKeys:               false,
				ESCSupport:          true,
			},
			[]string{"windows", "cmd", "None", "16-color", "No"},
		},
	}

	for _, test := range tests {
		result := test.features.String()
		for _, expected := range test.contains {
			if !strings.Contains(result, expected) {
				t.Errorf("TerminalFeatures.String() should contain %q, got %q", expected, result)
			}
		}
	}
}

// Benchmark tests
func BenchmarkNewFallbackConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFallbackConfig()
	}
}

func BenchmarkSpinnerCharacters(b *testing.B) {
	config := FallbackConfig{SpinnerStyle: SpinnerStyleUnicode}
	for i := 0; i < b.N; i++ {
		config.SpinnerCharacters()
	}
}

func BenchmarkProgressCharacters(b *testing.B) {
	config := FallbackConfig{ProgressStyle: ProgressStyleUnicode}
	for i := 0; i < b.N; i++ {
		config.ProgressCharacters()
	}
}

func BenchmarkSanitizeText(b *testing.B) {
	config := FallbackConfig{Capability: utils.UnicodeNone}
	text := "Loading... ⠋ Progress: █░"
	for i := 0; i < b.N; i++ {
		config.SanitizeText(text)
	}
}
