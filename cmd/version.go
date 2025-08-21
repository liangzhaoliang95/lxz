package cmd

import (
	"fmt"

	"github.com/liangzhaoliang95/lxz/internal/color"
	ver "github.com/liangzhaoliang95/lxz/internal/version"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	var short bool
	var checkUpdate bool

	command := cobra.Command{
		Use:   "version",
		Short: "Print version/build info",
		Long:  "Print version/build information and check for updates",
		Run: func(*cobra.Command, []string) {
			printVersion(short)
			if checkUpdate {
				checkForUpdates()
			}
		},
	}

	command.PersistentFlags().
		BoolVarP(&short, "short", "s", false, "Prints LXZ version info in short format")
	command.PersistentFlags().
		BoolVarP(&checkUpdate, "check-update", "c", false, "Check for available updates")

	return &command
}

func printVersion(short bool) {
	const fmat = "%-20s %s\n"
	var outputColor color.Paint

	if short {
		outputColor = -1
	} else {
		outputColor = color.Cyan
		printLogo(outputColor)
	}

	v := ver.GetVersion()
	printTuple(fmat, "Version", v.Version, outputColor)
	printTuple(fmat, "Commit", v.Commit, outputColor)
	printTuple(fmat, "Date", v.Date, outputColor)
	printTuple(fmat, "Go Version", v.GoVersion, outputColor)
	printTuple(fmat, "Platform", fmt.Sprintf("%s/%s", v.Platform, v.Architecture), outputColor)
}

func printTuple(fmat, section, value string, outputColor color.Paint) {
	if outputColor != -1 {
		_, _ = fmt.Fprintf(out, fmat, color.Colorize(section+":", outputColor), value)
		return
	}
	_, _ = fmt.Fprintf(out, fmat, section, value)
}

func checkForUpdates() {
	fmt.Println()
	fmt.Println("检查更新中...")

	updateInfo, err := ver.CheckForUpdates()
	if err != nil {
		fmt.Printf("检查更新失败: %v\n", err)
		return
	}

	if updateInfo == nil {
		fmt.Println("当前已是最新版本！")
		return
	}

	fmt.Printf("发现新版本: %s (当前: %s)\n", updateInfo.LatestVersion, updateInfo.CurrentVersion)
	fmt.Printf("发布日期: %s\n", updateInfo.PublishedAt.Format("2006-01-02 15:04:05"))

	if updateInfo.DownloadURL != "" {
		fmt.Printf("下载地址: %s\n", updateInfo.DownloadURL)
	}

	if updateInfo.ReleaseNotes != "" {
		fmt.Println("\n更新说明:")
		fmt.Println(updateInfo.ReleaseNotes)
	}
}
