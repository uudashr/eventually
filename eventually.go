package eventually

import (
	"context"
)

// Event represent specific event. It must be a struct type.
//
// Example:
//
//	type OrderCompleted struct {
//	  OrderID string
//	}
type Event any

type contextKey struct{}

// Publisher interface for event publishing.
type Publisher interface {
	// Publish the event.
	Publish(Event)
}

// ContextWithPub wrap the Publisher into the context.
func ContextWithPub(ctx context.Context, pub Publisher) context.Context {
	return context.WithValue(ctx, contextKey{}, pub)
}

// Publish the event to the Publisher in the context.
func Publish(ctx context.Context, event Event) {
	pub, ok := ctx.Value(contextKey{}).(Publisher)
	if !ok {
		return
	}

	pub.Publish(event)
}
