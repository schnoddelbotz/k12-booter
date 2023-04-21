package cui

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

func (app *App) userCommandExecutor() {
	time.Sleep(250 * time.Millisecond)
	e := app.views[ViewMain].view
	for {
		x := <-app.userCommands
		// TODO handle commands ... map?
		if x == "quit\n" {
			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintf(e, "> \033[33;1mbye, peace!\033[0m")
				return nil
			})
			time.Sleep(1 * time.Second)
			app.gui.Update(func(g *gocui.Gui) error {
				return gocui.ErrQuit
			})
		}
		app.gui.Update(func(g *gocui.Gui) error {
			fmt.Fprintf(e, "> \033[33;1mTODO\033[0m exec user command: \033[31;1m%s\033[0m", x)
			return nil
		})
	}
}
