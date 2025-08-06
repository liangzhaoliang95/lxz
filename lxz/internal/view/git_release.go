/**
 * @author  zhaoliang.liang
 * @date  2025/7/30 17:35
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/apimachinery/pkg/labels"
	"log/slog"
	"lxz/internal/model"
	"lxz/internal/ui"
	"lxz/internal/view/cmd"
)

// GitRelease 项目Release视图
type GitRelease struct {
	*BaseFlex
	actions    *ui.KeyActions
	name       string
	fullScreen bool
}

func (p *GitRelease) Name() string {
	return p.name
}

func (p *GitRelease) Init(ctx context.Context) error {
	p.bindKeys()
	// 用于初始化组件的边框、标题、快捷键等信息
	p.SetInputCapture(p.keyboard)
	//panic("TestComp init not implemented")
	return nil
}

func (p *GitRelease) Start() {
	slog.Info("GitRelease Start", "name", p.Name())
}

func (p *GitRelease) Stop() {
	slog.Info("GitRelease Stop", "name", p.Name())
}

func (p *GitRelease) Hints() model.MenuHints {
	return p.actions.Hints()
}

func (p *GitRelease) ExtraHints() map[string]string {
	return nil
}

func (p *GitRelease) InCmdMode() bool {
	return false
}

func (p *GitRelease) SetFilter(s string) {
	//TODO implement me
	panic("implement me")
}

func (p *GitRelease) SetLabelSelector(selector labels.Selector) {
	//TODO implement me
	panic("implement me")
}

func (p *GitRelease) SetCommand(interpreter *cmd.Interpreter) {
	//TODO implement me
	panic("implement me")
}

// Actions returns menu actions.
func (p *GitRelease) Actions() *ui.KeyActions {
	return p.actions
}

// Helpers

func (p *GitRelease) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info(
		"GitRelease keyboard",
		"key",
		evt.Key(),
		"rune",
		evt.Rune(),
		"modifiers",
		evt.Modifiers(),
	)

	if a, ok := p.actions.Get(ui.AsKey(evt)); ok {
		return a.Action(evt)
	}

	return evt
}

func (p *GitRelease) bindKeys() {
	p.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
		tcell.KeyEscape: ui.NewKeyAction("Toggle FullScreen", p.toggleFullScreenCmd, true),
	})
}

func (p *GitRelease) toggleFullScreenCmd(evt *tcell.EventKey) *tcell.EventKey {
	if evt.Key() == tcell.KeyEscape {
		if p.fullScreen {
			p.fullScreen = false
			p.SetFullScreen(false)
		}
	} else {
		p.fullScreen = !p.fullScreen
		p.SetFullScreen(p.fullScreen)
	}
	return evt
}

func NewGitRelease() *GitRelease {
	var name = "Git Release"
	tc := &GitRelease{
		BaseFlex: NewBaseFlex(name),
		name:     name,
		actions:  ui.NewKeyActions(),
	}

	tc.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(fmt.Sprintf(" %s ", name))
	tc.AddItem(tview.NewTextView().SetText(name), 0, 1, false)

	tc.SetIdentifier(ui.GIT_RELEASE_ID)
	return tc
}
