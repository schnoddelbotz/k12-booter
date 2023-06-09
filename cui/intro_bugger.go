package cui

import (
	"fmt"
	"log"
	"os/user"
	"schnoddelbotz/k12-booter/diagnostics"
	"schnoddelbotz/k12-booter/dnssd"
	"schnoddelbotz/k12-booter/sounds"
	"schnoddelbotz/k12-booter/utility"
	"schnoddelbotz/k12-booter/ws"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/spf13/viper"
)

func (app *App) bugger() {
	sounds.PlayIt(sounds.Maelstrom_LaserMetalPling, app.otoCtx)
	time.Sleep(250 * time.Millisecond) // "wait" for main loop to be run once m(
	e, err := app.gui.View(ViewMain)
	if err != nil {
		log.Fatal(err)
	}

	app.gui.Update(func(g *gocui.Gui) error {
		e.Autoscroll = true
		g.SetCurrentView(ViewCommand)
		e.Write([]byte("\n *** NOTA BENE ***\n"))
		e.Write([]byte(" Enterprise version does this twice as fast\n"))
		e.Write([]byte(" Press F2 for menu, F10 to quit\n\n"))
		fmt.Fprintf(e, "\033[37;1m%s\033[0m", `         ____________ 
        < k12-booter >
         ------------ 
                \   ^__^
                 \  (oo)\_______
                    (__)\       )\/\
                        ||----w |
                        ||     ||`+"\n\n")
		return nil
	})
	sounds.PlayIt(sounds.Maelstrom_LaserZap, app.otoCtx)

	si := diagnostics.GetSysInfoData()
	if !viper.GetBool(ConfigDisablePOST) {
		app.gui.Update(func(g *gocui.Gui) error {
			fmt.Fprint(e, "Running POST diagnostics ...\n")
			return nil
		})
		go sounds.PlayIt(sounds.Maelstrom_Warp, app.otoCtx)
		time.Sleep(100 * time.Millisecond)
		diagnostics.RunPOST(e)
	}
	go sounds.PlayIt(sounds.Maelstrom_AllRight, app.otoCtx)

	if viper.GetBool(ConfigEnableTeacher) {
		u, err := user.Current()
		utility.Fatal(err)
		app.gui.Update(func(g *gocui.Gui) error {
			fmt.Fprintf(e, "Enabling DNS-SD/bonjour discovery of this node as teacher '%s'\n", u.Username)
			return nil
		})
		go dnssd.RegisterTeacherService(u.Username)

		app.gui.Update(func(g *gocui.Gui) error {
			fmt.Fprintf(e, "Running teacher websocket server on port '%s'\n", ":8888")
			return nil
		})
		app.chatServer = ws.NewChatServer()
		go ws.RunWebsocketServer(":8888", app.chatServer)
	}

	time.Sleep(100 * time.Millisecond)
	app.gui.Update(func(g *gocui.Gui) error {
		fmt.Fprintf(e, "\n> k12-booter v%s (%s/%s) is ready\n", app.version, si.OS, si.Architecture)
		return nil
	})
}
