package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/shotgun-cli/internal/screens/filetree"
)

// Init implements tea.Model interface
func (a *AppState) Init() tea.Cmd {
	// Initialize the first screen (FileTree)
	return tea.Batch(
		a.InitScreenCmd(),
		a.FileTree.Init(),
	)
}

// Update implements tea.Model interface
func (a *AppState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update window size for all components
		a.UpdateWindowSize(msg)
		return a, nil

	case tea.KeyMsg:
		// Handle help and exit dialogs first
		if a.ShowingHelp {
			return a.handleHelpDialog(msg)
		}
		if a.ShowingExit {
			return a.handleExitDialog(msg)
		}

		// Check for global keys first
		if IsGlobalKey(msg.String()) {
			return a.GlobalKeyHandler(msg)
		}

		// Let current screen handle the key
		return a.handleScreenInput(msg)

	default:
		// Pass message to current screen
		return a.handleScreenMessage(msg)
	}

	return a, tea.Batch(cmds...)
}

// handleScreenInput routes keyboard input to current screen
func (a *AppState) handleScreenInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch a.CurrentScreen {
	case FileTreeScreen:
		updatedModel, cmd := a.FileTree.Update(msg)
		a.FileTree = updatedModel.(filetree.FileTreeModel)
		return a, cmd

	case TemplateScreen:
		return a.handleTemplateInput(msg)

	case TaskScreen:
		return a.handleTaskInput(msg)

	case RulesScreen:
		return a.handleRulesInput(msg)

	case ConfirmScreen:
		return a.handleConfirmationInput(msg)

	default:
		return a, nil
	}
}

// handleScreenMessage routes other messages to current screen
func (a *AppState) handleScreenMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch a.CurrentScreen {
	case FileTreeScreen:
		updatedModel, cmd := a.FileTree.Update(msg)
		a.FileTree = updatedModel.(filetree.FileTreeModel)
		return a, cmd

	default:
		return a, nil
	}
}

// Screen-specific input handlers

func (a *AppState) handleTemplateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if a.Template.cursor > 0 {
			a.Template.cursor--
		}
	case "down", "j":
		if a.Template.cursor < len(a.Template.items)-1 {
			a.Template.cursor++
		}
	case "enter", " ":
		if len(a.Template.items) > 0 && a.Template.cursor < len(a.Template.items) {
			a.Template.selected = a.Template.items[a.Template.cursor]
			a.SelectedTemplate = a.Template.selected
		}
	}
	return a, nil
}

func (a *AppState) handleTaskInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+a":
		a.TaskInput.cursor = 0
	case "ctrl+e":
		a.TaskInput.cursor = len(a.TaskInput.content)
	case "left":
		if a.TaskInput.cursor > 0 {
			a.TaskInput.cursor--
		}
	case "right":
		if a.TaskInput.cursor < len(a.TaskInput.content) {
			a.TaskInput.cursor++
		}
	case "backspace":
		if a.TaskInput.cursor > 0 && len(a.TaskInput.content) > 0 {
			a.TaskInput.content = a.TaskInput.content[:a.TaskInput.cursor-1] + a.TaskInput.content[a.TaskInput.cursor:]
			a.TaskInput.cursor--
			a.TaskContent = a.TaskInput.content
		}
	case "delete":
		if a.TaskInput.cursor < len(a.TaskInput.content) {
			a.TaskInput.content = a.TaskInput.content[:a.TaskInput.cursor] + a.TaskInput.content[a.TaskInput.cursor+1:]
			a.TaskContent = a.TaskInput.content
		}
	default:
		// Handle regular character input
		if len(msg.Runes) > 0 {
			char := string(msg.Runes[0])
			if len(char) == 1 && char >= " " { // Printable character
				a.TaskInput.content = a.TaskInput.content[:a.TaskInput.cursor] + char + a.TaskInput.content[a.TaskInput.cursor:]
				a.TaskInput.cursor++
				a.TaskContent = a.TaskInput.content
			}
		}
	}
	return a, nil
}

func (a *AppState) handleRulesInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+a":
		a.RulesInput.cursor = 0
	case "ctrl+e":
		a.RulesInput.cursor = len(a.RulesInput.content)
	case "left":
		if a.RulesInput.cursor > 0 {
			a.RulesInput.cursor--
		}
	case "right":
		if a.RulesInput.cursor < len(a.RulesInput.content) {
			a.RulesInput.cursor++
		}
	case "backspace":
		if a.RulesInput.cursor > 0 && len(a.RulesInput.content) > 0 {
			a.RulesInput.content = a.RulesInput.content[:a.RulesInput.cursor-1] + a.RulesInput.content[a.RulesInput.cursor:]
			a.RulesInput.cursor--
			a.RulesContent = a.RulesInput.content
		}
	case "delete":
		if a.RulesInput.cursor < len(a.RulesInput.content) {
			a.RulesInput.content = a.RulesInput.content[:a.RulesInput.cursor] + a.RulesInput.content[a.RulesInput.cursor+1:]
			a.RulesContent = a.RulesInput.content
		}
	default:
		// Handle regular character input
		if len(msg.Runes) > 0 {
			char := string(msg.Runes[0])
			if len(char) == 1 && char >= " " { // Printable character
				a.RulesInput.content = a.RulesInput.content[:a.RulesInput.cursor] + char + a.RulesInput.content[a.RulesInput.cursor:]
				a.RulesInput.cursor++
				a.RulesContent = a.RulesInput.content
			}
		}
	}
	return a, nil
}

func (a *AppState) handleConfirmationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		// Scroll up if supported
	case "down", "j":
		// Scroll down if supported
	}
	return a, nil
}

// Dialog handlers

func (a *AppState) handleHelpDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "f1", "q":
		a.ShowingHelp = false
		a.HelpContent = ""
		return a, nil
	}
	return a, nil
}

func (a *AppState) handleExitDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return a, tea.Quit
	case "n", "N", "esc":
		a.ShowingExit = false
		return a, nil
	}
	return a, nil
}

