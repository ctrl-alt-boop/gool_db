package database

import "github.com/ctrl-alt-boop/gooldb/internal/database/query"

type GoolDbDriver interface {
	TableNamesQuery() string
	DatabaseNameQuery() string
	CountQuery(table string) string
	SelectAllQuery(table string, opts query.QueryOptions) string
}
