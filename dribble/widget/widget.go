package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("widgets.log")

const (
	PanelWidthRatio       int = 6
	BorderThickness       int = 1
	BorderThicknessDouble int = 2
)

type Kind int

const (
	KindPanel Kind = iota
	KindWorkspace
	KindHelp
	KindCommandLine

	KindPopupHandler
	KindQueryOptions
	KindTableCell
)

type RequestFocus Kind

func ChangeFocus(id Kind) tea.Cmd {
	return func() tea.Msg {
		return RequestFocus(id)
	}
}

type (
	WidgetNames struct {
		Command   string
		Workspace string
		Panel     string
		Help      string

		Popups PopupNames
	}

	PopupNames struct {
		Handler      string
		QueryOptions string
		TableCell    string
	}
)

var Widgets = WidgetNames{
	Command:   "Command",
	Panel:     "Panel",
	Workspace: "Workspace",

	Help: "Help",

	Popups: PopupNames{
		Handler:      "Popups",
		QueryOptions: "Popup_QueryOptions",
		TableCell:    "Popup_TableCell",
	},
}
