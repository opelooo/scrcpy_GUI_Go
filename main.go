package main

// #cgo linux LDFLAGS: -L${SRCDIR}/lib
import "C"

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
)

//go:embed icon.png
var iconData []byte

func main() {
	a := app.NewWithID("com.opelooo.scrcpygui")
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("Minimal GUI Scrcpy")
	if len(iconData) > 0 {
		appIcon := fyne.NewStaticResource("icon.png", iconData)
		a.SetIcon(appIcon)
		w.SetIcon(appIcon)
	}

	w.Resize(fyne.NewSize(670, 360))
	w.SetFixedSize(true)

	runner, err := NewADBRunner()
	if err != nil {
		dialog.ShowError(err, w)
		w.ShowAndRun()
		return
	}

	_ = NewGUI(a, w, runner)
	w.CenterOnScreen()
	w.ShowAndRun()
}
