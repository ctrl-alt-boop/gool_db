package database

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
)

type DriverName = string

const (
	DriverPostgreSQL DriverName = "postgres"
	DriverMySql      DriverName = "mysql"
	DriverSQLite     DriverName = "sqlite3"
)

type DbDriver interface {
	SupportsJsonResult() bool

	Load() error
	ConnectionString(settings *connection.Settings) string

	// CustomQuery(opts query.Statement)

	DatabasesQuery() string
	DatabaseNameQuery() string
	TableNamesQuery() string
	CountQuery(table string) string
	SelectAllQuery(table string, opts query.Statement) string
	ResolveDatabaseType(dbType string, value []byte) (any, error)
}
