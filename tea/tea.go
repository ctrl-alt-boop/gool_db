package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/ctrl-alt-boop/gooldb/tea/event"
	"github.com/ctrl-alt-boop/gooldb/tea/widgets"
)

var logger = logging.NewLogger("tea.log")

type appModel struct {
	gooldb *gooldb.GoolDb

	panel   *widgets.Panel
	command *widgets.CommandBar
	data    *widgets.Workspace
	help    *widgets.HelpFooter

	popupHandler *widgets.PopupHandler

	focusedView string

	programSend func(msg tea.Msg)
}

func NewModel(gool *gooldb.GoolDb) appModel {
	return appModel{
		gooldb: gool,

		panel:        widgets.CreateListPanel(gool),
		command:      widgets.CreateCommandBar(gool),
		data:         widgets.CreateDataArea(gool),
		help:         widgets.CreateHelp(),
		popupHandler: widgets.CreatePopupHandler(gool),

		focusedView: widgets.Names.Panel,
	}
}

func (m appModel) SetProgramSend(send func(msg tea.Msg)) {
	m.programSend = send

	m.gooldb.RegisterEventHandler(gooldb.DriverSet, m.onEventFunc(gooldb.DriverSet))
	m.gooldb.RegisterEventHandler(gooldb.DatabaseSet, m.onEventFunc(gooldb.DatabaseSet))
	m.gooldb.RegisterEventHandler(gooldb.TableSet, m.onEventFunc(gooldb.TableSet))
}

func (m appModel) Init() tea.Cmd {

	m.help.Init()
	m.panel.Init()
	m.data.Init()
	m.command.Init()

	return nil
}

func (m appModel) onEventFunc(eventType gooldb.EventType) func(a any, err error) {
	return func(a any, err error) {
		logger.Info("Got event: ", eventType)
		event := event.GoolDbEventMsg{
			Type: eventType,
			Args: a,
			Err:  err,
		}
		m.programSend(event)
	}
}
