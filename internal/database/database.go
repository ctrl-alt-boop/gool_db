package database

import (
	"log"

	"github.com/ctrl-alt-boop/gooldb/internal/database/connection"
)

func Connect(driver GoolDbDriver, settings connection.Settings) (*DatabaseContext, error) {
	connString := driver.ConnectionString(settings)
	log.Printf("Connecting to %s:%d", settings.Ip, settings.Port)
	return CreateDatabaseContext(settings.DriverName, connString)
}

func SetDatabase(name string) {
	go func() {

	}()
}
