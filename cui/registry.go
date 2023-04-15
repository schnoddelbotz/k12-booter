package cui

import (
	"github.com/jroimartin/gocui"
)

// View registry.
// Auto-generated form code must register here.
// All the views. I hope this is a good idea.

type Registry struct {
	Views []BooterView
}

type BooterView struct {
	Name    string
	Title   string
	HotKeys []HotKey
}

type HotKey struct {
	ViewName string // F1 .. F10
	Label    string
	Key      gocui.Key
	Modifier gocui.Modifier
	Handler  func(g *gocui.Gui, v *gocui.View) error
}

func (r *Registry) Init() {
	r.register(&helpView)
	r.register(&hotkeysWidget)
}

func (r *Registry) register(bv *BooterView) error {
	// ensure it did not exist. append to r.Views
	return nil
}

// unregister: delete keybindings ...
