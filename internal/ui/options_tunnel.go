package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (g *GUI) initTunnelWidgets() {
	g.adbServerHostEntry = widget.NewEntry()
	g.adbServerHostEntry.SetPlaceHolder("e.g. 192.168.1.2:5037 or localhost:5038")

	g.tunnelHostEntry = widget.NewEntry()
	g.tunnelHostEntry.SetPlaceHolder("e.g. 192.168.1.2")

	g.tunnelPortEntry = widget.NewEntry()
	g.tunnelPortEntry.SetPlaceHolder("e.g. 27183")

	g.forceAdbForwardCheck = widget.NewCheck("Force ADB Forward (No reverse connection)", nil)
	g.forceAdbForwardCheck.SetChecked(false)

	g.portRangeEntry = widget.NewEntry()
	g.portRangeEntry.SetPlaceHolder("e.g. 27183 or 27183:27199")
}

func (g *GUI) buildTunnelTab() *fyne.Container {
	adbServerRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Remote ADB Server:"),
		nil,
		g.adbServerHostEntry,
	)

	tunnelHostRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Tunnel Host IP:"),
		nil,
		g.tunnelHostEntry,
	)

	tunnelPortRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Tunnel Port:"),
		nil,
		g.tunnelPortEntry,
	)

	portRangeRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("TCP Port Range:"),
		nil,
		g.portRangeEntry,
	)

	return container.NewVBox(
		adbServerRow,
		widget.NewSeparator(),
		tunnelHostRow,
		tunnelPortRow,
		g.forceAdbForwardCheck,
		widget.NewSeparator(),
		portRangeRow,
	)
}
