// 启动页

package database_ui

import (
	"github.com/rivo/tview"
	"lxz/internal/ui"
)

type LaunchPage struct {
	*ui.BaseFlex
	connList *tview.Table // connList 用于显示连接列表
}

func NewLaunchPage() *LaunchPage {
	lp := LaunchPage{
		BaseFlex: ui.NewBaseFlex("LaunchPage"),
	}
	lp.SetBorder(true)
	lp.SetTitle("Database Launch")
	lp.SetTitleAlign(tview.AlignCenter)
	lp.SetBorderPadding(1, 1, 2, 2)

	return &lp
}
