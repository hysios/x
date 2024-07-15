package mq

import (
	"fmt"

	"github.com/hysios/x/providers"
)

type Publisher interface {
	Publish(topic string, payload []byte, opts ...PubOpt) error
}

type Subscriber interface {
	Subscribe(topic string, opts ...SubOpt) (<-chan Message, error)
}

type Driver interface {
	Publisher
	Subscriber
}

type Message interface {
	Payload() []byte
	Ack() bool
}

type Config map[string]interface{}

var provider providers.Provider[string, providers.Ctor[Config, Driver]]

func Register(name string, ctor providers.Ctor[Config, Driver]) {
	provider.Register(name, ctor)
}

func Open(name string, cfg Config) (driver Driver, err error) {
	ctor, ok := provider.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("mq: driver %s not found", name)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("mq: driver %s panic: %v", name, r)
		}
	}()

	driver = ctor(cfg)
	return
}
