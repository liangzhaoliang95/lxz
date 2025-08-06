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
	"log/slog"
	"lxz/internal/ui"
)

// RedisBrowser 项目Release视图
type RedisBrowser struct {
	*BaseFlex
	actions    *ui.KeyActions
	name       string
	fullScreen bool
}

func (p *RedisBrowser) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info(
		"RedisBrowser keyboard",
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

func (p *RedisBrowser) bindKeys() {
	p.Actions().Bulk(ui.KeyMap{})
}

func (p *RedisBrowser) Init(ctx context.Context) error {
	p.bindKeys()
	// 用于初始化组件的边框、标题、快捷键等信息
	p.SetInputCapture(p.keyboard)
	//panic("TestComp init not implemented")
	return nil
}

func (p *RedisBrowser) Start() {
	slog.Info("RedisBrowser Start", "name", p.Name())
}

func (p *RedisBrowser) Stop() {
	slog.Info("RedisBrowser Stop", "name", p.Name())
}

// Helpers

func NewRedisBrowser() *RedisBrowser {
	var name = "Redis Browser"
	tc := &RedisBrowser{
		BaseFlex: NewBaseFlex(name),
		name:     name,
		actions:  ui.NewKeyActions(),
	}

	tc.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(fmt.Sprintf(" %s ", name))
	tc.AddItem(tview.NewTextView().SetText(name), 0, 1, false)

	tc.SetIdentifier(ui.REDIS_BROWSER_ID)
	return tc
}
