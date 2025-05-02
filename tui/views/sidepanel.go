package views

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/jesseduffield/gocui"
)

var logger = logging.NewLogger("sidepanel.log")

type sidePanelMode int

const SidePanelViewName string = "side_panel"

const (
	DriverList sidePanelMode = iota
	DatabaseList
	TableList
)

func (m sidePanelMode) Name() string {
	switch m {
	case DriverList:
		return "Drivers"
	case DatabaseList:
		return "Databases"
	case TableList:
		return "Tables"
	default:
		panic("unknown SidePanelMode")
	}
}

type SidePanelView struct {
	view        *gocui.View
	gui         *gocui.Gui
	currentMode sidePanelMode
}

func (s *SidePanelView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(SidePanel(maxX, maxY))
	if err != nil {
		if err != gocui.ErrUnknownView {
			//panic(err)
		}

		view.Frame = true
		view.Editable = false
		view.Highlight = true

		g.SetCurrentView(SidePanelViewName)
		s.gui = g
		s.view = view
		s.SetSidePanelMode(DriverList)
	}

	return nil
}

func (s *SidePanelView) OnEnterPressed(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		selection, ok := currentView.Word(currentView.Cursor())
		if !ok {
			return nil
		}
		switch s.currentMode {
		case DriverList:
			gool.SetDriver(selection)
		case DatabaseList:
			gool.SetDatabase(selection)
		case TableList:
			gool.SetTable(selection)
		}
		return nil
	}
}

// selected string, databases []string
func (s *SidePanelView) OnDriverSet(eventArgs any, err error) {
	logger.Info("OnDriverSet: ", eventArgs, err)
	args, ok := eventArgs.(gooldb.DriverSetEvent)
	if !ok {
		logger.Warn(eventArgs, args, ok)
		return
	}
	logger.Info("OnDriverSet: ")
	if err != nil {
		logger.Warn(err)
		return
	}
	logger.Info("OnDriverSet: ")
	s.gui.Update(func(g *gocui.Gui) error {
		logger.Info("OnDriverSet: ", "SetSidePanelMode(DatabaseList)")
		s.SetSidePanelMode(DatabaseList)
		s.view.Title = args.Selected
		fmt.Fprint(s.view, strings.Join(args.Databases, "\n"))
		return nil
	})
}

// selected string, tables []string
func (s *SidePanelView) OnDatabaseSet(eventArgs any, err error) {
	logger.Info("OnDatabaseSet: ", eventArgs, err)
	args, ok := eventArgs.(gooldb.DatabaseSetEvent)
	if !ok {
		logger.Warn(err)
		return
	}
	if err != nil {
		logger.Warn(err)
		return
	}
	s.gui.Update(func(g *gocui.Gui) error {
		s.SetSidePanelMode(TableList)
		s.view.Title = args.Selected
		fmt.Fprint(s.view, strings.Join(args.Tables, "\n"))
		return nil
	})
}

func (s *SidePanelView) SetSidePanelMode(mode sidePanelMode) {
	s.currentMode = mode
	s.view.Title = mode.Name()
	s.view.Highlight = true
	s.view.SelFgColor = gocui.AttrReverse
	s.view.SelBgColor = gocui.AttrReverse
	s.view.SetCursor(0, 0)

	s.view.Clear()

	if mode == DriverList {
		fmt.Fprint(s.view, strings.Join(gooldb.SupportedDrivers, "\n"))
	}
}

// maxX, _ := s.gui.Size()
// _, tableSizeY := dataView.Size()
// currentX := maxX / 6
// s.gui.Update(func(g *gocui.Gui) error {
// 	for i, col := range table.Columns() {
// 		data, width := table.GetColumnRows(i)
// 		if width <= len(col.Name) {
// 			width = columnPadding + len(col.Name) + columnPadding
// 		}
// 		colView, err := s.gui.SetView(DataColumnViewPrefix+col.Name, currentX, 0, currentX+width+columnPadding, tableSizeY+1)
// 		if err != nil {
// 			if err != gocui.ErrUnknownView {
// 				s.OnError(err)
// 				return err
// 			}
// 			colView.Title = col.Name
// 			colView.Frame = true
// 			colView.Editable = false
// 			fmt.Fprint(colView, strings.Join(data, "\n"))
// 		}
// 		currentX += width + columnPadding
// 		AddCurrentDataColumns(colView)
// 	}
// 	return nil
// })
