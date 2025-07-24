/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 16:46
 */

package model

import (
	"context"
	"github.com/rivo/tview"
	"k8s.io/apimachinery/pkg/labels"
	"lxz/internal/view/cmd"
)

// Primitive represents a UI primitive.
type Primitive interface {
	tview.Primitive

	// Name returns the view name.
	Name() string
}

// Igniter represents a runnable view.
type Igniter interface {
	// Start starts a component.
	Init(ctx context.Context) error

	// Start starts a component.
	Start()

	// Stop terminates a component.
	Stop()
}

// Hinter represent a menu mnemonic provider.
type Hinter interface {
	// Hints returns a collection of menu hints.
	Hints() MenuHints

	// ExtraHints returns additional hints.
	ExtraHints() map[string]string
}

// Commander tracks prompt status.
type Commander interface {
	// InCmdMode checks if prompt is active.
	InCmdMode() bool
}

// Viewer represents a resource viewer.
type Viewer interface {
	// SetCommand sets the current command.
	SetCommand(*cmd.Interpreter)
}

type Filterer interface {
	SetFilter(string)
	SetLabelSelector(labels.Selector)
}

// Component represents a ui component.
type Component interface {
	Primitive
	Igniter
	Hinter
	Commander
	Filterer
	Viewer
}
