package tui

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/tui/config"
	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jesseduffield/gocui"
)

// ApplyKeybindingConfig applies keybindings from the provided configuration.
func (tui *Tui) ApplyKeybindingConfig(appConfig *config.AppConfig) {
	actionHandlers := tui.getActionHandlers()

	for _, kb := range appConfig.Keybindings {
		handler, ok := actionHandlers[kb.Action]
		if !ok {
			logger.Warnf("Unknown action '%s' in keybinding config", kb.Action)
			continue
		}

		gocuiKey, gocuiMod, err := config.ParseKeyString(kb.Key)
		if err != nil {
			logger.Warnf("Failed to parse key '%s' for action '%s': %v", kb.Key, kb.Action, err)
			continue
		}

		viewName := kb.View

		if err := tui.SetKeybinding(viewName, gocuiKey, gocuiMod, handler); err != nil {
			// Log error instead of panic for individual keybinding failures
			logger.Errorf("Failed to set keybinding for action '%s' (key: '%s', view: '%s'): %v", kb.Action, kb.Key, viewName, err)
		}
	}
}

// getActionHandlers returns a map of action strings to their corresponding handler functions.
func (tui *Tui) getActionHandlers() map[string]func(*gocui.Gui, *gocui.View) error {
	return map[string]func(*gocui.Gui, *gocui.View) error{
		"quit":         Quit,
		"cycle_view":   tui.cycleCurrentView(),
		"refresh_view": tui.onF5Pressed(),

		"sidepanel_select":   tui.sidePanel.OnSelect(tui.goolDb),
		"sidepanel_back":     tui.sidePanel.OnBack(tui.goolDb),
		"sidepanel_up":       tui.sidePanel.MoveCursorUp(),
		"sidepanel_down":     tui.sidePanel.MoveCursorDown(),
		"sidepanel_up_alt":   tui.sidePanel.MoveCursorUp(),   // For 'k'
		"sidepanel_down_alt": tui.sidePanel.MoveCursorDown(), // For 'j'

		"dataview_up":        tui.dataView.MoveRowUp(),
		"dataview_down":      tui.dataView.MoveRowDown(),
		"dataview_left":      tui.dataView.MoveColumnLeft(),
		"dataview_right":     tui.dataView.MoveColumnRight(),
		"dataview_up_alt":    tui.dataView.MoveRowUp(),       // For 'k'
		"dataview_down_alt":  tui.dataView.MoveRowDown(),     // For 'j'
		"dataview_left_alt":  tui.dataView.MoveColumnLeft(),  // For 'h'
		"dataview_right_alt": tui.dataView.MoveColumnRight(), // For 'l'

		"tablecell_open":        tui.onTableCellOpen(),
		"tablecell_close":       tui.onTableCellClose(),
		"tablecell_scroll_up":   tui.onTableCellScrollUp(),
		"tablecell_scroll_down": tui.onTableCellScrollDown(),

		"commandbar_open":  tui.onCommandBarOpen(),
		"commandbar_close": tui.onCommandBarClose(),
		"commandbar_enter": tui.onCommandBarEnter(),

		"queryoptions_open":  tui.onQueryOptionsOpen(),
		"queryoptions_close": tui.onQueryOptionsClose(),
		"queryoptions_enter": tui.onQueryOptionsEnter(),
	}
}

func (tui *Tui) onTableCellOpen() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, dataTableView *gocui.View) error {
		tui.Update(func(_ *gocui.Gui) error {
			data := tui.dataView.GetSelectedCellData()
			x0, y0, x1, y1 := dataTableView.Dimensions()

			popupView, err := tui.SetView(views.TableCellViewName, x0+5, y0+5, x1-5, y1-5, 0)
			if err != nil {
				if !gocui.IsUnknownView(err) {
					return err
				}
				popupView.Frame = true
				popupView.Wrap = true
				popupView.Editable = false
			}
			popupView.Clear()
			fmt.Fprint(popupView, data)

			tui.SetCurrentView(views.TableCellViewName)
			tui.SetViewOnTop(views.TableCellViewName)

			return nil
		})
		return nil
	}
}

func (tui *Tui) onTableCellClose() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		if err := tui.DeleteView(views.TableCellViewName); err != nil {
			if !gocui.IsUnknownView(err) {
				logger.Warnf("Error deleting TableCellView: %v", err)
			}
		}
		tui.SetCurrentView(tui.prevView)
		return nil
	}
}

func (tui *Tui) onTableCellScrollUp() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, tableCellView *gocui.View) error {
		if tableCellView != nil {
			tableCellView.ScrollUp(3)
		}
		return nil
	}
}

func (tui *Tui) onTableCellScrollDown() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, tableCellView *gocui.View) error {
		if tableCellView != nil {
			tableCellView.ScrollDown(3)
		}
		return nil
	}
}

func (tui *Tui) onCommandBarOpen() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		tui.commandBar.Open(tui.Gui)

		tui.SetCurrentView(views.CommandBarViewName)
		return nil
	}
}

func (tui *Tui) onCommandBarClose() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		tui.commandBar.Close(tui.Gui)

		// tui.prevView was set when CommandBar was activated.
		// This restores to the view active before command bar.
		tui.previousView()
		return nil
	}
}

func (tui *Tui) onCommandBarEnter() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		executed, err := tui.commandBar.KeyEnter()
		if executed {
			tui.commandBar.Close(tui.Gui)

			// tui.prevView was set when CommandBar was activated.
			// This restores to the view active before command bar.
			tui.previousView()
		}
		if err != nil {
			logger.Warnf("Command execution error: %v", err)
			return err
		}

		return nil
	}
}

func (tui *Tui) onQueryOptionsOpen() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		selection, mode := tui.sidePanel.CurrentSelection()
		if mode != "Tables" {
			return nil
		}
		newQuery := true
		dataTableSelection := tui.dataView.CurrentTable()
		if dataTableSelection != "" {
			selection = dataTableSelection
			newQuery = false
		}

		tui.queryOptions.CreatePopup(newQuery, selection)
		// tui.SetCurrentView(views.QueryOptionsViewName)
		tui.SetViewOnBottom(views.DataTableViewName)
		return nil
	}
}

func (tui *Tui) onQueryOptionsClose() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		if err := tui.DeleteView(views.QueryOptionsViewName); err != nil {
			if !gocui.IsUnknownView(err) {
				logger.Warnf("Error deleting QueryOptionsView: %v", err)
			}
		}
		tui.SetCurrentView(tui.prevView)
		return nil
	}
}

func (tui *Tui) onQueryOptionsEnter() func(_ *gocui.Gui, _ *gocui.View) error {
	return func(_ *gocui.Gui, _ *gocui.View) error {
		err := tui.queryOptions.KeyEnter()
		if err != nil {
			// somethigs check exec
		}
		tui.queryOptions.Close(tui.Gui)

		tui.previousView()
		return nil
	}
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
