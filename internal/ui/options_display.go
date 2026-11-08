package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initDisplayWidgets() {
	g.displayIDEntry = widget.NewEntry()
	g.displayIDEntry.SetPlaceHolder("0 = main display")

	g.newDisplayCheck = widget.NewCheck("Create Virtual Display (Secondary Screen)", func(checked bool) {
		if checked && g.displayIDEntry.Text != "" {
			g.displayIDEntry.SetText("")
		}
	})
	g.newDisplayCheck.SetChecked(false)

	g.newDisplayResEntry = widget.NewEntry()
	g.newDisplayResEntry.SetPlaceHolder("e.g. 1920x1080 or /240 (empty = auto)")

	g.noVDSystemDecorationsCheck = widget.NewCheck("Disable System Decorations (No Launcher)", nil)
	g.noVDSystemDecorationsCheck.SetChecked(false)

	g.noVDDestroyContentCheck = widget.NewCheck("Preserve Apps on Close (Move to Main)", nil)
	g.noVDDestroyContentCheck.SetChecked(false)

	g.displayIMEPolicyCombo = widget.NewSelect([]string{"default", "local", "fallback", "hide"}, nil)
	g.displayIMEPolicyCombo.SetSelected("default")

	g.flexDisplayCheck = widget.NewCheck("Flex Display (Resize dynamically, scrcpy >= 3.4)", nil)
	g.flexDisplayCheck.SetChecked(false)
	if !g.caps.HasFlexDisplay {
		g.flexDisplayCheck.Hide()
	}
}

func (g *GUI) buildDisplayRows() []fyne.CanvasObject {
	dispIDRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Mirror Display ID:"),
		nil,
		g.displayIDEntry,
	)

	resRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Virtual Resolution/DPI:"),
		nil,
		g.newDisplayResEntry,
	)

	imeRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Virtual IME Policy:"),
		nil,
		g.displayIMEPolicyCombo,
	)

	return []fyne.CanvasObject{
		dispIDRow,
		widget.NewSeparator(),
		g.newDisplayCheck,
		resRow,
		g.noVDSystemDecorationsCheck,
		g.noVDDestroyContentCheck,
		imeRow,
		g.flexDisplayCheck,
	}
}

func (g *GUI) buildDisplayAndWindowTab() fyne.CanvasObject {
	windowHeader := widget.NewLabelWithStyle("Window Modes & Geometry", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	displayHeader := widget.NewLabelWithStyle("Display & Virtual Screen", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	items := make([]fyne.CanvasObject, 0, 25)
	items = append(items, windowHeader)
	items = append(items, g.buildWindowRows()...)
	items = append(items, widget.NewSeparator(), displayHeader)
	items = append(items, g.buildDisplayRows()...)

	return container.NewVScroll(container.NewVBox(items...))
}
