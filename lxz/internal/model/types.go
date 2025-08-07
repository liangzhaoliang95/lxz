/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 16:46
 */

package model

import (
	"context"
	"github.com/rivo/tview"
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

type Identifier interface {
	SetIdentifier(string)
	GetIdentifier() string
}

// Component represents a ui component.
type Component interface {
	Primitive
	Igniter
	Hinter
	Identifier
}
