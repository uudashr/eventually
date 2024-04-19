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
//	  OrderID string
//	}
type Event any

// EventHandle is a function type that handle the event.
//
// The form of function is:
//
//	func(event Event)
//
// where the [Event] is the event type (struct) that will be handled.
//
// Example:
//
//	eventually.HandleEvent(func(event OrderCompleted) {
//	  // handle the event
//	})
type EventHandler any

// Eventually is responsible for managing the lifecycle of events and their
// handlers.
type Eventually struct {
	events   []Event
	handlers map[reflect.Type][]EventHandler
}

// React reacting to the event with the given fn as the handler. It will only
// react to the new raised event. Past events will not be handled by the fn
// handler.
func (e *Eventually) React(fn EventHandler) {
	ensureValidHandler(fn)
	if e.handlers == nil {
		e.handlers = make(map[reflect.Type][]EventHandler)
	}

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)
	e.handlers[fnTypeIn] = append(e.handlers[fnTypeIn], fn)
}

// RemoveHandler remove the handler from the list of handlers which previously
// registed via [Eventually.React]. Subsequent raised events will not longer be
// handled by the fn handler.
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

// HandleEvent handle the raised events with the given handler. The fn handler
// will not be registerd, it will not be called for the new raised event.
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

// Events return the list of events that have been raised.
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

type contextKey struct{}

// WithContext return a new context with the evtl.
func WithContext(ctx context.Context, evtl *Eventually) context.Context {
	return context.WithValue(ctx, contextKey{}, evtl)
}

// RaiseEvent raise the event to the eventually in the ctx.
func RaiseEvent(ctx context.Context, event Event) error {
	evtl, ok := ctx.Value(contextKey{}).(*Eventually)
	if !ok {
		return errors.New("eventually: context does not have eventually")
	}

	evtl.RaiseEvent(event)
	return nil
}

// React reacting to the event with the given fn as the handler of [Eventually]
// in the ctx.
func React(ctx context.Context, fn EventHandler) error {
	evtl, ok := ctx.Value(contextKey{}).(*Eventually)
	if !ok {
		return errors.New("eventually: context does not have eventually")
	}

	evtl.React(fn)
	return nil
}

// HandleEvent handle the raised events with the given handler in the eventually
// in the ctx.
func HandleEvent(ctx context.Context, fn EventHandler) error {
	evtl, ok := ctx.Value(contextKey{}).(*Eventually)
	if !ok {
		return errors.New("eventually: context does not have eventually")
	}

	evtl.HandleEvent(fn)
	return nil
}
