package common

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopedro/shotgun/internal/styles"
)

// Animation and timing constants for consistent behavior across components
const (
	// Animation frame rate targeting 60 FPS for smooth animations
	AnimationFrameRate = 60
	AnimationTickRate  = time.Millisecond * (1000 / AnimationFrameRate) // ~16ms for 60 FPS

	// Progress update rate for real-time feedback
	ProgressUpdateRate = time.Millisecond * 100 // 10 FPS for progress updates

	// Anti-flicker minimum display duration
	MinDisplayDuration = time.Millisecond * 500

	// Spinner animation rate - slightly slower for better readability
	SpinnerTickRate = time.Millisecond * 80 // ~12.5 FPS
)

// Color palette for consistent theming across all components
var (
	// Primary colors
	ColorPrimary   = lipgloss.Color("12") // Blue
	ColorSecondary = lipgloss.Color("14") // Cyan
	ColorAccent    = lipgloss.Color("10") // Green

	// Status colors
	ColorSuccess = lipgloss.Color("10") // Green
	ColorWarning = lipgloss.Color("11") // Yellow
	ColorError   = lipgloss.Color("9")  // Red
	ColorInfo    = lipgloss.Color("12") // Blue

	// Text colors
	ColorText     = lipgloss.Color("15") // White
	ColorTextDim  = lipgloss.Color("8")  // Gray
	ColorTextBold = lipgloss.Color("15") // White (for bold text)

	// Background colors for highlights
	ColorHighlight    = lipgloss.Color("12") // Blue background
	ColorHighlightDim = lipgloss.Color("8")  // Gray background

	// Progress specific colors
	ColorProgressBar  = lipgloss.Color("10") // Green
	ColorProgressBg   = lipgloss.Color("8")  // Gray background
	ColorProgressText = lipgloss.Color("14") // Cyan for info text
)

// Common style definitions for consistent appearance
var (
	// Base styles
	StyleBase = lipgloss.NewStyle()

	// Title styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	StyleSubtitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	// Text styles
	StyleText = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleTextDim = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Italic(true)

	StyleTextBold = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTextBold)

	// Status styles
	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	// Highlight styles
	StyleHighlight = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText).
			Background(ColorHighlight).
			Padding(0, 1)

	// Progress styles
	StyleProgressBar = lipgloss.NewStyle().
				Foreground(ColorProgressBar)

	StyleProgressInfo = lipgloss.NewStyle().
				Foreground(ColorProgressText)

	StyleProgressETA = lipgloss.NewStyle().
				Foreground(ColorTextDim).
				Italic(true)

	// Spinner styles
	StyleSpinner = lipgloss.NewStyle().
			Foreground(ColorPrimary)

	StyleSpinnerMessage = lipgloss.NewStyle().
				Foreground(ColorText).
				Italic(true)

	// Cancel hint style
	StyleCancelHint = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Italic(true)
)

// Progress bar characters for consistent appearance
const (
	ProgressBarFilled = "█"
	ProgressBarEmpty  = "░"
	ProgressBarBorder = "│"
	ProgressArrow     = "→"
	ProgressCheck     = "✓"
	ProgressDot       = "•"
)

// Animation characters and symbols
var (
	SpinnerDots   = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	SpinnerLine   = []string{"|", "/", "-", "\\"}
	SpinnerCircle = []string{"◐", "◓", "◑", "◒"}
)

// GetThemeAwareStyles returns styles configured for the current terminal
func GetThemeAwareStyles() styles.StylePalette {
	return styles.GetGlobalTheme().Styles
}

// GetThemeAwareColors returns colors configured for the current terminal
func GetThemeAwareColors() styles.ColorPalette {
	return styles.GetGlobalTheme().Colors
}

// InitializeTheme ensures theme is initialized and available
func InitializeTheme() {
	// This will create the global theme if it doesn't exist
	_ = styles.GetGlobalTheme()
}
