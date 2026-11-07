package main

// #cgo linux LDFLAGS: -L${SRCDIR}/lib
import "C"

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"github.com/opelooo/scrcpy_GUI_Go/internal/platform/adb"
	"github.com/opelooo/scrcpy_GUI_Go/internal/service"
	"github.com/opelooo/scrcpy_GUI_Go/internal/ui"
)

//go:embed icon.png
var iconData []byte

const (
	appName    = "Minimal GUI Scrcpy"
	appVersion = "v1.0.0"
	appID      = "com.opelooo.scrcpygui"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Initialize Fyne desktop application
	a := app.NewWithID(appID)
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow(appName)
	if len(iconData) > 0 {
		appIcon := fyne.NewStaticResource("icon.png", iconData)
		a.SetIcon(appIcon)
		w.SetIcon(appIcon)
	}

	w.Resize(fyne.NewSize(670, 360))
	w.SetFixedSize(true)

	// Dependency wiring
	client, err := adb.NewClient()
	if err != nil {
		dialog.ShowError(err, w)
		w.CenterOnScreen()
		w.ShowAndRun()
		return
	}

	deviceService := service.NewScrcpyService(client)
	_ = ui.NewGUI(a, w, deviceService)

	w.CenterOnScreen()
	w.ShowAndRun()
}
