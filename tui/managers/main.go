package managers

import "github.com/jroimartin/gocui"

type MainManager struct {
}

func (tlm *MainManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView("main", maxX/5, 0, maxX-1, maxY-7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "gool_db"
		view.Frame = true
		view.Editable = false
	}
	return nil
}
