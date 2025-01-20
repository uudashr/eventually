package eventually_test

import (
	"context"
	"fmt"

	"github.com/uudashr/eventually/v2"
)

func ExampleRecorder() {
	// Event definition
	type OrderCompleted struct {
		OrderID string
	}

	// Setup Recorder
	rec := eventually.NewRecorder()
	ctx := eventually.ContextWithPub(context.Background(), rec)

	// Publish event through Recorder inside the context
	eventually.Publish(ctx, OrderCompleted{OrderID: "123"})
	eventually.Publish(ctx, OrderCompleted{OrderID: "456"})

	// Print the recorded events
	fmt.Printf("Record %d events\n", len(rec.Events))
	for i, event := range rec.Events {
		fmt.Printf("%d. %+v\n", i+1, event)
	}

	// Output:
	// Record 2 events
	// 1. {OrderID:123}
	// 2. {OrderID:456}

}
