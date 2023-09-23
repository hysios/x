package providers

import "github.com/hysios/x/maps"

type Ctor[C, A any] func(C) A

type Provider[T, A any] struct {
	store maps.Map[T, A]
}

// Register a provider
func (p *Provider[T, A]) Register(t T, ctor A) {
	p.store.Store(t, ctor)
}

// Lookup a provider
func (p *Provider[T, A]) Lookup(t T) (ctor A, ok bool) {
	if ctor, ok := p.store.Load(t); ok {
		return ctor, true
	}

	return
}

// Range
func (p *Provider[T, A]) Range(fn func(t T, ctor A) bool) {
	p.store.Range(func(key T, value A) bool {
		return fn(key, value)
	})
}

type SliceProvider[T, A any] struct {
	store maps.Map[T, []A]
}

// Register a provider
func (p *SliceProvider[T, A]) Register(t T, ctor A) {
	if ctors, ok := p.store.Load(t); !ok {
		p.store.Store(t, []A{ctor})
	} else {
		p.store.Store(t, append(ctors, ctor))
	}
}

// Lookup a provider
func (p *SliceProvider[T, A]) Lookup(t T) (ctors []A, ok bool) {
	if ctors, ok := p.store.Load(t); ok {
		return ctors, true
	}

	return
}

// Range
func (p *SliceProvider[T, A]) Range(fn func(t T, ctor A) bool) {
	p.store.Range(func(key T, value []A) bool {
		for _, ctor := range value {
			if !fn(key, ctor) {
				return false
			}
		}

		return true
	})
}
