package gooldb

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database"
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var sqlLogger = logging.NewLogger("sql.log")

const (
	_user string = "valmatics"
	_pass string = "valmatics"
)

var supportedDrivers []string = []string{
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

	eventHandlers map[EventType][]EventHandler
}

func (gool *GoolDb) GetDrivers() []string {
	return supportedDrivers
}

func Create(log *logging.Logger, notifier Notifier, ip string) *GoolDb {
	logger = log
	return &GoolDb{
		logger:        log,
		notifier:      notifier,
		ip:            ip,
		eventHandlers: make(map[EventType][]EventHandler),
	}
}

func (gool *GoolDb) RegisterEventHandler(eventType EventType, handler EventHandler) {
	handlers, exists := gool.eventHandlers[eventType]
	if !exists {
		handlers = make([]EventHandler, 0)
	}
	handlers = append(handlers, handler)
	gool.eventHandlers[eventType] = handlers
}

func (gool *GoolDb) SelectDriver(name database.DriverName) {
	handlers := gool.eventHandlers[DriverSet]
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
			connection.WithDb("postgres"),
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
			logger.Warn(err, " ip: ", gool.settings.Ip, " port: ", gool.settings.Port)
			gool.notifier.Notify(Warning, err)
			return
		}
		context.OnQueryExecuted(gool.onQueryExecuted)
		gool.databaseContext = context
		databases, err := gool.databaseContext.FetchDatabases()
		eventArgs := DriverSetEvent{
			Selected:  name,
			Databases: databases,
		}
		for _, handler := range handlers {
			handler(eventArgs, err)
		}
	}()
}

func (gool *GoolDb) onQueryExecuted(query string, err error) {
	sqlLogger.Info(query)
	if err != nil {
		sqlLogger.Error("\t", err)
	}
}

func (gool *GoolDb) SelectDatabase(name string) {
	handlers, exists := gool.eventHandlers[DatabaseSet]
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
		context.OnQueryExecuted(gool.onQueryExecuted)
		gool.databaseContext = context
		tables, err := gool.databaseContext.FetchTableNames()
		eventArgs := DatabaseSetEvent{
			Selected: name,
			Tables:   tables,
		}
		for _, handler := range handlers {
			handler(eventArgs, err)
		}
	}()
}

func (gool *GoolDb) SelectTable(name string) {
	handlers, exists := gool.eventHandlers[TableSet]
	if !exists {
		logger.Panic("no handler registered for TableSet event")
	}
	go func() {
		table, err := gool.FetchTable(name)
		eventArgs := TableSetEvent{
			Selected: name,
			Table:    table,
		}
		for _, handler := range handlers {
			handler(eventArgs, err)
		}
	}()
}

func (gool *GoolDb) FetchCounts(tables []string) map[string]int {
	return gool.databaseContext.FetchCounts(tables)
}

func (gool *GoolDb) FetchTable(selection string) (*DataTable, error) {
	columns, rows, err := gool.databaseContext.FetchTable(selection)
	return CreateDataTable(columns, rows), err
}
