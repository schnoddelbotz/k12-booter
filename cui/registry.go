package cui

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
	// all needed to set up cui hotkey
}

func (r *Registry) Init() {
	r.register(&helpView)
}

func (r *Registry) register(bv *BooterView) error {
	// ensure it did not exist. append to r.Views
	return nil
}

// unregister: delete keybindings ...
