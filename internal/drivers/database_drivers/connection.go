/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 11:19
 */

package database_drivers

import (
	"fmt"

	"github.com/liangzhaoliang95/lxz/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConnIns struct {
	IDatabaseConn IDatabaseConn
}

type DatabaseConn struct {
	cfg    *config.DBConnection
	dbConn *gorm.DB
}

func (_this *DatabaseConn) InitConnect() error {
	if _this.dbConn != nil {
		return nil
	}
	var dialector gorm.Dialector
	if _this.cfg.Provider == config.DatabaseProviderMySQL {
		if _this.cfg.DBName != "" {
			_this.cfg.URL = fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
				_this.cfg.UserName,
				_this.cfg.Password,
				_this.cfg.Host,
				_this.cfg.Port,
				_this.cfg.DBName,
			)
		} else {
			_this.cfg.URL = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
				_this.cfg.UserName,
				_this.cfg.Password,
				_this.cfg.Host,
				_this.cfg.Port)
		}
		dialector = mysql.Open(_this.cfg.URL)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	_this.dbConn = db
	return nil
}

func (_this *DatabaseConn) CloseConnect() error {
	if _this.dbConn == nil {
		return nil
	}

	sqlDB, err := _this.dbConn.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	_this.dbConn = nil
	connMap.Delete(_this.cfg.GetUniqKey())
	return nil
}
