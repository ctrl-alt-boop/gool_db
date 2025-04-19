package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type Column struct {
	Name     string
	ScanType reflect.Type
	DbType   string
}

type Row struct {
	Values []any
}

type DataTable struct {
	columns []Column
	rows    []Row
}

func (dt *DataTable) NumColumns() int {
	return len(dt.columns)
}

func (dt *DataTable) NumRows() int {
	return len(dt.rows)
}

func (dt *DataTable) Columns() []Column {
	return dt.columns
}

func (dt *DataTable) Rows() []Row {
	return dt.rows
}

func CreateDataTable(columnTypes []*sql.ColumnType) *DataTable {
	columns := SliceTransform(columnTypes, func(columnType *sql.ColumnType) Column {
		return Column{
			Name:     columnType.Name(),
			ScanType: columnType.ScanType(),
			DbType:   columnType.DatabaseTypeName(),
		}
	})

	return &DataTable{
		columns:        columns,
		rows:           make([]Row, 0),
	}
}

func (dt *DataTable) AddRow(rowScan func(dest ...any) error) error {
	row := Row{
		make([]any, len(dt.columns)),
	}
	scanArr := make([]any, len(dt.columns))
	for i := range row.Values {
		scanArr[i] = &row.Values[i]
	}

	err := rowScan(scanArr...)
	if err != nil {
		return err
	}

	dt.rows = append(dt.rows, row)
	return nil
}

func (dt *DataTable) GetRowString(index int) string {
	row := make([]string, len(dt.columns))
	for columnIndex, column := range dt.columns {
		rowValue := dt.rows[index].Values[columnIndex]
		switch value := rowValue.(type) {
		case string, int, int32, int64, float32, float64, uint, bool:
			row[columnIndex] = fmt.Sprint(value)
		case time.Time:
			row[columnIndex] = fmt.Sprint(value.Format("2006-01-02 15:04:05.000000-07"))
		case []byte:
			resolved, err := ResolveDatabaseType(column.DbType, value)
			if err != nil {
				row[columnIndex] = err.Error()
			} else {
				row[columnIndex] = fmt.Sprintf("%v", resolved)
			}
		case nil:
			row[columnIndex] = "null"
		default:
			log.Fatalln(errors.New(fmt.Sprintf("unknown value type %T for %s", value, column.DbType)))
		}
	}
	return strings.Join(row, " | ")
}

func (dt *DataTable) GetColumnRows(columnIndex int) ([]string, int) {
	rows := make([]string, dt.NumRows())
	var columnWidth int
	for rowIndex, row := range dt.rows {
		rows[rowIndex] = fmt.Sprint(row.Values[columnIndex])
		columnWidth = max(columnWidth, len(rows[rowIndex]))
	}
	return rows, columnWidth
}

// Names, Types, DbTypes
func (dt *DataTable) ColumnSlices() ([]string, []string, []string) {
	return dt.ColumnNames(), dt.ColumnTypeStrings(), dt.ColumnDatabaseTypeStrings()
}

func (dt *DataTable) ColumnNames() []string {
	return SliceTransform(dt.columns, func(col Column) string {
		return col.Name
	})
}

func (dt *DataTable) ColumnTypeStrings() []string {
	return SliceTransform(dt.columns, func(col Column) string {
		return col.ScanType.Kind().String()
	})
}

func (dt *DataTable) ColumnDatabaseTypeStrings() []string {
	return SliceTransform(dt.columns, func(col Column) string {
		return col.DbType
	})
}

func (dt *DataTable) ClearRows() error {
	dt.rows = make([]Row, 0)
	return nil
}
