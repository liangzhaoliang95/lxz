package ui

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/liangzhaoliang95/lxz/internal/config"
	ver "github.com/liangzhaoliang95/lxz/internal/version"
	"github.com/liangzhaoliang95/tview"
)

// Logo represents a LXZ logo.
type Logo struct {
	*tview.Flex

	logo, status, version *tview.TextView
	styles                *config.Styles
	mx                    sync.Mutex
}

// NewLogo returns a new logo.
func NewLogo(styles *config.Styles) *Logo {
	l := Logo{
		Flex:    tview.NewFlex(),
		logo:    logo(),
		status:  status(),
		version: version(),
		styles:  styles,
	}
	l.SetDirection(tview.FlexRow)
	l.AddItem(l.logo, 0, 7, false)
	l.AddItem(l.version, 0, 1, false) // 添加version组件
	//l.AddItem(l.status, 0, 1, false)
	//	l.refreshLogo(styles.Body().LogoColor)
	l.refreshLogo("dodgerblue")
	l.refreshVersion("dodgerblue") // 刷新version组件
	l.SetBackgroundColor(tcell.ColorBlue)
	styles.AddListener(&l)

	// 启动时自动检查版本更新
	go l.autoCheckVersion()

	return &l
}

// Logo returns the logo viewer.
func (l *Logo) Logo() *tview.TextView {
	return l.logo
}

// Status returns the status viewer.
func (l *Logo) Status() *tview.TextView {
	return l.status
}

// StylesChanged notifies the skin changed.
func (l *Logo) StylesChanged(s *config.Styles) {
	l.styles = s
	l.SetBackgroundColor(l.styles.BgColor())
	l.status.SetBackgroundColor(l.styles.BgColor())
	l.logo.SetBackgroundColor(l.styles.BgColor())
	l.version.SetBackgroundColor(l.styles.BgColor())
	l.refreshLogo(l.styles.Body().LogoColor)
	l.refreshVersion(l.styles.Body().LogoColor)
}

// IsBenchmarking checks if benchmarking is active or not.
func (l *Logo) IsBenchmarking() bool {
	txt := l.Status().GetText(true)
	return strings.Contains(txt, "Bench")
}

// Reset clears out the logo view and resets colors.
func (l *Logo) Reset() {
	l.status.Clear()
	l.StylesChanged(l.styles)
}

// Err displays a log error state.
func (l *Logo) Err(msg string) {
	l.update(msg, l.styles.Body().LogoColorError)
}

// Warn displays a log warning state.
func (l *Logo) Warn(msg string) {
	l.update(msg, l.styles.Body().LogoColorWarn)
}

// Info displays a log info state.
func (l *Logo) Info(msg string) {
	l.update(msg, l.styles.Body().LogoColorInfo)
}

func (l *Logo) update(msg string, c config.Color) {
	l.refreshStatus(msg, c)
	l.refreshLogo(c)
}

func (l *Logo) refreshStatus(msg string, c config.Color) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.status.SetBackgroundColor(c.Color())
	l.status.SetText(
		fmt.Sprintf("[%s::b]%s", l.styles.Body().LogoColorMsg, msg),
	)
}

func (l *Logo) refreshLogo(c config.Color) {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.logo.Clear()

	// 打印logo
	for i, s := range LogoSmall {
		_, _ = fmt.Fprintf(l.logo, "[%s::b]%s", c, s)
		if i+1 < len(LogoSmall) {
			_, _ = fmt.Fprintf(l.logo, "\n")
		}
	}
}

// refreshVersion 刷新版本信息显示
func (l *Logo) refreshVersion(c config.Color) {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.version.Clear()

	// 获取版本信息
	v := ver.GetVersion()

	// 检查是否有新版本
	updateInfo, err := ver.CheckForUpdates()
	hasUpdate := err == nil && updateInfo != nil

	// 打印版本信息
	versionText := fmt.Sprintf("%s", v.Version)
	if hasUpdate {
		// 如果有新版本，显示黄色提示和升级箭头
		_, _ = fmt.Fprintf(l.version, "[%s::b]%s [yellow::b]→ %s ↑",
			c, versionText, updateInfo.LatestVersion)
	} else {
		// 如果没有新版本，只显示当前版本
		_, _ = fmt.Fprintf(l.version, "[%s::b]%s", c, versionText)
	}
}

// autoCheckVersion 自动检查版本更新
func (l *Logo) autoCheckVersion() {
	// 延迟2秒后开始检查，避免阻塞启动
	time.Sleep(2 * time.Second)

	// 检查版本更新
	updateInfo, err := ver.CheckForUpdates()
	if err != nil {
		// 记录错误但不显示给用户
		slog.Error("自动检查版本更新失败", "error", err)
		return
	}

	// 如果有更新，刷新version组件显示
	if updateInfo != nil {
		// 使用主线程刷新UI
		l.refreshVersion(l.styles.Body().LogoColor)
	}
}

func logo() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignCenter)
	v.SetDynamicColors(true)
	//v.SetBorder(true)

	return v
}

func status() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignCenter)
	v.SetDynamicColors(true)
	return v
}

func version() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignCenter)
	v.SetDynamicColors(true)
	return v
}
