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
	caps           domain.ScrcpyCapabilities

	// Core & List widgets
	deviceList      *widget.List
	currDeviceLabel *widget.Label

	// Basic & Record widgets
	bitRateEntry     *widget.Entry
	bitRateSlider    *widget.Slider
	maxSizeEntry     *widget.Entry
	videoOnCheck     *widget.Check
	recordVideoCheck *widget.Check
	videoFormatCombo *widget.Select
	recordAudioCheck *widget.Check
	audioFormatCombo *widget.Select

	// Device widgets (device.md)
	screenOnCheck         *widget.Check
	stayAwakeCheck        *widget.Check
	showTouchesCheck      *widget.Check
	powerOffOnCloseCheck  *widget.Check
	noPowerOnCheck        *widget.Check
	keepActiveCheck       *widget.Check
	screenOffTimeoutEntry *widget.Entry
	startAppEntry         *widget.Entry

	// Video Engine widgets
	videoSourceCombo  *widget.Select
	cameraFacingCombo *widget.Select
	cameraSizeEntry   *widget.Entry
	cameraFPSCombo    *widget.Select
	videoCodecCombo   *widget.Select
	maxFPSCombo       *widget.Select
	orientationCombo  *widget.Select
	videoBufferEntry  *widget.Entry
	cropEntry         *widget.Entry
	v4l2SinkEntry     *widget.Entry
	v4l2BufferEntry   *widget.Entry
	noPlaybackCheck   *widget.Check

	// Audio Engine widgets
	audioOnCheck         *widget.Check
	audioSourceCombo     *widget.Select
	audioDupCheck        *widget.Check
	audioCodecCombo      *widget.Select
	audioBitRateCombo    *widget.Select
	audioBufferEntry     *widget.Entry
	noAudioPlaybackCheck *widget.Check

	// Control Engine widgets
	noControlCheck           *widget.Check
	keyboardModeCombo        *widget.Select
	mouseModeCombo           *widget.Select
	gamepadModeCombo         *widget.Select
	noMouseHoverCheck        *widget.Check
	noClipboardAutosyncCheck *widget.Check
	legacyPasteCheck         *widget.Check
	pushTargetEntry          *widget.Entry
	preferTextCheck          *widget.Check
	rawKeyEventsCheck        *widget.Check
	noKeyRepeatCheck         *widget.Check
	shortcutModCombo         *widget.Select
	otgModeCheck             *widget.Check

	// Window Engine widgets (window.md)
	alwaysOnTopCheck             *widget.Check
	fullscreenCheck              *widget.Check
	windowBorderlessCheck        *widget.Check
	disableScreensaverCheck      *widget.Check
	noWindowCheck                *widget.Check
	windowTitleEntry             *widget.Entry
	windowXEntry                 *widget.Entry
	windowYEntry                 *widget.Entry
	windowWidthEntry             *widget.Entry
	windowHeightEntry            *widget.Entry
	noWindowAspectRatioLockCheck *widget.Check
	backgroundColorEntry         *widget.Entry
	renderFitCombo               *widget.Select
}

// NewGUI constructs the GUI and sets up the layout.
func NewGUI(a fyne.App, w fyne.Window, svc domain.DeviceService) *GUI {
	caps, _ := svc.GetScrcpyCapabilities()
	g := &GUI{
		app:         a,
		window:      w,
		service:     svc,
		caps:        caps,
		stopWatcher: make(chan struct{}),
	}
	g.buildUI()
	g.refreshDevices()
	go g.startDeviceWatcher()
	return g
}

func (g *GUI) initWidgets() {
	g.initBasicWidgets()
	g.initDeviceWidgets()
	g.initWindowWidgets()
	g.initVideoWidgets()
	g.initAudioWidgets()
	g.initControlWidgets()
}

func (g *GUI) buildOptionsPanel() *fyne.Container {
	tabs := container.NewAppTabs(
		container.NewTabItem("Basic & Rec", g.buildBasicTab()),
		container.NewTabItem("Device", g.buildDeviceTab()),
		container.NewTabItem("Window", g.buildWindowTab()),
		container.NewTabItem("Video", g.buildVideoEngineTab()),
		container.NewTabItem("Audio", g.buildAudioEngineTab()),
		container.NewTabItem("Control", g.buildControlTab()),
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
			container.NewGridWrap(fyne.NewSize(230, 390), leftPanel),
			container.NewGridWrap(fyne.NewSize(125, 390), buttonPanel),
			container.NewGridWrap(fyne.NewSize(480, 390), optionsPanel),
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
