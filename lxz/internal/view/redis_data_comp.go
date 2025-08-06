/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 10:50
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/redis_drivers"
	"lxz/internal/view/base"
)

type RedisDataComponent struct {
	*BaseFlex
	app             *App
	dbNum           int // 数据库名称
	redisConnConfig *config.RedisConnConfig
	rdbClient       *redis_drivers.RedisClient
	// ui组件
	filterFlex  *tview.Flex       // 用于布局过滤条件输入框和标签
	filterLabel *tview.TextView   // 用于显示过滤条件标签
	filterInput *tview.InputField // 用于输入过滤条件
	dataTable   *tview.Table      // 用于显示表数据

}

func (_this *RedisDataComponent) focusSearch() {
	// 设置当前焦点为搜索框
	_this.app.UI.SetFocus(_this.filterFlex)
	_this.filterFlex.SetBorderColor(base.ActiveBorderColor)
}

func (_this *RedisDataComponent) focusTable() {
	// 设置当前焦点为搜索框
	_this.app.UI.SetFocus(_this.dataTable)
}

func (_this *RedisDataComponent) Init(ctx context.Context) error {
	slog.Info("database table component init", "dbNum", _this.dbNum)
	// 初始化数据库连接
	iRedisConn, err := redis_drivers.GetConnectOrInit(_this.redisConnConfig)

	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_this.rdbClient = iRedisConn

	// 初始化filterFlex
	_this.filterFlex = tview.NewFlex()
	_this.filterFlex.SetDirection(tview.FlexColumn)
	_this.filterFlex.SetBorder(true)
	_this.filterFlex.SetBorderPadding(0, 0, 1, 1)

	// 初始化filterLabel
	_this.filterLabel = tview.NewTextView()
	_this.filterLabel.SetText("WHERE")
	_this.filterLabel.SetTextAlign(tview.AlignCenter)
	_this.filterLabel.SetTextColor(tcell.ColorGreen)
	_this.filterLabel.SetBorderPadding(0, 0, 0, 0)
	_this.filterFlex.AddItem(_this.filterLabel, 6, 1, false)

	// 初始化filterInput
	_this.filterInput = tview.NewInputField()
	_this.filterInput.SetPlaceholder("Enter a WHERE clause to filter the results")
	_this.filterInput.SetFieldBackgroundColor(tcell.ColorBlack)
	_this.filterInput.SetFieldTextColor(tcell.ColorRed)
	_this.filterInput.SetFocusFunc(func() {
		_this.filterFlex.SetBorderColor(base.ActiveBorderColor)
	})
	_this.filterInput.SetBlurFunc(func() {
		_this.filterFlex.SetBorderColor(base.InactiveBorderColor)
	})
	_this.filterInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			whereClause := ""
			if _this.filterInput.GetText() != "" {
				whereClause = fmt.Sprintf("WHERE %s", _this.filterInput.GetText())
			}
			slog.Info("filter input done", "where", whereClause)
			// 初始化表格数据
			records, err := _this.rdbClient.GetRecords(
				_this.dbNum,
				"",
			)
			if err != nil {
				slog.Error(
					"Failed to get records for db",
					"dbNum",
					_this.dbNum,
					"error",
					err,
				)
				_this.app.UI.Flash().
					Err(fmt.Errorf("%w", err))

				// 焦点切换到表格
				_this.focusSearch()
			} else {
				_this.SetTableData([][]string{
					records,
				})

				// 焦点切换到表格
				_this.focusTable()
			}
		case tcell.KeyEscape:

		}
	})
	_this.filterFlex.AddItem(_this.filterInput, 0, 5, true)
	_this.AddItem(_this.filterFlex, 3, 1, false)

	// 初始化表格
	_this.dataTable = tview.NewTable()
	_this.dataTable.SetBorders(true)
	_this.dataTable.SetBorder(false)
	_this.dataTable.SetSeparator(tview.Borders.Vertical)
	_this.dataTable.SetSelectedStyle(
		tcell.StyleDefault.Background(tcell.ColorRed).
			Foreground(tview.Styles.ContrastSecondaryTextColor),
	)
	_this.dataTable.SetSelectable(true, false)
	_this.dataTable.SetFixed(1, 0)
	_this.AddItem(_this.dataTable, 0, 7, true)
	return nil
}

func (_this *RedisDataComponent) Start() {

	// 初始化表格数据
	records, err := _this.rdbClient.GetRecords(
		_this.dbNum,
		"",
	)
	if err != nil {
		_this.app.UI.Flash().
			Err(fmt.Errorf("failed to get records for db %d: %w", _this.dbNum, err))
		return
	}

	// 渲染表格数据 不知道为什么，就是需要刷两遍才能定位到第一行
	_this.SetTableData([][]string{records})

}

func (_this *RedisDataComponent) Stop() {

}

// --- data helpers ---

// SetTableData 设置表格数据
func (_this *RedisDataComponent) SetTableData(rows [][]string) {
	// 清空旧数据
	_this.dataTable.Clear()
	TableAddRows(_this.dataTable, rows)
	_this.dataTable.Select(1, 1)
}

func NewRedisDataComponent(
	a *App,
	dbNum int,
	redisConnConfig *config.RedisConnConfig,
) *RedisDataComponent {
	var name = ""
	lp := RedisDataComponent{
		BaseFlex:        NewBaseFlex(name),
		app:             a,
		dbNum:           dbNum,
		redisConnConfig: redisConnConfig,
	}
	lp.SetDirection(tview.FlexRow)
	lp.SetBorder(false)

	return &lp
}
