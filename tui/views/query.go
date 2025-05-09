package views

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	"github.com/jesseduffield/gocui"
)

const QueryOptionsViewName string = "query_options"

type QueryOptionsView struct {
	view *gocui.View
	gui  *gocui.Gui

	table          string
	currentOptions query.Statement
}

func (q *QueryOptionsView) InitQueryOptions(view *gocui.View, newQuery bool, selection string) {
	panic("unimplemented")
}

func (q *QueryOptionsView) Open(gui *gocui.Gui) {
	panic("unimplemented")
}

func (q *QueryOptionsView) KeyEnter() error {
	panic("unimplemented")
}

func (q *QueryOptionsView) Close(gui *gocui.Gui) {
	panic("unimplemented")
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

func (q *QueryOptionsView) CreatePopup(newQuery bool, selection string) func(_ *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		dataTableView, err := g.View(DataTableViewName)
		if err != nil {
			return err
		}
		popupView, err := g.SetView(Popup(QueryOptionsViewName, dataTableView.Dimensions))
		if err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}
			popupView.Frame = true
			popupView.Wrap = true
			popupView.Editable = false
			if newQuery {
				popupView.Title = fmt.Sprintf("New Query for %s", selection)
			} else {
				popupView.Title = fmt.Sprintf("SELECT * FROM %s WHERE ...", selection)
			}
			popupView.SetContent(q.CreateContent(newQuery, selection))
		}
		return nil
	}
}

func (q *QueryOptionsView) CreateContent(newQuery bool, selection string) string {
	panic("unimplemented")
}
