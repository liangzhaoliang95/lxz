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

type DatabaseTableView struct {
	*BaseFlex
	tableChangeChan chan tableChangeSubscribe // 表格数据变更订阅通道
	app             *App
	dbCfg           *config.DBConnection               // 数据库连接配置
	tablePages      *tview.Pages                       // 表格数据页面容器
	tableComponents map[string]*DatabaseTableComponent // 表格数据页面
	currentPageKey  string                             // 当前页面的键，用于切换页面
}

func (_this *DatabaseTableView) selfFocus() {

	comp := _this.tableComponents[_this.currentPageKey]
	if comp == nil {
		_this.app.UI.SetFocus(_this)
	} else {
		// 设置当前焦点为表格组件
		_this.app.UI.SetFocus(comp.dataTable)
	}

}

func (_this *DatabaseTableView) LunchPage(dbName, tableName string) error {
	slog.Info("Launching page for table", "tableName", tableName, "dbName", dbName)

	loading := dialog.ShowLoadingDialog(
		appViewInstance.Content.Pages,
		"",
		appUiInstance.ForceDraw,
	)
	pageKey := fmt.Sprintf("%s-%s", dbName, tableName)
	if _, ok := _this.tableComponents[pageKey]; !ok {
		// 当前页不存在，创建一个新的表格组件
		tableComponent := NewDatabaseTableComponent(_this.app, dbName, tableName, _this.dbCfg)
		slog.Info("Initializing table component", "pageKey", pageKey)
		err := tableComponent.Init(context.Background())
		if err != nil {
			slog.Error("Failed to initialize table component", "err", err)
			return err
		}
		_this.tableComponents[pageKey] = tableComponent
		_this.tablePages.AddPage(pageKey, tableComponent, true, true)
	}
	// 当前页已存在，切换到该页面
	_this.currentPageKey = pageKey
	_this.tablePages.SwitchToPage(pageKey)

	comp := _this.tableComponents[pageKey]
	comp.Start()

	loading.Hide()
	_this.selfFocus()
	appUiInstance.ForceDraw()
	return nil

}

func (_this *DatabaseTableView) Init(ctx context.Context) error {
	_this.AddItem(_this.tablePages, 0, 1, true)
	return nil
}

func (_this *DatabaseTableView) Start() {
	go func() {
		slog.Info("Starting DatabaseTableView...")
		// 监听changeChan通道，处理表格数据变更
		for change := range _this.tableChangeChan {
			slog.Info(
				"Received table change notification",
				"dbName",
				change.dbName,
				"tableName",
				change.tableName,
			)
			if change.dbName == "" || change.tableName == "" {
				continue // 无效的变更通知
			}
			if err := _this.LunchPage(change.dbName, change.tableName); err != nil {
				fmt.Printf(
					"Error launching page for %s.%s: %v\n",
					change.dbName,
					change.tableName,
					err,
				)
			}
		}
		slog.Info("886")
	}()
}

func (_this *DatabaseTableView) Stop() {

}

// --- data helpers ---

func NewDatabaseTableView(
	a *App,
	dbCfg *config.DBConnection,
	tableChangeChan chan tableChangeSubscribe,
) *DatabaseTableView {
	var name = "Table View"
	lp := DatabaseTableView{
		BaseFlex:        NewBaseFlex(name),
		app:             a,
		tablePages:      tview.NewPages(),
		tableComponents: make(map[string]*DatabaseTableComponent),
		dbCfg:           dbCfg,
		tableChangeChan: tableChangeChan,
	}
	lp.SetDirection(tview.FlexRow)
	lp.SetBorder(true)

	lp.SetBorderColor(base.BoarderDefaultColor)
	return &lp
}
