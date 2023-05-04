package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (app *App) DisplayDialogBool(description string, current bool, onSubmit func(bool) error) {
	app.views[ViewDialog] = &AppView{
		name:         ViewDialog,
		text:         description,
		currentBool:  current,
		onSubmitBool: onSubmit,
		layoutFunc:   app.dialogLayoutFunc,
	}
	app.dialogLayoutFunc() // unsure... needed to have views registered, no?
	if err := app.gui.SetKeybinding(ViewDialog, gocui.KeyEnter, gocui.ModNone, app.SaveDialogBool); err != nil {
		panic(err)
	}
	if err := app.gui.SetKeybinding(ViewDialog, gocui.MouseLeft, gocui.ModNone, app.SaveDialogBool); err != nil {
		panic(err)
	}
	if err := app.gui.SetKeybinding(ViewDialog, gocui.KeyEsc, gocui.ModNone, app.DestroyDialogBool); err != nil {
		panic(err)
	}
	if err := app.gui.SetKeybinding(ViewDialog, gocui.KeyArrowUp, gocui.ModNone, app.dialogBoolUp); err != nil {
		panic(err)
	}
	if err := app.gui.SetKeybinding(ViewDialog, gocui.KeyArrowDown, gocui.ModNone, app.dialogBoolDown); err != nil {
		panic(err)
	}
	app.gui.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(ViewDialog)
		return nil
	})
}

func (app *App) dialogBoolUp(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	if y > 0 {
		v.SetCursor(0, y-1)
	}
	return nil
}

func (app *App) dialogBoolDown(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	if y < 1 {
		v.SetCursor(0, y+1)
	}
	return nil
}

func (app *App) SaveDialogBool(g *gocui.Gui, v *gocui.View) error {
	app.gui.Update(func(g *gocui.Gui) error {
		_, y := v.Cursor()
		val := false
		if y == 1 {
			val = true
		}
		app.views[ViewDialog].onSubmitBool(val)
		return nil
	})
	app.currentMenuItem = 0
	return nil
}

func (app *App) DestroyDialogBool(g *gocui.Gui, v *gocui.View) error {
	app.currentMenuItem = 0
	app.gui.DeleteView(ViewDialog)
	app.gui.DeleteKeybindings(ViewDialog)
	app.gui.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(ViewCommand)
		return nil
	})
	delete(app.views, ViewDialog)
	return nil
}

func (app *App) dialogLayoutFunc() (*gocui.View, error) {
	g := app.gui
	d := app.views[ViewDialog]
	maxX, _ := g.Size()
	if v, err := g.SetView(ViewDialog, 30, 4, maxX-30, 7); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Wrap = false
		v.Frame = true
		v.Title = "[ " + d.text + " ]"
		v.Highlight = true
		v.Autoscroll = false
		fmt.Fprintln(v, "False          ")
		fmt.Fprintln(v, "True           ")
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if d.currentBool {
			v.SetCursor(0, 1)
		}
		return v, nil
	}
	return nil, nil
}
