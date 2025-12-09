/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/drivers/docker_drivers"
	"github.com/liangzhaoliang95/lxz/internal/helper"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/lxz/internal/ui/dialog"
	"github.com/liangzhaoliang95/tview"
)

type DockerBrowser struct {
	*BaseFlex
	app                 *App
	containerTableUI    *tview.Table
	selectedContainerID string // 当前选中的容器ID
	selectContainerName string // 当前选中的容器名称
}

func (_this *DockerBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF:        ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		ui.KeyI:        ui.NewKeyAction("Detail Info", _this.ShowDetail, true),
		ui.KeyS:        ui.NewKeyAction("Shell", _this.ShellExec, true),
		tcell.KeyEnter: ui.NewKeyAction("Logs", _this.EmptyKeyEvent, true),
		tcell.KeyCtrlR: ui.NewKeyAction("Restart", _this.restartContainer, true),
		tcell.KeyCtrlD: ui.NewKeyAction("Stop Or Delete", _this.stopDeleteContainer, true),
	})
}

// ShellExec 执行容器的Shell
func (_this *DockerBrowser) ShellExec(evt *tcell.EventKey) *tcell.EventKey {
	if _this.selectedContainerID == "" {
		_this.app.UI.Flash().Err(fmt.Errorf("please select a container first"))
		return nil
	}
	// 执行容器的Shell
	slog.Info(
		"Executing shell in container",
		"name",
		_this.selectContainerName,
		"containerID",
		_this.selectedContainerID,
	)
	err := _this.containerShellIn()
	if err != nil {
		_this.app.UI.Flash().
			Err(fmt.Errorf("failed to execute shell in container %s: %w", _this.selectedContainerID, err))
		return nil
	}

	return nil
}

// ShowDetail 显示容器的详细信息
func (_this *DockerBrowser) ShowDetail(evt *tcell.EventKey) *tcell.EventKey {
	if _this.selectedContainerID == "" {
		_this.app.UI.Flash().Err(fmt.Errorf("please select a container first"))
		return nil
	}
	if err := _this.app.inject(NewDockerInspectView(_this.app, "container", _this.selectedContainerID, _this.selectContainerName), false); err != nil {
		_this.app.UI.Flash().Err(fmt.Errorf("failed to inject docker inspect view: %w", err))
	}

	return nil
}

// restartContainer
func (_this *DockerBrowser) restartContainer(evt *tcell.EventKey) *tcell.EventKey {
	dialog.ShowConfirm(&config.Dialog{},
		_this.app.Content.Pages,
		"Are you sure you want to restart the container?",
		_this.selectContainerName,
		func(force bool) {
			loading := dialog.ShowLoadingDialog(
				appViewInstance.Content.Pages,
				"",
				appUiInstance.ForceDraw,
			)
			var timeout *int
			if force {
				timeout = helper.Ptr(0)
			}
			err := docker_drivers.RestartContainer(_this.selectedContainerID, timeout)
			if err != nil {
				_this.app.UI.Flash().Err(err)
			} else {
				_this.app.UI.Flash().Info("Container restarted successfully")
				_this._refreshData() // 刷新数据
			}
			loading.Hide()
		},
		func() {

		})

	return nil
}

// stopDeleteContainer
func (_this *DockerBrowser) stopDeleteContainer(evt *tcell.EventKey) *tcell.EventKey {
	detail, err := docker_drivers.InspectContainer(_this.selectedContainerID)
	if err != nil {
		_this.app.UI.Flash().
			Err(fmt.Errorf("failed to inspect container %s: %w", _this.selectedContainerID, err))
		return nil
	}
	operation := "stop"
	if detail.State.Running {
		operation = "stop"
	} else {
		operation = "delete"
	}

	dialog.ShowConfirm(&config.Dialog{},
		_this.app.Content.Pages,
		fmt.Sprintf("Are you sure you want to %s the container?", operation),
		_this.selectContainerName,
		func(force bool) {
			loading := dialog.ShowLoadingDialog(
				appViewInstance.Content.Pages,
				"",
				appUiInstance.ForceDraw,
			)
			var timeout *int
			if force {
				timeout = helper.Ptr(0)
			}
			var err error
			if operation == "delete" {
				err = docker_drivers.RemoveContainer(_this.selectedContainerID, force)
			} else {
				if err = docker_drivers.StopContainer(_this.selectedContainerID, timeout); err == nil {
					err = docker_drivers.WaitContainerStopped(_this.selectedContainerID, time.Duration(60)*time.Second)
				}
			}

			if err != nil {
				_this.app.UI.Flash().Err(err)
			} else {
				_this.app.UI.Flash().Info(fmt.Sprintf("Container:%s %s successfully", _this.selectContainerName, operation))
				_this._refreshData() // 刷新数据
			}
			loading.Hide()
		},
		func() {

		})
	return nil
}

func (_this *DockerBrowser) _initHeader() {
	// 初始化header
	_this.containerTableUI.SetCell(0, 0, tview.NewTableCell("ID").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.containerTableUI.SetCell(0, 1, tview.NewTableCell("Name").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.containerTableUI.SetCell(0, 2, tview.NewTableCell("Image").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.containerTableUI.SetCell(0, 3, tview.NewTableCell("Created").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.containerTableUI.SetCell(0, 4, tview.NewTableCell("Status").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.containerTableUI.SetCell(0, 5, tview.NewTableCell("State").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.containerTableUI.SetCell(0, 6, tview.NewTableCell("Port").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
}

func (_this *DockerBrowser) _refreshData() {
	ctrList, err := docker_drivers.ListContainers()
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return
	}
	// 填充容器数据
	for i, ctr := range ctrList {
		if i == 0 {
			_this.selectedContainerID = ctr.ID
			_this.selectContainerName = ctr.Names[0][1:] // 默认选中
		}
		_this.containerTableUI.SetCell(i+1, 0, tview.NewTableCell(ctr.ID[:12]).
			SetReference(ctr.ID).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		_this.containerTableUI.SetCell(i+1, 1, tview.NewTableCell(ctr.Names[0][1:]).
			SetReference(ctr.Names[0]).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		_this.containerTableUI.SetCell(i+1, 2, tview.NewTableCell(ctr.Image).
			SetReference(ctr.Status).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		_this.containerTableUI.SetCell(i+1, 3, tview.NewTableCell(helper.TimeFormat(ctr.Created)).
			SetReference(ctr.Status).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		_this.containerTableUI.SetCell(i+1, 4, tview.NewTableCell(ctr.Status).
			SetReference(ctr.Status).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		_this.containerTableUI.SetCell(i+1, 5, tview.NewTableCell(ctr.State).
			SetReference(ctr.Status).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		ports := ""
		if len(ctr.Ports) > 0 {
			for _, port := range ctr.Ports {
				if ports != "" {
					ports += ", "
				}
				if port.PublicPort > 0 {
					ports += fmt.Sprintf(
						"%s:%d -> %d/%s",
						port.IP,
						port.PublicPort,
						port.PrivatePort,
						port.Type,
					)
				} else {
					ports += fmt.Sprintf("%d/%s", port.PrivatePort, port.Type)
				}
			}
		} else {
			ports = "None"
		}
		_this.containerTableUI.SetCell(i+1, 6, tview.NewTableCell(ports).
			SetReference(ctr.Status).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
	}
}

func (_this *DockerBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)
	_this.containerTableUI = tview.NewTable()
	_this.containerTableUI.SetBorder(false)
	_this.containerTableUI.SetBorders(false)
	_this.containerTableUI.SetTitle("🌍 Connections")
	_this.containerTableUI.SetBorderPadding(1, 1, 2, 2)
	_this.containerTableUI.SetSelectable(true, false)
	_this.containerTableUI.Select(1, 0)
	// 配置回车函数
	_this.containerTableUI.SetSelectedFunc(func(row, column int) {
		slog.Info("Selected connection", "row", row, "col", column)
		// 获取选中的连接信息
		if row < 1 || row >= _this.containerTableUI.GetRowCount() {
			slog.Warn("Selected row is out of range", "row", row)
			return
		}
		connName := _this.containerTableUI.GetCell(row, 1).Text
		// 渲染日志页面
		connID := _this.containerTableUI.GetCell(row, 0).Text
		if err := _this.app.inject(NewDockerLogsPage(_this.app, connID, connName), false); err != nil {
			_this.app.UI.Flash().Err(fmt.Errorf("failed to inject docker logs page: %w", err))
		}
	})
	// 设置表格的选择模式
	_this.containerTableUI.SetSelectionChangedFunc(func(row, column int) {
		slog.Info("Selection changed", "row", row, "col", column)
		if row < 1 || row >= _this.containerTableUI.GetRowCount() {
			slog.Warn("Selection changed row is out of range", "row", row)
			return
		}
		_this.selectedContainerID = _this.containerTableUI.GetCell(row, 0).Text
		_this.selectContainerName = _this.containerTableUI.GetCell(row, 1).Text
	})
	_this._initHeader()

	_this.AddItem(_this.containerTableUI, 0, 1, true)

	return nil
}

func (_this *DockerBrowser) Start() {
	// ✅ 设置默认边框颜色 + 焦点 + 强制刷新
	_this._refreshData()

	_this.app.UI.SetFocus(_this.containerTableUI)
}

func (_this *DockerBrowser) Stop() {
	// if _this.stopDebounceCh != nil {
	//	close(_this.stopDebounceCh) // 停止防抖协程
	//}
}

// --- HELPER FUNCTIONS ---

func (_this *DockerBrowser) containerShellIn() error {
	_this.Stop()
	defer _this.Start()

	_this.shellIn()
	return nil
}

func (_this *DockerBrowser) shellIn() {
	args := make([]string, 0, 15)

	args = append(args, "exec", "-it", _this.selectedContainerID)
	args = append(args, "sh", "-c", shellCheck)

	slog.Info("Shell args", "args", args)
	c := color.New(color.BgGreen).Add(color.FgBlack).Add(color.Bold)
	err := runDockerExec(_this.app, &shellOpts{
		clear:  true,
		banner: c.Sprintf(bannerFmt, ""),
		args:   args},
	)
	if err != nil {
		_this.app.UI.Flash().Errf("Shell exec failed: %s", err)
	}
}

func NewDockerBrowser(app *App) *DockerBrowser {
	var name = "Docker Browser"
	f := &DockerBrowser{
		BaseFlex: NewBaseFlex(name),
		app:      app,
	}
	f.SetIdentifier(ui.DOCKER_BROWSER_ID)
	return f
}
