package view

import (
	"context"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/database_drivers"
	"lxz/internal/slogs"
	"lxz/internal/ui"
	"lxz/internal/ui/dialog"
	"strconv"
)

type DatabaseBrowser struct {
	*BaseFlex
	app       *App
	config    *config.DatabaseConfig
	connList  *tview.Table                    // connList 用于显示连接列表
	connMap   map[string]*config.DBConnection // connMap 用于存储连接信息的映射
	selectKey string                          // selectNum 用于记录选中的连接索引
}

func (_this *DatabaseBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 组件初始化
	// 连接列表
	_this.connList = tview.NewTable()
	_this.connList.SetBorder(false)
	_this.connList.SetBorders(false)
	_this.connList.SetTitle("🌍 Connections")
	_this.connList.SetBorderPadding(1, 1, 2, 2)
	_this.connList.SetSelectable(true, false)
	// 配置回车函数
	_this.connList.SetSelectedFunc(func(row, column int) {
		slog.Info("Selected connection", "row", row, "col", column)
		// 获取选中的连接信息
		if row < 1 || row >= _this.connList.GetRowCount() {
			slog.Warn("Selected row is out of range", "row", row)
			_this.app.UI.Flash().Warn("Please select a valid connection.")
			return
		}
		connName := _this.connList.GetCell(row, 0).Text
		slog.Info("Selected connection name", "name", connName)
		// 初始化数据库页面
	})
	// 设置表格的选择模式
	_this.connList.SetSelectionChangedFunc(func(row, column int) {
		slog.Info("Selection changed", "row", row, "col", column)
		if row < 1 || row >= _this.connList.GetRowCount() {
			slog.Warn("Selection changed row is out of range", "row", row)
			_this.app.UI.Flash().Warn("Please select a valid connection.")
			return
		}
	})

	// 设置布局 将连接列表居中
	_this.AddItem(tview.NewBox(), 3, 0, false)
	middlerFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	middlerFlex.AddItem(_this.connList, 0, 1, true)
	_this.AddItem(middlerFlex, 0, 1, true)

	return nil
}

func (_this *DatabaseBrowser) _initConfigTableHeader() {
	// 给列表设置列表头 name provider
	_this.connList.SetCell(
		0,
		0,
		tview.NewTableCell("Name").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connList.SetCell(
		0,
		1,
		tview.NewTableCell("Provider").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connList.SetCell(
		0,
		2,
		tview.NewTableCell("Host").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connList.SetCell(
		0,
		3,
		tview.NewTableCell("UserName").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connList.SetCell(
		0,
		4,
		tview.NewTableCell("Port").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
}

func (_this *DatabaseBrowser) _refreshTableData() {
	connMap := make(map[string]*config.DBConnection)
	for _, conn := range _this.config.DBConnections {
		connMap[conn.GetUniqKey()] = conn
	}
	_this.connMap = connMap

	_this.app.UI.QueueUpdateDraw(func() {
		// 清空表格
		_this.connList.Clear()
		_this._initConfigTableHeader()
		// 设置数据
		for i, connection := range _this.config.DBConnections {
			_this.connList.SetCell(
				i+1,
				0,
				tview.NewTableCell(connection.Name).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connList.SetCell(
				i+1,
				1,
				tview.NewTableCell(connection.Provider).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connList.SetCell(
				i+1,
				2,
				tview.NewTableCell(connection.Host).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connList.SetCell(
				i+1,
				3,
				tview.NewTableCell(connection.UserName).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connList.SetCell(
				i+1,
				4,
				tview.NewTableCell(strconv.FormatInt(connection.Port, 10)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
		}
	})
}

func (_this *DatabaseBrowser) Start() {
	// 设置数据
	_this._refreshTableData()

}

func (_this *DatabaseBrowser) Stop() {
	// 组件停止

}

// bindKeys
func (_this *DatabaseBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		tcell.KeyCtrlN: ui.NewKeyAction("New Connect", _this.createDatabaseConnectionModel, true),
		tcell.KeyCtrlD: ui.NewKeyAction(
			"Delete Connect",
			_this.deleteDatabaseConnectionModel,
			true,
		),
		tcell.KeyCtrlT: ui.NewKeyAction("Test Connect", _this.testConnect, true),
		ui.KeyE:        ui.NewKeyAction("Edit Connect", _this.createDatabaseConnectionModel, true),
		tcell.KeyEnter: ui.NewKeyAction("Connect", _this.startConnect, true),
		ui.KeyF:        ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
	})
}

// startConnect 处理连接事件
func (_this *DatabaseBrowser) startConnect(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info("Starting connection...")
	_this._getCurrentSelectKey()
	// 初始化main页面
	mainPage := NewDatabaseMainPage(_this.app, _this.connMap[_this.selectKey])
	loading := dialog.ShowLoadingDialog(appViewInstance.Content.Pages, "", appUiInstance.ForceDraw)
	_this.app.inject(mainPage, false)
	loading.Hide()
	return nil
}

func (_this *DatabaseBrowser) createDatabaseConnectionModel(evt *tcell.EventKey) *tcell.EventKey {
	var opts dialog.CreateDatabaseConnectionOpts
	if evt.Key() == tcell.KeyCtrlN {
		// 新建连接
		opts = dialog.CreateDatabaseConnectionOpts{
			Title:   "New Connection",
			Message: "",
			Ack: func(opts *config.DBConnection) bool {
				if opts.Name == "" {
					_this.app.UI.Flash().Warn("Connection name cannot be empty.")
					return false
				}
				key := opts.GetUniqKey()
				slog.Info("Creating new connection", "key", key)
				if _, exists := _this.connMap[key]; exists {
					slog.Warn("Connection already exists", "key", key)
					_this.app.UI.Flash().
						Warn("Connection already exists. Please choose a different name.")
					return false
				}
				if opts.Host == "" {
					_this.app.UI.Flash().Warn("Host cannot be empty.")
					return false
				}
				if opts.Port <= 0 {
					_this.app.UI.Flash().Warn("Port must be a positive integer.")
					return false
				}
				if opts.UserName == "" {
					_this.app.UI.Flash().Warn("Username cannot be empty.")
					return false
				}
				if opts.Password == "" {
					_this.app.UI.Flash().Warn("Password cannot be empty.")
					return false
				}

				_this.config.DBConnections = append(_this.config.DBConnections, opts)
				_this.config.Save(true)
				_this._refreshTableData()
				return true
			},
			Test: func(conn *config.DBConnection) bool {
				err := database_drivers.TestConnection(conn)
				if err != nil {
					slog.Error("Failed to test connection", slogs.Error, err)
					_this.app.UI.Flash().Warn(fmt.Sprintf("Failed to connect: %s", err.Error()))
					return false
				} else {
					slog.Info("Connection test successful", "connection", conn.GetUniqKey())
					_this.app.UI.Flash().Info("Connection test successful.")
					return true
				}
			},
			DBConnection: &config.DBConnection{
				Port: 3306,
			},
			Cancel: func() {},
		}
	}

	switch evt.Rune() {
	case 'e':
		// 编辑连接
		_this._getCurrentSelectKey()
		slog.Info("Editing connection", "selectKey", _this.selectKey)
		opts = dialog.CreateDatabaseConnectionOpts{
			Title:   "Edit Connection",
			Message: "",
			Ack: func(newConfig *config.DBConnection) bool {
				if newConfig.Name == "" {
					_this.app.UI.Flash().Warn("Connection name cannot be empty.")
					return false
				}
				key := newConfig.GetUniqKey()
				if key != _this.selectKey {
					if _, exists := _this.connMap[key]; exists {
						_this.app.UI.Flash().
							Warn("Connection already exists. Please choose a different name.")
						return false
					}
				}

				if newConfig.Host == "" {
					_this.app.UI.Flash().Warn("Host cannot be empty.")
					return false
				}
				if newConfig.Port <= 0 {
					_this.app.UI.Flash().Warn("Port must be a positive integer.")
					return false
				}
				if newConfig.UserName == "" {
					_this.app.UI.Flash().Warn("Username cannot be empty.")
					return false
				}
				if newConfig.Password == "" {
					_this.app.UI.Flash().Warn("Password cannot be empty.")
					return false
				}

				_this.config.Save(true)
				_this._refreshTableData()
				return true
			},
			Test: func(conn *config.DBConnection) bool {
				err := database_drivers.TestConnection(conn)
				if err != nil {
					slog.Error("Failed to test connection", slogs.Error, err)
					_this.app.UI.Flash().Warn(fmt.Sprintf("Failed to connect"))
					return false
				} else {
					slog.Info("Connection test successful", "connection", conn.GetUniqKey())
					_this.app.UI.Flash().Info("Connect success")
					return true
				}
			},
			DBConnection: _this.connMap[_this.selectKey],
			Cancel:       func() {},
		}
	}
	dialog.ShowCreateCreateDatabaseConnection(&config.Dialog{}, _this.app.Content.Pages, &opts)
	return nil
}

// deleteDatabaseConnectionModel 删除连接
func (_this *DatabaseBrowser) deleteDatabaseConnectionModel(evt *tcell.EventKey) *tcell.EventKey {
	_this._getCurrentSelectKey()
	opts := dialog.DeleteDatabaseConnectionOpts{
		Title:     "Delete Connection",
		Message:   "Are you sure you want to delete this connection?",
		SelectKey: _this.selectKey,
		Ack: func(key string) bool {
			// 删除连接

			// 删除选中的连接 根据索引的位置
			newConnections := make([]*config.DBConnection, 0, len(_this.config.DBConnections)-1)
			for i := 0; i < len(_this.config.DBConnections); i++ {
				item := _this.config.DBConnections[i]
				if item.GetUniqKey() == key {
					slog.Info("Deleting connection", "key", key)
					continue // 跳过删除的连接
				}
				newConnections = append(newConnections, item)
			}
			_this.config.DBConnections = newConnections

			_this.config.Save(true)
			_this._refreshTableData()
			return true
		},
		DBConnection: &config.DBConnection{},
		Cancel:       func() {},
	}
	dialog.ShowDeleteCreateDatabaseConnection(&config.Dialog{}, _this.app.Content.Pages, &opts)
	return nil
}

func (_this *DatabaseBrowser) _getCurrentSelectKey() {
	row, _ := _this.connList.GetSelection()
	currentSelectedName := _this.connList.GetCell(row, 0).Text
	currentSelectedProvider := _this.connList.GetCell(row, 1).Text
	_this.selectKey = fmt.Sprintf("%s@%s", currentSelectedProvider, currentSelectedName)
}

func (_this *DatabaseBrowser) testConnect(evt *tcell.EventKey) *tcell.EventKey {
	_this._getCurrentSelectKey()
	conn := _this.connMap[_this.selectKey]
	err := database_drivers.TestConnection(conn)
	if err != nil {
		slog.Error("Failed to test connection", slogs.Error, err)
		_this.app.UI.Flash().Warn(fmt.Sprintf("Failed to connect"))
	} else {
		slog.Info("Connection test successful", "connection", conn.GetUniqKey())
		_this.app.UI.Flash().Info("Connect success")
	}
	return nil
}

// ------------------helpers------------------

func loadConfiguration() (*config.DatabaseConfig, error) {
	slog.Info("🐶 lxz database browser loading configuration...")

	databaseCfg := config.NewDatabaseConfig()
	var errs error

	// 读取配置文件中的值,序列化到配置对象中 主要是将配置文件中的配置覆盖默认配置
	if err := databaseCfg.Load(config.AppDatabaseConfigFile, false); err != nil {
		errs = errors.Join(errs, err)
	}

	if err := databaseCfg.Save(true); err != nil {
		slog.Error("lxz config save failed", slogs.Error, err)
		errs = errors.Join(errs, err)
	} else {
		slog.Info("lxz config saved successfully", slogs.Path, config.AppDatabaseConfigFile)
	}

	return databaseCfg, errs
}

func _getDbConnectionKey(name, provider string) string {
	return fmt.Sprintf("%s@%s", provider, name)
}

func NewDatabaseBrowser(app *App) *DatabaseBrowser {
	databaseConfig, err := loadConfiguration()
	if err != nil {
		slog.Error("Failed to load database configuration", slogs.Error, err)
	} else {
		slog.Info(fmt.Sprintf("databaseConfig => %s", databaseConfig))
	}
	connMap := make(map[string]*config.DBConnection)
	for _, conn := range databaseConfig.DBConnections {
		connMap[conn.GetUniqKey()] = conn
	}
	db := DatabaseBrowser{
		BaseFlex: NewBaseFlex("DB Connections"),
		app:      app,
		config:   databaseConfig,
		connMap:  connMap,
	}

	db.SetIdentifier(ui.DB_BROWSER_ID)
	return &db
}
