package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) helpLayoutFunc() (*gocui.View, error) {
	g := app.gui
	maxX, maxY := g.Size()
	app.helpWidth = 0
	// toggle
	if app.showHelp {
		app.helpWidth = 31
		if v, err := g.SetView(ViewHelp, maxX-app.helpWidth-1, 1, maxX-1, maxY-3); err != nil {
			if err != gocui.ErrUnknownView {
				return nil, err
			}
			v.Title = "[ Context help ]"
			fmt.Fprintln(v, "This is k12-booter, yay!")
			app.helpView = v
		}
	} else {
		g.DeleteView(ViewHelp)
	}
	return app.views[ViewHelp].view, nil
}
