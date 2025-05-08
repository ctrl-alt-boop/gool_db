package views

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/jesseduffield/gocui"
	"github.com/mattn/go-runewidth"
)

const DataTableViewName string = "data"
const TableCellViewName string = "data_cell"

const firstRow = 1
const maxCellWidth = 36 // Guid length, including the '-'s

type dataViewState struct {
	currentColumnIndex int
	currentRowIndex    int
	columnWidths       []int
	table              *gooldb.DataTable
	contentWidth       int
	maxRows            int
}

type DataTableView struct {
	view *gocui.View
	gui  *gocui.Gui
	// gool  *gooldb.GoolDb
	state dataViewState
}

func (d *DataTableView) IsTableSet() bool {
	return d.state.table != nil
}

func (d *DataTableView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(DataView(maxX, maxY))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			logger.Panic(err)
		}
		view.Frame = false
		view.Editable = false
		d.view = view
		d.gui = g
		d.state = dataViewState{}
	}
	d.state.maxRows = view.InnerHeight() - firstRow
	return nil
}

func (d *DataTableView) OnEnterPressed(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {

		return nil
	}
}

func (d *DataTableView) OnTableSet(eventArgs any, err error) {
	logger.Info("OnTableSet: ", eventArgs, err)
	if err != nil {
		return
	}
	args, ok := eventArgs.(gooldb.TableSetEvent)
	if !ok {
		return
	}
	d.state.table = args.Table
	d.gui.Update(func(g *gocui.Gui) error {
		d.view.Clear()
		formatedRows := d.getFormatedRows()
		formatedHeader := d.getFormatedTitle()
		// d.view.Title = formatedHeader
		fmt.Fprint(d.view, formatedHeader)

		d.state.currentColumnIndex = 0
		d.state.currentRowIndex = firstRow

		logger.Info("header: ", formatedHeader)
		logger.Info("rows: ", strings.Join(formatedRows, "\n"))

		// fmt.Fprint(d.view, formatedHeader)
		fmt.Fprint(d.view, strings.Join(formatedRows, "\n"))
		return nil
	})
}

const (
	hiStart         = "\x1b[7m"
	hiEnd           = "\x1b[0m"
	columnSeparator = " \u2502 "
)

func (d *DataTableView) HighlightSelectedCell() {
	buf := d.view.Buffer()
	buf = strings.ReplaceAll(buf, hiStart, "")
	buf = strings.ReplaceAll(buf, hiEnd, "")
	data := strings.Split(buf, "\n")
	selectionStartIndex := 1 // 0 is a ' '
	for i := range d.state.columnWidths[:d.state.currentColumnIndex] {
		selectionStartIndex += d.state.columnWidths[i] + len(columnSeparator)
	}
	selectionEndIndex := selectionStartIndex + d.state.columnWidths[d.state.currentColumnIndex]
	highlightedData := data[d.state.currentRowIndex][:selectionStartIndex] + hiStart + data[d.state.currentRowIndex][selectionStartIndex:selectionEndIndex] + hiEnd + data[d.state.currentRowIndex][selectionEndIndex:]
	data[d.state.currentRowIndex] = highlightedData

	d.gui.Update(func(g *gocui.Gui) error {
		d.view.Clear()
		fmt.Fprint(d.view, strings.Join(data, "\n"))
		return nil
	})
}

func (d *DataTableView) ClearHighlight() {
	buf := d.view.Buffer()
	buf = strings.ReplaceAll(buf, hiStart, "")
	buf = strings.ReplaceAll(buf, hiEnd, "")
	d.gui.Update(func(g *gocui.Gui) error {
		d.view.Clear()
		fmt.Fprint(d.view, buf)
		return nil
	})
}

func (d *DataTableView) GetSelectedCellData() string {
	row, column := d.state.currentRowIndex-1, d.state.currentColumnIndex
	data, err := d.state.table.GetRowColumn(row, column)
	if err != nil {
		logger.Warn(err)
		return err.Error()
	}
	if d.state.table.Columns()[d.state.currentColumnIndex].DbType == "JSONB" {
		return PrettifyJson(data)
	}
	return data
}

func (d *DataTableView) MoveColumnLeft() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, dataTableView *gocui.View) error {
		if dataTableView == nil || d.state.currentColumnIndex <= 0 {
			return nil
		}
		d.state.currentColumnIndex--
		if d.state.currentColumnIndex == 0 {
			dataTableView.SetOriginX(0)
			dataTableView.SetCursorX(1)
			d.HighlightSelectedCell()
			return nil
		}
		selectionStartX := 1 // 0 is a ' '
		for i := range d.state.currentColumnIndex {
			selectionStartX += d.state.columnWidths[i] + runewidth.StringWidth(columnSeparator)
		}

		ox := dataTableView.OriginX()

		if selectionStartX < ox {
			dataTableView.ScrollLeft(d.state.columnWidths[d.state.currentColumnIndex] + 6)
		}

		newOx := dataTableView.OriginX()
		dataTableView.SetCursorX(selectionStartX - newOx)

		d.HighlightSelectedCell()
		return nil
	}
}

func (d *DataTableView) MoveColumnRight() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, dataTableView *gocui.View) error {
		if dataTableView == nil || d.state.currentColumnIndex >= d.state.table.NumColumns()-1 {
			return nil
		}
		ox := dataTableView.OriginX()
		viewWidth := dataTableView.InnerWidth()
		d.state.currentColumnIndex++

		selectionStartX := 1 // 0 is a ' '
		for i := range d.state.currentColumnIndex {
			selectionStartX += d.state.columnWidths[i] + runewidth.StringWidth(columnSeparator)
		}
		selectionEndX := selectionStartX + d.state.columnWidths[d.state.currentColumnIndex]

		if d.state.currentColumnIndex == d.state.table.NumColumns()-1 {
			dataTableView.SetOriginX(d.state.contentWidth - viewWidth)
		} else if selectionEndX > viewWidth+ox {
			dataTableView.ScrollRight(d.state.columnWidths[d.state.currentColumnIndex] + 6)
		}

		newOx := dataTableView.OriginX()
		dataTableView.SetCursorX(selectionStartX - newOx)

		d.HighlightSelectedCell()
		return nil
	}
}

func (d *DataTableView) MoveRowUp() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, dataTableView *gocui.View) error {
		if dataTableView == nil || d.state.currentRowIndex <= firstRow {
			return nil
		}
		curY := dataTableView.CursorY()
		d.state.currentRowIndex--
		dataTableView.SetCursorY(curY - 1)

		d.HighlightSelectedCell()
		return nil
	}
}

func (d *DataTableView) MoveRowDown() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, dataTableView *gocui.View) error {
		if dataTableView == nil || d.state.currentRowIndex >= len(d.state.table.Rows()) {
			return nil
		}
		curY := dataTableView.CursorY()
		d.state.currentRowIndex++
		dataTableView.SetCursorY(curY + 1)

		d.HighlightSelectedCell()
		return nil
	}
}

func (d *DataTableView) getFormatedRows() []string {
	d.state.columnWidths = make([]int, d.state.table.NumColumns())
	if d.state.table.NumRows() == 0 {
		return []string{}
	}
	for i := range d.state.table.Rows() {
		row := d.state.table.GetRowStrings(i)
		for columnIndex, value := range row {
			if len(value) >= maxCellWidth {
				d.state.columnWidths[columnIndex] = maxCellWidth
			}
			if len(value) > d.state.columnWidths[columnIndex] && len(value) <= maxCellWidth {
				d.state.columnWidths[columnIndex] = len(value)
			}
			if len(d.state.table.Columns()[columnIndex].Name) > d.state.columnWidths[columnIndex] {
				d.state.columnWidths[columnIndex] = len(d.state.table.Columns()[columnIndex].Name)
			}
		}
	}
	formatedRows := make([]string, d.state.table.NumRows())
	for i := range d.state.table.Rows() {
		row := d.state.table.GetRowStrings(i)
		for columnIndex, value := range row {
			cell := ""
			if len(value) > d.state.columnWidths[columnIndex] {
				cell = value[:d.state.columnWidths[columnIndex]-3] + "..."
				logger.Info(cell)
			} else {
				cell = value
			}
			formatedRows[i] += " " + cell
			formatedRows[i] += strings.Repeat(" ", d.state.columnWidths[columnIndex]-len(cell)+1)
			formatedRows[i] += "\u2502"
		}
	}
	// logger.Info("len(formatedRows) = ", len(formatedRows))
	// logger.Info("maxWidths = ", d.state.columnWidths)
	// logger.Info(formatedRows)
	return formatedRows
}

func (d *DataTableView) getFormatedTitle() string {
	formatedHeader := ""
	for columnIndex, name := range d.state.table.ColumnNames() {
		formatedHeader += "\u2500"

		formatedHeader += name
		if d.state.columnWidths[columnIndex] == 0 {
			d.state.columnWidths[columnIndex] = len(name) + 2
		}
		// logger.Info("column = ", name, ", width = ", d.state.columnWidths[columnIndex]-len(name)+1)
		formatedHeader += strings.Repeat("\u2500", d.state.columnWidths[columnIndex]-len(name)+1)
		if columnIndex == len(d.state.table.Columns())-1 {
			formatedHeader += "\u2510"
		} else {
			formatedHeader += "\u252c"
		}
	}
	d.state.contentWidth = runewidth.StringWidth(formatedHeader)
	formatedHeader += "\n"
	return formatedHeader
}
