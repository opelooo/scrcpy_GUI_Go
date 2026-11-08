package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

func (g *GUI) initDeviceList() {
	g.deviceList = widget.NewList(
		func() int {
			g.mu.Lock()
			defer g.mu.Unlock()
			return len(g.devices)
		},
		func() fyne.CanvasObject {
			return newClickableLabel(
				g.window.Canvas(),
				g.selectDeviceBySerial,
				func() { g.runScrcpy() },
				g.showDeviceInfo,
				g.showDeviceIP,
				g.disconnectDevice,
			)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			g.mu.Lock()
			defer g.mu.Unlock()
			if id < len(g.devices) {
				lbl := o.(*clickableLabel)
				lbl.updateDevice(g.devices[id])
			}
		},
	)
	g.deviceList.OnSelected = func(id widget.ListItemID) {
		g.mu.Lock()
		var serial string
		if id >= 0 && id < len(g.devices) {
			serial = g.devices[id].Serial
		}
		g.mu.Unlock()

		if serial != "" {
			g.selectDevice(serial)
		}
	}
}

func (g *GUI) selectDeviceBySerial(serial string) {
	g.mu.Lock()
	idx := -1
	for i, d := range g.devices {
		if d.Serial == serial {
			idx = i
			break
		}
	}
	g.mu.Unlock()

	if idx >= 0 {
		g.deviceList.Select(idx)
	} else {
		g.selectDevice(serial)
	}
}

func (g *GUI) selectDevice(serial string) {
	g.mu.Lock()
	g.selectedDevice = serial
	g.mu.Unlock()

	if serial == "" {
		g.currDeviceLabel.SetText("No Device Selected")
	} else {
		g.currDeviceLabel.SetText("Device: " + serial)
	}
}

func (g *GUI) refreshDevices() {
	devices, err := g.service.ListDevices()
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}
	g.setDevices(devices)
}

func (g *GUI) setDevices(devices []domain.Device) {
	g.mu.Lock()
	g.devices = devices
	count := len(devices)
	var firstSerial string
	if count > 0 {
		firstSerial = devices[0].Serial
	}
	found := false
	for _, d := range devices {
		if d.Serial == g.selectedDevice {
			found = true
			break
		}
	}
	g.mu.Unlock()

	g.deviceList.Refresh()

	if count > 0 {
		if !found {
			g.selectDevice(firstSerial)
			g.deviceList.Select(0)
		}
	} else {
		g.selectDevice("")
	}
}

func (g *GUI) startDeviceWatcher() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-g.stopWatcher:
			return
		case <-ticker.C:
			devs, err := g.service.ListDevices()
			if err != nil {
				continue
			}

			g.mu.Lock()
			changed := len(devs) != len(g.devices)
			if !changed {
				for i := range devs {
					if devs[i].Serial != g.devices[i].Serial || devs[i].State != g.devices[i].State {
						changed = true
						break
					}
				}
			}
			g.mu.Unlock()

			if changed {
				g.setDevices(devs)
			}
		}
	}
}

