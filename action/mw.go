package action

import (
	"context"
	"errors"
	"fmt"
)

const MaxRetry = 3

type RestFunc func(ctx context.Context, action any)

func ReplayMiddleware[Action any](rest RestFunc) Middleware[Action] {
	return func(next HandlerFunc[Action]) HandlerFunc[Action] {
		return func(ctx context.Context, action Action) (val any, err error) {

			for i := 0; i < MaxRetry; i++ {
				val, err = next(ctx, action)
				if errors.Is(err, ErrRetry) {
					rest(ctx, action)
				} else {
					return
				}
			}

			return nil, fmt.Errorf("retry too many times: %w", err)
		}
	}
}
