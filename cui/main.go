package cui

import (
	"log"
	"schnoddelbotz/k12-booter/internationalization"

	"github.com/jroimartin/gocui"
)

const (
	ViewTopMenu   = "topmenu"
	ViewMain      = "main"
	ViewShortcuts = "shortcuts"
	ViewCommand   = "command"
	ViewLocale    = "locale"
	ViewHelp      = "help"
	ViewMenu      = "menu"
)

type App struct {
	gui           *gocui.Gui
	views         map[string]*AppView
	hotkeysWidget *Widget
	helpView      *gocui.View // sidebar; should toggle with F1
	showHelp      bool
	helpWidth     int
	localeInfo    internationalization.LocaleInfo
}

type ViewIdentifier string

type AppView struct {
	name          ViewIdentifier
	view          *gocui.View
	hotkeys       []HotKey
	isCurrentView bool
	layoutFunc    func() (*gocui.View, error)
}

func Zain() {
	var (
		app App
		err error
	)
	app.localeInfo = internationalization.GetLocaleInfo()
	app.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.gui.Close()
	app.gui.Cursor = true
	app.gui.Mouse = true
	app.gui.ASCII = true
	// next two (?) lines will enable color upon g.SetCurrentView(nextview)
	app.gui.Highlight = true

	app.views = map[string]*AppView{
		// this should be built dynamically, based on forms etc?
		// keep "paint order" in mind? m(
		ViewCommand:   {name: ViewCommand, layoutFunc: app.commandLayoutFunc},
		ViewMain:      {name: ViewMain, layoutFunc: app.mainLayoutFunc},
		ViewHelp:      {name: ViewHelp, layoutFunc: app.helpLayoutFunc},
		ViewLocale:    {name: ViewLocale, layoutFunc: app.localeLayoutFunc},
		ViewShortcuts: {name: ViewShortcuts, layoutFunc: app.hotkeysLayoutFunc},
	}

	app.gui.SetManagerFunc(app.layout)

	app.InitHotkeysWidget()
	if err := app.SetHotkeyKeybindings(); err != nil {
		log.Fatalln(err)
	}

	go app.bugger() // lorem ipsum intro flooding main view

	if err := app.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func (app *App) layout(g *gocui.Gui) error {
	for _, av := range app.views {
		v, err := av.layoutFunc()
		if err != nil {
			return (err)
		}
		av.view = v
	}
	// TODO:
	// add view "blinkenlights" = Status LEDs -> [Alt] pressed? [CAPS]? Download RX/TX?
	// hovers border of main view (one-liner @ bottom, maybe one @top as well - clock? flag, kbd/locale->click?)
	return nil
}
