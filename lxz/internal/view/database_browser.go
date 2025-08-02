package view

import (
	"context"
	"errors"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/slogs"
	"lxz/internal/ui"
)

type DatabaseBrowser struct {
	*ui.BaseFlex
	app    *App
	config *config.DatabaseConfig
}

func (d *DatabaseBrowser) Init(ctx context.Context) error {
	// 组件初始化
	return nil
}

func (d *DatabaseBrowser) Start() {
	// 组件启动 拉起启动页 inject组件

}

func (d *DatabaseBrowser) Stop() {
	// 组件停止

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

	if err := databaseCfg.Save(false); err != nil {
		slog.Error("lxz config save failed", slogs.Error, err)
		errs = errors.Join(errs, err)
	}

	return databaseCfg, errs
}

func NewDatabaseBrowser(app *App) *DatabaseBrowser {
	databaseConfig, err := loadConfiguration()
	if err != nil {
		slog.Error("Failed to load database configuration", slogs.Error, err)
	}
	db := DatabaseBrowser{
		BaseFlex: ui.NewBaseFlex("DatabaseBrowser"),
		app:      app,
		config:   databaseConfig,
	}
	db.SetBorder(true)
	db.SetTitle("Database Browser")
	db.SetTitleAlign(tview.AlignCenter)

	return &db
}
