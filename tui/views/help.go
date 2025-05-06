package views

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type helpMode int

const HelpFooterName string = "help"

const (
	DriversHelp helpMode = iota
	DatabasesHelp
	TablesHelp

	DataTableHelp
	TableCellHelp
)

type HelpFooterView struct {
	view *gocui.View

	currentMode helpMode
}

func (h *HelpFooterView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(HelpBar(maxX, maxY))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			logger.Panic(err)
		}
		view.Frame = false
		view.Editable = false
		h.view = view
		fmt.Fprint(view, getHelpText())
	}

	return nil
}

func (h *HelpFooterView) SetMode(mode helpMode) {
	h.currentMode = mode
}

func getHelpText() string {
	return "\tCtrl-c quit\tF5 Fetch counts\t"
}
