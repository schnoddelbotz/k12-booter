package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) promptLayoutFunc() (*gocui.View, error) {
	g := app.gui
	_, maxY := g.Size()
	if v, err := g.SetView(ViewPrompt, -1, maxY-3, 1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Wrap = false
		v.Frame = false
		fmt.Fprintf(v, "]") // conforming to LSL1 standard
		return v, nil
	}
	return nil, nil
}
