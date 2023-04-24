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
	Label_CLS  = "CLS"
)

func (app *App) voidKeyHandler(g *gocui.Gui, v *gocui.View) error {
	log.Printf("VOID hotkeyhandler called [FIXME:hotkeys.go]")
	return nil
}

func (app *App) keyHandlerHelp(g *gocui.Gui, v *gocui.View) error {
	app.showHelp = !app.showHelp // toggle
	return nil
}

func (app *App) keyHandlerToggleMenu(g *gocui.Gui, v *gocui.View) error {
	// todo
	//app.views[ViewMain].view.Clear() // why? don't?
	//e := app.views[ViewMain].view
	if _, ok := app.views[ViewMenu]; ok {
		app.hideMenu()
	} else {
		app.showMenu(app.currentMenu)
	}
	return nil
}

func (app *App) showMenu(menuName string) {
	app.views[ViewMenu] = &AppView{
		name:       ViewMenu,
		layoutFunc: app.menuLayoutFunc,
	}
	/// FIXME ^^^
	/// take it into intial list and then delete it on startup again?
	/// or create .view prior to layout()

	app.currentMenu = menuName
	app.setupMenuKeybindings()
	app.gui.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(ViewMenu)
		if v, err := g.View(ViewMenu); err == nil {
			v.SetCursor(0, app.currentMenuItem)
		}
		return nil
	})
}

func (app *App) hideMenu() {
	app.gui.DeleteKeybindings(ViewMenu)
	app.gui.Update(func(g *gocui.Gui) error {
		app.gui.DeleteView(ViewMenu)
		return nil
	})
	delete(app.views, ViewMenu)
	app.gui.Update(func(g *gocui.Gui) error {
		app.gui.SetCurrentView(ViewCommand)
		return nil
	})
}

func (app *App) switchMenu(menuName string) {
	app.currentMenu = menuName
	app.currentMenuItem = 0
	app.gui.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(ViewMenu)
		if v, err := g.View(ViewMenu); err == nil {
			v.SetCursor(0, 0)
			// HACK
			g.DeleteView(ViewMenu)
			g.DeleteKeybindings(ViewMenu)
			app.showMenu(app.currentMenu)
		}
		return nil
	})
}

func (app *App) keyHandlerQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (app *App) keyHandlerMask(g *gocui.Gui, v *gocui.View) error {
	// try on view "" using F9 ;) ... Fine via mouse. Demo relict.
	v.Mask ^= '*'
	return nil
}

func (app *App) sendUserCommandCLS(g *gocui.Gui, v *gocui.View) error {
	app.userCommands <- "cls"
	return nil
}

func (app *App) setActiveViewCommand(g *gocui.Gui, v *gocui.View) error {
	app.setActiveView(ViewCommand)
	return nil
}

func (app *App) setActiveViewMain(g *gocui.Gui, v *gocui.View) error {
	app.setActiveView(ViewMain)
	return nil
}

func (app *App) setActiveView(id string) error {
	app.gui.Update(func(g *gocui.Gui) error {
		app.gui.SetCurrentView(id)
		return nil
	})
	return nil
}

func (app *App) mainScrollUp(g *gocui.Gui, v *gocui.View) error {
	app.gui.Update(func(g *gocui.Gui) error {
		v.MoveCursor(0, 0, false)
		return nil
	})
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
			{ViewName: HotKeyView_F2, Key: gocui.KeyF2, Label: Label_Menu, Handler: app.keyHandlerToggleMenu},
			{ViewName: HotKeyView_F3, Key: gocui.KeyF3, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F4, Key: gocui.KeyF4, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F5, Key: gocui.KeyF5, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F6, Key: gocui.KeyF6, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F7, Key: gocui.KeyF7, Handler: app.voidKeyHandler},
			{ViewName: HotKeyView_F8, Key: gocui.KeyF8, Label: Label_CLS, Handler: app.sendUserCommandCLS},
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
	if err := app.gui.SetKeybinding(ViewMain, gocui.KeyArrowUp, gocui.ModNone, app.mainScrollUp); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(ViewMain, gocui.KeyTab, gocui.ModNone, app.setActiveViewCommand); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(ViewCommand, gocui.KeyTab, gocui.ModNone, app.setActiveViewMain); err != nil {
		return err
	}
	return nil
}
