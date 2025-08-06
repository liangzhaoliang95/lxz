/**
 * @author  zhaoliang.liang
 * @date  2025/7/31 16:23
 */

package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/apimachinery/pkg/labels"
	"log/slog"
	"lxz/internal/model"
	"lxz/internal/ui"
	"lxz/internal/view/base"
	"lxz/internal/view/cmd"
)

type BaseFlex struct {
	*tview.Flex
	actions    *ui.KeyActions
	identity   string // 用于标识该Flex的唯一性
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

func (_this *BaseFlex) ExtraHints() map[string]string {
	return nil
}

func (_this *BaseFlex) InCmdMode() bool {
	return false
}

func (_this *BaseFlex) SetFilter(s string) {
	//TODO implement me
	panic("implement me")
}

func (_this *BaseFlex) SetLabelSelector(selector labels.Selector) {
	//TODO implement me
	panic("implement me")
}

func (_this *BaseFlex) SetCommand(interpreter *cmd.Interpreter) {
	//TODO implement me
	panic("implement me")
}

func (_this *BaseFlex) ToggleFullScreenCmd(evt *tcell.EventKey) *tcell.EventKey {
	if ui.IsInputPrimitive(appUiInstance.GetFocus()) {
		return evt
	}

	if evt.Key() == tcell.KeyEscape {
		if _this.fullScreen {
			_this.fullScreen = false
			_this.SetFullScreen(false)
		}
	} else {
		_this.fullScreen = !_this.fullScreen
		_this.SetFullScreen(_this.fullScreen)
	}
	return nil
}

func (_this *BaseFlex) EmptyKeyEvent(evt *tcell.EventKey) *tcell.EventKey {
	return evt
}

func (_this *BaseFlex) Keyboard(evt *tcell.EventKey) *tcell.EventKey {
	slog.Info(
		"BaseFlex keyboard",
		"name",
		_this.name,
		"key",
		evt.Key(),
		"rune",
		evt.Rune(),
		"modifiers",
		evt.Modifiers(),
	)

	if a, ok := _this.Actions().Get(ui.AsKey(evt)); ok {
		return a.Action(evt)
	}

	return evt
}

func (_this *BaseFlex) SetIdentifier(identity string) {
	_this.identity = identity
}

func (_this *BaseFlex) GetIdentifier() string {
	return _this.identity
}

func NewBaseFlex(name string) *BaseFlex {
	b := &BaseFlex{
		Flex:    tview.NewFlex(),
		name:    name,
		actions: ui.NewKeyActions(),
	}
	b.SetDirection(tview.FlexColumn)
	b.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(fmt.Sprintf(" %s ", name)).
		SetTitleAlign(tview.AlignCenter)
	b.SetBorderColor(base.FlexBorderColor)
	return b
}
