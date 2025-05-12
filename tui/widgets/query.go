package widgets

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	"github.com/ctrl-alt-boop/gooldb/tui/templates"
	"github.com/jesseduffield/gocui"
)

const QueryOptionsViewName string = "query_options"

type queryOptionsState struct {
	query *query.Statement
}

type QueryOptions struct {
	view *gocui.View
	// gui  *gocui.Gui

	state queryOptionsState
}

func (q *QueryOptions) Open(gui *gocui.Gui) {
	tmpl := templates.QueryOptions()
	tmpl.Execute(q.view, q.state.query)
}

func (q *QueryOptions) KeyEnter() error {
	//
	return nil
}

func (q *QueryOptions) Close(gui *gocui.Gui) {
	//
}

func (q *QueryOptions) CreatePopup(g *gocui.Gui, newQuery bool, selection string) error {
	dataTableView, err := g.View(DataAreaViewName)
	logger.Info("CreatePopup")
	if err != nil {
		return err
	}
	logger.Info("CreatePopup")
	logger.Info(dataTableView.Dimensions())
	popupView, err := g.SetView(Popup(QueryOptionsViewName, dataTableView.Dimensions))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		q.view = popupView
		q.state.query, _ = query.New(selection, query.Select)
		popupView.Frame = true
		popupView.Wrap = true
		popupView.Editable = false
		// popupView.HasLoader = true
		if newQuery {
			popupView.Title = fmt.Sprintf("New Query for %s", selection)

		} else {
			popupView.Title = fmt.Sprintf("SELECT FROM %s", selection)
		}
		q.CreateContent(newQuery, selection)
	}
	return nil
}

func (q *QueryOptions) CreateContent(newQuery bool, selection string) {
	tmpl := templates.QueryOptions()
	logger.Info(tmpl.Execute(logger, q.state.query))
	tmpl.Execute(q.view, q.state.query)
}
