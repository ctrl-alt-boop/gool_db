package tui

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/ctrl-alt-boop/gooldb/tui/config"
	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jesseduffield/gocui"
)

var logger = logging.NewLogger("tui.log")

type Tui struct {
	*gocui.Gui
	goolDb *gooldb.GoolDb

	sidePanel    *views.SidePanelView
	commandBar   *views.CommandBarView
	dataView     *views.DataTableView
	queryOptions *views.QueryOptionsView
	helpFooter   *views.HelpFooterView

	notificationHandler *views.NotificationHandler
	prevView            string
}

func Create(notifier *Notifier, goolDb *gooldb.GoolDb) *Tui {
	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:      gocui.Output256,
		SupportOverlaps: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create new gocui.Gui: %v", err))
	}

	sidePanel := &views.SidePanelView{}
	commandBar := &views.CommandBarView{
		GoolDb: goolDb,
	}
	dataView := &views.DataTableView{}
	queryOptions := &views.QueryOptionsView{}
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
		queryOptions:        queryOptions,
		helpFooter:          helpFooter,
		notificationHandler: notificationHandler,
		prevView:            "",
	}

	appConfig := config.GetDefaultKeybindings()

	// Load keybindings configuration
	// keybindingFilePath := ".config/gooldb/binds.yaml"
	// appConfig, err := config.LoadKeybindingConfig(keybindingFilePath)
	// if err != nil {
	// 	logger.Warnf("Failed to load keybindings from %s: %v. Using default keybindings.", keybindingFilePath, err)
	// 	appConfig = config.GetDefaultKeybindings()
	// }

	tui.ApplyKeybindingConfig(appConfig) // Apply the loaded or default keybindings

	tui.ShowListFooter = true

	if notifier != nil {
		notifier.SetOnNotification(func(level gooldb.NotificationLevel, args ...any) {
			if tui.notificationHandler != nil {
				tui.notificationHandler.NewNotification(level, args...)
			} else {
				logger.Warn("Notification received, but notificationHandler is nil.")
			}
		})
	} else {
		logger.Warn("Notifier is nil, notifications will not be displayed in TUI.")
	}

	return tui
}

func (tui *Tui) previousView() {
	if tui.prevView != "" && tui.prevView != views.CommandBarViewName {
		tui.SetCurrentView(tui.prevView)
	} else {
		tui.SetCurrentView(views.SidePanelViewName)
	}
}

// SetCurrentView updates the current view and stores the previous one.
func (tui *Tui) SetCurrentView(name string) *gocui.View {
	currentViewName := ""
	if cv := tui.CurrentView(); cv != nil {
		currentViewName = cv.Name()
	}

	view, err := tui.Gui.SetCurrentView(name)
	if err != nil {
		logger.Errorf("Failed to set current view to '%s': %v", name, err)
		return view
	}

	// Only update prevView if the view switch was successful and different
	if currentViewName != name {
		tui.prevView = currentViewName
	}
	return view
}

func (tui *Tui) Run() error {
	defer tui.Gui.Close()

	if err := tui.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		logger.Errorf("Gocui MainLoop error: %v", err)
		return err
	}
	return nil
}
