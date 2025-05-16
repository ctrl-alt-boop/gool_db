package widgets

import "github.com/jesseduffield/gocui"

type Drawer struct {
	gui   *gocui.Gui
	draws []func()
}

func (d *Drawer) Layout(g *gocui.Gui) error {
	for _, draw := range d.draws {
		draw()
	}
	return nil
}

func CreateDrawer(g *gocui.Gui) *Drawer {
	return &Drawer{
		gui:   g,
		draws: make([]func(), 0),
	}
}

func (d *Drawer) SetLeftBorder(dimensions func() (int, int, int, int), ch ...rune) {
	x0, y0, _, y1 := dimensions()

	d.draws = append(d.draws, func() {
		d.DrawVerticalLine(x0, y0, y1, ch...)
	})
}

func (d *Drawer) SetRightBorder(dimensions func() (int, int, int, int), ch ...rune) {
	_, y0, x1, y1 := dimensions()

	d.draws = append(d.draws, func() {
		d.DrawVerticalLine(x1, y0, y1, ch...)
	})
}

func (d *Drawer) SetTopBorder(dimensions func() (int, int, int, int), ch ...rune) {
	x0, y0, x1, _ := dimensions()

	d.draws = append(d.draws, func() {
		d.DrawHorizontalLine(x0, x1, y0, ch...)
	})
}

func (d *Drawer) SetTitledTopBorder(dimensions func() (int, int, int, int), title func() string, ch ...rune) {
	x0, y0, x1, _ := dimensions()

	d.draws = append(d.draws, func() {
		d.DrawTitledHorizontalLine(x0, x1, y0, title(), ch...)
	})
}

func (d *Drawer) SetBottomBorder(dimensions func() (int, int, int, int), ch ...rune) {
	x0, _, x1, y1 := dimensions()
	d.draws = append(d.draws, func() {
		d.DrawHorizontalLine(x0, x1, y1, ch...)
	})
}

// first rune: line, second: left corner, third: right corner
// line left + x0+1 -> x1-1 + right
// default line rune: ─ (\u2500)
// if only right wanted second arg = ”
func (d *Drawer) DrawHorizontalLine(x0, x1, y int, ch ...rune) {
	width, height := d.gui.Size()
	if y < 0 || y >= height {
		return
	}
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if x0 < 0 {
		x0 = 0
	}
	if x1 > width {
		x1 = width
	}
	if len(ch) == 0 {
		ch = append(ch, '─')
	}
	line := ch[0]
	var left, right rune
	if len(ch) > 1 {
		left = ch[1]
	}
	if len(ch) > 2 {
		right = ch[2]
	}
	for x := x0; x <= x1; x++ {
		if x == x0 {
			d.gui.SetRune(x, y, left, DefaultForegroundColor, DefaultBackgroundColor)
		} else if x == x1 {
			d.gui.SetRune(x, y, right, DefaultForegroundColor, DefaultBackgroundColor)
		} else {
			d.gui.SetRune(x, y, line, DefaultForegroundColor, DefaultBackgroundColor)
		}
	}
}

// first rune: line, second: top corner, third: bottom corner
// line left + y0+1 -> y1-1 + right
// default line rune: │ (\u2502)
// if only bottom wanted second arg = ”
func (d *Drawer) DrawVerticalLine(x, y0, y1 int, ch ...rune) {
	width, height := d.gui.Size()
	if x < 0 || x >= width {
		return
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if y0 < 0 {
		y0 = 0
	}
	if y1 > height {
		y1 = height
	}
	if len(ch) == 0 {
		ch = append(ch, '│')
	}
	line := ch[0]
	var top, bottom rune
	if len(ch) > 1 {
		top = ch[1]
	}
	if len(ch) > 2 {
		bottom = ch[2]
	}
	for y := y0; y <= y1; y++ {
		if y == y0 {
			d.gui.SetRune(x, y, top, 0, 0)
		} else if y == y1 {
			d.gui.SetRune(x, y, bottom, 0, 0)
		} else {
			d.gui.SetRune(x, y, line, 0, 0)
		}
	}
}

func (d *Drawer) DrawTitledHorizontalLine(x0, x1, y int, title string, ch ...rune) {
	width, height := d.gui.Size()
	if y < 0 || y >= height {
		return
	}
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if x0 < 0 {
		x0 = 0
	}
	if x1 > width {
		x1 = width
	}
	if len(ch) == 0 {
		ch = append(ch, '─')
	}
	line := ch[0]
	var left, right rune
	if len(ch) > 1 {
		left = ch[1]
	}
	if len(ch) > 2 {
		right = ch[2]
	}

	x := x0
	d.gui.SetRune(x, y, left, DefaultForegroundColor, DefaultBackgroundColor)
	x++
	d.gui.SetRune(x, y, line, DefaultForegroundColor, DefaultBackgroundColor)
	x++
	for _, r := range title {
		d.gui.SetRune(x, y, r, DefaultForegroundColor, DefaultBackgroundColor)
		x++
	}
	for x <= x1-1 {
		d.gui.SetRune(x, y, line, DefaultForegroundColor, DefaultBackgroundColor)
		x++
	}
	d.gui.SetRune(x, y, right, DefaultForegroundColor, DefaultBackgroundColor)
}
