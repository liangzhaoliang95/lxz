package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

type (
	okFunc     func(force bool)
	cancelFunc func()
)

// ShowDelete pops a resource deletion dialog.
func ShowDelete(styles *config.Dialog, pages *ui.Pages, msg string, ok okFunc, cancel cancelFunc) {
	force := false
	f := newBaseModelForm(styles)
	f.SetItemPadding(0)

	f.AddButton("Cancel", func() {
		dismiss(pages)
		cancel()
	})
	f.AddButton("OK", func() {
		ok(force)
		dismiss(pages)
		cancel()
	})
	for i := range 2 {
		b := f.GetButton(i)
		if b == nil {
			continue
		}
		b.SetBackgroundColor(tcell.ColorYellow)
	}
	f.SetFocus(2)

	confirm := tview.NewModalForm("<Delete>", f.Form)
	confirm.SetText(msg)
	confirm.SetTextColor(styles.FgColor.Color())
	confirm.SetDoneFunc(func(int, string) {
		dismiss(pages)
		cancel()
	})
	pages.AddPage(dialogKey, confirm, false, false)
	pages.ShowPage(dialogKey)
}
