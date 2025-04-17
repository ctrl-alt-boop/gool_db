package tui

import "github.com/jroimartin/gocui"

func MoveCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, -1, false)
	}
	return nil
}

func MoveCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, 1, false)
	}
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
