package ui

import (
	"fmt"
	"strings"

	"github.com/liangzhaoliang95/lxz/internal/config"
	ver "github.com/liangzhaoliang95/lxz/internal/version"

	"github.com/liangzhaoliang95/tview"
)

// LogoSmall LXZ small log.
var LogoSmall2 = []string{
	` _     __   __ ______`,
	`| |    \ \ / /|___  /`,
	`| |     \ V /    / / `,
	`| |     /   \   / /  `,
	`| |____/ /^\ \./ /___`,
	`\_____/\/   \/\_____/`,
}

// LogoBig LXZ big logo for splash page.
var LogoBig = []string{
	` _     __   __ ______`,
	`| |    \ \ / /|___  /`,
	`| |     \ V /    / / `,
	`| |     /   \   / /  `,
	`| |____/ /^\ \./ /___`,
	`\_____/\/   \/\_____/`,
}

var LogoSmall = []string{
	`▗▖   ▗▖  ▗▖▗▄▄▄▄▖`,
	`▐▌    ▝▚▞▘    ▗▞▘`,
	`▐▌     ▐▌   ▗▞▘  `,
	`▐▙▄▄▖▗▞▘▝▚▖▐▙▄▄▄▖`,
}

// Splash represents a splash screen.
type Splash struct {
	*tview.Flex
}

// NewSplash instantiates a new splash screen with product and company info.
func NewSplash(styles *config.Styles, version string) *Splash {
	s := Splash{Flex: tview.NewFlex()}
	s.SetBackgroundColor(styles.BgColor())

	logo := tview.NewTextView()
	logo.SetDynamicColors(true)
	logo.SetTextAlign(tview.AlignCenter)
	s.layoutLogo(logo, styles)

	vers := tview.NewTextView()
	vers.SetDynamicColors(true)
	vers.SetTextAlign(tview.AlignCenter)
	s.layoutRev(vers, version, styles)

	s.SetDirection(tview.FlexRow)
	s.AddItem(logo, 10, 1, false)
	s.AddItem(vers, 1, 1, false)

	return &s
}

func (*Splash) layoutLogo(t *tview.TextView, styles *config.Styles) {
	logo := strings.Join(LogoBig, fmt.Sprintf("\n[%s::b]", styles.Body().LogoColor))
	_, _ = fmt.Fprintf(t, "%s[%s::b]%s\n",
		strings.Repeat("\n", 2),
		styles.Body().LogoColor,
		logo)
}

func (*Splash) layoutRev(t *tview.TextView, rev string, styles *config.Styles) {
	// 获取版本信息
	v := ver.GetVersion()

	// 检查是否有新版本
	updateInfo, err := ver.CheckForUpdates()
	hasUpdate := err == nil && updateInfo != nil

	// 显示版本信息
	if hasUpdate {
		// 如果有新版本，显示黄色提示和升级箭头
		_, _ = fmt.Fprintf(t, "[%s::b]Version [yellow::b]%s [yellow::b]→ %s ↑",
			styles.Body().FgColor, v.Version, updateInfo.LatestVersion)
	} else {
		// 如果没有新版本，只显示当前版本
		_, _ = fmt.Fprintf(t, "[%s::b]Version [yellow::b]%s",
			styles.Body().FgColor, v.Version)
	}
}
