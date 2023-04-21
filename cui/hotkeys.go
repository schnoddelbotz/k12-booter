package cui

import (
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type HotKey struct {
	ViewName string // F1 .. F10
	Label    string
	Key      gocui.Key
	Modifier gocui.Modifier
	Handler  func(g *gocui.Gui, v *gocui.View) error
}

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
	Label_Menu = "Menu"
)

func (app *App) voidKeyHandler(g *gocui.Gui, v *gocui.View) error {
	log.Printf("VOID hotkeyhandler called [FIXME:hotkeys.go]")
	return nil
}

func (app *App) keyHandlerHelp(g *gocui.Gui, v *gocui.View) error {
	app.showHelp = !app.showHelp // toggle
	return nil
}

func (app *App) keyHandlerMainMenu(g *gocui.Gui, v *gocui.View) error {
	// todo
	//app.views[ViewMain].view.Clear() // why? don't?
	if _, ok := app.views[ViewMenu]; ok {
		app.gui.DeleteView(ViewMenu)
		delete(app.views, ViewMenu)
	} else {
		app.views[ViewMenu] = &AppView{name: ViewMenu, layoutFunc: app.menuLayoutFunc}
	}
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

func (app *App) handleUserCommand(g *gocui.Gui, v *gocui.View) error {
	app.userCommands <- strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

func (app *App) InitHotkeysWidget() {
	app.hotkeysWidget = &Widget{
		Title: "",
		Name:  "",
		HotKeys: []HotKey{
			{ViewName: HotKeyView_F1, Key: gocui.KeyF1, Label: Label_Help, Handler: app.keyHandlerHelp},
			{ViewName: HotKeyView_F2, Key: gocui.KeyF2, Label: Label_Menu, Handler: app.keyHandlerMainMenu},
			{ViewName: HotKeyView_F3, Key: gocui.KeyF3, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F4, Key: gocui.KeyF4, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F5, Key: gocui.KeyF5, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F6, Key: gocui.KeyF6, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F7, Key: gocui.KeyF7, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F8, Key: gocui.KeyF8, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F9, Key: gocui.KeyF9, Label: Label_Mask, Handler: app.keyHandlerMask},
			{ViewName: HotKeyView_F10, Key: gocui.KeyF10, Label: Label_Quit, Handler: app.keyHandlerQuit},
		},
	}
}

func (app *App) SetHotkeyKeybindings() error {
	// imaginable? if app.View == main, then:
	for _, key := range app.hotkeysWidget.HotKeys {
		if err := app.gui.SetKeybinding("" /* global hotkey */, key.Key, gocui.ModNone, key.Handler); err != nil {
			return err
		}
		if err := app.gui.SetKeybinding(key.ViewName, gocui.MouseLeft, gocui.ModNone, key.Handler); err != nil {
			return err
		}
	}
	if err := app.gui.SetKeybinding(ViewCommand, gocui.KeyEnter, gocui.ModNone, app.handleUserCommand); err != nil {
		return err
	}
	return nil
}
