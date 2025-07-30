/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 16:31
 */

package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/model"
	"lxz/internal/slogs"
)

type Pages struct {
	*tview.Pages
	*model.Stack
}

// NewPages return a new view.
func NewPages() *Pages {
	p := Pages{
		Pages: tview.NewPages(),
		Stack: model.NewStack(),
	}
	p.AddListener(&p)

	return &p
}

// Add adds a new page.
func (p *Pages) add(c model.Component) {
	p.AddPage(componentID(c), c, true, true)
}

// Show displays a given page.
func (p *Pages) Show(c model.Component) {
	p.SwitchToPage(componentID(c))
}

// AddAndShow adds a new page and bring it to front.
func (p *Pages) addAndShow(c model.Component) {
	slog.Info("LXZ Pages addAndShow 页面切换 💥", slogs.Component, c.Name())
	p.add(c)
	p.Show(c)
}

// StackPushed notifies a new component was pushed.
func (p *Pages) StackPushed(c model.Component) {
	slog.Info("LXZ Pages 入栈", slogs.Component, c.Name())
	p.addAndShow(c)
}

// Delete removes a page.
func (p *Pages) delete(c model.Component) {
	p.RemovePage(componentID(c))
}

// StackPopped notifies a component was removed.
func (p *Pages) StackPopped(o, _ model.Component) {
	slog.Info("LXZ Pages 出栈", slogs.Component, o.Name())
	p.delete(o)
}

// StackTop notifies a new component is at the top of the stack.
func (p *Pages) StackTop(top model.Component) {
	if top == nil {
		return
	}
	p.Show(top)
}

// IsTopDialog checks if front page is a dialog.
func (p *Pages) IsTopDialog() bool {
	_, pa := p.GetFrontPage()
	switch pa.(type) {
	case *tview.Modal, *ModalList:
		return true
	default:
		return false
	}
}

func componentID(c model.Component) string {
	if c.Name() == "" {
		slog.Error("Component has no name", slogs.Component, fmt.Sprintf("%T", c))
	}
	return fmt.Sprintf("%s-%p", c.Name(), c)
}
