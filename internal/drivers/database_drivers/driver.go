package database_drivers

import (
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"log/slog"
	"sync"
)

const (
	DefaultRowLimit = 100
)

var connMap sync.Map

type IDatabaseConn interface {
	GetDBConn() (*DatabaseConn, error)
	GetDbList() ([]string, error)
	GetTableList(dbName string) ([]string, error)
	GetRecords(database, table, where, sort string, offset, limit int) ([][]string, int, error)
	ExecuteQuery(query string) ([][]string, int, error)
}

// ---helpers

func _initDriver(cfg *config.DBConnection) (IDatabaseConn, error) {
	var dbDriver IDatabaseConn
	switch cfg.Provider {
	case config.DatabaseProviderMySQL:
		dbDriver = &MySQLDriver{
			DatabaseConn: &DatabaseConn{
				cfg:    cfg,
				dbConn: nil,
			},
		}
	default:
		return nil, fmt.Errorf("unsupported database provider: %s", cfg.Provider)
	}
	return dbDriver, nil
}

func GetConnect(cfg *config.DBConnection) (IDatabaseConn, error) {
	if db, exists := connMap.Load(cfg.GetUniqKey()); exists {
		return db.(IDatabaseConn), nil
	} else {
		return nil, fmt.Errorf("database connection not found for key: %s", cfg.GetUniqKey())
	}
}

func GetConnectOrInit(cfg *config.DBConnection) (IDatabaseConn, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in GetConnectOrInit", "error", r)
		}
	}()

	key := cfg.GetUniqKey()
	if db, exists := connMap.Load(key); exists {
		if conn, ok := db.(IDatabaseConn); ok {
			return conn, nil
		} else {
			return nil, fmt.Errorf("invalid type stored in connMap for key %s", key)
		}
	}

	iDriver, err := _initDriver(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database driver: %w", err)
	}
	connMap.Store(key, iDriver)

	return GetConnect(cfg)
}

func TestConnection(cfg *config.DBConnection) error {
	if cfg == nil {
		return fmt.Errorf("database connection configuration is nil")
	}
	iDriver, err := _initDriver(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database driver: %w", err)
	}
	conn, err := iDriver.GetDBConn()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB, err := conn.dbConn.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	conn.CloseConnect()
	return nil
}
