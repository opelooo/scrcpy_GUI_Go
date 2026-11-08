package ui

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initCameraWidgets() {
	g.videoSourceCombo = widget.NewSelect([]string{"display (default)", "camera"}, nil)
	g.videoSourceCombo.SetSelected("display (default)")

	g.cameraIDEntry = widget.NewEntry()
	g.cameraIDEntry.SetPlaceHolder("e.g. 0 or 1 (overrides facing)")

	g.cameraFacingCombo = widget.NewSelect([]string{"front", "back", "external"}, nil)
	g.cameraFacingCombo.SetSelected("front")

	g.cameraSizeEntry = widget.NewEntry()
	g.cameraSizeEntry.SetText("1920x1080")
	g.cameraSizeEntry.SetPlaceHolder("e.g. 1920x1080")

	g.cameraARCombo = widget.NewSelect([]string{"Default", "sensor", "16:9", "4:3", "1.6"}, nil)
	g.cameraARCombo.SetSelected("Default")

	g.cameraFPSCombo = widget.NewSelect([]string{"Default", "60", "30", "120"}, nil)
	g.cameraFPSCombo.SetSelected("Default")

	g.cameraHighSpeedCheck = widget.NewCheck("High-Speed Capture (120/240 fps)", nil)
	g.cameraHighSpeedCheck.SetChecked(false)
}

func (g *GUI) initVideoCodecWidgets() {
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

func (g *GUI) initVideoWidgets() {
	g.initCameraWidgets()
	g.initVideoCodecWidgets()
}

func (g *GUI) buildCameraRows() []fyne.CanvasObject {
	srcRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("video source:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.videoSourceCombo),
		nil,
	)
	idRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("camera id:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.cameraIDEntry),
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
	arRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("camera ar:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.cameraARCombo),
		nil,
	)
	fpsRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("camera fps:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.cameraFPSCombo),
		nil,
	)
	return []fyne.CanvasObject{srcRow, idRow, facingRow, sizeRow, arRow, fpsRow, g.cameraHighSpeedCheck}
}

func (g *GUI) buildV4L2Rows() []fyne.CanvasObject {
	if runtime.GOOS != "linux" {
		notice := widget.NewLabelWithStyle(
			"V4L2 Virtual Webcam is only supported on Linux",
			fyne.TextAlignLeading,
			fyne.TextStyle{Italic: true},
		)
		return []fyne.CanvasObject{notice}
	}

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

func (g *GUI) buildMediaTab() fyne.CanvasObject {
	videoHeader := widget.NewLabelWithStyle("Video Engine & Camera", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	audioHeader := widget.NewLabelWithStyle("Audio Engine", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	recordHeader := widget.NewLabelWithStyle("Recording", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	items := make([]fyne.CanvasObject, 0, 30)
	items = append(items, videoHeader)
	items = append(items, g.buildBasicVideoRows()...)
	items = append(items, g.buildCameraRows()...)
	items = append(items, g.buildVideoSettingsRows()...)
	items = append(items, g.buildV4L2Rows()...)
	items = append(items, g.noPlaybackCheck)

	items = append(items, widget.NewSeparator(), audioHeader)
	items = append(items, g.buildAudioEngineRows()...)

	items = append(items, widget.NewSeparator(), recordHeader)
	items = append(items, g.buildRecordingRows()...)

	return container.NewVScroll(container.NewVBox(items...))
}
