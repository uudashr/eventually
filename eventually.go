package eventually

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Event is a type of event. It can be struct of anything represent the event.
// The struct name describe the event name, the struct fields describe the event data.
//
// Example:
//
//	type OrderCompleted struct {
//		OrderID string
//	}
type Event any

// EventHandle is a function type that handle the event.
// The function form is:
//
//		func(event Event)
//
//	 where the Event is the event type (struct) that will be handled.
//
// Example:
//
//	eventually.HandleEvent(func(event OrderCompleted) {
//		// handle the event
//	})
type EventHandler any

type Eventually struct {
	events   []Event
	handlers map[reflect.Type][]EventHandler
}

func (e *Eventually) React(fn EventHandler) {
	ensureValidHandler(fn)
	if e.handlers == nil {
		e.handlers = make(map[reflect.Type][]EventHandler)
	}

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)
	e.handlers[fnTypeIn] = append(e.handlers[fnTypeIn], fn)
}

func (e *Eventually) RemoveHandler(fn EventHandler) {
	ensureValidHandler(fn)

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)
	handlers := e.handlers[fnTypeIn]
	if handlers == nil {
		return
	}

	for i, handler := range handlers {
		if reflect.ValueOf(handler).Pointer() == reflect.ValueOf(fn).Pointer() {
			e.handlers[fnTypeIn] = append(handlers[:i], handlers[i+1:]...)
			return
		}
	}
}

// RaiseEvent add the event to the list of events.
func (e *Eventually) RaiseEvent(event Event) error {
	err := validateEvent(event)
	if err != nil {
		return err
	}

	e.events = append(e.events, event)
	e.emit(event)
	return nil
}

func (e *Eventually) emit(event Event) {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() != reflect.Struct {
		panic(fmt.Errorf("eventually: event should be a struct (got: %v)", eventType.Kind()))
	}

	handlers := e.handlers[eventType]
	for _, handler := range handlers {
		invokeHandler(handler, event)
	}
}

// HandleEvent handle the raised events with the given handler.
func (e *Eventually) HandleEvent(fn EventHandler) error {
	err := validateHandler(fn)
	if err != nil {
		return err
	}

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)
	for _, event := range e.events {
		if reflect.TypeOf(event) == fnTypeIn {
			invokeHandler(fn, event)
		}
	}
	return nil
}

func (e *Eventually) Events() []Event {
	return e.events
}

func validateEvent(event Event) error {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() != reflect.Struct {
		return fmt.Errorf("eventually: event should be a struct (got: %v)", eventType.Kind())
	}

	return nil
}

func invokeHandler(handler EventHandler, event Event) {
	fnValue := reflect.ValueOf(handler)
	eventValue := reflect.ValueOf(event)
	fnValue.Call([]reflect.Value{eventValue})
}

func ensureValidHandler(fn EventHandler) {
	err := validateHandler(fn)
	if err != nil {
		panic(err)
	}
}

func validateHandler(fn EventHandler) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("eventually: fn EventHandler is not a function (got: %v)", fnType.Kind())
	}

	if fnType.NumIn() != 1 {
		return fmt.Errorf("eventually: fn EventHandler should have 1 input parameter (got: %d)", fnType.NumIn())
	}

	if fnType.NumOut() != 0 {
		return fmt.Errorf("eventually: fn EventHandler should have 0 output parameter (got: %d)", fnType.NumOut())
	}

	if fnType.In(0).Kind() != reflect.Struct {
		return fmt.Errorf("eventually: fn EventHandler input parameter should be a struct (got: %v)", fnType.In(0).Kind())
	}

	return nil
}

type contextKey string

const (
	contextKeyEventually contextKey = "eventually"
)

func WithContext(ctx context.Context, evtl *Eventually) context.Context {
	return context.WithValue(ctx, contextKeyEventually, evtl)
}

// RaiseEvent raise the event to the eventually in the context.
func RaiseEvent(ctx context.Context, event Event) error {
	evtl, ok := ctx.Value(contextKeyEventually).(*Eventually)
	if !ok {
		return errors.New("eventually: context does not have eventually")
	}

	evtl.RaiseEvent(event)
	return nil
}

func React(ctx context.Context, fn EventHandler) error {
	evtl, ok := ctx.Value(contextKeyEventually).(*Eventually)
	if !ok {
		return errors.New("eventually: context does not have eventually")
	}

	evtl.React(fn)
	return nil
}

// HandleEvent handle the raised events with the given handler in the eventually in the context.
func HandleEvent(ctx context.Context, fn EventHandler) error {
	evtl, ok := ctx.Value(contextKeyEventually).(*Eventually)
	if !ok {
		return errors.New("eventually: context does not have eventually")
	}

	evtl.HandleEvent(fn)
	return nil
}
