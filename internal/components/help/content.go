package help

// GetHelpContent returns context-sensitive help for each screen
func GetHelpContent(screen ScreenType) []HelpItem {
	switch screen {
	case FileTreeScreen:
		return []HelpItem{
			{"Space", "Select/deselect file", FileTreeScreen},
			{"Enter", "Toggle directory", FileTreeScreen},
			{"↑/↓", "Navigate file list", FileTreeScreen},
			{"j/k", "Navigate file list (vim)", FileTreeScreen},
		}
	case TemplateScreen:
		return []HelpItem{
			{"Enter", "Select template", TemplateScreen},
			{"↑/↓", "Navigate templates", TemplateScreen},
			{"j/k", "Navigate templates (vim)", TemplateScreen},
		}
	case TaskScreen:
		return []HelpItem{
			{"Tab", "Navigate between fields", TaskScreen},
			{"Ctrl+Enter", "Advance to next screen", TaskScreen},
		}
	case RulesScreen:
		return []HelpItem{
			{"Tab", "Navigate between fields", RulesScreen},
			{"Ctrl+Enter", "Advance to next screen", RulesScreen},
		}
	case ConfirmScreen:
		return []HelpItem{
			{"Enter", "Edit selected section", ConfirmScreen},
			{"↑/↓", "Navigate sections", ConfirmScreen},
		}
	case GenerateScreen:
		return []HelpItem{
			{"q", "Quit after generation", GenerateScreen},
		}
	default:
		return []HelpItem{}
	}
}