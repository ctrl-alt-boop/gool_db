package database

import (
	go_sql "database/sql"
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/pkg/query"
)

var DefaultFetchLimit = 10

type DatabaseContext struct {
	Driver           DbDriver
	DB               *go_sql.DB
	DriverName       DriverName
	ConnectionString string

	FetchLimit       int
	FetchLimitOffset int
}

func CreateDatabaseContext(driverName DriverName, connectionString string) (*DatabaseContext, error) {
	db, err := go_sql.Open(driverName, connectionString)
	if err != nil {
		return nil, fmt.Errorf("err when go_sql.Open(): %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("err when db.Ping(): %w", err)
	}

	context := &DatabaseContext{
		DriverName:       driverName,
		DB:               db,
		Driver:           NameToDriver(driverName),
		ConnectionString: connectionString,

		FetchLimit:       DefaultFetchLimit,
		FetchLimitOffset: 0,
	}
	return context, nil
}

func (context *DatabaseContext) FetchDatabases() ([]string, error) {
	var databases []string
	rows, err := context.DB.Query(context.Driver.DatabasesQuery())
	if err != nil {
		return nil, fmt.Errorf("err when context.DB.Query(): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var databaseName string
		if err := rows.Scan(&databaseName); err != nil {
			logger.Warn(err)
		}
		databases = append(databases, databaseName)
	}

	return databases, nil
}

func (context *DatabaseContext) FetchDatabaseName() (string, error) {
	var dbName string
	err := context.DB.QueryRow(context.Driver.DatabaseNameQuery()).Scan(&dbName)
	if err != nil {
		return "", fmt.Errorf("err when context.DB.QueryRow(): %w", err)
	}
	return dbName, nil
}

func (context *DatabaseContext) FetchTableNames() ([]string, error) {
	var tableNames []string
	rows, err := context.DB.Query(context.Driver.TableNamesQuery())
	if err != nil {
		return nil, fmt.Errorf("err when context.DB.Query(): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			logger.Warn(err)
		}
		tableNames = append(tableNames, tableName)
	}

	return tableNames, nil
}

func (context *DatabaseContext) FetchCount(table string) (int, error) {
	var count int
	err := context.DB.QueryRow(context.Driver.CountQuery(table)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("err when context.DB.QueryRow(): %w", err)
	}
	return count, nil
}

func (context *DatabaseContext) FetchCounts(tables []string) []string {
	for index, table := range tables {
		count, err := context.FetchCount(table)
		if err != nil {
			logger.Warn(err)
			tables[index] = fmt.Sprintf("%s (%s))", tables[index], err.Error())
			continue
		}
		tables[index] = fmt.Sprintf("%s (%d entries)", tables[index], count)
	}
	return tables
}

func (context *DatabaseContext) FetchTable(selectedTable string) ([]Column, []Row, error) { // context.FetchLimitOffset += context.FetchLimit
	opts := query.Statement{Limit: context.FetchLimit, Offset: context.FetchLimitOffset}
	dbRows, err := context.DB.Query(context.Driver.SelectAllQuery(selectedTable, opts))
	if err != nil {
		logger.Warn(err)
		return nil, nil, err
	}
	defer dbRows.Close()

	return ParseRows(context.Driver, dbRows)

	// dbColumns, err := dbRows.ColumnTypes()
	// if err != nil {
	// 	logger.Panic(err)
	// }
	// columns := make([]Column, len(dbColumns))
	// for i := range dbColumns {
	// 	columns[i] = Column{
	// 		Name:     dbColumns[i].Name(),
	// 		ScanType: dbColumns[i].ScanType(),
	// 		DbType:   dbColumns[i].DatabaseTypeName(),
	// 	}
	// }

	// rows := make([]Row, 0)
	// for dbRows.Next() {
	// 	row := Row{
	// 		make([]any, len(dbColumns)),
	// 	}
	// 	scanArr := make([]any, len(dbColumns))
	// 	for i := range row.Values {
	// 		scanArr[i] = &row.Values[i]
	// 	}

	// 	err := dbRows.Scan(scanArr...)
	// 	if err != nil {
	// 		logger.Warn(err)
	// 		continue
	// 	}

	// 	rows = append(rows, row)
	// }
	// return columns, rows
}
