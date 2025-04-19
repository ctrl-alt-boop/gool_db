package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ctrl-alt-boop/gooldb/internal/database/query"
)

type DriverName = string

const (
	DriverMySql      DriverName = "mysql"
	DriverPostgreSQL DriverName = "postgres"
	DriverSQLite     DriverName = "sqlite3"
)

type DatabaseContext struct {
	Driver           GoolDbDriver
	DB               *sql.DB
	DriverName       DriverName
	ConnectionString string

	FetchLimit       int
	FetchLimitOffset int
}

func (context *DatabaseContext) FetchTable(selectedTable string) *DataTable {
	opts := query.QueryOptions{Limit: context.FetchLimit, Offset: context.FetchLimitOffset}
	rows, err := context.DB.Query(context.Driver.SelectAllQuery(selectedTable, opts))
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	columns, err := rows.ColumnTypes()
	if err != nil {
		log.Panicln(err)
	}
	table := CreateDataTable(columns)
	for rows.Next() {
		table.AddRow(rows.Scan)
		log.Println("added row")
	}

	// context.FetchLimitOffset += context.FetchLimit
	return table
}

func (context *DatabaseContext) FetchDatabaseName() string {
	var dbName string
	err := context.DB.QueryRow(context.Driver.DatabaseNameQuery()).Scan(&dbName)
	if err != nil {
		log.Panicln(err)
	}
	return dbName
}

func (context *DatabaseContext) FetchTableNames() []string {
	var tableNames []string
	rows, err := context.DB.Query(context.Driver.TableNamesQuery())
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Panicln(err)
		}
		tableNames = append(tableNames, tableName)
	}

	return tableNames
}

func (context *DatabaseContext) FetchCount(table string) int {
	var count int
	err := context.DB.QueryRow(context.Driver.CountQuery(table)).Scan(&count)
	if err != nil {
		log.Println(err)
		return 0
	}
	return count
}

func (context *DatabaseContext) FetchCounts(tables []string) []string {
	for index, table := range tables {
		count := context.FetchCount(table)
		tables[index] = fmt.Sprintf("%s (%d entries)", tables[index], count)
	}
	return tables
}
