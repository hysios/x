package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type natsDriver struct {
	conn     *nats.Conn
	js       jetstream.JetStream
	consumes map[jetstream.ConsumeContext]bool
}

var Default *natsDriver

// CreateStream
func (n *natsDriver) CreateStream(name string, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	return n.js.CreateStream(context.TODO(), cfg)
}

func (n *natsDriver) Publish(topic string, payload []byte) error {
	_, err := n.js.Publish(context.TODO(), topic, payload)
	return err
}

func Publish(topic string, payload []byte) error {
	return Default.Publish(topic, payload)
}
