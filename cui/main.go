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

	app.gui.SetManagerFunc(app.layout)

	if err := initKeybindings(app.gui); err != nil {
		log.Fatalln(err)
	}

	go app.bugger()

	if err := app.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func (app *App) bugger() {
	time.Sleep(250 * time.Millisecond) // "wait" for main loop to be run once m(

	e := app.mainView

	app.gui.Update(func(g *gocui.Gui) error {
		e.Write([]byte(" * Load I18N data\n"))
		return nil
	})

	time.Sleep(1 * time.Second)
	e.Autoscroll = true
	for _, flag := range internationalization.CountryFlags {
		app.gui.Update(func(g *gocui.Gui) error {
			e.Write([]byte(" " + flag + " "))
			return nil
		})
		time.Sleep(40 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)

	app.gui.Update(func(g *gocui.Gui) error {
		e.Write([]byte("\n\n * NOTA BENE\n"))
		e.Write([]byte("Enterprise version does this twice as fast\n"))
		e.Write([]byte("Press F10 now to quit, then Ctrl-C\n\n"))

		e.Write([]byte("Should look for local config file now, ...\n"))
		e.Write([]byte("Or start a wizard to collect basic data\n"))
		return nil
	})
}

func (app *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(ViewTopMenu, 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " * k12-booter * "
		fmt.Fprintln(v, " FIRSTBOOT ")
		v.FgColor = gocui.ColorGreen
		app.topmenuView = v
	}

	if v, err := g.SetView(ViewMain, 0, 2, maxX-21, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		app.mainView = v
	}

	// toggle
	if v, err := g.SetView(ViewHelp, maxX-20, 2, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[ Context help ]"
		fmt.Fprintln(v, "This is k12-booter, yay!")
		app.helpView = v
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
		fmt.Fprintln(v, "F1 Help | F2 Run | F3 View | F4 Edit | F5 Copy | F6-F8 For SALE | F9 Mask | F10 Exit")
		app.shortcutsView = v
	}

	return nil
}

func initKeybindings(g *gocui.Gui) error {
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
	return nil
}
