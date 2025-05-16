package tui

import "github.com/charmbracelet/lipgloss"

func (m appModel) View() string {
	panelView := m.panel.View()
	commandBarView := m.command.View()
	dataView := m.data.View()
	helpFooterView := m.help.View()
	// popup := m.popupHandler.PopupView()

	mainView := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		panelView,
		dataView,
	)

	render := lipgloss.JoinVertical(
		lipgloss.Top,
		mainView,
		commandBarView,
		helpFooterView,
	)

	return render
}
