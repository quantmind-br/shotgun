package filetree

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the file tree
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Toggle   key.Binding
	VimUp    key.Binding
	VimDown  key.Binding
	VimLeft  key.Binding
	VimRight key.Binding
	Help     key.Binding
	Quit     key.Binding
}

// DefaultKeyMap returns the default keybindings
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
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "collapse directory"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "expand directory"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle selection"),
		),
		VimUp: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "move up (vim)"),
		),
		VimDown: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "move down (vim)"),
		),
		VimLeft: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "collapse directory (vim)"),
		),
		VimRight: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "expand directory (vim)"),
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

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Toggle, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.VimUp, k.VimDown},
		{k.Left, k.Right, k.VimLeft, k.VimRight},
		{k.Toggle, k.Help, k.Quit},
	}
}
