package popups

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

const (
	QueryOptions PopupType = "query"
	TableCell    PopupType = "cell"
)

type (
	PopupType string

	PopupModel interface {
		tea.Model
		Exec()
		Cancel()
	}

	PopupHandler struct {
		goolDb *gooldb.GoolDb

		currentPopup PopupModel
	}
)

func NewHandler(gool *gooldb.GoolDb) *PopupHandler {
	return &PopupHandler{
		goolDb: gool,
	}
}

func (p *PopupHandler) Init() tea.Cmd {
	return nil
}

func (p *PopupHandler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if p.currentPopup != nil {
		_, cmd = p.currentPopup.Update(msg)
	}
	return p, cmd
}

func (p *PopupHandler) View() string {
	if p.currentPopup != nil {
		return p.currentPopup.View()
	}
	return ""
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
