package eventually

import (
	"context"
)

// Event represents specific event. It must be a struct type.
//
// Example:
//
//	type OrderCompleted struct {
//	  OrderID string
//	}
type Event any

type contextKey struct{}

// Publisher publishes [Event].
type Publisher interface {
	// Publish event.
	Publish(Event)
}

// ContextWithPub wraps the pub and makes it available to the returned [context.Context].
func ContextWithPub(ctx context.Context, pub Publisher) context.Context {
	return context.WithValue(ctx, contextKey{}, pub)
}

// Publish event using [Publisher] available inside the ctx.
//
// If the [Publisher] is not available in the ctx, it will do nothing.
func Publish(ctx context.Context, event Event) {
	pub, ok := ctx.Value(contextKey{}).(Publisher)
	if !ok {
		return
	}

	pub.Publish(event)
}
