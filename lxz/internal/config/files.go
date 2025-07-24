/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 11:41
 */

package config

import (
	_ "embed"
	"github.com/adrg/xdg"
	"log/slog"
	"lxz/internal/config/data"
	"lxz/internal/slogs"
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

	// AppConfigFile tracks k9s config file.
	AppConfigFile string
)

// InitLogLoc initializes lxz logs location.
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
		appLogDir, err = xdg.StateFile(AppName)
		if err != nil {
			return err
		}
	}
	if err := data.EnsureFullPath(appLogDir, data.DefaultDirMod); err != nil {
		return err
	}
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

	AppConfigFile, err = xdg.ConfigFile(filepath.Join(AppName, data.MainConfigFile))
	if err != nil {
		return err
	}

	AppHotKeysFile = filepath.Join(AppConfigDir, "hotkeys.yaml")
	AppAliasesFile = filepath.Join(AppConfigDir, "aliases.yaml")
	AppPluginsFile = filepath.Join(AppConfigDir, "plugins.yaml")
	AppViewsFile = filepath.Join(AppConfigDir, "views.yaml")

	AppSkinsDir = filepath.Join(AppConfigDir, "skins")
	if e := data.EnsureFullPath(AppSkinsDir, data.DefaultDirMod); e != nil {
		slog.Warn("No skins dir detected", slogs.Error, e)
	}

	AppDumpsDir, err = xdg.StateFile(filepath.Join(AppName, "screen-dumps"))
	if err != nil {
		return err
	}

	_, err = xdg.DataFile(AppName)
	if err != nil {
		return err
	}

	return nil
}
