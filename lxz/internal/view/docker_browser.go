/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/rivo/tview"
	"io"
	"log"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/drivers/docker_drivers"
	"lxz/internal/helper"
	"lxz/internal/ui"
	"lxz/internal/ui/dialog"
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
		tcell.KeyEnter: ui.NewKeyAction("Logs", _this.EmptyKeyEvent, true),
		tcell.KeyCtrlR: ui.NewKeyAction("Restart", _this.restartContainer, true),
	})
}

// restartContainer
func (_this *DockerBrowser) restartContainer(evt *tcell.EventKey) *tcell.EventKey {
	dialog.ShowConfirm(&config.Dialog{},
		_this.app.Content.Pages,
		"Are you sure you want to restart the container?",
		_this.selectContainerName,
		func(force bool) {
			loading := dialog.ShowLoadingDialog(appViewInstance.Content.Pages, "", appUiInstance.ForceDraw)
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
}

func (_this *DockerBrowser) _refreshData() {
	ctrList, err := docker_drivers.ListContainers()
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return
	}
	// 填充容器数据
	for i, ctr := range ctrList {
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
		_this.app.inject(NewDockerLogsPage(_this.app, connID, connName), false)
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
	//if _this.stopDebounceCh != nil {
	//	close(_this.stopDebounceCh) // 停止防抖协程
	//}
}

// --- HELPER FUNCTIONS ---

func execIntoContainer(app *tview.Application, containerID string) {
	cli, _ := client.NewClientWithOpts(client.FromEnv)

	execConfig := container.ExecOptions{
		Cmd:          []string{"sh"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := cli.ContainerExecCreate(context.TODO(), containerID, execConfig)
	if err != nil {
		log.Fatal(err)
	}

	attachResp, err := cli.ContainerExecAttach(context.TODO(), execIDResp.ID,
		container.ExecAttachOptions{Tty: true})
	if err != nil {
		log.Fatal(err)
	}

	term := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	go func() {
		io.Copy(term, attachResp.Reader)
	}()

	inputField := tview.NewInputField().
		SetLabel("> ").
		SetDoneFunc(func(key tcell.Key) {
			//input := inputField.GetText()
			//if input == "exit" {
			//	app.SetRoot(containerTable, true) // 回到主界面
			//	return
			//}
			//attachResp.Conn.Write([]byte(input + "\n"))
			//inputField.SetText("")
		})

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(term, 0, 1, false).
		AddItem(inputField, 1, 0, true)

	app.SetRoot(layout, true)
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
