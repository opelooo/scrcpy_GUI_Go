package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initWindowWidgets() {
	g.alwaysOnTopCheck = widget.NewCheck("Always On Top", nil)
	g.alwaysOnTopCheck.SetChecked(false)

	g.fullscreenCheck = widget.NewCheck("Fullscreen on Start", nil)
	g.fullscreenCheck.SetChecked(false)

	g.windowBorderlessCheck = widget.NewCheck("Borderless Window", nil)
	g.windowBorderlessCheck.SetChecked(false)

	g.disableScreensaverCheck = widget.NewCheck("Disable Desktop Screensaver", nil)
	g.disableScreensaverCheck.SetChecked(false)

	g.noWindowCheck = widget.NewCheck("Headless Mode (No Window)", func(checked bool) {
		if checked {
			g.alwaysOnTopCheck.SetChecked(false)
			g.fullscreenCheck.SetChecked(false)
			g.windowBorderlessCheck.SetChecked(false)
		}
	})
	g.noWindowCheck.SetChecked(false)

	g.windowTitleEntry = widget.NewEntry()
	g.windowTitleEntry.SetPlaceHolder("Default: device model")

	g.windowXEntry = widget.NewEntry()
	g.windowXEntry.SetPlaceHolder("X coord")
	g.windowYEntry = widget.NewEntry()
	g.windowYEntry.SetPlaceHolder("Y coord")

	g.windowWidthEntry = widget.NewEntry()
	g.windowWidthEntry.SetPlaceHolder("Width")
	g.windowHeightEntry = widget.NewEntry()
	g.windowHeightEntry.SetPlaceHolder("Height")

	g.initWindowAdvancedWidgets()
}

func (g *GUI) initWindowAdvancedWidgets() {
	g.noWindowAspectRatioLockCheck = widget.NewCheck("Unlock Aspect Ratio (scrcpy >= 3.4)", nil)
	g.noWindowAspectRatioLockCheck.SetChecked(false)

	g.backgroundColorEntry = widget.NewEntry()
	g.backgroundColorEntry.SetPlaceHolder("Padding HEX color (e.g. #000000)")

	g.renderFitCombo = widget.NewSelect([]string{"default", "letterbox", "unscaled", "stretched"}, nil)
	g.renderFitCombo.SetSelected("default")

	if !g.caps.HasWindowAdvanced {
		g.noWindowAspectRatioLockCheck.Hide()
		g.backgroundColorEntry.Hide()
		g.renderFitCombo.Hide()
	}
}

func (g *GUI) buildWindowRows() []fyne.CanvasObject {
	titleRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Window Title:"),
		nil,
		g.windowTitleEntry,
	)

	posRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Initial Pos (X, Y):"),
		nil,
		container.NewGridWithColumns(2, g.windowXEntry, g.windowYEntry),
	)

	sizeRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Initial Size (W, H):"),
		nil,
		container.NewGridWithColumns(2, g.windowWidthEntry, g.windowHeightEntry),
	)

	return []fyne.CanvasObject{
		g.alwaysOnTopCheck,
		g.fullscreenCheck,
		g.windowBorderlessCheck,
		g.disableScreensaverCheck,
		g.noWindowCheck,
		titleRow,
		posRow,
		sizeRow,
		g.noWindowAspectRatioLockCheck,
		g.backgroundColorEntry,
		g.renderFitCombo,
	}
}
