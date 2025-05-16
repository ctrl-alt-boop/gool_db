package widgets

import "github.com/ctrl-alt-boop/gooldb/pkg/logging"

var logger = logging.NewLogger("widgets.log")

var Names = struct {
	CommandBar string
	DataArea   string
	Panel      string
	Help       string

	QueryOptionsPopup string
	TableCellPopup    string
}{
	CommandBar: "commandBar",
	DataArea:   "dataArea",
	Panel:      "panel",
	Help:       "help",

	QueryOptionsPopup: "queryOptionsPopup",
	TableCellPopup:    "tableCellPopup",
}
