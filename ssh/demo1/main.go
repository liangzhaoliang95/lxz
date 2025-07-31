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

type HostItem struct {
	Host     string
	HostName string
	User     string
	Port     string
	Jump     string
	Source   string
	FilePath string
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
	var currentHostName string
	inHostBlock := false

	base := filepath.Base(path)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过注释 和 空行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "Host ") {
			// 新 Host 块开始
			if len(currentHosts) > 0 && currentHostName != "" {
				for _, h := range currentHosts {
					results = append(results, HostItem{
						Host:     h,
						HostName: currentHostName,
						Source:   base,
						FilePath: path,
					})
				}
			}
			currentHosts = strings.Fields(line)[1:]
			currentHostName = ""
			inHostBlock = true
		} else if inHostBlock {
			if strings.HasPrefix(strings.ToLower(line), "hostname") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					currentHostName = parts[1]
				}
			}
		}
	}

	// 处理文件末尾的最后一个 Host block
	if len(currentHosts) > 0 && currentHostName != "" {
		for _, h := range currentHosts {
			results = append(results, HostItem{
				Host:     h,
				HostName: currentHostName,
				Source:   base,
				FilePath: path,
			})
		}
	}

	return results, nil
}

// map[来源文件名][]HostItem
type HostMap map[string][]HostItem

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
	hostList := tview.NewList()
	hostList.ShowSecondaryText(false).SetBorder(true).SetTitle("🔐 Hosts")

	// 加载所有配置主机
	hostMap, sourceOrder, err := loadAllHostsGrouped()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading hosts: %v\n", err)
		os.Exit(1)
	}

	// 按来源文件名排序展示
	// 左侧配置源列表项
	for _, source := range sourceOrder {
		items := hostMap[source]
		srcName := source

		envList.AddItem(srcName, "", 0, func() {
			hostList.Clear()
			for _, h := range items {
				host := h.Host
				label := fmt.Sprintf("%s ➜ %s", host, h.HostName)
				hostList.AddItem(label, "", 0, func() {
					app.Suspend(func() {
						fmt.Printf("🔗 Connecting to %s...\n", host)
						cmd := exec.Command("ssh", host)
						cmd.Stdin = os.Stdin
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Run()
					})
				})
			}
			app.SetFocus(hostList)
		})
	}
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			if app.GetFocus() == envList {
				app.SetFocus(hostList)
			} else {
				app.SetFocus(envList)
			}
			return nil
		}
		return event
	})

	// layout: 左列表 + 右列表
	flex := tview.NewFlex().
		AddItem(envList, 30, 1, true). // 左侧环境宽 30 字符
		AddItem(hostList, 0, 2, false) // 右侧主机拉伸填满

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
