package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

func (g *GUI) showDeviceInfo(serial string) {
	if serial == "" {
		dialog.ShowInformation("Info", "Please select device", g.window)
		return
	}
	info, err := g.service.GetDeviceInfo(serial)
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}
	msg := fmt.Sprintf("Device Info for: %s\nDevice Manufacturer: %s", serial, info)
	dialog.ShowInformation("Device Info", msg, g.window)
}

func (g *GUI) showDeviceIP(serial string) {
	if serial == "" {
		dialog.ShowInformation("Info", "Please select device", g.window)
		return
	}
	ip, err := g.service.GetDeviceIP(serial)
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}
	msg := fmt.Sprintf("Device Info for: %s\nDevice IP Address: %s", serial, ip)
	dialog.ShowInformation("Device Info", msg, g.window)
}

func (g *GUI) showConnectIPDialog() {
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("192.168.1.50")

	portEntry := widget.NewEntry()
	portEntry.SetText("5555")

	forceCheck := widget.NewCheck("Force Reconnect (+)", nil)

	formItems := []*widget.FormItem{
		widget.NewFormItem("IP Address", ipEntry),
		widget.NewFormItem("Port", portEntry),
		widget.NewFormItem("", forceCheck),
	}

	d := dialog.NewForm("Connect via TCP/IP", "Connect", "Cancel", formItems, func(ok bool) {
		if !ok {
			return
		}
		ip := strings.TrimSpace(ipEntry.Text)
		if ip == "" {
			dialog.ShowError(fmt.Errorf("IP address cannot be empty"), g.window)
			return
		}
		port := 5555
		if pStr := strings.TrimSpace(portEntry.Text); pStr != "" {
			if p, err := strconv.Atoi(pStr); err == nil && p > 0 && p <= 65535 {
				port = p
			}
		}

		progress := dialog.NewCustomWithoutButtons(
			"Connecting Wireless",
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Connecting to %s:%d...", ip, port)),
				widget.NewProgressBarInfinite(),
			),
			g.window,
		)
		progress.Show()

		go func() {
			err := g.service.DirectConnect(domain.ConnectOptions{
				IP:             ip,
				Port:           port,
				ForceReconnect: forceCheck.Checked,
			})
			time.Sleep(1200 * time.Millisecond)
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(err, g.window)
				} else {
					time.Sleep(500 * time.Millisecond)
					g.refreshDevices()
				}
			})
		}()
	}, g.window)

	d.Resize(fyne.NewSize(380, 240))
	d.Show()
}

func (g *GUI) showPairDialog() {
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("192.168.1.50")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("e.g. 37123")

	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("6-digit code")

	formItems := []*widget.FormItem{
		widget.NewFormItem("IP Address", ipEntry),
		widget.NewFormItem("Pairing Port", portEntry),
		widget.NewFormItem("Pairing Code", codeEntry),
	}

	d := dialog.NewForm("Android 11+ Wireless Pairing", "Pair", "Cancel", formItems, func(ok bool) {
		if !ok {
			return
		}
		ip := strings.TrimSpace(ipEntry.Text)
		code := strings.TrimSpace(codeEntry.Text)
		port, err := strconv.Atoi(strings.TrimSpace(portEntry.Text))
		if ip == "" || code == "" || err != nil || port <= 0 {
			dialog.ShowError(fmt.Errorf("please provide valid IP, pairing port, and 6-digit code"), g.window)
			return
		}

		progress := dialog.NewCustomWithoutButtons(
			"Pairing Device",
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Pairing with %s:%d...", ip, port)),
				widget.NewProgressBarInfinite(),
			),
			g.window,
		)
		progress.Show()

		go func() {
			err := g.service.PairDevice(domain.PairOptions{
				IP:   ip,
				Port: port,
				Code: code,
			})
			time.Sleep(1200 * time.Millisecond)
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(err, g.window)
				} else {
					dialog.ShowInformation("Pairing Successful", "Device paired successfully! Now connect via port 5555.", g.window)
					g.refreshDevices()
				}
			})
		}()
	}, g.window)

	d.Resize(fyne.NewSize(380, 240))
	d.Show()
}

