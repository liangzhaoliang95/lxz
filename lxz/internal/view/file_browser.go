/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/ui"
	"lxz/internal/ui/dialog"
	"lxz/internal/view/cmd"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type FileBrowser struct {
	*BaseFlex
	app              *App
	rootDir          string
	rootNode         *tview.TreeNode // 根目录节点
	preview          *tview.TextView
	tree             *tview.TreeView
	lastFocusedPath  string
	lastFocusedAt    time.Time
	debounceInterval time.Duration
	stopDebounceCh   chan struct{}
	createFileModel  *tview.Modal
	deleteFileModel  *tview.Modal
	renameFileModel  *tview.Modal
}

func (_this *FileBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyE:         ui.NewKeyAction("Edit", _this.fileCtrl, true),
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.toggleFullScreenCmd, true),
		tcell.KeyCtrlN:  ui.NewKeyAction("Create File", _this.fileCtrl, true),
		tcell.KeyCtrlD:  ui.NewKeyAction("Delete File", _this.fileCtrl, true),
		tcell.KeyEscape: ui.NewKeyAction("Quit FullScreen", _this.toggleFullScreenCmd, false),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
		tcell.KeyEnter:  ui.NewKeyAction("Confirm", _this.TabFocusChange, true),
		tcell.KeyLeft:   ui.NewKeyAction("Focus Change", _this.TabFocusChange, false),
		tcell.KeyRight:  ui.NewKeyAction("Focus Change", _this.TabFocusChange, false),
	})
}

func (_this *FileBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.keyboard)

	// 组件初始化
	_this.initRootNode()
	_this.initTree()
	_this.initPreview()

	_this.AddItem(_this.tree, 30, 1, true)
	_this.AddItem(_this.preview, 0, 2, false)

	// 初始展开根目录
	_this.addChildren(_this.rootNode, _this.rootDir)

	return nil
}

func (_this *FileBrowser) Start() {
	// ✅ 设置默认边框颜色 + 焦点 + 强制刷新
	_this.tree.SetBorderColor(activeBorderColor)
	go func(a *FileBrowser) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 如果停留时间 >= debounceInterval，触发预览
				if _this.lastFocusedPath != "" &&
					time.Since(_this.lastFocusedAt) >= _this.debounceInterval {
					path := _this.lastFocusedPath

					// 读取和刷新 UI 必须在主线程
					_this.app.UI.QueueUpdateDraw(func() {
						fi, err := os.Stat(path)
						if err != nil || fi.IsDir() {
							return
						}
						data, err := os.ReadFile(path)
						if err != nil {
							_this.preview.SetText(fmt.Sprintf("[red]无法读取文件: %v", err))
							return
						}
						_this.preview.SetTitle(filepath.Base(path))
						_this.preview.SetText(string(data))
					})

					// 渲染完重置，避免重复渲染
					_this.lastFocusedPath = ""
				}

			case <-_this.stopDebounceCh:
				return
			}
		}
	}(_this)
}

func (_this *FileBrowser) Stop() {
	//if _this.stopDebounceCh != nil {
	//	close(_this.stopDebounceCh) // 停止防抖协程
	//}
}

func (_this *FileBrowser) ExtraHints() map[string]string {
	//TODO implement me
	panic("implement me")
}

func (_this *FileBrowser) InCmdMode() bool {
	//TODO implement me
	panic("implement me")
}

func (_this *FileBrowser) SetFilter(s string) {
	//TODO implement me
	panic("implement me")
}

func (_this *FileBrowser) SetLabelSelector(selector labels.Selector) {
	//TODO implement me
	panic("implement me")
}

func (_this *FileBrowser) SetCommand(interpreter *cmd.Interpreter) {
	//TODO implement me
	panic("implement me")
}

// helpers
func (_this *FileBrowser) addChildren(node *tview.TreeNode, path string) {
	// 清空旧的 Children
	node.ClearChildren()
	files, err := os.ReadDir(path)
	if err != nil {
		// 显示错误信息
		_this.preview.SetText(fmt.Sprintf("[red]读取目录失败: %v", err))
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

// initRootNode
func (_this *FileBrowser) initRootNode() {
	// 文件树节点
	_this.rootNode = tview.NewTreeNode(_this.rootDir).
		SetColor(tcell.ColorRed).
		SetReference(_this.rootDir)
}

// initTree
func (_this *FileBrowser) initTree() {

	_this.tree = tview.NewTreeView().
		SetRoot(_this.rootNode).
		SetCurrentNode(_this.rootNode)
	_this.tree.SetBorder(true)

	_this.tree.SetChangedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		path := ref.(string)

		fi, err := os.Stat(path)
		if err != nil || fi.IsDir() {
			_this.lastFocusedPath = ""
			return
		}

		// 设置当前选中和时间
		_this.lastFocusedPath = path
		_this.lastFocusedAt = time.Now()
	})
	// 节点选择时事件
	_this.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			_this.preview.SetText("[red]请选择一个文件或目录")
			return
		}
		path := ref.(string)

		info, err := os.Stat(path)
		if err != nil {
			_this.preview.SetText(fmt.Sprintf("[red]读取失败: %v", err))
			return
		}

		if info.IsDir() {
			if len(node.GetChildren()) == 0 {
				_this.addChildren(node, path)
				node.SetExpanded(true) // ✅ 第一次加载子目录 -> 展开它
			} else {
				node.SetExpanded(!node.IsExpanded())
			}

		}
	})

}

// initPreview
func (_this *FileBrowser) initPreview() {
	// 文件预览区域
	_this.preview = tview.NewTextView()
	_this.preview.
		SetDynamicColors(false).
		SetWordWrap(true).
		SetBorder(true).
		SetTitle("")
}

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

func (_this *FileBrowser) TabFocusChange(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTAB {

	} else if event.Key() == tcell.KeyEnter {

		ref := _this.tree.GetCurrentNode().GetReference()
		if ref == nil {
			_this.preview.SetText("[red]请选择一个文件或目录")
			return nil
		}
		path := ref.(string)

		info, err := os.Stat(path)
		if err != nil {
			_this.preview.SetText(fmt.Sprintf("[red]读取失败: %v", err))
			return nil
		}

		if info.IsDir() {
			return event
		}

		if _this.app.UI.GetFocus() == _this.preview {
			return nil // 已经在 preview 上了
		}
	} else if event.Key() == tcell.KeyLeft {
		if _this.app.UI.GetFocus() == _this.tree {
			return nil // 已经在 tree 上了
		}
	} else if event.Key() == tcell.KeyRight {
		if _this.app.UI.GetFocus() == _this.preview {
			return nil // 已经在 preview 上了
		}
	}
	if _this.app.UI.GetFocus() == _this.tree {
		_this.app.UI.SetFocus(_this.preview)
		_this.tree.SetBorderColor(inactiveBorderColor)
		_this.preview.SetBorderColor(activeBorderColor)
	} else {
		_this.app.UI.SetFocus(_this.tree)
		_this.tree.SetBorderColor(activeBorderColor)
		_this.preview.SetBorderColor(inactiveBorderColor)
	}
	return nil
}

// fileCtrl
func (_this *FileBrowser) fileCtrl(event *tcell.EventKey) *tcell.EventKey {
	slog.Info("[fileCtrl] ", "event", event.Key())
	node := _this.tree.GetCurrentNode()
	if node == nil {
		slog.Info("[fileCtrl失败] ", "err", "未选中文件或目录")
		return nil
	}
	ref := node.GetReference()
	if ref == nil {
		slog.Info("[fileCtrl失败] ", "err", "未选中文件或目录")
		return nil
	}

	path := ref.(string)
	info, err := os.Stat(path)
	if err != nil {
		slog.Info("[fileCtrl失败] ", "err", err)
		return nil
	}

	// 文件的新建 删除
	if event.Key() == tcell.KeyCtrlD {
		// 删除文件
		slog.Info("Will delete file", "path", path)
		if info.IsDir() {
		}
		return nil
	} else if event.Key() == tcell.KeyCtrlN {
		slog.Info("Will create a new file", "path", path)
		// 新建文件
		if info.IsDir() {
			// 在此目录下新建文件
			opts := dialog.CreateFileOpts{
				Title:        "Create New File",
				Message:      "Enter the name of the new file:",
				FieldManager: "",
				Ack: func(opts *metav1.PatchOptions) bool {
					slog.Info("[文件新建] \n", "path", path)
					fileName := opts.FieldManager
					if fileName == "" {
						slog.Info("[文件新建失败] ", "err", "文件名不能为空")
						return false
					}
					newFilePath := filepath.Join(path, fileName)
					// 检查文件是否已存在
					if _, err := os.Stat(newFilePath); err == nil {
						slog.Info("[文件新建失败] ", "err", "文件已存在")
						return false
					}
					// 创建新文件
					file, err := os.Create(newFilePath)
					if err != nil {
						slog.Error("[文件新建失败] ", "err", err)
						return false
					}
					defer file.Close()
					_this.addChildren(node, path)
					return true
				},
				Cancel: func() {},
			}
			dialog.ShowCreateFile(&config.Dialog{}, _this.app.Content.Pages, &opts)
		} else {
			// TODO 使用flash组件通知不允许
		}
		return nil
	}

	switch event.Rune() {
	case 'e':
		if info.IsDir() {
			return nil
		}

		//🟡 暂停 tview UI 进入外部编辑器（阻塞直到编辑完成）
		_this.app.UI.Suspend(func() {
			err := openSystemEditor(path)
			if err != nil {
				slog.Error("[编辑器打开失败] %v\n", "err", os.Stderr)
				fmt.Scanln()
			}
		})

		// 🟢 回到 TUI，自动刷新预览内容
		data, err := os.ReadFile(path)
		if err != nil {
			_this.preview.SetText("[red]文件读取失败")
		} else {
			_this.preview.SetTitle(filepath.Base(path))
			_this.preview.SetText(string(data))
			// 👇自动切焦点到右侧预览区域
			_this.app.UI.SetFocus(_this.preview)
		}
	}
	return event
}

func NewFileBrowser(app *App) *FileBrowser {
	var name = "File Browser"
	f := &FileBrowser{
		BaseFlex:         newBaseFlex(name),
		app:              app,
		debounceInterval: 200 * time.Millisecond,
		stopDebounceCh:   make(chan struct{}),
		rootDir:          ".",
	}
	f.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(fmt.Sprintf(" %s ", name))

	return f
}
