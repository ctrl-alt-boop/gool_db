package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jroimartin/gocui"
)

type Tui struct {
	*gocui.Gui
	GoolDb *gooldb.GoolDb

	SidePanel          string
	MainView           string
	CurrentDataColumns []string
}

func Create(goolDb *gooldb.GoolDb) *Tui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}

	g.SetManagerFunc(func(g *gocui.Gui) error {

		views.InitHelpFooter(g)

		views.InitCommandBar(g)

		views.InitSidePanel(g)
		views.SetSidePanelMode(views.DriverList)

		g.SetCurrentView(views.SidePanelViewName)

		views.InitDataView(g)

		return nil
	})

	tui := &Tui{
		Gui:                g,
		GoolDb:             goolDb,
		CurrentDataColumns: nil,
	}
	tui.setKeybinds()
	return tui
}

func (tui *Tui) setKeybinds() {
	if err := tui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, MoveCursorUp); err != nil {
		log.Panicln(err)
	}
	if err := tui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, MoveCursorDown); err != nil {
		log.Panicln(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, OnEnterPressed(tui)); err != nil {
		log.Panicln(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, CycleCurrentView(tui)); err != nil {
		log.Panicln(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyF5, gocui.ModNone, OnF5Pressed(tui)); err != nil {
		log.Panicln(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}
}

func (tui *Tui) onEnterPressed(_ *gocui.Gui, v *gocui.View) error {
	selection, err := v.Word(v.Cursor())
	if err != nil {
		return err
	}
	for _, dataColumn := range tui.CurrentDataColumns {
		tui.DeleteView(dataColumn)
	}

	table := tui.GoolDb.FetchTable(selection)
	// names := table.ColumnNames()
	// log.Println(strings.Join(names, ", "))

	dataColumns := make([]string, 0)
	columnPadding := 2

	mainView, err := tui.View(views.DataTableViewName)
	if err != nil {
		return err
	}
	maxX, _ := tui.Size()
	_, mainSizeY := mainView.Size()
	currentX := maxX / 6
	for i, col := range table.Columns() {
		data, width := table.GetColumnRows(i)
		if width <= len(col.Name) {
			width = columnPadding + len(col.Name) + columnPadding
		}
		colView, err := tui.SetView(col.Name, currentX, 0, currentX+width+columnPadding, mainSizeY)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			colView.Title = col.Name
			colView.Frame = true
			colView.Editable = false
			fmt.Fprint(colView, strings.Join(data, "\n"))
		}
		currentX += width + columnPadding
		dataColumns = append(dataColumns, col.Name)
	}
	tui.CurrentDataColumns = dataColumns
	return nil
}

// func (tui *Tui) SetDatabaseContext(context *database.DatabaseContext) {
// 	tui.DatabaseContext = context

// 	tui.onSetDatabaseContext()
// }

func (tui *Tui) onSetDatabaseContext() {
	// dbName := tui.DatabaseContext.FetchDatabaseName()
	// tableNames := tui.DatabaseContext.FetchTableNames()

	// tui.Update(func(g *gocui.Gui) error {
	// 	view, err := g.View("table_list")
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	view.Clear()

	// 	view.Title = dbName
	// 	view.Highlight = true
	// 	view.SelFgColor = gocui.AttrReverse
	// 	view.SelBgColor = gocui.AttrReverse

	// 	fmt.Fprint(view, strings.Join(tableNames, "\n"))

	// 	g.SetCurrentView("table_list")
	// 	if err := view.SetCursor(0, 0); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
}
