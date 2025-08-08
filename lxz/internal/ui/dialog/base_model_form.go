package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"lxz/internal/config"
)

type baseModelForm struct {
	*tview.Form
	force bool
}

func newBaseModelForm(styles *config.Dialog) *baseModelForm {

	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.ColorBlue).
		SetButtonTextColor(tcell.ColorBlack).
		SetLabelColor(styles.LabelFgColor.Color()).
		SetFieldTextColor(styles.FieldFgColor.Color())
	b := &baseModelForm{
		Form: f,
	}
	return b
}
