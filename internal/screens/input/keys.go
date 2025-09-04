package input

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines key bindings for input screens (task and rules)
type KeyMap struct {
	Tab        key.Binding
	ShiftTab   key.Binding
	Submit     key.Binding
	Help       key.Binding
	Quit       key.Binding
}

// DefaultKeyMap returns the default key mappings for input screens
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "navigate fields"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "navigate backwards"),
		),
		Submit: key.NewBinding(
			key.WithKeys("ctrl+enter"),
			key.WithHelp("ctrl+enter", "advance to next screen"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns key help summary for input screens
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Submit, k.Help, k.Quit}
}

// FullHelp returns extended key help for input screens
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.ShiftTab},
		{k.Submit, k.Help, k.Quit},
	}
}