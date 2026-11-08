package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initVideoWidgets() {
	g.videoSourceCombo = widget.NewSelect([]string{"display (default)", "camera"}, nil)
	g.videoSourceCombo.SetSelected("display (default)")

	g.cameraFacingCombo = widget.NewSelect([]string{"front", "back", "external"}, nil)
	g.cameraFacingCombo.SetSelected("front")

	g.cameraSizeEntry = widget.NewEntry()
	g.cameraSizeEntry.SetText("1920x1080")
	g.cameraSizeEntry.SetPlaceHolder("e.g. 1920x1080")

	g.cameraFPSCombo = widget.NewSelect([]string{"Default", "60", "30"}, nil)
	g.cameraFPSCombo.SetSelected("Default")

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

	g.v4l2SinkEntry = widget.NewEntry()
	g.v4l2SinkEntry.SetPlaceHolder("e.g. /dev/video0")

	g.v4l2BufferEntry = widget.NewEntry()
	g.v4l2BufferEntry.SetPlaceHolder("e.g. 50ms")

	g.noPlaybackCheck = widget.NewCheck("No Video Playback (Record Only)", nil)
	g.noPlaybackCheck.SetChecked(false)
}

func (g *GUI) buildCameraRows() []fyne.CanvasObject {
	srcRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("video source:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.videoSourceCombo),
		nil,
	)
	facingRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("camera facing:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.cameraFacingCombo),
		nil,
	)
	sizeRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("camera size:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.cameraSizeEntry),
		nil,
	)
	return []fyne.CanvasObject{srcRow, facingRow, sizeRow}
}

func (g *GUI) buildV4L2Rows() []fyne.CanvasObject {
	sinkRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("v4l2 sink:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.v4l2SinkEntry),
		nil,
	)
	bufRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("v4l2 buffer:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.v4l2BufferEntry),
		nil,
	)
	return []fyne.CanvasObject{sinkRow, bufRow}
}

func (g *GUI) buildVideoSettingsRows() []fyne.CanvasObject {
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
		container.NewGridWrap(fyne.NewSize(180, 36), g.videoBufferEntry),
		nil,
	)
	cropRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("screen crop:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.cropEntry),
		nil,
	)
	return []fyne.CanvasObject{codecRow, fpsRow, orientRow, bufRow, cropRow}
}

func (g *GUI) buildVideoEngineTab() fyne.CanvasObject {
	items := make([]fyne.CanvasObject, 0, 12)
	items = append(items, g.buildCameraRows()...)
	items = append(items, g.buildVideoSettingsRows()...)
	items = append(items, g.buildV4L2Rows()...)
	items = append(items, g.noPlaybackCheck)

	return container.NewVScroll(container.NewVBox(items...))
}
