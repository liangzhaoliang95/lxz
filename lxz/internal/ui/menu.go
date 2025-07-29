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
)

type Menu struct {
	*tview.Table

	styles *config.Styles
}

// NewMenu return a new view.
func NewMenu(styles *config.Styles) *Menu {
	p := Menu{
		styles: styles,
		Table:  tview.NewTable(),
	}
	p.SetFixed(1, 1)
	p.SetBorders(true)

	for i := 0; i < 10; i++ {
		p.SetCell(i, 0, &tview.TableCell{
			Text:  fmt.Sprintf("Key %d", i+1),
			Color: tcell.ColorBlue,
			Align: tview.AlignLeft,
		})
		p.SetCell(i, 1, &tview.TableCell{
			Text:  fmt.Sprintf("Item %d", i+1),
			Color: tcell.ColorDefault,
			Align: tview.AlignLeft,
		})
	}

	return &p
}

// StylesChanged notifies skin changed.
func (c *Menu) StylesChanged(s *config.Styles) {
	c.styles = s
	c.SetBackgroundColor(s.BgColor())
	c.refresh([]string{})
}

// StackPushed indicates a new item was added.
func (c *Menu) StackPushed(comp model.Component) {

}

// StackPopped indicates an item was deleted.
func (c *Menu) StackPopped(_, _ model.Component) {

}

// StackTop indicates the top of the stack.
func (*Menu) StackTop(model.Component) {}

// Refresh updates view with new crumbs.
func (c *Menu) refresh(crumbs []string) {
	c.Clear()

}
