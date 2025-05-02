package views

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/jesseduffield/gocui"
)

const DataTableViewName string = "data"

type dataViewState struct {
	currentColumnIndex int
	currentRowIndex    int
	columnWidths       []int
	table              *gooldb.DataTable
}

type DataTableView struct {
	view *gocui.View
	gui  *gocui.Gui
	// gool  *gooldb.GoolDb
	state dataViewState
}

func (d *DataTableView) RightColumn() int {
	if d.state.currentColumnIndex >= len(d.state.columnWidths)-1 {
		return 0
	}
	moveWidth := d.state.columnWidths[d.state.currentColumnIndex]
	d.state.currentColumnIndex++
	return moveWidth
}

func (d *DataTableView) LeftColumn() int {
	if d.state.currentColumnIndex <= 0 {
		return 0
	}
	moveWidth := d.state.columnWidths[d.state.currentColumnIndex-1]
	d.state.currentColumnIndex--
	return moveWidth
}

func (d *DataTableView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(DataView(maxX, maxY))
	if err != nil {
		if err != gocui.ErrUnknownView {
			//panic(err)
		}
		view.Frame = true
		view.Editable = false
		d.view = view
		d.gui = g
		d.state = dataViewState{}
	}

	return nil
}

func (d *DataTableView) OnEnterPressed(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {

		return nil
	}
}

// selected string, table *gooldb.DataTable
func (d *DataTableView) OnTableSet(eventArgs any, err error) {
	logger.Info("OnTableSet: ", eventArgs, err)
	args, ok := eventArgs.(gooldb.TableSetEvent)
	if !ok {
		return
	}
	table := args.Table
	if err != nil {
		return
	}

	d.state.table = table
	dataView := d.view

	d.gui.Update(func(g *gocui.Gui) error {
		dataView.Clear()
		formatedRows, columnWidths := table.GetFormatedRows()
		formatedHeader := table.GetFormatedTitle(columnWidths)
		dataView.Title = formatedHeader

		d.state.columnWidths = columnWidths
		d.state.currentColumnIndex = 0
		d.state.currentRowIndex = 0

		logger.Info("header: ", formatedHeader)
		logger.Info("rows: ", strings.Join(formatedRows, "\n"))

		// fmt.Fprint(dataView, formatedHeader)
		fmt.Fprint(dataView, strings.Join(formatedRows, "\n"))
		return nil
	})
}
