package widgets

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/jesseduffield/gocui"
)

var logger = logging.NewLogger("sidepanel.log")

type SidePanelMode int

const SidePanelViewName string = "sidepanel"

const (
	DriverList SidePanelMode = iota
	DatabaseList
	TableList
)

func (m SidePanelMode) String() string {
	switch m {
	case DriverList:
		return "drivers"
	case DatabaseList:
		return "databases"
	case TableList:
		return "tables"
	default:
		logger.Panic("unknown SidePanelMode")
		return ""
	}
}

func (m SidePanelMode) Name() string {
	switch m {
	case DriverList:
		return "Drivers"
	case DatabaseList:
		return "Databases"
	case TableList:
		return "Tables"
	default:
		logger.Panic("unknown SidePanelMode")
		return ""
	}
}

type sidePanelState struct {
	mode                  SidePanelMode
	selectedDriver        string
	selectedDriverIndex   int
	selectedDatabase      string
	selectedDatabaseIndex int
	selectedTable         string
	selectedTableIndex    int

	drivers    []string
	databases  []string
	tables     []string
	showCounts bool
	counts     map[string]int
}

type SidePanel struct {
	view *gocui.View
	gui  *gocui.Gui

	state sidePanelState
}

func CreateSidePanel() *SidePanel {
	return &SidePanel{
		state: sidePanelState{
			mode:    DriverList,
			drivers: gooldb.SupportedDrivers,
		},
	}
}

func (s *SidePanel) SetCounts(counts map[string]int) {
	s.state.counts = counts
}

func (s *SidePanel) SetTableCounts(tablesWithCounts map[string]int) {
	counts := make([]string, 0)
	for table, count := range tablesWithCounts {
		counts = append(counts, fmt.Sprintf("%s (%d rows)", table, count))
	}
	s.gui.Update(func(_ *gocui.Gui) error {
		s.view.SetContent(strings.Join(counts, "\n"))
		return nil
	})
}

func (s *SidePanel) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(SidePanelLayout(maxX, maxY))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			logger.Panic(err)
		}

		view.Frame = true
		view.Editable = false
		view.Highlight = true
		view.SelFgColor = gocui.AttrReverse
		view.SelBgColor = gocui.AttrReverse

		s.gui = g
		s.view = view
		s.SetMode(s.state.mode)
		g.SetCurrentView(SidePanelViewName)
	}

	return nil
}

func (s *SidePanel) ShowNoCounts() {
	s.gui.Update(func(g *gocui.Gui) error {
		s.view.Title = s.state.selectedTable
		s.view.SetContent(strings.Join(s.state.tables, "\n"))
		return nil
	})
}

func (s *SidePanel) ShowCounts() {
	tablesWithCounts := make([]string, len(s.state.counts))
	for i, table := range s.state.tables {
		tablesWithCounts[i] = fmt.Sprintf("%s (%d rows)", table, s.state.counts[table])
	}
	s.gui.Update(func(g *gocui.Gui) error {
		s.view.Title = s.state.selectedTable
		s.view.SetContent(strings.Join(tablesWithCounts, "\n"))
		return nil
	})
}

func (s *SidePanel) CurrentSelection() (selection string, mode SidePanelMode) {
	selection, _ = s.view.Word(s.view.Cursor())
	mode = s.state.mode
	return
}

func (s *SidePanel) OnSelect(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		selection, ok := currentView.Word(currentView.Cursor())
		if !ok {
			return nil
		}
		switch s.state.mode {
		case DriverList:
			s.state.selectedDriverIndex = currentView.CursorY()
			gool.SelectDriver(selection)
			currentView.SetCursorY(0)

		case DatabaseList:
			s.state.selectedDatabaseIndex = currentView.CursorY()
			gool.SelectDatabase(selection)
			currentView.SetCursorY(0)

		case TableList:
			s.state.selectedTableIndex = currentView.CursorY()
			s.state.selectedTable = selection
			gool.SelectTable(selection)

		}
		return nil
	}
}

func (s *SidePanel) OnBack(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		switch s.state.mode {
		case TableList:
			s.SetMode(DatabaseList)
			currentView.SetCursorY(s.state.selectedDatabaseIndex)
		case DatabaseList:
			s.SetMode(DriverList)
			currentView.SetCursorY(s.state.selectedDriverIndex)
		}
		return nil
	}
}

func (s *SidePanel) ToggleCounts(fetch func(tables []string) map[string]int) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		if s.state.counts == nil {
			s.state.counts = fetch(s.state.tables)
		}
		s.state.showCounts = !s.state.showCounts
		if s.state.showCounts {
			s.ShowCounts()
		} else {
			s.ShowNoCounts()
		}
		return nil
	}
}

// selected string, databases []string
func (s *SidePanel) OnDriverSet(eventArgs any, err error) {
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
	s.state.databases = args.Databases
	s.state.selectedDriver = args.Selected
	s.state.selectedDatabase = ""
	s.gui.Update(func(g *gocui.Gui) error {
		s.SetMode(DatabaseList)
		s.view.Title = args.Selected
		return nil
	})
}

// selected string, tables []string
func (s *SidePanel) OnDatabaseSet(eventArgs any, err error) {
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
	s.state.tables = args.Tables
	s.SetMode(TableList)
	s.ShowNoCounts()
	s.state.selectedDatabase = args.Selected
}

func (s *SidePanel) GetModeString() string {
	return s.state.mode.String()
}

func (s *SidePanel) SetMode(mode SidePanelMode) {
	s.state.mode = mode
	content := ""
	switch mode {
	case DriverList:
		Help.SetCurrentView(DriverList.String())
		content = strings.Join(s.state.drivers, "\n")
	case DatabaseList:
		Help.SetCurrentView(DatabaseList.String())
		content = strings.Join(s.state.databases, "\n")
	case TableList:
		Help.SetCurrentView(TableList.String())
		content = strings.Join(s.state.tables, "\n")
	}

	s.gui.Update(func(g *gocui.Gui) error {
		s.view.Clear()
		s.view.SetContent(content)
		s.view.Title = mode.Name()
		return nil
	})
}

func (s *SidePanel) MoveCursorUp() func(*gocui.Gui, *gocui.View) error {
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

func (s *SidePanel) MoveCursorDown() func(*gocui.Gui, *gocui.View) error {
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

func (s *SidePanel) SetInactiveColors() {
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

func (s *SidePanel) SetActiveColors() {
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
