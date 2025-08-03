package database_drivers

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lxz/internal/config"
	"sync"
)

var connMap sync.Map

func InitConnect(cfg *config.DBConnection) error {
	var dialector gorm.Dialector
	switch cfg.Provider {
	case config.DatabaseProviderMySQL:
		if cfg.DBName != "" {
			cfg.URL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
				cfg.UserName,
				cfg.Password,
				cfg.Host,
				cfg.Port,
				cfg.DBName)
		} else {
			cfg.URL = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
				cfg.UserName,
				cfg.Password,
				cfg.Host,
				cfg.Port)
		}
		dialector = mysql.Open(cfg.URL)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	connMap.Store(cfg.GetUniqKey(), db)
	return nil
}

func GetConnect(cfg *config.DBConnection) (*gorm.DB, error) {
	if db, exists := connMap.Load(cfg.GetUniqKey()); exists {
		return db.(*gorm.DB), nil
	} else {
		return nil, fmt.Errorf("database connection not found for key: %s", cfg.GetUniqKey())
	}

}

func CloseConnect(cfg *config.DBConnection) error {
	db, err := GetConnect(cfg)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	connMap.Delete(cfg.GetUniqKey())
	return nil
}

func TestConnection(cfg *config.DBConnection) error {
	if cfg == nil {
		return fmt.Errorf("database connection configuration is nil")
	}
	err := InitConnect(cfg)
	if err != nil {
		return fmt.Errorf("failed to test connection: %w", err)
	}
	db, err := GetConnect(cfg)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	CloseConnect(cfg)
	return nil
}
