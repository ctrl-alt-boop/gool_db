package dribble

import (
	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) View() string {

	panelView := m.panel.View()
	commandBarView := m.command.View()
	workspaceView := m.workspace.View()
	helpFooterView := m.help.View()
	popupView := m.popupHandler.View()

	mainView := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		panelView,
		workspaceView,
	)

	panelViewWidth := lipgloss.Width(panelView)
	commandBarViewWidth := lipgloss.Width(commandBarView)

	panelBorder := lipgloss.PlaceHorizontal(
		panelViewWidth-1,
		lipgloss.Left,
		"─", lipgloss.WithWhitespaceChars("─"))

	workspaceBorder := lipgloss.PlaceHorizontal(
		commandBarViewWidth-panelViewWidth-2,
		lipgloss.Left,
		"─", lipgloss.WithWhitespaceChars("─"))

	rightSeparatorCorner := "┤"
	if m.workspace.IsTableSet() {
		rightSeparatorCorner = "┬"
	}

	separator := lipgloss.JoinHorizontal(
		lipgloss.Left,
		"├",
		panelBorder,
		"┴",
		workspaceBorder,
		rightSeparatorCorner,
	)

	render := lipgloss.JoinVertical(
		lipgloss.Top,
		mainView,
		separator,
		commandBarView,
		helpFooterView,
	)

	if popupView != "" {
		render = lipgloss.Place(lipgloss.Width(popupView), lipgloss.Height(popupView),
			lipgloss.Center, lipgloss.Center, popupView)
	}

	return render
}
