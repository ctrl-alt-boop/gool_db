package database

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("db.log")

func Connect(driver DbDriver, settings *connection.Settings) (*DatabaseContext, error) {
	connString := driver.ConnectionString(settings)
	logger.Info("Connecting: ", settings.Ip, " ", settings.Port, " ", settings.DbName)
	return CreateDatabaseContext(settings.DriverName, connString)
}
