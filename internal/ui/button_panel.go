package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

func (g *GUI) buildButtonPanel() *fyne.Container {
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

	return container.NewVBox(
		refreshBtn,
		infoBtn,
		runBtn,
		tcpipBtn,
		connectIPBtn,
		pairBtn,
	)
}

func (g *GUI) collectScrcpyOptions() domain.ScrcpyOptions {
	return domain.ScrcpyOptions{
		DeviceSerial:    g.selectedDevice,
		MaxSize:         g.maxSizeEntry.Text,
		BitRate:         g.bitRateEntry.Text,
		VideoOn:         g.videoOnCheck.Checked,
		ScreenOn:        g.screenOnCheck.Checked,
		StayAwake:       g.stayAwakeCheck.Checked,
		RecordVideo:         g.recordVideoCheck.Checked,
		VideoFormat:         g.videoFormatCombo.Selected,
		RecordAudio:         g.recordAudioCheck.Checked,
		AudioFormat:         g.audioFormatCombo.Selected,
		RecordPath:          g.recordPathEntry.Text,
		RecordOrientation:   g.recordOrientationCombo.Selected,
		TimeLimit:           g.timeLimitEntry.Text,
		MaxFPS:          g.maxFPSCombo.Selected,
		VideoCodec:      g.videoCodecCombo.Selected,
		Orientation:     g.orientationCombo.Selected,
		VideoBuffer:     g.videoBufferEntry.Text,
		Crop:            g.cropEntry.Text,
		V4L2Sink:        g.v4l2SinkEntry.Text,
		V4L2Buffer:      g.v4l2BufferEntry.Text,
		VideoSource:     g.videoSourceCombo.Selected,
		CameraFacing:    g.cameraFacingCombo.Selected,
		CameraSize:      g.cameraSizeEntry.Text,
		CameraFPS:       g.cameraFPSCombo.Selected,
		NoPlayback:      g.noPlaybackCheck.Checked,
		AudioOn:         g.audioOnCheck.Checked,
		AudioSource:     g.audioSourceCombo.Selected,
		AudioDup:        g.audioDupCheck.Checked,
		AudioCodec:      g.audioCodecCombo.Selected,
		AudioBitRate:    g.audioBitRateCombo.Selected,
		AudioBuffer:         g.audioBufferEntry.Text,
		NoAudioPlayback:     g.noAudioPlaybackCheck.Checked,
		NoControl:           g.noControlCheck.Checked,
		KeyboardMode:        g.keyboardModeCombo.Selected,
		MouseMode:           g.mouseModeCombo.Selected,
		GamepadMode:         g.gamepadModeCombo.Selected,
		NoMouseHover:        g.noMouseHoverCheck.Checked,
		NoClipboardAutosync: g.noClipboardAutosyncCheck.Checked,
		LegacyPaste:         g.legacyPasteCheck.Checked,
		PushTarget:          g.pushTargetEntry.Text,
		PreferText:          g.preferTextCheck.Checked,
		RawKeyEvents:        g.rawKeyEventsCheck.Checked,
		NoKeyRepeat:         g.noKeyRepeatCheck.Checked,
		ShortcutMod:         g.shortcutModCombo.Selected,
		OTGMode:             g.otgModeCheck.Checked,
		ShowTouches:         g.showTouchesCheck.Checked,
		PowerOffOnClose:     g.powerOffOnCloseCheck.Checked,
		NoPowerOn:           g.noPowerOnCheck.Checked,
		KeepActive:          g.keepActiveCheck.Visible() && g.keepActiveCheck.Checked,
		ScreenOffTimeout:        g.screenOffTimeoutEntry.Text,
		StartApp:                g.startAppEntry.Text,
		AlwaysOnTop:             g.alwaysOnTopCheck.Checked,
		Fullscreen:              g.fullscreenCheck.Checked,
		WindowBorderless:        g.windowBorderlessCheck.Checked,
		WindowTitle:             g.windowTitleEntry.Text,
		WindowX:                 g.windowXEntry.Text,
		WindowY:                 g.windowYEntry.Text,
		WindowWidth:             g.windowWidthEntry.Text,
		WindowHeight:            g.windowHeightEntry.Text,
		DisableScreensaver:      g.disableScreensaverCheck.Checked,
		NoWindow:                g.noWindowCheck.Checked,
		NoWindowAspectRatioLock: g.noWindowAspectRatioLockCheck.Visible() && g.noWindowAspectRatioLockCheck.Checked,
		BackgroundColor:         g.backgroundColorEntry.Text,
		RenderFit:               g.renderFitCombo.Selected,
	}
}

func (g *GUI) runScrcpy() {
	if g.selectedDevice == "" {
		dialog.ShowInformation("Info", "Please select device", g.window)
		return
	}

	opts := g.collectScrcpyOptions()
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
		fyne.Do(func() {
			progress.Hide()
			if err != nil {
				dialog.ShowError(err, g.window)
			}
		})
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
		fyne.Do(func() {
			progress.Hide()
			if err != nil {
				dialog.ShowError(err, g.window)
			} else {
				dialog.ShowInformation("Success", "Wireless scrcpy launched!", g.window)
				g.refreshDevices()
			}
		})
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
		fyne.Do(func() {
			progress.Hide()
			if err != nil {
				dialog.ShowError(err, g.window)
			} else {
				g.refreshDevices()
			}
		})
	}()
}
