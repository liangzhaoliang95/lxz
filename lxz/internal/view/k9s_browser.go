/**
 * @author  zhaoliang.liang
 * @date  2025/8/1 10:31
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/tview"
	"log/slog"
	"lxz/internal/helper"
	"lxz/internal/ui"
	"os"
)

type K9SBrowser struct {
	*BaseFlex
	app              *App
	k9sConfigTableUI *tview.Table
	selectConfigPath string // 当前选中的配置文件路径
}

func (_this *K9SBrowser) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		tcell.KeyEnter: ui.NewKeyAction("Launch K9s", _this.K9SShellIn, true),
	})
}

func (_this *K9SBrowser) Init(ctx context.Context) error {
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)
	_this.k9sConfigTableUI = tview.NewTable()
	_this.k9sConfigTableUI.SetBorder(false)
	_this.k9sConfigTableUI.SetBorders(false)
	_this.k9sConfigTableUI.SetTitle("k9s Browser")
	_this.k9sConfigTableUI.SetBorderPadding(1, 1, 2, 2)
	_this.k9sConfigTableUI.SetSelectable(true, false)
	// 配置回车函数
	_this.k9sConfigTableUI.SetSelectedFunc(func(row, column int) {

	})
	// 设置表格的选择模式
	_this.k9sConfigTableUI.SetSelectionChangedFunc(func(row, column int) {
		slog.Info("Selection changed", "row", row, "col", column)
		if row < 1 || row >= _this.k9sConfigTableUI.GetRowCount() {
			slog.Warn("Selection changed row is out of range", "row", row)
			return
		}
		_this.selectConfigPath = _this.k9sConfigTableUI.GetCell(row, 1).Text

	})
	_this._initHeader()
	_this.AddItem(_this.k9sConfigTableUI, 0, 1, true)
	return nil
}

func (_this *K9SBrowser) _initHeader() {
	// 初始化header
	_this.k9sConfigTableUI.SetCell(0, 0, tview.NewTableCell("Num").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
	_this.k9sConfigTableUI.SetCell(0, 1, tview.NewTableCell("Path").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(false))
}

func (_this *K9SBrowser) _refreshData() {
	list, err := _this.listDefaultK8sConfigFiles()
	if err != nil {
		_this.app.UI.Flash().Err(fmt.Errorf("failed to list k8s config files: %w", err))
		return
	}
	// 填充容器数据
	for i, ctr := range list {
		_this.k9sConfigTableUI.SetCell(i+1, 0, tview.NewTableCell(helper.IntToString(i)).
			SetReference(helper.IntToString(i)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		_this.k9sConfigTableUI.SetCell(i+1, 1, tview.NewTableCell(ctr).
			SetReference(ctr).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
	}

}

func (_this *K9SBrowser) Start() {
	_this._refreshData()
	_this.app.UI.SetFocus(_this.k9sConfigTableUI)

}

func (_this *K9SBrowser) Stop() {
	//if _this.stopDebounceCh != nil {
	//	close(_this.stopDebounceCh) // 停止防抖协程
	//}
}

// --- HELPER FUNCTIONS ---

func (_this *K9SBrowser) K9SShellIn(evt *tcell.EventKey) *tcell.EventKey {

	_this.Stop()
	defer _this.Start()

	_this.shellIn()
	return nil

}

func (_this *K9SBrowser) shellIn() {
	args := make([]string, 0, 15)
	args = append(args, "--kubeconfig", _this.selectConfigPath)
	slog.Info("k9s Shell args", "args", args)

	err := runK9sExec(_this.app,
		&shellOpts{
			clear:  true,
			banner: "",
			args:   args,
		},
	)
	if err != nil {
		_this.app.UI.Flash().Errf("Shell exec failed: %s", err)
	}
}

// listDefaultK8sConfigFiles 列出所有的k8s配置文件
func (_this *K9SBrowser) listDefaultK8sConfigFiles() ([]string, error) {
	// 获取当前用户家目录的.kube目录下的文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	kubeDir := homeDir + "/.kube"
	files, err := helper.ListDefaultK8sConfigFiles(kubeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read kube directory: %w", err)
	}

	return files, nil
}

func NewK9SBrowser(app *App) *K9SBrowser {
	var name = "K9S Browser"
	f := &K9SBrowser{
		BaseFlex: NewBaseFlex(name),
		app:      app,
	}
	f.SetIdentifier(ui.K9S_BROWSER_ID)
	return f
}
