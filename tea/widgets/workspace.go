package widgets

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/tea/event"
	"github.com/ctrl-alt-boop/gooldb/tea/ui"
)

type Workspace struct {
	width, height int
	goolDb        *gooldb.GoolDb

	table        []list.Model
	columnWidths []int

	delegate  list.ItemDelegate
	isLoading bool
	spinner   spinner.Model
}

func CreateDataArea(gool *gooldb.GoolDb) *Workspace {
	return &Workspace{
		goolDb:       gool,
		table:        make([]list.Model, 0),
		columnWidths: make([]int, 0),
	}
}

func (d *Workspace) Init() tea.Cmd {
	d.spinner = spinner.New()
	d.spinner.Spinner = ui.MovingBlock
	d.isLoading = false

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle.Height(1)
	delegate.Styles.SelectedTitle.Height(1)
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	d.delegate = delegate

	return nil
}

func (d *Workspace) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.UpdateSize(msg.Width, msg.Height)
	case event.GoolDbEventMsg:
		d.isLoading = false
		if msg.Err != nil {
			logger.Error(msg.Err)
			return d, event.NewGoolDbError(msg.Err)
		}
		switch msg.Type {
		case gooldb.TableSet:
			args, ok := msg.Args.(gooldb.TableSetEvent)
			if ok {
				d.SetTable(args.Table)
			}
		}
	}

	return d, nil
}

func (d *Workspace) SetTable(table *gooldb.DataTable) {
	column, width := table.GetColumnRows(0)
	items := make([]list.Item, len(column))
	for i, item := range column {
		items[i] = ui.ListItem(item)
	}
	list := list.New(items, d.delegate, width, 0)
	list.SetShowHelp(false)
	list.SetShowStatusBar(false)
	list.SetShowFilter(false)
	list.SetFilteringEnabled(false)
	list.SetShowPagination(false)
	list.SetShowTitle(false)
	d.table = append(d.table, list)
	d.columnWidths = append(d.columnWidths, width)
}

func (d *Workspace) UpdateSize(termWidth, termHeight int) {
	panelWidth := termWidth/PanelWidthRatio - BorderThicknessDouble
	d.width, d.height = termWidth-panelWidth-BorderThicknessDouble-1, termHeight-5
}

func (d *Workspace) View() string {
	dataBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Right:       "│",
		TopLeft:     "┬",
		TopRight:    "┐",
		BottomLeft:  "┴",
		BottomRight: "┤",
	}

	dataStyle := lipgloss.NewStyle().
		Width(d.width).
		Height(d.height).
		Border(dataBorder, true, true, true, false).
		Align(lipgloss.Left, lipgloss.Top)

	if d.isLoading {
		return dataStyle.
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(d.spinner.View())
	}

	views := make([]string, len(d.table))
	for i, table := range d.table {
		table.SetHeight(d.height)
		table.SetWidth(d.columnWidths[i])
		views[i] = table.View()
	}

	table := lipgloss.JoinHorizontal(lipgloss.Top, views...)

	return dataStyle.Render(table)
}
