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

	type OrderCanceled struct {
		OrderID string
	}

	// Setup the PubMux
	mux := new(eventually.PubMux)
	ctx := eventually.ContextWithPub(context.Background(), mux)

	// React to the events
	mux.React(func(event OrderCompleted) {
		fmt.Printf("OrderCompleted{OrderID:%s}\n", event.OrderID)
	})

	mux.React(func(event OrderCanceled) {
		fmt.Printf("OrderCanceled{OrderID:%s}\n", event.OrderID)
	})

	// Publish event through Reactor inside the context
	eventually.Publish(ctx, OrderCompleted{OrderID: "123"})
	eventually.Publish(ctx, OrderCanceled{OrderID: "456"})

	// Output:
	// OrderCompleted{OrderID:123}
	// OrderCanceled{OrderID:456}
}
