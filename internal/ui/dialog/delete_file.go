package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

type DeleteFn func() bool

type DeleteFileOpts struct {
	Title, Message string
	FieldManager   string
	Ack            DeleteFn
	Cancel         cancelFunc
}

func ShowDeleteFile(styles *config.Dialog, pages *ui.Pages, opts *DeleteFileOpts) {
	f := newBaseModelForm(styles)

	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	modal := tview.NewModalForm("<"+opts.Title+">", f.Form)

	f.AddButton("OK", func() {
		if !opts.Ack() {
			return
		}
		dismissConfirm(pages)
		opts.Cancel()
	})
	for i := range 2 {
		b := f.GetButton(i)
		if b == nil {
			continue
		}
		b.SetBackgroundColorActivated(tcell.ColorRed)
		b.SetLabelColorActivated(tcell.ColorWhite)
	}
	f.SetFocus(0)

	message := opts.Message
	modal.SetText(message)
	modal.SetTextColor(styles.FgColor.Color())
	modal.SetDoneFunc(func(int, string) {
		dismissConfirm(pages)
		opts.Cancel()
	})
	pages.AddPage(confirmKey, modal, false, false)
	pages.ShowPage(confirmKey)
}
