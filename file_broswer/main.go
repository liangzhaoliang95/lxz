package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// 打开系统默认编辑器
func openSystemEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// fallback 到常见编辑器
		for _, candidate := range []string{"nano", "vim", "vi", "code"} {
			if _, err := exec.LookPath(candidate); err == nil {
				editor = candidate
				break
			}
		}
	}
	if editor == "" {
		return fmt.Errorf("未设置 $EDITOR，且未找到可用编辑器（如 vim/nano）")
	}
	// 检查命令是否存在
	if _, err := exec.LookPath(editor); err != nil {
		return fmt.Errorf("未设置 $EDITOR，且未找到可用编辑器（如 vim/nano）")
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var (
		lastFocusedPath  string
		lastFocusedAt    time.Time
		debounceInterval = 200 * time.Millisecond
		stopDebounceCh   = make(chan struct{})
	)

	rootDir := "/Users/liang" // 用当前路径，可以替换为"/Users/liang"等绝对路径

	// 文件树节点
	rootNode := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed).
		SetReference(rootDir)

	tree := tview.NewTreeView().
		SetRoot(rootNode).
		SetCurrentNode(rootNode)

	// 文件预览区域
	preview := tview.NewTextView()
	preview.
		SetDynamicColors(false).
		SetWordWrap(true).
		SetBorder(true).
		SetTitle("预览")

	// 编辑区域（进入编辑模式时显示）
	editor := tview.NewTextView().SetChangedFunc(func() {

	})
	editor.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("编辑中...(按 Esc 退出)")

	// 页面布局
	flex := tview.NewFlex().
		AddItem(tree, 40, 1, true).
		AddItem(preview, 0, 2, false)

	app := tview.NewApplication()

	// 添加节点的辅助函数
	addChildren := func(node *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			// 显示错误信息
			preview.SetText(fmt.Sprintf("[red]读取目录失败: %v", err))
			return
		}
		for _, file := range files {
			fullPath := filepath.Join(path, file.Name())
			child := tview.NewTreeNode(file.Name()).
				SetReference(fullPath).
				SetSelectable(true)
			if file.IsDir() {
				child.SetColor(tcell.ColorGreen)
			}
			node.AddChild(child)
		}
	}

	// 初始展开根目录
	addChildren(rootNode, rootDir)
	tree.SetChangedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		path := ref.(string)

		fi, err := os.Stat(path)
		if err != nil || fi.IsDir() {
			lastFocusedPath = ""
			return
		}

		// 设置当前选中和时间
		lastFocusedPath = path
		lastFocusedAt = time.Now()
	})
	// 节点选择时事件
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			preview.SetText("[red]请选择一个文件或目录")
			return
		}
		path := ref.(string)

		info, err := os.Stat(path)
		if err != nil {
			preview.SetText(fmt.Sprintf("[red]读取失败: %v", err))
			return
		}

		if info.IsDir() {
			if len(node.GetChildren()) == 0 {
				addChildren(node, path)
				node.SetExpanded(true) // ✅ 第一次加载子目录 -> 展开它
			} else {
				node.SetExpanded(!node.IsExpanded())
			}

		}
	})

	// 键盘事件捕获逻辑
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'e':
			node := tree.GetCurrentNode()
			if node == nil {
				return nil
			}
			ref := node.GetReference()
			if ref == nil {
				return nil
			}
			path := ref.(string)

			info, err := os.Stat(path)
			if err != nil || info.IsDir() {
				return nil
			}

			//🟡 暂停 tview UI 进入外部编辑器（阻塞直到编辑完成）
			app.Suspend(func() {
				err := openSystemEditor(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[编辑器打开失败] %v\n", err)
					fmt.Println("按回车继续...")
					fmt.Scanln()
				}
			})

			// 🟢 回到 TUI，自动刷新预览内容
			data, err := os.ReadFile(path)
			if err != nil {
				preview.SetText("[red]文件读取失败")
			} else {
				preview.SetTitle("预览: " + filepath.Base(path))
				preview.SetText(string(data))
				// 👇自动切焦点到右侧预览区域
				app.SetFocus(preview)
			}
		}
		return event
	})

	// 编辑器按 Esc 退出编辑
	editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.SetRoot(flex, true).SetFocus(tree)
			return nil
		}
		return event
	})

	var inPreviewFocus = false // 当前焦点标记

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if inPreviewFocus {
				app.SetFocus(tree)
			} else {
				app.SetFocus(preview)
			}
			inPreviewFocus = !inPreviewFocus
			return nil
		}
		return event
	})

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 如果停留时间 >= debounceInterval，触发预览
				if lastFocusedPath != "" && time.Since(lastFocusedAt) >= debounceInterval {
					path := lastFocusedPath

					// 读取和刷新 UI 必须在主线程
					app.QueueUpdateDraw(func() {
						fi, err := os.Stat(path)
						if err != nil || fi.IsDir() {
							return
						}
						data, err := os.ReadFile(path)
						if err != nil {
							preview.SetText(fmt.Sprintf("[red]无法读取文件: %v", err))
							return
						}
						preview.SetTitle("预览: " + filepath.Base(path))
						preview.SetText(string(data))
					})

					// 渲染完重置，避免重复渲染
					lastFocusedPath = ""
				}

			case <-stopDebounceCh:
				return
			}
		}
	}()

	// 启动程序
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
