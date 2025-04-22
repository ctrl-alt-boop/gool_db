package drivers

import (
	"github.com/ctrl-alt-boop/gooldb/internal/database/connection"
	_ "github.com/go-sql-driver/mysql"
)

type MySql struct{}

// ConnectionString implements database.GoolDbDriver.
func (d *MySql) ConnectionString(settings connection.Settings) string {
	panic("unimplemented")
}

// ResolveDatabaseType implements database.GoolDbDriver.
func (d *MySql) ResolveDatabaseType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}

// Load implements database.GoolDbDriver.
func (d MySql) Load() error {
	panic("unimplemented")
}

// DatabasesQuery implements database.GoolDbDriver.
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
func (d *MySql) SelectAllQuery(table string, opts QueryOptions) string {
	panic("unimplemented")
}
