package popups

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
)

type QueryOptionsPopup struct {
	Query *query.Statement
}

// Exec implements PopupModel.
func (q *QueryOptionsPopup) Exec() {
	panic("unimplemented")
}

// Cancel implements PopupModel.
func (q *QueryOptionsPopup) Cancel() {
	panic("unimplemented")
}

// Init implements tea.Model.
func (q *QueryOptionsPopup) Init() tea.Cmd {
	panic("unimplemented")
}

// Update implements tea.Model.
func (q *QueryOptionsPopup) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return q, nil
}

// View implements tea.Model.
func (q *QueryOptionsPopup) View() string {
	panic("unimplemented")
}
