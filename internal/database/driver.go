package database

import (
	"github.com/ctrl-alt-boop/gooldb/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/internal/database/drivers"
)

var SupportedDrivers []string = []string{
	DriverMySql,
	DriverPostgreSQL,
	DriverSQLite,
}

type DriverName = string

const (
	DriverMySql      DriverName = "mysql"
	DriverPostgreSQL DriverName = "postgres"
	DriverSQLite     DriverName = "sqlite3"
)

type GoolDbDriver interface {
	Load() error
	ConnectionString(settings connection.Settings) string
	DatabasesQuery() string
	DatabaseNameQuery() string
	TableNamesQuery() string
	CountQuery(table string) string
	SelectAllQuery(table string, opts drivers.QueryOptions) string
	ResolveDatabaseType(dbType string, value []byte) (any, error)
}
