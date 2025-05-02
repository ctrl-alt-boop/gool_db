package database

import (
	"database/sql"
	"fmt"
	"reflect"
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

func ParseRows(driver DbDriver, dbRows *sql.Rows) ([]Column, []Row, error) {
	dbColumns, err := dbRows.ColumnTypes()
	if err != nil {
		logger.Warn(err)
		return nil, nil, err
	}
	columns := make([]Column, len(dbColumns))
	for i := range dbColumns {
		columns[i] = Column{
			Name:     dbColumns[i].Name(),
			ScanType: dbColumns[i].ScanType(),
			DbType:   dbColumns[i].DatabaseTypeName(),
		}
	}

	rows := make([]Row, 0)
	for dbRows.Next() {
		row := Row{
			make([]any, len(dbColumns)),
		}
		scanArr := make([]any, len(dbColumns))
		for i := range row.Values {
			scanArr[i] = &row.Values[i]
		}

		err := dbRows.Scan(scanArr...)
		if err != nil {
			logger.Warn(err)
			continue
		}
		for i := range row.Values {
			row.Values[i], err = ResolveTypes(driver, row.Values[i], columns[i])
			if err != nil {
				logger.Warn(err)
				continue
			}
		}

		rows = append(rows, row)
	}
	logger.Info(columns)
	return columns, rows, nil
}

func ResolveTypes(driver DbDriver, rowValue any, column Column) (any, error) {
	switch value := rowValue.(type) {
	case string, int, int32, int64, float32, float64, uint, bool:
		return value, nil
	case time.Time:
		return fmt.Sprint(value.Format("2006-01-02 15:04:05.000000-07")), nil
	case []byte:
		resolved, err := driver.ResolveDatabaseType(column.DbType, value)
		//logger.Info("resolving type", driver, column, string(value))
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", resolved), nil
		}
	case nil:
		return "null", nil
	default:
		err := fmt.Errorf("unknown value type %T for %s", value, column.DbType)
		logger.Info(err)
		return "", err
	}
}
