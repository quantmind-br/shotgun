package utils

import (
	"testing"
	"unicode/utf8"
)

func TestUnicodeCapabilityString(t *testing.T) {
	tests := []struct {
		capability UnicodeCapability
		expected   string
	}{
		{UnicodeNone, "none"},
		{UnicodeBasic, "basic"},
		{UnicodeFull, "full"},
		{UnicodeCapability(99), "unknown"},
	}

	for _, test := range tests {
		result := test.capability.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s for capability %d", test.expected, result, test.capability)
		}
	}
}

func TestDetectUnicodeCapability(t *testing.T) {
	capability := DetectUnicodeCapability()

	// Should return one of the valid capabilities
	if capability != UnicodeNone && capability != UnicodeBasic && capability != UnicodeFull {
		t.Errorf("Invalid Unicode capability returned: %d", capability)
	}
}

func TestIsValidUTF8(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"hello world", true},
		{"⠋⠙⠹⠸", true},
		{"█▉▊▋", true},
		{string([]byte{0xff, 0xfe}), false}, // Invalid UTF-8
	}

	for _, test := range tests {
		result := IsValidUTF8(test.input)
		if result != test.expected {
			t.Errorf("IsValidUTF8(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestContainsUnicode(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", false},
		{"hello world", false},
		{"⠋⠙⠹⠸", true},
		{"hello █ world", true},
		{"123!@#", false},
		{"", false},
	}

	for _, test := range tests {
		result := ContainsUnicode(test.input)
		if result != test.expected {
			t.Errorf("ContainsUnicode(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToASCIIFallback(t *testing.T) {
	tests := []struct {
		input    rune
		expected rune
	}{
		{'a', 'a'}, // ASCII character should remain unchanged
		{'█', '#'}, // Unicode block should become #
		{'┌', '+'}, // Box drawing should become +
		{'⠋', '|'}, // Braille should become |
		{'✓', '+'}, // Check mark should become +
		{'€', '?'}, // Unknown Unicode should become ?
	}

	for _, test := range tests {
		result := ToASCIIFallback(test.input)
		if result != test.expected {
			t.Errorf("ToASCIIFallback('%c') = '%c', expected '%c'", test.input, result, test.expected)
		}
	}
}

func TestConvertStringToASCII(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"█▉▊▋", "####"},
		{"┌─┐", "+-+"},
		{"Loading... ⠋", "Loading... |"},
		{"Progress: █░", "Progress: #-"},
	}

	for _, test := range tests {
		result := ConvertStringToASCII(test.input)
		if result != test.expected {
			t.Errorf("ConvertStringToASCII(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestSanitizeForTerminal(t *testing.T) {
	input := "Loading... ⠋ Progress: █░"

	tests := []struct {
		capability UnicodeCapability
		contains   []string // Strings that should be present in output
	}{
		{UnicodeNone, []string{"Loading...", "|", "Progress:", "#", "-"}},
		{UnicodeBasic, []string{"Loading...", "|", "Progress:", "█", "░"}}, // Basic Unicode chars allowed
		{UnicodeFull, []string{"Loading...", "⠋", "Progress:", "█", "░"}},
	}

	for _, test := range tests {
		result := SanitizeForTerminal(input, test.capability)
		for _, expected := range test.contains {
			if !contains(result, expected) {
				t.Errorf("SanitizeForTerminal with %s capability should contain %q, got %q",
					test.capability.String(), expected, result)
			}
		}
	}
}

func TestGetSpinnerChars(t *testing.T) {
	tests := []struct {
		capability UnicodeCapability
		expected   []string
	}{
		{UnicodeNone, []string{"|", "/", "-", "\\"}},
		{UnicodeBasic, []string{"|", "/", "-", "\\"}},
		{UnicodeFull, []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}},
	}

	for _, test := range tests {
		result := GetSpinnerChars(test.capability)
		if len(result) != len(test.expected) {
			t.Errorf("GetSpinnerChars(%s) length = %d, expected %d",
				test.capability.String(), len(result), len(test.expected))
			continue
		}

		for i, expected := range test.expected {
			if result[i] != expected {
				t.Errorf("GetSpinnerChars(%s)[%d] = %q, expected %q",
					test.capability.String(), i, result[i], expected)
			}
		}
	}
}

func TestGetProgressChars(t *testing.T) {
	tests := []struct {
		capability     UnicodeCapability
		expectedFilled string
		expectedEmpty  string
	}{
		{UnicodeNone, "#", "-"},
		{UnicodeBasic, "#", "-"},
		{UnicodeFull, "█", "░"},
	}

	for _, test := range tests {
		filled, empty := GetProgressChars(test.capability)
		if filled != test.expectedFilled {
			t.Errorf("GetProgressChars(%s) filled = %q, expected %q",
				test.capability.String(), filled, test.expectedFilled)
		}
		if empty != test.expectedEmpty {
			t.Errorf("GetProgressChars(%s) empty = %q, expected %q",
				test.capability.String(), empty, test.expectedEmpty)
		}
	}
}

func TestGetBorderChars(t *testing.T) {
	tests := []struct {
		capability UnicodeCapability
		unicode    bool
	}{
		{UnicodeNone, false},
		{UnicodeBasic, false},
		{UnicodeFull, true},
	}

	for _, test := range tests {
		tl, tr, bl, br, h, v := GetBorderChars(test.capability)

		if test.unicode {
			// Should have Unicode box drawing characters
			if tl != "┌" || tr != "┐" || bl != "└" || br != "┘" || h != "─" || v != "│" {
				t.Errorf("GetBorderChars(%s) should return Unicode box drawing characters",
					test.capability.String())
			}
		} else {
			// Should have ASCII fallbacks
			if tl != "+" || tr != "+" || bl != "+" || br != "+" || h != "-" || v != "|" {
				t.Errorf("GetBorderChars(%s) should return ASCII fallback characters",
					test.capability.String())
			}
		}
	}
}

func TestTestUnicodeRendering(t *testing.T) {
	results := TestUnicodeRendering()

	// Should have results for basic and advanced
	if _, exists := results["basic"]; !exists {
		t.Error("TestUnicodeRendering should test basic Unicode")
	}

	if _, exists := results["advanced"]; !exists {
		t.Error("TestUnicodeRendering should test advanced Unicode")
	}

	// Results should be boolean
	for key, value := range results {
		if _, isBool := interface{}(value).(bool); !isBool {
			t.Errorf("TestUnicodeRendering result for %s should be boolean", key)
		}
	}
}

func TestContainsInvalidChars(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello world", false},
		{"hello\tworld\n", false}, // Tab and newline are allowed
		{"hello\x00world", true},  // Null character is invalid
		{"hello\x7fworld", true},  // DEL character is invalid
		{"⠋⠙⠹⠸", false},           // Valid Unicode
		{"", false},               // Empty string is valid
	}

	for _, test := range tests {
		result := ContainsInvalidChars(test.input)
		if result != test.expected {
			t.Errorf("ContainsInvalidChars(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// Test that all fallback mappings are valid
func TestFallbackMappings(t *testing.T) {
	for unicode, ascii := range ASCIIFallbacks {
		// Unicode character should be valid
		if !utf8.ValidRune(unicode) {
			t.Errorf("Invalid Unicode rune in fallback: %U", unicode)
		}

		// ASCII fallback should be printable ASCII
		if ascii > 127 || ascii < 32 {
			if ascii != '\n' && ascii != '\t' && ascii != '\r' {
				t.Errorf("Invalid ASCII fallback for %U: %U", unicode, ascii)
			}
		}
	}
}

// Benchmark tests
func BenchmarkToASCIIFallback(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToASCIIFallback('█')
		ToASCIIFallback('a')
		ToASCIIFallback('⠋')
	}
}

func BenchmarkConvertStringToASCII(b *testing.B) {
	testString := "Loading... ⠋ Progress: █▉▊▋▌▍▎▏░"
	for i := 0; i < b.N; i++ {
		ConvertStringToASCII(testString)
	}
}

func BenchmarkSanitizeForTerminal(b *testing.B) {
	testString := "Loading... ⠋ Progress: █▉▊▋▌▍▎▏░"
	for i := 0; i < b.N; i++ {
		SanitizeForTerminal(testString, UnicodeFull)
		SanitizeForTerminal(testString, UnicodeNone)
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
