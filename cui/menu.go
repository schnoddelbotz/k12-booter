package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) menuLayoutFunc() (*gocui.View, error) {
	g := app.gui
	maxX, maxY := g.Size()
	if v, err := g.SetView(ViewMenu, 6, 6, maxX-6, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Title = "[ MENU ]"
		fmt.Fprintln(v, "SELECT one ... up/dn ... wip")
		fmt.Fprintln(v, "Press F2 again for now. Thanks.")
	}
	return app.views[ViewMenu].view, nil
}

func (app *App) handleMenuKeyUp(g *gocui.Gui, v *gocui.View) error {
	e := app.views[ViewMain].view
	fmt.Fprintln(e, "# menu debug: UP ↑")
	return nil
}

func (app *App) handleMenuKeyDown(g *gocui.Gui, v *gocui.View) error {
	e := app.views[ViewMain].view
	fmt.Fprintln(e, "# menu debug: DN ↓")
	return nil
}

func (app *App) setupMenuKeybindings() error {
	if err := app.gui.SetKeybinding(ViewMenu, gocui.KeyArrowUp, gocui.ModNone, app.handleMenuKeyUp); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(ViewMenu, gocui.KeyArrowDown, gocui.ModNone, app.handleMenuKeyDown); err != nil {
		return err
	}
	return nil
}
