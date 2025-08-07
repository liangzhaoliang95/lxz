package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"lxz/internal/config"
	"lxz/internal/ui"
)

type CreateFileFn func(name string, isDir bool) bool

type CreateFileOpts struct {
	Title, Message string
	FileName       string
	IsDir          bool
	Ack            CreateFileFn
	Cancel         cancelFunc
}

func ShowCreateFile(styles *config.Dialog, pages *ui.Pages, opts *CreateFileOpts) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.ColorBlue).
		SetButtonTextColor(tcell.ColorBlack).
		SetLabelColor(styles.LabelFgColor.Color()).
		SetFieldTextColor(styles.FieldFgColor.Color())
	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	modal := tview.NewModalForm("<"+opts.Title+">", f)

	f.AddInputField("Name:", opts.FileName, 0, nil, func(v string) {
		opts.FileName = v
	})
	f.AddCheckbox("IsDir:", false, func(c bool) {
		opts.IsDir = c
	})

	f.AddButton("OK", func() {
		if !opts.Ack(opts.FileName, opts.IsDir) {
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
