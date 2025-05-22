package config

import "github.com/charmbracelet/bubbles/key"

const (
	SymbolArrowUp    = "↑"
	SymbolArrowDown  = "↓"
	SymbolArrowLeft  = "←"
	SymbolArrowRight = "→"
	SymbolEnter      = "⏎"
	SymbolBackspace  = "⌫"
	SymbolSpace      = "␣"
	SymbolTab        = "⇥"
	SymbolEscape     = "⎋"
)

const (
	ArrowNav = SymbolArrowLeft + SymbolArrowDown + SymbolArrowUp + SymbolArrowRight
	VimNav   = "hjkl"
)

const (
	ArrowUp    = "up"
	ArrowDown  = "down"
	ArrowLeft  = "left"
	ArrowRight = "right"
)

type KeyMap struct {
	Nav       key.Binding // Used as combined navigation keys for help
	CycleView key.Binding

	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Back   key.Binding

	Help key.Binding
	Quit key.Binding
}

func (keys KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Help, keys.Quit}
}

func (keys KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.Help},
		{keys.Quit},
		{keys.CycleView},
		{keys.Nav},
		{keys.Select},
		{keys.Back},
	}
}

var Keys = createKeyMap()

func createKeyMap() KeyMap {
	return KeyMap{
		Nav: key.NewBinding(
			key.WithKeys("nil"),
			key.WithHelp(ArrowNav+"/"+VimNav, "navigate"),
		),
		CycleView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(SymbolTab, "cycle view"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp(SymbolArrowLeft, "left"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp(SymbolArrowDown, "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp(SymbolArrowUp, "up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp(SymbolArrowRight, "right"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(SymbolEnter, "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp(SymbolEscape+"/"+SymbolBackspace, "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
