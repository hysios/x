package mq

type Publisher interface {
	Publish(topic string, payload []byte) error
}

type Subscriber interface {
	Subscribe(topic string) (<-chan []byte, error)
}

type Driver interface {
	Publisher
	Subscriber
}
