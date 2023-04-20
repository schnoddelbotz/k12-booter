package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) localeLayoutFunc() (*gocui.View, error) {
	g := app.gui
	maxX, _ := g.Size()
	if v, err := g.SetView(ViewLocale, maxX-33, -1, maxX-2, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Wrap = false
		v.Frame = false
		msg := fmt.Sprintf("[ %5s %-20s ]", app.localeInfo.Locale, app.localeInfo.LanguageLocalName)
		v.Write([]byte(msg))
		v.BgColor = gocui.ColorBlue
		v.FgColor = gocui.ColorYellow
		return v, nil
	}
	return app.views[ViewLocale].view, nil
}
