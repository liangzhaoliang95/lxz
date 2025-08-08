package dialog

import (
	"github.com/rivo/tview"
	"lxz/internal/config"
	"lxz/internal/ui"
)

const dialogKey = "dialog"

type confirmFunc func(force bool)

func ShowConfirmAck(
	app *ui.App,
	pages *ui.Pages,
	acceptStr string,
	override bool,
	title, msg string,
	ack confirmFunc,
	cancel cancelFunc,
) {
	styles := app.Styles.Dialog()

	f := newBaseModelForm(&styles)
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(styles.ButtonBgColor.Color()).
		SetButtonTextColor(styles.ButtonFgColor.Color()).
		SetLabelColor(styles.LabelFgColor.Color()).
		SetFieldTextColor(styles.FieldFgColor.Color())
	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		cancel()
	})

	var accept bool
	if override {
		changedFn := func(t string) {
			accept = (t == acceptStr)
		}
		f.AddInputField("Confirm:", "", 30, nil, changedFn)
	} else {
		accept = true
	}

	f.AddButton("OK", func() {
		if !accept {
			return
		}
		ack(f.force)
		dismissConfirm(pages)
		cancel()
	})
	for i := range 2 {
		b := f.GetButton(i)
		if b == nil {
			continue
		}
		b.SetBackgroundColorActivated(styles.ButtonFocusBgColor.Color())
		b.SetLabelColorActivated(styles.ButtonFocusFgColor.Color())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm("<"+title+">", f.Form)
	modal.SetText(msg)
	modal.SetTextColor(styles.FgColor.Color())
	modal.SetDoneFunc(func(int, string) {
		dismissConfirm(pages)
		cancel()
	})
	pages.AddPage(confirmKey, modal, false, false)
	pages.ShowPage(confirmKey)
}

// ShowConfirm pops a confirmation dialog.
func ShowConfirm(
	styles *config.Dialog,
	pages *ui.Pages,
	title, msg string,
	ack confirmFunc,
	cancel cancelFunc,
) {
	f := newBaseModelForm(styles)
	f.AddCheckbox("Force", false, func(checked bool) {
		f.force = checked
	})
	f.AddButton("Cancel", func() {
		dismiss(pages)
		cancel()
	})
	f.AddButton("OK", func() {
		ack(f.force)
		dismiss(pages)
		cancel()
	})
	for i := range 2 {
		if b := f.GetButton(i); b != nil {
			b.SetBackgroundColorActivated(styles.ButtonFocusBgColor.Color())
			b.SetLabelColorActivated(styles.ButtonFocusFgColor.Color())
		}
	}
	f.SetFocus(0)
	modal := tview.NewModalForm("<"+title+">", f.Form)
	modal.SetText(msg)
	modal.SetTextColor(styles.FgColor.Color())
	modal.SetDoneFunc(func(int, string) {
		dismiss(pages)
		cancel()
	})
	pages.AddPage(dialogKey, modal, false, false)
	pages.ShowPage(dialogKey)
}

func dismiss(pages *ui.Pages) {
	pages.RemovePage(dialogKey)
}
