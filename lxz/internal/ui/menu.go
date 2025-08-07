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
	"lxz/internal/helper"
	"lxz/internal/model"
)

const maxRow = 4

const (
	SSH_CONNECT_ID    = "SSH_CONNECT"
	FILE_BROWSER_ID   = "FILE_BROWSER"
	DB_BROWSER_ID     = "DB_BROWSER"
	REDIS_BROWSER_ID  = "REDIS_BROWSER"
	DOCKER_BROWSER_ID = "DOCKER_BROWSER_ID"
)

var menuMap = map[string]map[string]string{
	"<F1>": {
		"name": "🖥️ SSH Connect",
		"id":   SSH_CONNECT_ID,
		"sort": "1",
	},
	"<F2>": {
		"name": "🗂️ File Browser",
		"id":   FILE_BROWSER_ID,
		"sort": "2",
	},
	"<F3>": {
		"name": "🎯 Redis Browser",
		"id":   REDIS_BROWSER_ID,
		"sort": "3",
	},
	"<F4>": {
		"name": "📊 DB Browser",
		"id":   DB_BROWSER_ID,
		"sort": "4",
	},
	"<F5>": {
		"name": "🐳 Docker Browser",
		"id":   DOCKER_BROWSER_ID,
		"sort": "5",
	},
}

type Menu struct {
	*tview.Table
	nowIdentifier string // 当前选中的组件标识符
	styles        *config.Styles
}

// NewMenu return a new view.
func NewMenu(styles *config.Styles) *Menu {
	p := Menu{
		styles: styles,
		Table:  tview.NewTable(),
	}
	p.SetFixed(1, 1)
	p.SetBorderPadding(0, 1, 1, 1)
	p.SetBackgroundColor(tcell.ColorBlack)
	//p.SetBorders(true)
	p.refresh("")

	return &p
}

// StylesChanged notifies skin changed.
func (c *Menu) StylesChanged(s *config.Styles) {
	c.styles = s
	c.SetBackgroundColor(s.BgColor())
	c.refresh("")
}

// StackPushed indicates a new item was added.
func (c *Menu) StackPushed(comp model.Component) {
	slog.Info("Menu StackPushed", "component", comp.GetIdentifier())
	if comp.GetIdentifier() != "" {
		c.nowIdentifier = comp.GetIdentifier()
	}
	c.refresh(c.nowIdentifier)

}

// StackPopped indicates an item was deleted.
func (c *Menu) StackPopped(_, _ model.Component) {

}

// StackTop indicates the top of the stack.
func (*Menu) StackTop(model.Component) {}

// Refresh updates view with new crumbs.
func (c *Menu) refresh(compId string) {
	c.Clear()
	menuKeys := make([]string, 0)

	for key, _ := range menuMap {
		menuKeys = append(menuKeys, key)
	}
	// 根据sort字段对菜单进行排序
	// 这里使用简单的冒泡排序，实际应用中可以使用更高效的排序算法
	for i := 0; i < len(menuKeys)-1; i++ {
		for j := 0; j < len(menuKeys)-i-1; j++ {
			if menuMap[menuKeys[j]]["sort"] > menuMap[menuKeys[j+1]]["sort"] {
				menuKeys[j], menuKeys[j+1] = menuKeys[j+1], menuKeys[j]
			}
		}
	}

	row, col := 0, 0
	for i := 0; i < len(menuKeys); i++ {
		if i > 0 && i%maxRow == 0 {
			// 新列开头
			col += 3 // 2列为一个组：key + value
			row = 0
		}
		item := menuMap[menuKeys[i]]
		id := item["id"]
		c.SetCell(row, col, &tview.TableCell{
			Text:            helper.If(id == compId, "👉", "  "),
			Color:           tcell.ColorGreen,
			Align:           tview.AlignLeft,
			BackgroundColor: tcell.ColorBlack,
		})

		c.SetCell(row, col+1, &tview.TableCell{
			Text:            fmt.Sprintf("%s", menuKeys[i]),
			Color:           tcell.ColorFuchsia,
			Align:           tview.AlignLeft,
			BackgroundColor: tcell.ColorBlack,
		})
		c.SetCell(row, col+2, &tview.TableCell{
			Text:            item["name"],
			Color:           tcell.ColorDefault,
			Align:           tview.AlignLeft,
			BackgroundColor: tcell.ColorBlack,
		})
		row++
	}
}
