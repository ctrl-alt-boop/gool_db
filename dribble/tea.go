package dribble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/dribble/message"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget/popups"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("tea.log")

type AppModel struct {
	gooldb *gooldb.GoolDb

	panel        *widget.Panel
	command      *widget.CommandLine
	workspace    *widget.Workspace
	help         *widget.Help
	popupHandler *popups.PopupHandler

	inFocus widget.Kind

	programSend func(msg tea.Msg)
}

func NewModel(gool *gooldb.GoolDb) AppModel {
	return AppModel{
		gooldb:       gool,
		panel:        widget.NewPanel(gool),
		command:      widget.NewCommandBar(gool),
		workspace:    widget.NewWorkspace(gool),
		help:         widget.NewHelp(),
		popupHandler: popups.NewHandler(gool),
	}
}

func (m AppModel) SetProgramSend(send func(msg tea.Msg)) {
	m.programSend = send

	m.gooldb.RegisterEventHandler(gooldb.DriverSet, m.onEventFunc(gooldb.DriverSet))
	m.gooldb.RegisterEventHandler(gooldb.DatabaseSet, m.onEventFunc(gooldb.DatabaseSet))
	m.gooldb.RegisterEventHandler(gooldb.TableSet, m.onEventFunc(gooldb.TableSet))
}

func (m AppModel) Init() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	cmd = m.panel.Init()
	cmds = append(cmds, cmd)
	cmd = m.command.Init()
	cmds = append(cmds, cmd)
	cmd = m.workspace.Init()
	cmds = append(cmds, cmd)
	cmd = m.help.Init()
	cmds = append(cmds, cmd)
	cmd = m.popupHandler.Init()
	cmds = append(cmds, cmd)

	cmds = append(cmds, widget.ChangeFocus(widget.KindPanel))

	return tea.Batch(cmds...)
}

func (m AppModel) onEventFunc(eventType gooldb.EventType) func(a any, err error) {
	return func(a any, err error) {
		event := message.GoolDbEventMsg{
			Type: eventType,
			Args: a,
			Err:  err,
		}
		m.programSend(event)
	}
}
