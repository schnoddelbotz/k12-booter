package cui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

const (
	HotKeyView_F1  = "F1"
	HotKeyView_F2  = "F2"
	HotKeyView_F3  = "F3"
	HotKeyView_F4  = "F4"
	HotKeyView_F5  = "F5"
	HotKeyView_F6  = "F6"
	HotKeyView_F7  = "F7"
	HotKeyView_F8  = "F8"
	HotKeyView_F9  = "F9"
	HotKeyView_F10 = "F10"

	Label_Help = "Help"
	Label_Mask = "Mask"
	Label_Quit = "Quit"
)

var hotkeysWidget = BooterView{
	Title: "",
	Name:  "",
	HotKeys: []HotKey{
		{ViewName: HotKeyView_F1, Key: gocui.KeyF1, Label: Label_Help},
		{ViewName: HotKeyView_F2, Key: gocui.KeyF2},
		{ViewName: HotKeyView_F3, Key: gocui.KeyF3},
		{ViewName: HotKeyView_F4, Key: gocui.KeyF4},
		{ViewName: HotKeyView_F5, Key: gocui.KeyF5},
		{ViewName: HotKeyView_F6, Key: gocui.KeyF6},
		{ViewName: HotKeyView_F7, Key: gocui.KeyF7},
		{ViewName: HotKeyView_F8, Key: gocui.KeyF8},
		{ViewName: HotKeyView_F9, Key: gocui.KeyF9, Label: Label_Mask},
		{ViewName: HotKeyView_F10, Key: gocui.KeyF10, Label: Label_Quit},
	},
}

func (app *App) voidKeyHandler(g *gocui.Gui, v *gocui.View) error {
	log.Printf("VOID hotkeyhandler called [FIXME:hotkeys.go]")
	return nil
}

func (app *App) keyHandlerHelp(g *gocui.Gui, v *gocui.View) error {
	// they must depend on current view ...?
	// ESC should ...? Get into main view and let it consume further ESCs.
	// F2 in mainview should select CMD as active win

	//if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone,
	//	func(g *gocui.Gui, v *gocui.View) error {
	app.showHelp = !app.showHelp // toggle
	//		return nil
	//	}); err != nil {
	//	return err
	//}
	return nil
}

func (app *App) keyHandlerQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (app *App) keyHandlerMask(g *gocui.Gui, v *gocui.View) error {
	// try on view "" using F9 ;) ... Fine via mouse. Demo relict.
	v.Mask ^= '*'
	return nil
}

func (app *App) SetHotkeyKeybindings() error {
	for _, key := range hotkeysWidget.HotKeys {
		handler := app.voidKeyHandler
		switch key.ViewName {
		case HotKeyView_F1:
			handler = app.keyHandlerHelp
		case HotKeyView_F9:
			handler = app.keyHandlerMask
		case HotKeyView_F10:
			handler = app.keyHandlerQuit
		}
		if err := app.gui.SetKeybinding("" /* global hotkey */, key.Key, gocui.ModNone, handler); err != nil {
			return err
		}
		if err := app.gui.SetKeybinding(key.ViewName, gocui.MouseLeft, gocui.ModNone, handler); err != nil {
			return err
		}
	}
	return nil
}

func LayoutHotkeys(g *gocui.Gui) error {
	_, maxY := g.Size()
	for i, key := range hotkeysWidget.HotKeys {
		if v, err := g.SetView(key.ViewName, i*9, maxY-2, i*9+9, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.BgColor = gocui.ColorBlue
			v.FgColor = gocui.ColorWhite
			v.Frame = false
			// https://github.com/gczgcz2015/gocui/blob/master/_examples/colors256.go
			// NC was: black empty bg instead of | char
			fmt.Fprintf(v, "%s %s", key.ViewName, key.Label)
		}
	}
	return nil
}
