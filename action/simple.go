package action

import (
	"context"
	"fmt"
	"reflect"
)

type SimpleSet[Action comparable] struct {
	handlers map[reflect.Type][]Handler[Action]
	mws      []Middleware[Action]
}

type Middleware[Action any] func(HandlerFunc[Action]) HandlerFunc[Action]

type Handler[Action any] interface {
	Handle(ctx context.Context, action Action) (any, error)
}

type HandlerFunc[Action any] func(ctx context.Context, action Action) (any, error)

// init
func (b *SimpleSet[Action]) init() {
	if b.mws == nil {
		b.mws = make([]Middleware[Action], 0)
		b.Use(b.mwReturnAction)
	}
}

// Use
func (b *SimpleSet[Action]) Use(mw Middleware[Action]) {
	b.mws = append(b.mws, mw)
}

// buildMwChain
func (b *SimpleSet[Action]) buildMwChain(fn HandlerFunc[Action]) HandlerFunc[Action] {
	for i := len(b.mws) - 1; i >= 0; i-- {
		fn = b.mws[i](fn)
	}

	return fn
}

// HandleI handle action
func (b *SimpleSet[Action]) HandleI(action Action, handler Handler[Action]) {
	// 识别 action 的接口的 action 的具体类型
	var t = reflect.TypeOf(action).Elem()

	if b.handlers == nil {
		b.handlers = make(map[reflect.Type][]Handler[Action])
	}

	b.handlers[t] = append(b.handlers[t], handler)
}

// Handle handle action
func (b *SimpleSet[Action]) Handle(action Action, fn HandlerFunc[Action]) {
	b.HandleI(action, &funcHandler[Action]{fn: fn})
}

// Invoke invoke action
func (b *SimpleSet[Action]) Invoke(ctx context.Context, action Action) (any, error) {
	b.init()

	var t = reflect.TypeOf(action).Elem()

	if handlers, ok := b.handlers[t]; ok {
		for _, handler := range handlers {
			fn := b.buildMwChain(handler.Handle)
			v, err := fn(ctx, action)
			if v != nil || err != nil {
				return v, err
			}
		}
	}

	return nil, fmt.Errorf("not found handler for action %T %+v", action, action)
}

func (b *SimpleSet[Action]) mwReturnAction(next HandlerFunc[Action]) HandlerFunc[Action] {
	return func(ctx context.Context, action Action) (any, error) {
		v, err := next(ctx, action)
		if err != nil {
			return nil, err
		}

		if v != nil {
			if act, ok := v.(Action); ok {
				return b.Invoke(ctx, act)
			}
		}
		return v, nil
	}
}

type funcHandler[Action any] struct {
	fn HandlerFunc[Action]
}

func (f *funcHandler[Action]) Handle(ctx context.Context, action Action) (any, error) {
	return f.fn(ctx, action)
}
