package amqp

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/hysios/x/events/common"
	"github.com/hysios/x/events/driver"
)

type amqpDriver struct{}

// CreatePublisher implements driver.Driver.
func (a *amqpDriver) CreatePublisher() (common.Publisher, error) {
	return amqp.NewPublisher(amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: DefaultConfig.URL,
			Reconnect: &amqp.ReconnectConfig{
				BackoffInitialInterval:     DefaultConfig.BackoffInitialInterval,
				BackoffRandomizationFactor: DefaultConfig.BackoffRandomizationFactor,
				BackoffMultiplier:          DefaultConfig.BackoffMultiplier,
				BackoffMaxInterval:         DefaultConfig.BackoffMaxInterval,
			},
		},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				return topic
			},
			Type:    "topic",
			Durable: true,
		},
		Queue: amqp.QueueConfig{
			Durable: true,
		},
		Marshaler: DefaultConfig.Marshaler,
	}, watermill.NewStdLogger(false, false))
}

// CreateSubscriber implements driver.Driver.
func (a *amqpDriver) CreateSubscriber() (common.Subscriber, error) {
	return amqp.NewSubscriber(amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: DefaultConfig.URL,
			Reconnect: &amqp.ReconnectConfig{
				BackoffInitialInterval:     DefaultConfig.BackoffInitialInterval,
				BackoffRandomizationFactor: DefaultConfig.BackoffRandomizationFactor,
				BackoffMultiplier:          DefaultConfig.BackoffMultiplier,
				BackoffMaxInterval:         DefaultConfig.BackoffMaxInterval,
			},
		},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				return topic
			},
			Type:    "topic",
			Durable: true,
		},
		Queue: amqp.QueueConfig{
			Durable: true,
		},
		Marshaler: DefaultConfig.Marshaler,
	}, watermill.NewStdLogger(false, false))
}

type Config struct {
	URL                        string
	Marshaler                  amqp.Marshaler
	BackoffInitialInterval     time.Duration
	BackoffRandomizationFactor float64
	BackoffMultiplier          float64
	BackoffMaxInterval         time.Duration
}

var DefaultConfig = Config{
	URL:                        "amqp://guest:guest@localhost:5672",
	Marshaler:                  amqp.DefaultMarshaler{},
	BackoffInitialInterval:     500 * time.Millisecond,
	BackoffRandomizationFactor: 0.15,
	BackoffMultiplier:          1.5,
	BackoffMaxInterval:         5 * time.Second,
}

var _ driver.Driver = &amqpDriver{}

func init() {
	driver.SetDriver(&amqpDriver{})
}
