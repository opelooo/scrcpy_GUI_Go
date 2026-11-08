package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

func resolveActiveModifier(selected string) string {
	clean := strings.TrimSpace(selected)
	if clean == "" || strings.HasPrefix(strings.ToLower(clean), "default") {
		return "Alt/Super"
	}
	return clean
}

func getScrcpyShortcuts(mod string) []domain.ShortcutItem {
	return []domain.ShortcutItem{
		// Window & Display
		{Category: "Window", Action: "Switch Fullscreen", KeySequence: fmt.Sprintf("%s + f (or F11)", mod), Description: "Toggle fullscreen display mode"},
		{Category: "Window", Action: "Switch Maximized", KeySequence: fmt.Sprintf("%s + m", mod), Description: "Maximize mirror window"},
		{Category: "Window", Action: "Rotate Display Left", KeySequence: fmt.Sprintf("%s + Left", mod), Description: "Rotate video stream left"},
		{Category: "Window", Action: "Rotate Display Right", KeySequence: fmt.Sprintf("%s + Right", mod), Description: "Rotate video stream right"},
		{Category: "Window", Action: "Flip Orientation", KeySequence: fmt.Sprintf("%s + x", mod), Description: "Flip portrait / landscape"},
		{Category: "Window", Action: "Pixel-Perfect (1:1)", KeySequence: fmt.Sprintf("%s + g", mod), Description: "Resize window to exact device resolution"},
		{Category: "Window", Action: "Remove Black Borders", KeySequence: fmt.Sprintf("%s + w (or double-click)", mod), Description: "Resize window to fit content bounds"},

		// Navigation
		{Category: "Navigation", Action: "HOME Button", KeySequence: fmt.Sprintf("%s + h (or middle-click)", mod), Description: "Navigate to home screen"},
		{Category: "Navigation", Action: "BACK Button", KeySequence: fmt.Sprintf("%s + b (or right-click)", mod), Description: "Navigate back"},
		{Category: "Navigation", Action: "App Switch / Recents", KeySequence: fmt.Sprintf("%s + s", mod), Description: "Open recent application switcher"},
		{Category: "Navigation", Action: "POWER Button", KeySequence: fmt.Sprintf("%s + p", mod), Description: "Press device power button"},
		{Category: "Navigation", Action: "Turn Screen Off", KeySequence: fmt.Sprintf("%s + o", mod), Description: "Keep mirroring while device screen is off"},
		{Category: "Navigation", Action: "Turn Screen On", KeySequence: fmt.Sprintf("%s + Shift + o", mod), Description: "Turn device physical screen back on"},

		// System & Panels
		{Category: "System", Action: "Expand Notifications", KeySequence: fmt.Sprintf("%s + n (or 5th-click)", mod), Description: "Pull down notification drawer"},
		{Category: "System", Action: "Expand Quick Settings", KeySequence: fmt.Sprintf("%s + Shift + n", mod), Description: "Pull down quick settings panel"},
		{Category: "System", Action: "Collapse Panels", KeySequence: fmt.Sprintf("%s + c", mod), Description: "Close notification and settings panels"},
		{Category: "System", Action: "Volume Up", KeySequence: fmt.Sprintf("%s + Up", mod), Description: "Increase audio volume"},
		{Category: "System", Action: "Volume Down", KeySequence: fmt.Sprintf("%s + Down (or 4th-click)", mod), Description: "Decrease audio volume"},
		{Category: "System", Action: "Physical Keyboard Settings", KeySequence: fmt.Sprintf("%s + k", mod), Description: "Open device physical keyboard settings"},

		// Clipboard & Input
		{Category: "Clipboard", Action: "Copy PC to Device", KeySequence: fmt.Sprintf("%s + v", mod), Description: "Paste computer clipboard to device"},
		{Category: "Clipboard", Action: "Copy Device to PC", KeySequence: fmt.Sprintf("%s + Shift + v", mod), Description: "Sync device clipboard back to computer"},
		{Category: "Clipboard", Action: "Inject Clipboard Text", KeySequence: fmt.Sprintf("%s + Ctrl + v", mod), Description: "Inject clipboard content as raw key events"},

		// Mouse Controls
		{Category: "Mouse", Action: "Right-Click", KeySequence: "Right-Click", Description: "Triggers BACK button (or turns screen on)"},
		{Category: "Mouse", Action: "Middle-Click", KeySequence: "Middle-Click", Description: "Triggers HOME button"},
		{Category: "Mouse", Action: "Double-Click", KeySequence: "Double-Click", Description: "Resizes mirror window to remove borders"},
		{Category: "Mouse", Action: "4th Button (Back)", KeySequence: "Mouse Button 4", Description: "Lowers volume"},
		{Category: "Mouse", Action: "5th Button (Forward)", KeySequence: "Mouse Button 5", Description: "Expands notifications"},
	}
}

func buildShortcutRow(it domain.ShortcutItem) fyne.CanvasObject {
	actionLabel := widget.NewLabelWithStyle(it.Action, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	keyLabel := widget.NewLabelWithStyle(it.KeySequence, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	descLabel := widget.NewLabel(it.Description)
	descLabel.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		nil, nil,
		container.NewGridWrap(fyne.NewSize(200, 32), actionLabel),
		nil,
		container.NewBorder(
			nil, nil,
			container.NewGridWrap(fyne.NewSize(230, 32), keyLabel),
			nil,
			descLabel,
		),
	)
}

func buildCategoryTable(items []domain.ShortcutItem) fyne.CanvasObject {
	rows := make([]fyne.CanvasObject, 0, len(items)+2)

	headerAction := widget.NewLabelWithStyle("Action", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	headerKey := widget.NewLabelWithStyle("Shortcut", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	headerDesc := widget.NewLabelWithStyle("Description", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	header := container.NewBorder(
		nil, nil,
		container.NewGridWrap(fyne.NewSize(200, 32), headerAction),
		nil,
		container.NewBorder(
			nil, nil,
			container.NewGridWrap(fyne.NewSize(230, 32), headerKey),
			nil,
			headerDesc,
		),
	)
	rows = append(rows, header, widget.NewSeparator())

	for _, it := range items {
		rows = append(rows, buildShortcutRow(it))
	}
	return container.NewVBox(rows...)
}

func buildShortcutsCategoryTabs(shortcuts []domain.ShortcutItem) *container.AppTabs {
	tabs := container.NewAppTabs()
	categories := []string{"Window", "Navigation", "System", "Clipboard", "Mouse"}

	for _, cat := range categories {
		var catItems []domain.ShortcutItem
		for _, it := range shortcuts {
			if it.Category == cat {
				catItems = append(catItems, it)
			}
		}
		table := buildCategoryTable(catItems)
		tabs.Append(container.NewTabItem(cat, container.NewVScroll(table)))
	}
	return tabs
}

func (g *GUI) showShortcutsCheatsheetDialog() {
	currentMod := "Default (Alt/Super)"
	if g.shortcutModCombo != nil && g.shortcutModCombo.Selected != "" {
		currentMod = g.shortcutModCombo.Selected
	}

	activeMod := resolveActiveModifier(currentMod)
	shortcuts := getScrcpyShortcuts(activeMod)

	headerText := fmt.Sprintf("Active Modifier: %s", activeMod)
	headerLabel := widget.NewLabelWithStyle(headerText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	tabs := buildShortcutsCategoryTabs(shortcuts)
	content := container.NewBorder(headerLabel, nil, nil, nil, tabs)

	dlg := dialog.NewCustom("scrcpy Shortcuts Cheatsheet", "Close", content, g.window)
	dlg.Resize(fyne.NewSize(840, 520))
	dlg.Show()
}

