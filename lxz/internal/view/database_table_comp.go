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
	"lxz/internal/database_drivers"
	"lxz/internal/ui"
	"strings"
)

type DatabaseTableComponent struct {
	*ui.BaseFlex
	app       *App
	dbName    string // 数据库名称
	tableName string // 表名称
	dbCfg     *config.DBConnection
	dbConn    database_drivers.IDatabaseConn // 数据库连接接口
	// ui组件
	filterFlex  *tview.Flex       // 用于布局过滤条件输入框和标签
	filterLabel *tview.TextView   // 用于显示过滤条件标签
	filterInput *tview.InputField // 用于输入过滤条件
	dataTable   *tview.Table      // 用于显示表数据

}

func (_this *DatabaseTableComponent) Init(ctx context.Context) error {
	slog.Info("database table component init", "tableName", _this.tableName)
	// 初始化数据库连接 TODO 有bug
	iDatabaseConn, err := database_drivers.GetConnectOrInit(_this.dbCfg)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	slog.Info("123")
	_this.dbConn = iDatabaseConn

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
	_this.filterFlex.AddItem(_this.filterInput, 0, 5, true)

	_this.AddItem(_this.filterFlex, 3, 1, true)

	// 初始化表格
	_this.dataTable = tview.NewTable()
	_this.dataTable.SetBorders(true)
	_this.dataTable.SetBorder(true)
	_this.AddItem(_this.dataTable, 0, 7, false)

	return nil
}

func (_this *DatabaseTableComponent) Start() {
	// 初始化表格数据
	records, _, err := _this.dbConn.GetRecords(
		_this.dbName,
		_this.tableName,
		"",
		"",
		0,
		0,
	)
	if err != nil {
		_this.app.UI.Flash().
			Err(fmt.Errorf("failed to get records for table %s: %w", _this.tableName, err))
		return
	}
	// 渲染表格数据
	_this.SetTableData(records)

}

func (_this *DatabaseTableComponent) Stop() {

}

// --- data helpers ---

func (_this *DatabaseTableComponent) AddRows(rows [][]string) {
	for i, row := range rows {
		for j, cell := range row {
			tableCell := tview.NewTableCell(cell)
			tableCell.SetTextColor(tcell.ColorBlue)

			if cell == "EMPTY&" || cell == "NULL&" || cell == "DEFAULT&" {
				tableCell.SetText(strings.Replace(cell, "&", "", 1))
				tableCell.SetStyle(tcell.Style{})
				tableCell.SetReference(cell)
			}

			tableCell.SetSelectable(i > 0)
			tableCell.SetExpansion(1)

			_this.dataTable.SetCell(i, j, tableCell)
		}
	}
}

// SetTableData 设置表格数据
func (_this *DatabaseTableComponent) SetTableData(rows [][]string) {
	// 清空旧数据
	_this.dataTable.Clear()
	_this.AddRows(rows)
	_this.dataTable.Select(1, 0)
	_this.app.UI.QueueUpdateDraw(func() {})
}

func NewDatabaseTableComponent(
	a *App,
	dbName string,
	tableName string,
	dbCfg *config.DBConnection,
) *DatabaseTableComponent {
	var name = "Table View"
	lp := DatabaseTableComponent{
		BaseFlex:  ui.NewBaseFlex(name),
		app:       a,
		dbName:    dbName,
		tableName: tableName,
		dbCfg:     dbCfg,
	}
	lp.SetDirection(tview.FlexRow)
	lp.SetBorder(false)

	return &lp
}
