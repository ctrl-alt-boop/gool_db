package widgets

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type CommandBar struct {
	width, height int
	goolDb        *gooldb.GoolDb
}

func CreateCommandBar(gool *gooldb.GoolDb) *CommandBar {
	return &CommandBar{
		goolDb: gool,
	}
}

func (c *CommandBar) Init() tea.Cmd {
	return nil
}

func (c *CommandBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.UpdateSize(msg.Width, 1)
	}
	return c, nil
}

func (c *CommandBar) UpdateSize(termWidth, termHeight int) {
	c.width, c.height = termWidth-BorderThicknessDouble, termHeight
}

func (c *CommandBar) View() string {
	commandBorder := lipgloss.Border{
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "├",
		TopRight:    "┤",
		BottomLeft:  "└",
		BottomRight: "┘",
		MiddleTop:   "┴",
	}

	input := huh.NewInput().
		Prompt(">")

	inputStyle := lipgloss.NewStyle().
		Width(c.width).
		Height(c.height).
		Align(lipgloss.Left, lipgloss.Bottom).
		Border(commandBorder, false, true, true, true)

	return inputStyle.Render(input.View())
}
