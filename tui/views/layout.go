package views

import "github.com/jesseduffield/gocui"

const (
	edgeOffset          int = 1
	sidePanelWidthRatio int = 6
	commandBarHeight    int = 2
	helpBarHeight       int = 2
	columnPadding       int = 2
)

var DefaultTextColor = gocui.ColorDefault
var DefaultBackgroundColor = gocui.ColorDefault
var InvTextColor = gocui.ColorDefault | gocui.AttrReverse
var InvBackgroundColor = gocui.ColorDefault | gocui.AttrReverse

func SidePanel(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	termSizeX /= sidePanelWidthRatio
	termSizeY -= commandBarHeight + helpBarHeight
	x0, x1 := 0, termSizeX
	y0, y1 := 0, termSizeY
	return SidePanelViewName, x0, y0, x1, y1, 0
}

func DataView(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	termSizeY -= commandBarHeight + helpBarHeight
	x0, x1 := termSizeX/sidePanelWidthRatio, termSizeX
	y0, y1 := 0, termSizeY
	return DataTableViewName, x0, y0, x1, y1, 0
}

func CommandBar(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	termSizeY -= commandBarHeight
	x0, x1 := 0, termSizeX
	y0, y1 := termSizeY-commandBarHeight, termSizeY
	return CommandBarViewName, x0, y0, x1, y1, 0
}

func HelpBar(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	x0, x1 := 0, termSizeX
	y0, y1 := termSizeY-helpBarHeight, termSizeY
	return HelpFooterName, x0, y0, x1, y1, 0
}

func SetSidePanelInactiveColors(view *gocui.View) {
	view.FgColor = gocui.AttrDim
	view.BgColor = gocui.AttrDim
	view.SelFgColor = DefaultTextColor
	view.SelBgColor = DefaultBackgroundColor
}

func SetSidePanelColors(view *gocui.View) {
	view.FgColor = DefaultTextColor
	view.BgColor = DefaultBackgroundColor
	view.SelFgColor = gocui.AttrReverse
	view.SelBgColor = gocui.AttrReverse
}
