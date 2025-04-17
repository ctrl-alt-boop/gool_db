package managers

import "github.com/jroimartin/gocui"

type StatusBarManager struct {
}

func (tlm *StatusBarManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView("status_bar", 0, maxY-6, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Editable = false
	}
	return nil
}
