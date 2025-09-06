package utils

import (
	"unicode"
	"unicode/utf8"
)

// UnicodeCapability represents different levels of Unicode support
type UnicodeCapability int

const (
	UnicodeNone UnicodeCapability = iota
	UnicodeBasic
	UnicodeFull
)

// String returns a string representation of the Unicode capability level
func (u UnicodeCapability) String() string {
	switch u {
	case UnicodeNone:
		return "none"
	case UnicodeBasic:
		return "basic"
	case UnicodeFull:
		return "full"
	default:
		return "unknown"
	}
}

// TestUnicodeCharacters contains characters for testing Unicode support levels
var TestUnicodeCharacters = struct {
	Basic []rune // Basic Unicode characters (should work in most environments)
	Full  []rune // Advanced Unicode characters (require full Unicode support)
}{
	Basic: []rune{
		'█', '▉', '▊', '▋', '▌', '▍', '▎', '▏', '░', // Block elements
		'←', '↑', '→', '↓', // Basic arrows
		'•', '◦', // Bullets
	},
	Full: []rune{
		'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏', // Braille patterns (spinner)
		'┌', '┐', '└', '┘', '│', '─', // Box drawing
		'▲', '▼', '◆', '●', '○', // Advanced symbols
		'✓', '✗', '⚠', // Status symbols
	},
}

// ASCIIFallbacks provides ASCII alternatives for Unicode characters
var ASCIIFallbacks = map[rune]rune{
	// Block elements
	'█': '#',
	'▉': '#',
	'▊': '#',
	'▋': '#',
	'▌': '#',
	'▍': '#',
	'▎': '#',
	'▏': '#',
	'░': '-',

	// Box drawing
	'┌': '+',
	'┐': '+',
	'└': '+',
	'┘': '+',
	'│': '|',
	'─': '-',

	// Arrows
	'←': '<',
	'↑': '^',
	'→': '>',
	'↓': 'v',

	// Symbols
	'•': '*',
	'◦': 'o',
	'▲': '^',
	'▼': 'v',
	'◆': '*',
	'●': '*',
	'○': 'o',

	// Status symbols
	'✓': '+',
	'✗': 'x',
	'⚠': '!',

	// Braille patterns (spinner fallbacks)
	'⠋': '|',
	'⠙': '/',
	'⠹': '-',
	'⠸': '\\',
	'⠼': '|',
	'⠴': '/',
	'⠦': '-',
	'⠧': '\\',
	'⠇': '|',
	'⠏': '/',
}

// DetectUnicodeCapability determines the level of Unicode support
func DetectUnicodeCapability() UnicodeCapability {
	caps := DetectTerminalCapabilities()

	if !caps.HasUnicode {
		return UnicodeNone
	}

	// For now, assume full Unicode if HasUnicode is true
	// In the future, we could add more sophisticated testing
	return UnicodeFull
}

// IsValidUTF8 checks if a string is valid UTF-8
func IsValidUTF8(s string) bool {
	return utf8.ValidString(s)
}

// ContainsUnicode checks if a string contains Unicode characters beyond ASCII
func ContainsUnicode(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

// ToASCIIFallback converts Unicode characters to ASCII fallbacks
func ToASCIIFallback(r rune) rune {
	if fallback, exists := ASCIIFallbacks[r]; exists {
		return fallback
	}

	// If no specific fallback, check if it's a printable ASCII character
	if r <= unicode.MaxASCII && unicode.IsPrint(r) {
		return r
	}

	// Default fallback for unprintable or unknown Unicode
	return '?'
}

// ConvertStringToASCII converts a string to ASCII using fallbacks
func ConvertStringToASCII(s string) string {
	var result []rune
	for _, r := range s {
		result = append(result, ToASCIIFallback(r))
	}
	return string(result)
}

// SanitizeForTerminal ensures a string is safe for terminal output
func SanitizeForTerminal(s string, capability UnicodeCapability) string {
	switch capability {
	case UnicodeNone:
		return ConvertStringToASCII(s)
	case UnicodeBasic:
		// Convert advanced Unicode to ASCII, keep basic Unicode and ASCII
		var result []rune
		for _, r := range s {
			if isAdvancedUnicode(r) {
				result = append(result, ToASCIIFallback(r))
			} else {
				result = append(result, r)
			}
		}
		return string(result)
	case UnicodeFull:
		return s // Return as-is for full Unicode support
	default:
		return ConvertStringToASCII(s)
	}
}

// isAdvancedUnicode determines if a character requires full Unicode support
func isAdvancedUnicode(r rune) bool {
	for _, advanced := range TestUnicodeCharacters.Full {
		if r == advanced {
			return true
		}
	}
	// Also check if it's a basic Unicode that should be converted in basic mode
	for _, basic := range TestUnicodeCharacters.Basic {
		if r == basic {
			return false // Basic Unicode characters are allowed in basic mode
		}
	}
	// If it's Unicode but not in our test sets, consider it advanced
	if r > 127 {
		return true
	}
	return false
}

// GetSpinnerChars returns appropriate spinner characters based on capability
func GetSpinnerChars(capability UnicodeCapability) []string {
	switch capability {
	case UnicodeNone:
		return []string{"|", "/", "-", "\\"}
	case UnicodeBasic:
		return []string{"|", "/", "-", "\\"}
	case UnicodeFull:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	default:
		return []string{"|", "/", "-", "\\"}
	}
}

// GetProgressChars returns appropriate progress bar characters
func GetProgressChars(capability UnicodeCapability) (filled, empty string) {
	switch capability {
	case UnicodeNone:
		return "#", "-"
	case UnicodeBasic:
		return "#", "-"
	case UnicodeFull:
		return "█", "░"
	default:
		return "#", "-"
	}
}

// GetBorderChars returns appropriate border characters for boxes
func GetBorderChars(capability UnicodeCapability) (topLeft, topRight, bottomLeft, bottomRight, horizontal, vertical string) {
	switch capability {
	case UnicodeNone:
		return "+", "+", "+", "+", "-", "|"
	case UnicodeBasic:
		return "+", "+", "+", "+", "-", "|"
	case UnicodeFull:
		return "┌", "┐", "└", "┘", "─", "│"
	default:
		return "+", "+", "+", "+", "-", "|"
	}
}

// TestUnicodeRendering attempts to render test characters to verify support
func TestUnicodeRendering() map[string]bool {
	results := make(map[string]bool)

	// Test basic Unicode characters
	basicTest := "█░•←→"
	results["basic"] = IsValidUTF8(basicTest) && !ContainsInvalidChars(basicTest)

	// Test advanced Unicode characters
	advancedTest := "⠋┌┐└┘✓✗"
	results["advanced"] = IsValidUTF8(advancedTest) && !ContainsInvalidChars(advancedTest)

	return results
}

// ContainsInvalidChars checks for characters that might cause rendering issues
func ContainsInvalidChars(s string) bool {
	for _, r := range s {
		// Check for control characters (except common ones like \n, \t)
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return true
		}
		// Check for non-printable characters
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
