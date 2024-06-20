package maps

import "sync"

type Map[T, V any] struct {
	sync.Map
}

func NewMap[T, V any]() *Map[T, V] {
	return &Map[T, V]{Map: sync.Map{}}
}

func (m *Map[T, V]) Load(key T) (val V, ok bool) {
	if v, ok := m.Map.Load(key); !ok {
		return val, false
	} else if val, ok = v.(V); ok {
		return val, ok
	} else {
		return val, false
	}
}

func (m *Map[T, V]) Store(key T, val V) {
	m.Map.Store(key, val)
}

func (m *Map[T, V]) LoadOrStore(key T, value V) (actual V, loaded bool) {
	v, loaded := m.Map.LoadOrStore(key, value)
	if !loaded {
		return value, false
	} else if actual, ok := v.(V); ok {
		return actual, true
	} else {
		return value, true
	}
}

func (m *Map[T, V]) Delete(key T) {
	m.Map.Delete(key)
}

func (m *Map[T, V]) LoadAndDelete(key T) (value V, loaded bool) {
	if v, loaded := m.Map.LoadAndDelete(key); !loaded {
		return value, false
	} else if value, ok := v.(V); ok {
		return value, true
	} else {
		return value, false
	}
}

func (m *Map[T, V]) Range(fn func(key T, value V) bool) {
	m.Map.Range(func(key, value any) bool {
		var (
			k, ok1 = key.(T)
			v, ok2 = value.(V)
		)

		if ok1 && ok2 {
			return fn(k, v)
		} else {
			return true
		}
	})

}
