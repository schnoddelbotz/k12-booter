package cui

import (
	"context"
	"fmt"
	"schnoddelbotz/k12-booter/diagnostics"
	"schnoddelbotz/k12-booter/dnssd"
	"schnoddelbotz/k12-booter/sounds"
	"schnoddelbotz/k12-booter/utility"
	"schnoddelbotz/k12-booter/ws"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/pkg/browser"
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
				url, _ := strings.CutPrefix(x, "all ")
				app.chatServer.Publish([]byte(url))
			} else {
				fmt.Fprintf(e, "> Error: \033[31;1m%s\033[0m not enabled in preferences\n", "bonjour/teacher-mode")
			}
		} else if strings.HasPrefix(x, "join") {
			if app.chatServer != nil {
				app.gui.Update(func(g *gocui.Gui) error {
					fmt.Fprintln(e, "> teacher shall not / can not join other class")
					return nil
				})
				continue
			}
			if app.wsClientConn != nil {
				app.gui.Update(func(g *gocui.Gui) error {
					fmt.Fprintln(e, "> Cannot join more - already connected.")
					return nil
				})
				continue
			}
			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintln(e, "> PoC: Joining first teacher found on local network - browsing ...")
				return nil
			})
			x := dnssd.BrowseSingle()
			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintf(e, "> Connecting to teacher %s at %s:%d\n", x.Instance, x.AddrIPv4, x.Port)
				return nil
			})
			// conntect to ws IP now ...
			var err error
			app.wsClientConn, err = ws.ConnectToTeacherWS(fmt.Sprintf("http://%s:%d", x.AddrIPv4, x.Port))
			utility.Fatal(err)

			app.gui.Update(func(g *gocui.Gui) error {
				fmt.Fprintln(e, "> Client connected successfully")
				return nil
			})

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
				defer cancel()

				for {
					_, message, err := app.wsClientConn.Read(ctx)
					if err != nil {
						utility.Fatal(err)
					}
					// log.Printf("Received %v", message)
					app.gui.Update(func(g *gocui.Gui) error {
						fmt.Fprintf(e, "> Received %v\n", string(message))
						// HORROR DEMO FIXME SECURITY ETC. Thanks.
						browser.OpenURL(string(message))
						return nil
					})
				}
			}()
			// log.Println("Note: PoC. Not secure. Not ready for general use.")
			// log.Println("Should provide chooser, if no /well-known/ found or so.")
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
