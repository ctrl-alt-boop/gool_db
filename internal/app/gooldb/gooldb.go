package gooldb

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database"
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

const (
	_user string = "valmatics"
	_pass string = "valmatics"
)

var SupportedDrivers []string = []string{
	database.DriverPostgreSQL,
	database.DriverMySql,
	database.DriverSQLite,
}

var logger *logging.Logger

type GoolDb struct {
	logger   *logging.Logger
	notifier Notifier
	ip       string

	settings        *connection.Settings
	databaseContext *database.DatabaseContext

	eventHandlers map[EventType]EventHandler
}

func Create(log *logging.Logger, notifier Notifier, ip string) *GoolDb {
	logger = log
	return &GoolDb{
		logger:        log,
		notifier:      notifier,
		ip:            ip,
		eventHandlers: make(map[EventType]EventHandler),
	}
}

func (gool *GoolDb) RegisterEventHandler(eventType EventType, handler EventHandler) {
	gool.eventHandlers[eventType] = handler
}

func (gool *GoolDb) SelectDriver(name database.DriverName) {
	handler, exists := gool.eventHandlers[DriverSet]
	if !exists {
		logger.Panic("no handler registered for DriverSet event")
	}
	go func() {
		driver := database.NameToDriver(name)
		err := driver.Load()
		if err != nil {
			logger.Warn(err)
			gool.notifier.Notify(Warning, err)
			return
		}
		if gool.databaseContext != nil {
			gool.databaseContext.DB.Close()
		}
		gool.settings, err = connection.NewSettings(
			connection.WithDriver(name),
			connection.WithIp(gool.ip),
			connection.WithUser(_user),
			connection.WithPass(_pass),
			connection.WithSslMode("disable"),
		)
		if err != nil {
			logger.Warn(err)
			gool.notifier.Notify(Warning, err)
			return
		}
		context, err := database.Connect(driver, gool.settings)
		if err != nil {
			logger.Warn(err, "ip:", gool.settings.Ip, "port:", gool.settings.Port)
			gool.notifier.Notify(Warning, err)
			return
		}

		gool.databaseContext = context
		databases, err := gool.databaseContext.FetchDatabases()
		eventArgs := DriverSetEvent{
			Selected:  name,
			Databases: databases,
		}
		handler(eventArgs, err)
	}()
}

func (gool *GoolDb) SelectDatabase(name string) {
	handler, exists := gool.eventHandlers[DatabaseSet]
	if !exists {
		logger.Panic("no handler registered for DatabaseSet event")
	}
	go func() {
		gool.settings.DbName = name
		if gool.databaseContext != nil {
			gool.databaseContext.DB.Close()
		}
		context, err := database.Connect(gool.databaseContext.Driver, gool.settings)
		if err != nil {
			logger.Warn(err)
			gool.notifier.Notify(Warning, err)
			return
		}
		gool.databaseContext = context
		tables, err := gool.databaseContext.FetchTableNames()
		eventArgs := DatabaseSetEvent{
			Selected: name,
			Tables:   tables,
		}
		handler(eventArgs, err)
	}()
}

func (gool *GoolDb) SelectTable(name string) {
	handler, exists := gool.eventHandlers[TableSet]
	if !exists {
		logger.Panic("no handler registered for TableSet event")
	}
	go func() {
		table, err := gool.FetchTable(name)
		eventArgs := TableSetEvent{
			Selected: name,
			Table:    table,
		}
		handler(eventArgs, err)
	}()
}

func (gool *GoolDb) FetchCounts(tables []string) []string {
	return gool.databaseContext.FetchCounts(tables)
}

func (gool *GoolDb) FetchTable(selection string) (*DataTable, error) {
	columns, rows, err := gool.databaseContext.FetchTable(selection)
	return CreateDataTable(columns, rows), err
}
