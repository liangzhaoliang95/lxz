// 启动页

package database_ui

import (
	"context"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/ui"
	"strings"
)

type LaunchPage struct {
	*ui.BaseFlex
	connList *tview.List // connList 用于显示连接列表
}

func (_this *LaunchPage) Init(ctx context.Context) error {
	return nil
}

func (_this *LaunchPage) Start() {
	slog.Info("🐶 lxz database browser launch page starting...")
	// 启动页的初始化逻辑
	// 这里可以添加一些组件，比如连接列表、按钮等
	// 目前只是一个空的启动页
	_this.connList = tview.NewList()

	_this.connList.SetBorder(true)
	lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	for i := 0; i < len(lorem); i++ {
		_this.connList.AddItem(lorem[i], "", 0, func() {})
	}
	_this.AddItem(_this.connList, 0, 1, true)
}

func (_this *LaunchPage) Stop() {

}

func NewLaunchPage() *LaunchPage {
	var name = "Database Launch"
	lp := LaunchPage{
		BaseFlex: ui.NewBaseFlex(name),
	}
	lp.SetBorder(true)
	lp.SetTitleAlign(tview.AlignCenter)

	return &lp
}
