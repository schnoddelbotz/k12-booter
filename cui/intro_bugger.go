package cui

import (
	"fmt"
	"log"
	"schnoddelbotz/k12-booter/diagnostics"
	"time"

	"github.com/jroimartin/gocui"
)

func (app *App) bugger() {
	time.Sleep(250 * time.Millisecond) // "wait" for main loop to be run once m(
	e, err := app.gui.View(ViewMain)
	if err != nil {
		log.Fatal(err)
	}

	app.gui.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(ViewCommand)
		return nil
	})

	e.Autoscroll = true

	app.gui.Update(func(g *gocui.Gui) error {
		e.Write([]byte("\n *** NOTA BENE ***\n"))
		e.Write([]byte(" Enterprise version does this twice as fast\n"))
		e.Write([]byte(" Press F2 for menu, F10 to quit\n\n"))
		fmt.Fprintf(e, "\033[37;1m%s\033[0m", `         ____________ 
        < k12-booter >
         ------------ 
                \   ^__^
                 \  (oo)\_______
                    (__)\       )\/\
                        ||----w |
                        ||     ||`+"\n\n")

		si := diagnostics.GetSysInfoData()
		fmt.Fprintf(e, "\n> k12-booter v%s (%s/%s) is ready\n", app.version, si.OS, si.Architecture)
		return nil
	})
}
