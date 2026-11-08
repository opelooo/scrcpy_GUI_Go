package ui

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
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

// GUI manages the graphical user interface.
type GUI struct {
	app     fyne.App
	window  fyne.Window
	service domain.DeviceService

	mu             sync.Mutex
	devices        []domain.Device
	selectedDevice string
	stopWatcher    chan struct{}

	// UI widgets
	deviceList       *widget.List
	currDeviceLabel  *widget.Label
	bitRateEntry     *widget.Entry
	bitRateSlider    *widget.Slider
	maxSizeEntry     *widget.Entry
	videoOnCheck     *widget.Check
	screenOnCheck    *widget.Check
	stayAwakeCheck   *widget.Check
	recordVideoCheck *widget.Check
	videoFormatCombo *widget.Select
	recordAudioCheck *widget.Check
	audioFormatCombo *widget.Select
}

// NewGUI constructs the GUI and sets up the layout.
func NewGUI(a fyne.App, w fyne.Window, svc domain.DeviceService) *GUI {
	g := &GUI{
		app:         a,
		window:      w,
		service:     svc,
		stopWatcher: make(chan struct{}),
	}
	g.buildUI()
	g.refreshDevices()
	go g.startDeviceWatcher()
	return g
}

func (g *GUI) buildUI() {
	g.currDeviceLabel = widget.NewLabelWithStyle("No Device Selected", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Device List widget
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

	// Bit rate controls
	g.bitRateEntry = widget.NewEntry()
	g.bitRateEntry.SetText("2M")
	g.bitRateEntry.OnChanged = func(s string) {
		clean := strings.TrimSuffix(strings.TrimSpace(strings.ToUpper(s)), "M")
		if val, err := strconv.Atoi(clean); err == nil && val >= 1 && val <= 50 {
			if g.bitRateSlider.Value != float64(val) {
				g.bitRateSlider.SetValue(float64(val))
			}
		}
	}

	g.bitRateSlider = widget.NewSlider(1, 50)
	g.bitRateSlider.Step = 1
	g.bitRateSlider.SetValue(2)
	g.bitRateSlider.OnChanged = func(val float64) {
		text := fmt.Sprintf("%dM", int(val))
		if g.bitRateEntry.Text != text {
			g.bitRateEntry.SetText(text)
		}
	}

	// Max size control
	g.maxSizeEntry = widget.NewEntry()
	g.maxSizeEntry.SetText("1024")

	// Toggles
	g.videoOnCheck = widget.NewCheck("Video On", nil)
	g.videoOnCheck.SetChecked(true)

	g.screenOnCheck = widget.NewCheck("Turn Screen On", nil)
	g.screenOnCheck.SetChecked(true)

	g.stayAwakeCheck = widget.NewCheck("Stay Awake", nil)
	g.stayAwakeCheck.SetChecked(false)

	// Recording options
	g.recordVideoCheck = widget.NewCheck("Record Video", nil)
	g.videoFormatCombo = widget.NewSelect([]string{"mp4", "mkv"}, nil)
	g.videoFormatCombo.SetSelected("mp4")

	g.recordAudioCheck = widget.NewCheck("Record Audio", nil)
	g.audioFormatCombo = widget.NewSelect([]string{"m4a", "wav", "flac", "opus", "aac", "mka"}, nil)
	g.audioFormatCombo.SetSelected("m4a")

	// Action Buttons
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		g.refreshDevices()
	})
	infoBtn := widget.NewButtonWithIcon("Info", theme.InfoIcon(), func() {
		g.showDeviceInfo(g.selectedDevice)
	})
	runBtn := widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() {
		g.runScrcpy()
	})
	runBtn.Importance = widget.HighImportance

	tcpipBtn := widget.NewButtonWithIcon("USB→WiFi", theme.NavigateNextIcon(), func() {
		g.connectTCPIP()
	})

	connectIPBtn := widget.NewButtonWithIcon("Connect IP", theme.LoginIcon(), func() {
		g.showConnectIPDialog()
	})

	pairBtn := widget.NewButtonWithIcon("Pair (A11+)", theme.SettingsIcon(), func() {
		g.showPairDialog()
	})

	buttonPanel := container.NewVBox(
		refreshBtn,
		infoBtn,
		runBtn,
		tcpipBtn,
		connectIPBtn,
		pairBtn,
	)

	// Options Form Layout
	bitrateRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("video bit rate:"),
		container.NewGridWrap(fyne.NewSize(70, 36), g.bitRateEntry),
		g.bitRateSlider,
	)
	maxSizeRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("max mirroring size:"),
		container.NewGridWrap(fyne.NewSize(70, 36), g.maxSizeEntry),
		nil,
	)

	recVideoRow := container.NewBorder(
		nil, nil,
		g.recordVideoCheck,
		container.NewGridWrap(fyne.NewSize(85, 36), g.videoFormatCombo),
		nil,
	)
	recAudioRow := container.NewBorder(
		nil, nil,
		g.recordAudioCheck,
		container.NewGridWrap(fyne.NewSize(85, 36), g.audioFormatCombo),
		nil,
	)

	optionsPanel := container.NewVBox(
		g.currDeviceLabel,
		widget.NewSeparator(),
		bitrateRow,
		maxSizeRow,
		g.videoOnCheck,
		g.screenOnCheck,
		g.stayAwakeCheck,
		recVideoRow,
		recAudioRow,
	)

	// Left list panel with header
	listHeader := widget.NewLabelWithStyle("Attached Devices", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	leftPanel := container.NewBorder(listHeader, nil, nil, nil, g.deviceList)

	// Main 3-column layout
	content := container.NewBorder(
		nil, nil, nil, nil,
		container.NewHBox(
			container.NewGridWrap(fyne.NewSize(230, 360), leftPanel),
			container.NewGridWrap(fyne.NewSize(125, 360), buttonPanel),
			container.NewGridWrap(fyne.NewSize(315, 360), optionsPanel),
		),
	)

	g.window.SetContent(content)
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

func (g *GUI) Close() {
	select {
	case <-g.stopWatcher:
	default:
		close(g.stopWatcher)
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

func (g *GUI) runScrcpy() {
	if g.selectedDevice == "" {
		dialog.ShowInformation("Info", "Please select device", g.window)
		return
	}

	opts := domain.ScrcpyOptions{
		DeviceSerial: g.selectedDevice,
		MaxSize:      g.maxSizeEntry.Text,
		BitRate:      g.bitRateEntry.Text,
		VideoOn:      g.videoOnCheck.Checked,
		ScreenOn:     g.screenOnCheck.Checked,
		StayAwake:    g.stayAwakeCheck.Checked,
		RecordVideo:  g.recordVideoCheck.Checked,
		VideoFormat:  g.videoFormatCombo.Selected,
		RecordAudio:  g.recordAudioCheck.Checked,
		AudioFormat:  g.audioFormatCombo.Selected,
	}

	progress := dialog.NewCustomWithoutButtons(
		"Running scrcpy",
		container.NewVBox(
			widget.NewLabel("Launching scrcpy for "+g.selectedDevice+"..."),
			widget.NewProgressBarInfinite(),
		),
		g.window,
	)
	progress.Show()

	go func() {
		err := g.service.LaunchScrcpy(opts)
		time.Sleep(1200 * time.Millisecond)
		progress.Hide()
		if err != nil {
			dialog.ShowError(err, g.window)
		}
	}()
}

func (g *GUI) connectTCPIP() {
	if g.selectedDevice == "" {
		dialog.ShowInformation("Info", "Please select device", g.window)
		return
	}

	progress := dialog.NewCustomWithoutButtons(
		"TCP/IP Connection",
		container.NewVBox(
			widget.NewLabel("Connecting to "+g.selectedDevice+" over TCP/IP..."),
			widget.NewProgressBarInfinite(),
		),
		g.window,
	)
	progress.Show()

	go func() {
		err := g.service.ConnectTCPIP(g.selectedDevice)
		time.Sleep(1500 * time.Millisecond)
		progress.Hide()
		if err != nil {
			dialog.ShowError(err, g.window)
		} else {
			time.Sleep(500 * time.Millisecond)
			g.refreshDevices()
		}
	}()
}

func (g *GUI) disconnectDevice(target string) {
	progress := dialog.NewCustomWithoutButtons(
		"Disconnecting",
		container.NewVBox(
			widget.NewLabel("Disconnecting "+target+"..."),
			widget.NewProgressBarInfinite(),
		),
		g.window,
	)
	progress.Show()

	go func() {
		err := g.service.DisconnectDevice(target)
		time.Sleep(500 * time.Millisecond)
		progress.Hide()
		if err != nil {
			dialog.ShowError(err, g.window)
		} else {
			g.refreshDevices()
		}
	}()
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
			progress.Hide()
			if err != nil {
				dialog.ShowError(err, g.window)
			} else {
				time.Sleep(500 * time.Millisecond)
				g.refreshDevices()
			}
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
			progress.Hide()
			if err != nil {
				dialog.ShowError(err, g.window)
			} else {
				dialog.ShowInformation("Pairing Successful", "Device paired successfully! Now connect via port 5555.", g.window)
				g.refreshDevices()
			}
		}()
	}, g.window)

	d.Resize(fyne.NewSize(380, 240))
	d.Show()
}
