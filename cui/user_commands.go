package cui

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

func (app *App) userCommandExecutor() {
	time.Sleep(250 * time.Millisecond)
	e, _ := app.gui.View(ViewMain)
	actionMap := map[string]func(){
		"cls": func() {
			app.gui.Update(func(g *gocui.Gui) error {
				e.Clear()
				e.SetOrigin(0, 0)
				return nil
			})
		},
		"quit": func() {
			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintf(e, "\n> OK \033[33;1mbye, peace!\033[0m\n\n")
				return nil
			})
			time.Sleep(1 * time.Second)
			app.gui.Update(func(g *gocui.Gui) error {
				return gocui.ErrQuit
			})
		},
	}
	for {
		x := <-app.userCommands
		if x == "" {
			continue
		}
		// log received command to main view
		fmt.Fprintf(e, "< \033[33;1m%s\033[0m\n", x)
		if f, ok := actionMap[x]; ok {
			f()
		} else {
			// next ... command with arguments m(

			// for now ...
			app.gui.Update(func(g *gocui.Gui) error {
				cmds := "cls, quit" // cough.
				fmt.Fprintf(e, "> Unknown command: \033[31;1m%s\033[0m\n", x)
				fmt.Fprintf(e, "> Known commands: \033[33;1m%s\033[0m\n", cmds)
				fmt.Fprintln(e, "> AI tip of the day: \033[37;1mPress F1 for help\033[0m")
				return nil
			})
		}
	}
}
