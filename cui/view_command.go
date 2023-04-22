package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) commandLayoutFunc() (*gocui.View, error) {
	g := app.gui
	maxX, maxY := g.Size()
	if v, err := g.SetView(ViewCommand, 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Editable = true
		//v.Highlight = true
		v.Frame = false
		fmt.Fprintln(v, "")
		return v, nil
	}
	return app.views[ViewCommand].view, nil
}
