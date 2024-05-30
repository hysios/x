package events

import (
	"context"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/hysios/x/events/common"
	"github.com/hysios/x/events/driver"
)

type Router struct {
	*message.Router
}

var DefaultRoute = &Router{}
var once sync.Once

func NewRouter() *Router {
	var logger = watermill.NewStdLogger(true, true)

	router, _ := message.NewRouter(message.RouterConfig{}, logger)
	router.AddPlugin(plugin.SignalsHandler)

	// Router level middleware are executed for every message sent to the router
	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,

		// The handler function is retried if it returns an error.
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,

		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
	)

	return &Router{
		Router: router,
	}
}

var l sync.Mutex

func Start(ctx context.Context, ready func(router *Router)) *Router {
	l.Lock()
	defer l.Unlock()

	if DefaultRoute == nil {
		DefaultRoute = NewRouter()
		go func() {
			if ready != nil {
				ready(DefaultRoute)
			}
			DefaultRoute.Run(ctx)
		}()
	}

	return DefaultRoute
}

func (r *Router) Close() error {
	return r.Router.Close()
}

func (r *Router) Publish(topic string, messages ...*common.Message) error {
	publish, err := driver.CreatePublisher()
	if err != nil {
		return err
	}
	return publish.Publish(topic, messages...)
}

// Subscribe subscribes to the given topic.
func (r *Router) Subscribe(ctx context.Context, topic string) (<-chan *common.Message, error) {
	subscribe, err := driver.CreateSubscriber()
	if err != nil {
		return nil, err
	}
	return subscribe.Subscribe(ctx, topic)
}

// Run
func (r *Router) Run(ctx context.Context) error {
	return r.Router.Run(ctx)
}

func Publish(topic string, messages ...*common.Message) error {
	return DefaultRoute.Publish(topic, messages...)
}

func Subscribe(ctx context.Context, topic string) (<-chan *common.Message, error) {
	return DefaultRoute.Subscribe(ctx, topic)
}
