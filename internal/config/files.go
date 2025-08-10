/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 11:41
 */

package config

import (
	_ "embed"
	"github.com/adrg/xdg"
	"github.com/liangzhaoliang95/lxz/internal/config/data"
	"github.com/liangzhaoliang95/lxz/internal/slogs"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	AppName         = "lxz"
	LXZEnvLogsDir   = "LXZ_LOGS_DIR"
	LXZEnvConfigDir = "LXZ_CONFIG_DIR"
	LXZLogsFile     = "lxz.log"
)

var (
	//go:embed templates/benchmarks.yaml
	// benchmarkTpl tracks benchmark default config template
	benchmarkTpl []byte

	//go:embed templates/aliases.yaml
	// aliasesTpl tracks aliases default config template
	aliasesTpl []byte

	//go:embed templates/hotkeys.yaml
	// hotkeysTpl tracks hotkeys default config template
	hotkeysTpl []byte

	//go:embed templates/stock-skin.yaml
	// stockSkinTpl tracks stock skin template
	stockSkinTpl []byte
)

var (
	AppLogFile string

	// AppConfigDir tracks main lxz config home directory.
	AppConfigDir string

	// AppDumpsDir tracks screen dumps data directory.
	AppDumpsDir string

	// AppSkinsDir tracks skins data directory.
	AppSkinsDir string

	// AppViewsFile tracks custom views config file.
	AppViewsFile string

	// AppAliasesFile tracks aliases config file.
	AppAliasesFile string

	// AppPluginsFile tracks plugins config file.
	AppPluginsFile string

	// AppHotKeysFile tracks hotkeys config file.
	AppHotKeysFile string

	// AppConfigFile tracks LXZ config file.
	AppConfigFile string

	// AppDatabaseConfigFile tracks LXZ database config file.
	AppDatabaseConfigFile string

	// AppRedisConfigFile tracks LXZ redis config file.
	AppRedisConfigFile string
)

// InitLogLoc initializes LXZ logs location.
func InitLogLoc() error {
	var appLogDir string
	switch {
	case isEnvSet(LXZEnvLogsDir):
		appLogDir = os.Getenv(LXZEnvLogsDir)
	case isEnvSet(LXZEnvConfigDir):
		tmpDir, err := UserTmpDir()
		if err != nil {
			return err
		}
		appLogDir = tmpDir
	default:
		var err error
		// 使用xfg的StateFile方法获取配置目录
		appLogDir, err = xdg.StateFile(AppName)
		if err != nil {
			return err
		}
	}
	// 确保日志目录存在 不存在则创建
	if err := data.EnsureFullPath(appLogDir, data.DefaultDirMod); err != nil {
		return err
	}
	// 设置日志目录
	AppLogFile = filepath.Join(appLogDir, LXZLogsFile)

	return nil
}

// InitLocs initializes lxz artifacts locations.
func InitLocs() error {
	// 如果设置了环境变量 LXZ_CONFIG_DIR，则使用该目录作为配置目录
	if isEnvSet(LXZEnvConfigDir) {
		return initLxzEnvLocs()
	}

	// 默认使用 xdg 目录结构
	return initXDGLocs()
}

func initLxzEnvLocs() error {
	AppConfigDir = os.Getenv(LXZEnvConfigDir)
	if err := data.EnsureFullPath(AppConfigDir, data.DefaultDirMod); err != nil {
		return err
	}

	AppDumpsDir = filepath.Join(AppConfigDir, "screen-dumps")
	if err := data.EnsureFullPath(AppDumpsDir, data.DefaultDirMod); err != nil {
		slog.Warn("Unable to create screen-dumps dir", slogs.Dir, AppDumpsDir, slogs.Error, err)
	}

	AppSkinsDir = filepath.Join(AppConfigDir, "skins")
	if err := data.EnsureFullPath(AppSkinsDir, data.DefaultDirMod); err != nil {
		slog.Warn("Unable to create skins dir",
			slogs.Dir, AppSkinsDir,
			slogs.Error, err,
		)
	}
	AppConfigFile = filepath.Join(AppConfigDir, data.MainConfigFile)
	AppDatabaseConfigFile = filepath.Join(AppConfigDir, data.AppDatabaseConfigFile)
	AppRedisConfigFile = filepath.Join(AppConfigDir, data.AppRedisConfigFile)
	AppHotKeysFile = filepath.Join(AppConfigDir, "hotkeys.yaml")
	AppAliasesFile = filepath.Join(AppConfigDir, "aliases.yaml")
	AppPluginsFile = filepath.Join(AppConfigDir, "plugins.yaml")
	AppViewsFile = filepath.Join(AppConfigDir, "views.yaml")

	return nil
}

func initXDGLocs() error {
	var err error

	AppConfigDir, err = xdg.ConfigFile(AppName)
	if err != nil {
		return err
	}

	// 获取配置文件路径
	AppConfigFile, err = xdg.ConfigFile(filepath.Join(AppName, data.MainConfigFile))
	if err != nil {
		return err
	}

	// 快捷键配置文件路径
	AppHotKeysFile = filepath.Join(AppConfigDir, "hotkeys.yaml")
	// 别名配置文件路径
	AppAliasesFile = filepath.Join(AppConfigDir, "aliases.yaml")
	// 插件配置文件路径
	AppPluginsFile = filepath.Join(AppConfigDir, "plugins.yaml")
	// 视图配置文件路径
	AppViewsFile = filepath.Join(AppConfigDir, "views.yaml")
	// 皮肤配置文件夹路径
	AppSkinsDir = filepath.Join(AppConfigDir, "skins")
	if e := data.EnsureFullPath(AppSkinsDir, data.DefaultDirMod); e != nil {
		slog.Warn("No skins dir detected", slogs.Error, e)
	}
	// 堆栈截图保存路径
	AppDumpsDir, err = xdg.StateFile(filepath.Join(AppName, "screen-dumps"))
	if err != nil {
		return err
	}

	// --- 以下是具体应用的相关配置

	// 数据库管理配置文件路径
	AppDatabaseConfigFile = filepath.Join(AppConfigDir, data.AppDatabaseConfigFile)

	// Redis配置文件路径
	AppRedisConfigFile = filepath.Join(AppConfigDir, data.AppRedisConfigFile)

	// 检查配置文件夹
	_, err = xdg.DataFile(AppName)
	if err != nil {
		return err
	}

	return nil
}
