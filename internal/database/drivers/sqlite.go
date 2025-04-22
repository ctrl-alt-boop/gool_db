package drivers

import (
	"github.com/ctrl-alt-boop/gooldb/internal/database/connection"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite3 struct{}

// ConnectionString implements database.GoolDbDriver.
func (d *SQLite3) ConnectionString(settings connection.Settings) string {
	panic("unimplemented")
}

// ResolveDatabaseType implements database.GoolDbDriver.
func (d *SQLite3) ResolveDatabaseType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}

// Load implements database.GoolDbDriver.
func (d SQLite3) Load() error {
	panic("unimplemented")
}

// DatabasesQuery implements database.GoolDbDriver.
func (d *SQLite3) DatabasesQuery() string {
	panic("unimplemented")
}

func (d *SQLite3) DatabaseNameQuery() string {
	panic("unimplemented")
}

func (d *SQLite3) TableNamesQuery() string {
	return "SELECT name FROM sqlite_master WHERE type='table';"
}

func (d *SQLite3) CountQuery(table string) string {
	panic("unimplemented")
}

// SelectAllQuery implements database.GoolDbDriver.
func (d *SQLite3) SelectAllQuery(table string, opts QueryOptions) string {
	panic("unimplemented")
}
