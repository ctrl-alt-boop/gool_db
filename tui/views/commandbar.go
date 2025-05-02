package views

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/jesseduffield/gocui"
)

const (
	Block     string = "█"
	SemiBlock string = "▒"
)

const CommandBarViewName string = "command_bar"

type CommandBarView struct {
	view   *gocui.View
	GoolDb *gooldb.GoolDb
}

func (c *CommandBarView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(CommandBar(maxX, maxY))
	if err != nil {
		if err != gocui.ErrUnknownView {
			//panic(err)
		}

		view.Frame = true
		view.Editable = false
		c.view = view
		fmt.Fprint(view, "...")
	}
	return nil
}

func (c *CommandBarView) OnEnterPressed(gool *gooldb.GoolDb) func(*gocui.Gui, *gocui.View) error {
	return func(_ *gocui.Gui, currentView *gocui.View) error {
		c.OnEnter(currentView.Buffer())
		return nil
	}
}

func (c *CommandBarView) OnEnter(text string) error {
	return nil
}
