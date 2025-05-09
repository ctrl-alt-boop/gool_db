package gooldb

type EventType string

const (
	Startup     EventType = "Startup"
	DriverSet   EventType = "DriverSet"
	DatabaseSet EventType = "DatabaseSet"
	TableSet    EventType = "TableSet"

	RowSet    EventType = "RowSet"
	ColumnSet EventType = "ColumnSet"
)

type EventHandler func(any, error)

type DriverSetEvent struct {
	Selected  string
	Databases []string
}

type DatabaseSetEvent struct {
	Selected string
	Tables   []string
}

type TableSetEvent struct {
	Selected string
	Count    int
	Table    *DataTable
}

type CellSelectedEvent struct {
	Selected string
	Value    string
}
