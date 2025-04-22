package database

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/database/drivers"
)

func NameToDriver(driver DriverName) GoolDbDriver {
	switch driver {
	case DriverMySql:
		return &drivers.MySql{}
	case DriverPostgreSQL:
		return &drivers.Postgres{}
	case DriverSQLite:
		return &drivers.SQLite3{}
	default:
		panic("Driver not implemented")
	}
}

func Limit(count int) string {
	return fmt.Sprintf("LIMIT %d", count)
}

func Offset(count int) string {
	return fmt.Sprintf("OFFSET %d", count)
}

func SliceTransform[T any, U any](slice []T, selector func(T) U) []U {
	results := make([]U, len(slice))
	for i, value := range slice {
		results[i] = selector(value)
	}
	return results
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
