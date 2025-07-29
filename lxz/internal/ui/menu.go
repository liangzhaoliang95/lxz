/**
 * @author  zhaoliang.liang
 * @date  2025/7/24 14:18
 */

package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"lxz/internal/config"
	"lxz/internal/model"
	"strings"
)

type Menu struct {
	*tview.TextView

	styles *config.Styles
	stack  *model.Stack
}

// NewMenu return a new view.
func NewMenu(styles *config.Styles) *Menu {
	p := Menu{
		stack:    model.NewStack(),
		styles:   styles,
		TextView: tview.NewTextView(),
	}
	p.SetBorder(true)
	p.SetBorderPadding(0, 0, 1, 1)
	p.SetBackgroundColor(tcell.ColorYellow)
	p.SetBorderColor(tcell.ColorRed)
	p.SetBorderAttributes(tcell.AttrDim)
	fmt.Fprintf(&p, "我是状态指示器")

	return &p
}

// StylesChanged notifies skin changed.
func (c *Menu) StylesChanged(s *config.Styles) {
	c.styles = s
	c.SetBackgroundColor(s.BgColor())
	c.refresh(c.stack.Flatten())
}

// StackPushed indicates a new item was added.
func (c *Menu) StackPushed(comp model.Component) {
	c.stack.Push(comp)
	c.refresh(c.stack.Flatten())
}

// StackPopped indicates an item was deleted.
func (c *Menu) StackPopped(_, _ model.Component) {
	c.stack.Pop()
	c.refresh(c.stack.Flatten())
}

// StackTop indicates the top of the stack.
func (*Menu) StackTop(model.Component) {}

// Refresh updates view with new crumbs.
func (c *Menu) refresh(crumbs []string) {
	c.Clear()
	last, bgColor := len(crumbs)-1, c.styles.Frame().Crumb.BgColor
	for i, crumb := range crumbs {
		if i == last {
			bgColor = c.styles.Frame().Crumb.ActiveColor
		}
		_, _ = fmt.Fprintf(c, "[%s:%s:b] <%s> [-:%s:-] ",
			c.styles.Frame().Crumb.FgColor,
			bgColor, strings.ReplaceAll(strings.ToLower(crumb), " ", ""),
			c.styles.Body().BgColor)
	}
}
