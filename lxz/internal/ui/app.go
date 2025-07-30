/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 16:26
 */

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/slogs"
	"sync"
)

type App struct {
	*tview.Application
	Main    *Pages // 主页面容器（目前就两个页面，启动动画、功能页）
	actions *KeyActions
	views   map[string]tview.Primitive
	running bool
	// 应用配置
	Config *config.Config
	// 应用样式
	Styles *config.Styles

	mx sync.RWMutex
}

func NewApp(cfg *config.Config) *App {
	a := App{
		// 应用本体
		Application: tview.NewApplication(),
		// 存放快捷键的行为
		actions: NewKeyActions(),
		// 主页面容器
		Main:   NewPages(),
		Config: cfg,
		Styles: config.NewStyles(),
	}

	// 初始化应用的部分小组件
	a.views = map[string]tview.Primitive{
		"logo":    NewLogo(a.Styles),    // logo
		"status":  NewStatus(a.Styles),  // status
		"menu":    NewMenu(a.Styles),    // menu
		"subMenu": NewSubMenu(a.Styles), // submenu
	}

	return &a
}

// Views return the application root views.
func (a *App) Views() map[string]tview.Primitive {
	return a.views
}

func (a *App) Status() *Status {
	return a.views["status"].(*Status)
}

// SubMenu 用于显示各个功能页面的快捷键描述
func (a *App) SubMenu() *SubMenu {
	if v, ok := a.views["subMenu"]; ok {
		return v.(*SubMenu)
	}
	slog.Error("SubMenu not found", slogs.Subsys, "ui", slogs.Component, "app")
	return nil
}

func (a *App) Logo() *Logo {
	return a.views["logo"].(*Logo)
}

func (a *App) Menu() *Menu {
	return a.views["menu"].(*Menu)
}

func (a *App) bindKeys() {
	a.actions = NewKeyActionsFromMap(KeyMap{
		// 激活命令行模式
		//KeyColon:       NewKeyAction("Cmd", a.activateCmd, false),
		//tcell.KeyCtrlR: NewKeyAction("Redraw", a.redrawCmd, false),
		//tcell.KeyCtrlP: NewKeyAction("Persist", a.saveCmd, false),
		//tcell.KeyCtrlU: NewSharedKeyAction("Clear Filter", a.clearCmd, false),
		//tcell.KeyCtrlQ: NewSharedKeyAction("Clear Filter", a.clearCmd, false),
	})
}

// StylesChanged notifies the skin changed.
func (a *App) StylesChanged(s *config.Styles) {
	a.Main.SetBackgroundColor(s.BgColor())
	if f, ok := a.Main.GetPrimitive("main").(*tview.Flex); ok {
		f.SetBackgroundColor(s.BgColor())
		if !a.Config.LXZ.IsHeadless() {
			if h, ok := f.ItemAt(0).(*tview.Flex); ok {
				h.SetBackgroundColor(s.BgColor())
			} else {
				slog.Warn("Header not found", slogs.Subsys, "styles", slogs.Component, "app")
			}
		}
	} else {
		slog.Error("Main panel not found", slogs.Subsys, "styles", slogs.Component, "app")
	}
}

// Init initializes the application.
func (a *App) Init() {
	// 绑定快捷键
	a.bindKeys()
	a.Styles.AddListener(a)
	// 设置应用程序的根视图 是一个Main容器,里面存了很多个页面
	a.SetRoot(a.Main, true).EnableMouse(true)
}

// HasAction checks if key matches a registered binding.
func (a *App) HasAction(key tcell.Key) (KeyAction, bool) {
	return a.actions.Get(key)
}

// AsKey converts rune to keyboard key.
func AsKey(evt *tcell.EventKey) tcell.Key {
	// 如果按键不是字符键，则直接返回按键
	if evt.Key() != tcell.KeyRune {
		return evt.Key()
	}
	key := tcell.Key(evt.Rune())
	// 如果按键是Alt键，则将按键值与修饰符相乘
	if evt.Modifiers() == tcell.ModAlt {
		key = tcell.Key(int16(evt.Rune()) * int16(evt.Modifiers()))
	}
	return key
}

// AddActions returns the application actions.
func (a *App) AddActions(aa *KeyActions) {
	a.actions.Merge(aa)
}

// QueueUpdateDraw queues up a ui action and redraw the ui.
func (a *App) QueueUpdateDraw(f func()) {
	if a.Application == nil {
		return
	}
	go func() {
		a.Application.QueueUpdateDraw(f)
	}()
}

// IsRunning checks if app is actually running.
func (a *App) IsRunning() bool {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return a.running
}

// SetRunning sets the app run state.
func (a *App) SetRunning(f bool) {
	a.mx.Lock()
	defer a.mx.Unlock()
	a.running = f
}
