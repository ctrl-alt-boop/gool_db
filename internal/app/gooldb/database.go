package gooldb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ctrl-alt-boop/gooldb/internal/database"
)

var DefaultFetchLimit = 10

func Connect(ip string) (*database.DatabaseContext, error) {
	connString := fmt.Sprintf("host=%s user=suncat dbname=gool_db_dev password=gool_db sslmode=disable", ip)
	log.Printf("Connecting to %s", ip)
	return createDatabaseContext(database.DriverPostgreSQL, connString)
}

func createDatabaseContext(driverName database.DriverName, connectionString string) (*database.DatabaseContext, error) {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, fmt.Errorf("err when sql.Open(): %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("err when db.Ping(): %w", err)
	}

	context := &database.DatabaseContext{
		DriverName:       driverName,
		DB:               db,
		Driver:           database.NameToDriver(driverName),
		ConnectionString: connectionString,

		FetchLimit:       DefaultFetchLimit,
		FetchLimitOffset: 0,
	}
	return context, nil
}
