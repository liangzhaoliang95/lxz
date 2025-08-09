package render

import (
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

// Truncate a string to the given l and suffix ellipsis if needed.
func Truncate(str string, width int) string {
	return runewidth.Truncate(str, width, string(tview.SemigraphicsHorizontalEllipsis))
}
