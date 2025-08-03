package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"lxz/internal/config"
)

func newBaseModelForm(styles *config.Dialog) *tview.Form {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.ColorBlue).
		SetButtonTextColor(tcell.ColorBlack).
		SetLabelColor(styles.LabelFgColor.Color()).
		SetFieldTextColor(styles.FieldFgColor.Color())
	return f
}
