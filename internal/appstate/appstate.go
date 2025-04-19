package appstate

type ContentType int

const (
	Driver ContentType = iota
	Database
	Table

)

type AppLayout struct {

}

type ListPanelState struct {
	Type ContentType

}

type DataPanelState struct {

}

type HelpBarState struct {
	Hidden bool
}

type StatusBarState struct {
}

type ApplicationState struct {
	Connected bool

	Driver string
}

var AppState *ApplicationState

