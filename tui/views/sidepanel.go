package views

import (
	"fmt"
	"log"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/internal/database"
	"github.com/jroimartin/gocui"
)

type SidePanelMode int

const SidePanelViewName string = "side_panel"

const (
	DriverList SidePanelMode = iota
	DatabaseList
	TableList
)

func (m SidePanelMode) Name() string {
	switch m {
	case DriverList:
		return "Drivers"
	case DatabaseList:
		return "Databases"
	case TableList:
		return "Tables"
	default:
		log.Panic("unknown SidePanelMode")
		return ""
	}
}

type SidePanelView struct {
	*gocui.View
	CurrentMode SidePanelMode
}

var SidePanel *SidePanelView

func InitSidePanel(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView(SidePanelViewName, 0, 0, maxX/6, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		SidePanel = &SidePanelView{
			View:        view,
			CurrentMode: 0,
		}

		view.Frame = true
		view.Editable = false
	}
	return nil
}

func (s *SidePanelView) OnEnter(gool *gooldb.GoolDb, selection string) error {
	switch SidePanel.CurrentMode {
	case DriverList:
		err := gool.SetDriver(selection)
		if err != nil {
			return err
		}
		SetSidePanelMode(DatabaseList)
	case DatabaseList:
		err := gool.SetDatabase(selection)
		if err != nil {
			return err
		}
		SetSidePanelMode(TableList)
	case TableList:

	}
	return nil
}

func SetSidePanelMode(mode SidePanelMode) error {
	SidePanel.CurrentMode = mode
	SidePanel.Title = mode.Name()
	SidePanel.Highlight = true
	SidePanel.SelFgColor = gocui.AttrReverse
	SidePanel.SelBgColor = gocui.AttrReverse

	SidePanel.Clear()
	fmt.Fprint(SidePanel.View, strings.Join(database.SupportedDrivers, "\n"))

	return nil
}
