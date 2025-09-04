package template

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines key bindings for the template screen
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Select  key.Binding
	VimUp   key.Binding
	VimDown key.Binding
	Help    key.Binding
	Quit    key.Binding
}

// DefaultKeyMap returns the default key mappings for template screen
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select template"),
		),
		VimUp: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "move up (vim)"),
		),
		VimDown: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "move down (vim)"),
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

// ShortHelp returns key help summary for template screen
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Help, k.Quit}
}

// FullHelp returns extended key help for template screen
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.VimUp, k.VimDown},
		{k.Select, k.Help, k.Quit},
	}
}