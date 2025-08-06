/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 10:50
 */

package view

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/ui/dialog"
	"lxz/internal/view/base"
)

type RedisDataView struct {
	*BaseFlex
	redisDbChangeChan   chan redisDbChangeSubscribe // 表格数据变更订阅通道
	app                 *App
	redisConnConfig     *config.RedisConnConfig        // 数据库连接配置
	redisDbPages        *tview.Pages                   // 表格数据页面容器
	redisDataComponents map[string]*RedisDataComponent // 表格数据页面
	currentPageKey      string                         // 当前页面的键，用于切换页面
}

func (_this *RedisDataView) selfFocus() {

	comp := _this.redisDataComponents[_this.currentPageKey]
	if comp == nil {
		_this.app.UI.SetFocus(_this)
	} else {
		// 设置当前焦点为表格组件
		_this.app.UI.SetFocus(comp.keyGroupTree)
	}

}

func (_this *RedisDataView) LunchPage(dbNum int) error {
	slog.Info("Launching page for table", "dbNum", dbNum)

	loading := dialog.ShowLoadingDialog(
		appViewInstance.Content.Pages,
		"",
		appUiInstance.ForceDraw,
	)
	pageKey := fmt.Sprintf("%d", dbNum)
	if _, ok := _this.redisDataComponents[pageKey]; !ok {
		// 当前页不存在，创建一个新的表格组件
		tableComponent := NewRedisDataComponent(_this.app, dbNum, _this.redisConnConfig)
		slog.Info("Initializing table component", "pageKey", pageKey)
		err := tableComponent.Init(context.Background())
		if err != nil {
			slog.Error("Failed to initialize table component", "err", err)
			return err
		}
		_this.redisDataComponents[pageKey] = tableComponent
		_this.redisDbPages.AddPage(pageKey, tableComponent, true, true)
	}
	// 当前页已存在，切换到该页面
	_this.currentPageKey = pageKey
	_this.redisDbPages.SwitchToPage(pageKey)

	comp := _this.redisDataComponents[pageKey]
	comp.Start()

	loading.Hide()
	_this.selfFocus()
	appUiInstance.ForceDraw()
	return nil

}

func (_this *RedisDataView) Init(ctx context.Context) error {
	_this.AddItem(_this.redisDbPages, 0, 1, true)
	return nil
}

func (_this *RedisDataView) Start() {
	go func() {
		slog.Info("Starting RedisDataView...")
		// 监听changeChan通道，处理表格数据变更
		for change := range _this.redisDbChangeChan {
			slog.Info(
				"Received redis db change notification",
				"dbName",
				change.dbNum,
			)
			if err := _this.LunchPage(change.dbNum); err != nil {
				slog.Info("Failed to launch page for dbNum", "dbNum", change.dbNum, "error", err)
			}
		}
		slog.Info("886")
	}()
}

func (_this *RedisDataView) Stop() {

}

// --- data helpers ---

func NewRedisDataView(
	a *App,
	redisConnConfig *config.RedisConnConfig,
	redisDbChangeChan chan redisDbChangeSubscribe,
) *RedisDataView {
	var name = "Table View"
	lp := RedisDataView{
		BaseFlex:            NewBaseFlex(name),
		app:                 a,
		redisDbPages:        tview.NewPages(),
		redisDataComponents: make(map[string]*RedisDataComponent),
		redisConnConfig:     redisConnConfig,
		redisDbChangeChan:   redisDbChangeChan,
	}
	lp.SetDirection(tview.FlexRow)
	lp.SetBorder(true)

	lp.SetBorderColor(base.BoarderDefaultColor)
	return &lp
}
