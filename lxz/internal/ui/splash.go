// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"fmt"
	"lxz/internal/config"
	"strings"

	"github.com/rivo/tview"
)

// LogoSmall K9s small log.
var LogoSmall2 = []string{
	` _     __   __ ______`,
	`| |    \ \ / /|___  /`,
	`| |     \ V /    / / `,
	`| |     /   \   / /  `,
	`| |____/ /^\ \./ /___`,
	`\_____/\/   \/\_____/`,
}

// LogoBig K9s big logo for splash page.
var LogoBig = []string{
	`  _    __  __ _____     _  _____ ____  `,
	` | |   \ \/ /|__  /    | |/ ( _ ) ___| `,
	` | |    \  /   / /_____| ' // _ \___ \ `,
	` | |___ /  \  / /|_____| . \ (_) |__) |`,
	` |_____/_/\_\/____|    |_|\_\___/____/ `,
	`                                       `,
}

var LogoSmall = []string{
	` _     __   __ ______         _____                _	   `,
	`| |    \ \ / /|___  /        |_   _|              | |	   `,
	`| |     \ V /    / /  ______   | |    ___    ___  | | ___ `,
	`| |     /   \   / /  |______|  | |   / _ \  / _ \ | |/ __|`,
	`| |____/ /^\ \./ /___          | |  | (_) || (_) || |\__ \`,
	`\_____/\/   \/\_____/          \_/   \___/  \___/ |_||___/`,
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
