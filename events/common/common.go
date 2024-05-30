package common

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Message = message.Message

type Publisher interface {
	Publish(topic string, messages ...*Message) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
	Close() error
}
