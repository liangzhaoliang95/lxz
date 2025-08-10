/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 10:50
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/drivers/database_drivers"
	"lxz/internal/ui"
	"lxz/internal/view/base"
)

type DatabaseQueryView struct {
	*BaseFlex
	app    *App
	dbCfg  *config.DBConnection // 数据库连接配置
	dbConn database_drivers.IDatabaseConn

	filterInput *tview.InputField // 用于输入过滤条件
	dataTable   *tview.Table      // 用于显示表数据

}

func (_this *DatabaseQueryView) TabFocusChange(event *tcell.EventKey) *tcell.EventKey {
	if _this.app.UI.GetFocus() == _this.filterInput {
		_this.focusTable()
	} else {
		_this.focusSearch()
	}
	return nil
}

func (_this *DatabaseQueryView) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		ui.KeySlash:     ui.NewKeyAction("Search", _this.ToggleSearch, true),
		tcell.KeyEscape: ui.NewKeyAction("Last Page", _this.EmptyKeyEvent, true),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
	})
}

// ToggleSearch
func (_this *DatabaseQueryView) ToggleSearch(evt *tcell.EventKey) *tcell.EventKey {
	if ui.IsInputPrimitive(appUiInstance.GetFocus()) {
		slog.Info("Search input is already focused, ignoring toggle search event")
		return evt
	}
	// 切换焦点到搜索框
	if _this.app.UI.GetFocus() == _this.filterInput {
		_this.focusTable()
	} else {
		_this.app.UI.SetFocus(_this.filterInput)
	}
	return nil
}

// SetTableData 设置表格数据
func (_this *DatabaseQueryView) SetTableData(rows [][]string) {
	// 清空旧数据
	_this.dataTable.Clear()
	TableAddRows(_this.dataTable, rows)
	_this.dataTable.Select(1, 1)
}

func (_this *DatabaseQueryView) focusTable() {
	// 设置当前焦点为搜索框
	_this.app.UI.SetFocus(_this.dataTable)
}

// focusSearch
func (_this *DatabaseQueryView) focusSearch() {
	// 切换焦点到搜索框
	_this.app.UI.SetFocus(_this.filterInput)
}

func (_this *DatabaseQueryView) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 初始化数据库连接
	iDatabaseConn, err := database_drivers.GetConnectOrInit(_this.dbCfg)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_this.dbConn = iDatabaseConn

	// 初始化组件
	_this.filterInput = tview.NewInputField()
	_this.filterInput.SetLabel("Filter: ")
	_this.filterInput.SetBorderPadding(0, 0, 1, 1)
	_this.filterInput.SetBorder(true)
	_this.filterInput.SetFocusFunc(func() {
		_this.filterInput.SetBorderColor(base.ActiveBorderColor)
	})
	_this.filterInput.SetBlurFunc(func() {
		_this.filterInput.SetBorderColor(base.InactiveBorderColor)
	})
	_this.filterInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			whereClause := ""
			if _this.filterInput.GetText() != "" {
				whereClause = fmt.Sprintf("%s", _this.filterInput.GetText())
			}
			// 初始化表格数据
			//loading := dialog.ShowLoadingDialog(_this.app.Content.Pages, "", _this.app.UI.Draw)
			records, _, err := _this.dbConn.ExecuteQuery(
				whereClause,
			)
			slog.Info("Executing query with where clause: " + whereClause)
			if err != nil {
				_this.app.UI.Flash().
					Err(fmt.Errorf("%w", err))
				_this.focusSearch()
			} else {

				_this.SetTableData(records)
				_this.focusTable()
			}
			//loading.Hide()
		case tcell.KeyEscape:

		}
	})

	_this.dataTable = tview.NewTable()
	_this.dataTable.SetBorders(true)
	_this.dataTable.SetBorder(false)
	_this.dataTable.SetSeparator(tview.Borders.Vertical)
	_this.dataTable.SetFocusFunc(func() {
		_this.dataTable.SetBorderColor(base.ActiveBorderColor)
	})
	_this.dataTable.SetBlurFunc(func() {
		_this.dataTable.SetBorderColor(base.InactiveBorderColor)
	})
	_this.dataTable.SetSelectedStyle(
		tcell.StyleDefault.Background(tcell.ColorRed).
			Foreground(tview.Styles.ContrastSecondaryTextColor),
	)
	_this.dataTable.SetSelectable(true, false)
	_this.dataTable.SetFixed(1, 0)

	_this.AddItem(_this.filterInput, 3, 0, true)
	_this.AddItem(_this.dataTable, 0, 1, false)
	return nil
}

func (_this *DatabaseQueryView) Start() {

	_this.SetTableData([][]string{
		{"No data", "Please enter a query"},
	})

}

func (_this *DatabaseQueryView) Stop() {

}

// --- data helpers ---

func NewDatabaseQueryView(
	a *App,
	dbCfg *config.DBConnection,
) *DatabaseQueryView {
	var name = "Manul Query"
	lp := DatabaseQueryView{
		BaseFlex: NewBaseFlex(name),
		app:      a,
		dbCfg:    dbCfg,
	}
	lp.SetDirection(tview.FlexRow)
	lp.SetBorder(true)

	lp.SetBorderColor(base.BoarderDefaultColor)
	return &lp
}
