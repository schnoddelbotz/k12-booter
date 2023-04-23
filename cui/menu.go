package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type MenuItemType int

const (
	LinkMenu MenuItemType = iota
	LinkUserCommand
	LinkFunc
	LinkShellCommand
)

type MenuItem struct {
	Label    string // i18n ?!?!
	Parent   string
	itemType MenuItemType
	target   string
	linkFunc *func() error
	// attributes ...
}

const (
	Menu_Main         = "k12-booter Main Menu"
	Menu_Preferences  = "Preferences"
	Menu_Admin        = "Administration"
	Menu_Applications = "Applications"
	Menu_Deployment   = "Lab management & deployment"
)

// AppMenu wants to be outsourced and auto-generated.
var AppMenu = []MenuItem{
	{
		Label:    "Discover applications", // launch or install
		Parent:   Menu_Main,
		itemType: LinkMenu,
		target:   Menu_Applications,
	},
	{
		Label:    "School administration", // API ? Buckets ? /country/canton|bundesland|state.../city/zip/schoolname
		Parent:   Menu_Main,
		itemType: LinkMenu,
		target:   Menu_Admin,
	},
	{
		Label:    "Manage & Deploy labs",
		Parent:   Menu_Main,
		itemType: LinkMenu,
		target:   Menu_Deployment,
	},
	{
		Label:    "Preferences", // configure API endpoints/params, buckets,
		Parent:   Menu_Main,
		itemType: LinkMenu,
		target:   Menu_Preferences,
		// should write to main if essential stuff not set!
	},
	{
		Label:    "CLS",
		Parent:   Menu_Main,
		itemType: LinkUserCommand,
		target:   "cls",
	},
	{
		Label:    "Quit",
		Parent:   Menu_Main,
		itemType: LinkUserCommand,
		target:   "quit",
	},
	{
		Label:    "Set display language... tbd",
		Parent:   Menu_Preferences,
		itemType: LinkUserCommand,
		target:   "quit", //
	},
	// hmm.
	{
		Label:    "Return to main menu",
		Parent:   Menu_Applications,
		itemType: LinkMenu,
		target:   Menu_Main,
	},
	{
		Label:    "Return to main menu",
		Parent:   Menu_Admin,
		itemType: LinkMenu,
		target:   Menu_Main,
	},
	{
		Label:    "Return to main menu",
		Parent:   Menu_Deployment,
		itemType: LinkMenu,
		target:   Menu_Main,
	},
	{
		Label:    "Return to main menu",
		Parent:   Menu_Preferences,
		itemType: LinkMenu,
		target:   Menu_Main,
	},
}

func (app *App) menuLayoutFunc() (*gocui.View, error) {
	g := app.gui
	if app.currentMenu == "" {
		app.currentMenu = Menu_Main
	}
	_, maxY := g.Size()
	if v, err := g.SetView(ViewMenu, 20, 4, 55, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Title = fmt.Sprintf("[ %s ]", app.currentMenu)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		items := app.menuItems(app.currentMenu)
		for i, mi := range items {
			fmt.Fprintf(v, "%3d > %-32s\n", i, mi.Label)
		}
		return v, nil
	}
	return nil, nil
}

func (app *App) handleMenuKeyUp(g *gocui.Gui, v *gocui.View) error {
	//e := app.views[ViewMain].view
	if app.currentMenuItem > 0 {
		app.currentMenuItem--
		// fmt.Fprintf(e, "# menu debug: UP ↑ -> %d\n", app.currentMenuItem)
	}
	if app.currentMenuItem < len(app.menuItems(app.currentMenu)) {
		// fmt.Fprintf(e, "cur %d len %d\n", app.currentMenuItem, len(app.menuItems(app.currentMenu)))
		if v, err := app.gui.View(ViewMenu); err == nil {
			v.SetCursor(0, app.currentMenuItem)
		}
	}

	return nil
}

func (app *App) handleMenuKeyDown(g *gocui.Gui, v *gocui.View) error {
	//e := app.views[ViewMain].view
	if app.currentMenuItem+1 < len(app.menuItems(app.currentMenu)) {
		app.currentMenuItem++
		// fmt.Fprintf(e, "# menu debug: DN ↓-> %d\n", app.currentMenuItem)
	}
	if app.currentMenuItem < len(app.menuItems(app.currentMenu)) {
		if v, err := app.gui.View(ViewMenu); err == nil {
			v.SetCursor(0, app.currentMenuItem)
		}
	}

	return nil
}

func (app *App) menuEnterKeyHandler(g *gocui.Gui, v *gocui.View) error {
	e, _ := app.gui.View(ViewMain)
	items := app.menuItems(app.currentMenu)
	if app.currentMenuItem <= len(items) {
		clickedItem := items[app.currentMenuItem]
		fmt.Fprintf(e, "# menu debug: ENTER -> %d = %v\n", app.currentMenuItem, clickedItem)
		app.gui.Update(func(*gocui.Gui) error { return nil })
		return app.execMenuItemAction(clickedItem)
	}
	fmt.Fprintf(e, "# menu debug: IGNORE ENTER -> %d\n", app.currentMenuItem)
	return nil
}

func (app *App) menuMouseClickHandler(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if _, err := v.Line(cy); err != nil {
		// ignore click on empty lines
		return nil
	}
	items := app.menuItems(app.currentMenu)
	if cy+1 > len(items) {
		return nil
	}
	clickedItem := items[cy]
	// fmt.Fprintf(e, "> *clicked* %d %s %v\n", cy, l, clickedItem)
	return app.execMenuItemAction(clickedItem)
}

func (app *App) execMenuItemAction(mi MenuItem) error {
	e, _ := app.gui.View(ViewMain)
	switch mi.itemType {
	case LinkUserCommand:
		app.userCommands <- mi.target
		// todo: close menu now?
	case LinkMenu:
		fmt.Fprintf(e, "# TODO Scotty, beam me into sub-menu %s\n", mi.Label)
		//app.showMenu(mi.target)
		app.switchMenu(mi.target)
		//app.currentMenu = mi.target
		//app.currentMenuItem = 0
		//app.gui.SetCurrentView(ViewMenu)
	default:
		fmt.Fprintf(e, "# MI_EXEC TODO impl %+v\n", mi)
	}
	return nil
}

func (app *App) menuItems(menuName string) (result []MenuItem) {
	for _, mi := range AppMenu {
		if mi.Parent != menuName {
			continue
		}
		result = append(result, mi)
	}
	return
}

func (app *App) setupMenuKeybindings() error {
	if err := app.gui.SetKeybinding(ViewMenu, gocui.KeyArrowUp, gocui.ModNone, app.handleMenuKeyUp); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(ViewMenu, gocui.KeyArrowDown, gocui.ModNone, app.handleMenuKeyDown); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(ViewMenu, gocui.KeyEnter, gocui.ModNone, app.menuEnterKeyHandler); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(ViewMenu, gocui.MouseLeft, gocui.ModNone, app.menuMouseClickHandler); err != nil {
		return err
	}
	// 0..9 ?
	return nil
}
