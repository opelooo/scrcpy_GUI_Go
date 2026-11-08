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
		RecordVideo:     g.recordVideoCheck.Checked,
		VideoFormat:     g.videoFormatCombo.Selected,
		RecordAudio:     g.recordAudioCheck.Checked,
		AudioFormat:     g.audioFormatCombo.Selected,
		MaxFPS:          g.maxFPSCombo.Selected,
		VideoCodec:      g.videoCodecCombo.Selected,
		Orientation:     g.orientationCombo.Selected,
		VideoBuffer:     g.videoBufferEntry.Text,
		Crop:            g.cropEntry.Text,
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
		NoMouseHover:        g.noMouseHoverCheck.Checked,
		NoClipboardAutosync: g.noClipboardAutosyncCheck.Checked,
		LegacyPaste:         g.legacyPasteCheck.Checked,
		PushTarget:          g.pushTargetEntry.Text,
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
			dialog.ShowInformation("Success", "Wireless scrcpy launched!", g.window)
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
