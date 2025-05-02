package gooldb

import (
	"fmt"
	"strings"
	"time"

	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database"
	shared "github.com/ctrl-alt-boop/gooldb/pkg"
)

const maxCellWidth = 36 // Guid length, including the '-'s

type DataTable struct {
	columns []database.Column
	rows    []database.Row

	Driver database.DbDriver
}

func (dt *DataTable) NumColumns() int {
	return len(dt.columns)
}

func (dt *DataTable) NumRows() int {
	return len(dt.rows)
}

func (dt *DataTable) Columns() []database.Column {
	return dt.columns
}

func (dt *DataTable) Rows() []database.Row {
	return dt.rows
}

func CreateDataTable(columns []database.Column, rows []database.Row) *DataTable {
	return &DataTable{
		columns: columns,
		rows:    rows,
	}
}

func (dt *DataTable) GetRowColumn(row, column int) (string, error) { // Needs a lookielook to see if other drivers are at least similar to this
	rowColumn := dt.rows[row].Values[column]
	switch value := rowColumn.(type) {
	case string, int, int32, int64, float32, float64, uint, bool:
		return fmt.Sprint(value), nil
	case time.Time:
		return fmt.Sprint(value.Format("2006-01-02 15:04:05.000000-07")), nil
	case []byte:
		logger.Info(dt.Driver, dt.columns, value)
		resolved, err := dt.Driver.ResolveDatabaseType(dt.columns[column].DbType, value)
		if err != nil {
			logger.Warn(err)
			return "", err
		} else {
			return fmt.Sprintf("%v", resolved), nil
		}
	case nil:
		return "null", nil
	default:
		err := fmt.Errorf("unknown value type %T for %s", value, dt.columns[column].DbType)
		logger.Info(err)
		return "", err
	}
}

func (dt *DataTable) GetRowStrings(index int) []string {
	row := make([]string, len(dt.columns))
	for columnIndex := range dt.columns {
		value, err := dt.GetRowColumn(index, columnIndex)
		if err != nil {
			row[columnIndex] = err.Error()
		} else {
			row[columnIndex] = value
		}
	}
	return row
}

func (dt *DataTable) GetRowString(index int) string {
	row := dt.GetRowStrings(index)
	return strings.Join(row, " | ")
}

func (dt *DataTable) GetColumnRows(columnIndex int) ([]string, int) {
	rows := make([]string, dt.NumRows())
	var columnWidth int
	for rowIndex := range dt.rows {
		value, err := dt.GetRowColumn(rowIndex, columnIndex)
		if err != nil {
			rows[rowIndex] = err.Error()
		} else {
			rows[rowIndex] = value
		}

		columnWidth = max(columnWidth, len(rows[rowIndex]))
	}
	return rows, columnWidth
}

// Names, Types, DbTypes
func (dt *DataTable) ColumnSlices() ([]string, []string, []string) {
	return dt.ColumnNames(), dt.ColumnTypeStrings(), dt.ColumnDatabaseTypeStrings()
}

func (dt *DataTable) ColumnNames() []string {
	return shared.SliceTransform(dt.columns, func(col database.Column) string {
		return col.Name
	})
}

func (dt *DataTable) ColumnTypeStrings() []string {
	return shared.SliceTransform(dt.columns, func(col database.Column) string {
		return col.ScanType.Kind().String()
	})
}

func (dt *DataTable) ColumnDatabaseTypeStrings() []string {
	return shared.SliceTransform(dt.columns, func(col database.Column) string {
		return col.DbType
	})
}

func (dt *DataTable) ClearRows() error {
	dt.rows = make([]database.Row, 0)
	return nil
}

func (dt *DataTable) GetFormatedRows() ([]string, []int) {
	columnWidths := make([]int, dt.NumColumns())
	if dt.NumRows() == 0 {
		return []string{}, columnWidths
	}
	for i := range dt.rows {
		row := dt.GetRowStrings(i)
		for columnIndex, value := range row {
			if len(value) >= maxCellWidth {
				columnWidths[columnIndex] = maxCellWidth
			}
			if len(value) > columnWidths[columnIndex] && len(value) <= maxCellWidth {
				columnWidths[columnIndex] = len(value)
			}
		}
	}
	formatedRows := make([]string, dt.NumRows())
	for i := range dt.rows {
		row := dt.GetRowStrings(i)
		for columnIndex, value := range row {
			cell := ""
			if len(value) > columnWidths[columnIndex] {
				cell = value[:columnWidths[columnIndex]-3] + "..."
				logger.Info(cell)
			} else {
				cell = value
			}
			formatedRows[i] += " " + cell
			formatedRows[i] += strings.Repeat(" ", columnWidths[columnIndex]-len(cell)+1)
			formatedRows[i] += "\u2502"
		}
	}
	logger.Info("len(formatedRows) = ", len(formatedRows))
	logger.Info("maxWidths = ", columnWidths)
	logger.Info(formatedRows)
	return formatedRows, columnWidths
}

func (dt *DataTable) GetFormatedTitle(columnWidths []int) string {
	formatedHeader := ""
	for columnIndex, name := range dt.ColumnNames() {
		if columnIndex > 0 {
			formatedHeader += "\u2500"
		}
		formatedHeader += name
		if columnWidths[columnIndex] == 0 {
			columnWidths[columnIndex] = len(name) + 2
		}
		logger.Info("column = ", name, ", width = ", columnWidths[columnIndex]-len(name))
		formatedHeader += strings.Repeat("\u2500", columnWidths[columnIndex]-len(name)+1)
		formatedHeader += "\u252c"
	}
	return formatedHeader
}
