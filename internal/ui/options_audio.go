package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initAudioWidgets() {
	g.audioOnCheck = widget.NewCheck("Audio Forwarding On", nil)
	g.audioOnCheck.SetChecked(true)

	g.audioSourceCombo = widget.NewSelect([]string{"output", "mic", "playback"}, nil)
	g.audioSourceCombo.SetSelected("output")

	g.audioDupCheck = widget.NewCheck("Duplicate Audio to Device (A13+)", nil)
	g.audioDupCheck.SetChecked(false)

	g.audioCodecCombo = widget.NewSelect([]string{"opus", "aac", "flac", "raw"}, nil)
	g.audioCodecCombo.SetSelected("opus")

	g.audioBitRateCombo = widget.NewSelect([]string{"Default", "64K", "128K", "192K", "256K"}, nil)
	g.audioBitRateCombo.SetSelected("Default")

	g.audioBufferEntry = widget.NewEntry()
	g.audioBufferEntry.SetPlaceHolder("e.g. 50ms")

	g.noAudioPlaybackCheck = widget.NewCheck("Mute PC Playback", nil)
	g.noAudioPlaybackCheck.SetChecked(false)
}

func (g *GUI) buildAudioEngineTab() *fyne.Container {
	sourceRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("audio source:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.audioSourceCombo),
		nil,
	)
	codecRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("audio codec:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.audioCodecCombo),
		nil,
	)
	bitrateRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("audio bit rate:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.audioBitRateCombo),
		nil,
	)
	bufRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("audio buffer:"),
		container.NewGridWrap(fyne.NewSize(150, 36), g.audioBufferEntry),
		nil,
	)

	return container.NewVBox(
		g.audioOnCheck,
		sourceRow,
		codecRow,
		bitrateRow,
		bufRow,
		g.audioDupCheck,
		g.noAudioPlaybackCheck,
	)
}

