package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

type RenameFn func(newFileName string) bool

type RenameFileOpts struct {
	Title, Message string
	FileName       string
	Ack            RenameFn
	Cancel         cancelFunc
}

func ShowRenameFile(styles *config.Dialog, pages *ui.Pages, opts *RenameFileOpts) {
	f := newBaseModelForm(styles)
	f.SetItemPadding(0)

	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	modal := tview.NewModalForm("<"+opts.Title+">", f.Form)

	f.AddInputField("FileName:", opts.FileName, 0, nil, func(v string) {
		opts.FileName = v
	})

	f.AddButton("OK", func() {
		if !opts.Ack(opts.FileName) {
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
