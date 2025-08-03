// 核心页面 承载数据库表列表和表数据

package view

import (
	"context"
	"github.com/rivo/tview"
	"lxz/internal/ui"
)

type TablePage struct {
	*ui.BaseFlex
	app          *App            // 应用实例
	databaseTree *tview.TreeView // 用于显示数据库树
}

func (_this *TablePage) Init(ctx context.Context) error {
	// 初始化数据库连接
	// 获取当前连接下的数据库列表

	return nil
}

func (_this *TablePage) Start() {

}

func (_this *TablePage) Stop() {

}

func NewTablePage(a *App) *TablePage {
	var name = "Table Browser"
	lp := TablePage{
		BaseFlex: ui.NewBaseFlex(name),
		app:      a,
	}

	return &lp
}
