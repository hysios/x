package nats

import "github.com/nats-io/nats.go"

func (n *natsDriver) Subscribe(topic string) (<-chan []byte, error) {
	var ch = make(chan []byte)
	_, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		ch <- msg.Data
	})
	if err != nil {
		return nil, err
	}

	return ch, nil
}

// Subscribe subscribes to a topic and returns a channel to receive messages.
func Subscribe(topic string) (<-chan []byte, error) {
	return Default.Subscribe(topic)
}
