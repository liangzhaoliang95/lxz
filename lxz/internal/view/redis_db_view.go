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
	"lxz/internal/redis_drivers"
	"lxz/internal/view/base"
)

type RedisDbListView struct {
	*BaseFlex
	app               *App
	redisDbChangeChan chan redisDbChangeSubscribe // 用于订阅表变化的通道
	// 数据库
	connConfig *config.RedisConnConfig    // redis连接配置
	rdbClient  *redis_drivers.RedisClient // redis连接
	dbList     []string                   // 当前连接下的数据库列表
	selectDB   string                     // 当前选中的数据库
	// UI组件
	dbListUI *tview.Table // 用于显示库列表
}

func (_this *RedisDbListView) selfFocus() {
	// 设置当前焦点为表格组件
	_this.app.UI.SetFocus(_this)
}

func (_this *RedisDbListView) Init(ctx context.Context) error {
	// 获取数据库连接配置
	// 初始化数据库连接
	rdbClient, err := redis_drivers.GetConnectOrInit(_this.connConfig)
	if err != nil {
		return fmt.Errorf("failed to get redis connection: %w", err)
	}
	_this.rdbClient = rdbClient

	// 初始化tree view
	_this.dbListUI = tview.NewTable()
	//_this.dbListUI.SetBorder(false)
	//_this.dbListUI.SetSelectedFunc(func(row, col int) {
	//	slog.Info("Selected Node is a table node")
	//	// 启动表视图 发送表变化订阅
	//	dbChangeChan := redisDbChangeSubscribe{
	//		dbNum: row - 1,
	//	}
	//	_this.redisDbChangeChan <- dbChangeChan
	//})
	// 设置头
	//_this.dbListUI.SetCell(0, 0,
	//	tview.NewTableCell("DB NUM").
	//		SetTextColor(tcell.ColorYellow).
	//		SetAlign(tview.AlignCenter).
	//		SetExpansion(1).
	//		SetSelectable(false))
	//_this.AddItem(_this.dbListUI, 0, 1, false)

	// 获取当前连接下的数据库列表
	dbNum, err := _this.rdbClient.ListDB()
	if err != nil {
		return err
	}

	for i := 0; i < dbNum; i++ {
		_this.dbList = append(_this.dbList, fmt.Sprintf("db%d", i))

		//_this.dbListUI.SetCell(i+1, 0,
		//	tview.NewTableCell(fmt.Sprintf("%d", i)).
		//		SetTextColor(tcell.ColorYellow).
		//		SetAlign(tview.AlignLeft).
		//		SetExpansion(1).
		//		SetSelectable(true))
	}
	return nil
}

func (_this *RedisDbListView) Start() {
	// TODO 此处可能有Bug 会导致页面卡
	slog.Info("DatabaseDbTree Start", "redis", _this.connConfig.Name)
	for i := 0; i < len(_this.dbList); i++ {
		_this.dbListUI.SetCell(i+1, 0,
			tview.NewTableCell(_this.dbList[i]))
	}
}

func (_this *RedisDbListView) Stop() {

}

func NewRedisDbTree(
	a *App,
	dbCfg *config.RedisConnConfig,
	redisDbChangeChan chan redisDbChangeSubscribe,
) *RedisDbListView {
	var name = dbCfg.Name
	lp := RedisDbListView{
		BaseFlex:          NewBaseFlex(name),
		app:               a,
		connConfig:        dbCfg,
		redisDbChangeChan: redisDbChangeChan,
		dbList:            make([]string, 0),
	}
	lp.SetBorder(true)
	lp.SetBorderColor(base.BoarderDefaultColor)
	return &lp
}
