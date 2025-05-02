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
	Load() error
	ConnectionString(settings *connection.Settings) string
	DatabasesQuery() string
	DatabaseNameQuery() string
	TableNamesQuery() string
	CountQuery(table string) string
	SelectAllQuery(table string, opts query.Option) string
	ResolveDatabaseType(dbType string, value []byte) (any, error)
}
