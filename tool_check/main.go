package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
)

func main() {
	// 要检查的命令列表
	commands := []string{"git", "docker", "kubectl", "helm", "make", "node", "go", "npm"}

	// 创建 tview 应用实例
	app := tview.NewApplication()

	// 创建表格
	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0) // 固定顶部 header 行

	// 设置表头
	headers := []string{"状态", "命令", "路径"}
	for i, h := range headers {
		table.SetCell(0, i,
			tview.NewTableCell("[::b]"+h).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	// 填充命令状态数据
	for row, cmd := range commands {
		// exec.LookPath 会返回命令的路径或错误
		path, err := exec.LookPath(cmd)
		var status string
		var color string
		if err == nil {
			status = "✔"
			color = "green"
		} else {
			status = "✘"
			color = "red"
			path = "N/A"
		}

		// 设置单元格
		table.SetCell(row+1, 0, tview.NewTableCell("["+color+"]"+status+"[white]").
			SetAlign(tview.AlignCenter))
		table.SetCell(row+1, 1, tview.NewTableCell(cmd).
			SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 2, tview.NewTableCell(path).
			SetAlign(tview.AlignLeft))
	}

	// 表格样式
	table.SetTitle("📦 系统命令检测").
		SetBorder(true)

	// 支持 Esc 或 q/Q 退出
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			app.Stop()
			return nil
		}
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
			return nil
		}
		return event
	})

	// 启动应用
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
