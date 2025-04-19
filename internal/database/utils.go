package database

import (
	"fmt"
	"log"

	"github.com/ctrl-alt-boop/gooldb/internal/database/query"
	"github.com/google/uuid"
)

func NameToDriver(driver DriverName) GoolDbDriver {
	switch driver {
	case DriverMySql:
		return &query.MySql{}
	case DriverPostgreSQL:
		return &query.Postgres{}
	case DriverSQLite:
		return &query.SQLite3{}
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

func ResolveDatabaseType(dbType string, value []byte) (any, error) {
	log.Printf("Resolving %s", dbType)
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
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
