package domain

// ShortcutItem represents an individual scrcpy keyboard or mouse shortcut entry.
type ShortcutItem struct {
	Category    string
	Action      string
	KeySequence string
	Description string
}
