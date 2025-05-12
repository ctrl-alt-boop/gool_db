package tui

import (
	"github.com/ctrl-alt-boop/gooldb/tui/widgets"
	"github.com/jesseduffield/gocui"
)

func (tui *Tui) cycleCurrentView() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(g *gocui.Gui, currentView *gocui.View) error {
		var nextViewName string

		switch currentView.Name() {
		case widgets.SidePanelViewName:
			if !tui.dataView.IsTableSet() {
				return nil
			}
			nextViewName = widgets.DataAreaViewName
			tui.sidePanel.SetInactiveColors()

		case widgets.DataAreaViewName:
			nextViewName = widgets.SidePanelViewName
			tui.dataView.ClearHighlight()
			tui.sidePanel.SetActiveColors()
		}

		newActiveView := tui.SetCurrentView(nextViewName)

		if newActiveView != nil {
			g.Cursor = false
			if nextViewName == widgets.DataAreaViewName {
				tui.dataView.HighlightSelectedCell()
			}
		}

		return nil
	}
}

func (tui *Tui) previousView() {
	if tui.prevView != "" && tui.prevView != widgets.CommandBarViewName {
		tui.SetCurrentView(tui.prevView)
	} else {
		tui.SetCurrentView(widgets.SidePanelViewName)
	}
}

// SetCurrentView updates the current view and stores the previous one.
func (tui *Tui) SetCurrentView(name string) *gocui.View {
	currentViewName := ""
	if cv := tui.CurrentView(); cv != nil {
		currentViewName = cv.Name()
	}

	view, err := tui.Gui.SetCurrentView(name)
	if err != nil {
		logger.Fatalf("Failed to set current view to '%s': %v", name, err)
		return view
	}

	// Only update prevView if the view switch was successful and different
	if currentViewName != name {
		tui.prevView = currentViewName
	}
	if name == widgets.SidePanelViewName {
		tui.helpFooter.SetCurrentView(tui.sidePanel.GetModeString())
	} else {
		tui.helpFooter.SetCurrentView(name)
	}
	return view
}
