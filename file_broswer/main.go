package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false).SetBorders(false)

	currentPath := "/Users/liang"

	// 外层的 frame，用于显示当前路径标题
	frame := tview.NewFrame(table).
		SetBorders(1, 1, 1, 1, 2, 2)

	// 更新 table 的内容和路径标题
	var updateTable func(string)
	updateTable = func(path string) {
		table.Clear()

		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Printf("读取目录失败: %v", err)
			return
		}

		// 更新当前路径变量
		currentPath = path

		// 更新标题信息
		frame.Clear()
		frame.AddText("终端文件浏览器 - ↑↓选择, Enter进入, Backspace返回, q退出", true, tview.AlignCenter, tcell.ColorGreen).
			AddText("当前路径: "+currentPath, false, tview.AlignLeft, tcell.ColorYellow)

		// 设置 ".." 返回上级
		table.SetCell(0, 0, tview.NewTableCell("..").
			SetTextColor(tcell.ColorYellow).
			SetSelectable(true))

		// 排序，保持目录在上面
		sort.SliceStable(files, func(i, j int) bool {
			if files[i].IsDir() && !files[j].IsDir() {
				return true
			}
			return files[i].Name() < files[j].Name()
		})

		for i, file := range files {
			name := file.Name()
			cell := tview.NewTableCell(name)
			if file.IsDir() {
				cell.SetTextColor(tcell.ColorSkyblue)
			} else {
				cell.SetTextColor(tcell.ColorWhite)
			}
			table.SetCell(i+1, 0, cell)
		}
	}

	// 初始化第一次加载
	updateTable(currentPath)

	table.SetSelectedFunc(func(row, column int) {
		cell := table.GetCell(row, column)
		if cell == nil {
			return
		}
		selected := cell.Text

		var selectedPath string
		if selected == ".." {
			selectedPath = filepath.Dir(currentPath)
		} else {
			selectedPath = filepath.Join(currentPath, selected)
		}

		// 判断是不是目录
		info, err := os.Stat(selectedPath)
		if err != nil {
			log.Printf("无法读取 %s: %v", selectedPath, err)
			return
		}
		if info.IsDir() {
			updateTable(selectedPath)
			app.Draw()
		}
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBackspace2, tcell.KeyBackspace:
			updateTable(filepath.Dir(currentPath))
			app.Draw()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' || event.Rune() == 'Q' {
				app.Stop()
				return nil
			}
		}
		return event
	})

	// 启动应用
	if err := app.SetRoot(frame, true).Run(); err != nil {
		panic(err)
	}
}
