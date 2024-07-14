package nats

import "github.com/nats-io/nats.go"

type natsDriver struct {
	conn *nats.Conn
}

var Default = &natsDriver{}

func (n *natsDriver) Publish(topic string, payload []byte) error {
	return n.conn.Publish(topic, payload)
}

func Publish(topic string, payload []byte) error {
	return Default.Publish(topic, payload)
}
