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
	//Types  []string
}

var DefaultColumnMaxWidth int = 12

type DataTable struct {
	Columns []Column
	Rows    []Row

	ColumnMaxWidth int
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
		Columns:        columns,
		Rows:           make([]Row, 0),
		ColumnMaxWidth: DefaultColumnMaxWidth,
	}
}

func (dt *DataTable) AddRow(rowScan func(dest ...any) error) error {
	row := Row{
		make([]any, len(dt.Columns)),
	}
	scanArr := make([]any, len(dt.Columns))
	for i := range row.Values {
		scanArr[i] = &row.Values[i]
	}

	err := rowScan(scanArr...)
	if err != nil {
		return err
	}

	dt.Rows = append(dt.Rows, row)
	return nil
}

func (dt *DataTable) GetRowString(index int) string {
	row := make([]string, len(dt.Columns))
	for columnIndex, column := range dt.Columns {
		rowValue := dt.Rows[index].Values[columnIndex]
		switch value := rowValue.(type) {
		case string, int, float64, uint, bool:
			row[columnIndex] = fmt.Sprint(value)
		case time.Time:
			row[columnIndex] = fmt.Sprint(value.Format("2006-01-02 15:04:05.000000-07"))
		case []byte:
			resolved, err := ResolveDatabaseType(column.DbType, value)
			if err != nil {
				row[columnIndex] = err.Error()
			} else {
				row[columnIndex] = fmt.Sprintf("%T(%v)", resolved, resolved)
			}
		default:
			log.Fatalln(errors.New("unknown row value type"))
		}
		if len(row[columnIndex]) > dt.ColumnMaxWidth {
			row[columnIndex] = Abbr(row[columnIndex], dt.ColumnMaxWidth)
		}
	}
	return strings.Join(row, " | ")
}

// Names, Types, DbTypes
func (dt *DataTable) ColumnSlices() ([]string, []string, []string) {
	return dt.ColumnNames(), dt.ColumnTypeStrings(), dt.ColumnDatabaseTypeStrings()
}

func (dt *DataTable) ColumnNames() []string {
	return SliceTransform(dt.Columns, func(col Column) string {
		return col.Name
	})
}

func (dt *DataTable) ColumnTypeStrings() []string {
	return SliceTransform(dt.Columns, func(col Column) string {
		return col.ScanType.Kind().String()
		// return col.ScanType.Name()
	})
}

func (dt *DataTable) ColumnDatabaseTypeStrings() []string {
	return SliceTransform(dt.Columns, func(col Column) string {
		return col.DbType
	})
}
