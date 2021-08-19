package app

import (
	_ "embed"
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sagernet/materia-gtk3-theme-go"
	"github.com/sagernet/sagernetdesktop/core"
	"github.com/sagernet/sagernetdesktop/core/gtkx"
	"github.com/sagernet/sagernetdesktop/res/drawable"
	"os"
)

var mainApp *SagerApplication

type SagerApplication struct {
	*gtk.Application
}

func (a *SagerApplication) onCreate() {
	gtk.WindowSetDefaultIcon(drawable.MustPixbuf(drawable.IconPng))

	initMain()
}

func Launch() {
	cacheDir, err := os.UserCacheDir()
	core.Must("get user cache dir", err)
	cacheDir = fmt.Sprint(cacheDir, "/SagerNet")
	core.Must("init themes", materia_gtk3_theme_go.InitTheme(cacheDir))

	gtk.Init(nil)
	settings, err := gtk.SettingsGetDefault()
	core.Must("get gtk settings", err)
	err = settings.Set("gtk-theme-name", "Materia")
	core.Must("set gtk theme", err)

	mainApp = &SagerApplication{}
	app, err := gtk.ApplicationNew("io.nekohasekai.sagernet", glib.APPLICATION_FLAGS_NONE)
	gtkx.Must("create main window", err)
	mainApp.Application = app
	mainApp.Connect("activate", mainApp.onCreate)
	mainApp.Run(os.Args)
}
