package popups

import tea "github.com/charmbracelet/bubbletea"

type TableCellPopup struct {
	Value string
}

// Exec implements PopupModel.
func (t *TableCellPopup) Exec() {
	panic("unimplemented")
}

// Cancel implements PopupModel.
func (t *TableCellPopup) Cancel() {
	panic("unimplemented")
}

// Init implements tea.Model.
func (t *TableCellPopup) Init() tea.Cmd {
	panic("unimplemented")
}

// Update implements tea.Model.
func (t *TableCellPopup) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

// View implements tea.Model.
func (t *TableCellPopup) View() string {
	panic("unimplemented")
}
