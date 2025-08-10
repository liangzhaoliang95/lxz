package dialog

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/drivers/redis_drivers"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
	"log/slog"
	"strconv"
)

type CreateUpdateRedisDataFn func(data *redis_drivers.RedisData) bool

type CreateUpdateRedisDataOpts struct {
	Title, Message string
	Data           *redis_drivers.RedisData
	Ack            CreateUpdateRedisDataFn
	Cancel         cancelFunc
}

func ShowCreateUpdateRedisData(
	styles *config.Dialog,
	pages *ui.Pages,
	opts *CreateUpdateRedisDataOpts,
) {
	f := newBaseModelForm(styles)

	f.AddTextArea("Key:", opts.Data.KeyName, 0, 2, 0, func(text string) {
		opts.Data.KeyName = text
	})
	f.AddTextView("Type:", opts.Data.KetType, 0, 1, true, false)

	f.AddInputField(
		"TTL:",
		fmt.Sprintf("%d", opts.Data.KeyTTL),
		0,
		tview.InputFieldInteger,
		func(v string) {
			ttl, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				slog.Error("Invalid port number", "port", v, "error", err)
				return
			}
			opts.Data.KeyTTL = ttl
		},
	)

	f.AddTextArea("Value:", opts.Data.KeyValue, 0, 3, 0, func(v string) {
		opts.Data.KeyValue = v
	})

	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	f.AddButton("OK", func() {
		if !opts.Ack(opts.Data) {
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
		b.SetBackgroundColor(tcell.ColorYellow)
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
