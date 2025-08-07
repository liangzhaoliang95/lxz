// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dialog

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/redis_drivers"
	"lxz/internal/ui"
	"strconv"
)

type CreateUpdateRedisDataFn func(connection *config.RedisConnConfig) bool

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

	f.AddTextView("Type:", opts.Config.Name, 0, 0, true, false)

	f.AddInputField("TTL:", opts.Config.Host, 0, nil, func(v string) {
		opts.Config.Host = v
	})

	f.AddTextView("Value:", opts.Config.Host, 0, nil, func(v string) {
		opts.Config.Host = v
	})

	f.AddInputField(
		"Port:",
		fmt.Sprintf("%d", opts.Config.Port),
		0,
		tview.InputFieldInteger,
		func(v string) {
			port, err := strconv.Atoi(v)
			if err != nil {
				slog.Error("Invalid port number", "port", v, "error", err)
			}
			opts.Config.Port = int64(port)
		},
	)

	f.AddInputField("UserName:", opts.Config.UserName, 0, nil, func(v string) {
		opts.Config.UserName = v
	})

	f.AddInputField("Password:", opts.Config.Password, 0, nil, func(v string) {
		opts.Config.Password = v
	})

	f.AddButton("Test", func() {
		// 测试数据库能否连接
		opts.Test(opts.Config)
	})

	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	f.AddButton("OK", func() {
		if !opts.Ack(opts.Config) {
			return
		}
		dismissConfirm(pages)
		opts.Cancel()
	})
	for i := range 3 {
		b := f.GetButton(i)
		if b == nil {
			continue
		}
		b.SetBackgroundColor(tcell.ColorYellow)
	}
	f.SetFocus(0)

	message := opts.Message

	modal := tview.NewModalForm("<"+opts.Title+">", f)
	modal.SetText(message)
	modal.SetTextColor(styles.FgColor.Color())
	modal.SetDoneFunc(func(int, string) {
		dismissConfirm(pages)
		opts.Cancel()
	})
	pages.AddPage(confirmKey, modal, false, false)
	pages.ShowPage(confirmKey)
}
