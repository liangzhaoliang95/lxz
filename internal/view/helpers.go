package view

import (
	"context"
	"errors"

	"github.com/liangzhaoliang95/lxz/internal"
)

func extractApp(ctx context.Context) (*App, error) {
	app, ok := ctx.Value(internal.KeyApp).(*App)
	if !ok {
		return nil, errors.New("no application found in context")
	}

	return app, nil
}
