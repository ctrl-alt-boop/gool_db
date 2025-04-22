package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jroimartin/gocui"
)

func MoveCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, -1, false)
	}
	return nil
}

func MoveCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, 1, false)
	}
	return nil
}

func OnEnterPressed(tui *Tui) func(g *gocui.Gui, v *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		switch currentView.Name() {
		case views.SidePanelViewName:
			word, err := currentView.Word(currentView.Cursor())
			if err != nil {
				return err
			}
			log.Println("OnEnter:", word)
			views.SidePanel.OnEnter(tui.GoolDb, word)
		case views.CommandBarViewName:
			log.Println("OnEnter:", currentView.Buffer())
			views.CommandBar.OnEnter(currentView.Buffer())
		default:
			if strings.HasPrefix(currentView.Name(), "column_") {
				_, y := currentView.Cursor()
				log.Println("OnEnter:", currentView.Name(), "row", y)
				views.DataTable.OnEnter(strings.TrimPrefix(currentView.Name(), "column_"))
			}
		}

		return nil
	}
}

func OnF5Pressed(tui *Tui) func(g *gocui.Gui, v *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		tablesWithCounts := tui.GoolDb.FetchCounts(currentView.BufferLines())
		tui.Update(func(g *gocui.Gui) error {
			currentView.Clear()
			fmt.Fprint(currentView, strings.Join(tablesWithCounts, "\n"))
			return nil
		})
		return nil
	}
}

func CycleCurrentView(tui *Tui) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, currentView *gocui.View) error {

		if tui.CurrentDataColumns == nil {
			return nil
		}

		switch currentView.Name() {
		case views.SidePanelViewName:
			view, err := tui.SetCurrentView(tui.CurrentDataColumns[0])
			if err != nil {
				log.Fatal(err)
			}
			view.Highlight = true
			view.SelFgColor = gocui.AttrReverse
			view.SelBgColor = gocui.AttrReverse
			if err := view.SetCursor(0, 0); err != nil {
				return err
			}
		default:
			view, err := tui.SetCurrentView(views.SidePanelViewName)
			if err != nil {
				log.Fatal(err)
			}
			currentView.Highlight = false

			view.Highlight = true
			view.SelFgColor = gocui.AttrReverse
			view.SelBgColor = gocui.AttrReverse
		}

		return nil
	}
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
