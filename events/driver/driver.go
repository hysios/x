package driver

import (
	"errors"

	"github.com/hysios/x/events/common"
)

var (
	publisher func() (common.Publisher, error)
	subscribe func() (common.Subscriber, error)
)

// CreatePublisher creates a new Publisher.
func CreatePublisher() (common.Publisher, error) {
	if publisher == nil {
		return nil, errors.New("publisher is not set")
	}

	return publisher()
}

// CreateSubscriber creates a new Subscriber.
func CreateSubscriber() (common.Subscriber, error) {
	if subscribe == nil {
		return nil, errors.New("subscriber is not set")
	}

	return subscribe()
}

func SetPublisher(fn func() (common.Publisher, error)) {
	publisher = fn
}

func SetSubscriber(fn func() (common.Subscriber, error)) {
	subscribe = fn
}

func SetDriver(d Driver) {
	SetPublisher(d.CreatePublisher)
	SetSubscriber(d.CreateSubscriber)
}

type Driver interface {
	CreateSubscriber() (common.Subscriber, error)
	CreatePublisher() (common.Publisher, error)
	// Close() error
}
