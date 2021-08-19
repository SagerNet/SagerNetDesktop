package gtkx

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sagernet/sagernetdesktop/core"
	"github.com/sagernet/sagernetdesktop/res/layout"
	"github.com/xjasonlyu/tun2socks/log"
)

func Must(action string, err error) {
	if err != nil {
		builder, errN := gtk.BuilderNewFromString(layout.DialogError)
		if errN != nil {
			log.Errorf("failed to create gtkx builder: %v", errN)
			core.Must(action, err)
			return
		}

		obj, errN := builder.GetObject("dialog")
		if errN != nil {
			log.Errorf("failed to get error dialog: %v", errN)
			core.Must(action, err)
			return
		}

		dialog := *obj.(*gtk.Dialog)

		obj, errN = builder.GetObject("message")
		if errN != nil {
			log.Errorf("failed to get error label: %v", errN)
			core.Must(action, err)
			return
		}

		message := *obj.(*gtk.Label)
		message.SetText(fmt.Sprintf("Failed to %s: %v", action, err))

		obj, errN = builder.GetObject("close")
		if errN != nil {
			log.Errorf("failed to get close button: %v", errN)
			core.Must(action, err)
			return
		}

		closeButton := *obj.(*gtk.Button)
		closeButton.Connect("clicked", dialog.Close)

		dialog.Show()
	}
}

func Mustf(action string, err error, args ...interface{}) {
	if err != nil {
		Mustf(fmt.Sprintf(action, args), err)
	}
}
