/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 11:38
 */

package cmd

import (
	"errors"
	"fmt"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"log/slog"
	"lxz/internal/color"
	"lxz/internal/config"
	"lxz/internal/config/data"
	"lxz/internal/slogs"
	"lxz/internal/view"
	"os"
	"runtime/debug"
	"time"

	"github.com/mattn/go-colorable"
)

const (
	appName      = config.AppName
	shortAppDesc = "A graphical CLI for lxz devops"
	longAppDesc  = "lxz is a CLI to view and manage devops"
)

var (
	version, commit, date = "dev", "dev", "0"
	lxzFlags              *config.Flags

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Long:  longAppDesc,
		RunE:  run,
	}

	out = colorable.NewColorableStdout()
)

type flagError struct{ err error }

func (e flagError) Error() string { return e.err.Error() }

func init() {
	if err := config.InitLogLoc(); err != nil {
		fmt.Printf("Fail to init LXZ logs location %s\n", err)
	}

	rootCmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return flagError{err: err}
	})

	rootCmd.AddCommand(versionCmd(), infoCmd())
	initLXZFlags()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if !errors.As(err, &flagError{}) {
			panic(err)
		}
	}
}

func run(*cobra.Command, []string) error {
	if err := config.InitLocs(); err != nil {
		return err
	}
	logFile, err := os.OpenFile(
		*lxzFlags.LogFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		data.DefaultFileMod,
	)
	if err != nil {
		return fmt.Errorf("log file %q init failed: %w", *lxzFlags.LogFile, err)
	}
	defer func() {
		if logFile != nil {
			_ = logFile.Close()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			slog.Error("Boom!! LXZ init failed", slogs.Error, err)
			slog.Error("", slogs.Stack, string(debug.Stack()))
			printLogo(color.Red)
			fmt.Printf("%s", color.Colorize("Boom!! ", color.Red))
			fmt.Printf("%v.\n", err)
		}
	}()

	slog.SetDefault(slog.New(tint.NewHandler(logFile, &tint.Options{
		Level:      parseLevel(*lxzFlags.LogLevel),
		TimeFormat: time.Kitchen,
	})))

	cfg, err := loadConfiguration()
	if err != nil {
		slog.Warn("Fail to load global/context configuration", slogs.Error, err)
	}
	slog.Info("🐶 lxz version", "config", cfg)

	// 新建lxz应用实例
	app := view.NewApp(cfg)
	// 应用初始化
	if err := app.Init(version, 2); err != nil {
		return err
	}

	// 应用启动
	if err := app.Run(); err != nil {
		return err
	}
	//if view.ExitStatus != "" {
	//	return fmt.Errorf("view exit status %s", view.ExitStatus)
	//}

	return nil
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func initLXZFlags() {
	lxzFlags = config.NewFlags()
	rootCmd.Flags().IntVarP(
		lxzFlags.RefreshRate,
		"refresh", "r",
		config.DefaultRefreshRate,
		"Specify the default refresh rate as an integer (sec)",
	)
	rootCmd.Flags().StringVarP(
		lxzFlags.LogLevel,
		"logLevel", "l",
		config.DefaultLogLevel,
		"Specify a log level (error, warn, info, debug)",
	)
	rootCmd.Flags().StringVarP(
		lxzFlags.LogFile,
		"logFile", "",
		config.AppLogFile,
		"Specify the log file",
	)
	rootCmd.Flags().BoolVar(
		lxzFlags.Headless,
		"headless",
		false,
		"Turn LXZ header off",
	)
	rootCmd.Flags().BoolVar(
		lxzFlags.Logoless,
		"logoless",
		false,
		"Turn LXZ logo off",
	)
	rootCmd.Flags().BoolVar(
		lxzFlags.Crumbsless,
		"crumbsless",
		false,
		"Turn LXZ crumbs off",
	)
	rootCmd.Flags().BoolVar(
		lxzFlags.Splashless,
		"splashless",
		false,
		"Turn LXZ splash screen off",
	)

	rootCmd.Flags()
}

func loadConfiguration() (*config.Config, error) {
	slog.Info("🐶 lxz starting up...")

	lxzCfg := config.NewConfig()
	var errs error

	if err := lxzCfg.Load(config.AppConfigFile, false); err != nil {
		errs = errors.Join(errs, err)
	}

	return lxzCfg, errs
}
