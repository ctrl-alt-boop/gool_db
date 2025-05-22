package gooldb

import (
	"fmt"
	"strings"
	"time"

	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database"
	shared "github.com/ctrl-alt-boop/gooldb/pkg"
)

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
		logger.Warn(err)
		return "", err
	}
}

func (dt *DataTable) GetRowStringsAll() [][]string {
	rows := make([][]string, len(dt.columns))
	for i := range dt.NumRows() {
		rows[i] = dt.GetRowStrings(i)
	}

	return rows
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

func (dt *DataTable) GetColumnRows(columnIndex int) (rows []string, columnWidth int) {
	columnRows := make([]string, dt.NumRows())
	for rowIndex := range dt.rows {
		value, err := dt.GetRowColumn(rowIndex, columnIndex)
		if err != nil {
			columnRows[rowIndex] = err.Error()
		} else {
			columnRows[rowIndex] = value
		}

		columnWidth = max(columnWidth, len(columnRows[rowIndex]))
	}
	return columnRows, columnWidth
}

func (dt *DataTable) ColumnSlices() (names []string, types []string, dbTypes []string) {
	names = dt.ColumnNames()
	types = dt.ColumnTypeStrings()
	dbTypes = dt.ColumnDatabaseTypeStrings()
	return
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
