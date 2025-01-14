package eventually

import (
	"context"
	"errors"
)

// Event is a type of event. It can be struct of anything represent the event.
// The struct name describe the event name, the struct fields describe the event data.
//
// Example:
//
//	type OrderCompleted struct {
//	  OrderID string
//	}
type Event any

type contextKey struct{}

type Publisher interface {
	Publish(Event) error
}

// ContextWithPub wrap the Publisher into the context.
func ContextWithPub(ctx context.Context, pub Publisher) context.Context {
	return context.WithValue(ctx, contextKey{}, pub)
}

// Publish the event to the Publisher in the context.
func Publish(ctx context.Context, event Event) error {
	pub, ok := ctx.Value(contextKey{}).(Publisher)
	if !ok {
		return errors.New("eventually: context does not have Publisher")
	}

	return pub.Publish(event)
}

// OptPublish publish the event if the Publisher is available in the context.
func OptPublish(ctx context.Context, event Event) {
	_ = Publish(ctx, event)
}

// MustPublish publish the event. It will panic upon error found.
func MustPublish(ctx context.Context, event Event) {
	if err := Publish(ctx, event); err != nil {
		panic(err)
	}
}
