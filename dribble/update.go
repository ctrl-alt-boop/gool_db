package dribble

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/dribble/message"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// AppModel messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		_, cmd = m.panel.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.workspace.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.popupHandler.Update(msg)
		cmds = append(cmds, cmd)

	case widget.RequestFocus:
		m.inFocus = widget.Kind(msg)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.CycleView):
			switch m.inFocus {
			case widget.KindPanel:
				m.inFocus = widget.KindWorkspace
			case widget.KindWorkspace:
				m.inFocus = widget.KindCommandLine
			case widget.KindCommandLine:
				m.inFocus = widget.KindPopupHandler
			case widget.KindPopupHandler:
				m.inFocus = widget.KindPanel
			}
			return m, nil
		case key.Matches(msg, config.Keys.Quit):
			return m, tea.Quit
			// case key.Matches(msg, config.Keys.CycleView)
		}
		// case message.CommandExecMsg:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.inFocus == widget.KindPanel:
			_, cmd = m.panel.Update(msg)
			cmds = append(cmds, cmd)
		case m.inFocus == widget.KindWorkspace:
			_, cmd = m.workspace.Update(msg)
			cmds = append(cmds, cmd)
		case m.inFocus == widget.KindCommandLine:
			_, cmd = m.command.Update(msg)
			cmds = append(cmds, cmd)
		case m.inFocus == widget.KindPopupHandler:
			_, cmd = m.popupHandler.Update(msg)
			cmds = append(cmds, cmd)
		}
		_, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
	case message.GoolDbEventMsg:
		_, cmd = m.panel.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.workspace.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.popupHandler.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
