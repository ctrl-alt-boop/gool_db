package widgets

import (
	"github.com/jesseduffield/gocui"
)

const (
	edgeOffset          int = 1
	sidePanelWidthRatio int = 6
	commandBarHeight    int = 2
	helpBarHeight       int = 2
	columnPadding       int = 2
	modalMargin         int = 5
)

var DefaultForegroundColor = gocui.ColorDefault
var DefaultBackgroundColor = gocui.ColorDefault
var InvForegroundColor = gocui.ColorDefault | gocui.AttrReverse
var InvBackgroundColor = gocui.ColorDefault | gocui.AttrReverse

func SidePanelLayout(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	termSizeX /= sidePanelWidthRatio
	termSizeY -= commandBarHeight + helpBarHeight
	x0, x1 := 0, termSizeX
	y0, y1 := 0, termSizeY
	return SidePanelViewName, x0, y0, x1, y1, 0 //gocui.RIGHT
}

func DataViewLayout(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	termSizeY -= commandBarHeight + helpBarHeight
	x0, x1 := termSizeX/sidePanelWidthRatio, termSizeX
	y0, y1 := -1, termSizeY
	return DataAreaViewName, x0, y0, x1, y1, 0
}

func CommandBarLayout(termSizeX, termSizeY, extraHeight int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	termSizeY -= commandBarHeight
	x0, x1 := 0, termSizeX
	y0, y1 := termSizeY-commandBarHeight, termSizeY
	return CommandBarViewName, x0, y0 - extraHeight, x1, y1, 0
}

func HelpBarLayout(termSizeX, termSizeY int) (string, int, int, int, int, byte) {
	termSizeX -= edgeOffset
	x0, x1 := 0, termSizeX
	y0, y1 := termSizeY-helpBarHeight, termSizeY
	return HelpFooterName, x0, y0, x1, y1, 0
}

func Popup(name string, dimensionsFunc func() (int, int, int, int)) (string, int, int, int, int, byte) {
	x0, y0, x1, y1 := dimensionsFunc()
	return name, x0 + modalMargin, y0 + modalMargin, x1 - modalMargin, y1 - modalMargin, 0
}

// []rune{'─', '│', '┌', '┐', '└', '┘'}
func RoundedCorners() []rune {
	return []rune{'─', '│', '╭', '╮', '╰', '╯'}
}

// func Popup(name string, dataTableX0, dataTableY0, dataTableX1, dataTableY1 int) (string, int, int, int, int, byte) {
// 	x0, x1 := dataTableX0+modalMargin, dataTableX1-modalMargin
// 	y0, y1 := dataTableY0+modalMargin, dataTableY1-modalMargin
// 	return name, x0, y0, x1, y1, 0
// }

func SetSidePanelInactiveColors(view *gocui.View) {
	view.FgColor = gocui.AttrDim
	view.BgColor = gocui.AttrDim
	view.SelFgColor = DefaultForegroundColor
	view.SelBgColor = DefaultBackgroundColor
}

func SetSidePanelColors(view *gocui.View) {
	view.FgColor = DefaultForegroundColor
	view.BgColor = DefaultBackgroundColor
	view.SelFgColor = gocui.AttrReverse
	view.SelBgColor = gocui.AttrReverse
}
