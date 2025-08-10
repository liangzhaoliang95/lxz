package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/tview"
	"lxz/internal/config"
	"lxz/internal/ui"
)

type DeleteDatabaseConnectionFn func(key string) bool

type DeleteDatabaseConnectionOpts struct {
	Title, Message string
	DBConnection   *config.DBConnection
	Ack            DeleteDatabaseConnectionFn
	Cancel         cancelFunc
	SelectKey      string
}

func ShowDeleteCreateDatabaseConnection(
	styles *config.Dialog,
	pages *ui.Pages,
	opts *DeleteDatabaseConnectionOpts,
) {
	f := newBaseModelForm(styles)

	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	f.AddButton("OK", func() {
		if !opts.Ack(opts.SelectKey) {
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

	modal := tview.NewModalForm("<"+opts.Title+">", f.Form)
	modal.SetText(message)
	modal.SetTextColor(styles.FgColor.Color())
	modal.SetDoneFunc(func(int, string) {
		dismissConfirm(pages)
		opts.Cancel()
	})
	pages.AddPage(confirmKey, modal, false, false)
	pages.ShowPage(confirmKey)
}
