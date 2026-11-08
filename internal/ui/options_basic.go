package ui

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initBasicWidgets() {
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

	g.maxSizeEntry = widget.NewEntry()
	g.maxSizeEntry.SetText("1024")

	g.videoOnCheck = widget.NewCheck("Video On", nil)
	g.videoOnCheck.SetChecked(true)

	g.screenOnCheck = widget.NewCheck("Turn Screen On", nil)
	g.screenOnCheck.SetChecked(true)

	g.stayAwakeCheck = widget.NewCheck("Stay Awake", nil)
	g.stayAwakeCheck.SetChecked(false)

	g.recordVideoCheck = widget.NewCheck("Record Video", nil)
	g.videoFormatCombo = widget.NewSelect([]string{"mp4", "mkv"}, nil)
	g.videoFormatCombo.SetSelected("mp4")

	g.recordAudioCheck = widget.NewCheck("Record Audio", nil)
	g.audioFormatCombo = widget.NewSelect([]string{"m4a", "wav", "flac", "opus", "aac", "mka"}, nil)
	g.audioFormatCombo.SetSelected("m4a")
}

func (g *GUI) buildBasicTab() *fyne.Container {
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

	return container.NewVBox(
		bitrateRow,
		maxSizeRow,
		g.videoOnCheck,
		g.screenOnCheck,
		g.stayAwakeCheck,
		recVideoRow,
		recAudioRow,
	)
}

