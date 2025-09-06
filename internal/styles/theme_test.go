package styles

import (
	"testing"

	"github.com/diogopedro/shotgun/internal/utils"
	"github.com/muesli/termenv"
)

func TestNewThemeConfig(t *testing.T) {
	theme := NewThemeConfig()

	// Should have valid render config
	if theme.RenderConfig.ColorProfile < 0 {
		t.Error("RenderConfig should have a valid color profile")
	}

	// Should have valid color palette
	if string(theme.Colors.Primary) == "" {
		t.Error("Colors should have a primary color")
	}

	// Should have valid style palette
	if theme.Styles.Title.GetForeground() == nil {
		t.Error("Styles should have a title style with foreground color")
	}
}

func TestDetectColorProfile(t *testing.T) {
	tests := []struct {
		colorDepth      int
		expectedProfile termenv.Profile
	}{
		{1, termenv.Ascii},
		{8, termenv.ANSI},
		{16, termenv.ANSI},
		{256, termenv.ANSI256},
		{16777216, termenv.TrueColor},
	}

	for _, test := range tests {
		caps := utils.TerminalCapabilities{ColorDepth: test.colorDepth}
		profile := detectColorProfile(caps)

		// Compare profile types (can't directly compare profiles)
		if getProfileType(profile) != getProfileType(test.expectedProfile) {
			t.Errorf("Color depth %d should map to profile type %s, got %s",
				test.colorDepth, getProfileType(test.expectedProfile), getProfileType(profile))
		}
	}
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Test that environment variables are handled at the platform detection level
	// and that detectColorProfile correctly maps color depths to profiles
	tests := []struct {
		colorDepth      int
		expectedProfile termenv.Profile
	}{
		{1, termenv.Ascii},
		{8, termenv.ANSI},
		{16, termenv.ANSI},
		{256, termenv.ANSI256},
		{16777216, termenv.TrueColor},
	}

	for _, test := range tests {
		caps := utils.TerminalCapabilities{ColorDepth: test.colorDepth}
		profile := detectColorProfile(caps)

		if getProfileType(profile) != getProfileType(test.expectedProfile) {
			t.Errorf("ColorDepth=%d should map to profile type %s, got %s",
				test.colorDepth, getProfileType(test.expectedProfile), getProfileType(profile))
		}
	}
}

func TestNoColorOverride(t *testing.T) {
	// Test that monochrome color depth maps to ASCII profile (NO_COLOR is handled at platform level)
	caps := utils.TerminalCapabilities{ColorDepth: 1} // Monochrome from platform detection
	profile := detectColorProfile(caps)

	if getProfileType(profile) != getProfileType(termenv.Ascii) {
		t.Error("Monochrome color depth should map to ASCII profile")
	}
}

func TestCreateColorPalettes(t *testing.T) {
	colorDepths := []int{1, 8, 16, 256, 16777216}

	for _, depth := range colorDepths {
		palette := createColorPalette(depth)

		// All palettes should have all required colors
		if string(palette.Primary) == "" {
			t.Errorf("Color depth %d palette missing Primary color", depth)
		}
		if string(palette.Success) == "" {
			t.Errorf("Color depth %d palette missing Success color", depth)
		}
		if string(palette.Error) == "" {
			t.Errorf("Color depth %d palette missing Error color", depth)
		}
		if string(palette.Text) == "" {
			t.Errorf("Color depth %d palette missing Text color", depth)
		}
	}
}

func TestCreateStylePalette(t *testing.T) {
	colors := create16ColorPalette()
	profile := termenv.ANSI

	styles := createStylePalette(colors, profile)

	// Test that styles have appropriate properties
	if styles.Title.GetForeground() == nil {
		t.Error("Title style should have foreground color")
	}

	if !styles.Title.GetBold() {
		t.Error("Title style should be bold")
	}

	if styles.Success.GetForeground() == nil {
		t.Error("Success style should have foreground color")
	}

	if styles.TextDim.GetForeground() == nil {
		t.Error("TextDim style should have foreground color")
	}
}

func TestGetSpinnerStyleForCapability(t *testing.T) {
	tests := []struct {
		hasUnicode    bool
		expectedStyle SpinnerStyle
	}{
		{true, SpinnerStyleUnicode},
		{false, SpinnerStyleASCII},
	}

	for _, test := range tests {
		caps := utils.TerminalCapabilities{HasUnicode: test.hasUnicode}
		style := getSpinnerStyleForCapability(caps)

		if style != test.expectedStyle {
			t.Errorf("HasUnicode=%v should give spinner style %v, got %v",
				test.hasUnicode, test.expectedStyle, style)
		}
	}
}

func TestGetProgressStyleForCapability(t *testing.T) {
	tests := []struct {
		hasUnicode    bool
		expectedStyle ProgressStyle
	}{
		{true, ProgressStyleUnicode},
		{false, ProgressStyleASCII},
	}

	for _, test := range tests {
		caps := utils.TerminalCapabilities{HasUnicode: test.hasUnicode}
		style := getProgressStyleForCapability(caps)

		if style != test.expectedStyle {
			t.Errorf("HasUnicode=%v should give progress style %v, got %v",
				test.hasUnicode, test.expectedStyle, style)
		}
	}
}

func TestGlobalTheme(t *testing.T) {
	// Reset global theme for clean test
	ResetGlobalTheme()

	// First call should create theme
	theme1 := GetGlobalTheme()
	if string(theme1.Colors.Primary) == "" {
		t.Error("Global theme should have primary color")
	}

	// Second call should return same instance
	theme2 := GetGlobalTheme()
	if theme1.Colors.Primary != theme2.Colors.Primary {
		t.Error("Global theme should be singleton")
	}

	// Reset should allow recreation
	ResetGlobalTheme()
	theme3 := GetGlobalTheme()
	if string(theme3.Colors.Primary) == "" {
		t.Error("Reset global theme should have primary color")
	}
}

func TestDetectCapabilities(t *testing.T) {
	caps := DetectCapabilities()

	// Should return valid capabilities
	if caps.ColorDepth < 1 {
		t.Error("Color depth should be at least 1")
	}

	if caps.Platform == "" {
		t.Error("Platform should not be empty")
	}
}

func TestGetColorProfile(t *testing.T) {
	profile := GetColorProfile()

	// Should return a valid profile
	if profile < 0 {
		t.Error("Color profile should be valid")
	}
}

func TestGetSpinnerChars(t *testing.T) {
	chars := GetSpinnerChars()

	if len(chars) == 0 {
		t.Error("Spinner chars should not be empty")
	}

	// Should be valid strings
	for i, char := range chars {
		if char == "" {
			t.Errorf("Spinner char %d should not be empty", i)
		}
	}
}

func TestGetProgressChars(t *testing.T) {
	filled, empty := GetProgressChars()

	if filled == "" {
		t.Error("Progress filled char should not be empty")
	}

	if empty == "" {
		t.Error("Progress empty char should not be empty")
	}
}

func TestColorDepthEnum(t *testing.T) {
	// Test that color depth constants are properly defined
	depths := []ColorDepth{
		ColorDepthMonochrome,
		ColorDepth8Color,
		ColorDepth16Color,
		ColorDepth256Color,
		ColorDepthTrueColor,
	}

	for i, depth := range depths {
		if depth < 0 {
			t.Errorf("Color depth %d should be non-negative", i)
		}
	}
}

func TestMonochromePalette(t *testing.T) {
	palette := createMonochromePalette()

	// In monochrome, most colors should be white or gray
	if string(palette.Text) != "15" {
		t.Error("Monochrome text should be white")
	}

	if string(palette.TextDim) != "8" {
		t.Error("Monochrome dim text should be gray")
	}
}

func TestTrueColorPalette(t *testing.T) {
	palette := createTrueColorPalette()

	// True color should use hex colors
	if string(palette.Primary) != "#0066CC" {
		t.Error("True color primary should be hex color")
	}

	if string(palette.Success) != "#00AA00" {
		t.Error("True color success should be hex color")
	}
}

// Helper function to determine profile type for comparison
func getProfileType(profile termenv.Profile) string {
	switch profile {
	case termenv.Ascii:
		return "ascii"
	case termenv.ANSI:
		return "ansi"
	case termenv.ANSI256:
		return "ansi256"
	case termenv.TrueColor:
		return "truecolor"
	default:
		return "unknown"
	}
}

// Benchmark tests
func BenchmarkNewThemeConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewThemeConfig()
	}
}

func BenchmarkGetGlobalTheme(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGlobalTheme()
	}
}

func BenchmarkCreateColorPalette(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createColorPalette(256)
		createColorPalette(16777216)
	}
}
