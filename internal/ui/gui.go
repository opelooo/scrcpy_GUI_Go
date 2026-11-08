package ui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

// GUI coordinates the graphical user interface.
type GUI struct {
	app     fyne.App
	window  fyne.Window
	service domain.DeviceService

	mu             sync.Mutex
	devices        []domain.Device
	selectedDevice string
	stopWatcher    chan struct{}

	// Core & List widgets
	deviceList      *widget.List
	currDeviceLabel *widget.Label

	// Basic & Record widgets
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

	// Video Engine widgets
	videoCodecCombo  *widget.Select
	maxFPSCombo      *widget.Select
	orientationCombo *widget.Select
	videoBufferEntry *widget.Entry
	cropEntry        *widget.Entry
	noPlaybackCheck  *widget.Check

	// Audio Engine widgets
	audioOnCheck         *widget.Check
	audioSourceCombo     *widget.Select
	audioDupCheck        *widget.Check
	audioCodecCombo      *widget.Select
	audioBitRateCombo    *widget.Select
	audioBufferEntry     *widget.Entry
	noAudioPlaybackCheck *widget.Check
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

func (g *GUI) initWidgets() {
	g.initBasicWidgets()
	g.initVideoWidgets()
	g.initAudioWidgets()
}

func (g *GUI) buildOptionsPanel() *fyne.Container {
	tabs := container.NewAppTabs(
		container.NewTabItem("Basic & Rec", g.buildBasicTab()),
		container.NewTabItem("Video", g.buildVideoEngineTab()),
		container.NewTabItem("Audio", g.buildAudioEngineTab()),
	)
	return container.NewBorder(
		container.NewVBox(g.currDeviceLabel, widget.NewSeparator()),
		nil, nil, nil,
		tabs,
	)
}

func (g *GUI) buildUI() {
	g.currDeviceLabel = widget.NewLabelWithStyle("No Device Selected", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	g.initDeviceList()
	g.initWidgets()

	buttonPanel := g.buildButtonPanel()
	optionsPanel := g.buildOptionsPanel()

	listHeader := widget.NewLabelWithStyle("Attached Devices", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	leftPanel := container.NewBorder(listHeader, nil, nil, nil, g.deviceList)

	content := container.NewBorder(
		nil, nil, nil, nil,
		container.NewHBox(
			container.NewGridWrap(fyne.NewSize(230, 360), leftPanel),
			container.NewGridWrap(fyne.NewSize(125, 360), buttonPanel),
			container.NewGridWrap(fyne.NewSize(420, 360), optionsPanel),
		),
	)

	g.window.SetContent(content)
}

// Close gracefully terminates background watchers.
func (g *GUI) Close() {
	select {
	case <-g.stopWatcher:
	default:
		close(g.stopWatcher)
	}
}
