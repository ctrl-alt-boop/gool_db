package widgets

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

func (keys keyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Help, keys.Quit}
}

func (keys keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.Help},
		{keys.Quit},
		{keys.Nav},
		{keys.Select},
		{keys.Back},
	}
}

type HelpFooter struct {
	help help.Model
	Keys keyMap
}

// Init implements tea.Model.
func (h *HelpFooter) Init() tea.Cmd {
	h.help.FullSeparator = " \u2502 "
	h.help.ShortSeparator = " \u2502 "
	return nil
}

func CreateHelp() *HelpFooter {
	return &HelpFooter{
		help: help.New(),
		Keys: KeyMap,
	}
}

func (h *HelpFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.UpdateSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Help):
			h.help.ShowAll = !h.help.ShowAll
		}
	}
	return h, nil
}

func (h *HelpFooter) UpdateSize(width int, _ int) {
	h.help.Width = width
}

func (h *HelpFooter) View() string {
	helpStyle := lipgloss.NewStyle().
		Width(h.help.Width).
		Height(1).
		Align(lipgloss.Left, lipgloss.Bottom).
		MarginLeft(1)
	return helpStyle.Render(h.help.View(h.Keys))
}
