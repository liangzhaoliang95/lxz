package view

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/liangzhaoliang95/lxz/internal/drivers/docker_drivers"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

type DockerLogsPage struct {
	*BaseFlex
	app      *App            // 应用实例
	logsView *tview.TextView // 日志视图

	containerId   string // 当前选中的容器ID
	containerName string // 当前选中的容器名称
}

func (_this *DockerLogsPage) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF: ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
	})
}

func (_this *DockerLogsPage) showLogsForSelectedContainer(containerID string) {
	reader, err := docker_drivers.ContainerLogs(containerID)
	if err != nil {
		slog.Error("Failed to get logs for container", "containerID", containerID, "error", err)
		return
	}

	go func() {
		defer func() {
			_ = reader.Close()
		}()
		buf := make([]byte, 4096)
		for {
			n, err := reader.Read(buf)
			slog.Info("Reading logs", "containerID", containerID, "bytesRead", n, "error", err)
			if n > 0 {
				_, _ = fmt.Fprintf(_this.logsView, "%s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()
}

func (_this *DockerLogsPage) Init(ctx context.Context) error {
	slog.Info(
		"Initializing DockerLogsPage",
		"containerId",
		_this.containerId,
		"containerName",
		_this.containerName,
	)
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	logsView := tview.NewTextView()
	logsView.SetDynamicColors(true).
		SetWrap(false).
		SetRegions(true).
		SetScrollable(true).
		SetBorder(false)
	logsView.SetWordWrap(true)
	logsView.SetWrap(true)
	logsView.SetChangedFunc(func() {
		_this.app.UI.QueueUpdateDraw(func() {
			logsView.ScrollToEnd() // 👈 自动滚动到最后
		})
	})

	_this.logsView = logsView

	_this.
		AddItem(logsView, 0, 1, true)
	return nil
}

func (_this *DockerLogsPage) Start() {
	slog.Info(
		"Starting DockerLogsPage",
		"containerId",
		_this.containerId,
		"containerName",
		_this.containerName,
	)
	_this.showLogsForSelectedContainer(_this.containerId)
	_this.app.UI.SetFocus(_this)
}

func (_this *DockerLogsPage) Stop() {}

func NewDockerLogsPage(app *App, containerId string, containerName string) *DockerLogsPage {
	page := &DockerLogsPage{
		BaseFlex:      NewBaseFlex(containerName),
		app:           app,
		containerId:   containerId,
		containerName: containerName,
	}

	return page
}
