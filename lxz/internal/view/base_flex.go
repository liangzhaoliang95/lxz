/**
 * @author  zhaoliang.liang
 * @date  2025/7/31 16:23
 */

package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/model"
	"lxz/internal/ui"
)

type BaseFlex struct {
	*tview.Flex
	actions    *ui.KeyActions
	name       string
	fullScreen bool
}

func (_this *BaseFlex) Name() string {
	return _this.name
}

func (_this *BaseFlex) Actions() *ui.KeyActions {
	return _this.actions
}

func (_this *BaseFlex) Hints() model.MenuHints {
	return _this.actions.Hints()
}

func (_this *BaseFlex) toggleFullScreenCmd(evt *tcell.EventKey) *tcell.EventKey {
	if evt.Key() == tcell.KeyEscape {
		if _this.fullScreen {
			_this.fullScreen = false
			_this.SetFullScreen(false)
		}
	} else {
		_this.fullScreen = !_this.fullScreen
		_this.SetFullScreen(_this.fullScreen)
	}
	return evt
}

func (_this *BaseFlex) emptyKeyEvent(evt *tcell.EventKey) *tcell.EventKey {
	return evt
}

func (_this *BaseFlex) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info(
		"BaseFlex keyboard",
		"key",
		evt.Key(),
		"rune",
		evt.Rune(),
		"modifiers",
		evt.Modifiers(),
	)

	if a, ok := _this.actions.Get(ui.AsKey(evt)); ok {
		return a.Action(evt)
	}

	return evt
}

func newBaseFlex(name string) *BaseFlex {
	b := &BaseFlex{
		Flex:    tview.NewFlex(),
		name:    name,
		actions: ui.NewKeyActions(),
	}

	b.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrNone).
		SetTitle(name)

	return b
}
