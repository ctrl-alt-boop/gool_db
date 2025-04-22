package views

import "github.com/jroimartin/gocui"

const (
	Block     string = "█"
	SemiBlock string = "▒"
)

const CommandBarViewName string = "command_bar"

type CommandBarView struct {
	*gocui.Gui
}

var CommandBar *CommandBarView

func InitCommandBar(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if view, err := g.SetView(CommandBarViewName, 0, maxY-4, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		CommandBar = &CommandBarView{
			Gui: g,
		}

		view.Frame = true
		view.Editable = false
	}
	return nil
}

func (c *CommandBarView) OnEnter(text string) error {
	return nil
}
