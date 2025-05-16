package widgets

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
)

type PopupType string

const (
	QueryOptions PopupType = "query_options"
	TableCell    PopupType = "table_cell"
)

type PopupHandler struct {
	width, height int
	goolDb        *gooldb.GoolDb

	currentPopup PopupModel
}

func CreatePopupHandler(gool *gooldb.GoolDb) *PopupHandler {
	return &PopupHandler{
		goolDb: gool,
	}
}

func (p *PopupHandler) PopupView() string {
	if p.currentPopup == nil {
		return ""
	}
	return p.currentPopup.View()
}

func (p *PopupHandler) UpdateSize(termWidth, termHeight int) {
	p.width, p.height = termWidth/2, termHeight/2
}

func (p *PopupHandler) Popup(popupType PopupType, args ...any) {
	switch popupType {
	case QueryOptions:
		p.currentPopup = &QueryOptionsPopup{}
	case TableCell:
		value := args[0].(string)
		p.currentPopup = &TableCellPopup{
			Value: value,
		}
	default:
		p.currentPopup = nil
	}
}

type PopupModel interface {
	tea.Model
	Exec()
	Cancel()
}

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
