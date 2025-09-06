package styles

import (
	"github.com/diogopedro/shotgun/internal/utils"
)

// SpinnerStyle represents different spinner rendering modes
type SpinnerStyle int

const (
	SpinnerStyleUnicode SpinnerStyle = iota
	SpinnerStyleASCII
)

// ProgressStyle represents different progress bar rendering modes
type ProgressStyle int

const (
	ProgressStyleUnicode ProgressStyle = iota
	ProgressStyleASCII
	ProgressStyleMinimal
)

// FallbackConfig holds configuration for ASCII fallbacks
type FallbackConfig struct {
	SpinnerStyle  SpinnerStyle
	ProgressStyle ProgressStyle
	UseUnicode    bool
	Capability    utils.UnicodeCapability
}

// NewFallbackConfig creates a configuration based on terminal capabilities
func NewFallbackConfig() FallbackConfig {
	capability := utils.DetectUnicodeCapability()

	config := FallbackConfig{
		Capability: capability,
		UseUnicode: capability != utils.UnicodeNone,
	}

	// Set spinner style based on capability
	switch capability {
	case utils.UnicodeNone:
		config.SpinnerStyle = SpinnerStyleASCII
		config.ProgressStyle = ProgressStyleASCII
	case utils.UnicodeBasic:
		config.SpinnerStyle = SpinnerStyleASCII     // Use ASCII for spinners even in basic mode
		config.ProgressStyle = ProgressStyleUnicode // Use Unicode for progress bars in basic mode
	case utils.UnicodeFull:
		config.SpinnerStyle = SpinnerStyleUnicode
		config.ProgressStyle = ProgressStyleUnicode
	}

	return config
}

// SpinnerCharacters returns appropriate spinner characters
func (c FallbackConfig) SpinnerCharacters() []string {
	switch c.SpinnerStyle {
	case SpinnerStyleUnicode:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	case SpinnerStyleASCII:
		return []string{"|", "/", "-", "\\"}
	default:
		return []string{"|", "/", "-", "\\"}
	}
}

// ProgressCharacters returns appropriate progress bar characters
func (c FallbackConfig) ProgressCharacters() (filled, empty string) {
	switch c.ProgressStyle {
	case ProgressStyleUnicode:
		return "█", "░"
	case ProgressStyleASCII:
		return "#", "-"
	case ProgressStyleMinimal:
		return "=", " "
	default:
		return "#", "-"
	}
}

// BorderCharacters returns appropriate border characters
func (c FallbackConfig) BorderCharacters() BorderChars {
	if c.UseUnicode && c.Capability == utils.UnicodeFull {
		return BorderChars{
			TopLeft:     "┌",
			TopRight:    "┐",
			BottomLeft:  "└",
			BottomRight: "┘",
			Horizontal:  "─",
			Vertical:    "│",
		}
	}

	return BorderChars{
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Horizontal:  "-",
		Vertical:    "|",
	}
}

// BorderChars holds all border character definitions
type BorderChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

// StatusSymbols returns appropriate status symbols
func (c FallbackConfig) StatusSymbols() StatusSymbols {
	if c.UseUnicode && c.Capability == utils.UnicodeFull {
		return StatusSymbols{
			Success: "✓",
			Error:   "✗",
			Warning: "⚠",
			Info:    "ⓘ",
			Arrow:   "→",
			Dot:     "•",
		}
	}

	return StatusSymbols{
		Success: "+",
		Error:   "x",
		Warning: "!",
		Info:    "i",
		Arrow:   ">",
		Dot:     "*",
	}
}

// StatusSymbols holds all status symbol definitions
type StatusSymbols struct {
	Success string
	Error   string
	Warning string
	Info    string
	Arrow   string
	Dot     string
}

// SanitizeText ensures text is compatible with the terminal
func (c FallbackConfig) SanitizeText(text string) string {
	return utils.SanitizeForTerminal(text, c.Capability)
}

// GetFallbackSpinnerFrames returns spinner frames with proper fallbacks
func GetFallbackSpinnerFrames() map[SpinnerStyle][]string {
	return map[SpinnerStyle][]string{
		SpinnerStyleUnicode: {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		SpinnerStyleASCII:   {"|", "/", "-", "\\", "|", "/", "-", "\\", "|", "/"},
	}
}

// GetFallbackProgressChars returns progress bar characters with fallbacks
func GetFallbackProgressChars() map[ProgressStyle][2]string {
	return map[ProgressStyle][2]string{
		ProgressStyleUnicode: {"█", "░"},
		ProgressStyleASCII:   {"#", "-"},
		ProgressStyleMinimal: {"=", " "},
	}
}

// TestTerminalFeatures tests what Unicode features the current terminal supports
func TestTerminalFeatures() TerminalFeatures {
	capability := utils.DetectUnicodeCapability()
	platformInfo := utils.GetPlatformInfo()

	return TerminalFeatures{
		SupportsUnicode:     capability != utils.UnicodeNone,
		SupportsFullUnicode: capability == utils.UnicodeFull,
		SupportsColor:       platformInfo.Capabilities.ColorDepth > 1,
		ColorDepth:          platformInfo.Capabilities.ColorDepth,
		Platform:            platformInfo.OS,
		Terminal:            platformInfo.Terminal,
		FKeys:               platformInfo.Capabilities.HasF1F10Keys,
		ESCSupport:          platformInfo.Capabilities.ESCSupport,
	}
}

// TerminalFeatures describes what features the terminal supports
type TerminalFeatures struct {
	SupportsUnicode     bool
	SupportsFullUnicode bool
	SupportsColor       bool
	ColorDepth          int
	Platform            string
	Terminal            string
	FKeys               bool
	ESCSupport          bool
}

// String returns a human-readable description of terminal features
func (tf TerminalFeatures) String() string {
	unicode := "None"
	if tf.SupportsFullUnicode {
		unicode = "Full"
	} else if tf.SupportsUnicode {
		unicode = "Basic"
	}

	return "Platform: " + tf.Platform +
		", Terminal: " + tf.Terminal +
		", Unicode: " + unicode +
		", Colors: " + colorDepthString(tf.ColorDepth) +
		", F-Keys: " + boolString(tf.FKeys) +
		", ESC: " + boolString(tf.ESCSupport)
}

func colorDepthString(depth int) string {
	switch depth {
	case 1:
		return "Monochrome"
	case 8:
		return "8-color"
	case 16:
		return "16-color"
	case 256:
		return "256-color"
	case 16777216:
		return "True-color"
	default:
		return "Unknown"
	}
}

func boolString(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
