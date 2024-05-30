package nats

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/hysios/x/events/common"
	"github.com/hysios/x/events/driver"
	nc "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type Config struct {
	URL              string
	SubscribersCount int
	QueueGroupPrefix string
	CloseTimeout     time.Duration
	AckWaitTimeout   time.Duration
	StanOptions      []stan.Option
	Marshaler        nats.Marshaler
	Unmarshaler      nats.Unmarshaler
}

type PublishConfig struct {
	URL string
}

var (
	DefaultURL    = "nats://localhost:4222"
	DefaultConfig = Config{
		URL:              DefaultURL,
		QueueGroupPrefix: "events",
		SubscribersCount: 4,
		CloseTimeout:     time.Minute,
		AckWaitTimeout:   time.Second * 30,
		Marshaler:        nats.JSONMarshaler{},
		Unmarshaler:      nats.JSONMarshaler{},
	}
	DefaultPublicConfig = PublishConfig{
		URL: DefaultURL,
	}
)

type natsDriver struct {
}

// CreatePublisher implements driver.Driver.
func (n *natsDriver) CreatePublisher() (common.Publisher, error) {
	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}
	return nats.NewPublisher(
		nats.PublisherConfig{
			URL:       DefaultURL,
			JetStream: jsConfig,
			Marshaler: DefaultConfig.Marshaler,
			// NatsOptions: DefaultPublicConfig.NatsOptions,
		},
		watermill.NewStdLogger(false, false),
	)
}

// CreateSubscriber implements driver.Driver.
func (n *natsDriver) CreateSubscriber() (common.Subscriber, error) {
	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}

	return nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:              DefaultURL,
			QueueGroupPrefix: DefaultConfig.QueueGroupPrefix,
			SubscribersCount: DefaultConfig.SubscribersCount,
			CloseTimeout:     DefaultConfig.CloseTimeout,
			AckWaitTimeout:   DefaultConfig.AckWaitTimeout,
			Unmarshaler:      DefaultConfig.Unmarshaler,
			JetStream:        jsConfig,
		},
		watermill.NewStdLogger(false, false),
	)
}

var _ driver.Driver = &natsDriver{}

func init() {
	driver.SetDriver(&natsDriver{})
}
