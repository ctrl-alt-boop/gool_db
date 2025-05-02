package tui

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jesseduffield/gocui"
)

var logger = logging.NewLogger("tui.log")

type Tui struct {
	*gocui.Gui
	goolDb *gooldb.GoolDb

	sidePanel  *views.SidePanelView
	commandBar *views.CommandBarView
	dataView   *views.DataTableView
	helpFooter *views.HelpFooterView

	notificationHandler *views.NotificationHandler
}

func Create(notifier *Notifier, goolDb *gooldb.GoolDb) *Tui {
	g, err := gocui.NewGui(gocui.NewGuiOpts{OutputMode: gocui.Output256})
	if err != nil {
		panic(err)
	}
	sidePanel := &views.SidePanelView{}
	commandBar := &views.CommandBarView{}
	dataView := &views.DataTableView{}
	helpFooter := &views.HelpFooterView{}
	notificationHandler := &views.NotificationHandler{}

	goolDb.RegisterEventHandler(gooldb.DriverSet, sidePanel.OnDriverSet)
	goolDb.RegisterEventHandler(gooldb.DatabaseSet, sidePanel.OnDatabaseSet)
	goolDb.RegisterEventHandler(gooldb.TableSet, dataView.OnTableSet)

	g.SetManager(sidePanel, commandBar, dataView, helpFooter)

	tui := &Tui{
		Gui:                 g,
		goolDb:              goolDb,
		sidePanel:           sidePanel,
		commandBar:          commandBar,
		dataView:            dataView,
		helpFooter:          helpFooter,
		notificationHandler: notificationHandler,
	}

	tui.setKeybinds()
	notifier.SetOnNotification(func(level gooldb.NotificationLevel, args ...any) {
		tui.notificationHandler.NewNotification(level, args...)
	})

	return tui
}

func (tui *Tui) setKeybinds() {
	if err := tui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, MoveCursorUp()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, MoveCursorDown()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", 'k', gocui.ModNone, MoveCursorUp()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", 'j', gocui.ModNone, MoveCursorDown()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, MoveCursorLeft()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, MoveCursorRight()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", 'h', gocui.ModNone, MoveCursorLeft()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding("", 'l', gocui.ModNone, MoveCursorRight()); err != nil {
		logger.Panic(err)
	}

	// DataView binds
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyArrowLeft, gocui.ModNone, tui.MoveColumnLeft()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyArrowRight, gocui.ModNone, tui.MoveColumnRight()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, 'h', gocui.ModNone, tui.MoveColumnLeft()); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, 'l', gocui.ModNone, tui.MoveColumnRight()); err != nil {
		logger.Panic(err)
	}

	if err := tui.SetKeybinding(views.SidePanelViewName, gocui.KeyEnter, gocui.ModNone, tui.sidePanel.OnEnterPressed(tui.goolDb)); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.CommandBarViewName, gocui.KeyEnter, gocui.ModNone, tui.commandBar.OnEnterPressed(tui.goolDb)); err != nil {
		logger.Panic(err)
	}
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeyEnter, gocui.ModNone, tui.dataView.OnEnterPressed(tui.goolDb)); err != nil {
		logger.Panic(err)
	}
	// if err := tui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, OnEnterPressed(tui)); err != nil {
	// 	logger.Panic(err)
	// }
	if err := tui.SetKeybinding(views.DataTableViewName, gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.SetViewOnTop(views.DataTableViewName)
		g.SetCurrentView(views.DataTableViewName)
		x0, y0, x1, y1, err := g.ViewPosition(views.DataTableViewName)
		logger.Info("view: ", v.Name())
		logger.Info("view pos: ", x0, y0, y1, x1)
		if err != nil {
			logger.Panic(err)
		}
		cursorX, cursorY := v.Cursor()
		correctedX, correctedY := x0+cursorX+1, y0+cursorY+1
		r, _ := g.Rune(correctedX, correctedY)
		logger.Info("before: cursor: ", correctedX, correctedY, " rune: ", string(r))
		v.Visible = false
		g.SetRune(correctedX, correctedY, 'O', gocui.ColorBlue, gocui.ColorGreen)
		r, _ = g.Rune(correctedX, correctedY)
		logger.Info("after: cursor: ", correctedX, correctedY, " rune: ", string(r))
		return nil
	}); err != nil {
		logger.Panic(err)
	}

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

func (tui *Tui) cycleCurrentView() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, currentView *gocui.View) error {

		switch currentView.Name() {
		case views.SidePanelViewName:
			view, err := tui.SetCurrentView(views.DataTableViewName)
			if err != nil {
				logger.Panic(err)
			}
			views.SetSidePanelInactiveColors(currentView)
			g.Cursor = true
			view.SetCursor(1, 0)

		default:
			view, err := tui.SetCurrentView(views.SidePanelViewName)
			if err != nil {
				logger.Panic(err)
			}
			views.SetSidePanelColors(view)
			g.Cursor = false
			currentView.Highlight = false
		}

		return nil
	}
}

func (tui *Tui) MoveColumnLeft() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		if currentView != nil {
			cx := tui.dataView.LeftColumn()
			MoveCursorWithScrolling(currentView, -cx, 0)
		}
		return nil
	}
}

func (tui *Tui) MoveColumnRight() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		if currentView != nil {
			cx := tui.dataView.RightColumn()
			MoveCursorWithScrolling(currentView, cx, 0)
		}
		return nil
	}
}
