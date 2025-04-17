package managers

import "reflect"

type Widget struct {
	Name        string
	X0          int
	X1          int
	Y0          int
	Y1          int
	ManagerType reflect.Type // maybe, i don't even know
}

type Layout struct {
	Main      Widget
	Help      Widget
	StatusBar Widget
	TableList Widget
}

type Application struct {
	Layout Layout
}

var App = Application{
	Layout: Layout{
		Main: Widget{
			Name:        "main",
			ManagerType: reflect.TypeFor[MainManager](),
		},
		Help: Widget{
			Name:        "help",
			ManagerType: reflect.TypeFor[HelpManager](),
		},
		StatusBar: Widget{
			Name:        "status_bar",
			ManagerType: reflect.TypeFor[StatusBarManager](),
		},
		TableList: Widget{
			Name:        "table_list",
			ManagerType: reflect.TypeFor[TableListManager](),
		},
	},
}

func (app *Application) GetHelpText() string {
	return "\tCtrl-c quit\tF5 Fetch counts\t"
}
