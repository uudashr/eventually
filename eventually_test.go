package eventually_test

import (
	"reflect"
	"testing"
)

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
