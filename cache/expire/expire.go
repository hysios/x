package expire

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/hysios/x/cache"
)

type expireCache[Key comparable, Value any] struct {
	cache *expirable.LRU[Key, Value]
}

func New[Key comparable, Value any](size int, m expirable.EvictCallback[Key, Value], ttl time.Duration) cache.Cache[Key, Value] {
	l := expirable.NewLRU[Key, Value](size, m, ttl)

	return &expireCache[Key, Value]{
		cache: l,
	}
}

// Load returns the value stored in the cache for a key, or nil if no value is present.
func (l *expireCache[Key, Value]) Load(key Key, opts ...cache.LoadOpt) (val Value, ok bool) {
	return l.cache.Get(key)
}

// Update sets the value for a key.
func (l *expireCache[Key, Value]) Update(key Key, val Value, opts ...cache.UpdateOpt) {
	l.cache.Add(key, val)
}

// Clear removes the value for a key.
func (l *expireCache[Key, Value]) Clear(key Key) {
	l.cache.Remove(key)
}

func (l *expireCache[Key, Value]) Keys() []Key {
	return l.cache.Keys()
}

func (l *expireCache[Key, Value]) Range(fn func(k Key, v Value) bool) {
	keys := l.cache.Keys()
	for _, key := range keys {
		val, ok := l.cache.Get(key)
		if !ok {
			continue
		}

		if !fn(key, val) {
			break
		}
	}
}
