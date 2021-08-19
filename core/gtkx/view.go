package gtkx

import "github.com/gotk3/gotk3/gtk"

type View struct {
	*gtk.Widget
}

func (v View) SetOnClickListener(listener func()) {
	v.Widget.Connect("clicked", listener)
}

type MenuItem struct {
	View
	*gtk.MenuItem
}

func NewMenuItem(v *gtk.MenuItem) MenuItem {
	return MenuItem{
		View{v.ToWidget()},
		v,
	}
}

func (v MenuItem) SetOnClickListener(listener func()) {
	v.Connect("activate", listener)
}
