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

type redisDbChangeSubscribe struct {
	dbNum int
}

type RedisMainPage struct {
	*BaseFlex
	app               *App                        // 应用实例
	redisDbChangeChan chan redisDbChangeSubscribe // 用于接收表数据变更的消息
	// 数据库连接配置
	redisConnConfig *config.RedisConnConfig
	// UI组件
	dbListUI    *RedisDbListView // 库列表, 用于显示db列表
	tableViewUI *RedisDataView   // 数据库表格视图，用于显示表数据

}

func (_this *RedisMainPage) TabFocusChange(event *tcell.EventKey) *tcell.EventKey {
	if _this.app.UI.GetFocus() == _this.dbListUI.dbListUI {
		_this.tableViewUI.selfFocus()
	} else {
		_this.dbListUI.selfFocus()
	}
	return nil
}

// ToggleSearch 触发搜索功能
func (_this *RedisMainPage) ToggleSearch(evt *tcell.EventKey) *tcell.EventKey {
	// 触发搜索功能
	slog.Info("Search triggered", "event", evt)
	currentPage := _this.tableViewUI.redisDataComponents[_this.tableViewUI.currentPageKey]
	if currentPage == nil {
		appUiInstance.Flash().Err(fmt.Errorf("select one table first"))
	} else {
		// 将焦点定位到输入框上
		_this.tableViewUI.redisDataComponents[_this.tableViewUI.currentPageKey].focusSearch()
	}
	return nil

}

func (_this *RedisMainPage) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		ui.KeySlash:     ui.NewKeyAction("Search", _this.ToggleSearch, true),
		tcell.KeyEscape: ui.NewKeyAction("Last Page", _this.EmptyKeyEvent, true),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
	})
}

func (_this *RedisMainPage) Init(ctx context.Context) error {
	slog.Info("RedisMainPage Init")
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 左侧数据库列表
	_this.dbListUI = NewRedisDbTree(_this.app, _this.redisConnConfig, _this.redisDbChangeChan)
	err := _this.dbListUI.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return err
	}
	_this.AddItem(_this.dbListUI, 20, 0, true)
	//
	// 初始化右侧的表格视图
	_this.tableViewUI = NewRedisDataView(_this.app, _this.redisConnConfig, _this.redisDbChangeChan)
	err = _this.tableViewUI.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return err
	}
	_this.AddItem(_this.tableViewUI, 0, 10, false)
	return nil
}

func (_this *RedisMainPage) Start() {
	slog.Info("RedisMainPage Start")

	//// 启动数据库树的初始化
	_this.dbListUI.Start()
	//
	//// 启动表格视图的初始化
	_this.tableViewUI.Start()
}

func (_this *RedisMainPage) Stop() {

}

func NewRedisMainPage(a *App, connConfig *config.RedisConnConfig) *RedisMainPage {
	var name = "Redis View"
	lp := RedisMainPage{
		BaseFlex:          NewBaseFlex(name),
		app:               a,
		redisConnConfig:   connConfig,
		redisDbChangeChan: make(chan redisDbChangeSubscribe, 10), // 初始化消息通道
	}
	lp.SetDirection(tview.FlexColumn)

	return &lp
}
