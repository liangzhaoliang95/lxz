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
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const (
	menuIndexFmt = " [key:-:b]<%d> [fg:-:fgstyle]%s "
	maxRows      = 3
)

var menuRX = regexp.MustCompile(`\d`)

// SubMenu 用于显示各个功能页面的快捷键描述
type SubMenu struct {
	*tview.Table

	styles *config.Styles
}

// NewSubMenu return a new view.
func NewSubMenu(styles *config.Styles) *SubMenu {
	p := SubMenu{
		styles: styles,
		Table:  tview.NewTable(),
	}
	p.SetBorder(true)
	p.SetBorderPadding(0, 0, 0, 0)
	p.SetBackgroundColor(tcell.ColorBlack)
	p.SetBorderColor(tcell.ColorOrange)
	// 虚线边框
	p.SetBorderAttributes(tcell.AttrBold)
	p.SetTitle(" Sub Menu ")
	p.SetTitleAlign(tview.AlignCenter)
	p.SetTitleColor(tcell.ColorBlue)

	slog.Info("LXZ SubMenu NewSubMenu Done")
	return &p
}

// StylesChanged notifies skin changed.
func (_this *SubMenu) StylesChanged(s *config.Styles) {
	return
	_this.styles = s
	_this.SetBackgroundColor(s.BgColor())
}

// StackPushed indicates a new item was added.
func (_this *SubMenu) StackPushed(comp model.Component) {
	_this.HydrateMenu(comp.Hints())
}

func (_this *SubMenu) hasDigits(hh model.MenuHints) bool {
	for _, h := range hh {
		if !h.Visible {
			continue
		}
		if menuRX.MatchString(h.Mnemonic) {
			return true
		}
	}
	return false
}

func (_this *SubMenu) buildMenuTable(
	hh model.MenuHints,
	table []model.MenuHints,
	colCount int,
) [][]string {
	var row, col int
	firstCmd := true
	maxKeys := make([]int, colCount)
	for _, h := range hh {
		if !h.Visible {
			continue
		}

		if !menuRX.MatchString(h.Mnemonic) && firstCmd {
			row, col, firstCmd = 0, col+1, false
			if table[0][0].IsBlank() {
				col = 0
			}
		}
		if maxKeys[col] < len(h.Mnemonic) {
			maxKeys[col] = len(h.Mnemonic)
		}
		table[row][col] = h
		row++
		if row >= maxRows {
			row, col = 0, col+1
		}
	}

	out := make([][]string, len(table))
	for r := range out {
		out[r] = make([]string, len(table[r]))
	}
	_this.layout(table, maxKeys, out)

	return out
}

func (_this *SubMenu) layout(table []model.MenuHints, mm []int, out [][]string) {
	for r := range table {
		for c := range table[r] {
			out[r][c] = _this.formatMenu(table[r][c], mm[c])
		}
	}
}

func (_this *SubMenu) formatMenu(h model.MenuHint, size int) string {
	if h.Mnemonic == "" || h.Description == "" {
		return ""
	}
	styles := _this.styles.Frame()
	i, err := strconv.Atoi(h.Mnemonic)
	if err == nil {
		return formatNSMenu(i, h.Description, &styles)
	}

	return formatPlainMenu(h, size, &styles)
}

func (_this *SubMenu) HydrateMenu(hh model.MenuHints) {
	_this.Clear()
	sort.Sort(hh)

	table := make([]model.MenuHints, maxRows+1)
	colCount := (len(hh) / maxRows) + 1
	if _this.hasDigits(hh) {
		colCount++
	}
	for row := range maxRows {
		table[row] = make(model.MenuHints, colCount)
	}
	t := _this.buildMenuTable(hh, table, colCount)

	for row := range t {
		for col := range len(t[row]) {
			c := tview.NewTableCell(t[row][col])
			if t[row][col] == "" {
				c = tview.NewTableCell("")
			}
			c.SetBackgroundColor(tcell.ColorBlack)
			_this.SetCell(row, col, c)
		}
	}
}

// StackPopped indicates an item was deleted.
func (_this *SubMenu) StackPopped(_, comp model.Component) {
	if comp != nil {
		_this.HydrateMenu(comp.Hints())
	} else {
		_this.Clear()
	}

}

// StackTop indicates the top of the stack.
func (*SubMenu) StackTop(model.Component) {}

func keyConv(s string) string {
	if s == "" || !strings.Contains(s, "alt") {
		return s
	}
	if runtime.GOOS != "darwin" {
		return s
	}

	return strings.Replace(s, "alt", "opt", 1)
}

func ToMnemonic(s string) string {
	if s == "" {
		return s
	}

	return "<" + keyConv(strings.ToLower(s)) + ">"
}

func formatNSMenu(i int, name string, styles *config.Frame) string {
	fmat := strings.Replace(menuIndexFmt, "[key", "["+styles.Menu.NumKeyColor.String(), 1)
	fmat = strings.ReplaceAll(fmat, ":bg:", ":"+styles.Title.BgColor.String()+":")
	fmat = strings.Replace(fmat, "[fg", "["+styles.Menu.FgColor.String(), 1)
	fmat = strings.Replace(fmat, "fgstyle]", styles.Menu.FgStyle.ToShortString()+"]", 1)

	return fmt.Sprintf(fmat, i, name)
}

func formatPlainMenu(h model.MenuHint, size int, styles *config.Frame) string {
	menuFmt := " [key:-:b]%-" + strconv.Itoa(size+2) + "s [fg:-:fgstyle]%s "
	fmat := strings.Replace(menuFmt, "[key", "["+styles.Menu.KeyColor.String(), 1)
	fmat = strings.Replace(fmat, "[fg", "["+styles.Menu.FgColor.String(), 1)
	fmat = strings.ReplaceAll(fmat, ":bg:", ":"+styles.Title.BgColor.String()+":")
	fmat = strings.Replace(fmat, "fgstyle]", styles.Menu.FgStyle.ToShortString()+"]", 1)

	return fmt.Sprintf(fmat, ToMnemonic(h.Mnemonic), h.Description)
}
