package amqp

import (
	"context"
	"time"

	"github.com/hysios/x/mq"
	"github.com/mitchellh/mapstructure"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL            string        `mapstructure:"url"`
	ExchangeName   string        `mapstructure:"exchange_name"`
	QueueName      string        `mapstructure:"queue_name"`
	PublishTimeout time.Duration `mapstructure:"publish_timeout"`
	Durable        bool          `mapstructure:"durable"`
}

type amqpDriver struct {
	conn           *amqp.Connection
	ExchangeName   string
	QueueName      string
	PublishTimeout time.Duration
	Durable        bool
}

var (
	DefaultConfig = Config{
		URL:            "amqp://guest:guest@localhost:5672/",
		ExchangeName:   "events",
		PublishTimeout: 5 * time.Second,
		Durable:        true,
	}
)

var Default *amqpDriver

func Open(cfg Config) (*amqpDriver, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	var driver = &amqpDriver{
		conn:           conn,
		ExchangeName:   cfg.ExchangeName,
		PublishTimeout: cfg.PublishTimeout,
		QueueName:      cfg.QueueName,
		Durable:        cfg.Durable,
	}

	if Default == nil {
		Default = driver
		return Default, nil
	}

	return driver, nil
}

func (a *amqpDriver) Close() error {
	return a.conn.Close()
}

func (a *amqpDriver) Channel() (*amqp.Channel, error) {
	return a.conn.Channel()
}

// CreateExchange
func (a *amqpDriver) CreateExchange(name string) error {
	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		name,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// createExchange
func (a *amqpDriver) createExchange(ch *amqp.Channel, name string) error {
	err := ch.ExchangeDeclare(
		name,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// queueBind
func (a *amqpDriver) queueBind(ch *amqp.Channel, queue, topic string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queue,
		a.Durable, // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		topic,  // routing key
		a.ExchangeName,
		false,
		nil)
	if err != nil {
		return amqp.Queue{}, err
	}

	return q, nil
}

// Publish
func (a *amqpDriver) Publish(topic string, payload []byte, opts ...mq.PubOpt) error {
	var opt = &mq.PubOption{}
	for _, o := range opts {
		o(opt)
	}

	anch, err := a.conn.Channel()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.PublishTimeout)
	defer cancel()

	if err := a.createExchange(anch, a.ExchangeName); err != nil {
		return err
	}

	var deliveryMode uint8
	if a.Durable {
		deliveryMode = amqp.Persistent
	}

	return anch.PublishWithContext(ctx,
		a.ExchangeName, // exchange
		topic,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         payload,
			DeliveryMode: deliveryMode,
		})
}

// Subscribe
func (a *amqpDriver) Subscribe(topic string, opts ...mq.SubOpt) (<-chan mq.Message, error) {
	var opt = &mq.SubOption{}
	for _, o := range opts {
		o(opt)
	}

	anch, err := a.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := a.createExchange(anch, a.ExchangeName); err != nil {
		return nil, err
	}

	q, err := a.queueBind(anch, opt.Queue, topic)
	if err != nil {
		return nil, err
	}

	msgs, err := anch.Consume(
		q.Name,      // queue
		opt.Consume, // consumer
		false,       // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	if err != nil {
		return nil, err
	}

	var ch = make(chan mq.Message)
	go func() {
		for d := range msgs {
			ch <- &message{d}
		}
	}()

	return ch, nil
}

type message struct {
	amqp.Delivery
}

func (m *message) Payload() []byte {
	return m.Body
}

func (m *message) Ack() bool {
	return m.Delivery.Ack(false) == nil
}

func init() {
	mq.Register("amqp", func(c mq.Config) mq.Driver {
		var cfg Config

		// 配置 mapstructure 解码器以支持 time.Duration
		decoderConfig := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
			Result:     &cfg,
		}
		decoder, err := mapstructure.NewDecoder(decoderConfig)
		if err != nil {
			panic(err)
		}

		if err := decoder.Decode(c); err != nil {
			panic(err)
		}

		// 手动合并配置，确保用户配置能够覆盖默认配置，包括零值
		dst := mergeConfigs(DefaultConfig, cfg, c)

		driver, err := Open(dst)
		if err != nil {
			panic(err)
		}

		return driver
	})
}

// mergeConfigs 手动合并配置，确保用户设置的零值能够覆盖默认值
func mergeConfigs(defaultCfg, userCfg Config, rawConfig mq.Config) Config {
	result := defaultCfg

	// 如果 rawConfig 为 nil，直接返回默认配置
	if rawConfig == nil {
		return result
	}

	// 检查原始配置中是否明确设置了字段，如果设置了就使用用户配置的值（包括零值）
	if _, exists := rawConfig["url"]; exists {
		result.URL = userCfg.URL
	}
	if _, exists := rawConfig["exchange_name"]; exists {
		result.ExchangeName = userCfg.ExchangeName
	}
	if _, exists := rawConfig["queue_name"]; exists {
		result.QueueName = userCfg.QueueName
	}
	if _, exists := rawConfig["publish_timeout"]; exists {
		result.PublishTimeout = userCfg.PublishTimeout
	}
	if _, exists := rawConfig["durable"]; exists {
		result.Durable = userCfg.Durable
	}

	return result
}
