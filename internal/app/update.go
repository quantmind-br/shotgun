package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/shotgun-cli/internal/screens/filetree"
	"github.com/user/shotgun-cli/internal/screens/input"
	"github.com/user/shotgun-cli/internal/screens/template"
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
	var cmds []tea.Cmd

	switch a.CurrentScreen {
	case FileTreeScreen:
		updatedModel, cmd := a.FileTree.Update(msg)
		a.FileTree = updatedModel.(filetree.FileTreeModel)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case TemplateScreen:
		updatedModel, cmd := a.Template.Update(msg)
		a.Template = updatedModel
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle template screen specific messages
		switch msg := msg.(type) {
		case template.TemplateSelectedMsg:
			a.SelectedTemplate = msg.Template
			a.SetCurrentScreen(TaskScreen)
		case template.BackToFileTreeMsg:
			a.SetCurrentScreen(FileTreeScreen)
		case template.RefreshTemplatesMsg:
			// Handle template refresh - would need template service
			// cmds = append(cmds, template.RefreshTemplatesCmd(a.templateService, a.ctx))
		}

	case TaskScreen:
		updatedModel, cmd := a.TaskInput.Update(msg)
		a.TaskInput = updatedModel
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Update AppState.TaskContent with current textarea content
		a.TaskContent = a.TaskInput.GetContent()

		// Handle task screen specific messages
		switch msg := msg.(type) {
		case input.TaskInputMsg:
			// Advance to rules screen when task input is complete
			a.SetCurrentScreen(RulesScreen)
		case input.BackToTemplateMsg:
			a.SetCurrentScreen(TemplateScreen)
		}

	case RulesScreen:
		updatedModel, cmd := a.RulesInput.Update(msg)
		a.RulesInput = updatedModel
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Update AppState.RulesContent with current textarea content
		a.RulesContent = a.RulesInput.GetContent()

		// Handle rules screen specific messages
		switch msg := msg.(type) {
		case input.RulesInputMsg:
			// Advance to confirmation screen when rules input is complete
			a.SetCurrentScreen(ConfirmScreen)
		case input.BackToTaskMsg:
			a.SetCurrentScreen(TaskScreen)
		case input.SkipRulesMsg:
			// Skip rules screen entirely and go to confirmation
			a.SetCurrentScreen(ConfirmScreen)
		}
	}

	return a, tea.Batch(cmds...)
}

// Screen-specific input handlers

func (a *AppState) handleTemplateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := a.Template.Update(msg)
	a.Template = updatedModel
	return a, cmd
}

func (a *AppState) handleTaskInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := a.TaskInput.Update(msg)
	a.TaskInput = updatedModel
	
	// Update AppState.TaskContent with current textarea content
	a.TaskContent = a.TaskInput.GetContent()
	
	return a, cmd
}

func (a *AppState) handleRulesInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := a.RulesInput.Update(msg)
	a.RulesInput = updatedModel
	
	// Update AppState.RulesContent with current textarea content
	a.RulesContent = a.RulesInput.GetContent()
	
	return a, cmd
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
