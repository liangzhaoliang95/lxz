/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 16:26
 */

package view

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"github.com/liangzhaoliang95/lxz/internal/model"
	"github.com/liangzhaoliang95/lxz/internal/slogs"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

const (
	splashDelay      = 1 * time.Second
	clusterRefresh   = 15 * time.Second
	clusterInfoWidth = 25
	clusterInfoPad   = 15
)

var appViewInstance *App
var appUiInstance *ui.App

// App represents the application view layer.
type App struct {
	version string
	UI      *ui.App

	Content *PageStack

	cancelFn   context.CancelFunc
	showHeader bool
	showLogo   bool
}

// NewApp returns a LXZ app instance.
func NewApp(cfg *config.Config) *App {
	a := App{
		// 初始化UI-APP
		UI: ui.NewApp(cfg),
		// 应用内容主体 本质是个pages容器，是核心
		Content: NewPageStack(),
	}
	a.ReloadStyles()

	// 初始化应用的部分小组件
	// 集群信息组件
	// a.UI.Views()["clusterInfo"] = NewClusterInfo(&a)
	slog.Info("LXZ 目前只完成了实例的初始化")
	appViewInstance = &a
	appUiInstance = a.UI
	return &a
}

// ReloadStyles reloads skin file.
func (a *App) ReloadStyles() {

}

func (a *App) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info(
		"LXZ App keyboard event",
		"key",
		evt.Key(),
		"rune",
		evt.Rune(),
		"modifiers",
		evt.Modifiers(),
	)
	var key tcell.Key
	if evt.Modifiers() == tcell.ModCtrl {
		key = tcell.Key(evt.Rune())
	} else {
		key = ui.AsKey(evt)
	}
	if k, ok := a.UI.HasAction(key); ok && !a.Content.IsTopDialog() {
		return k.Action(evt)
	}

	return evt
}

func (a *App) menuPageChange(evt *tcell.EventKey) *tcell.EventKey {
	changeSuccess := false
	pageName := ""
	slog.Info("Menu page change event", "key", evt.Key(), "rune", evt.Rune())
	switch evt.Key() {
	case tcell.KeyF1:
		comp := NewSshConnect(a)
		if err := a.inject(comp, true); err != nil {
			slog.Error("Failed to inject SshConnect component", slogs.Error, err)
		} else {
			changeSuccess = true
			pageName = comp.Name()
		}
	case tcell.KeyF2:
		comp := NewFileBrowser(a)
		if err := a.inject(comp, true); err != nil {
			slog.Error("Failed to inject FileBrowser component", slogs.Error, err)
		} else {
			changeSuccess = true
			pageName = comp.Name()
		}
	case tcell.KeyF3:
		comp := NewRedisBrowser(a)
		if err := a.inject(comp, true); err != nil {
			slog.Error("Failed to inject Redis component", slogs.Error, err)
		} else {
			changeSuccess = true
			pageName = comp.name
		}
	case tcell.KeyF4:
		comp := NewDatabaseBrowser(a)
		if err := a.inject(comp, true); err != nil {
			slog.Error("Failed to inject RedisBrowser component", slogs.Error, err)
		} else {
			changeSuccess = true
			pageName = comp.Name()
		}
	case tcell.KeyF5:
		comp := NewDockerBrowser(a)
		if err := a.inject(comp, true); err != nil {
			slog.Error("Failed to inject DockerBrowser component", slogs.Error, err)
		} else {
			changeSuccess = true
			pageName = comp.Name()
		}
	case tcell.KeyF6:
		comp := NewK9SBrowser(a)
		if err := a.inject(comp, true); err != nil {
			slog.Error("Failed to inject DockerBrowser component", slogs.Error, err)
		} else {
			changeSuccess = true
			pageName = comp.Name()
		}
	default:
		slog.Warn("Unknown menu page change key", "key", evt.Rune())
		a.UI.Flash().Err(fmt.Errorf("unknown menu page change key: %d", evt.Rune()))
		return evt
	}
	if changeSuccess {
		a.UI.Flash().Info(fmt.Sprintf("Switched to <%s> page", pageName))
	}

	return nil
}

// PrevCmd pops the command stack.
func (a *App) PrevCmd(*tcell.EventKey) *tcell.EventKey {
	if !a.Content.IsLast() && !a.Content.IsTopDialog() {
		a.Content.Pop()
	}

	return nil
}

func (a *App) bindKeys() {
	a.UI.AddActions(ui.NewKeyActionsFromMap(ui.KeyMap{
		// tcell.KeyCtrlE: ui.NewSharedKeyAction("ToggleHeader", a.toggleHeaderCmd, false),
		tcell.KeyEscape: ui.NewSharedKeyAction("Go Back", a.PrevCmd, false),
		tcell.KeyF1:     ui.NewSharedKeyAction("SSH Connect", a.menuPageChange, false),
		tcell.KeyF2:     ui.NewSharedKeyAction("File Browser", a.menuPageChange, false),
		tcell.KeyF3:     ui.NewSharedKeyAction("Redis Browser", a.menuPageChange, false),
		tcell.KeyF4:     ui.NewSharedKeyAction("DB Browser", a.menuPageChange, false),
		tcell.KeyF5:     ui.NewSharedKeyAction("Docker Browser", a.menuPageChange, false),
		tcell.KeyF6:     ui.NewSharedKeyAction("K9S", a.menuPageChange, false),
	}))
}

func (a *App) buildHeader() tview.Primitive {
	header := tview.NewFlex()
	header.SetBackgroundColor(a.UI.Styles.BgColor())
	header.SetDirection(tview.FlexColumn)
	header.SetBorder(false)
	header.SetBorderPadding(0, 0, 0, 1)
	if !a.showHeader {
		return header
	}
	//
	header.AddItem(a.UI.Menu(), 0, 1, false)
	if a.showLogo {
		header.AddItem(a.UI.Logo(), 17, 1, false)
	}
	// go func() {
	//	for {
	//		a.UI.Logo().Status().Clear()
	//		fmt.Fprint(a.UI.Logo().Status(), time.Now().Format("15:04:05"))
	//		a.UI.QueueUpdateDraw(func() {})
	//		time.Sleep(1 * time.Second)
	// }
	// }()

	return header
}

func (a *App) layout(ctx context.Context) {
	flash := ui.NewFlash(a.UI)
	go flash.Watch(ctx, a.UI.Flash().Channel())
	// 主页是一个flex布局应用 方向是垂直的
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	a.showHeader = true
	a.showLogo = true

	// message 组件 用于显示一些提示信息
	main.AddItem(flash, 1, 1, false)
	// header 组件
	main.AddItem(a.buildHeader(), 5, 1, false)

	// 状态指示器
	// main.AddItem(a.UI.Status(), 5, 1, false)
	main.AddItem(a.UI.SubMenu(), 5, 1, false)

	// 内容区域 展示集群资源信息
	main.AddItem(a.Content, 0, 10, true)

	// 往main容器里面添加main组件
	a.UI.Main.AddPage("main", main, true, false)

	// 启动动画
	a.UI.Main.AddPage("splash", ui.NewSplash(a.UI.Styles, a.version), true, true)
}

func (*App) initSignals() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)

	go func(sig chan os.Signal) {
		<-sig
		os.Exit(0)
	}(sig)
}

// Halt stop the application event loop.
func (a *App) Halt() {
	if a.cancelFn != nil {
		a.cancelFn()
		a.cancelFn = nil
	}
}

// Resume restarts the app event loop.
func (a *App) Resume() {
	/*var ctx context.Context
	ctx, a.cancelFn = context.WithCancel(context.Background())

	go a.clusterUpdater(ctx)

	if a.Config.K9s.UI.Reactive {
		if err := a.ConfigWatcher(ctx, a); err != nil {
			slog.Warn("ConfigWatcher failed", slogs.Error, err)
		}
		if err := a.SkinsDirWatcher(ctx, a); err != nil {
			slog.Warn("SkinsWatcher failed", slogs.Error, err)
		}
		if err := a.CustomViewsWatcher(ctx, a); err != nil {
			slog.Warn("CustomView watcher failed", slogs.Error, err)
		}
	}*/
}

// Init initializes the application.
func (a *App) Init(version string, _ int) error {
	a.version = model.NormalizeVersion(version)

	ctx := context.WithValue(context.Background(), internal.KeyApp, a)
	if err := a.Content.Init(ctx); err != nil {
		return err
	}

	// 面包屑组件添加监听
	a.Content.AddListener(a.UI.Menu())
	a.Content.AddListener(a.UI.SubMenu())
	// 快捷键+数据初始化
	a.UI.Init()

	a.UI.SetInputCapture(a.keyboard)
	// 绑定快捷键
	a.bindKeys()

	// 初始化布局
	a.layout(ctx)

	// 初始化信号监听
	a.initSignals()

	// 刷新样式
	a.ReloadStyles()

	slog.Info("LXZ 此时UI-APP已经初始化完成")
	a.UI.Flash().Info("Launch LXZ done")

	return nil
}

// Run starts the application loop.
func (a *App) Run() error {
	go func() {
		if !a.UI.Config.LXZ.IsSplashless() {
			<-time.After(splashDelay)
		}
		a.UI.QueueUpdateDraw(func() {
			a.UI.Main.SwitchToPage("main")

			// 定位到第一个功能
			if err := a.inject(NewSshConnect(a), true); err != nil {
				slog.Error("Failed to inject SshConnect component", slogs.Error, err)
			}
		})
	}()

	slog.Info("🚀 LXZ 开始运行UI-APP")
	// if err := a.command.defaultCmd(true); err != nil {
	//	return err
	//}
	a.UI.SetRunning(true)
	return a.UI.Application.Run()
}

// 往应用中注入一个组件 一般用于激活某个页面
func (a *App) inject(c model.Component, clearStack bool) error {
	ctx := context.WithValue(context.Background(), internal.KeyApp, a)
	if err := c.Init(ctx); err != nil {
		slog.Error("Component init failed",
			slogs.Error, err,
			slogs.CompName, c.Name(),
		)
		return err
	}
	slog.Info("LXZ Injecting component 💉", "component", c.Name())
	if clearStack {
		a.Content.Clear()
	}
	// 将组件添加到应用的页面栈中
	a.Content.Push(c)

	return nil
}
