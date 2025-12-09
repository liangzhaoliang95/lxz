/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 10:50
 */

package view

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/drivers/redis_drivers"
	"github.com/liangzhaoliang95/lxz/internal/helper"
	"github.com/liangzhaoliang95/lxz/internal/view/base"
	"github.com/liangzhaoliang95/tview"
)

type RedisDbListView struct {
	*BaseFlex
	app               *App
	redisDbChangeChan chan redisDbChangeSubscribe // 用于订阅表变化的通道
	// 数据库
	dbNum       int                     // 数据库数量
	connConfig  *config.RedisConnConfig // redis连接配置
	dbKeyNumMap map[string]int64        // 用于存储数据库名称和对应的索引
	dbKeyMu     sync.Mutex              // 键数量映射的互斥锁，确保线程安全
	// UI组件
	dbListUI *tview.Table // 用于显示库列表
}

func (_this *RedisDbListView) selfFocus() {
	// 设置当前焦点为表格组件
	_this.app.UI.SetFocus(_this.dbListUI)
}

func (_this *RedisDbListView) _setDbKeyNumMap(db string, num int64) {
	// 安全设置数据库键数量映射
	_this.dbKeyMu.Lock()
	defer _this.dbKeyMu.Unlock()
	_this.dbKeyNumMap[db] = num
}

func (_this *RedisDbListView) Init(ctx context.Context) error {
	// 获取数据库连接配置
	// 初始化数据库连接
	rdbClient, err := redis_drivers.GetConnectOrInit(_this.connConfig, 0)
	if err != nil {
		return fmt.Errorf("failed to get redis connection: %w", err)
	}
	// 获取当前连接下的数据库列表
	dbNum, err := rdbClient.ListDB()
	if err != nil {
		return err
	}
	_this.dbNum = dbNum

	// 获取真实有key的数据库
	hasKeyDbs, _ := rdbClient.GetHasKeyDbNum()

	// 使用并发处理
	wg := sync.WaitGroup{}

	for i := 0; i < dbNum; i++ {
		wg.Add(1)
		go func(dbIndex int) {
			defer wg.Done()
			var dbSize int64
			// 如果当前数据库没有key，直接跳过
			if !helper.Contains(hasKeyDbs, dbIndex) {
				dbSize = 0
			} else {
				// 初始化数据库连接
				conn, err := redis_drivers.GetConnectOrInit(_this.connConfig, i)
				if err != nil {
					dbSize = 0
				} else {
					dbSize, _ = conn.GetDBKeyNum()
				}
			}

			_this._setDbKeyNumMap(fmt.Sprintf("%d", i), dbSize)
		}(i)
	}
	wg.Wait()

	// 初始化tree view
	_this.dbListUI = tview.NewTable()
	_this.dbListUI.SetBorder(false)
	_this.dbListUI.SetFixed(0, 0)
	_this.dbListUI.SetFixed(0, 1)
	_this.dbListUI.Select(1, 0) // 默认选中第一个数据库
	_this.dbListUI.SetSelectable(true, true)

	// 设置表格的选择模式
	_this.dbListUI.SetSelectionChangedFunc(func(row, column int) {
		slog.Info("Selection changed", "row", row, "col", column)
		if row < 1 || row >= _this.dbListUI.GetRowCount() {
			slog.Warn("Selection changed row is out of range", "row", row)
			return
		}
	})

	_this.dbListUI.SetSelectedFunc(func(row, col int) {
		slog.Info("Selected Node is a table node")
		// 启动表视图 发送表变化订阅
		dbChangeChan := redisDbChangeSubscribe{
			dbNum: row - 1,
		}
		_this.redisDbChangeChan <- dbChangeChan
	})
	// 设置头
	_this.dbListUI.SetCell(0, 0,
		tview.NewTableCell("DB").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetExpansion(1).
			SetSelectable(false))
	_this.dbListUI.SetCell(0, 1,
		tview.NewTableCell("Size").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetExpansion(1).
			SetSelectable(false))

	_this.AddItem(_this.dbListUI, 0, 1, true)

	return nil
}

func (_this *RedisDbListView) Start() {
	slog.Info("DatabaseDbTree Start", "redis", _this.connConfig.Name)
	for i := 0; i < _this.dbNum; i++ {
		_this.dbListUI.SetCell(i+1, 0,
			tview.NewTableCell(fmt.Sprintf("%d", i)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetSelectable(true))
		keySize := _this.dbKeyNumMap[fmt.Sprintf("%d", i)]
		_this.dbListUI.SetCell(i+1, 1,
			tview.NewTableCell(fmt.Sprintf("%d", keySize)).
				SetTextColor(helper.If[tcell.Color](keySize > 0, tcell.ColorRed, tcell.ColorGreen)).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetSelectable(true))
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
		dbNum:             0,
		dbKeyNumMap:       make(map[string]int64),
	}
	lp.SetBorder(true)
	lp.SetBorderColor(base.BoarderDefaultColor)
	return &lp
}
