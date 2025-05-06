package tui

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jesseduffield/gocui"
)

func (tui *Tui) setKeybinds() {
	tui.setSidePanelKeybinds()

	tui.setDataviewKeybinds()

	tui.setTableCellKeybinds()

	tui.setCommandBarKeybinds()

	if err := tui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, tui.cycleCurrentView()); err != nil {
		logger.Panic(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyF5, gocui.ModNone, OnF5Pressed(tui)); err != nil {
		logger.Panic(err)
	}

	if err := tui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		logger.Panic(err)
	}
}

func (tui *Tui) setSidePanelKeybinds() {
	if err := tui.SetKeybinding(views.SidePanelViewName, gocui.KeyEnter, gocui.ModNone, tui.sidePanel.OnEnterPressed(tui.goolDb)); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.SidePanelViewName, gocui.KeyArrowUp, gocui.ModNone, tui.sidePanel.MoveCursorUp()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.SidePanelViewName, gocui.KeyArrowDown, gocui.ModNone, tui.sidePanel.MoveCursorDown()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.SidePanelViewName, 'k', gocui.ModNone, tui.sidePanel.MoveCursorUp()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.SidePanelViewName, 'j', gocui.ModNone, tui.sidePanel.MoveCursorDown()); err != nil {
		logger.Panic(err)
	}
}

func (tui *Tui) setDataviewKeybinds() {
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyArrowUp, gocui.ModNone, tui.dataView.MoveRowUp()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyArrowDown, gocui.ModNone, tui.dataView.MoveRowDown()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyArrowLeft, gocui.ModNone, tui.dataView.MoveColumnLeft()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyArrowRight, gocui.ModNone, tui.dataView.MoveColumnRight()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, 'k', gocui.ModNone, tui.dataView.MoveRowUp()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, 'j', gocui.ModNone, tui.dataView.MoveRowDown()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, 'h', gocui.ModNone, tui.dataView.MoveColumnLeft()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, 'l', gocui.ModNone, tui.dataView.MoveColumnRight()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.Update(func(g *gocui.Gui) error {
			data := tui.dataView.GetSelectedCellData()

			x0, y0, x1, y1 := v.Dimensions()
			dataView, err := g.SetView(views.TableCellViewName, x0+5, y0+5, x1-5, y1-5, 0)

			if err != nil {
				if !gocui.IsUnknownView(err) {
					logger.Panic(err)
				}
				dataView.Frame = true
				dataView.Wrap = true
				dataView.Editable = false
			}
			dataView.Clear()
			fmt.Fprint(dataView, data)
			g.SetCurrentView(views.TableCellViewName)
			g.SetViewOnTop(views.TableCellViewName)
			return nil
		})
		return nil
	}); err != nil {
		logger.Panic(err)
	}

}

func (tui *Tui) setTableCellKeybinds() {
	if err := tui.SetKeybinding(views.TableCellViewName, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		if err := g.DeleteView(views.TableCellViewName); err != nil {
			logger.Panic(err)
		}
		g.SetCurrentView(views.DataTableViewName)
		return nil
	}); err != nil {
		logger.Panic(err)
	}

	// TableCell movement
	if err := tui.SetKeybinding(views.TableCellViewName, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, tableCellView *gocui.View) error {
		tableCellView.ScrollUp(3)
		return nil
	}); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.TableCellViewName, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, tableCellView *gocui.View) error {
		tableCellView.ScrollDown(3)
		return nil
	}); err != nil {
		logger.Panic(err)
	}
}

func (tui *Tui) setCommandBarKeybinds() {
	if err := tui.SetKeybinding("", ':', gocui.ModNone, func(g *gocui.Gui, prevView *gocui.View) error {
		commandBarView, err := g.View(views.CommandBarViewName)
		if err != nil {
			logger.Panic(err)
		}
		g.Update(func(g *gocui.Gui) error {
			commandBarView.Clear()
			commandBarView.SetCursor(0, 0)
			return nil
		})
		g.Cursor = true

		tui.prevView = prevView.Name()
		g.SetCurrentView(views.CommandBarViewName)
		return nil
	}); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.CommandBarViewName, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, commandBarView *gocui.View) error {
		g.Cursor = false
		g.Update(func(g *gocui.Gui) error {
			commandBarView.Clear()
			fmt.Fprint(commandBarView, ">")
			return nil
		})
		g.SetCurrentView(tui.prevView)
		return nil
	}); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.CommandBarViewName, gocui.KeyEnter, gocui.ModNone, tui.commandBar.OnEnterPressed(tui.goolDb)); err != nil {
		logger.Panic(err)
	}
}
