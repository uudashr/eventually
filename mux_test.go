package eventually_test

import (
	"context"
	"fmt"

	"github.com/uudashr/eventually/v2"
)

func ExamplePubMux() {
	// Event definition
	type OrderCompleted struct {
		OrderID string
	}

	// Setup the PubMux
	mux := eventually.NewPubMux()
	ctx := eventually.ContextWithPub(context.Background(), mux)

	// React to the events
	mux.React(func(event OrderCompleted) {
		fmt.Printf("OrderCompleted{OrderID:%s} \n", event.OrderID)
	})

	// Publish event through Reactor inside the context
	eventually.Publish(ctx, OrderCompleted{OrderID: "123"})

	// Output:
	// OrderCompleted{OrderID:123}
}
