/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"lxz/internal/drivers/docker_drivers"
	"lxz/internal/helper"
	"lxz/internal/ui"
)

type DockerInspectView struct {
	*BaseFlex
	app           *App
	inspectViewUI *tview.TextView // 用于显示容器的详细信息
	inspectId     string          // 当前选中的容器ID
	inspectName   string          // 当前选中的容器名称
	inspectType   string          // image container network volume
}

func (_this *DockerInspectView) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF: ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
	})
}

func (_this *DockerInspectView) _refreshInspectView() {
	if _this.inspectType == "" {
		_this.app.UI.Flash().Err(fmt.Errorf("inspectType is empty"))
		return
	}
	switch _this.inspectType {
	case "container":
		_this._inspectContainer()
	case "image":
		_this._inspectImage()
	case "network":
		_this._inspectNetwork()
	case "volume":
		_this._inspectVolume()
	default:
		_this.app.UI.Flash().Err(fmt.Errorf("unknown inspect type: %s", _this.inspectType))
	}
}

func (_this *DockerInspectView) _inspectContainer() {
	res, err := docker_drivers.InspectContainer(_this.inspectId)
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return
	}
	var output string
	if res == nil {
		output = "No container found with the given ID."
	} else {
		output = helper.Prettify(res)
	}
	_this.inspectViewUI.SetText(output)
	_this.inspectViewUI.SetText(helper.Prettify(res))
}

func (_this *DockerInspectView) _inspectImage() {
	res, err := docker_drivers.InspectImage(_this.inspectId)
	if err != nil {
		_this.app.UI.Flash().Err(err)
		return
	}
	// 处理镜像的Inspect结果
	var output string
	if res == nil {
		output = "No image found with the given ID."
	} else {
		output = helper.Prettify(res)
	}
	_this.inspectViewUI.SetText(output)

}

func (_this *DockerInspectView) _inspectNetwork() {

}

func (_this *DockerInspectView) _inspectVolume() {

}

func (_this *DockerInspectView) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)
	_this.inspectViewUI = tview.NewTextView()
	_this.inspectViewUI.SetBorder(false)
	_this.inspectViewUI.SetTitle("")
	_this.inspectViewUI.SetWrap(true)
	_this.inspectViewUI.SetBorderPadding(1, 1, 2, 2)
	_this.AddItem(_this.inspectViewUI, 0, 1, true)

	return nil
}

func (_this *DockerInspectView) Start() {
	// ✅ 设置默认边框颜色 + 焦点 + 强制刷新
	_this._refreshInspectView()

	_this.app.UI.SetFocus(_this.inspectViewUI)
}

func (_this *DockerInspectView) Stop() {

}

// --- HELPER FUNCTIONS ---

func NewDockerInspectView(app *App, inspectType string, inspectId, inspectName string) *DockerInspectView {
	var name = fmt.Sprintf("Inspect %s: %s", inspectType, inspectName)
	f := &DockerInspectView{
		BaseFlex:    NewBaseFlex(name),
		app:         app,
		inspectId:   inspectId,
		inspectName: inspectName,
		inspectType: inspectType,
	}
	return f
}
