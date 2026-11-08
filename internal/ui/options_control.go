package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initControlCombos() {
	g.keyboardModeCombo = widget.NewSelect([]string{"default", "uhid", "aoa", "disabled"}, nil)
	g.keyboardModeCombo.SetSelected("default")

	g.mouseModeCombo = widget.NewSelect([]string{"default", "uhid", "aoa", "disabled"}, nil)
	g.mouseModeCombo.SetSelected("default")

	g.gamepadModeCombo = widget.NewSelect([]string{"default", "uhid", "aoa", "disabled"}, nil)
	g.gamepadModeCombo.SetSelected("default")

	g.shortcutModCombo = widget.NewSelect([]string{
		"Default (Alt/Super)",
		"rctrl",
		"lctrl+lalt",
		"lalt",
		"ralt",
		"lsuper",
		"rsuper",
		"lctrl+lsuper",
	}, nil)
	g.shortcutModCombo.SetSelected("Default (Alt/Super)")

	g.pushTargetEntry = widget.NewEntry()
	g.pushTargetEntry.SetText("/sdcard/Download/")
}

func (g *GUI) initControlModes() {
	g.noControlCheck = widget.NewCheck("Read-Only (Disable Control)", func(checked bool) {
		if !checked {
			return
		}
		if g.otgModeCheck != nil && g.otgModeCheck.Checked {
			g.otgModeCheck.SetChecked(false)
		}
		if g.controlOnlyCheck != nil && g.controlOnlyCheck.Checked {
			g.controlOnlyCheck.SetChecked(false)
		}
	})
	g.noControlCheck.SetChecked(false)

	g.otgModeCheck = widget.NewCheck("OTG Mode (Hardware KBM via USB)", func(checked bool) {
		if !checked {
			return
		}
		if g.noControlCheck != nil && g.noControlCheck.Checked {
			g.noControlCheck.SetChecked(false)
		}
		if g.controlOnlyCheck != nil && g.controlOnlyCheck.Checked {
			g.controlOnlyCheck.SetChecked(false)
		}
	})
	g.otgModeCheck.SetChecked(false)

	g.controlOnlyCheck = widget.NewCheck("Control Only (No Video/Audio, Wireless UHID)", func(checked bool) {
		if !checked {
			return
		}
		if g.noControlCheck != nil && g.noControlCheck.Checked {
			g.noControlCheck.SetChecked(false)
		}
		if g.otgModeCheck != nil && g.otgModeCheck.Checked {
			g.otgModeCheck.SetChecked(false)
		}
	})
	g.controlOnlyCheck.SetChecked(false)
}

func (g *GUI) initControlModifiers() {
	g.preferTextCheck = widget.NewCheck("Prefer Text Injection (Accents/Symbols)", func(checked bool) {
		if checked && g.rawKeyEventsCheck != nil && g.rawKeyEventsCheck.Checked {
			g.rawKeyEventsCheck.SetChecked(false)
		}
	})
	g.preferTextCheck.SetChecked(false)

	g.rawKeyEventsCheck = widget.NewCheck("Raw Key Events (Game/WASD Mode)", func(checked bool) {
		if checked && g.preferTextCheck != nil && g.preferTextCheck.Checked {
			g.preferTextCheck.SetChecked(false)
		}
	})
	g.rawKeyEventsCheck.SetChecked(false)

	g.noKeyRepeatCheck = widget.NewCheck("Disable Key Repeat on Hold", nil)
	g.noKeyRepeatCheck.SetChecked(false)

	g.noMouseHoverCheck = widget.NewCheck("Disable Mouse Hover", nil)
	g.noMouseHoverCheck.SetChecked(false)

	g.noClipboardAutosyncCheck = widget.NewCheck("Disable Clipboard Autosync", nil)
	g.noClipboardAutosyncCheck.SetChecked(false)

	g.legacyPasteCheck = widget.NewCheck("Legacy Paste (Inject Keys)", nil)
	g.legacyPasteCheck.SetChecked(false)
}

func (g *GUI) initControlWidgets() {
	g.initControlCombos()
	g.initControlModes()
	g.initControlModifiers()
}

func (g *GUI) buildControlModeRows() []fyne.CanvasObject {
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
	gpRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("gamepad mode:"),
		container.NewGridWrap(fyne.NewSize(120, 36), g.gamepadModeCombo),
		nil,
	)
	cheatsheetBtn := widget.NewButton("Shortcuts Cheatsheet", func() {
		g.showShortcutsCheatsheetDialog()
	})
	modRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("shortcut mod:"),
		container.NewHBox(
			container.NewGridWrap(fyne.NewSize(160, 36), g.shortcutModCombo),
			cheatsheetBtn,
		),
		nil,
	)
	pushRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("push target:"),
		container.NewGridWrap(fyne.NewSize(180, 36), g.pushTargetEntry),
		nil,
	)
	return []fyne.CanvasObject{kbdRow, mouseRow, gpRow, modRow, pushRow}
}

func (g *GUI) buildControlCheckboxes() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		g.noControlCheck,
		g.otgModeCheck,
		g.controlOnlyCheck,
		g.preferTextCheck,
		g.rawKeyEventsCheck,
		g.noKeyRepeatCheck,
		g.noMouseHoverCheck,
		g.noClipboardAutosyncCheck,
		g.legacyPasteCheck,
	}
}

func (g *GUI) buildControlTab() fyne.CanvasObject {
	hintLabel := widget.NewLabelWithStyle(
		"Tip: Press Left Alt (or Alt+Tab) to release mouse capture",
		fyne.TextAlignLeading,
		fyne.TextStyle{Italic: true},
	)

	items := make([]fyne.CanvasObject, 0, 17)
	items = append(items, g.buildControlModeRows()...)
	items = append(items, hintLabel)
	items = append(items, g.buildControlCheckboxes()...)

	return container.NewVScroll(container.NewVBox(items...))
}
