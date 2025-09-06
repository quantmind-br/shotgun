package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/diogopedro/shotgun/internal/utils"
	"github.com/muesli/termenv"
)

// RenderConfig contains all configuration for terminal-compatible rendering
type RenderConfig struct {
	UseUnicode    bool
	ColorProfile  termenv.Profile
	SpinnerStyle  SpinnerStyle
	ProgressStyle ProgressStyle
	Capability    utils.TerminalCapabilities
}

// ThemeConfig represents the complete theme configuration for the application
type ThemeConfig struct {
	RenderConfig RenderConfig
	Colors       ColorPalette
	Styles       StylePalette
}

// ColorPalette contains all colors used in the application with fallbacks
type ColorPalette struct {
	// Primary brand colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color

	// Status colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color

	// Text colors
	Text     lipgloss.Color
	TextDim  lipgloss.Color
	TextBold lipgloss.Color

	// Background colors
	Background   lipgloss.Color
	Highlight    lipgloss.Color
	HighlightDim lipgloss.Color

	// Progress specific colors
	ProgressBar  lipgloss.Color
	ProgressBg   lipgloss.Color
	ProgressText lipgloss.Color
}

// StylePalette contains all pre-configured lipgloss styles
type StylePalette struct {
	// Base styles
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Text     lipgloss.Style
	TextDim  lipgloss.Style
	TextBold lipgloss.Style

	// Status styles
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style

	// Interactive styles
	Highlight lipgloss.Style
	Selected  lipgloss.Style

	// Progress styles
	ProgressBar  lipgloss.Style
	ProgressInfo lipgloss.Style
	ProgressETA  lipgloss.Style

	// Spinner styles
	Spinner        lipgloss.Style
	SpinnerMessage lipgloss.Style

	// UI element styles
	CancelHint lipgloss.Style
	Border     lipgloss.Style
}

// ColorDepth represents different color capabilities
type ColorDepth int

const (
	ColorDepthMonochrome ColorDepth = 1
	ColorDepth8Color     ColorDepth = 8
	ColorDepth16Color    ColorDepth = 16
	ColorDepth256Color   ColorDepth = 256
	ColorDepthTrueColor  ColorDepth = 16777216
)

// NewThemeConfig creates a theme configuration based on terminal capabilities
func NewThemeConfig() ThemeConfig {
	caps := utils.DetectTerminalCapabilities()
	colorProfile := detectColorProfile(caps)

	return ThemeConfig{
		RenderConfig: RenderConfig{
			UseUnicode:    caps.HasUnicode,
			ColorProfile:  colorProfile,
			SpinnerStyle:  getSpinnerStyleForCapability(caps),
			ProgressStyle: getProgressStyleForCapability(caps),
			Capability:    caps,
		},
		Colors: createColorPalette(caps.ColorDepth),
		Styles: createStylePalette(createColorPalette(caps.ColorDepth), colorProfile),
	}
}

// DetectCapabilities returns comprehensive terminal capabilities
func DetectCapabilities() utils.TerminalCapabilities {
	return utils.DetectTerminalCapabilities()
}

// GetColorProfile returns the appropriate lipgloss color profile
func GetColorProfile() termenv.Profile {
	caps := DetectCapabilities()
	return detectColorProfile(caps)
}

// GetSpinnerChars returns spinner characters appropriate for terminal
func GetSpinnerChars() []string {
	fallback := NewFallbackConfig()
	return fallback.SpinnerCharacters()
}

// GetProgressChars returns progress bar characters appropriate for terminal
func GetProgressChars() (filled, empty string) {
	fallback := NewFallbackConfig()
	return fallback.ProgressCharacters()
}

// detectColorProfile maps terminal capabilities to lipgloss profiles
func detectColorProfile(caps utils.TerminalCapabilities) termenv.Profile {
	// Use the color depth from terminal capabilities (which already handles env vars)
	switch ColorDepth(caps.ColorDepth) {
	case ColorDepthMonochrome:
		return termenv.Ascii
	case ColorDepth8Color:
		return termenv.ANSI
	case ColorDepth16Color:
		return termenv.ANSI
	case ColorDepth256Color:
		return termenv.ANSI256
	case ColorDepthTrueColor:
		return termenv.TrueColor
	default:
		return termenv.ANSI // Safe default
	}
}

// createColorPalette creates a color palette appropriate for the color depth
func createColorPalette(colorDepth int) ColorPalette {
	switch ColorDepth(colorDepth) {
	case ColorDepthMonochrome:
		return createMonochromePalette()
	case ColorDepth8Color:
		return create8ColorPalette()
	case ColorDepth16Color:
		return create16ColorPalette()
	case ColorDepth256Color:
		return create256ColorPalette()
	case ColorDepthTrueColor:
		return createTrueColorPalette()
	default:
		return create16ColorPalette() // Safe default
	}
}

// createMonochromePalette creates a monochrome color palette
func createMonochromePalette() ColorPalette {
	return ColorPalette{
		Primary:   lipgloss.Color("15"), // White (bold)
		Secondary: lipgloss.Color("15"), // White
		Accent:    lipgloss.Color("15"), // White

		Success: lipgloss.Color("15"), // White
		Warning: lipgloss.Color("15"), // White (bold)
		Error:   lipgloss.Color("15"), // White (bold)
		Info:    lipgloss.Color("15"), // White

		Text:     lipgloss.Color("15"), // White
		TextDim:  lipgloss.Color("8"),  // Gray/dim
		TextBold: lipgloss.Color("15"), // White

		Background:   lipgloss.Color("0"),  // Black
		Highlight:    lipgloss.Color("15"), // White background
		HighlightDim: lipgloss.Color("8"),  // Gray

		ProgressBar:  lipgloss.Color("15"), // White
		ProgressBg:   lipgloss.Color("8"),  // Gray
		ProgressText: lipgloss.Color("15"), // White
	}
}

// create8ColorPalette creates an 8-color ANSI palette
func create8ColorPalette() ColorPalette {
	return ColorPalette{
		Primary:   lipgloss.Color("4"), // Blue
		Secondary: lipgloss.Color("6"), // Cyan
		Accent:    lipgloss.Color("2"), // Green

		Success: lipgloss.Color("2"), // Green
		Warning: lipgloss.Color("3"), // Yellow
		Error:   lipgloss.Color("1"), // Red
		Info:    lipgloss.Color("4"), // Blue

		Text:     lipgloss.Color("7"),  // White
		TextDim:  lipgloss.Color("8"),  // Bright black (gray)
		TextBold: lipgloss.Color("15"), // Bright white

		Background:   lipgloss.Color("0"), // Black
		Highlight:    lipgloss.Color("4"), // Blue background
		HighlightDim: lipgloss.Color("8"), // Gray

		ProgressBar:  lipgloss.Color("2"), // Green
		ProgressBg:   lipgloss.Color("8"), // Gray
		ProgressText: lipgloss.Color("6"), // Cyan
	}
}

// create16ColorPalette creates a 16-color ANSI palette
func create16ColorPalette() ColorPalette {
	return ColorPalette{
		Primary:   lipgloss.Color("12"), // Bright Blue
		Secondary: lipgloss.Color("14"), // Bright Cyan
		Accent:    lipgloss.Color("10"), // Bright Green

		Success: lipgloss.Color("10"), // Bright Green
		Warning: lipgloss.Color("11"), // Bright Yellow
		Error:   lipgloss.Color("9"),  // Bright Red
		Info:    lipgloss.Color("12"), // Bright Blue

		Text:     lipgloss.Color("15"), // Bright White
		TextDim:  lipgloss.Color("8"),  // Gray
		TextBold: lipgloss.Color("15"), // Bright White

		Background:   lipgloss.Color("0"),  // Black
		Highlight:    lipgloss.Color("12"), // Bright Blue background
		HighlightDim: lipgloss.Color("8"),  // Gray

		ProgressBar:  lipgloss.Color("10"), // Bright Green
		ProgressBg:   lipgloss.Color("8"),  // Gray
		ProgressText: lipgloss.Color("14"), // Bright Cyan
	}
}

// create256ColorPalette creates a 256-color palette
func create256ColorPalette() ColorPalette {
	return ColorPalette{
		Primary:   lipgloss.Color("33"), // Blue
		Secondary: lipgloss.Color("51"), // Cyan
		Accent:    lipgloss.Color("46"), // Green

		Success: lipgloss.Color("46"),  // Green
		Warning: lipgloss.Color("226"), // Yellow
		Error:   lipgloss.Color("196"), // Red
		Info:    lipgloss.Color("33"),  // Blue

		Text:     lipgloss.Color("255"), // White
		TextDim:  lipgloss.Color("244"), // Gray
		TextBold: lipgloss.Color("255"), // White

		Background:   lipgloss.Color("0"),   // Black
		Highlight:    lipgloss.Color("33"),  // Blue background
		HighlightDim: lipgloss.Color("244"), // Gray

		ProgressBar:  lipgloss.Color("46"),  // Green
		ProgressBg:   lipgloss.Color("244"), // Gray
		ProgressText: lipgloss.Color("51"),  // Cyan
	}
}

// createTrueColorPalette creates a true color (24-bit) palette
func createTrueColorPalette() ColorPalette {
	return ColorPalette{
		Primary:   lipgloss.Color("#0066CC"), // Blue
		Secondary: lipgloss.Color("#00CCCC"), // Cyan
		Accent:    lipgloss.Color("#00AA00"), // Green

		Success: lipgloss.Color("#00AA00"), // Green
		Warning: lipgloss.Color("#FFAA00"), // Orange
		Error:   lipgloss.Color("#CC0000"), // Red
		Info:    lipgloss.Color("#0066CC"), // Blue

		Text:     lipgloss.Color("#FFFFFF"), // White
		TextDim:  lipgloss.Color("#888888"), // Gray
		TextBold: lipgloss.Color("#FFFFFF"), // White

		Background:   lipgloss.Color("#000000"), // Black
		Highlight:    lipgloss.Color("#0066CC"), // Blue background
		HighlightDim: lipgloss.Color("#444444"), // Dark Gray

		ProgressBar:  lipgloss.Color("#00AA00"), // Green
		ProgressBg:   lipgloss.Color("#444444"), // Dark Gray
		ProgressText: lipgloss.Color("#00CCCC"), // Cyan
	}
}

// createStylePalette creates lipgloss styles using the color palette
func createStylePalette(colors ColorPalette, profile termenv.Profile) StylePalette {
	return StylePalette{
		// Base styles
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Primary),

		Subtitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Secondary),

		Text: lipgloss.NewStyle().
			Foreground(colors.Text),

		TextDim: lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Italic(true),

		TextBold: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.TextBold),

		// Status styles
		Success: lipgloss.NewStyle().
			Foreground(colors.Success).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(colors.Warning).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(colors.Error).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(colors.Info),

		// Interactive styles
		Highlight: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Text).
			Background(colors.Highlight).
			Padding(0, 1),

		Selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Background).
			Background(colors.Primary),

		// Progress styles
		ProgressBar: lipgloss.NewStyle().
			Foreground(colors.ProgressBar),

		ProgressInfo: lipgloss.NewStyle().
			Foreground(colors.ProgressText),

		ProgressETA: lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Italic(true),

		// Spinner styles
		Spinner: lipgloss.NewStyle().
			Foreground(colors.Primary),

		SpinnerMessage: lipgloss.NewStyle().
			Foreground(colors.Text).
			Italic(true),

		// UI element styles
		CancelHint: lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Italic(true),

		Border: lipgloss.NewStyle().
			Foreground(colors.TextDim),
	}
}

// getSpinnerStyleForCapability determines appropriate spinner style
func getSpinnerStyleForCapability(caps utils.TerminalCapabilities) SpinnerStyle {
	if caps.HasUnicode {
		return SpinnerStyleUnicode
	}
	return SpinnerStyleASCII
}

// getProgressStyleForCapability determines appropriate progress style
func getProgressStyleForCapability(caps utils.TerminalCapabilities) ProgressStyle {
	if caps.HasUnicode {
		return ProgressStyleUnicode
	}
	return ProgressStyleASCII
}

// GetGlobalTheme returns the global theme instance
var globalTheme *ThemeConfig

func GetGlobalTheme() ThemeConfig {
	if globalTheme == nil {
		theme := NewThemeConfig()
		globalTheme = &theme
	}
	return *globalTheme
}

// ResetGlobalTheme forces recreation of the global theme
func ResetGlobalTheme() {
	globalTheme = nil
}
