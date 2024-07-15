package nats

import (
	"context"
	"log"
	"time"

	"github.com/hysios/x/mq"
	"github.com/nats-io/nats.go/jetstream"
)

type SubOption struct {
	Subjects []string
}

type SubOpt func(*SubOption)

type msgWarp struct {
	jetstream.Msg
}

func (m *msgWarp) Payload() []byte {
	return m.Msg.Data()
}

func (m *msgWarp) Ack() bool {
	return m.Msg.Ack() == nil
}

func (n *natsDriver) Subscribe(topic string, opts ...SubOpt) (<-chan mq.Message, error) {
	var (
		ch     = make(chan mq.Message)
		ctx, _ = context.WithTimeout(context.Background(), 10000*time.Hour)
		opt    = &SubOption{}
	)

	for _, o := range opts {
		o(opt)
	}

	// Create a stream
	s, err := n.js.Stream(context.TODO(), topic)
	if err != nil {
		s, _ = n.js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     topic,
			Subjects: opt.Subjects,
		})
	}

	// Create durable consumer
	c, err := s.CreateOrUpdateConsumer(context.TODO(), jetstream.ConsumerConfig{
		// Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Printf("create consumer error: %v", err)
		return nil, err
	}

	// Receive messages continuously in a callback
	cons, err := c.Consume(func(msg jetstream.Msg) {
		ch <- &msgWarp{msg}
	})
	if err != nil {
		log.Printf("consume error: %v", err)
		return nil, err
	}
	n.consumes[cons] = true

	return ch, nil
}

// Close
func (n *natsDriver) Close() error {
	n.conn.Close()

	for cons := range n.consumes {
		cons.Stop()
	}
	return nil
}

// Subscribe subscribes to a topic and returns a channel to receive messages.
func Subscribe(topic string, opts ...SubOpt) (<-chan mq.Message, error) {
	return Default.Subscribe(topic, opts...)
}
