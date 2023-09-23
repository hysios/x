package lru

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hysios/x/cache"
)

type lruCache[Key comparable, Value any] struct {
	cache *lru.Cache[Key, Value]
}

func New[Key comparable, Value any](size int) cache.Cache[Key, Value] {
	l, _ := lru.New[Key, Value](size)

	return &lruCache[Key, Value]{
		cache: l,
	}
}

func (l *lruCache[Key, Value]) Load(key Key, opts ...cache.LoadOpt) (val Value, ok bool) {
	return l.cache.Get(key)
}

func (l *lruCache[Key, Value]) Update(key Key, val Value, opts ...cache.UpdateOpt) {
	l.cache.Add(key, val)
}

func (l *lruCache[Key, Value]) Clear(key Key) {
	l.cache.Remove(key)
}
