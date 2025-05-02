package database

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/sql"
)

func NameToDriver(driver DriverName) DbDriver {
	switch driver {
	case DriverMySql:
		return &sql.MySql{}
	case DriverPostgreSQL:
		return &sql.Postgres{}
	case DriverSQLite:
		return &sql.SQLite3{}
	default:
		panic("Driver not implemented")
	}
}

func Abbr(s string, maxWidth int) string {
	n := (maxWidth / 2)
	return fmt.Sprintf("%s...%s", s[:n], s[len(s)-n:])
}

func FirstN(s string, n int) string {
	return fmt.Sprintf("%s...", s[:n])
}

func LastN(s string, n int) string {
	return fmt.Sprintf("...%s", s[len(s)-n:])
}
