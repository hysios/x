package cache

var (
	Namespace = "$$cache"
)

type Cache[Key, Value any] interface {
	Load(key Key, opts ...LoadOpt) (val Value, ok bool)
	Update(key Key, val Value, opts ...UpdateOpt)
	Clear(key Key)
}

type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, v interface{}) error
}

var (
	DefaultEncoder Encoder = &jsonEncoder{}
	DefaultDecoder Decoder = &jsonEncoder{}
)

func SetNamespace(ns string) {
	Namespace = ns
}

func With[Key, Value any](cache Cache[Key, Value], set func(key Key) (Value, error)) func(key Key) (Value, error) {
	return func(key Key) (Value, error) {
		var z Value
		if val, ok := cache.Load(key); ok {
			return val, nil
		}

		val, err := set(key)
		if err != nil {
			return z, err
		}
		cache.Update(key, val)
		return val, nil
	}
}
