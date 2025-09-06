package utils

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
)

// PlatformInfo contains platform and terminal information
type PlatformInfo struct {
	OS           string
	Terminal     string
	Capabilities TerminalCapabilities
}

// TerminalCapabilities represents terminal feature support
type TerminalCapabilities struct {
	HasUnicode   bool
	ColorDepth   int // 1, 8, 16, 256, 16777216
	HasF1F10Keys bool
	Platform     string
	ESCSupport   bool
}

// PlatformSupport represents compatibility matrix entry
type PlatformSupport struct {
	Platform     string // windows, darwin, linux
	Terminal     string // cmd, powershell, terminal, iterm2, etc.
	UnicodeLevel string // none, basic, full
	ColorSupport int    // 1, 8, 16, 256, 16777216
	FKeySupport  bool   // F1-F10 detection capability
	ESCSupport   bool   // ESC key handling
	Tested       bool   // Verification status
}

// GetPlatformInfo returns comprehensive platform and terminal information
func GetPlatformInfo() PlatformInfo {
	return PlatformInfo{
		OS:           runtime.GOOS,
		Terminal:     detectTerminal(),
		Capabilities: DetectTerminalCapabilities(),
	}
}

// DetectTerminalCapabilities performs runtime terminal capability detection
func DetectTerminalCapabilities() TerminalCapabilities {
	caps := TerminalCapabilities{
		Platform:   runtime.GOOS,
		ESCSupport: true, // Assume ESC support for most modern terminals
	}

	// Detect color depth - check environment variables first, then termenv
	caps.ColorDepth = detectColorDepth()

	// Detect Unicode support
	caps.HasUnicode = detectUnicodeSupport()

	// Detect F-key support (platform-specific heuristics)
	caps.HasF1F10Keys = detectFKeySupport()

	// Windows-specific adjustments
	if runtime.GOOS == "windows" {
		caps = adjustForWindows(caps)
	}

	return caps
}

// IsWindowsCMD detects if running in Windows Command Prompt
func IsWindowsCMD() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	// Check common CMD environment indicators
	comspec := os.Getenv("COMSPEC")
	if strings.Contains(strings.ToLower(comspec), "cmd.exe") {
		// Additional check: PowerShell sets PSModulePath
		if os.Getenv("PSModulePath") == "" {
			return true
		}
	}

	return false
}

// HasUnicodeSupport tests Unicode rendering capability
func HasUnicodeSupport() bool {
	return DetectTerminalCapabilities().HasUnicode
}

// GetColorDepth returns detected color support level
func GetColorDepth() int {
	return DetectTerminalCapabilities().ColorDepth
}

// detectColorDepth determines color support level with environment variable overrides
func detectColorDepth() int {
	// Check FORCE_COLOR environment variable first (takes highest precedence)
	if colorForce := os.Getenv("FORCE_COLOR"); colorForce != "" {
		if forced, err := strconv.Atoi(colorForce); err == nil {
			switch forced {
			case 0:
				return 1 // Monochrome
			case 1:
				return 8 // 8 colors
			case 2:
				return 256 // 256 colors
			case 3:
				return 16777216 // True color
			}
		}
	}

	// Check NO_COLOR environment variable (disables colors)
	if os.Getenv("NO_COLOR") != "" {
		return 1 // Monochrome
	}

	// Fall back to termenv detection
	profile := termenv.ColorProfile()
	switch profile {
	case termenv.Ascii:
		return 1
	case termenv.ANSI:
		return 8
	case termenv.ANSI256:
		return 256
	case termenv.TrueColor:
		return 16777216
	default:
		return 16 // Default to 16 colors for unknown profiles
	}
}

// detectTerminal identifies the specific terminal application
func detectTerminal() string {
	// Check environment variables that identify terminals
	if term := os.Getenv("TERM_PROGRAM"); term != "" {
		switch strings.ToLower(term) {
		case "iterm.app":
			return "iterm2"
		case "apple_terminal":
			return "terminal"
		case "vscode":
			return "vscode"
		}
	}

	if os.Getenv("WT_SESSION") != "" {
		return "windows_terminal"
	}

	// Check for PowerShell
	if os.Getenv("PSModulePath") != "" {
		if psVersion := os.Getenv("PSVersionTable"); psVersion != "" {
			if strings.Contains(psVersion, "7.") || strings.Contains(psVersion, "Core") {
				return "powershell_7"
			}
		}
		return "powershell_5"
	}

	// Check for CMD on Windows
	if IsWindowsCMD() {
		return "cmd"
	}

	// Check TERM for Linux terminals
	term := os.Getenv("TERM")
	switch {
	case strings.Contains(term, "xterm"):
		return "xterm"
	case strings.Contains(term, "gnome"):
		return "gnome_terminal"
	case strings.Contains(term, "konsole"):
		return "konsole"
	default:
		return term
	}
}

// detectUnicodeSupport tests for Unicode capability
func detectUnicodeSupport() bool {
	// Check environment variables
	lang := os.Getenv("LANG")
	lcAll := os.Getenv("LC_ALL")
	lcCtype := os.Getenv("LC_CTYPE")

	// UTF-8 in environment usually indicates Unicode support
	for _, env := range []string{lang, lcAll, lcCtype} {
		if strings.Contains(strings.ToUpper(env), "UTF-8") ||
			strings.Contains(strings.ToUpper(env), "UTF8") {
			return true
		}
	}

	// Windows-specific checks
	if runtime.GOOS == "windows" {
		// Windows Terminal and PowerShell 7 generally support Unicode
		terminal := detectTerminal()
		if terminal == "windows_terminal" || terminal == "powershell_7" {
			return true
		}
		// CMD has very limited Unicode support
		if terminal == "cmd" {
			return false
		}
		// PowerShell 5.1 has limited Unicode support
		if terminal == "powershell_5" {
			return false
		}
	}

	// macOS and modern Linux terminals generally support Unicode
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return true
	}

	// Default to false for unknown terminals to be safe
	return false
}

// detectFKeySupport detects F1-F10 key support
func detectFKeySupport() bool {
	terminal := detectTerminal()

	// Terminals with good F-key support
	goodFKeyTerminals := []string{
		"iterm2",
		"terminal", // macOS Terminal
		"windows_terminal",
		"powershell_7",
		"gnome_terminal",
		"konsole",
		"xterm",
	}

	for _, goodTerminal := range goodFKeyTerminals {
		if terminal == goodTerminal {
			return true
		}
	}

	// CMD and PowerShell 5.1 may have issues with F-keys
	if terminal == "cmd" || terminal == "powershell_5" {
		return false
	}

	// Default to true for unknown terminals
	return true
}

// adjustForWindows applies Windows-specific capability adjustments
func adjustForWindows(caps TerminalCapabilities) TerminalCapabilities {
	terminal := detectTerminal()

	switch terminal {
	case "cmd":
		// CMD has very limited capabilities
		caps.HasUnicode = false
		caps.ColorDepth = min(caps.ColorDepth, 16) // Limited to 16 colors
		caps.HasF1F10Keys = false
		caps.ESCSupport = true // CMD does support ESC

	case "powershell_5":
		// PowerShell 5.1 has limited Unicode
		caps.HasUnicode = false
		caps.ColorDepth = min(caps.ColorDepth, 256) // Better than CMD
		caps.HasF1F10Keys = false
		caps.ESCSupport = true

	case "powershell_7":
		// PowerShell 7 has good support
		caps.HasUnicode = true
		// Keep detected color depth
		caps.HasF1F10Keys = true
		caps.ESCSupport = true

	case "windows_terminal":
		// Windows Terminal has excellent support
		caps.HasUnicode = true
		// Keep detected color depth
		caps.HasF1F10Keys = true
		caps.ESCSupport = true
	}

	return caps
}

// GetTerminalMatrix returns the complete compatibility matrix
func GetTerminalMatrix() []PlatformSupport {
	return []PlatformSupport{
		// Windows platforms
		{
			Platform:     "windows",
			Terminal:     "cmd",
			UnicodeLevel: "none",
			ColorSupport: 16,
			FKeySupport:  false,
			ESCSupport:   true,
			Tested:       false,
		},
		{
			Platform:     "windows",
			Terminal:     "powershell_5",
			UnicodeLevel: "basic",
			ColorSupport: 256,
			FKeySupport:  false,
			ESCSupport:   true,
			Tested:       false,
		},
		{
			Platform:     "windows",
			Terminal:     "powershell_7",
			UnicodeLevel: "full",
			ColorSupport: 16777216,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
		{
			Platform:     "windows",
			Terminal:     "windows_terminal",
			UnicodeLevel: "full",
			ColorSupport: 16777216,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
		// macOS platforms
		{
			Platform:     "darwin",
			Terminal:     "terminal",
			UnicodeLevel: "full",
			ColorSupport: 16777216,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
		{
			Platform:     "darwin",
			Terminal:     "iterm2",
			UnicodeLevel: "full",
			ColorSupport: 16777216,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
		// Linux platforms
		{
			Platform:     "linux",
			Terminal:     "gnome_terminal",
			UnicodeLevel: "full",
			ColorSupport: 16777216,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
		{
			Platform:     "linux",
			Terminal:     "konsole",
			UnicodeLevel: "full",
			ColorSupport: 16777216,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
		{
			Platform:     "linux",
			Terminal:     "xterm",
			UnicodeLevel: "basic",
			ColorSupport: 256,
			FKeySupport:  true,
			ESCSupport:   true,
			Tested:       false,
		},
	}
}

// Helper function for Go versions that don't have min built-in
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
