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
	dbListViewUI *RedisDbListView // 库列表, 用于显示db列表
	dataViewUI   *RedisDataView   // 数据库表格视图，用于显示表数据

}

func (_this *RedisMainPage) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		ui.KeySlash:     ui.NewKeyAction("Search", _this.ToggleSearch, true),
		tcell.KeyEscape: ui.NewKeyAction("Last Page", _this.EmptyKeyEvent, true),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
		tcell.KeyCtrlD:  ui.NewKeyAction("Delete Key", _this.DeleteKey, true),
		tcell.KeyCtrlN:  ui.NewKeyAction("Delete Key", _this.NewKey, true),
		tcell.KeyCtrlR:  ui.NewKeyAction("Refresh", _this.Refresh, true),
		ui.KeyE:         ui.NewKeyAction("Edit Key", _this.EditKey, true),
	})
}

// Refresh 刷新当前页面
func (_this *RedisMainPage) Refresh(event *tcell.EventKey) *tcell.EventKey {
	slog.Info("Refreshing RedisMainPage")
	// 刷新当前页面
	currentCompPage := _this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey]
	if currentCompPage == nil {
		_this.app.UI.Flash().Err(fmt.Errorf("select one table first"))
		return nil
	}
	// 刷新当前页面的表格数据
	err := currentCompPage.refreshData()
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return nil
	}
	return event
}

// NewKey 创建一个新的键
func (_this *RedisMainPage) NewKey(event *tcell.EventKey) *tcell.EventKey {
	currentCompPage := _this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey]
	if _this.app.UI.GetFocus() != currentCompPage.keyGroupTree {
		_this.app.UI.Flash().Err(fmt.Errorf("please select a key first"))
		return nil
	}
	// 当前焦点在键列表上，可以执行新建操作
	// 给keyGroupTree组件发送新建键的命令
	currentCompPage.newKey()
	return nil
}

// EditKey 编辑当前选中的键
func (_this *RedisMainPage) EditKey(event *tcell.EventKey) *tcell.EventKey {
	currentCompPage := _this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey]
	if _this.app.UI.GetFocus() != currentCompPage.keyGroupTree {
		_this.app.UI.Flash().Err(fmt.Errorf("please select a key first"))
		return nil
	}
	// 当前焦点在键列表上，可以执行编辑操作
	// 给keyGroupTree组件发送编辑键的命令
	currentCompPage.editKey()
	return nil
}

// DeleteKey 删除当前选中的键
func (_this *RedisMainPage) DeleteKey(event *tcell.EventKey) *tcell.EventKey {
	currentCompPage := _this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey]
	if _this.app.UI.GetFocus() != currentCompPage.keyGroupTree {
		_this.app.UI.Flash().Err(fmt.Errorf("please select a key first"))
		return nil
	}
	// 当前焦点在键列表上，可以执行删除操作
	// 给keyGroupTree组件发送删除键的命令
	currentCompPage.deleteKey()
	return nil
}

func (_this *RedisMainPage) TabFocusChange(event *tcell.EventKey) *tcell.EventKey {
	currentCompPage := _this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey]
	if _this.app.UI.GetFocus() == _this.dbListViewUI.dbListUI ||
		_this.app.UI.GetFocus() == _this.dbListViewUI {
		_this.dataViewUI.selfFocus()
	} else if _this.app.UI.GetFocus() == currentCompPage.KeyValue {
		_this.dataViewUI.selfFocus()
	} else {
		_this.dbListViewUI.selfFocus()
	}
	return nil
}

// ToggleSearch 触发搜索功能
func (_this *RedisMainPage) ToggleSearch(evt *tcell.EventKey) *tcell.EventKey {
	// 触发搜索功能
	slog.Info("Search triggered", "event", evt)
	currentPage := _this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey]
	if currentPage == nil {
		appUiInstance.Flash().Err(fmt.Errorf("select one table first"))
	} else {
		// 将焦点定位到输入框上
		_this.dataViewUI.redisDataComponents[_this.dataViewUI.currentPageKey].focusSearch()
	}
	return nil

}

func (_this *RedisMainPage) Init(ctx context.Context) error {
	slog.Info("RedisMainPage Init")
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 左侧数据库列表
	_this.dbListViewUI = NewRedisDbTree(_this.app, _this.redisConnConfig, _this.redisDbChangeChan)
	err := _this.dbListViewUI.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return err
	}
	_this.AddItem(_this.dbListViewUI, 20, 0, true)

	// 初始化右侧的表格视图
	_this.dataViewUI = NewRedisDataView(_this.app, _this.redisConnConfig, _this.redisDbChangeChan)
	err = _this.dataViewUI.Init(context.Background())
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return err
	}
	_this.AddItem(_this.dataViewUI, 0, 10, false)
	return nil
}

func (_this *RedisMainPage) Start() {
	slog.Info("RedisMainPage Start")

	// 启动数据库树的初始化
	_this.dbListViewUI.Start()

	// 启动表格视图的初始化
	_this.dataViewUI.Start()

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
