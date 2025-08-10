package ui

import (
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/config"
	"strings"

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
	_, _ = fmt.Fprintf(t, "[%s::b]Revision [red::b]%s", styles.Body().FgColor, rev)
}
