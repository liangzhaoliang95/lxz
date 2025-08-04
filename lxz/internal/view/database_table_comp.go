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
	"lxz/internal/view/base"
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

func (_this *DatabaseTableComponent) focusSearch() {
	// 设置当前焦点为搜索框
	_this.app.UI.SetFocus(_this.filterFlex)
	_this.filterFlex.SetBorderColor(base.ActiveBorderColor)
}

func (_this *DatabaseTableComponent) focusTable() {
	// 设置当前焦点为搜索框
	_this.app.UI.SetFocus(_this.dataTable)
}

func (_this *DatabaseTableComponent) Init(ctx context.Context) error {
	slog.Info("database table component init", "tableName", _this.tableName, "config", _this.dbCfg)
	// 初始化数据库连接
	iDatabaseConn, err := database_drivers.GetConnectOrInit(_this.dbCfg)
	slog.Info("get database connection", "dbName", _this.dbName, "tableName", _this.tableName, "CONN", iDatabaseConn)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
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
			// 初始化表格数据
			records, _, err := _this.dbConn.GetRecords(
				_this.dbName,
				_this.tableName,
				whereClause,
				"",
				0,
				0,
			)
			if err != nil {
				slog.Error("Failed to get records for table", "tableName", _this.tableName, "error", err)
				_this.app.UI.Flash().
					Err(fmt.Errorf("%w", err))
			} else {
				// 渲染表格数据
				_this.SetTableData(records)
			}
			// 焦点切换到表格
			_this.focusTable()
		case tcell.KeyEscape:

		}
	})
	_this.filterFlex.AddItem(_this.filterInput, 0, 5, true)

	_this.AddItem(_this.filterFlex, 3, 1, true)

	// 初始化表格
	_this.dataTable = tview.NewTable()
	_this.dataTable.SetBorders(true)
	_this.dataTable.SetBorder(false)
	_this.dataTable.SetSelectedStyle(tcell.Style{}.Background(tcell.ColorRed))
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
			if i == 0 {
				// 设置表头样式
				tableCell.SetTextColor(tcell.ColorYellow)
			}

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
	// 固定表头
	_this.dataTable.SetFixed(0, 0)
	// 选中第一行
	_this.dataTable.Select(1, 0)
	_this.app.UI.QueueUpdateDraw(func() {})
}

func NewDatabaseTableComponent(
	a *App,
	dbName string,
	tableName string,
	dbCfg *config.DBConnection,
) *DatabaseTableComponent {
	var name = ""
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
