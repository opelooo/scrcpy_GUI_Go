package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initDeviceWidgets() {
	g.screenOnCheck = widget.NewCheck("Turn Screen On on Start", nil)
	g.screenOnCheck.SetChecked(true)

	g.stayAwakeCheck = widget.NewCheck("Stay Awake (While Plugged In)", nil)
	g.stayAwakeCheck.SetChecked(false)

	g.showTouchesCheck = widget.NewCheck("Show Physical Touches", nil)
	g.showTouchesCheck.SetChecked(false)

	g.powerOffOnCloseCheck = widget.NewCheck("Power Off Screen on Close", nil)
	g.powerOffOnCloseCheck.SetChecked(false)

	g.noPowerOnCheck = widget.NewCheck("Do Not Power On on Start", nil)
	g.noPowerOnCheck.SetChecked(false)

	g.keepActiveCheck = widget.NewCheck("Keep Active (scrcpy >= 3.4)", nil)
	g.keepActiveCheck.SetChecked(false)
	if !g.caps.HasKeepActive {
		g.keepActiveCheck.Hide()
	}

	g.screenOffTimeoutEntry = widget.NewEntry()
	g.screenOffTimeoutEntry.SetPlaceHolder("seconds (0 = keep current)")

	g.startAppEntry = widget.NewEntry()
	g.startAppEntry.SetPlaceHolder("e.g. org.mozilla.firefox or +pkg")
}

func (g *GUI) buildDeviceRows() []fyne.CanvasObject {
	timeoutRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Screen Off Timeout:"),
		nil,
		g.screenOffTimeoutEntry,
	)

	startAppRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Auto-Start App:"),
		nil,
		g.startAppEntry,
	)

	return []fyne.CanvasObject{
		g.screenOnCheck,
		g.stayAwakeCheck,
		g.showTouchesCheck,
		g.powerOffOnCloseCheck,
		g.noPowerOnCheck,
		g.keepActiveCheck,
		timeoutRow,
		startAppRow,
	}
}

