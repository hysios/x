package action

import (
	"context"
	"strings"
	"testing"

	"github.com/tj/assert"
)

func TestActionSet_Handle(t *testing.T) {
	// Initialize the ActionSet
	as := &ActionSet[string, string]{
		Broadcast: true,
	}

	// Add a handler function
	as.Handle("upper", func(ctx context.Context, action string) (interface{}, error) {
		// Test the handler logic here
		return strings.ToUpper(action), nil
	})

	as.Handle("lower", func(ctx context.Context, action string) (interface{}, error) {
		// Test the handler logic here
		return strings.ToLower(action), nil
	})

	as.Handle("lower", func(ctx context.Context, action string) (interface{}, error) {
		// Test the handler logic here
		t.Logf("lower: %s", action)
		return strings.ToLower(action), nil
	})

	// Invoke the action
	val, err := as.Invoke(context.Background(), "upper", "action")
	assert.NoError(t, err)
	assert.Equal(t, "ACTION", val)

	val, err = as.Invoke(context.Background(), "lower", "ACTION")
	assert.NoError(t, err)
	assert.Equal(t, "action", val)
}

func TestActionSet_Invoke(t *testing.T) {
	// Initialize the ActionSet
	as := &ActionSet[string, string]{
		handlers: make(map[string][]Handler[string]),
		mws:      []ActionMiddleware[string, string]{},
	}

	// Add a handler function
	as.Handle("key", func(ctx context.Context, action string) (interface{}, error) {
		// Test the handler logic here
		return nil, nil
	})

	// Invoke the action
	_, err := as.Invoke(context.Background(), "key", "action")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestActionSet_Use(t *testing.T) {
	// Initialize the ActionSet
	as := &ActionSet[string, string]{
		handlers: make(map[string][]Handler[string]),
		mws:      []ActionMiddleware[string, string]{},
	}

	// Add a middleware function
	as.Use(func(next ActionHandlerFunc[string, string]) ActionHandlerFunc[string, string] {
		return func(ctx context.Context, key string, action string) (interface{}, error) {
			// Test the middleware logic here
			t.Logf("key: %s, action: %s", key, action)
			return next(ctx, key, action)
		}
	})

	// Add a handler function
	as.Handle("key", func(ctx context.Context, action string) (interface{}, error) {
		// Test the handler logic here
		return nil, nil
	})

	// Invoke the action
	_, err := as.Invoke(context.Background(), "key", "action")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
