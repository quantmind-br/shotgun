package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/core/builder"
	"github.com/diogopedro/shotgun/internal/screens/confirm"
	"github.com/diogopedro/shotgun/internal/screens/filetree"
	"github.com/diogopedro/shotgun/internal/screens/generate"
	"github.com/diogopedro/shotgun/internal/screens/input"
	"github.com/diogopedro/shotgun/internal/screens/template"
)

// Init implements tea.Model interface
func (a *AppState) Init() tea.Cmd {
	// Initialize the first screen (FileTree)
	return tea.Batch(
		a.InitScreenCmd(),
		a.FileTree.Init(),
		// Auto-start file scanning for the initial screen
		a.FileTree.StartScanning(),
		a.FileTree.LoadFromScanner(a.ctx, "."),
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

		// Check for global keys (robust to platform-specific key types)
		if IsGlobalKey(normalizeKey(msg)) || isFunctionKeyMsg(msg) {
			return a.GlobalKeyHandler(msg)
		}

		// Screen-specific shortcuts (non-text screens can also use plain Enter)
		// FileTree: Alt+C (and Enter on some terminals) advances if files are selected
		if a.CurrentScreen == FileTreeScreen && (normalizeKey(msg) == "alt+c" || normalizeKey(msg) == "enter") {
			if a.canGoToNextScreen() {
				return a.goToNextScreen()
			}
			return a, nil
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

	case GenerateScreen:
		return a.handleGenerationInput(msg)

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
			// Refresh templates via service with discovery spinner
			cmds = append(cmds,
				a.Template.StartDiscovery(),
				template.RefreshTemplatesCmd(a.templateService, a.ctx),
			)
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
		switch msg.(type) {
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
		switch msg.(type) {
		case input.RulesInputMsg:
			// Advance to confirmation screen when rules input is complete
			a.SetCurrentScreen(ConfirmScreen)
		case input.BackToTaskMsg:
			a.SetCurrentScreen(TaskScreen)
		case input.SkipRulesMsg:
			// Skip rules screen entirely and go to confirmation
			a.SetCurrentScreen(ConfirmScreen)
		}

	case ConfirmScreen:
		updatedModel, cmd := a.Confirmation.Update(msg)
		a.Confirmation = updatedModel
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle confirmation screen specific messages
		switch msg.(type) {
		case confirm.ConfirmGenerationMsg:
			// Confirm and proceed to generation - removed size check to allow generation
			if !a.Confirmation.IsCalculating() {
				return a, a.generatePrompt()
			}
		case confirm.NavigateToRulesMsg:
			a.SetCurrentScreen(RulesScreen)
		}

	case GenerateScreen:
		updatedModel, cmd := a.Generation.Update(msg)
		a.Generation = updatedModel
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle generation screen specific messages
		switch msg.(type) {
		case generate.NavigateBackMsg:
			a.SetCurrentScreen(ConfirmScreen)
		case generate.NavigateToFileTreeMsg:
			a.SetCurrentScreen(FileTreeScreen)
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
	// Initialize confirmation screen with current data if needed
	if !a.Confirmation.IsReady() {
		a.Confirmation.SetData(
			a.SelectedTemplate,
			a.SelectedFiles,
			a.TaskContent,
			a.RulesContent,
		)

		// Trigger size calculation and filename generation
		return a, tea.Batch(
			confirm.GenerateFilenameCmd(),
			a.startSizeCalculation(),
		)
	}

	// Handle confirmation screen updates
	updatedModel, cmd := a.Confirmation.Update(msg)
	a.Confirmation = updatedModel

	return a, cmd
}

func (a *AppState) handleGenerationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := a.Generation.Update(msg)
	a.Generation = updatedModel
	return a, cmd
}

// startSizeCalculation initiates size calculation for the confirmation screen
func (a *AppState) startSizeCalculation() tea.Cmd {
	if a.SelectedTemplate == nil {
		return nil
	}

	// Create estimator (would need to be imported from builder package)
	// For now, return a placeholder command
	return func() tea.Msg {
		return confirm.SizeCalculationStartMsg{}
	}
}

// generatePrompt triggers the prompt generation process
func (a *AppState) generatePrompt() tea.Cmd {
	// Switch to generation screen
	a.SetCurrentScreen(GenerateScreen)

	// Create generation configuration from app state
	config := builder.GenerationConfig{
		Template:      a.SelectedTemplate,
		Variables:     make(map[string]string),
		SelectedFiles: a.SelectedFiles,
		TaskContent:   a.TaskContent,
		RulesContent:  a.RulesContent,
		OutputPath:    "", // Use current directory
	}

	// Start generation process
	return func() tea.Msg {
		return generate.StartGenerationMsg{
			Config: config,
		}
	}
}

// Dialog handlers

func (a *AppState) handleHelpDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Update help component first
	var cmd tea.Cmd
	a.Help, cmd = a.Help.Update(msg)

	// Update our state based on help component
	a.ShowingHelp = a.Help.IsVisible()

	return a, cmd
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
