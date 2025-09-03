package filetree

import (
	"strings"
	"testing"
)

func TestDefaultKeyMap(t *testing.T) {
	keyMap := DefaultKeyMap()

	// Test up key
	if !keyMap.Up.Enabled() {
		t.Error("Expected Up key to be enabled")
	}

	// Test that up key bindings include expected keys
	upKeys := keyMap.Up.Keys()
	hasUp := false
	for _, k := range upKeys {
		if k == "up" {
			hasUp = true
		}
	}
	if !hasUp {
		t.Error("Expected Up key binding to include 'up'")
	}

	// Test vim up key separately
	if !keyMap.VimUp.Enabled() {
		t.Error("Expected VimUp key to be enabled")
	}
	
	vimUpKeys := keyMap.VimUp.Keys()
	hasK := false
	for _, k := range vimUpKeys {
		if k == "k" {
			hasK = true
		}
	}
	if !hasK {
		t.Error("Expected VimUp key binding to include 'k'")
	}

	// Test down key
	if !keyMap.Down.Enabled() {
		t.Error("Expected Down key to be enabled")
	}

	// Test right key (expand)
	if !keyMap.Right.Enabled() {
		t.Error("Expected Right key to be enabled")
	}

	// Test left key (collapse)
	if !keyMap.Left.Enabled() {
		t.Error("Expected Left key to be enabled")
	}

	// Test toggle key
	if !keyMap.Toggle.Enabled() {
		t.Error("Expected Toggle key to be enabled")
	}

	// Test quit key
	if !keyMap.Quit.Enabled() {
		t.Error("Expected Quit key to be enabled")
	}
}

func TestShortHelp(t *testing.T) {
	keyMap := DefaultKeyMap()
	help := keyMap.ShortHelp()

	// Should return key bindings
	if len(help) == 0 {
		t.Error("Expected ShortHelp to return key bindings")
	}

	// Should contain 5 bindings as defined
	if len(help) != 5 {
		t.Errorf("Expected ShortHelp to return 5 key bindings, got %d", len(help))
	}

	// Should contain some expected key bindings
	found := false
	for _, binding := range help {
		if binding.Enabled() {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected at least one enabled key binding in ShortHelp")
	}
}

func TestFullHelp(t *testing.T) {
	keyMap := DefaultKeyMap()
	help := keyMap.FullHelp()

	// Should return key bindings organized in sections
	if len(help) == 0 {
		t.Error("Expected FullHelp to return key binding sections")
	}

	// Should have 3 sections as defined in keys.go
	if len(help) != 3 {
		t.Errorf("Expected FullHelp to return 3 sections, got %d", len(help))
	}

	// Should have at least one section with key bindings
	found := false
	for _, section := range help {
		if len(section) > 0 {
			for _, binding := range section {
				if binding.Enabled() {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("Expected at least one enabled key binding in FullHelp")
	}
}

func TestKeyMapHelp(t *testing.T) {
	keyMap := DefaultKeyMap()
	
	// Test that we can get help string from a key
	upHelp := keyMap.Up.Help()
	if upHelp.Key == "" {
		t.Error("Expected Up key to have help text")
	}
	
	// Test space key help
	spaceHelp := keyMap.Toggle.Help()
	if spaceHelp.Key == "" {
		t.Error("Expected Toggle key to have help text")
	}
	if !strings.Contains(strings.ToLower(spaceHelp.Desc), "toggle") && 
	   !strings.Contains(strings.ToLower(spaceHelp.Desc), "select") {
		t.Error("Expected Toggle key help to mention toggle or select")
	}
}