package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/database"
	"github.com/ctrl-alt-boop/gooldb/tui/managers"
	"github.com/jroimartin/gocui"
)

type Tui struct {
	*gocui.Gui

	DatabaseContext *database.DatabaseContext
	CurrentData     *database.DataTable
}

func Create() *Tui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}

	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if view, err := g.SetView("help", 0, maxY-2, maxX-1, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			view.Frame = false
			view.Editable = false
			fmt.Fprint(view, managers.App.GetHelpText())
		}

		if view, err := g.SetView("table_list", 0, 0, maxX/6, maxY-4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			view.Title = "Tables"
			view.Frame = true
			view.Editable = false
		}
		if view, err := g.SetView("main", maxX/6, 0, maxX-1, maxY-4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			view.Title = "gool_db"
			view.Frame = true
			view.Editable = false
		}
		if view, err := g.SetView("status_bar", 0, maxY-4, maxX-1, maxY-2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			view.Frame = true
			view.Editable = false
		}

		return nil
	})

	tui := &Tui{
		g,
		nil,
		nil,
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

	if err := tui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, tui.onEnterPressed); err != nil {
		log.Panicln(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyF5, gocui.ModNone, tui.onF5Pressed); err != nil {
		log.Panicln(err)
	}
	if err := tui.SetKeybinding("", gocui.KeyBackspace, gocui.ModNone, tui.onBackspacePressed); err != nil {
		log.Panicln(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}
}

func (tui *Tui) onEnterPressed(g *gocui.Gui, v *gocui.View) error {
	selection, err := v.Word(v.Cursor())
	if err != nil {
		return err
	}

	table := tui.DatabaseContext.FetchTable(selection)

	names := table.ColumnNames()
	log.Println(strings.Join(names, ", "))
	// log.Println(strings.Join(types, ", "))
	// log.Println(strings.Join(dbTypes, ", "))
	tui.Update(func(g *gocui.Gui) error {
		mainView, err := g.View("main")
		if err != nil {
			return err
		}
		maxX, _ := g.Size()
		_, mainSizeY := mainView.Size()
		currentX := maxX / 6
		for i, col := range table.Columns() {
			g.SetViewOnBottom(col.Name)
			data, width := table.GetColumnRows(i)
			colView, err := g.SetView(col.Name, currentX, 0, currentX+width, mainSizeY) //mainSizeX/len(table.Columns()
			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				colView.Title = col.Name
				colView.Frame = true
				colView.Editable = false
				fmt.Fprint(colView, strings.Join(data, "\n"))
			}
			currentX += width
		}
		return nil
	})

	return nil
}

func (tui *Tui) onF5Pressed(g *gocui.Gui, v *gocui.View) error {
	tablesWithCounts := tui.DatabaseContext.FetchCounts(v.BufferLines())
	tui.Update(func(g *gocui.Gui) error {
		v.Clear()
		fmt.Fprint(v, strings.Join(tablesWithCounts, "\n"))
		return nil
	})
	return nil
}

func (tui *Tui) onBackspacePressed(g *gocui.Gui, v *gocui.View) error {

	tui.Update(func(g *gocui.Gui) error {
		return nil
	})
	return nil
}

func (tui *Tui) SetDatabaseContext(context *database.DatabaseContext) {
	tui.DatabaseContext = context

	tui.onSetDatabaseContext()
}

func (tui *Tui) onSetDatabaseContext() {
	dbName := tui.DatabaseContext.FetchDatabaseName()
	tableNames := tui.DatabaseContext.FetchTableNames()

	tui.Update(func(g *gocui.Gui) error {
		view, err := g.View("table_list")
		if err != nil {
			log.Println(err)
			return err
		}
		view.Clear()

		view.Title = dbName
		view.Highlight = true
		view.SelFgColor = gocui.ColorBlack | gocui.AttrBold
		view.SelBgColor = gocui.ColorWhite

		fmt.Fprint(view, strings.Join(tableNames, "\n"))

		g.SetCurrentView("table_list")
		if err := view.SetCursor(0, 0); err != nil {
			return err
		}
		return nil
	})
}
