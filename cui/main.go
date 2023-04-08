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
}

func Zain() {
	var (
		app App
		err error
	)
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
	app.gui.FgColor = gocui.ColorRed

	app.gui.SetManagerFunc(app.layout)

	if err := app.initKeybindings(app.gui); err != nil {
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
	for country := range internationalization.CountryFlags {
		app.gui.Update(func(g *gocui.Gui) error {
			//e.Write([]byte(flag + " "))
			fmt.Fprintf(e, "%s ", country)
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

		fmt.Fprintf(e, "%s", ` ____________ 
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

	if v, err := g.SetView(ViewTopMenu, 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[ k12-booter ]"
		fmt.Fprintln(v, " * FIRSTBOOT ")
		v.FgColor = gocui.ColorGreen
		app.topmenuView = v
	}

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

	if v, err := g.SetView(ViewMain, 0, 2, maxX-helpWidth-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		app.mainView = v
	}

	if v, err := g.SetView(ViewCommand, 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(ViewCommand); err != nil {
			return err
		}
		v.Title = "[ Command ]"
		v.Editable = true
		v.Highlight = true
		fmt.Fprintln(v, "")
		app.commandView = v
	}

	if v, err := g.SetView(ViewShortcuts, 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//v.Title = "Shortcuts"
		v.BgColor = gocui.ColorBlue
		v.Frame = false
		// https://github.com/gczgcz2015/gocui/blob/master/_examples/colors256.go
		// NC was: black empty bg instead of | char
		fmt.Fprintf(v, "F1 Help | F2 CMD | F3 View | F4 Edit | F5 Search | F9 Mask | F10 Exit")
		app.shortcutsView = v
	}

	return nil
}

func (app *App) initKeybindings(g *gocui.Gui) error {
	// they must depend on current view ...?
	// ESC should ...? Get into main view and let it consume further ESCs.
	// F2 in mainview should select CMD as active win

	if err := g.SetKeybinding("", gocui.KeyF10, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(ViewCommand, gocui.KeyF9, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.Mask ^= '*'
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			app.showHelp = !app.showHelp // toggle
			return nil
		}); err != nil {
		return err
	}
	return nil
}
