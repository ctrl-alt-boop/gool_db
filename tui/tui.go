package tui

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/ctrl-alt-boop/gooldb/tui/config"
	"github.com/ctrl-alt-boop/gooldb/tui/widgets"
	"github.com/jesseduffield/gocui"
)

var logger = logging.NewLogger("tui.log")

type Tui struct {
	*gocui.Gui
	goolDb *gooldb.GoolDb

	sidePanel    *widgets.SidePanel
	commandBar   *widgets.CommandBar
	dataView     *widgets.DataArea
	queryOptions *widgets.QueryOptions
	helpFooter   *widgets.HelpFooter

	notificationHandler *widgets.NotificationHandler
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
	drawer := widgets.CreateDrawer(g)

	sidePanel := widgets.CreateSidePanel(drawer, goolDb)
	commandBar := widgets.CreateCommandBar(drawer, goolDb)
	dataArea := widgets.CreateDataArea(drawer)
	queryOptions := &widgets.QueryOptions{}
	// helpFooter := widgets.CreateHelpFooter()
	notificationHandler := &widgets.NotificationHandler{}

	goolDb.RegisterEventHandler(gooldb.DriverSet, sidePanel.OnDriverSet)
	goolDb.RegisterEventHandler(gooldb.DatabaseSet, sidePanel.OnDatabaseSet)
	goolDb.RegisterEventHandler(gooldb.TableSet, dataArea.OnTableSet)

	g.SetManager(sidePanel, commandBar, dataArea, widgets.Help, drawer)

	tui := &Tui{
		Gui:                 g,
		goolDb:              goolDb,
		sidePanel:           sidePanel,
		commandBar:          commandBar,
		dataView:            dataArea,
		queryOptions:        queryOptions,
		helpFooter:          widgets.Help,
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

func (tui *Tui) Run() error {
	defer tui.Gui.Close()

	if err := tui.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		logger.Fatalf("Gocui MainLoop error: %v", err)
		return err
	}
	return nil
}
