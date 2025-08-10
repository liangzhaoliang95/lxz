package data

import (
	"os"

	"github.com/liangzhaoliang95/lxz/internal/config/json"
)

// JSONValidator validate yaml configurations.
var JSONValidator = json.NewValidator()

const (
	// DefaultDirMod default unix perms for LXZ directory.
	DefaultDirMod os.FileMode = 0744

	// DefaultFileMod default unix perms for LXZ files.
	DefaultFileMod os.FileMode = 0600

	// MainConfigFile track main configuration file.
	MainConfigFile        = "config.yaml"
	AppDatabaseConfigFile = "app_database_config.yaml" // 数据库应用的配置文件名称
	AppRedisConfigFile    = "app_redis_config.yaml"    // Redis应用的配置文件名称
)
