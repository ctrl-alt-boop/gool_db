package gooldb

import (
	"log"

	"github.com/ctrl-alt-boop/gooldb/internal/database"
	"github.com/ctrl-alt-boop/gooldb/internal/database/connection"
)

const (
	_user string = "valmatics"
	_pass string = "valmatics"
)

type GoolDb struct {
	ip              string
	DatabaseContext *database.DatabaseContext
	settings        connection.Settings
}

func Create(ip string) *GoolDb {
	return &GoolDb{
		ip: ip,
	}
}

func (gool *GoolDb) SetDriver(name database.DriverName) error {
	driver := database.NameToDriver(name)
	err := driver.Load()
	if err != nil {
		return err
	}
	gool.settings = connection.NewSettings(
		connection.WithDriver(name),
		connection.WithIp("localhost"),
		connection.WithUser(_user),
		connection.WithPass(_pass),
		connection.WithSslMode("disable"),
	)
	context, err := database.Connect(driver, gool.settings)
	if err != nil {
		log.Panicln(err, "Ip:", gool.settings.Ip, "Port:", gool.settings.Port)
	}

	gool.DatabaseContext = context
	return nil
}

func (gool *GoolDb) SetDatabase(name string) error {
	gool.settings.DbName = name
	if gool.DatabaseContext != nil {
		gool.DatabaseContext.DB.Close()
	}
	context, err := database.Connect(gool.DatabaseContext.Driver, gool.settings)
	if err != nil {
		return err
	}
	gool.DatabaseContext = context
	return nil
}

func (gool *GoolDb) SetTable(name string) error {
	return nil
}

func (gool *GoolDb) FetchCounts(tables []string) []string {
	return gool.DatabaseContext.FetchCounts(tables)
}

func (gool *GoolDb) FetchTable(selection string) *database.DataTable {
	return gool.DatabaseContext.FetchTable(selection)
}
