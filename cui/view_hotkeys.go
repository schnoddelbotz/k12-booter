package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) hotkeysLayoutFunc() (*gocui.View, error) {
	g := app.gui
	maxX, maxY := g.Size()
	width := maxX / 10
	for i, key := range app.hotkeysWidget.HotKeys {
		if v, err := g.SetView(key.ViewName, i*width-1, maxY-2, i*width+width-1, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return nil, err
			}
			v.Frame = false
			fmt.Fprintf(v, "\033[37;1m%d\033[0m\033[36;7m%-10s\033[0m\033[37;1m \033[0m", i+1, key.Label)
		}
	}
	return nil, nil
}
