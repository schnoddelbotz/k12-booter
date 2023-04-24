package cui

import (
	"fmt"
	"log"
	"schnoddelbotz/k12-booter/internationalization"
	"strings"
	"time"

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
	gui             *gocui.Gui
	views           map[string]*AppView
	hotkeysWidget   *Widget
	helpView        *gocui.View // sidebar; should toggle with F1
	showHelp        bool
	helpWidth       int
	localeInfo      internationalization.LocaleInfo
	currentMenu     string
	currentMenuItem int
	userCommands    chan string
	version         string
}

type ViewIdentifier string

type AppView struct {
	name       ViewIdentifier
	hotkeys    []HotKey
	layoutFunc func() (*gocui.View, error)
}

func Zain(version string) {
	var (
		app App
		err error
	)
	app.version = version
	app.printBanner()
	app.localeInfo = internationalization.GetLocaleInfo()
	app.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.gui.Close()
	app.gui.Cursor = true
	app.gui.Mouse = true
	app.gui.ASCII = true
	app.gui.InputEsc = true
	// next two (?) lines will enable color upon g.SetCurrentView(nextview)
	app.gui.Highlight = true
	app.userCommands = make(chan string)

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
	go app.userCommandExecutor()

	if err := app.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func (app *App) layout(g *gocui.Gui) error {
	for _, av := range app.views {
		_, err := av.layoutFunc()
		if err != nil {
			return (err)
		}
	}
	// TODO:
	// add view "blinkenlights" = Status LEDs -> [Alt] pressed? [CAPS]? Download RX/TX?
	// hovers border of main view (one-liner @ bottom, maybe one @top as well - clock? flag, kbd/locale->click?)
	return nil
}

func (app *App) printBanner() {
	const k12booterBanner = `🇦🇫🇦🇽           🇦🇱🇩🇿       🇦🇸🇦🇩
🇦🇴      🇦🇮  🇦🇶  🇦🇬🇦🇷       🇦🇲
🇦🇼      🇦🇺      🇦🇹       🇦🇿
🇧🇸     🇧🇭🇧🇩      🇧🇧🇧🇾      🇧🇪
🇧🇿     🇧🇯🇧🇲      🇧🇹🇧🇴      🏳️
🇧🇦     🇧🇼🇧🇻      🇧🇷🇮🇴      🇧🇳                 🇧🇬
🇧🇫     🇧🇮🏳️      🇰🇭🇨🇲      🇨🇦                 🇰🇾
🇨🇫 🇹🇩🇨🇱  🇨🇳🇨🇽      🇨🇨       🇨🇴 🇰🇲🇨🇩   🇨🇬🇨🇰🇨🇷   🇨🇮🇭🇷🇨🇺 🇨🇼🇨🇾🇨🇿   🇩🇰🇩🇯🇩🇲 🇩🇴🇪🇨 🇪🇬🇸🇻
🇬🇶 🇪🇷   🇪🇪🇸🇿     🇪🇹🇫🇰       🇫🇴🇫🇯 🇫🇮🇫🇷 🇬🇫  🇵🇫  🇹🇫  🇬🇦  🇬🇲   🇬🇪  🇩🇪  🇬🇭🇬🇮 🇬🇷🇬🇱
🇬🇩     🇬🇵🇬🇺    🇬🇹🇬🇬        🇬🇳   🇬🇼 🇬🇾   🇭🇹 🇭🇲   🇻🇦 🇭🇳   🇭🇰  🇭🇺🇮🇸 🇮🇳  🇮🇩
🇮🇷🇮🇶    🇮🇪🇮🇲    🇮🇱   🇮🇹🇯🇲🇯🇵🇯🇪🇯🇴 🇰🇿   🇰🇪🇰🇮🇰🇵   🇰🇷🇰🇼🇰🇬   🇱🇦 🇱🇻  🇱🇧🇱🇸   🇱🇷 🇱🇾
🇱🇮🇱🇹🇱🇺   🇲🇴🇲🇰   🇲🇬          🇲🇼   🇲🇾🇲🇻🇲🇱   🇲🇹🇲🇭🇲🇶   🇲🇷 🇲🇺  🇾🇹🇲🇽🇫🇲🇲🇩🇲🇨🇲🇳 🇲🇪
🇲🇸 🇲🇦   🇲🇿🇲🇲              🇳🇦   🇳🇷🇳🇵🇳🇱   🇳🇨🇳🇿🇳🇮   🇳🇪 🇳🇬  🇳🇺🇳🇫     🇲🇵
🇳🇴 🇴🇲🇵🇰  🇵🇼🇵🇸  🇵🇦    🇵🇬      🇵🇾   🇵🇪 🇵🇭   🇵🇳 🇵🇱   🇵🇹 🇵🇷   🇶🇦     🇷🇪
🇷🇴  🇷🇺  🇷🇼🇧🇱  🇹🇦🇰🇳🇱🇨🇲🇫🇵🇲🇻🇨      🇼🇸  🇸🇲🏳️ 🇸🇦  🇸🇳  🇷🇸  🇸🇨  🇸🇱🇸🇬  🇸🇽     🇸🇰
🇸🇮  🇸🇧🇸🇴 🇿🇦🇬🇸  🇸🇸🇪🇸🇱🇰🇸🇩🇸🇷       🇸🇯 🇸🇪🇨🇭   🇸🇾🇹🇼🇹🇯   🇹🇿🇹🇭🇹🇱   🇹🇬🇹🇰  🇹🇴 🇹🇹  🇹🇳`
	fmt.Printf("Welcome to k-12 booter version %s\n", app.version)
	lines := strings.Split(k12booterBanner, "\n")
	for i, l := range lines {
		fmt.Println(l)
		time.Sleep(time.Duration(i*5) * time.Millisecond)
	}
	time.Sleep(250 * time.Millisecond)
}
