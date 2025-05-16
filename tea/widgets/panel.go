package widgets

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/tea/data"
	"github.com/ctrl-alt-boop/gooldb/tea/event"
	"github.com/ctrl-alt-boop/gooldb/tea/ui"
)

type PanelMode string

const (
	DriverList   PanelMode = "driverList"
	DatabaseList PanelMode = "databaseList"
	TableList    PanelMode = "tableList"
)

type selection struct {
	index     int
	cachePath string
}

type Panel struct {
	list          *ui.List
	width, height int
	goolDb        *gooldb.GoolDb

	mode PanelMode

	isLoading bool
	spinner   spinner.Model

	selectHistory []selection
	cache         data.Cache
}

func CreateListPanel(gool *gooldb.GoolDb) *Panel {
	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = ui.MovingBlock

	return &Panel{
		list:          ui.NewList(),
		goolDb:        gool,
		mode:          DriverList,
		spinner:       loadingSpinner,
		selectHistory: make([]selection, 0),
		cache:         data.NewCache(),
	}
}

func (p *Panel) UpdateSize(termWidth, termHeight int) {
	p.width, p.height = termWidth/PanelWidthRatio-BorderThicknessDouble, termHeight-5
	p.list.SetSize(p.width, p.height)
}

func (p *Panel) SetMode(mode PanelMode) {
	p.mode = mode
}

func (p *Panel) GetMode() PanelMode {
	return p.mode
}

func (p *Panel) OnSelect() tea.Cmd {
	selection, ok := p.list.SelectedItem().(ui.ListItem)
	if !ok {
		return nil
	}
	switch p.mode {
	case DriverList:
		p.goolDb.SelectDriver(string(selection))
	case DatabaseList:
		p.goolDb.SelectDatabase(string(selection))
	case TableList:
		p.goolDb.SelectTable(string(selection))
	}
	p.isLoading = true

	return p.spinner.Tick
}

func (p *Panel) Init() tea.Cmd {
	driverList := p.goolDb.GetDrivers()
	p.list.SetItems(driverList)
	return nil
}

func (p *Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.UpdateSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Up):
			p.list.CursorUp()
		case key.Matches(msg, KeyMap.Down):
			p.list.CursorDown()
		case key.Matches(msg, KeyMap.Select):
			return p, p.OnSelect()
		}
	case event.GoolDbEventMsg:
		p.isLoading = false
		if msg.Err != nil {
			logger.Error(msg.Err)
			return p, nil
		}
		switch msg.Type {
		case gooldb.DriverSet:
			args, ok := msg.Args.(gooldb.DriverSetEvent)
			if ok {
				p.cache.Add(data.NewCacheable(args.Selected, args.Databases))
				p.list.SetItems(args.Databases)
				p.SetMode(DatabaseList)
			}
		case gooldb.DatabaseSet:
			args, ok := msg.Args.(gooldb.DatabaseSetEvent)
			if ok {
				p.cache.Add(data.NewCacheable(args.Selected, args.Tables))
				p.list.SetItems(args.Tables)
				p.SetMode(TableList)
			}
		case gooldb.TableSet:
			args, ok := msg.Args.(gooldb.TableSetEvent)
			if ok {
				p.cache.Forward(args.Selected)
			}
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd
	}

	return p, nil
}

func (p *Panel) View() string {
	panelBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┬",
		BottomLeft:  "├",
		BottomRight: "┴",
	}

	panelStyle := lipgloss.NewStyle().
		Height(p.height).
		Width(p.width).
		Border(panelBorder, true, true, true, true).
		Align(lipgloss.Left, lipgloss.Top)

	if p.isLoading {
		return panelStyle.
			Height(p.height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Top).
			Render(p.spinner.View())
	}

	return panelStyle.Render(p.list.View())
}
