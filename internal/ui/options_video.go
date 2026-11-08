package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initVideoWidgets() {
	g.videoCodecCombo = widget.NewSelect([]string{"default", "h264", "h265", "av1", "vp8", "vp9"}, nil)
	g.videoCodecCombo.SetSelected("default")

	g.maxFPSCombo = widget.NewSelect([]string{"Default", "60", "30", "15"}, nil)
	g.maxFPSCombo.SetSelected("Default")

	g.orientationCombo = widget.NewSelect([]string{"Default", "0", "90", "180", "270"}, nil)
	g.orientationCombo.SetSelected("Default")

	g.videoBufferEntry = widget.NewEntry()
	g.videoBufferEntry.SetPlaceHolder("e.g. 50ms")

	g.cropEntry = widget.NewEntry()
	g.cropEntry.SetPlaceHolder("width:height:x:y")

	g.noPlaybackCheck = widget.NewCheck("No Video Playback (Record Only)", nil)
	g.noPlaybackCheck.SetChecked(false)
}

func (g *GUI) buildVideoEngineTab() *fyne.Container {
	codecRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("video codec:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.videoCodecCombo),
		nil,
	)
	fpsRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("max fps:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.maxFPSCombo),
		nil,
	)
	orientRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("orientation:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.orientationCombo),
		nil,
	)
	bufRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("video buffer:"),
		container.NewGridWrap(fyne.NewSize(150, 36), g.videoBufferEntry),
		nil,
	)
	cropRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("screen crop:"),
		container.NewGridWrap(fyne.NewSize(150, 36), g.cropEntry),
		nil,
	)

	return container.NewVBox(
		codecRow,
		fpsRow,
		orientRow,
		bufRow,
		cropRow,
		g.noPlaybackCheck,
	)
}

