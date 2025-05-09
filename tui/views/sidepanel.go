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

const SidePanelViewName string = "sidepanel"

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

	selectedDriver        string
	selectedDriverIndex   int
	selectedDatabase      string
	selectedDatabaseIndex int
}

func (s *SidePanelView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(SidePanel(maxX, maxY))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			logger.Panic(err)
		}

		view.Frame = true
		view.Editable = false
		view.Highlight = true
		view.SelFgColor = gocui.AttrReverse
		view.SelBgColor = gocui.AttrReverse

		g.SetCurrentView(SidePanelViewName)
		s.gui = g
		s.view = view
		s.SetSidePanelMode(DriverList)

		view.SetContent(strings.Join(gooldb.SupportedDrivers, "\n"))
	}

	return nil
}

func (s *SidePanelView) CurrentSelection() (selection string, mode string) {
	selection, _ = s.view.Word(s.view.Cursor())
	mode = s.currentMode.Name()
	return
}

func (s *SidePanelView) OnSelect(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		selection, ok := currentView.Word(currentView.Cursor())
		if !ok {
			return nil
		}
		switch s.currentMode {
		case DriverList:
			s.selectedDriverIndex = currentView.CursorY()
			gool.SelectDriver(selection)
			currentView.SetCursorY(0)

		case DatabaseList:
			s.selectedDatabaseIndex = currentView.CursorY()
			gool.SelectDatabase(selection)
			currentView.SetCursorY(0)

		case TableList:
			gool.SelectTable(selection)
		}
		return nil
	}
}

func (s *SidePanelView) OnBack(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		switch s.currentMode {
		case TableList:
			gool.SelectDriver(s.selectedDriver)
			currentView.SetCursorY(s.selectedDatabaseIndex)
		case DatabaseList:
			s.SetSidePanelMode(DriverList)
			s.view.SetContent(strings.Join(gooldb.SupportedDrivers, "\n"))
			currentView.SetCursorY(s.selectedDriverIndex)
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
	if err != nil {
		logger.Warn(err)
		return
	}
	s.gui.Update(func(g *gocui.Gui) error {
		s.SetSidePanelMode(DatabaseList)
		s.view.Title = args.Selected
		fmt.Fprint(s.view, strings.Join(args.Databases, "\n"))
		return nil
	})
	s.selectedDriver = args.Selected
	s.selectedDatabase = ""
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
	s.selectedDatabase = args.Selected
}

func (s *SidePanelView) SetSidePanelMode(mode sidePanelMode) {
	s.currentMode = mode
	s.view.Title = mode.Name()

	s.view.Clear()
}

func (s *SidePanelView) MoveCursorUp() func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, sidePanelView *gocui.View) error {
		if sidePanelView != nil {
			curY := sidePanelView.CursorY()
			if curY == 0 {
				return nil
			}
			sidePanelView.SetCursor(0, curY-1)
		}
		return nil
	}
}

func (s *SidePanelView) MoveCursorDown() func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, sidePanelView *gocui.View) error {
		if sidePanelView != nil {
			curY := sidePanelView.CursorY()
			numLines := len(sidePanelView.BufferLines())
			if curY >= numLines-1 {
				return nil
			}
			sidePanelView.SetCursor(0, curY+1)
		}
		return nil
	}
}

func (s *SidePanelView) SetInactiveColors() {
	s.gui.Update(func(g *gocui.Gui) error {
		s.view.TitleColor = DefaultForegroundColor | gocui.AttrDim
		s.view.FrameColor = DefaultForegroundColor | gocui.AttrDim

		s.view.FgColor = DefaultForegroundColor | gocui.AttrDim
		s.view.BgColor = DefaultBackgroundColor | gocui.AttrDim
		s.view.SelFgColor = DefaultForegroundColor | gocui.AttrBold
		s.view.SelBgColor = DefaultBackgroundColor | gocui.AttrBold
		return nil
	})
}

func (s *SidePanelView) SetActiveColors() {
	s.gui.Update(func(g *gocui.Gui) error {
		s.view.TitleColor = DefaultForegroundColor
		s.view.FrameColor = DefaultForegroundColor

		s.view.FgColor = DefaultForegroundColor
		s.view.BgColor = DefaultBackgroundColor
		s.view.SelFgColor = InvForegroundColor
		s.view.SelBgColor = InvBackgroundColor
		return nil
	})
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
