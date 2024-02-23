package eventually_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/uudashr/eventually"
)

func ExampleEventually() {
	type OrderCompleted struct {
		OrderID string
	}

	var sendOrderCompleteNotification = func(orderID string) {
		fmt.Printf("Your order %q has been completed\n", orderID)
	}

	evtl := &eventually.Eventually{}

	// ReaiseEvent accept any kind of struct type
	evtl.RaiseEvent(OrderCompleted{OrderID: "order-123"})

	// HandleEvent accept any kind of func with any struct type as the argument
	evtl.HandleEvent(func(event OrderCompleted) {
		sendOrderCompleteNotification(event.OrderID)
	})

	// Output:
	// Your order "order-123" has been completed
}

func ExampleHandleEvent() {
	type OrderCompleted struct {
		OrderID string
	}

	var sendOrderCompleteNotification = func(orderID string) {
		fmt.Printf("Your order %q has been completed\n", orderID)
	}

	evtl := &eventually.Eventually{}

	ctx := eventually.WithContext(context.Background(), evtl)

	// react immediately
	eventually.React(ctx, func(event OrderCompleted) {
		fmt.Printf("Reacting to OrderCompleted event: %q\n", event.OrderID)
	})

	eventually.RaiseEvent(ctx, OrderCompleted{OrderID: "order-123"})

	// handle the raised events
	eventually.HandleEvent(ctx, func(event OrderCompleted) {
		sendOrderCompleteNotification(event.OrderID)
	})

	// Output:
	// Reacting to OrderCompleted event: "order-123"
	// Your order "order-123" has been completed
}

func TestReflect(t *testing.T) {
	type OrderCompleted struct {
		OrderID string
	}

	var handler = func(event OrderCompleted) error {
		return nil
	}

	ht := reflect.TypeOf(handler)
	if got, want := ht.Kind(), reflect.Func; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	if got, want := ht.NumIn(), 1; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	if got, want := ht.NumOut(), 1; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	if got, want := ht.In(0).Kind(), reflect.Struct; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	if got, want := ht.In(0).Name(), "OrderCompleted"; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	if got, want := ht.Out(0), reflect.TypeOf((*error)(nil)).Elem(); got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	if got, want := ht.Out(0).Name(), "error"; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	inVal := reflect.ValueOf(ht.In(0)).Interface()
	if got, want := inVal, reflect.TypeOf(OrderCompleted{}); got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
