/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 16:26
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal"
	"lxz/internal/config"
	"lxz/internal/model"
	"lxz/internal/slogs"
	"lxz/internal/ui"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	splashDelay      = 1 * time.Second
	clusterRefresh   = 15 * time.Second
	clusterInfoWidth = 25
	clusterInfoPad   = 15
)

// App represents the application view layer.
type App struct {
	version string
	UI      *ui.App

	Content *PageStack

	showHeader bool
	showLogo   bool
}

// NewApp returns a K9s app instance.
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
	//a.UI.Views()["clusterInfo"] = NewClusterInfo(&a)
	slog.Info("LXZ 目前只完成了实例的初始化")

	return &a
}

// ReloadStyles reloads skin file.
func (a *App) ReloadStyles() {

}

func (a *App) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info("LXZ App keyboard event", "key", evt.Key(), "rune", evt.Rune(), "modifiers", evt.Modifiers())
	if k, ok := a.UI.HasAction(ui.AsKey(evt)); ok && !a.Content.IsTopDialog() {
		return k.Action(evt)
	}

	return evt
}

func (a *App) toggleHeaderCmd(evt *tcell.EventKey) *tcell.EventKey {

	a.UI.QueueUpdateDraw(func() {
		a.showHeader = !a.showHeader
		a.toggleHeader(a.showHeader, a.showLogo)
	})

	return nil
}

func (a *App) testContentChange(evt *tcell.EventKey) *tcell.EventKey {
	a.inject(ui.NewTestComp(time.Now().Format("15:04:05")), false)
	return nil
}

func (a *App) menuPageChange(evt *tcell.EventKey) *tcell.EventKey {
	switch evt.Rune() {
	case rune(ui.Key1):
		if err := a.inject(NewProjectRelease(), true); err != nil {
			slog.Error("Failed to inject ProjectRelease component", slogs.Error, err)
		}
	default:
		slog.Warn("Unknown menu page change key", "key", evt.Rune())
		return evt
	}

	return nil
}

// PrevCmd pops the command stack.
func (a *App) PrevCmd(*tcell.EventKey) *tcell.EventKey {
	if !a.Content.IsLast() {
		a.Content.Pop()
	}

	return nil
}

func (a *App) bindKeys() {
	a.UI.AddActions(ui.NewKeyActionsFromMap(ui.KeyMap{
		tcell.KeyCtrlE: ui.NewSharedKeyAction("ToggleHeader", a.toggleHeaderCmd, false),
		//tcell.KeyCtrlG:     ui.NewSharedKeyAction("toggleCrumbs", a.toggleCrumbsCmd, false),
		//ui.KeyHelp: ui.NewSharedKeyAction("Help", a.helpCmd, false),
		ui.KeyHelp: ui.NewSharedKeyAction("Test", a.testContentChange, false),
		ui.Key1:    ui.NewSharedKeyAction("projectRelease", a.menuPageChange, false),
		//ui.KeyLeftBracket:  ui.NewSharedKeyAction("Go Back", a.previousCommand, false),
		//ui.KeyRightBracket: ui.NewSharedKeyAction("Go Forward", a.nextCommand, false),
		//ui.KeyDash:         ui.NewSharedKeyAction("Last View", a.lastCommand, false),
		//tcell.KeyCtrlA:     ui.NewSharedKeyAction("Aliases", a.aliasCmd, false),
		//tcell.KeyEnter:     ui.NewKeyAction("Goto", a.gotoCmd, false),
		//tcell.KeyCtrlC:     ui.NewKeyAction("Quit", a.quitCmd, false),
	}))
}

func (a *App) buildHeader() tview.Primitive {
	header := tview.NewFlex()
	header.SetBackgroundColor(a.UI.Styles.BgColor())
	header.SetDirection(tview.FlexColumn)
	header.SetBorder(true)
	if !a.showHeader {
		return header
	}
	//
	header.AddItem(a.UI.Menu(), 0, 1, false)
	if a.showLogo {
		header.AddItem(a.UI.Logo(), 24, 1, false)
	}
	go func() {
		for {
			a.UI.Logo().Status().Clear()
			fmt.Fprint(a.UI.Logo().Status(), time.Now().Format("15:04:05"))
			a.UI.QueueUpdateDraw(func() {})
			time.Sleep(1 * time.Second)
		}
	}()

	return header
}

func (a *App) toggleHeader(header, logo bool) {
	a.showHeader, a.showLogo = header, logo
	flex, ok := a.UI.Main.GetPrimitive("main").(*tview.Flex)
	if !ok {
		slog.Error("Expecting flex view main panel. Exiting!")
		os.Exit(1)
	}
	if a.showHeader {
		flex.AddItemAtIndex(0, a.buildHeader(), 10, 1, false)
	} else {
		flex.RemoveItemAtIndex(0)
	}
}

func (a *App) layout(ctx context.Context) {

	// 主页是一个flex布局应用 方向是垂直的
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	a.showHeader = true
	a.showLogo = true

	// header 组件
	main.AddItem(a.buildHeader(), 10, 1, false)

	// 状态指示器
	//main.AddItem(a.UI.Status(), 5, 1, false)
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

// Init initializes the application.
func (a *App) Init(version string, _ int) error {
	a.version = model.NormalizeVersion(version)

	ctx := context.WithValue(context.Background(), internal.KeyApp, a)
	if err := a.Content.Init(ctx); err != nil {
		return err
	}

	// 面包屑组件添加监听
	//a.Content.AddListener(a.UI.Status())
	a.Content.AddListener(a.UI.SubMenu())
	// 快捷键+数据初始化
	a.UI.Init()

	a.UI.SetInputCapture(a.keyboard)
	// 绑定快捷键
	a.bindKeys()

	// 初始化命令组件 通过命令实现页面的
	//a.command = NewCommand(a)
	//if err := a.command.Init(a.Config.ContextAliasesPath()); err != nil {
	//	return err
	//}

	// 初始化布局
	a.layout(ctx)

	// 初始化信号监听
	a.initSignals()

	// 刷新样式
	a.ReloadStyles()

	slog.Info("LXZ 此时UI-APP已经初始化完成")

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
			// if command bar is already active, focus it
			//if a.CmdBuff().IsActive() {
			//	a.SetFocus(a.Prompt())
			//}
		})
	}()

	slog.Info("🚀 LXZ 开始运行UI-APP, 执行默认命令 defaultCmd")
	//if err := a.command.defaultCmd(true); err != nil {
	//	return err
	//}
	a.UI.SetRunning(true)
	if err := a.UI.Application.Run(); err != nil {
		return err
	}

	return nil
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
