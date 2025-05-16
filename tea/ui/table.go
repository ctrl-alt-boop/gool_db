package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type Table struct {
	table []List
}

func NewTable(width, height int) *Table {
	return &Table{
		table: make([]List, 0),
	}
}

func (t *Table) SetTable(table *gooldb.DataTable) {
	t.table = make([]List, table.NumColumns())
	for i := range table.Columns() {
		column, width := table.GetColumnRows(i)
		t.table[i].Model.SetWidth(width)
		t.table[i].Model.Styles.Title = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			Padding(0, 1)
		t.table[i].SetItems(column)
		t.table[i].Model.SetShowTitle(true)
	}
}

func (t *Table) View() string {
	views := make([]string, len(t.table))
	for i, table := range t.table {
		views[i] = table.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}

func (t *Table) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Table) Init() tea.Cmd {
	return nil
}
