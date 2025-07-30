/**
 * @author  zhaoliang.liang
 * @date  2025/7/30 14:42
 */

package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/apimachinery/pkg/labels"
	"log/slog"
	"lxz/internal/model"
	"lxz/internal/view/cmd"
)

type TestComp struct {
	*tview.Flex
	actions    *KeyActions
	name       string
	fullScreen bool
}

func NewTestComp(boxName string) *TestComp {
	tc := &TestComp{
		Flex:    tview.NewFlex(),
		name:    boxName,
		actions: NewKeyActions(),
	}

	tc.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(boxName)
	tc.AddItem(tview.NewTextView().SetText("This is a test component"), 0, 1, false)
	return tc
}

func (t *TestComp) Name() string {
	return t.name
}

func (t *TestComp) Init(ctx context.Context) error {
	t.bindKeys()
	// 用于初始化组件的边框、标题、快捷键等信息
	t.SetInputCapture(t.keyboard)
	//panic("TestComp init not implemented")
	return nil
}

func (t *TestComp) Start() {
	// 组件启动
	// 通常是配置数据模型以及监听器
	slog.Info("Comp Start", "name", t.Name())
}

func (t *TestComp) Stop() {
	// 组件停止
	// 停止模型 监听器等
	slog.Info("Comp Stop", "name", t.Name())
}

func (t *TestComp) Hints() model.MenuHints {
	return nil
}

func (t *TestComp) ExtraHints() map[string]string {
	return nil
}

func (t *TestComp) InCmdMode() bool {
	return false
}

func (t *TestComp) SetFilter(s string) {
	//TODO implement me
	panic("implement me")
}

func (t *TestComp) SetLabelSelector(selector labels.Selector) {
	//TODO implement me
	panic("implement me")
}

func (t *TestComp) SetCommand(interpreter *cmd.Interpreter) {
	//TODO implement me
	panic("implement me")
}

// Actions returns menu actions.
func (t *TestComp) Actions() *KeyActions {
	return t.actions
}

// Helpers
func (t *TestComp) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	if a, ok := t.actions.Get(AsKey(evt)); ok {
		return a.Action(evt)
	}

	return evt
}

func (t *TestComp) bindKeys() {
	t.Actions().Bulk(KeyMap{
		KeyF: NewKeyAction("Toggle FullScreen", t.toggleFullScreenCmd, true),
	})
}

func (t *TestComp) toggleFullScreenCmd(evt *tcell.EventKey) *tcell.EventKey {
	t.fullScreen = !t.fullScreen
	t.SetFullScreen(t.fullScreen)
	return evt
}
