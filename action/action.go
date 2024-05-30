package action

import (
	"context"
)

type ActionSet[K comparable, Action any] struct {
	Broadcast bool
	handlers  map[K][]Handler[Action]
	mws       []ActionMiddleware[K, Action]
}

type (
	ActionMiddleware[K comparable, Action any]  func(ActionHandlerFunc[K, Action]) ActionHandlerFunc[K, Action]
	ActionHandlerFunc[K comparable, Action any] func(ctx context.Context, key K, action Action) (any, error)
)

// init
func (b *ActionSet[K, Action]) init() {
	if b.mws == nil {
		b.mws = make([]ActionMiddleware[K, Action], 0)
		// b.Use(b.mwReturnAction)
	}
}

// Use
func (b *ActionSet[K, Action]) Use(mw ActionMiddleware[K, Action]) {
	b.mws = append(b.mws, mw)
}

// buildMwChain
func (b *ActionSet[K, Action]) buildMwChain(fn ActionHandlerFunc[K, Action]) ActionHandlerFunc[K, Action] {
	for i := len(b.mws) - 1; i >= 0; i-- {
		fn = b.mws[i](fn)
	}

	return fn
}

// Handle handle action
func (b *ActionSet[K, Action]) Handle(key K, fn HandlerFunc[Action]) {
	if b.handlers == nil {
		b.handlers = make(map[K][]Handler[Action])
	}

	b.handlers[key] = append(b.handlers[key], &funcHandler[Action]{fn: fn})
}

// Invoke
func (b *ActionSet[K, Action]) Invoke(ctx context.Context, key K, act Action) (any, error) {
	b.init()

	if handlers, ok := b.handlers[key]; ok {
		for i, handler := range handlers {
			fn := b.buildMwChain(func(ctx context.Context, key K, action Action) (any, error) {
				return handler.Handle(ctx, action)
			})
			v, err := fn(ctx, key, act)
			if b.Broadcast && i < len(handlers)-1 {
				continue
			}
			if v != nil || err != nil {
				return v, err
			}
		}
	} else {
		fn := b.buildMwChain(func(ctx context.Context, key K, action Action) (any, error) {
			return nil, ErrActionNotFound
		})

		return fn(ctx, key, act)
	}

	return nil, nil
}
