package managers

import "github.com/jroimartin/gocui"

type TableListManager struct {
}

func (tlm *TableListManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView("table_list", 0, 0, maxX/5, maxY-7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "Tables"
		view.Frame = true
		view.Editable = false
	}
	return nil
}
