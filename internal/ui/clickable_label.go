package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

// clickableLabel provides primary selection, double tap to run, and right click context menu.
type clickableLabel struct {
	widget.Label
	device       domain.Device
	onSelected   func(serial string)
	onRun        func()
	onDeviceInfo func(serial string)
	onDeviceIP   func(serial string)
	onDisconnect func(serial string)
	canvas       fyne.Canvas
}

func newClickableLabel(canvas fyne.Canvas, onSel func(string), onRun func(), onInfo, onIP, onDisc func(string)) *clickableLabel {
	l := &clickableLabel{
		canvas:       canvas,
		onSelected:   onSel,
		onRun:        onRun,
		onDeviceInfo: onInfo,
		onDeviceIP:   onIP,
		onDisconnect: onDisc,
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *clickableLabel) updateDevice(dev domain.Device) {
	l.device = dev
	prefix := "[USB] "
	if dev.ConnectionType == domain.ConnTCPIP {
		prefix = "[Wi-Fi] "
	}
	l.SetText(prefix + dev.Serial)
}

func (l *clickableLabel) Tapped(_ *fyne.PointEvent) {
	if l.onSelected != nil {
		l.onSelected(l.device.Serial)
	}
}

func (l *clickableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if l.onRun != nil {
		l.onRun()
	}
}

func (l *clickableLabel) TappedSecondary(e *fyne.PointEvent) {
	if l.onSelected != nil {
		l.onSelected(l.device.Serial)
	}
	items := []*fyne.MenuItem{
		fyne.NewMenuItem("Show Device Info", func() {
			if l.onDeviceInfo != nil {
				l.onDeviceInfo(l.device.Serial)
			}
		}),
		fyne.NewMenuItem("Show Device IP", func() {
			if l.onDeviceIP != nil {
				l.onDeviceIP(l.device.Serial)
			}
		}),
	}

	if l.device.ConnectionType == domain.ConnTCPIP {
		items = append(items, fyne.NewMenuItem("Disconnect Wireless", func() {
			if l.onDisconnect != nil {
				l.onDisconnect(l.device.Serial)
			}
		}))
	}

	menu := fyne.NewMenu("", items...)
	widget.NewPopUpMenu(menu, l.canvas).ShowAtPosition(e.AbsolutePosition)
}

