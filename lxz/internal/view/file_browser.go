/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/ui"
	"lxz/internal/ui/dialog"
	"lxz/internal/view/base"
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
}

func (_this *FileBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyE:         ui.NewKeyAction("Edit", _this.fileCtrl, true),
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		tcell.KeyCtrlN:  ui.NewKeyAction("Create File", _this.fileCtrl, true),
		tcell.KeyCtrlD:  ui.NewKeyAction("Delete File", _this.fileCtrl, true),
		tcell.KeyCtrlR:  ui.NewKeyAction("Rename File", _this.fileCtrl, true),
		tcell.KeyEscape: ui.NewKeyAction("Quit FullScreen", _this.ToggleFullScreenCmd, false),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
		tcell.KeyEnter:  ui.NewKeyAction("Preview", _this.TabFocusChange, true),
		tcell.KeyLeft:   ui.NewKeyAction("Focus Change", _this.TabFocusChange, false),
		tcell.KeyRight:  ui.NewKeyAction("Focus Change", _this.TabFocusChange, false),
	})
}

func (_this *FileBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

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
	_this.tree.SetBorderColor(base.ActiveBorderColor)
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
						if err != nil {
							return
						}
						if fi.IsDir() {
							// 如果是目录，清空预览内容
							_this.preview.SetText("[yellow]请选择一个文件进行预览")
							_this.preview.SetTitle(filepath.Base(path))
							return
						}
						// 读取文件大小,只预览文本文件
						if fi.Size() > 1024*1024*100 { // 大于1MB不预览
							_this.preview.SetText("[red]文件过大，无法预览, 仅支持预览100MB以下的文本文件")
							_this.preview.SetTitle(filepath.Base(path))
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

		_, err := os.Stat(path)
		if err != nil {
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
		SetDynamicColors(true).
		SetWordWrap(true).
		SetBorder(true).
		SetTitle("")
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
		_this.tree.SetBorderColor(base.InactiveBorderColor)
		_this.preview.SetBorderColor(base.ActiveBorderColor)
	} else {
		_this.app.UI.SetFocus(_this.tree)
		_this.tree.SetBorderColor(base.ActiveBorderColor)
		_this.preview.SetBorderColor(base.InactiveBorderColor)
	}
	return nil
}

// fileCtrl
func (_this *FileBrowser) fileCtrl(event *tcell.EventKey) *tcell.EventKey {
	node := _this.tree.GetCurrentNode()
	if node == nil {
		_this.app.UI.Flash().Warn("请先选中文件或目录")
		return nil
	}
	ref := node.GetReference()
	if ref == nil {
		_this.app.UI.Flash().Warn("请先选中文件或目录")
		return nil
	}

	path := ref.(string)
	info, err := os.Stat(path)
	if err != nil {
		_this.app.UI.Flash().Err(fmt.Errorf("读取文件信息失败: %v", err))
		return nil
	}

	slog.Info("[fileCtrl] ", "path", path, "isDir", info.IsDir(), "node", node.GetText())

	if event.Key() == tcell.KeyCtrlD {
		slog.Info("Will delete file", "path", path)
		_this.deleteFileModel(node, path)
		return nil
	} else if event.Key() == tcell.KeyCtrlN {
		slog.Info("Will create a new file", "path", path)
		if info.IsDir() {
			// 在此目录下新建文件
			_this.createFileModel(node, path)
		} else {
			// 基于当前文件所在的目录新建文件
			_this.createFileModel(node.GetParentNode(), filepath.Dir(path))
		}
		return nil
	} else if event.Key() == tcell.KeyCtrlR {
		// 文件或者文件夹改名
		slog.Info("[fileCtrl] Rename File", "path", path, "isDir", info.IsDir(), "node", node.GetText())
		if info.IsDir() {

		} else {
			_this.renameFileModel(node.GetParentNode(), path)
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
				_this.app.UI.Flash().Err(fmt.Errorf("open editor failed"))
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

func (_this *FileBrowser) createFileModel(node *tview.TreeNode, path string) {
	opts := dialog.CreateFileOpts{
		Title:    "Create New File",
		Message:  "Enter the name of the new file:",
		FileName: "",
		Ack: func(fileName string, isDir bool) bool {
			slog.Info("[文件新建] ", "path", path, "fileName", fileName, "isDir", isDir)
			if fileName == "" {
				_this.app.UI.Flash().Warn("name cannot be empty")
				return false
			}
			newFilePath := filepath.Join(path, fileName)
			// 检查文件是否已存在
			if _, err := os.Stat(newFilePath); err == nil {
				_this.app.UI.Flash().Warn(fmt.Sprintf("[red]<%s>[-] has existed", fileName))
				return false
			}
			if isDir {
				// 创建新目录
				err := os.Mkdir(newFilePath, 0755)
				if err != nil {
					_this.app.UI.Flash().Err(fmt.Errorf("[red]<%s>[-] create failed", fileName))
					return false
				}
			} else {
				// 创建新文件
				file, err := os.Create(newFilePath)
				if err != nil {
					_this.app.UI.Flash().Err(fmt.Errorf("[red]<%s>[-] create failed", fileName))
					return false
				}
				defer file.Close()
			}
			_this.app.UI.Flash().Info(fmt.Sprintf("[red]<%s>[-] create success", fileName))

			_this.addChildren(node, path)
			return true
		},
		Cancel: func() {},
	}
	dialog.ShowCreateFile(&config.Dialog{}, _this.app.Content.Pages, &opts)
}

func (_this *FileBrowser) renameFileModel(node *tview.TreeNode, path string) {
	opts := dialog.RenameFileOpts{
		Title:    "Rename File",
		Message:  "Rename the file to:",
		FileName: filepath.Base(path),
		Ack: func(newFileName string) bool {
			fileName := newFileName
			slog.Info("[文件名修改] ", "path", path, "fileName", fileName)
			if fileName == "" {
				_this.app.UI.Flash().Warn("文件名不能为空")
				return false
			}
			parentPath := filepath.Dir(path)
			newFilePath := filepath.Join(parentPath, fileName)
			// 检查文件是否已存在
			if _, err := os.Stat(newFilePath); err == nil {
				_this.app.UI.Flash().Warn(fmt.Sprintf("[red]<%s>[-] has existed", fileName))
				return false
			}

			err := os.Rename(path, newFilePath)
			if err != nil {
				_this.app.UI.Flash().
					Err(errors.New(fmt.Sprintf("[red]<%s>[-] rename failed", filepath.Base(path))))
				return false
			}
			_this.addChildren(node, path)
			return true
		},
		Cancel: func() {},
	}
	dialog.ShowRenameFile(&config.Dialog{}, _this.app.Content.Pages, &opts)
}

func (_this *FileBrowser) deleteFileModel(node *tview.TreeNode, path string) {
	// 删除文件
	opts := dialog.DeleteFileOpts{
		Title:        "Delete File",
		Message:      fmt.Sprintf("Are you sure you want to delete <%s> ?", filepath.Base(path)),
		FieldManager: "",
		Ack: func() bool {
			slog.Info("[文件删除] ", "path", path)
			// 删除文件
			err := os.Remove(path)
			if err != nil {
				_this.app.UI.Flash().
					Err(fmt.Errorf("[red]<%s>[-] delete failed", filepath.Base(path)))
				return false
			}
			// 获取父目录路径
			parentPath := filepath.Dir(path)
			_this.app.UI.Flash().
				Info(fmt.Sprintf("[red]<%s>[-] delete success", filepath.Base(path)))
			_this.addChildren(node.GetParentNode(), parentPath)
			return true
		},
		Cancel: func() {},
	}
	dialog.ShowDeleteFile(&config.Dialog{}, _this.app.Content.Pages, &opts)
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

func NewFileBrowser(app *App) *FileBrowser {
	var name = "File Browser"
	f := &FileBrowser{
		BaseFlex:         NewBaseFlex(name),
		app:              app,
		debounceInterval: 200 * time.Millisecond,
		stopDebounceCh:   make(chan struct{}),
		rootDir:          ".",
	}
	f.SetIdentifier(ui.FILE_BROWSER_ID)
	return f
}
