package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strings"
	"sync"

	"github.com/rivo/tview"
	"lxz/internal/config"
)

// Logo represents a LXZ logo.
type Logo struct {
	*tview.Flex

	logo, status *tview.TextView
	styles       *config.Styles
	mx           sync.Mutex
}

// NewLogo returns a new logo.
func NewLogo(styles *config.Styles) *Logo {
	l := Logo{
		Flex:   tview.NewFlex(),
		logo:   logo(),
		status: status(),
		styles: styles,
	}
	l.SetDirection(tview.FlexRow)
	l.AddItem(l.logo, 0, 7, false)
	l.AddItem(l.status, 0, 1, false)
	//	l.refreshLogo(styles.Body().LogoColor)
	l.refreshLogo("red")
	l.SetBackgroundColor(tcell.ColorBlue)
	styles.AddListener(&l)

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
	l.refreshLogo(l.styles.Body().LogoColor)
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
	for i, s := range LogoSmall {
		_, _ = fmt.Fprintf(l.logo, "[%s::b]%s", c, s)
		if i+1 < len(LogoSmall) {
			_, _ = fmt.Fprintf(l.logo, "\n")
		}
	}
}

func logo() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignLeft)
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
