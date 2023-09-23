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
