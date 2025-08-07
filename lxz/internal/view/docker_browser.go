/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"github.com/rivo/tview"
	"lxz/internal/ui"
	"lxz/internal/view/base"
	"time"
)

type DockerBrowser struct {
	*BaseFlex
	app              *App
	rootDir          string
	rootNode         *tview.TreeNode // 根目录节点
	preview          *tview.TextView
	tree             *tview.TreeView
	lastFocusedPath  string
	lastFocusedAt    time.Time
	debounceInterval time.Duration
	stopDebounceCh   chan struct{}
}

func (_this *DockerBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF: ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
	})
}

func (_this *DockerBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)
	_this.tree = tview.NewTreeView()
	_this.preview = tview.NewTextView()
	_this.AddItem(_this.tree, 30, 1, true)
	_this.AddItem(_this.preview, 0, 2, false)

	return nil
}

func (_this *DockerBrowser) Start() {
	// ✅ 设置默认边框颜色 + 焦点 + 强制刷新
	_this.tree.SetBorderColor(base.ActiveBorderColor)
	_this.app.UI.SetFocus(_this)
}

func (_this *DockerBrowser) Stop() {
	//if _this.stopDebounceCh != nil {
	//	close(_this.stopDebounceCh) // 停止防抖协程
	//}
}

func NewDockerBrowser(app *App) *DockerBrowser {
	var name = "Docker Browser"
	f := &DockerBrowser{
		BaseFlex:         NewBaseFlex(name),
		app:              app,
		debounceInterval: 200 * time.Millisecond,
		stopDebounceCh:   make(chan struct{}),
		rootDir:          ".",
	}
	f.SetIdentifier(ui.DOCKER_BROWSER_ID)
	return f
}
