package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

var (
	activeBorderColor   = tcell.ColorRed  // ✅ 获得焦点
	inactiveBorderColor = tcell.ColorGray // ❌ 非焦点
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

func main() {
	app := tview.NewApplication()

	// 左侧配置来源列表
	envList := tview.NewList()
	envList.ShowSecondaryText(false).SetBorder(true).SetTitle("🌍 Environments")

	// 右侧主机列表
	//hostList := tview.NewList()
	//hostList.ShowSecondaryText(false).SetBorder(true).SetTitle("🔐 Hosts")

	hostTable := tview.NewTable().
		SetSelectable(true, false)
	hostTable.
		SetBorders(false).
		SetBorder(true).
		SetTitle("🔐 Hosts")

	// ✅ 设置 table 的回车行为，ssh 连接
	hostTable.SetSelectedFunc(func(row, _ int) {
		if row == 0 {
			return // 表头
		}
		cell := hostTable.GetCell(row, 0)
		if cell == nil {
			return
		}
		host := cell.Text
		app.Suspend(func() {
			fmt.Printf("🔗 Connecting to %s...\n", host)
			cmd := exec.Command("ssh", host)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		})
	})

	// 加载所有配置主机
	hostMap, sourceOrder, err := loadAllHostsGrouped()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading hosts: %v\n", err)
		os.Exit(1)
	}

	// 按来源文件名排序展示
	// 左侧配置源列表项
	// ✅ 替换 envList.AddItem 中的回调闭包部分：
	for _, source := range sourceOrder {
		//items := hostMap[source]
		srcName := source

		// ⚠️ 闭包要绑定当前 source / items
		envList.AddItem(srcName, "", 0, nil)
	}

	// 自动更新右侧 hostTable 内容
	envList.SetChangedFunc(func(index int, name string, _ string, _ rune) {
		items := hostMap[name]
		hostTable.Clear()

		// 表头
		headers := []string{"Name", "HostName", "User", "Port"}
		for i, h := range headers {
			hostTable.SetCell(0, i,
				tview.NewTableCell("[::b]"+h).
					SetTextColor(tcell.ColorYellow).
					SetSelectable(false),
			)
		}

		// 数据
		for row, item := range items {
			hostTable.SetCell(row+1, 0, tview.NewTableCell(item.Host))
			hostTable.SetCell(row+1, 1, tview.NewTableCell(item.HostName))
			hostTable.SetCell(row+1, 2, tview.NewTableCell(item.User))
			hostTable.SetCell(row+1, 3, tview.NewTableCell(item.Port))
		}

		// 自动选中首行主机
		if len(items) > 0 {
			hostTable.Select(1, 0)
		}
	})

	// 回车时只做焦点切换
	envList.SetSelectedFunc(func(index int, name string, _ string, _ rune) {
		app.SetFocus(hostTable)
		envList.SetBorderColor(inactiveBorderColor)
		hostTable.SetBorderColor(activeBorderColor)
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if app.GetFocus() == envList {
				app.SetFocus(hostTable)
				envList.SetBorderColor(inactiveBorderColor)
				hostTable.SetBorderColor(activeBorderColor)
			} else {
				app.SetFocus(envList)
				envList.SetBorderColor(activeBorderColor)
				hostTable.SetBorderColor(inactiveBorderColor)
			}
			return nil
		}
		return event
	})

	// layout: 左列表 + 右列表
	flex := tview.NewFlex().
		AddItem(envList, 30, 1, true).  // 左侧环境宽 30 字符
		AddItem(hostTable, 0, 2, false) // 右侧主机拉伸填满

	// ✅ 默认选中第一个环境并调用它
	if envList.GetItemCount() > 0 {
		envList.SetCurrentItem(1)
		envList.SetCurrentItem(0)
		selectFunc := envList.GetItemSelectedFunc(0)
		if selectFunc != nil {
			selectFunc()
		}

	}

	// ✅ 设置默认边框颜色 + 焦点 + 强制刷新
	envList.SetBorderColor(activeBorderColor)
	hostTable.SetBorderColor(inactiveBorderColor)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
