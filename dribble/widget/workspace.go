package widget

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/message"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type Workspace struct {
	width, height int
	goolDb        *gooldb.GoolDb

	table        *ui.BubblesTable
	columnWidths []int

	isLoading bool
	spinner   spinner.Model
}

func NewWorkspace(gool *gooldb.GoolDb) *Workspace {
	return &Workspace{
		goolDb:       gool,
		table:        ui.New(),
		columnWidths: make([]int, 0),
	}
}

func (d *Workspace) Init() tea.Cmd {
	d.spinner = spinner.New()
	d.spinner.Spinner = ui.MovingBlock
	d.isLoading = false

	return nil
}

func (d *Workspace) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.UpdateSize(msg.Width, msg.Height)
	case message.GoolDbEventMsg:
		d.isLoading = false
		if msg.Err != nil {
			logger.Error(msg.Err)
			return d, message.NewGoolDbError(msg.Err)
		}
		switch msg.Type {
		case gooldb.TableSet:
			args, ok := msg.Args.(gooldb.TableSetEvent)
			if ok {
				return d, d.SetTable(args.Table)
			}
		}
	}

	return d, nil
}

func (d *Workspace) SetTable(table *gooldb.DataTable) tea.Cmd {
	d.table.SetTable(table)
	return func() tea.Msg {
		return message.TableSet(true)
	}
}

func (d *Workspace) IsTableSet() bool {
	return d.table.IsTableSet()
}

func (d *Workspace) UpdateSize(termWidth, termHeight int) {
	panelWidth := termWidth/PanelWidthRatio - BorderThicknessDouble
	d.width, d.height = termWidth-panelWidth-BorderThicknessDouble-1, termHeight-5
	// d.table.SetHeight(d.height)
}

func (d *Workspace) View() string {
	dataBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Right:       "│",
		Left:        "│",
		TopLeft:     "┬",
		TopRight:    "┐",
		BottomLeft:  "┴",
		BottomRight: "┤",
	}

	width := d.width
	if len(d.table.Table.Columns()) > 0 {
		width = d.table.Table.Width()
	}

	dataStyle := lipgloss.NewStyle().
		Width(width).
		Height(d.height).
		Border(dataBorder, true, true, false, true).
		Align(lipgloss.Left, lipgloss.Top)

	if d.isLoading {
		return dataStyle.
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(d.spinner.View())
	}

	tableView := d.table.View()

	return dataStyle.Render(tableView)
}
