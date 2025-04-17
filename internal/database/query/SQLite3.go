package query

import (
	_ "github.com/mattn/go-sqlite3"
)

type SQLite3 struct{}

// SelectAllQuery implements database.GoolDbDriver.
func (d *SQLite3) SelectAllQuery(table string, opts QueryOptions) string {
	panic("unimplemented")
}

func (d *SQLite3) CountQuery(table string) string {
	panic("unimplemented")
}

func (d *SQLite3) DatabaseNameQuery() string {
	panic("unimplemented")
}

func (d *SQLite3) TableNamesQuery() string {
	return "SELECT name FROM sqlite_master WHERE type='table';"
}
