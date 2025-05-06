package sql

import (
	"fmt"
	"plugin"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	"github.com/google/uuid"
)

type Postgres struct{}

func (d *Postgres) ConnectionString(settings *connection.Settings) string {
	connString := ""
	connString += fmt.Sprintf("host=%s ", settings.Ip)
	if settings.Port == 0 {
		settings.Port = 5432
	}
	connString += fmt.Sprintf("port=%d ", settings.Port)
	connString += fmt.Sprintf("user=%s ", settings.Username)
	connString += fmt.Sprintf("password=%s ", settings.Password)
	if settings.DbName != "" {
		connString += fmt.Sprintf("dbname=%s ", settings.DbName)
	}
	if settings.SslMode == "" {
		settings.SslMode = "disable"
	}
	connString += fmt.Sprintf("sslmode=%s ", settings.SslMode)

	return connString
}

func (d Postgres) Load() error {
	plug, err := plugin.Open("./plugins/postgres.so")
	if err != nil {
		return err
	}

	_, err = plug.Lookup("Loaded")
	if err != nil {
		return err
	}
	return nil
}

func (d *Postgres) DatabasesQuery() string {
	return "SELECT datname FROM pg_database WHERE datistemplate = false"
}

func (d *Postgres) DatabaseNameQuery() string {
	return "SELECT current_database()"
}

func (d *Postgres) TableNamesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"
}

func (d *Postgres) CountQuery(table string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
}

func (d *Postgres) SelectAllQuery(table string, opts query.Option) string {
	queryOptsString := buildQueryOptions(opts)
	return fmt.Sprintf("SELECT * FROM %s%s", table, queryOptsString)
}

func (d *Postgres) ResolveDatabaseType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}

func buildQueryOptions(opts query.Option) string {
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
