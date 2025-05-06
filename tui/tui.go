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
	prevView            string
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
		prevView:            "",
	}

	tui.setKeybinds()
	notifier.SetOnNotification(func(level gooldb.NotificationLevel, args ...any) {
		tui.notificationHandler.NewNotification(level, args...)
	})

	return tui
}

func (tui *Tui) cycleCurrentView() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, currentView *gocui.View) error {

		switch currentView.Name() {
		case views.SidePanelViewName:
			view := tui.SetCurrentView(views.DataTableViewName)
			views.SetSidePanelInactiveColors(currentView)
			g.Cursor = false
			view.SetCursor(1, 1)
			tui.dataView.HighlightSelectedCell()

		default:
			view := tui.SetCurrentView(views.SidePanelViewName)
			views.SetSidePanelColors(view)
			g.Cursor = false
			currentView.Highlight = false
		}

		return nil
	}
}

func (tui *Tui) SetCurrentView(name string) *gocui.View {
	currentView := tui.CurrentView().Name()
	view, err := tui.Gui.SetCurrentView(name)
	if err != nil {
		logger.Panic(err)
	}
	tui.prevView = currentView
	return view
}
