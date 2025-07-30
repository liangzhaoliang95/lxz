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

const maxRow = 7

var menuKey = []string{
	"<alt-1> | 项目release",
	"<alt-2> | project relea",
	"<alt-3> | project",
	"<alt-4> | project release",
	"<alt-5> | project release",
	"<alt-6> | pr",
	"<alt-7> | project release",
	"<alt-8> | project releaseas asd",
}

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
	p.SetBorderPadding(0, 1, 1, 1)
	//p.SetBorders(true)
	menuKeys := make([]string, 0, len(menuKey))
	menuNames := make([]string, 0, len(menuKey))

	for i := 0; i < len(menuKey); i++ {
		m := strings.Split(menuKey[i], "|")
		menuKeys = append(menuKeys, strings.Trim(m[0], " "))
		menuNames = append(menuNames, strings.Trim(m[1], " "))
	}

	row, col := 0, 0
	for i := 0; i < len(menuKeys); i++ {
		if i > 0 && i%maxRow == 0 {
			// 新列开头
			col += 2 // 2列为一个组：key + value
			row = 0
		}
		p.SetCell(row, col, &tview.TableCell{
			Text:  fmt.Sprintf("%s", menuKeys[i]),
			Color: tcell.ColorFuchsia,
			Align: tview.AlignLeft,
		})
		p.SetCell(row, col+1, &tview.TableCell{
			Text:  menuNames[i],
			Color: tcell.ColorDefault,
			Align: tview.AlignLeft,
		})
		row++
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
