package cui

import (
	"fmt"
	"log"
	"schnoddelbotz/k12-booter/diagnostics"
	"schnoddelbotz/k12-booter/sounds"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

func (app *App) userCommandExecutor() {
	time.Sleep(250 * time.Millisecond)
	e, _ := app.gui.View(ViewMain)
	actionMap := map[string]func(){
		"cls": func() {
			go sounds.PlayIt(sounds.Maelstrom_AlertSound, app.otoCtx)
			app.gui.Update(func(g *gocui.Gui) error {
				e.Clear()
				e.SetOrigin(0, 0)
				return nil
			})
		},
		"quit": func() {
			go sounds.PlayIt(sounds.Maelstrom_Yahoo, app.otoCtx)
			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintf(e, "\n> OK \033[33;1mbye, peace!\033[0m\n\n")
				return nil
			})
			sounds.PlayIt(sounds.Maelstrom_Laughing, app.otoCtx)
			time.Sleep(1 * time.Second)
			app.gui.Update(func(g *gocui.Gui) error {
				return gocui.ErrQuit
			})
		},
		"sysinfo": func() {
			go sounds.PlayIt(sounds.Maelstrom_Yo, app.otoCtx)
			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintln(e, diagnostics.SysInfo())
				return nil
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
		} else if strings.HasPrefix(x, "all") || x == "who" {
			if app.chatServer != nil {
				if x == "who" {
					app.gui.Update(func(g *gocui.Gui) error {
						who := app.chatServer.Who()
						fmt.Fprintf(e, "> \033[37;1mwho\033[0m - %d users connected.\n", len(who))
						for _, w := range who {
							fmt.Fprintf(e, "> %s - details tbd\n", w)
						}
						return nil
					})
					continue
				}
				// todo HANDLE input, publish if OK ... and let clients act upon. TRUST!!!?
				app.chatServer.Publish([]byte("haha - teacher command " + x))
			} else {
				fmt.Fprintf(e, "> Error: \033[31;1m%s\033[0m not enabled in preferences\n", "bonjour/teacher-mode")
			}
		} else if strings.HasPrefix(x, "join") {
			if app.chatServer != nil {
				log.Printf("teacher shall not / can not join other class")
				continue
			}
			log.Println("Joining first teacher found on local network ...")
			log.Println("Note: PoC. Not secure. Not ready for general use.")
			log.Println("Should provide chooser, if no /well-known/ found or so.")
			// todo:
			// - browse / async
			// - connct to first/given browse result /subscribe, process commands from teacher
		} else {
			// for now ...
			sounds.PlayIt(sounds.Maelstrom_MaleOop, app.otoCtx)
			app.gui.Update(func(g *gocui.Gui) error {
				cmds := "cls, sysinfo, quit" // cough.
				fmt.Fprintf(e, "> Unknown command: \033[31;1m%s\033[0m\n", x)
				fmt.Fprintf(e, "> Known commands: \033[33;1m%s\033[0m\n", cmds)
				fmt.Fprintln(e, "> AI tip of the day: \033[37;1mPress F1 for help\033[0m")
				return nil
			})
		}
	}
}
