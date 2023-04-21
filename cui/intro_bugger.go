package cui

import (
	"fmt"
	"schnoddelbotz/k12-booter/internationalization"
	"time"

	"github.com/jroimartin/gocui"
)

func (app *App) bugger() {
	time.Sleep(250 * time.Millisecond) // "wait" for main loop to be run once m(

	e := app.views[ViewMain].view

	app.gui.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(e, " * Loading I18N data (fake) ...")
		return nil
	})

	time.Sleep(1 * time.Second)
	e.Autoscroll = true
	// WIP
	// Provide C# Cultures interface here.
	// Build POJOs, no structs holding data ... some typing work?

	for _, country := range internationalization.Cultures {
		app.gui.Update(func(g *gocui.Gui) error {
			//e.Write([]byte(flag + " "))
			fmt.Fprintf(e, "%s %s %s %s\n",
				country.Alpha2Code, country.Alpha3Code, country.InternetccTLD, country.CountryName)
			return nil
		})
		time.Sleep(7 * time.Millisecond)
	}

	time.Sleep(200 * time.Millisecond)

	app.gui.Update(func(g *gocui.Gui) error {
		e.Write([]byte("\n\n * NOTA BENE\n"))
		e.Write([]byte("Enterprise version does this twice as fast\n"))
		e.Write([]byte("Press F10 now to quit, then Ctrl-C\n\n"))

		e.Write([]byte("Should look for local config file now, ...\n"))
		e.Write([]byte("Or start a wizard to collect basic data\n"))

		fmt.Fprintf(e, "\033[37;1m%s\033[0m", `         ____________ 
        < k12-booter >
         ------------ 
                \   ^__^
                 \  (oo)\_______
                    (__)\       )\/\
                        ||----w |
                        ||     ||`+"\n\n")

		return nil
	})
}
