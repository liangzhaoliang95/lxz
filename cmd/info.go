package cmd

import (
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/color"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/spf13/cobra"
)

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "List lxz configurations info",
		RunE:  printInfo,
	}
}

func printInfo(*cobra.Command, []string) error {
	if err := config.InitLocs(); err != nil {
		return err
	}

	const fmat = "%-27s %s\n"
	printLogo(color.Cyan)
	printTuple(fmat, "Version", version, color.Cyan)
	printTuple(fmat, "Config", config.AppConfigFile, color.Cyan)
	printTuple(fmat, "Custom Views", config.AppViewsFile, color.Cyan)
	printTuple(fmat, "Plugins", config.AppPluginsFile, color.Cyan)
	printTuple(fmat, "Hotkeys", config.AppHotKeysFile, color.Cyan)
	printTuple(fmat, "Aliases", config.AppAliasesFile, color.Cyan)
	printTuple(fmat, "Skins", config.AppSkinsDir, color.Cyan)
	printTuple(fmat, "Logs", config.AppLogFile, color.Cyan)
	printTuple(fmat, "DatabaseConfig", config.AppDatabaseConfigFile, color.Cyan)
	printTuple(fmat, "RedisConfig", config.AppRedisConfigFile, color.Cyan)

	return nil
}

func printLogo(c color.Paint) {
	for _, l := range ui.LogoSmall {
		_, _ = fmt.Fprintln(out, color.Colorize(l, c))
	}
	_, _ = fmt.Fprintln(out)
}
