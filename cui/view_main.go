package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) printlnMain(s string) {
	e, err := app.gui.View(ViewMain)
	if err != nil {
		panic(err)
	}
	app.gui.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(e, s)
		return nil
	})
}

func (app *App) mainLayoutFunc() (*gocui.View, error) {
	g := app.gui
	maxX, maxY := g.Size()
	if v, err := g.SetView(ViewMain, 0, 0, maxX-app.helpWidth-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Wrap = true
		v.Title = " {[ k12-booter ]} "
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorGreen
		return v, nil
	}
	return nil, nil
}
