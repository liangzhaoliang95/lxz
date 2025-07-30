/**
 * @author  zhaoliang.liang
 * @date  2025/7/24 14:18
 */

package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/model"
	"strings"
)

// SubMenu 用于显示各个功能页面的快捷键描述
type SubMenu struct {
	*tview.TextView

	styles *config.Styles
}

// NewSubMenu return a new view.
func NewSubMenu(styles *config.Styles) *SubMenu {
	p := SubMenu{
		styles:   styles,
		TextView: tview.NewTextView(),
	}
	p.SetBorder(true)
	p.SetBorderPadding(0, 0, 1, 1)
	p.SetBackgroundColor(tcell.ColorYellow)
	p.SetBorderColor(tcell.ColorRed)
	p.SetBorderAttributes(tcell.AttrDim)
	p.SetChangedFunc(func() {
		// 每次内容发生变化后，自动滚动到底部
		p.ScrollToEnd()
	})
	slog.Info("LXZ SubMenu NewSubMenu Done")
	return &p
}

// StylesChanged notifies skin changed.
func (c *SubMenu) StylesChanged(s *config.Styles) {
	c.styles = s
	c.SetBackgroundColor(s.BgColor())
}

// StackPushed indicates a new item was added.
func (c *SubMenu) StackPushed(comp model.Component) {
	c.HydrateMenu(comp.Hints())
	fmt.Fprintf(c, "SubMenu StackPushed %s \n", comp.Name())
}
func (c *SubMenu) HydrateMenu(hh model.MenuHints) {}

// StackPopped indicates an item was deleted.
func (c *SubMenu) StackPopped(_, comp model.Component) {
	if comp != nil {
		c.HydrateMenu(comp.Hints())
		fmt.Fprintf(c, "SubMenu StackPopped %s \n", comp.Name())
	} else {
		c.Clear()
	}

}

// StackTop indicates the top of the stack.
func (*SubMenu) StackTop(model.Component) {}

// Refresh updates view with new crumbs.
func (c *SubMenu) refresh(crumbs []string) {
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
