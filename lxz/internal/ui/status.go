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

type Status struct {
	*tview.TextView

	styles *config.Styles
	stack  *model.Stack
}

// NewStatus return a new view.
func NewStatus(styles *config.Styles) *Status {
	p := Status{
		stack:    model.NewStack(),
		styles:   styles,
		TextView: tview.NewTextView(),
	}

	p.SetBackgroundColor(tcell.ColorRed)
	fmt.Fprintf(&p, "789 ")

	return &p
}

// StylesChanged notifies skin changed.
func (c *Status) StylesChanged(s *config.Styles) {
	c.styles = s
	c.SetBackgroundColor(s.BgColor())
	c.refresh(c.stack.Flatten())
}

// StackPushed indicates a new item was added.
func (c *Status) StackPushed(comp model.Component) {
	c.stack.Push(comp)
	c.refresh(c.stack.Flatten())
}

// StackPopped indicates an item was deleted.
func (c *Status) StackPopped(_, _ model.Component) {
	c.stack.Pop()
	c.refresh(c.stack.Flatten())
}

// StackTop indicates the top of the stack.
func (*Status) StackTop(model.Component) {}

// Refresh updates view with new crumbs.
func (c *Status) refresh(crumbs []string) {
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
