package store

type Store[Key, Value any] interface {
	Load(key Key) (val Value, ok bool)
	Store(key Key, val Value)
	LoadOrStore(key Key, value Value) (actual Value, loaded bool)
	Delete(key Key)
	LoadAndDelete(key Key) (value Value, loaded bool)
	Range(f func(key Key, value Value) bool)
}
