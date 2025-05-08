package views

import (
	"strings"
	"unicode"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/jesseduffield/gocui"
)

const (
	Block     string = "█"
	SemiBlock string = "▒"
)

const CommandBarViewName string = "command_bar"
const commandBarPrompt string = "> "

type CommandBarView struct {
	view   *gocui.View
	GoolDb *gooldb.GoolDb

	prevCommands []string

	prevChar    rune
	extraHeight int
}

func (c *CommandBarView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	view, err := g.SetView(CommandBar(maxX, maxY, c.extraHeight))
	if err != nil {
		if !gocui.IsUnknownView(err) {
			logger.Panic(err)
		}

		view.Frame = true
		view.Editable = true
		view.IgnoreCarriageReturns = true
		view.Editor = gocui.EditorFunc(c.commandBarEdit)

		c.view = view
		c.prevCommands = make([]string, 0)
		c.prevChar = ' '
		c.view.WriteString(commandBarPrompt)
		c.view.SetCursorX(2)
	}
	return nil
}

func (c *CommandBarView) Open(g *gocui.Gui) error {
	g.Cursor = true
	g.Update(func(_ *gocui.Gui) error {
		c.view.TextArea.TypeString(commandBarPrompt)
		return nil
	})

	return nil
}

func (c *CommandBarView) Close(g *gocui.Gui) error {
	g.Cursor = false
	c.extraHeight = 0
	g.Update(func(_ *gocui.Gui) error {
		if c.view != nil {
			c.view.Rewind()
			c.view.ClearTextArea()
			c.view.SetContent(commandBarPrompt)
			c.view.SetCursorX(2)
			c.prevChar = ' '
		}
		return nil
	})

	return nil
}

func (c *CommandBarView) KeyEnter() (executed bool, err error) {
	if c.prevChar == ';' {
		err = c.executeCommand()
		executed = true
		return
	}

	c.view.TextArea.TypeString("\n  ")
	c.extraHeight++
	c.view.RenderTextArea()
	executed = false
	return
}

func (c *CommandBarView) onBackSpace() {
	if c.view.Buffer() == commandBarPrompt {
		return
	}
	c.view.TextArea.BackSpaceChar()
}

func (c *CommandBarView) onLeftArrow() {
	if c.view.Buffer() == commandBarPrompt {
		return
	}
	c.view.TextArea.MoveCursorLeft()
}

func (c *CommandBarView) executeCommand() error {
	command := c.view.Buffer()[len(commandBarPrompt):]
	command = strings.ReplaceAll(command, "\n  ", " ")
	command = strings.TrimSpace(command)
	logger.Info(command)
	c.prevCommands = append(c.prevCommands, command)
	// err := c.GoolDb.ExecuteCommand(command)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (c *CommandBarView) commandBarEdit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		c.onBackSpace()
	case key == gocui.KeyCtrlD || key == gocui.KeyDelete:
		v.TextArea.DeleteChar()
	case key == gocui.KeyArrowLeft:
		c.onLeftArrow()
	case key == gocui.KeyArrowRight:
		v.TextArea.MoveCursorRight()
	case key == gocui.KeySpace:
		v.TextArea.TypeRune(' ')
	case key == gocui.KeyCtrlA || key == gocui.KeyHome:
		v.TextArea.GoToStartOfLine()
	case key == gocui.KeyCtrlE || key == gocui.KeyEnd:
		v.TextArea.GoToEndOfLine()
	case unicode.IsPrint(ch):
		c.prevChar = ch
		v.TextArea.TypeRune(ch)
	default:
		return false
	}
	v.RenderTextArea()

	return true
}
