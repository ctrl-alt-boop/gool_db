package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const HelpFooterName string = "help"

func getHelpText() string {
	return "\tCtrl-c quit\tF5 Fetch counts\t"
}

func InitHelpFooter(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView(HelpFooterName, 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = false
		view.Editable = false
		fmt.Fprint(view, getHelpText())
	}

	return nil
}
