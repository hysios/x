package action

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/fatih/structs"
	"github.com/hysios/x/maps"
	"github.com/tj/assert"
)

type Event interface {
	isEvent()
}

type UserCreated struct {
	UserId uint
	Rest   float64
}

func (*UserCreated) isEvent() {}

type UserDeleted struct {
}

func (*UserDeleted) isEvent() {}

// BalanceChanged
type BalanceChanged struct {
	UserId uint
	Money  float64
}

func (*BalanceChanged) isEvent() {}

// LogRecord log record
type LogRecord struct {
	UserId  uint
	Message string
	From    string
	Action  string
}

func (*LogRecord) isEvent() {}

// Login login
type Login struct {
	UserId   uint
	Password string
}

func (*Login) isEvent() {}

type User struct {
	Id         uint
	Balance    float64
	LoginCount int
}

func buildActions() *SimpleSet[Event] {
	var actions = &SimpleSet[Event]{}

	var users maps.Map[uint, *User]

	actions.Handle(&UserCreated{}, func(ctx context.Context, action Event) (any, error) {
		if action.(*UserCreated).Rest < 0 {
			return nil, fmt.Errorf("rest must be greater than 0")
		}

		_, _ = actions.Invoke(ctx, &LogRecord{
			UserId:  action.(*UserCreated).UserId,
			Message: "user created",
			From:    "user",
			Action:  "create",
		})

		return actions.Invoke(ctx, &BalanceChanged{
			UserId: action.(*UserCreated).UserId,
			Money:  100.0,
		})
	})

	actions.Handle(&BalanceChanged{}, func(ctx context.Context, action Event) (any, error) {
		_, _ = actions.Invoke(ctx, &LogRecord{
			UserId:  action.(*BalanceChanged).UserId,
			Message: "user balance changed",
			From:    "user",
			Action:  "balance.changed",
		})

		var user = &User{
			Id:      action.(*BalanceChanged).UserId,
			Balance: action.(*BalanceChanged).Money,
		}

		users.Store(user.Id, user)

		return user, nil
	})

	logMiddleware := func(next HandlerFunc[Event]) HandlerFunc[Event] {
		return func(ctx context.Context, action Event) (any, error) {
			log.Printf("mw: event is %+v", action)
			return next(ctx, action)
		}
	}

	actions.Use(logMiddleware)
	actions.Use(ReplayMiddleware[Event](func(ctx context.Context, action any) {
		log.Printf("replay: event is %+v", action)

		s := structs.New(action)
		userId := s.Field("UserId").Value().(uint)
		users.Store(userId, &User{Id: userId})
	}))

	actions.Handle(&LogRecord{}, func(ctx context.Context, action Event) (any, error) {
		return nil, nil
	})

	// Handle login
	actions.Handle(&Login{}, func(ctx context.Context, action Event) (any, error) {
		user, ok := users.Load(action.(*Login).UserId)
		if !ok {
			return nil, fmt.Errorf("not found user id %d: failed %w", action.(*Login).UserId, ErrRetry)
		}

		user.LoginCount++
		users.Store(user.Id, user)

		return user, nil
	})

	return actions
}

func TestActions(t *testing.T) {
	var (
		actions = buildActions()
		ctx     = context.Background()
	)

	v, err := actions.Invoke(ctx, &UserCreated{
		UserId: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, &User{
		Id:      1,
		Balance: 100.0,
	}, v)

	v, err = actions.Invoke(ctx, &UserCreated{
		UserId: 2,
		Rest:   -1.0,
	})

	assert.Error(t, err)
	assert.Nil(t, v)

	v, err = actions.Invoke(ctx, &Login{
		UserId: 1,
	})

	v, err = actions.Invoke(ctx, &Login{
		UserId: 2,
	})

	assert.NoError(t, err)
	t.Logf("user is %+v", v)
}
