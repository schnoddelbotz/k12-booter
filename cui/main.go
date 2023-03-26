package cui

import (
	"fmt"
	"log"
	"schnoddelbotz/k12-booter/internationalization"
	"time"

	"github.com/jroimartin/gocui"
)

const (
	ViewHelp   = "help"
	ViewEditor = "input"
)

type App struct {
	gui        *gocui.Gui
	helpView   *gocui.View
	editorView *gocui.View
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

	fmt.Println("I hope you could be convinced. OSS only has advantages.")
}

func (app *App) bugger() {
	time.Sleep(250 * time.Millisecond) // "wait" for main loop to be run once m(

	e := app.editorView
	e.Write([]byte(" * Load I18N data\n"))

	time.Sleep(1 * time.Second)

	for _, flag := range internationalization.CountryFlags {
		app.gui.Update(func(g *gocui.Gui) error {
			e.Write([]byte(flag))
			return nil
		})
		time.Sleep(40 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)

	app.gui.Update(func(g *gocui.Gui) error {
		e.Write([]byte("\n * NOTA BENE\n"))
		e.Write([]byte("Enterprise version does this twice as fast\n"))
		e.Write([]byte("Press Ctrl-A now, then Ctrl-C\n"))
		return nil
	})
}

func (app *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(ViewHelp, maxX-23, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keybindings"
		fmt.Fprintln(v, "^a: Set mask")
		fmt.Fprintln(v, "^c: Exit")
		app.helpView = v
	}

	if v, err := g.SetView(ViewEditor, 0, 0, maxX-24, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(ViewEditor); err != nil {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Title = "Enter your wishes"
		app.editorView = v
	}

	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyCtrlA, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.Mask ^= '*'
			return nil
		}); err != nil {
		return err
	}
	return nil
}
