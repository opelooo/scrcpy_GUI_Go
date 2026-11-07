package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// clickableLabel provides primary selection, double tap to run, and right click context menu.
type clickableLabel struct {
	widget.Label
	serial       string
	onSelected   func(serial string)
	onRun        func()
	onDeviceInfo func(serial string)
	onDeviceIP   func(serial string)
	canvas       fyne.Canvas
}

func newClickableLabel(serial string, canvas fyne.Canvas, onSel, onRun, onInfo, onIP func(string)) *clickableLabel {
	l := &clickableLabel{
		serial:       serial,
		canvas:       canvas,
		onSelected:   onSel,
		onRun: func() {
			if onRun != nil {
				onRun(serial)
			}
		},
		onDeviceInfo: onInfo,
		onDeviceIP:   onIP,
	}
	l.Text = serial
	l.ExtendBaseWidget(l)
	return l
}

func (l *clickableLabel) Tapped(_ *fyne.PointEvent) {
	if l.onSelected != nil {
		l.onSelected(l.serial)
	}
}

func (l *clickableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if l.onRun != nil {
		l.onRun()
	}
}

func (l *clickableLabel) TappedSecondary(e *fyne.PointEvent) {
	if l.onSelected != nil {
		l.onSelected(l.serial)
	}
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Show Device Info", func() {
			if l.onDeviceInfo != nil {
				l.onDeviceInfo(l.serial)
			}
		}),
		fyne.NewMenuItem("Show Device IP", func() {
			if l.onDeviceIP != nil {
				l.onDeviceIP(l.serial)
			}
		}),
	)
	widget.NewPopUpMenu(menu, l.canvas).ShowAtPosition(e.AbsolutePosition)
}

type GUI struct {
	app    fyne.App
	window fyne.Window
	runner *ADBRunner

	devices        []string
	selectedDevice string

	// UI controls
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

func NewGUI(a fyne.App, w fyne.Window, runner *ADBRunner) *GUI {
	g := &GUI{
		app:    a,
		window: w,
		runner: runner,
	}
	g.buildUI()
	g.refreshDevices()
	return g
}

func (g *GUI) buildUI() {
	// Selected device label
	g.currDeviceLabel = widget.NewLabelWithStyle("No Device Selected", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Device List widget
	g.deviceList = widget.NewList(
		func() int {
			return len(g.devices)
		},
		func() fyne.CanvasObject {
			return newClickableLabel("", g.window.Canvas(), g.selectDeviceBySerial, func(_ string) { g.runScrcpy() }, g.showDeviceInfo, g.showDeviceIP)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if id < len(g.devices) {
				lbl := o.(*clickableLabel)
				lbl.serial = g.devices[id]
				lbl.SetText(g.devices[id])
			}
		},
	)
	g.deviceList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(g.devices) {
			g.selectDevice(g.devices[id])
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

	tcpipBtn := widget.NewButtonWithIcon("TCP/IP", theme.NavigateNextIcon(), func() {
		g.connectTCPIP()
	})

	buttonPanel := container.NewVBox(
		refreshBtn,
		infoBtn,
		runBtn,
		tcpipBtn,
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
			container.NewGridWrap(fyne.NewSize(220, 340), leftPanel),
			container.NewGridWrap(fyne.NewSize(110, 340), buttonPanel),
			container.NewGridWrap(fyne.NewSize(310, 340), optionsPanel),
		),
	)

	g.window.SetContent(content)
}

func (g *GUI) selectDeviceBySerial(serial string) {
	for i, d := range g.devices {
		if d == serial {
			g.deviceList.Select(i)
			return
		}
	}
	g.selectDevice(serial)
}

func (g *GUI) selectDevice(serial string) {
	g.selectedDevice = serial
	if serial == "" {
		g.currDeviceLabel.SetText("No Device Selected")
	} else {
		g.currDeviceLabel.SetText("Device: " + serial)
	}
}

func (g *GUI) refreshDevices() {
	devices, err := g.runner.ListDevices()
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}
	g.devices = devices
	g.deviceList.Refresh()

	if len(g.devices) > 0 {
		g.selectDevice(g.devices[0])
		g.deviceList.Select(0)
	} else {
		g.selectDevice("")
	}
}

func (g *GUI) showDeviceInfo(serial string) {
	if serial == "" {
		dialog.ShowInformation("Info", "Please select device", g.window)
		return
	}
	info, err := g.runner.GetDeviceInfo(serial)
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
	ip, err := g.runner.GetDeviceIP(serial)
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

	opts := ScrcpyOptions{
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
		err := g.runner.RunScrcpy(opts)
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
		err := g.runner.ConnectTCPIP(g.selectedDevice)
		time.Sleep(1500 * time.Millisecond)
		progress.Hide()
		if err != nil {
			dialog.ShowError(err, g.window)
		} else {
			// Refresh list after wireless connection
			time.Sleep(500 * time.Millisecond)
			g.refreshDevices()
		}
	}()
}
