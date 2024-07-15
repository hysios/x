package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/stan.go"
)

type Config struct {
	URL              string
	SubscribersCount int
	QueueGroupPrefix string
	Stream           bool
	Subjects         []string
	CloseTimeout     time.Duration
	AckWaitTimeout   time.Duration
	StanOptions      []stan.Option
	NatsOptions      []nats.Option
}

var (
	DefaultURL    = "nats://localhost:4222"
	DefaultConfig = Config{
		URL:              DefaultURL,
		QueueGroupPrefix: "events",
		CloseTimeout:     time.Minute,
		AckWaitTimeout:   time.Second * 30,
	}
)

func Open(cfg Config) (*natsDriver, error) {
	conn, err := nats.Connect(cfg.URL, cfg.NatsOptions...)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to NATS: %w", err)
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("cannot create jetstream: %w", err)
	}

	if Default == nil {
		Default = &natsDriver{
			conn: conn,
			js:   js,
		}
		return Default, nil
	}

	return &natsDriver{
		conn: conn,
		js:   js,
	}, nil
}
