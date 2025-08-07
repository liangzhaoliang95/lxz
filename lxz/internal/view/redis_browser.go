package view

import (
	"context"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/drivers/redis_drivers"
	"lxz/internal/slogs"
	"lxz/internal/ui"
	"lxz/internal/ui/dialog"
	"strconv"
)

type RedisBrowser struct {
	*BaseFlex
	app    *App
	config *config.RedisConfig

	connMap   map[string]*config.RedisConnConfig // connMap 用于存储连接信息的映射
	selectKey string                             // selectNum 用于记录选中的连接索引

	// UI组件
	connListTableUI *tview.Table // connListTableUI 用于显示连接列表
}

func (_this *RedisBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 组件初始化
	// 连接列表
	_this.connListTableUI = tview.NewTable()
	_this.connListTableUI.SetBorder(false)
	_this.connListTableUI.SetBorders(false)
	_this.connListTableUI.SetTitle("🌍 Connections")
	_this.connListTableUI.SetBorderPadding(1, 1, 2, 2)
	_this.connListTableUI.SetSelectable(true, false)
	// 配置回车函数
	_this.connListTableUI.SetSelectedFunc(func(row, column int) {
		slog.Info("Selected connection", "row", row, "col", column)
		// 获取选中的连接信息
		if row < 1 || row >= _this.connListTableUI.GetRowCount() {
			slog.Warn("Selected row is out of range", "row", row)
			return
		}
		connName := _this.connListTableUI.GetCell(row, 0).Text
		slog.Info("Selected connection name", "name", connName)
		// 初始化数据库页面
	})
	// 设置表格的选择模式
	_this.connListTableUI.SetSelectionChangedFunc(func(row, column int) {
		slog.Info("Selection changed", "row", row, "col", column)
		if row < 1 || row >= _this.connListTableUI.GetRowCount() {
			slog.Warn("Selection changed row is out of range", "row", row)
			return
		}
	})

	// 设置布局 将连接列表居中
	_this.AddItem(tview.NewBox(), 3, 0, false)
	middlerFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	middlerFlex.AddItem(_this.connListTableUI, 0, 1, true)
	_this.AddItem(middlerFlex, 0, 1, true)

	return nil
}

func (_this *RedisBrowser) _initRedisConfigTableHeader() {
	// 给列表设置列表头 name provider
	_this.connListTableUI.SetCell(
		0,
		0,
		tview.NewTableCell("Name").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connListTableUI.SetCell(
		0,
		1,
		tview.NewTableCell("Host").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connListTableUI.SetCell(
		0,
		2,
		tview.NewTableCell("UserName").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
	_this.connListTableUI.SetCell(
		0,
		3,
		tview.NewTableCell("Port").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetExpansion(1).
			SetSelectable(false),
	)
}

func (_this *RedisBrowser) _refreshTableData() {
	connMap := make(map[string]*config.RedisConnConfig)
	for _, conn := range _this.config.RedisConnConfig {
		connMap[conn.Name] = conn
	}
	_this.connMap = connMap

	_this.app.UI.QueueUpdateDraw(func() {
		// 清空表格
		_this.connListTableUI.Clear()
		_this._initRedisConfigTableHeader()
		// 设置数据
		for i, connection := range _this.config.RedisConnConfig {
			_this.connListTableUI.SetCell(
				i+1,
				0,
				tview.NewTableCell(connection.Name).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connListTableUI.SetCell(
				i+1,
				1,
				tview.NewTableCell(connection.Host).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connListTableUI.SetCell(
				i+1,
				2,
				tview.NewTableCell(connection.UserName).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
			_this.connListTableUI.SetCell(
				i+1,
				3,
				tview.NewTableCell(strconv.FormatInt(connection.Port, 10)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetExpansion(1),
			)
		}
	})
}

func (_this *RedisBrowser) Start() {
	// 设置数据
	_this._refreshTableData()

}

func (_this *RedisBrowser) Stop() {
	// 组件停止

}

// bindKeys
func (_this *RedisBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		tcell.KeyCtrlN: ui.NewKeyAction("New Connect", _this.createRedisConfigModel, true),
		tcell.KeyCtrlD: ui.NewKeyAction(
			"Delete Connect",
			_this.deleteRedisConnectionModel,
			true,
		),
		tcell.KeyCtrlT: ui.NewKeyAction("Test Connect", _this.testConnect, true),
		ui.KeyE:        ui.NewKeyAction("Edit Connect", _this.createRedisConfigModel, true),
		tcell.KeyEnter: ui.NewKeyAction("Connect", _this.startConnect, true),
		ui.KeyF:        ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
	})
}

// startConnect 处理连接事件
func (_this *RedisBrowser) startConnect(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info("Starting connection redis...")
	_this._getCurrentSelectKey()
	// 初始化main页面
	loading := dialog.ShowLoadingDialog(appViewInstance.Content.Pages, "", appUiInstance.ForceDraw)

	mainPage := NewRedisMainPage(_this.app, _this.connMap[_this.selectKey])
	_this.app.inject(mainPage, false)
	loading.Hide()
	return nil
}

func (_this *RedisBrowser) createRedisConfigModel(evt *tcell.EventKey) *tcell.EventKey {
	var opts dialog.CreateRedisConnectionOpts
	if evt.Key() == tcell.KeyCtrlN {
		// 新建连接
		opts = dialog.CreateRedisConnectionOpts{
			Title:   "New Connection",
			Message: "",
			Ack: func(opts *config.RedisConnConfig) bool {
				if opts.Name == "" {
					_this.app.UI.Flash().Warn("Connection name cannot be empty.")
					return false
				}
				if _, exists := _this.connMap[opts.Name]; exists {
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

				_this.config.RedisConnConfig = append(_this.config.RedisConnConfig, opts)
				err := _this.config.Save(true)
				if err != nil {
					_this.app.UI.Flash().Warn("Failed to save configuration: " + err.Error())
					return false
				}
				_this._refreshTableData()
				return true
			},
			Test: func(conn *config.RedisConnConfig) bool {
				err := redis_drivers.TestConnection(conn)
				if err != nil {
					slog.Error("Failed to test connection", slogs.Error, err)
					_this.app.UI.Flash().Warn(fmt.Sprintf("Failed to connect: %s", err.Error()))
					return false
				} else {
					_this.app.UI.Flash().Info("Connection successful.")
					return true
				}
			},
			Config: &config.RedisConnConfig{
				Port: 6379,
			},
			Cancel: func() {},
		}
	}

	switch evt.Rune() {
	case 'e':
		// 编辑连接
		_this._getCurrentSelectKey()
		slog.Info("Editing connection", "selectKey", _this.selectKey)
		opts = dialog.CreateRedisConnectionOpts{
			Title:   "Edit Connection",
			Message: "",
			Ack: func(newConfig *config.RedisConnConfig) bool {
				if newConfig.Name == "" {
					_this.app.UI.Flash().Warn("Connection name cannot be empty.")
					return false
				}
				key := newConfig.Name
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

				err := _this.config.Save(true)
				if err != nil {
					_this.app.UI.Flash().Warn("Failed to save configuration: " + err.Error())
					return false
				}
				_this._refreshTableData()
				return true
			},
			Test: func(conn *config.RedisConnConfig) bool {
				err := redis_drivers.TestConnection(conn)
				if err != nil {
					_this.app.UI.Flash().Warn(fmt.Sprintf("Failed to connect"))
					return false
				} else {
					_this.app.UI.Flash().Info("Connect success")
					return true
				}
			},
			Config: _this.connMap[_this.selectKey],
			Cancel: func() {},
		}
	}
	dialog.ShowCreateRedisConnection(&config.Dialog{}, _this.app.Content.Pages, &opts)
	return nil
}

// deleteRedisConnectionModel 删除连接
func (_this *RedisBrowser) deleteRedisConnectionModel(evt *tcell.EventKey) *tcell.EventKey {
	_this._getCurrentSelectKey()
	opts := dialog.DeleteRedisConnectionOpts{
		Title:     "Delete Connection",
		Message:   "Are you sure you want to delete this connection?",
		SelectKey: _this.selectKey,
		Ack: func(key string) bool {
			// 删除连接

			// 删除选中的连接 根据索引的位置
			newConnections := make(
				[]*config.RedisConnConfig,
				0,
				len(_this.config.RedisConnConfig)-1,
			)
			for i := 0; i < len(_this.config.RedisConnConfig); i++ {
				item := _this.config.RedisConnConfig[i]
				if item.Name == key {
					continue // 跳过删除的连接
				}
				newConnections = append(newConnections, item)
			}
			_this.config.RedisConnConfig = newConnections

			err := _this.config.Save(true)
			if err != nil {
				_this.app.UI.Flash().Warn("Failed to save configuration: " + err.Error())
				return false
			} else {
				slog.Info("Connection deleted successfully", "key", key)
				_this.app.UI.Flash().Info("Connection deleted successfully.")
			}
			_this._refreshTableData()
			return true
		},
		Config: &config.RedisConnConfig{},
		Cancel: func() {},
	}
	dialog.ShowDeleteRedisConnection(&config.Dialog{}, _this.app.Content.Pages, &opts)
	return nil
}

func (_this *RedisBrowser) _getCurrentSelectKey() {
	row, _ := _this.connListTableUI.GetSelection()
	currentSelectedName := _this.connListTableUI.GetCell(row, 0).Text
	_this.selectKey = fmt.Sprintf("%s", currentSelectedName)
}

func (_this *RedisBrowser) testConnect(evt *tcell.EventKey) *tcell.EventKey {
	_this._getCurrentSelectKey()
	conn := _this.connMap[_this.selectKey]
	err := redis_drivers.TestConnection(conn)
	if err != nil {
		_this.app.UI.Flash().Warn(fmt.Sprintf("Failed to connect"))
	} else {
		_this.app.UI.Flash().Info("Connect success")
	}
	return nil
}

// ------------------helpers------------------

func loadRedisConfiguration() (*config.RedisConfig, error) {
	slog.Info("🐶 lxz redis browser loading configuration...")

	redisConfig := config.NewRedisConfig()
	var errs error

	// 读取配置文件中的值,序列化到配置对象中 主要是将配置文件中的配置覆盖默认配置
	if err := redisConfig.Load(config.AppRedisConfigFile, false); err != nil {
		errs = errors.Join(errs, err)
	}

	if err := redisConfig.Save(true); err != nil {
		slog.Error("lxz redis config save failed", slogs.Error, err)
		errs = errors.Join(errs, err)
	} else {
		slog.Info("lxz redis config saved successfully", slogs.Path, config.AppRedisConfigFile)
	}

	return redisConfig, errs
}

func NewRedisBrowser(app *App) *RedisBrowser {
	databaseConfig, err := loadRedisConfiguration()
	if err != nil {
		slog.Error("Failed to load redis configuration", slogs.Error, err)
	} else {
		slog.Info(fmt.Sprintf("redis config => %s", databaseConfig))
	}
	connMap := make(map[string]*config.RedisConnConfig)
	for _, conn := range databaseConfig.RedisConnConfig {
		connMap[conn.Name] = conn
	}
	db := RedisBrowser{
		BaseFlex: NewBaseFlex("Redis Connections"),
		app:      app,
		config:   databaseConfig,
		connMap:  connMap,
	}

	db.SetIdentifier(ui.REDIS_BROWSER_ID)
	return &db
}
