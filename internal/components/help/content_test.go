package help

import (
	"testing"
)

func TestHelpContentStructure(t *testing.T) {
	// Test that help content is defined for each screen
	screens := []ScreenType{
		FileTreeScreen,
		TemplateScreen,
		TaskScreen,
		RulesScreen,
		ConfirmScreen,
	}

	for _, screen := range screens {
		content := GetHelpContent(screen)
		
		if content == nil {
			t.Errorf("Expected help content for screen %v to be defined", screen)
			continue
		}

		// Check that each screen has at least some help items
		screenSpecificItems := 0
		globalItems := 0
		
		for _, item := range content {
			if item.Context == screen {
				screenSpecificItems++
			} else if item.Context == -1 {
				globalItems++
			}
		}

		// Each screen should have some specific items
		if screenSpecificItems == 0 && screen != GenerateScreen {
			t.Errorf("Screen %v has no screen-specific help items", screen)
		}
	}
}

func TestHelpItemKeys(t *testing.T) {
	// Test that help items have valid keys and descriptions
	content := GetHelpContent(FileTreeScreen)
	
	for _, item := range content {
		if item.Key == "" {
			t.Error("Found help item with empty key")
		}
		
		if item.Description == "" {
			t.Error("Found help item with empty description")
		}
		
		// Check that context is valid
		if item.Context < -1 || item.Context > GenerateScreen {
			t.Errorf("Invalid context %v for help item %s", item.Context, item.Key)
		}
	}
}

func TestGlobalHelpItems(t *testing.T) {
	// Test that global items (context -1) are consistent across screens
	globalKeys := make(map[string]string) // key -> description
	
	screens := []ScreenType{
		FileTreeScreen,
		TemplateScreen,
		TaskScreen,
		RulesScreen,
		ConfirmScreen,
	}
	
	for _, screen := range screens {
		content := GetHelpContent(screen)
		
		for _, item := range content {
			if item.Context == -1 { // Global item
				if desc, exists := globalKeys[item.Key]; exists {
					// Check that global items are consistent
					if desc != item.Description {
						t.Errorf("Inconsistent global help for key %s: '%s' vs '%s'",
							item.Key, desc, item.Description)
					}
				} else {
					globalKeys[item.Key] = item.Description
				}
			}
		}
	}
}