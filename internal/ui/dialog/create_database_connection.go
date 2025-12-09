package dialog

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

type CreateDatabaseConnectionFn func(connection *config.DBConnection) bool

type CreateDatabaseConnectionOpts struct {
	Title, Message string
	DBConnection   *config.DBConnection
	Ack            CreateDatabaseConnectionFn
	Test           func(connection *config.DBConnection) bool
	Cancel         cancelFunc
}

func ShowCreateCreateDatabaseConnection(
	styles *config.Dialog,
	pages *ui.Pages,
	opts *CreateDatabaseConnectionOpts,
) {
	f := newBaseModelForm(styles)

	providerMap := make(map[string]int)
	for i := 0; i < len(config.DatabaseProviderList); i++ {
		provider := config.DatabaseProviderList[i]
		providerMap[provider] = i
	}
	slog.Info("providerMap", "providerMap", providerMap)

	modal := tview.NewModalForm("<"+opts.Title+">", f.Form)
	f.AddDropDown(
		"Provider:",
		config.DatabaseProviderList,
		providerMap[opts.DBConnection.Provider],
		func(s string, i int) {
			opts.DBConnection.Provider = s
		},
	)

	f.AddInputField("Name:", opts.DBConnection.Name, 0, nil, func(v string) {
		opts.DBConnection.Name = v
	})
	f.AddInputField("UserName:", opts.DBConnection.UserName, 0, nil, func(v string) {
		opts.DBConnection.UserName = v
	})

	f.AddInputField("Password:", opts.DBConnection.Password, 0, nil, func(v string) {
		opts.DBConnection.Password = v
	})

	f.AddInputField("Host:", opts.DBConnection.Host, 0, nil, func(v string) {
		opts.DBConnection.Host = v
	})
	f.AddInputField(
		"Port:",
		fmt.Sprintf("%d", opts.DBConnection.Port),
		0,
		tview.InputFieldInteger,
		func(v string) {
			port, err := strconv.Atoi(v)
			if err != nil {
				slog.Error("Invalid port number", "port", v, "error", err)
			}
			opts.DBConnection.Port = int64(port)
		},
	)
	f.AddInputField("DBName:", "", 0, nil, func(v string) {
		opts.DBConnection.DBName = v
	})

	f.AddButton("Test", func() {
		// 测试数据库能否连接
		opts.Test(opts.DBConnection)
	})

	f.AddButton("Cancel", func() {
		dismissConfirm(pages)
		opts.Cancel()
	})

	f.AddButton("OK", func() {
		if !opts.Ack(opts.DBConnection) {
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
		//b.SetBackgroundColorActivated(tcell.ColorRed)
		//b.SetLabelColorActivated(tcell.ColorWhite)
		//b.SetBorder(true)
		b.SetBackgroundColor(tcell.ColorYellow)
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
