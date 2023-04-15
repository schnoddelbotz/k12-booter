package cui

import (
	"fmt"
	"log"
	"schnoddelbotz/k12-booter/internationalization"
	"time"

	"github.com/jroimartin/gocui"
)

const (
	ViewTopMenu   = "topmenu"
	ViewMain      = "main"
	ViewShortcuts = "shortcuts"
	ViewCommand   = "command"
	ViewHelp      = "help"
)

type App struct {
	gui           *gocui.Gui
	topmenuView   *gocui.View // AS 400 like first row menu info
	mainView      *gocui.View // the (?) main view
	commandView   *gocui.View // interactive commands like: go infra ...
	shortcutsView *gocui.View // Norton Commander / AS400-like bottom keymap
	helpView      *gocui.View // sidebar; should toggle with F1
	showHelp      bool
	registry      *Registry
}

func Zain() {
	var (
		app App
		err error
	)
	app.registry.Init()
	app.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.gui.Close()
	app.gui.Cursor = true
	app.gui.Mouse = true
	app.gui.ASCII = true
	// next two lines will enable color upon g.SetCurrentView(nextview)
	app.gui.Highlight = true

	app.gui.SetManagerFunc(app.layout)

	if err := app.SetHotkeyKeybindings(); err != nil {
		log.Fatalln(err)
	}

	go app.bugger()

	if err := app.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func (app *App) bugger() {
	time.Sleep(25 * time.Millisecond) // "wait" for main loop to be run once m(

	e := app.mainView

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
			fmt.Fprintf(e, "%4s %s %s %s %s\n",
				country.Flag, country.Alpha2Code, country.Alpha3Code, country.InternetccTLD, country.CountryName)
			return nil
		})
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)

	app.gui.Update(func(g *gocui.Gui) error {
		e.Write([]byte("\n\n * NOTA BENE\n"))
		e.Write([]byte("Enterprise version does this twice as fast\n"))
		e.Write([]byte("Press F10 now to quit, then Ctrl-C\n\n"))

		e.Write([]byte("Should look for local config file now, ...\n"))
		e.Write([]byte("Or start a wizard to collect basic data\n"))

		fmt.Fprintf(e, "\033[37;1m%s\033[0m", `        ____________ 
        < k12-booter >
         ------------ 
                \   ^__^
                 \  (oo)\_______
                    (__)\       )\/\
                        ||----w |
                        ||     ||`)

		return nil
	})
}

func (app *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	helpWidth := 0

	// toggle
	if app.showHelp {
		helpWidth = 30
		if v, err := g.SetView(ViewHelp, maxX-helpWidth-1, 2, maxX-1, maxY-4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "[ Context help ]"
			fmt.Fprintln(v, "This is k12-booter, yay!")
			app.helpView = v
		}
	} else {
		g.DeleteView(ViewHelp)
	}

	if v, err := g.SetView(ViewMain, 0, 0, maxX-helpWidth-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = " [ k12-booter - FirstBoot ] "
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorGreen
		app.mainView = v
	}

	if v, err := g.SetView(ViewCommand, 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(ViewCommand); err != nil {
			return err
		}
		v.Editable = true
		//v.Highlight = true

		v.Frame = false
		fmt.Fprintln(v, "")
		app.commandView = v
	}

	return LayoutHotkeys(g)
}
