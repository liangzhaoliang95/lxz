// 核心页面 承载数据库表列表和表数据

package view

import (
	"context"
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
	*ui.BaseFlex
	app             *App                      // 应用实例
	tableChangeChan chan tableChangeSubscribe // 用于接收表数据变更的消息
	// 数据库连接配置
	dbConnCfg *config.DBConnection
	// UI组件
	dbTree    *DatabaseDbTree    // 数据库树，用于显示数据库和表结构
	tableView *DatabaseTableView // 数据库表格视图，用于显示表数据

}

func (_this *DatabaseMainPage) Init(ctx context.Context) error {

	// 左侧数据库列表
	_this.dbTree = NewDatabaseDbTree(_this.app, _this.dbConnCfg, _this.tableChangeChan)
	_this.AddItem(_this.dbTree, 0, 3, true)

	// 初始化右侧的表格视图
	_this.tableView = NewDatabaseTableView(_this.app, _this.dbConnCfg, _this.tableChangeChan)
	_this.AddItem(_this.tableView, 0, 10, false)
	return nil
}

func (_this *DatabaseMainPage) Start() {
	slog.Info("DatabaseMainPage Start")
	// 启动数据库树的初始化
	err := _this.dbTree.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
	} else {
		_this.dbTree.Start()
	}

	// 启动表格视图的初始化
	err = _this.tableView.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
	} else {
		_this.tableView.Start()
	}
}

func (_this *DatabaseMainPage) Stop() {

}

func NewDatabaseMainPage(a *App, dbConnCfg *config.DBConnection) *DatabaseMainPage {
	var name = "DB Browser"
	lp := DatabaseMainPage{
		BaseFlex:        ui.NewBaseFlex(name),
		app:             a,
		dbConnCfg:       dbConnCfg,
		tableChangeChan: make(chan tableChangeSubscribe, 10), // 初始化消息通道
	}
	lp.SetDirection(tview.FlexColumn)

	return &lp
}
