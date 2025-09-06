package utils

import (
	"fmt"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

// KeyboardCapabilities represents platform-specific keyboard support
type KeyboardCapabilities struct {
	Platform        string
	Terminal        string
	SupportsF1F10   bool
	SupportsESC     bool
	FKeyMappings    map[string]string // Function key -> expected sequence
	KnownIssues     []string
	TestedTerminals []string
	RecommendedKeys []string
}

// GetKeyboardCapabilities returns keyboard capabilities for current platform
func GetKeyboardCapabilities() KeyboardCapabilities {
	platformInfo := GetPlatformInfo()

	caps := KeyboardCapabilities{
		Platform:        platformInfo.OS,
		Terminal:        platformInfo.Terminal,
		SupportsF1F10:   platformInfo.Capabilities.HasF1F10Keys,
		SupportsESC:     platformInfo.Capabilities.ESCSupport,
		FKeyMappings:    make(map[string]string),
		RecommendedKeys: []string{"ESC", "Enter", "Tab", "Arrow Keys"},
	}

	// Set platform-specific details
	switch runtime.GOOS {
	case "windows":
		caps = configureWindowsKeyboard(caps, platformInfo.Terminal)
	case "darwin":
		caps = configureMacOSKeyboard(caps, platformInfo.Terminal)
	case "linux":
		caps = configureLinuxKeyboard(caps, platformInfo.Terminal)
	default:
		caps = configureUnknownKeyboard(caps)
	}

	return caps
}

// configureWindowsKeyboard sets up Windows-specific keyboard configuration
func configureWindowsKeyboard(caps KeyboardCapabilities, terminal string) KeyboardCapabilities {
	switch terminal {
	case "cmd":
		caps.SupportsF1F10 = false
		caps.KnownIssues = []string{
			"F1-F10 keys may not be properly detected",
			"Some Unicode characters may display incorrectly",
			"Limited ANSI escape sequence support",
		}
		caps.TestedTerminals = []string{"Windows CMD"}
		caps.RecommendedKeys = []string{"ESC", "Enter", "Tab", "Arrow Keys"}

	case "powershell_5":
		caps.SupportsF1F10 = false
		caps.KnownIssues = []string{
			"F1-F10 keys may be intercepted by PowerShell",
			"Some function keys trigger PowerShell shortcuts",
			"Unicode support limited",
		}
		caps.TestedTerminals = []string{"Windows PowerShell 5.1"}
		caps.RecommendedKeys = []string{"ESC", "Enter", "Tab", "Arrow Keys", "Ctrl+C"}

	case "powershell_7":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"Generally good key support",
			"F1 may trigger help in some contexts",
		}
		caps.TestedTerminals = []string{"PowerShell 7.x"}
		caps.RecommendedKeys = []string{"F1", "ESC", "Enter", "Tab", "Arrow Keys", "Ctrl+C"}

	case "windows_terminal":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"Excellent key support",
			"Modern terminal features available",
		}
		caps.TestedTerminals = []string{"Windows Terminal"}
		caps.RecommendedKeys = []string{"F1-F10", "ESC", "Enter", "Tab", "Arrow Keys", "Ctrl+C"}

	default:
		caps.KnownIssues = []string{"Unknown Windows terminal, behavior may vary"}
		caps.TestedTerminals = []string{"Unknown"}
	}

	// Common Windows F-key mappings (when supported)
	if caps.SupportsF1F10 {
		caps.FKeyMappings = map[string]string{
			"F1":  "\\x1b[OP",  // Help
			"F2":  "\\x1b[OQ",  // Rename/Edit
			"F3":  "\\x1b[OR",  // Find/Search
			"F4":  "\\x1b[OS",  // Address bar
			"F5":  "\\x1b[15~", // Refresh
			"F6":  "\\x1b[17~", // Focus address bar
			"F7":  "\\x1b[18~", // Spell check
			"F8":  "\\x1b[19~", // Extend selection
			"F9":  "\\x1b[20~", // Calculate fields
			"F10": "\\x1b[21~", // Menu activation
		}
	}

	return caps
}

// configureMacOSKeyboard sets up macOS-specific keyboard configuration
func configureMacOSKeyboard(caps KeyboardCapabilities, terminal string) KeyboardCapabilities {
	switch terminal {
	case "terminal":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"F1-F4 may be intercepted by system shortcuts",
			"Generally excellent key support",
		}
		caps.TestedTerminals = []string{"macOS Terminal.app"}
		caps.RecommendedKeys = []string{"F5-F10", "ESC", "Enter", "Tab", "Arrow Keys", "Cmd+C"}

	case "iterm2":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"Excellent key support with customizable key mappings",
			"Advanced terminal features available",
		}
		caps.TestedTerminals = []string{"iTerm2"}
		caps.RecommendedKeys = []string{"F1-F10", "ESC", "Enter", "Tab", "Arrow Keys", "Cmd+C"}

	default:
		caps.KnownIssues = []string{"Unknown macOS terminal, generally good support expected"}
		caps.TestedTerminals = []string{"Unknown"}
	}

	// macOS F-key mappings
	caps.FKeyMappings = map[string]string{
		"F1":  "\\x1b[OP",  // Help (may be system intercepted)
		"F2":  "\\x1b[OQ",  // Rename (may be system intercepted)
		"F3":  "\\x1b[OR",  // Mission Control (may be system intercepted)
		"F4":  "\\x1b[OS",  // Launchpad (may be system intercepted)
		"F5":  "\\x1b[15~", // Available
		"F6":  "\\x1b[17~", // Available
		"F7":  "\\x1b[18~", // Available
		"F8":  "\\x1b[19~", // Available
		"F9":  "\\x1b[20~", // Available
		"F10": "\\x1b[21~", // Available
	}

	return caps
}

// configureLinuxKeyboard sets up Linux-specific keyboard configuration
func configureLinuxKeyboard(caps KeyboardCapabilities, terminal string) KeyboardCapabilities {
	switch terminal {
	case "gnome_terminal":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"Generally excellent key support",
			"F10 may activate menu bar in some configurations",
		}
		caps.TestedTerminals = []string{"GNOME Terminal"}
		caps.RecommendedKeys = []string{"F1-F9", "ESC", "Enter", "Tab", "Arrow Keys", "Ctrl+C"}

	case "konsole":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"Excellent key support with KDE integration",
			"Highly configurable key mappings",
		}
		caps.TestedTerminals = []string{"Konsole"}
		caps.RecommendedKeys = []string{"F1-F10", "ESC", "Enter", "Tab", "Arrow Keys", "Ctrl+C"}

	case "xterm":
		caps.SupportsF1F10 = true
		caps.KnownIssues = []string{
			"Basic but reliable key support",
			"May require additional configuration for best experience",
		}
		caps.TestedTerminals = []string{"xterm"}
		caps.RecommendedKeys = []string{"F1-F10", "ESC", "Enter", "Tab", "Arrow Keys", "Ctrl+C"}

	default:
		caps.KnownIssues = []string{"Unknown Linux terminal, generally good support expected"}
		caps.TestedTerminals = []string{"Unknown"}
	}

	// Linux F-key mappings (standard)
	caps.FKeyMappings = map[string]string{
		"F1":  "\\x1b[OP",  // Help
		"F2":  "\\x1b[OQ",  // Menu or context actions
		"F3":  "\\x1b[OR",  // Find/Search
		"F4":  "\\x1b[OS",  // Open/Edit
		"F5":  "\\x1b[15~", // Refresh
		"F6":  "\\x1b[17~", // Navigate
		"F7":  "\\x1b[18~", // Spell check
		"F8":  "\\x1b[19~", // Extend/Select
		"F9":  "\\x1b[20~", // Build/Execute
		"F10": "\\x1b[21~", // Menu activation
	}

	return caps
}

// configureUnknownKeyboard sets up default keyboard configuration
func configureUnknownKeyboard(caps KeyboardCapabilities) KeyboardCapabilities {
	caps.SupportsF1F10 = false
	caps.KnownIssues = []string{
		"Unknown platform, function key support uncertain",
		"ESC key should work on most platforms",
	}
	caps.TestedTerminals = []string{"Unknown"}
	caps.RecommendedKeys = []string{"ESC", "Enter", "Tab", "Arrow Keys"}
	caps.FKeyMappings = map[string]string{}

	return caps
}

// IsKeySupported checks if a specific key is supported on current platform
func IsKeySupported(key string) bool {
	caps := GetKeyboardCapabilities()

	switch key {
	case "ESC", "Escape":
		return caps.SupportsESC
	case "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10":
		return caps.SupportsF1F10
	case "Enter", "Tab", "Up", "Down", "Left", "Right":
		return true // These should work on all platforms
	case "Ctrl+C":
		return true // Universal interrupt
	default:
		return true // Assume other keys work unless specifically known not to
	}
}

// GetRecommendedKeys returns keys that are recommended for the current platform
func GetRecommendedKeys() []string {
	return GetKeyboardCapabilities().RecommendedKeys
}

// GetKeyMappingInfo returns human-readable key mapping information
func GetKeyMappingInfo() string {
	caps := GetKeyboardCapabilities()

	info := fmt.Sprintf("Platform: %s (%s)\n", caps.Platform, caps.Terminal)
	info += fmt.Sprintf("F1-F10 Support: %t\n", caps.SupportsF1F10)
	info += fmt.Sprintf("ESC Support: %t\n", caps.SupportsESC)

	if len(caps.RecommendedKeys) > 0 {
		info += fmt.Sprintf("Recommended Keys: %v\n", caps.RecommendedKeys)
	}

	if len(caps.KnownIssues) > 0 {
		info += "Known Issues:\n"
		for _, issue := range caps.KnownIssues {
			info += fmt.Sprintf("  - %s\n", issue)
		}
	}

	if len(caps.TestedTerminals) > 0 {
		info += fmt.Sprintf("Tested On: %v\n", caps.TestedTerminals)
	}

	return info
}

// TestKeyboardInput provides a simple test function for keyboard input
func TestKeyboardInput(msg tea.KeyMsg) string {
	caps := GetKeyboardCapabilities()

	keyName := msg.String()
	supported := IsKeySupported(keyName)

	result := fmt.Sprintf("Key: %s, Supported: %t", keyName, supported)

	if mapping, exists := caps.FKeyMappings[keyName]; exists {
		result += fmt.Sprintf(", Expected Sequence: %s", mapping)
	}

	if !supported {
		result += " (Not recommended on this platform)"
	}

	return result
}

// KeyboardTestMatrix returns the complete testing matrix for documentation
func KeyboardTestMatrix() []KeyboardCapabilities {
	// This would normally be filled with actual test results
	// For now, return the expected capabilities for each platform

	matrices := []KeyboardCapabilities{}

	platforms := []string{"windows", "darwin", "linux"}
	terminals := map[string][]string{
		"windows": {"cmd", "powershell_5", "powershell_7", "windows_terminal"},
		"darwin":  {"terminal", "iterm2"},
		"linux":   {"gnome_terminal", "konsole", "xterm"},
	}

	for _, platform := range platforms {
		for _, terminal := range terminals[platform] {
			caps := KeyboardCapabilities{
				Platform: platform,
				Terminal: terminal,
			}

			switch platform {
			case "windows":
				caps = configureWindowsKeyboard(caps, terminal)
			case "darwin":
				caps = configureMacOSKeyboard(caps, terminal)
			case "linux":
				caps = configureLinuxKeyboard(caps, terminal)
			}

			matrices = append(matrices, caps)
		}
	}

	return matrices
}
