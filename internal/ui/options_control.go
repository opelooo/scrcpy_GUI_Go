package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initControlWidgets() {
	g.noControlCheck = widget.NewCheck("Read-Only (Disable Control)", nil)
	g.noControlCheck.SetChecked(false)

	g.keyboardModeCombo = widget.NewSelect([]string{"default", "uhid", "aoa", "disabled"}, nil)
	g.keyboardModeCombo.SetSelected("default")

	g.mouseModeCombo = widget.NewSelect([]string{"default", "uhid", "aoa", "disabled"}, nil)
	g.mouseModeCombo.SetSelected("default")

	g.noMouseHoverCheck = widget.NewCheck("Disable Mouse Hover", nil)
	g.noMouseHoverCheck.SetChecked(false)

	g.noClipboardAutosyncCheck = widget.NewCheck("Disable Clipboard Autosync", nil)
	g.noClipboardAutosyncCheck.SetChecked(false)

	g.legacyPasteCheck = widget.NewCheck("Legacy Paste (Inject Keys)", nil)
	g.legacyPasteCheck.SetChecked(false)

	g.pushTargetEntry = widget.NewEntry()
	g.pushTargetEntry.SetText("/sdcard/Download/")
}

func (g *GUI) buildControlTab() *fyne.Container {
	kbdRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("keyboard mode:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.keyboardModeCombo),
		nil,
	)
	mouseRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("mouse mode:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.mouseModeCombo),
		nil,
	)
	pushRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("push target:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.pushTargetEntry),
		nil,
	)

	return container.NewVBox(
		g.noControlCheck,
		kbdRow,
		mouseRow,
		pushRow,
		g.noMouseHoverCheck,
		g.noClipboardAutosyncCheck,
		g.legacyPasteCheck,
	)
}

