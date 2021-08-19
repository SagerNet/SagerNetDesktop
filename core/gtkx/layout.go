package gtkx

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Layout struct {
	B *gtk.Builder
}

func (l *Layout) GetObject(name string) glib.IObject {
	obj, err := l.B.GetObject(name)
	Mustf("get obj %s", err, name)
	return obj
}

func (l *Layout) GetAppWindow(name string) *gtk.ApplicationWindow {
	return l.GetObject(name).(*gtk.ApplicationWindow)
}

func (l *Layout) GetWindow(name string) *gtk.Window {
	return l.GetObject(name).(*gtk.Window)
}

func (l *Layout) GetDialog(name string) *gtk.Dialog {
	return l.GetObject(name).(*gtk.Dialog)
}

func (l *Layout) GetAboutDialog(name string) *gtk.AboutDialog {
	return l.GetObject(name).(*gtk.AboutDialog)
}

func NewLayout(layout string) *Layout {
	builder, err := gtk.BuilderNewFromString(layout)
	Must("create layout builder", err)
	return &Layout{builder}
}

func (l *Layout) GetLabel(name string) *gtk.Label {
	return l.GetObject(name).(*gtk.Label)
}

func (l *Layout) GetMenuItem(name string) MenuItem {
	return NewMenuItem(l.GetObject(name).(*gtk.MenuItem))
}
