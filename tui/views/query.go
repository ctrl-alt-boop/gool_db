package views

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	"github.com/jesseduffield/gocui"
)

const QueryOptionsViewName = "query_options"

type QueryOptionsView struct {
	view *gocui.View
	gui  *gocui.Gui

	table          string
	currentOptions query.Statement
}

func (q *QueryOptionsView) Layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	view, err := gui.SetView(QueryOptionsViewName, maxX/2-30, maxY/2, maxX/2+30, maxY/2+2, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		view.Title = fmt.Sprintf("SELECT * FROM %s WHERE ...", q.table)
		view.Highlight = true
		view.SelBgColor = gocui.AttrReverse
		view.SelFgColor = gocui.AttrReverse

		q.view = view
		q.gui = gui
	}
	return nil
}

func (q *QueryOptionsView) OnEnterPressed() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {

		return nil
	}
}
