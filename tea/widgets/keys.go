package widgets

import "github.com/charmbracelet/bubbles/key"

const (
	ArrowUp    = "up"
	ArrowDown  = "down"
	ArrowLeft  = "left"
	ArrowRight = "right"
)

type keyMap struct {
	Nav key.Binding // Used as combined navigation keys for help

	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Back   key.Binding

	Help key.Binding
	Quit key.Binding
}

var KeyMap = createKeyMap()

func createKeyMap() keyMap {
	return keyMap{
		Nav: key.NewBinding(
			key.WithKeys("nil"),
			key.WithHelp(ArrowNav+"/"+VimNav, "navigate"),
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
