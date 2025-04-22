package views

import "github.com/jroimartin/gocui"

const DataTableViewName string = "data"

type DataTableView struct {
	*gocui.View
}

var DataTable *DataTableView

func InitDataView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView(DataTableViewName, maxX/6, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Editable = false
	}
	return nil
}

func (d *DataTableView) OnEnter(column string) error {
	return nil
}
