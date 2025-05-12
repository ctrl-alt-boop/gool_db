package widgets

import (
	"github.com/jesseduffield/gocui"
)

type helpMode int

const HelpFooterName string = "help"

var Help *HelpFooter = CreateHelpFooter()

const (
	DriversHelp helpMode = iota
	DatabasesHelp
	TablesHelp

	CommandBarHelp
	DataTableHelp
	TableCellHelp
	QueryOptionsHelp
)

func modeFromName(name string) helpMode {
	switch name {
	case DriverList.String():
		return DriversHelp
	case DatabaseList.String():
		return DatabasesHelp
	case TableList.String():
		return TablesHelp
	case CommandBarViewName:
		return CommandBarHelp
	case DataAreaViewName:
		return DataTableHelp
	case TableCellViewName:
		return TableCellHelp
	case QueryOptionsViewName:
		return QueryOptionsHelp
	default:
		return DriversHelp
	}
}

func (m helpMode) Text() string {
	helpText := "\tCtrl-c quit\t"
	switch m {
	case DriversHelp:
		helpText += "Drivers: <enter> Select"
	case DatabasesHelp:
		helpText += "Databases: <enter> Select"
	case TablesHelp:
		helpText += "Tables: <esc> Back \t<enter> Select\tc Toggle Counts"
	case DataTableHelp:
		helpText += "Data table: f Query options\t <enter> Open table cell"
	case TableCellHelp:
		helpText += "Table cell: <esc> Close"
	case QueryOptionsHelp:
		helpText += "Query options: <esc> Close"
	case CommandBarHelp:
		helpText += "Command bar: <esc> Close\t <enter> Execute"
	default:
	}
	return helpText
}

type HelpFooter struct {
	view *gocui.View

	currentMode helpMode
}

func CreateHelpFooter() *HelpFooter {
	return &HelpFooter{
		currentMode: DriversHelp,
	}
}

func (h *HelpFooter) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(HelpBarLayout(maxX, maxY))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			logger.Panic(err)
		}
		view.Frame = false
		view.Editable = false
		h.view = view

	}
	h.view.SetContent(h.currentMode.Text())
	return nil
}

func (h *HelpFooter) SetCurrentView(name string) {
	h.currentMode = modeFromName(name)
}
