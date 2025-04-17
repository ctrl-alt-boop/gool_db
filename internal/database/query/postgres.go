package query

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct{}

func (d *Postgres) SelectAllQuery(table string, opts QueryOptions) string {
	queryOptsString := buildQueryOptions(opts)
	return fmt.Sprintf("SELECT * FROM %s%s", table, queryOptsString)
}

func (d *Postgres) CountQuery(table string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
}

func (d *Postgres) DatabaseNameQuery() string {
	return "SELECT current_database()"
}

func (d *Postgres) TableNamesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"
}

func buildQueryOptions(opts QueryOptions) string {
	sb := strings.Builder{}

	// TODO: ADD WHERE THING HERE

	if opts.OrderBy.Column != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(opts.OrderBy.Column)
		if opts.OrderBy.Desc {
			sb.WriteString(" DESC")
		} else {
			sb.WriteString(" ASC")
		}
	}

	if opts.Limit > 0 {
		sb.WriteString(fmt.Sprint(" LIMIT ", opts.Limit))
	}
	if opts.Offset > 0 {
		sb.WriteString(fmt.Sprint(" OFFSET ", opts.Offset))
	}

	return sb.String()
}
