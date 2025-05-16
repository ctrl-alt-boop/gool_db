package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		//

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.help.Keys.Quit):
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd

	_, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = m.panel.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = m.command.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = m.data.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
