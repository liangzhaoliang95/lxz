/**
 * @author  zhaoliang.liang
 * @date  2025/7/30 17:35
 */

package view

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/apimachinery/pkg/labels"
	"log/slog"
	"lxz/internal/model"
	"lxz/internal/ui"
	"lxz/internal/view/cmd"
)

// ProjectRelease 项目Release视图
type ProjectRelease struct {
	*tview.Flex
	actions    *ui.KeyActions
	name       string
	fullScreen bool
}

func (p *ProjectRelease) Name() string {
	return p.name
}

func (p *ProjectRelease) Init(ctx context.Context) error {
	p.bindKeys()
	// 用于初始化组件的边框、标题、快捷键等信息
	p.SetInputCapture(p.keyboard)
	//panic("TestComp init not implemented")
	return nil
}

func (p *ProjectRelease) Start() {
	slog.Info("ProjectRelease Start", "name", p.Name())
}

func (p *ProjectRelease) Stop() {
	slog.Info("ProjectRelease Stop", "name", p.Name())
}

func (p *ProjectRelease) Hints() model.MenuHints {
	return p.actions.Hints()
}

func (p *ProjectRelease) ExtraHints() map[string]string {
	return nil
}

func (p *ProjectRelease) InCmdMode() bool {
	return false
}

func (p *ProjectRelease) SetFilter(s string) {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectRelease) SetLabelSelector(selector labels.Selector) {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectRelease) SetCommand(interpreter *cmd.Interpreter) {
	//TODO implement me
	panic("implement me")
}

// Actions returns menu actions.
func (p *ProjectRelease) Actions() *ui.KeyActions {
	return p.actions
}

// Helpers

func (p *ProjectRelease) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info("ProjectRelease keyboard", "key", evt.Key(), "rune", evt.Rune(), "modifiers", evt.Modifiers())

	if a, ok := p.actions.Get(ui.AsKey(evt)); ok {
		return a.Action(evt)
	}

	return evt
}

func (p *ProjectRelease) bindKeys() {
	p.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
		ui.KeyA:         ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
		ui.KeyB:         ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
		ui.KeyC:         ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
		tcell.KeyEscape: ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
	})
}

func (p *ProjectRelease) toggleFullScreenCmd(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info("ProjectRelease toggleFullScreenCmd", "key", evt.Key(), "rune", evt.Rune(), "modifiers", evt.Modifiers())
	if evt.Key() == tcell.KeyEscape && p.fullScreen {
		p.fullScreen = false
		p.SetFullScreen(false)
		return evt
	}
	p.fullScreen = !p.fullScreen
	p.SetFullScreen(p.fullScreen)
	return evt
}

func NewProjectRelease() *ProjectRelease {
	tc := &ProjectRelease{
		Flex:    tview.NewFlex(),
		name:    "projectRelease",
		actions: ui.NewKeyActions(),
	}

	tc.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle("project Release")
	tc.AddItem(tview.NewTextView().SetText("This is a projectRelease"), 0, 1, false)
	return tc
}
