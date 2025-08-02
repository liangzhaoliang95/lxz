/**
 * @author  zhaoliang.liang
 * @date  2025/7/31 16:02
 */

package view

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/apimachinery/pkg/labels"
	"log/slog"
	"lxz/internal/ui"
	"lxz/internal/view/cmd"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type SshConnect struct {
	*ui.BaseFlex
	app       *App
	envMap    map[string][]HostItem
	envOrder  []string
	envList   *tview.List
	hostTable *tview.Table
}

func (_this *SshConnect) bindKeys() {
	_this.Actions().Bulk(ui.KeyMap{
		ui.KeyF:         ui.NewKeyAction("FullScreen", _this.ToggleFullScreenCmd, true),
		tcell.KeyEscape: ui.NewKeyAction("Quit FullScreen", _this.ToggleFullScreenCmd, false),
		tcell.KeyTAB:    ui.NewKeyAction("Focus Change", _this.TabFocusChange, true),
		tcell.KeyEnter:  ui.NewKeyAction("Connect", _this.EmptyKeyEvent, true),
		tcell.KeyLeft:   ui.NewKeyAction("Focus Change", _this.TabFocusChange, false),
		tcell.KeyRight:  ui.NewKeyAction("Focus Change", _this.TabFocusChange, false),
	})
}

func (_this *SshConnect) Init(ctx context.Context) error {
	// 用于初始化组件的边框、标题、快捷键等信息
	var err error
	_this.envMap, _this.envOrder, err = loadAllHostsGrouped()
	if err != nil {
		slog.Error(fmt.Sprintf("%s loadAllHostsGrouped failed err => %s", _this.Name(), err.Error()))
		return err
	}
	_this.bindKeys()
	_this.SetInputCapture(_this.Keyboard)

	// 初始化左侧配置来源列表
	_this.envList = tview.NewList()
	_this.envList.ShowSecondaryText(false).
		SetBorder(true).
		SetBorderAttributes(tcell.AttrNone).
		SetBorderPadding(0, 1, 1, 1).
		SetBorderColor(tcell.ColorDefault)
	//_this.envList.SetTitle(" 🌍 Environments ")

	// 初始化hostTable
	_this.hostTable = tview.NewTable().
		SetSelectable(true, false)

	_this.hostTable.
		SetBorders(false).
		SetBorder(true).
		SetBorderAttributes(tcell.AttrNone).
		SetBorderPadding(0, 1, 1, 1).
		SetBorderColor(tcell.ColorDefault)
	//_this.hostTable.SetTitle(" 🔐 Hosts ")
	return nil
}

func (_this *SshConnect) Start() {
	slog.Info("SshConnect Start")
	// 按来源文件名排序展示
	// 左侧配置源列表项
	// ✅ 替换 envList.AddItem 中的回调闭包部分：
	for _, source := range _this.envOrder {
		//items := hostMap[source]
		srcName := source
		// ⚠️ 闭包要绑定当前 source / items
		_this.envList.AddItem(srcName, "", 0, nil)
	}

	// 自动更新右侧 hostTable 内容
	_this.envList.SetChangedFunc(func(index int, name string, _ string, _ rune) {
		items := _this.envMap[name]
		_this.hostTable.Clear()

		// 表头
		headers := []string{"Name", "HostName", "User", "Port"}
		for i, h := range headers {
			_this.hostTable.SetCell(0, i,
				tview.NewTableCell("[::b]"+h).
					SetTextColor(tcell.ColorYellow).
					SetSelectable(false),
			)
		}

		// 数据
		for row, item := range items {
			_this.hostTable.SetCell(row+1, 0, tview.NewTableCell(item.Host))
			_this.hostTable.SetCell(row+1, 1, tview.NewTableCell(item.HostName))
			_this.hostTable.SetCell(row+1, 2, tview.NewTableCell(item.User))
			_this.hostTable.SetCell(row+1, 3, tview.NewTableCell(item.Port))
		}

		// 自动选中首行主机
		if len(items) > 0 {
			_this.hostTable.Select(1, 0)
		}
	})

	// 回车时只做焦点切换
	_this.envList.SetSelectedFunc(func(index int, name string, _ string, _ rune) {
		_this.app.UI.SetFocus(_this.hostTable)
		_this.envList.SetBorderColor(inactiveBorderColor)
		_this.hostTable.SetBorderColor(activeBorderColor)
	})

	_this.AddItem(_this.envList, 30, 1, true)   // 左侧环境宽 30 字符
	_this.AddItem(_this.hostTable, 0, 2, false) // 右侧主机拉伸填满

	// ✅ 默认选中第一个环境并调用它
	if _this.envList.GetItemCount() > 0 {
		_this.envList.SetCurrentItem(1)
		_this.envList.SetCurrentItem(0)
		selectFunc := _this.envList.GetItemSelectedFunc(0)
		if selectFunc != nil {
			selectFunc()
		}

	}

	// ✅ 设置默认边框颜色 + 焦点 + 强制刷新
	_this.envList.SetBorderColor(activeBorderColor)
	_this.hostTable.SetBorderColor(inactiveBorderColor)

	// ✅ 设置 table 的回车行为，ssh 连接
	_this.hostTable.SetSelectedFunc(func(row, _ int) {
		if row == 0 {
			return // 表头
		}
		cell := _this.hostTable.GetCell(row, 0)
		if cell == nil {
			return
		}
		host := cell.Text
		_this.app.UI.Suspend(func() {
			slog.Info("🔗 Connecting to ", "host", host)
			cmd := exec.Command("ssh", host)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		})
	})
}

func (_this *SshConnect) Stop() {
	slog.Info("SshConnect Stop")
}

func (_this *SshConnect) ExtraHints() map[string]string {
	return nil
}

func (_this *SshConnect) InCmdMode() bool {
	return false
}

func (_this *SshConnect) SetFilter(s2 string) {
	//TODO implement me
	panic("implement me")
}

func (_this *SshConnect) SetLabelSelector(selector labels.Selector) {
	//TODO implement me
	panic("implement me")
}

func (_this *SshConnect) SetCommand(interpreter *cmd.Interpreter) {
	//TODO implement me
	panic("implement me")
}

func (_this *SshConnect) TabFocusChange(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTAB {

	} else if event.Key() == tcell.KeyLeft {
		if _this.app.UI.GetFocus() == _this.envList {
			return nil // 已经在 envList 上了
		}
	} else if event.Key() == tcell.KeyRight {
		if _this.app.UI.GetFocus() == _this.hostTable {
			return nil // 已经在 hostTable 上了
		}
	}
	if _this.app.UI.GetFocus() == _this.envList {
		_this.app.UI.SetFocus(_this.hostTable)
		_this.envList.SetBorderColor(inactiveBorderColor)
		_this.hostTable.SetBorderColor(activeBorderColor)
	} else {
		_this.app.UI.SetFocus(_this.envList)
		_this.envList.SetBorderColor(activeBorderColor)
		_this.hostTable.SetBorderColor(inactiveBorderColor)
	}
	return nil
}

// helpers

var (
	activeBorderColor   = tcell.ColorGreenYellow // ✅ 获得焦点
	inactiveBorderColor = tcell.ColorGray        // ❌ 非焦点
)

type HostItem struct {
	Host         string
	HostName     string
	User         string
	Port         string
	Jump         string
	Source       string
	FilePath     string
	IdentityFile string
	ProxyCommand string
}

// 解析配置文件，只返回含 HostName 字段的 Host
func parseConfigFileWithHostname(path string) ([]HostItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var results []HostItem
	var currentHosts []string
	var hostName, user, port, identityFile, proxyCommand string
	inHostBlock := false

	base := filepath.Base(path)

	saveCurrentHostBlock := func() {
		if len(currentHosts) > 0 && hostName != "" {
			for _, h := range currentHosts {
				results = append(results, HostItem{
					Host:         h,
					HostName:     hostName,
					User:         user,
					Port:         port,
					IdentityFile: identityFile,
					ProxyCommand: proxyCommand,
					Source:       base,
					FilePath:     path,
				})
			}
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "Host ") {
			// 保存上一个 Host 块
			saveCurrentHostBlock()

			// 开启新块
			currentHosts = strings.Fields(line)[1:]
			hostName, user, port, identityFile, proxyCommand = "", "", "", "", ""
			inHostBlock = true
		} else if inHostBlock {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			key := strings.ToLower(fields[0])
			value := strings.Join(fields[1:], " ")
			switch key {
			case "hostname":
				hostName = value
			case "user":
				user = value
			case "port":
				port = value
			case "identityfile":
				identityFile = value
			case "proxycommand":
				proxyCommand = value
			}
		}
	}

	// 保存最后一个块
	saveCurrentHostBlock()
	return results, nil
}

func loadAllHostsGrouped() (map[string][]HostItem, []string, error) {
	usr, _ := user.Current()
	mainConfig := filepath.Join(usr.HomeDir, ".ssh", "config")

	queue := []string{mainConfig}
	visited := make(map[string]bool)
	sourceOrder := []string{}
	hostMap := make(map[string][]HostItem)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		items, err := parseConfigFileWithHostname(current)
		if err == nil && len(items) > 0 {
			source := filepath.Base(current)
			hostMap[source] = append(hostMap[source], items...)
			sourceOrder = append(sourceOrder, source)
		}

		f, err := os.Open(current)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "Include ") {
				included := expandIncludes(line)
				queue = append(queue, included...)
			}
		}
	}

	return hostMap, sourceOrder, nil
}

func expandIncludes(line string) []string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}
	paths := fields[1:]
	var expanded []string
	for _, p := range paths {
		if strings.HasPrefix(p, "~") {
			if home, err := os.UserHomeDir(); err == nil {
				p = filepath.Join(home, p[1:])
			}
		}
		matches, err := filepath.Glob(p)
		if err == nil {
			expanded = append(expanded, matches...)
		}
	}
	return expanded
}

func NewSshConnect(app *App) *SshConnect {
	var name = "SSH Connect"
	tc := &SshConnect{
		BaseFlex: ui.NewBaseFlex(name),
		app:      app,
	}

	tc.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(fmt.Sprintf(" %s ", name))
	return tc
}
