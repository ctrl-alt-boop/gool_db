package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type CommandLine struct {
	width, height int
	goolDb        *gooldb.GoolDb
}

func NewCommandBar(gool *gooldb.GoolDb) *CommandLine {
	return &CommandLine{
		goolDb: gool,
	}
}

func (c *CommandLine) Init() tea.Cmd {
	return nil
}

func (c *CommandLine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.UpdateSize(msg.Width, 1)
	}
	return c, nil
}

func (c *CommandLine) UpdateSize(termWidth, termHeight int) {
	c.width, c.height = termWidth-BorderThicknessDouble, termHeight
}

func (c *CommandLine) View() string {
	commandBorder := lipgloss.Border{
		Bottom:      "─",
		Top:         "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "├",
		TopRight:    "┤",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	input := huh.NewInput().
		Prompt(">")

	inputStyle := lipgloss.NewStyle().
		Width(c.width).
		Height(c.height).
		Align(lipgloss.Left, lipgloss.Center).
		Border(commandBorder, false, true, true, true)

	return inputStyle.Render(input.View())
}
