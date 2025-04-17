package query

import (
	_ "github.com/go-sql-driver/mysql"
)

type MySql struct{}

// SelectAllQuery implements database.GoolDbDriver.
func (d *MySql) SelectAllQuery(table string, opts QueryOptions) string {
	panic("unimplemented")
}

func (d *MySql) CountQuery(table string) string {
	panic("unimplemented")
}

func (d *MySql) DatabaseNameQuery() string {
	panic("unimplemented")
}

func (d *MySql) TableNamesQuery() string {
	return "SHOW TABLES"
}
