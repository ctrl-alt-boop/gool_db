package tui

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jesseduffield/gocui"
)

func MoveCursorUp() func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, currentView *gocui.View) error {
		if currentView != nil {
			MoveCursorWithScrolling(currentView, 0, -1)
		}
		return nil
	}
}

func MoveCursorDown() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		if currentView != nil {
			MoveCursorWithScrolling(currentView, 0, 1)
		}
		return nil
	}
}

func MoveCursorLeft() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		MoveCursorWithScrolling(currentView, -1, 0)
		return nil
	}
}

func MoveCursorRight() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		MoveCursorWithScrolling(currentView, 1, 0)
		return nil
	}
}

func (tui *Tui) cycleCurrentView() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(g *gocui.Gui, currentView *gocui.View) error {
		var nextViewName string

		switch currentView.Name() {
		case views.SidePanelViewName:
			if !tui.dataView.IsTableSet() {
				return nil
			}
			nextViewName = views.DataTableViewName
			tui.sidePanel.SetInactiveColors()

		case views.DataTableViewName:
			nextViewName = views.SidePanelViewName
			tui.dataView.ClearHighlight()
			tui.sidePanel.SetActiveColors()
		}

		newActiveView := tui.SetCurrentView(nextViewName)

		if newActiveView != nil {
			g.Cursor = false
			if nextViewName == views.DataTableViewName {
				// newActiveView.SetCursor(1, 1) // TODO: Maybe not
				// newActiveView.SetOrigin(0, 0) // TODO: Maybe not
				tui.dataView.HighlightSelectedCell()
			}
		}

		return nil
	}
}

func (tui *Tui) onF5Pressed() func(g *gocui.Gui, v *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		tablesWithCounts := tui.goolDb.FetchCounts(currentView.BufferLines())
		tui.Update(func(g *gocui.Gui) error {
			currentView.Clear()
			fmt.Fprint(currentView, strings.Join(tablesWithCounts, "\n"))
			return nil
		})
		return nil
	}
}

func MoveCursorWithScrolling(v *gocui.View, dx, dy int) {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	maxX, maxY := v.Size()

	newCx := cx + dx
	newCy := cy + dy

	newOx := ox
	newOy := oy

	// Adjust origin for vertical scrolling
	if newCy < oy {
		newOy = newCy // Scroll up
	} else if newCy >= oy+maxY {
		newOy = newCy - maxY + 1 // Scroll down
	}

	// Adjust origin for horizontal scrolling
	if newCx < ox {
		newOx = newCx // Scroll left
	} else if newCx >= ox+maxX {
		newOx = newCx - maxX + 1 // Scroll right
	}

	// Ensure origin doesn't go below zero
	if newOx < 0 {
		newOx = 0
	}
	if newOy < 0 {
		newOy = 0
	}

	v.SetCursor(newCx, newCy)
	// Apply the new origin if it changed
	if newOx != ox || newOy != oy {
		v.SetOrigin(newOx, newOy)
	}
}
