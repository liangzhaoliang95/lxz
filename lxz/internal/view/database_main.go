// 核心页面 承载数据库表列表和表数据

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/ui"
)

type tableChangeSubscribe struct {
	tableName string // 表名
	dbName    string // 数据库名
}

type DatabaseMainPage struct {
	*BaseFlex
	app             *App                      // 应用实例
	tableChangeChan chan tableChangeSubscribe // 用于接收表数据变更的消息
	// 数据库连接配置
	dbConnCfg *config.DBConnection
	// UI组件
	dbTree    *DatabaseDbTree    // 数据库树，用于显示数据库和表结构
	tableView *DatabaseTableView // 数据库表格视图，用于显示表数据

}

func (_this *DatabaseMainPage) TabFocusChange(event *tcell.EventKey) *tcell.EventKey {
	if _this.app.UI.GetFocus() == _this.dbTree.databaseUiTree {
		_this.tableView.selfFocus()
	} else {
		_this.dbTree.selfFocus()
	}
	return nil
}

// ToggleSearch 触发搜索功能
func (_this *DatabaseMainPage) ToggleSearch(evt *tcell.EventKey) *tcell.EventKey {
	// 触发搜索功能
	slog.Info("DatabaseMainPage Search triggered", "event", evt)
	if ui.IsInputPrimitive(appUiInstance.GetFocus()) {
		slog.Info("Search input is already focused, ignoring toggle search event")
		return evt
	}
	currentPage := _this.tableView.tableComponents[_this.tableView.currentPageKey]
	if currentPage == nil {
		appUiInstance.Flash().Err(fmt.Errorf("select one table first"))
	} else {
		// 将焦点定位到输入框上
		_this.tableView.tableComponents[_this.tableView.currentPageKey].focusSearch()
	}
	return nil

}

func (_this *DatabaseMainPage) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		ui.KeySlash:     ui.NewKeyAction("Search", _this.ToggleSearch, true),
		tcell.KeyCtrlO:  ui.NewKeyAction("Open Query Page", _this.goToQueryPage, true),
		tcell.KeyEscape: ui.NewKeyAction("Last Page", _this.EmptyKeyEvent, true),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
	})
}

func (_this *DatabaseMainPage) goToQueryPage(evt *tcell.EventKey) *tcell.EventKey {

	// 跳到手动查询页面
	_this.app.inject(NewDatabaseQueryView(_this.app, _this.dbConnCfg), false)
	return nil
}

func (_this *DatabaseMainPage) Init(ctx context.Context) error {
	slog.Info("DatabaseMainPage Init")
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 左侧数据库列表
	_this.dbTree = NewDatabaseDbTree(_this.app, _this.dbConnCfg, _this.tableChangeChan)
	err := _this.dbTree.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return err
	}
	_this.AddItem(_this.dbTree, 0, 3, true)

	// 初始化右侧的表格视图
	_this.tableView = NewDatabaseTableView(_this.app, _this.dbConnCfg, _this.tableChangeChan)
	err = _this.tableView.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return err
	}
	_this.AddItem(_this.tableView, 0, 10, false)
	return nil
}

func (_this *DatabaseMainPage) Start() {
	slog.Info("DatabaseMainPage Start")

	// 启动数据库树的初始化
	_this.dbTree.Start()

	// 启动表格视图的初始化
	_this.tableView.Start()
}

func (_this *DatabaseMainPage) Stop() {

}

func NewDatabaseMainPage(a *App, dbConnCfg *config.DBConnection) *DatabaseMainPage {
	var name = "Table View"
	lp := DatabaseMainPage{
		BaseFlex:        NewBaseFlex(name),
		app:             a,
		dbConnCfg:       dbConnCfg,
		tableChangeChan: make(chan tableChangeSubscribe, 10), // 初始化消息通道
	}
	lp.SetDirection(tview.FlexColumn)

	return &lp
}
