// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"lxz/internal/config"
	"lxz/internal/ui"
)

type RestartFn func(*metav1.PatchOptions) bool

type CreateFileOpts struct {
	Title, Message string
	FieldManager   string
	Ack            RestartFn
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

	args := metav1.PatchOptions{
		FieldManager: opts.FieldManager,
	}
	f.AddInputField("FileName:", args.FieldManager, 0, nil, func(v string) {
		args.FieldManager = v
	})

	f.AddButton("OK", func() {
		if !opts.Ack(&args) {
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
