package sql

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	_ "github.com/go-sql-driver/mysql"
)

type MySql struct{}

// SupportsJsonResult implements database.DbDriver.
func (d *MySql) SupportsJsonResult() bool {
	panic("unimplemented")
}

func (d *MySql) ConnectionString(settings *connection.Settings) string {
	panic("unimplemented")
}

func (d *MySql) ResolveDatabaseType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}

func (d MySql) Load() error {
	panic("unimplemented")
}

func (d *MySql) DatabasesQuery() string {
	panic("unimplemented")
}

func (d *MySql) DatabaseNameQuery() string {
	panic("unimplemented")
}

func (d *MySql) TableNamesQuery() string {
	return "SHOW TABLES"
}

func (d *MySql) CountQuery(table string) string {
	panic("unimplemented")
}

// SelectAllQuery implements database.GoolDbDriver.
func (d *MySql) SelectAllQuery(table string, opts query.Statement) string {
	panic("unimplemented")
}
